package codegen

import (
	"cool-compiler/ast"
	"fmt"
	"sort"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// CaseResult represents the result of a case expression,
// including the value and the COOL types of each branch
type CaseResult struct {
	Value       value.Value // The actual LLVM value (PHI node result)
	BranchTypes []string    // The COOL types of each branch
}

// CodeGenerator handles the code generation process
type CodeGenerator struct {
	// Module represents the LLVM module being built
	Module *ir.Module

	// TypeMap maps COOL types to LLVM types
	TypeMap map[string]*types.StructType

	// VTables maps class names to their virtual method tables
	VTables map[string]*ir.Global

	// ClassHierarchy keeps track of inheritance relationships
	ClassHierarchy map[string]string

	// Current function being processed
	CurrentFunc *ir.Func

	CurrentBlock *ir.Block

	// Symbol table for variables in current scope
	Symbols map[string]value.Value

	// Standard library functions
	StdlibFuncs map[string]*ir.Func

	// Built-in classes
	BuiltInClasses []*ast.Class

	// Program classes from the source code
	ProgramClasses []*ast.Class

	// Map class names to their AST nodes for efficient lookup
	ClassNameToAST map[string]*ast.Class

	// AttributeIndices maps class names to a map of attribute names to their indices in the class struct
	AttributeIndices map[string]map[string]int

	// MethodIndices maps class names to a map of method names to their indices in the vtable
	MethodIndices map[string]map[string]int

	// Counters for generating unique block names in control structures
	IfCounter    int
	WhileCounter int
	CaseCounter  int

	// Common constants
	EmptyStringGlobal *ir.Global

	// CaseResults stores the results of case expressions for later retrieval
	CaseResults map[value.Value]*CaseResult
}

// Generate is the main entry point for code generation
func Generate(program *ast.Program) (*ir.Module, error) {
	// Create the code generator
	cg := NewCodeGenerator()

	// Define built-in classes
	cg.DefineBuiltInClasses()

	// Initialize standard library functions
	cg.initStdlib()

	// Store program classes for later reference
	cg.ProgramClasses = program.Classes

	// Build an efficient mapping of class names to AST nodes
	for _, class := range program.Classes {
		cg.ClassNameToAST[class.Name.Value] = class
	}

	// Add built-in classes to the map as well
	for _, builtInClass := range cg.BuiltInClasses {
		cg.ClassNameToAST[builtInClass.Name.Value] = builtInClass
	}

	// First, generate class structures
	cg.GenerateClassStructs(program)

	// Next, generate vtables
	cg.GenerateVTables(program)

	// Generate method implementations
	cg.GenerateMethods(program)

	// Ensure Main class exists before continuing
	if _, exists := cg.TypeMap["Main"]; !exists {
		panic("Program must have a Main class")
	}

	// Generate main function
	cg.GenerateMain(program)

	return cg.Module, nil
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	cg := &CodeGenerator{
		Module:           ir.NewModule(),
		TypeMap:          make(map[string]*types.StructType),
		VTables:          make(map[string]*ir.Global),
		ClassHierarchy:   make(map[string]string),
		Symbols:          make(map[string]value.Value),
		StdlibFuncs:      make(map[string]*ir.Func),
		BuiltInClasses:   []*ast.Class{},
		ProgramClasses:   []*ast.Class{},
		ClassNameToAST:   make(map[string]*ast.Class),
		AttributeIndices: make(map[string]map[string]int),
		MethodIndices:    make(map[string]map[string]int),
		IfCounter:        0,
		WhileCounter:     0,
		CaseCounter:      0,
	}

	emptyStr := constant.NewCharArrayFromString("\x00")
	cg.EmptyStringGlobal = cg.Module.NewGlobalDef(".str.empty", emptyStr)
	cg.EmptyStringGlobal.Immutable = true

	return cg
}

// initStdlib initializes standard library functions
func (cg *CodeGenerator) initStdlib() {
	// Memory management functions

	// malloc for object allocation (returns i8*)
	cg.StdlibFuncs["malloc"] = cg.Module.NewFunc(
		"malloc",
		types.NewPointer(types.I8),
		ir.NewParam("size", types.I64),
	)

	// free for manual memory deallocation (though COOL uses garbage collection)
	cg.StdlibFuncs["free"] = cg.Module.NewFunc(
		"free",
		types.Void,
		ir.NewParam("ptr", types.NewPointer(types.I8)),
	)

	// exit for terminating the program
	cg.StdlibFuncs["exit"] = cg.Module.NewFunc(
		"exit",
		types.Void,
		ir.NewParam("status", types.I32),
	)

	// IO functions

	// Declare C standard library functions
	printfFunc := cg.Module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true
	cg.StdlibFuncs["printf"] = printfFunc

	// scanf function for input
	scanfFunc := cg.Module.NewFunc(
		"scanf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	scanfFunc.Sig.Variadic = true
	cg.StdlibFuncs["scanf"] = scanfFunc

	// Standard C string functions
	cg.StdlibFuncs["strlen"] = cg.Module.NewFunc(
		"strlen",
		types.I32,
		ir.NewParam("str", types.NewPointer(types.I8)),
	)

	cg.StdlibFuncs["strcpy"] = cg.Module.NewFunc(
		"strcpy",
		types.NewPointer(types.I8),
		ir.NewParam("dest", types.NewPointer(types.I8)),
		ir.NewParam("src", types.NewPointer(types.I8)),
	)

	cg.StdlibFuncs["strcat"] = cg.Module.NewFunc(
		"strcat",
		types.NewPointer(types.I8),
		ir.NewParam("dest", types.NewPointer(types.I8)),
		ir.NewParam("src", types.NewPointer(types.I8)),
	)

	cg.StdlibFuncs["strncpy"] = cg.Module.NewFunc(
		"strncpy",
		types.NewPointer(types.I8),
		ir.NewParam("dest", types.NewPointer(types.I8)),
		ir.NewParam("src", types.NewPointer(types.I8)),
		ir.NewParam("n", types.I32),
	)
}

// GenerateClassStructs creates LLVM struct types for COOL classes
func (cg *CodeGenerator) GenerateClassStructs(program *ast.Program) {
	// First pass: Declare all class types
	for _, class := range program.Classes {
		cg.declareClassType(class)
	}

	// Second pass: Define class structures with fields
	for _, class := range program.Classes {
		cg.defineClassStruct(class)
	}
}

// declareClassType creates an LLVM struct type for a class
func (cg *CodeGenerator) declareClassType(class *ast.Class) {
	// Skip if this class is already declared
	if _, exists := cg.TypeMap[class.Name.Value]; exists {
		return
	}

	// Record inheritance relationship
	if class.Parent != nil {
		cg.ClassHierarchy[class.Name.Value] = class.Parent.Value
	} else {
		// If no parent specified, default to Object (unless this is the Object class itself)
		if class.Name.Value != "Object" {
			cg.ClassHierarchy[class.Name.Value] = "Object"
		}
	}

	// Create a placeholder struct type that will be defined later
	structType := types.NewStruct()
	cg.Module.NewTypeDef(class.Name.Value, structType)
	cg.TypeMap[class.Name.Value] = structType
}

// defineClassStruct defines the fields of a class struct
func (cg *CodeGenerator) defineClassStruct(class *ast.Class) {
	className := class.Name.Value
	classType := cg.TypeMap[className]

	// Create a map for this class's attributes if it doesn't exist
	if _, exists := cg.AttributeIndices[className]; !exists {
		cg.AttributeIndices[className] = make(map[string]int)
	}

	// Collect all fields for this class (including inherited fields)
	var fields []types.Type
	var parentFields []types.Type
	fieldIndex := 1 // Start at 1 because index 0 is vtable pointer

	// First field is always a pointer to the vtable
	fields = append(fields, types.NewPointer(types.I8)) // vtable pointer

	// If this class inherits from another class, include parent fields
	if parent, exists := cg.ClassHierarchy[className]; exists && parent != "" {
		parentType, exists := cg.TypeMap[parent]
		if !exists {
			panic(fmt.Sprintf("parent class %s not found for class %s", parent, className))
		}

		// Extract field types from parent class, skipping the vtable pointer
		// which will be at the beginning of our class
		for i := 1; i < len(parentType.Fields); i++ {
			parentFields = append(parentFields, parentType.Fields[i])
		}

		// Add parent fields to this class's fields
		fields = append(fields, parentFields...)

		// Copy the parent attribute indices to this class
		for attrName, attrIndex := range cg.AttributeIndices[parent] {
			cg.AttributeIndices[className][attrName] = attrIndex
		}

		fieldIndex += len(parentFields)
	}

	// Add class's own fields
	for _, feature := range class.Features {
		if attr, isAttr := feature.(*ast.Attribute); isAttr {
			// Determine LLVM type for the attribute
			var attrType types.Type

			switch attr.TypeDecl.Value {
			case "Int":
				attrType = types.I32
			case "Bool":
				attrType = types.I1
			case "String":
				attrType = types.NewPointer(types.I8)
			case "SELF_TYPE":
				// For self-referential types, use a pointer to the class itself
				attrType = types.NewPointer(classType)
			default:
				// For user-defined classes, use a pointer to the class
				referencedType, exists := cg.TypeMap[attr.TypeDecl.Value]
				if !exists {
					panic(fmt.Sprintf("undefined type %s in attribute %s of class %s",
						attr.TypeDecl.Value, attr.Name.Value, className))
				}
				attrType = types.NewPointer(referencedType)
			}

			// Add attribute to fields and track its index
			fields = append(fields, attrType)
			cg.AttributeIndices[className][attr.Name.Value] = fieldIndex
			fieldIndex++
		}
	}

	// Update the struct definition with all fields
	classType.Fields = fields
}

// GenerateVTables creates virtual method tables for all classes
func (cg *CodeGenerator) GenerateVTables(program *ast.Program) {
	// First pass: Create vtable types for each class
	for _, class := range program.Classes {
		cg.createVTableForClass(class, program)
	}
}

// createVTableForClass creates a vtable for a specific class
func (cg *CodeGenerator) createVTableForClass(class *ast.Class, program *ast.Program) {
	className := class.Name.Value

	// Skip if vtable already exists for this class
	if _, exists := cg.VTables[className]; exists {
		return
	}

	// Create method indices map for this class if it doesn't exist
	if _, exists := cg.MethodIndices[className]; !exists {
		cg.MethodIndices[className] = make(map[string]int)
	}

	// Initialize empty maps for our methods
	methods := make(map[string]*ir.Func)
	methodIndices := make(map[string]int)
	var methodNames []string

	// Process parent class first to inherit its methods
	if parent, exists := cg.ClassHierarchy[className]; exists && parent != "" {
		// Make sure parent vtable is created first
		parentClass := findClass(parent, program)
		if parentClass != nil {
			// First check if the parent vtable already exists to avoid recursive duplication
			if _, parentVTableExists := cg.VTables[parent]; !parentVTableExists {
				cg.createVTableForClass(parentClass, program)
			}
		}

		// Copy ALL methods from parent vtable and maintain the same indices
		// This ensures that overridden methods use the same slot in the vtable

		// We need to walk up the entire inheritance chain to gather all methods
		currentParent := parent
		for currentParent != "" {
			// Copy methods from the current parent
			for methodName, methodIndex := range cg.MethodIndices[currentParent] {
				// Only add if we haven't already added this method
				if _, exists := methods[methodName]; !exists {
					methodIndices[methodName] = methodIndex

					// Find the corresponding function in the parent
					parentFuncName := fmt.Sprintf("%s.%s", currentParent, methodName)
					for _, f := range cg.Module.Funcs {
						if f.Name() == parentFuncName {
							methods[methodName] = f
							if !containsString(methodNames, methodName) {
								methodNames = append(methodNames, methodName)
							}
							break
						}
					}
				}
			}

			// Move up to the next parent in the inheritance chain
			currentParent = cg.ClassHierarchy[currentParent]
		}
	}

	// Process methods of this class
	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod {
			methodName := method.Name.Value

			// Create function declaration
			var paramTypes []types.Type

			// First parameter is always 'self'
			selfType := types.NewPointer(cg.TypeMap[className])
			paramTypes = append(paramTypes, selfType)

			// Add formal parameter types
			for _, formal := range method.Formals {
				var paramType types.Type
				switch formal.TypeDecl.Value {
				case "Int":
					paramType = types.I32
				case "Bool":
					paramType = types.I1
				case "String":
					paramType = types.NewPointer(types.I8)
				case "SELF_TYPE":
					paramType = selfType
				default:
					// For class types, use a pointer to the class
					referencedType, exists := cg.TypeMap[formal.TypeDecl.Value]
					if !exists {
						panic(fmt.Sprintf("undefined type %s in method %s of class %s",
							formal.TypeDecl.Value, method.Name.Value, className))
					}
					paramType = types.NewPointer(referencedType)
				}
				paramTypes = append(paramTypes, paramType)
			}

			// Determine return type
			var returnType types.Type
			switch method.TypeDecl.Value {
			case "Int":
				returnType = types.I32
			case "Bool":
				returnType = types.I1
			case "String":
				returnType = types.NewPointer(types.I8)
			case "SELF_TYPE":
				returnType = selfType
			default:
				// For class types, use a pointer to the class
				referencedType, exists := cg.TypeMap[method.TypeDecl.Value]
				if !exists {
					panic(fmt.Sprintf("undefined return type %s in method %s of class %s",
						method.TypeDecl.Value, method.Name.Value, className))
				}
				returnType = types.NewPointer(referencedType)
			}

			// Create function with mangled name to avoid collisions
			funcName := fmt.Sprintf("%s.%s", className, methodName)

			// Create function parameters
			params := make([]*ir.Param, len(paramTypes))
			for i, paramType := range paramTypes {
				var paramName string
				if i == 0 {
					paramName = "self"
				} else {
					paramName = method.Formals[i-1].Name.Value
				}
				params[i] = ir.NewParam(paramName, paramType)
			}

			// Create the function
			function := cg.Module.NewFunc(funcName, returnType, params...)

			// Add to methods map - override any inherited method
			methods[methodName] = function

			// Add to method names if it doesn't exist already
			if !containsString(methodNames, methodName) {
				methodNames = append(methodNames, methodName)
			}
		}
	}

	// Sort method names for consistent ordering
	sort.Strings(methodNames)

	// Build the final vtable in the correct order
	finalMethods := make([]*ir.Func, 0)
	for _, name := range methodNames {
		if method, exists := methods[name]; exists {
			finalMethods = append(finalMethods, method)

			// Update method index for this class
			methodIndices[name] = len(finalMethods) - 1
		}
	}

	// Create an array of function pointers for the vtable
	methodCount := len(finalMethods)
	vtableType := types.NewArray(uint64(methodCount), types.NewPointer(types.I8))

	// Create global array with proper initialization
	vtableName := fmt.Sprintf("vtable.%s", className)

	// Create the initializers for the vtable
	initializers := make([]constant.Constant, methodCount)
	for i, method := range finalMethods {
		// Cast the function pointer to i8*
		initializers[i] = constant.NewBitCast(method, types.NewPointer(types.I8))
	}

	// Create the vtable with the function pointers directly initialized
	var vtable *ir.Global
	if methodCount > 0 {
		// Create a constant array with our initializers
		arrayConst := &constant.Array{
			Typ:   vtableType,
			Elems: initializers,
		}
		vtable = cg.Module.NewGlobalDef(vtableName, arrayConst)
	} else {
		vtable = cg.Module.NewGlobalDef(vtableName, constant.NewZeroInitializer(vtableType))
	}

	// Store in the VTables map
	cg.VTables[className] = vtable

	// Update the method indices map with our final indices
	cg.MethodIndices[className] = methodIndices
}

// Helper function to check if a string is in a slice
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Helper function to check if a method is a built-in method
func isBuiltInMethod(className, methodName string) bool {
	if className == "Object" {
		switch methodName {
		case "abort", "type_name", "copy":
			return true
		}
	} else if className == "IO" {
		switch methodName {
		case "out_string", "out_int", "in_string", "in_int":
			return true
		}
	} else if className == "String" {
		switch methodName {
		case "length", "concat", "substr":
			return true
		}
	}
	return false
}

// findClass finds a class by name in the program
func findClass(name string, program *ast.Program) *ast.Class {
	for _, class := range program.Classes {
		if class.Name.Value == name {
			return class
		}
	}
	return nil
}

// GenerateMethods generates LLVM IR for all class methods
func (cg *CodeGenerator) GenerateMethods(program *ast.Program) {
	// Generate implementation for each class method, including built-in classes
	allClasses := append(program.Classes, cg.BuiltInClasses...)
	for _, class := range allClasses {
		for _, feature := range class.Features {
			if method, isMethod := feature.(*ast.Method); isMethod {
				cg.generateMethod(class, method)
			}
		}
	}
}

// generateMethod creates an LLVM function for a class method
func (cg *CodeGenerator) generateMethod(class *ast.Class, method *ast.Method) {
	className := class.Name.Value
	methodName := method.Name.Value
	mangledName := fmt.Sprintf("%s.%s", className, methodName)

	// Handle built-in methods first
	if className == "Object" {
		switch methodName {
		case "abort":
			cg.generateObjectAbortMethod()
			return
		case "type_name":
			cg.generateTypeNameMethod()
			return
		case "copy":
			cg.generateCopyMethod()
			return
		}
	}

	if className == "IO" {
		switch methodName {
		case "out_string":
			cg.generateIOOutStringMethod()
			return
		case "out_int":
			cg.generateIOOutIntMethod()
			return
		case "in_string":
			cg.generateIOInStringMethod()
			return
		case "in_int":
			cg.generateIOInIntMethod()
			return
		}
	}

	// Handle String class methods
	if className == "String" {
		switch methodName {
		case "length":
			cg.generateStringLengthMethod()
			return
		case "concat":
			cg.generateStringConcatMethod()
			return
		case "substr":
			cg.generateStringSubstrMethod()
			return
		}
	}

	// Find the function declaration (should have been created during vtable generation)
	var methodFunc *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == mangledName {
			methodFunc = f
			break
		}
	}

	if methodFunc == nil {
		panic(fmt.Sprintf("method function not found: %s", mangledName))
	}

	// Skip if this method already has a body
	if len(methodFunc.Blocks) > 0 {
		return
	}

	// Create entry block
	entryBlock := methodFunc.NewBlock("entry")
	cg.CurrentBlock = entryBlock

	// Set up for code generation
	cg.CurrentFunc = methodFunc

	// Set up the symbol table with parameters
	cg.Symbols = make(map[string]value.Value)

	// The first parameter is 'self'
	selfParam := methodFunc.Params[0]
	cg.Symbols["self"] = selfParam

	// Add formal parameters to the symbol table
	for i, formal := range method.Formals {
		cg.Symbols[formal.Name.Value] = methodFunc.Params[i+1] // +1 to skip 'self'
	}

	// Generate code for the method body
	bodyValue := cg.generateExpression(method.Expression)

	// Handle SELF_TYPE return values (methods that return self)
	if method.TypeDecl.Value == "SELF_TYPE" {
		// When returning SELF_TYPE, return the 'self' parameter
		cg.CurrentBlock.NewRet(selfParam)
		return
	}

	// Convert the bodyValue to the correct return type if needed
	returnType := cg.getLLVMTypeForCOOLTypeWithContext(method.TypeDecl.Value, className)

	// Check if we need to convert from primitive type to object pointer
	// For example, converting from i1 (boolean) to %Object*
	_, bodyIsPtr := bodyValue.Type().(*types.PointerType)
	_, returnIsPtr := returnType.(*types.PointerType)

	if bodyValue.Type() != returnType {
		// Special handling for primitive to pointer conversion
		if !bodyIsPtr && returnIsPtr {
			// Handle primitive to pointer conversion
			if bodyValue.Type() == types.I1 || bodyValue.Type() == types.I32 {
				// First convert primitive to i8* with inttoptr
				tmpPtr := cg.CurrentBlock.NewIntToPtr(bodyValue, types.NewPointer(types.I8))
				// Then bitcast to the target type
				bodyValue = cg.CurrentBlock.NewBitCast(tmpPtr, returnType)
			} else {
				// Use ensureCorrectArgumentType for other conversions
				bodyValue = cg.ensureCorrectArgumentType(bodyValue, returnType)
			}
		} else {
			// Standard bitcast for pointer-to-pointer conversion
			bodyValue = cg.CurrentBlock.NewBitCast(bodyValue, returnType)
		}
	}

	// Return the computed value
	cg.CurrentBlock.NewRet(bodyValue)
}

// generateExpression generates code for an expression
func (cg *CodeGenerator) generateExpression(expr ast.Expression) value.Value {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return cg.generateIntegerLiteral(e)
	case *ast.StringLiteral:
		return cg.generateStringLiteral(e)
	case *ast.BooleanLiteral:
		return cg.generateBooleanLiteral(e)
	case *ast.ObjectIdentifier:
		return cg.generateObjectIdentifier(e)
	case *ast.IfExpression:
		return cg.generateIfExpression(e)
	case *ast.WhileExpression:
		return cg.generateWhileExpression(e)
	case *ast.BlockExpression:
		return cg.generateBlockExpression(e)
	case *ast.LetExpression:
		return cg.generateLetExpression(e)
	case *ast.NewExpression:
		return cg.generateNewExpression(e)
	case *ast.IsVoidExpression:
		return cg.generateIsVoidExpression(e)
	case *ast.UnaryExpression:
		return cg.generateUnaryExpression(e)
	case *ast.BinaryExpression:
		return cg.generateBinaryExpression(e)
	case *ast.CaseExpression:
		return cg.generateCaseExpression(e)
	case *ast.CallExpression:
		return cg.generateCallExpression(e)
	case *ast.AssignmentExpression:
		return cg.generateAssignmentExpression(e)
	case *ast.DotCallExpression:
		return cg.generateDotCallExpression(e)
	default:
		panic(fmt.Sprintf("Unsupported expression type: %T", expr))
	}
}

// generateIntegerLiteral creates an LLVM constant integer from a COOL integer literal
func (cg *CodeGenerator) generateIntegerLiteral(intLit *ast.IntegerLiteral) value.Value {
	// Create an LLVM constant integer of type i32 for the COOL integer literal.
	return constant.NewInt(types.I32, int64(intLit.Value))
}

// generateStringLiteral creates an LLVM global string constant from a COOL string literal
func (cg *CodeGenerator) generateStringLiteral(strLit *ast.StringLiteral) value.Value {
	// Create a unique global name for the string constant.
	globalName := fmt.Sprintf(".str%d", len(cg.Module.Globals))

	strConst := constant.NewCharArrayFromString(strLit.Value + "\x00")

	// Define a global constant for the string literal in the module.
	globalStr := cg.Module.NewGlobalDef(globalName, strConst)
	globalStr.Immutable = true
	globalStr.Linkage = enum.LinkageInternal

	// Generate a getelementptr constant to obtain an i8* pointer to the first character in the global array.
	gep := constant.NewGetElementPtr(strConst.Type(), globalStr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	return gep
}

// generateBooleanLiteral creates an LLVM constant integer (i1) representing a COOL boolean literal
func (cg *CodeGenerator) generateBooleanLiteral(boolLit *ast.BooleanLiteral) value.Value {
	// LLVM IR booleans are represented using the i1 type: 0 for false, 1 for true.
	var boolVal int64
	if boolLit.Value {
		boolVal = 1
	} else {
		boolVal = 0
	}
	return constant.NewInt(types.I1, boolVal)
}

// Helper function to calculate size of a type for malloc
func (cg *CodeGenerator) getSizeOf(typ types.Type, block *ir.Block) value.Value {
	// Create a null pointer of the given type
	nullPtr := constant.NewNull(types.NewPointer(typ))

	// Get a pointer to the next element (this gives us the size)
	gep := block.NewGetElementPtr(typ, nullPtr, constant.NewInt(types.I32, 1))

	// Convert this pointer to an integer to get the size in bytes
	return block.NewPtrToInt(gep, types.I64)
}

// generateObjectAllocation creates a new instance of a class
func (cg *CodeGenerator) generateObjectAllocation(typeName string) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Get the LLVM struct type for the class
	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("attempt to create an instance of unknown type: %s", typeName))
	}

	// Calculate the size of the class struct
	sizeValue := cg.getSizeOf(classType, block)

	// Call malloc with the size of the class
	mallocFunc, exists := cg.StdlibFuncs["malloc"]
	if !exists {
		panic("malloc function not found")
	}

	// malloc returns i8* which we'll cast to the appropriate type
	mallocCall := block.NewCall(mallocFunc, sizeValue)

	// Cast the i8* to the class pointer type
	objectPtr := block.NewBitCast(mallocCall, types.NewPointer(classType))

	// Get the vtable for this class
	vtable, exists := cg.VTables[typeName]
	if !exists {
		panic(fmt.Sprintf("vtable not found for type: %s", typeName))
	}

	// Store vtable pointer in the object
	vtableFieldPtr := block.NewGetElementPtr(
		classType,
		objectPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	block.NewStore(
		block.NewBitCast(vtable, types.NewPointer(types.I8)),
		vtableFieldPtr,
	)

	// Initialize attributes with default values
	cg.initializeAttributes(typeName, objectPtr)

	return objectPtr
}

// initializeAttributes initializes all attributes of a class with their default values
func (cg *CodeGenerator) initializeAttributes(className string, objectPtr value.Value) {
	block := cg.CurrentBlock

	// Directly look up the class by name using the map - O(1) operation
	_, exists := cg.ClassNameToAST[className]
	if !exists {
		// If not found in the map, this is unexpected
		return
	}

	// Get the struct type for the class
	classType := cg.TypeMap[className]

	// Save the old 'self' value
	oldSelf, hasSelf := cg.Symbols["self"]

	// Set 'self' to the new object to allow attribute init expressions to access it
	cg.Symbols["self"] = objectPtr

	// Process attributes including inherited ones
	ancestors := []string{className}
	current := className
	for {
		parent, exists := cg.ClassHierarchy[current]
		if !exists || parent == "" {
			break
		}
		ancestors = append([]string{parent}, ancestors...) // Add parent to the beginning
		current = parent
	}

	// Initialize attributes from parent to child
	for _, ancestorName := range ancestors {
		// Efficiently lookup the ancestor class - O(1) operation
		ancestor, exists := cg.ClassNameToAST[ancestorName]
		if !exists {
			continue // Skip if we can't find the ancestor
		}

		// Initialize this class's attributes
		for _, feature := range ancestor.Features {
			if attr, isAttr := feature.(*ast.Attribute); isAttr {
				// Get the attribute index
				attrIndex, exists := cg.AttributeIndices[className][attr.Name.Value]
				if !exists {
					continue // Skip if the attribute index isn't found
				}

				// Get a pointer to the attribute field
				attrPtr := block.NewGetElementPtr(
					classType,
					objectPtr,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, int64(attrIndex)),
				)

				var initValue value.Value

				if attr.Expression != nil {
					// Generate the init expression if it exists
					initValue = cg.generateExpression(attr.Expression)
				} else {
					// If there's no initialization, use a default value based on the type
					switch attr.TypeDecl.Value {
					case "Int":
						initValue = constant.NewInt(types.I32, 0)
					case "Bool":
						initValue = constant.NewInt(types.I1, 0)
					case "String":
						initValue = cg.EmptyStringGlobal
					default:
						// For objects, use null (nil pointer)
						objType := cg.getLLVMTypeForCOOLType(attr.TypeDecl.Value)
						// Check if it's a pointer type before calling NewNull
						if ptrType, ok := objType.(*types.PointerType); ok {
							initValue = constant.NewNull(ptrType)
						} else {
							// If not a pointer type, use a zero value appropriate for the type
							initValue = constant.NewZeroInitializer(objType)
						}
					}
				}

				// Get the attribute type
				attrType := classType.Fields[attrIndex]

				// Make sure types match
				if !initValue.Type().Equal(attrType) {
					// Need to cast if types don't match
					initValue = block.NewBitCast(initValue, attrType)
				}

				// Store the initial value
				block.NewStore(initValue, attrPtr)
			}
		}
	}

	// Restore the old 'self' value
	if hasSelf {
		cg.Symbols["self"] = oldSelf
	} else {
		delete(cg.Symbols, "self")
	}
}

func (cg *CodeGenerator) getObjectRuntimeType(object value.Value) string {
	// Handle primitive types first
	switch object.Type() {
	case types.I1: // Boolean type
		return "Bool"
	case types.I32: // Integer type
		return "Int"
	case types.I8Ptr: // String type (i8*)
		return "String"
	case types.I8: // Single character (might be part of a string)
		return "String"
	}

	// Check if we have a pointer to i8 (string)
	if ptrType, ok := object.Type().(*types.PointerType); ok {
		if ptrType.ElemType == types.I8 {
			return "String"
		}
	}

	// Existing struct type handling
	objPtrType, ok := object.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Sprintf("expected object to be a pointer type, got: %v", object.Type()))
	}

	objStructType, ok := objPtrType.ElemType.(*types.StructType)
	if !ok {
		panic(fmt.Sprintf("expected object to point to a struct type, got: %v", objPtrType.ElemType))
	}

	for name, typ := range cg.TypeMap {
		if typ == objStructType {
			return name
		}
	}

	return "Object"
}

// generateDynamicDispatch generates code for method dispatch
func (cg *CodeGenerator) generateDynamicDispatch(object value.Value, methodName string, args []value.Value) value.Value {
	block := cg.CurrentBlock

	// Get the actual type of the object
	runtimeType := cg.getObjectRuntimeType(object)

	// Get the vtable for this type
	vtable, exists := cg.VTables[runtimeType]
	if !exists {
		panic(fmt.Sprintf("vtable not found for type: %s", runtimeType))
	}

	// Look up the method index in the vtable
	methodIndices, exists := cg.MethodIndices[runtimeType]
	if !exists {
		panic(fmt.Sprintf("method indices not found for type: %s", runtimeType))
	}

	methodIndex, exists := methodIndices[methodName]
	if !exists {
		// Try to find the method in the parent classes
		currentClass := runtimeType
		for {
			parent, exists := cg.ClassHierarchy[currentClass]
			if !exists || parent == "" {
				break
			}
			parentMethodIndices, exists := cg.MethodIndices[parent]
			if !exists {
				currentClass = parent
				continue
			}
			parentMethodIndex, exists := parentMethodIndices[methodName]
			if !exists {
				currentClass = parent
				continue
			}
			// Found the method in a parent class
			methodIndex = parentMethodIndex
			exists = true
			break
		}
		if !exists {
			panic(fmt.Sprintf("method %s not found in class %s or its ancestors", methodName, runtimeType))
		}
	}

	// Get a pointer to the method in the vtable
	methodSlotPtr := block.NewGetElementPtr(
		vtable.ContentType,
		vtable,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(methodIndex)),
	)

	// Load the method function pointer
	methodPtr := block.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	// Find the target method to get its signature
	targetMethod := cg.findMethodByName(runtimeType, methodName)
	if targetMethod == nil {
		panic(fmt.Sprintf("method %s not found in class %s for signature lookup", methodName, runtimeType))
	}

	// Get the method's declaring class (where it's actually defined)
	methodDeclaringClass := cg.findMethodDeclaringClass(runtimeType, methodName)
	if methodDeclaringClass == "" {
		panic(fmt.Sprintf("declaring class for method %s not found", methodName))
	}

	// Create function type for the method (matching its declaration)
	classType, exists := cg.TypeMap[methodDeclaringClass]
	if !exists {
		panic(fmt.Sprintf("class type %s not found", methodDeclaringClass))
	}

	// Build the method's parameter types
	paramTypes := []types.Type{types.NewPointer(classType)} // Self parameter first
	for _, formal := range targetMethod.Formals {
		paramType := cg.getLLVMTypeForCOOLType(formal.TypeDecl.Value)
		paramTypes = append(paramTypes, paramType)
	}

	// Create function type with proper return type
	returnType := cg.getLLVMTypeForCOOLTypeWithContext(targetMethod.TypeDecl.Value, methodDeclaringClass)
	funcType := types.NewPointer(types.NewFunc(returnType, paramTypes...))

	// Cast the method pointer to the correct function type
	castedMethodPtr := block.NewBitCast(methodPtr, funcType)

	// Prepare the arguments, starting with 'self'
	// We need to cast 'self' to the method's declaring class type
	selfArg := block.NewBitCast(object, types.NewPointer(classType))

	// Combine self with the rest of the args, ensuring correct types
	allArgs := make([]value.Value, 0, len(args)+1)
	allArgs = append(allArgs, selfArg)

	// Ensure all arguments have the correct type based on the formal parameter types
	for i, arg := range args {
		if i < len(targetMethod.Formals) {
			paramType := cg.getLLVMTypeForCOOLType(targetMethod.Formals[i].TypeDecl.Value)
			convertedArg := cg.ensureCorrectArgumentType(arg, paramType)
			allArgs = append(allArgs, convertedArg)
		} else {
			// If we have more arguments than formals (shouldn't happen in a well-typed program),
			// just add the argument as is
			allArgs = append(allArgs, arg)
		}
	}

	// Call the method with the properly typed arguments
	return block.NewCall(castedMethodPtr, allArgs...)
}

// Helper method to find the actual class that defines a method
func (cg *CodeGenerator) findMethodDeclaringClass(className, methodName string) string {
	// First check if the method is defined in this class
	class, exists := cg.ClassNameToAST[className]
	if !exists {
		return ""
	}

	// Look for the method in this class's features
	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod && method.Name.Value == methodName {
			return className // Method found in this class
		}
	}

	// If not found, check the parent class
	parent, exists := cg.ClassHierarchy[className]
	if exists && parent != "" {
		return cg.findMethodDeclaringClass(parent, methodName)
	}

	// Method not found
	return ""
}

// Helper method to find a method's AST node
func (cg *CodeGenerator) findMethodByName(className, methodName string) *ast.Method {
	// First check if the method is defined in this class
	class, exists := cg.ClassNameToAST[className]
	if !exists {
		return nil
	}

	// Look for the method in this class's features
	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod && method.Name.Value == methodName {
			return method
		}
	}

	// If not found, check the parent class
	parent, exists := cg.ClassHierarchy[className]
	if exists && parent != "" {
		return cg.findMethodByName(parent, methodName)
	}

	// Method not found
	return nil
}

// Helper function to convert COOL types to LLVM types
func (cg *CodeGenerator) getLLVMTypeForCOOLType(coolType string) types.Type {
	switch coolType {
	case "Int":
		return types.I32
	case "Bool":
		return types.I1
	case "String":
		return types.NewPointer(types.I8)
	case "SELF_TYPE":
		// For SELF_TYPE, we need the current class context
		// Default to Object if we can't determine the current context
		classType, exists := cg.TypeMap["Object"]
		if !exists {
			panic("Object type not found when processing SELF_TYPE")
		}
		return types.NewPointer(classType)
	case "Object", "IO":
		fallthrough
	default:
		// For other class types, use a pointer to the class type
		classType, exists := cg.TypeMap[coolType]
		if !exists {
			panic(fmt.Sprintf("unknown type: %s", coolType))
		}
		return types.NewPointer(classType)
	}
}

// Helper function to convert COOL types to LLVM types with class context
func (cg *CodeGenerator) getLLVMTypeForCOOLTypeWithContext(coolType string, currentClass string) types.Type {
	if coolType == "SELF_TYPE" {
		// For SELF_TYPE, use the current class type
		classType, exists := cg.TypeMap[currentClass]
		if !exists {
			panic(fmt.Sprintf("Current class type %s not found when processing SELF_TYPE", currentClass))
		}
		return types.NewPointer(classType)
	}
	return cg.getLLVMTypeForCOOLType(coolType)
}

// GenerateMain generates the LLVM main function
func (cg *CodeGenerator) GenerateMain(program *ast.Program) {
	// Create the main function with signature: int main()
	mainFunc := cg.Module.NewFunc("main", types.I32)
	entryBlock := mainFunc.NewBlock("entry")

	// Set the current function and block for code generation
	cg.CurrentFunc = mainFunc
	cg.CurrentBlock = entryBlock

	// Create a new instance of the Main class
	mainClass, exists := cg.TypeMap["Main"]
	if !exists {
		panic("Program must have a Main class")
	}

	// Set up the symbol table with 'self' pointing to the Main instance
	cg.Symbols = make(map[string]value.Value)

	// Create a new Main object
	mainObj := cg.generateObjectAllocation("Main")

	// Store the Main object in a local variable for use as 'self'
	mainObjAlloca := entryBlock.NewAlloca(types.NewPointer(mainClass))
	entryBlock.NewStore(mainObj, mainObjAlloca)
	cg.Symbols["self"] = mainObjAlloca

	// Ensure all Main attributes are properly initialized
	cg.initializeAttributes("Main", mainObj)

	// Call the Main.main() method
	// First, find the main method in the vtable of Main
	vtable, exists := cg.VTables["Main"]
	if !exists {
		panic("Main class must have a vtable")
	}

	// Look up the index of the main method in the vtable
	mainMethodIndex, exists := cg.MethodIndices["Main"]["main"]
	if !exists {
		panic("Main class must have a main method")
	}

	// Get a pointer to the method slot in the vtable
	methodSlotPtr := entryBlock.NewGetElementPtr(
		vtable.ContentType,
		vtable,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(mainMethodIndex)),
	)

	// Load the method function pointer
	methodPtr := entryBlock.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	// Cast the i8* function pointer to the correct function type
	// The Main.main method takes no arguments and returns Object (void*)
	funcType := types.NewPointer(types.NewFunc(types.NewPointer(types.I8), types.NewPointer(mainClass)))
	castedMethodPtr := entryBlock.NewBitCast(methodPtr, funcType)

	// Call the Main.main method with 'self' as the argument
	entryBlock.NewCall(castedMethodPtr, mainObj)

	// The C main function should return 0 for success
	entryBlock.NewRet(constant.NewInt(types.I32, 0))
}

// generateObjectIdentifier creates LLVM IR to access a variable by its identifier
func (cg *CodeGenerator) generateObjectIdentifier(identifier *ast.ObjectIdentifier) value.Value {
	// Handle 'self' special case
	if identifier.Value == "self" {
		return cg.Symbols["self"] // 'self' should be in the symbol table already
	}

	block := cg.CurrentBlock

	// Try to find the identifier in the symbol table
	val, exists := cg.Symbols[identifier.Value]
	if exists {
		// Check if the identifier refers to a local variable or parameter (already stored in register)
		if _, isLocalVar := val.(*ir.InstAlloca); isLocalVar {
			// For local variables (alloca instructions), we need to load the value
			load := block.NewLoad(val.Type().(*types.PointerType).ElemType, val)
			return load
		}

		// For non-local variables (global, parameters, etc.)
		return val
	}

	// If not found in the symbol table, it might be a class attribute
	selfPtr, exists := cg.Symbols["self"]
	if !exists {
		panic("'self' not found in symbol table")
	}

	// Get the class name from self's type
	selfPtrType := selfPtr.Type().(*types.PointerType)
	structType := selfPtrType.ElemType.(*types.StructType)
	className := ""

	// Extract class name from struct type name
	for name, typ := range cg.TypeMap {
		if typ == structType {
			className = name
			break
		}
	}

	if className == "" {
		panic("couldn't determine class name for self")
	}

	// Find the attribute index
	attrIndex, exists := cg.AttributeIndices[className][identifier.Value]
	if !exists {
		// Before we panic, check if this identifier is a method in the class
		// This is a special case for when methods are called without self.method
		if _, methodExists := cg.MethodIndices[className][identifier.Value]; methodExists {
			// The identifier is a method in the current class
			// We'll let the call expression handle this as a method on self
			// This just returns self for now, so the call expression can use it
			return selfPtr
		}

		// This should not happen if semantic analysis was successful
		panic(fmt.Sprintf("undefined attribute in class %s: %s", className, identifier.Value))
	}

	// Get the attribute type
	attributeType := structType.Fields[attrIndex]

	// Get a pointer to the attribute
	attrPtr := block.NewGetElementPtr(structType, selfPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(attrIndex)))

	// Load the attribute value
	load := block.NewLoad(attributeType, attrPtr)
	return load
}

// generateAssignmentExpression creates LLVM IR for variable assignment
func (cg *CodeGenerator) generateAssignmentExpression(assign *ast.AssignmentExpression) value.Value {
	// First, generate code for the right-hand side expression
	rhsValue := cg.generateExpression(assign.Expression)

	// Get the current basic block
	block := cg.CurrentBlock

	// Check if it's a local variable (in the symbol table)
	if target, exists := cg.Symbols[assign.Identifier.Value]; exists {
		// For local variables (created with alloca), we use a store instruction
		if allocaInst, isLocalVar := target.(*ir.InstAlloca); isLocalVar {
			// Check if types are compatible, cast if needed
			targetType := allocaInst.Type().(*types.PointerType).ElemType
			if !targetType.Equal(rhsValue.Type()) {
				// We need to cast the value to match the destination type
				rhsValue = block.NewBitCast(rhsValue, targetType)
			}
			block.NewStore(rhsValue, allocaInst)
		} else if _, isParam := target.(*ir.Param); isParam {
			// Parameters should have local storage
			panic("assignment to parameter not properly handled - parameters should have local storage")
		} else {
			// Other cases (e.g., global variables) would be handled here
			panic(fmt.Sprintf("unsupported assignment target type: %T", target))
		}
	} else {
		// If not in symbol table, it must be a class attribute - access through 'self'
		selfPtr, exists := cg.Symbols["self"]
		if !exists {
			panic("'self' not found in symbol table")
		}

		// Get the class name from self's type
		selfPtrType := selfPtr.Type().(*types.PointerType)
		structType := selfPtrType.ElemType.(*types.StructType)
		className := "" // Need to extract class name from type

		// Extract class name from struct type name
		for name, typ := range cg.TypeMap {
			if typ == structType {
				className = name
				break
			}
		}

		if className == "" {
			panic("couldn't determine class name for self")
		}

		// Find the attribute index
		attrIndex, exists := cg.AttributeIndices[className][assign.Identifier.Value]
		if !exists {
			panic(fmt.Sprintf("undefined attribute in class %s: %s", className, assign.Identifier.Value))
		}

		// Get a pointer to the attribute and store the value
		attrPtr := block.NewGetElementPtr(structType, selfPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(attrIndex)))

		// Check if types are compatible, cast if needed
		attributeType := structType.Fields[attrIndex]
		if !attributeType.Equal(rhsValue.Type()) {
			// We need to cast the value to match the destination type
			rhsValue = block.NewBitCast(rhsValue, attributeType)
		}

		block.NewStore(rhsValue, attrPtr)
	}

	// In COOL, an assignment returns the assigned value
	return rhsValue
}

// generateIfExpression creates LLVM IR for a conditional expression
func (cg *CodeGenerator) generateIfExpression(ifExpr *ast.IfExpression) value.Value {
	// Generate code for the condition
	condValue := cg.generateExpression(ifExpr.Condition)

	// Check if condition is a boolean
	// Ideally this was already checked in semantic analysis
	if condValue.Type() != types.I1 {
		panic(fmt.Sprintf("condition in if expression must be of boolean type"))
	}

	// Increment the if counter for unique block names
	cg.IfCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.IfCounter)

	// Create the basic blocks for the true, false, and merge paths with unique names
	currentFunc := cg.CurrentFunc
	trueBlock := currentFunc.NewBlock("if.then" + counterSuffix)
	falseBlock := currentFunc.NewBlock("if.else" + counterSuffix)
	mergeBlock := currentFunc.NewBlock("if.end" + counterSuffix)

	// Create the conditional branch
	cg.CurrentBlock.NewCondBr(condValue, trueBlock, falseBlock)

	// Generate code for the true branch
	cg.CurrentBlock = trueBlock
	trueValue := cg.generateExpression(ifExpr.Consequence)
	trueEndBlock := cg.CurrentBlock
	cg.CurrentBlock.NewBr(mergeBlock)

	// Generate code for the false branch
	cg.CurrentBlock = falseBlock
	falseValue := cg.generateExpression(ifExpr.Alternative)
	falseEndBlock := cg.CurrentBlock
	cg.CurrentBlock.NewBr(mergeBlock)

	// Set current block to merge block
	cg.CurrentBlock = mergeBlock

	// Figure out the common type for the result
	// In COOL, this would be the least common ancestor of the two types
	var resultType types.Type

	// Check if types are the same
	if trueValue.Type().Equal(falseValue.Type()) {
		resultType = trueValue.Type()
	} else {
		// If one is a boolean (i1) and the other is a pointer, or if types don't match otherwise,
		// we need to choose a common type
		_, trueIsPtr := trueValue.Type().(*types.PointerType)
		_, falseIsPtr := falseValue.Type().(*types.PointerType)

		if trueValue.Type() == types.I1 && falseIsPtr {
			// Convert boolean to pointer
			resultType = falseValue.Type()
			// Cast true value to pointer
			intPtrType := types.NewPointer(types.I8)
			trueValue = trueEndBlock.NewIntToPtr(trueValue, intPtrType)
			if !intPtrType.Equal(resultType) {
				trueValue = trueEndBlock.NewBitCast(trueValue, resultType)
			}
		} else if falseValue.Type() == types.I1 && trueIsPtr {
			// Convert boolean to pointer
			resultType = trueValue.Type()
			// Cast false value to pointer
			intPtrType := types.NewPointer(types.I8)
			falseValue = falseEndBlock.NewIntToPtr(falseValue, intPtrType)
			if !intPtrType.Equal(resultType) {
				falseValue = falseEndBlock.NewBitCast(falseValue, resultType)
			}
		} else {
			// For simplicity, use i8* as a generic object pointer type
			// In a full implementation, you would calculate the least common ancestor type
			resultType = types.NewPointer(types.I8)

			// Cast both values to the common type if needed
			if !trueValue.Type().Equal(resultType) {
				if trueValue.Type() == types.I1 {
					trueValue = trueEndBlock.NewIntToPtr(trueValue, resultType)
				} else if trueIsPtr {
					trueValue = trueEndBlock.NewBitCast(trueValue, resultType)
				}
			}

			if !falseValue.Type().Equal(resultType) {
				if falseValue.Type() == types.I1 {
					falseValue = falseEndBlock.NewIntToPtr(falseValue, resultType)
				} else if falseIsPtr {
					falseValue = falseEndBlock.NewBitCast(falseValue, resultType)
				}
			}
		}
	}

	// Create a PHI node with incoming values right away
	phi := cg.CurrentBlock.NewPhi(
		&ir.Incoming{X: trueValue, Pred: trueEndBlock},
		&ir.Incoming{X: falseValue, Pred: falseEndBlock},
	)

	// Set the correct type for the PHI node
	phi.Typ = resultType

	// Return the PHI node as the result of the case expression
	return phi
}

// generateWhileExpression creates LLVM IR for a while loop expression
func (cg *CodeGenerator) generateWhileExpression(whileExpr *ast.WhileExpression) value.Value {
	// Increment the while counter for unique block names
	cg.WhileCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.WhileCounter)

	// Create the basic blocks for the loop
	currentFunc := cg.CurrentFunc
	condBlock := currentFunc.NewBlock("while.cond" + counterSuffix)
	bodyBlock := currentFunc.NewBlock("while.body" + counterSuffix)
	exitBlock := currentFunc.NewBlock("while.exit" + counterSuffix)

	// Branch from current block to condition block
	cg.CurrentBlock.NewBr(condBlock)

	// Set current block to condition block
	cg.CurrentBlock = condBlock

	// Generate code for the condition
	condValue := cg.generateExpression(whileExpr.Condition)

	// Check if condition is a boolean
	if condValue.Type() != types.I1 {
		panic(fmt.Sprintf("condition in while expression must be of boolean type"))
	}

	// Create conditional branch: if condition is true, enter body, otherwise exit
	cg.CurrentBlock.NewCondBr(condValue, bodyBlock, exitBlock)

	// Generate code for the loop body
	cg.CurrentBlock = bodyBlock
	cg.generateExpression(whileExpr.Body)

	// After executing the body, jump back to the condition block
	cg.CurrentBlock.NewBr(condBlock)

	// Set the current block to the exit block
	cg.CurrentBlock = exitBlock

	// In COOL, while loops return void (represented as null in our case)
	// For COOL semantics, this would typically be a reference to the "Object" class's null instance
	// We'll use a null pointer of type i8* as a generic object pointer
	nullValue := constant.NewNull(types.NewPointer(types.I8))

	return nullValue
}

// generateBlockExpression creates LLVM IR for a block of expressions
func (cg *CodeGenerator) generateBlockExpression(blockExpr *ast.BlockExpression) value.Value {
	// In COOL, a block consists of a sequence of expressions separated by semicolons
	// The value of a block is the value of the last expression

	var lastValue value.Value

	// Generate code for each expression in the block
	for _, expr := range blockExpr.Expressions {
		// Generate the expression
		lastValue = cg.generateExpression(expr)
	}

	// If the block is empty (no expressions), return a void value
	if lastValue == nil {
		// Use a null pointer as a generic "void" value
		return constant.NewNull(types.NewPointer(types.I8))
	}

	// Return the value of the last expression
	return lastValue
}

// generateBinaryExpression creates LLVM IR for binary operations
func (cg *CodeGenerator) generateBinaryExpression(binExpr *ast.BinaryExpression) value.Value {
	// Get the current basic block
	block := cg.CurrentBlock

	// Generate code for the left and right operands
	leftValue := cg.generateExpression(binExpr.Left)
	rightValue := cg.generateExpression(binExpr.Right)

	// Handle the operation based on the operator
	switch binExpr.Operator {
	// Arithmetic operations
	case "+":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewAdd(leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '+' must be integers"))
		}

	case "-":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewSub(leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '-' must be integers"))
		}

	case "*":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewMul(leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '*' must be integers"))
		}

	case "/":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			// In COOL, division by zero results in an error (implementation-dependent)
			// Here we'll just do the division directly, but in a more complete implementation,
			// we might want to add a runtime check for division by zero
			return block.NewSDiv(leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '/' must be integers"))
		}

	// Comparison operations
	case "<":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredSLT, leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '<' must be integers"))
		}

	case "<=":
		// Check if both operands are integers
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredSLE, leftValue, rightValue)
		} else {
			panic(fmt.Sprintf("operands of '<=' must be integers"))
		}

	case "=":
		// In COOL, equality is defined for all types
		// For primitive types like Int, String, and Bool, we compare values
		// For objects, we compare references

		// Check if types are the same
		if !leftValue.Type().Equal(rightValue.Type()) {
			// In COOL, comparing objects of different types is false
			return constant.NewInt(types.I1, 0)
		}

		// For integer types, compare with ICmp
		if leftValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		// For boolean types, compare with ICmp
		if leftValue.Type() == types.I1 {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		// For pointer types (objects, strings), compare the pointers
		if _, isPtr := leftValue.Type().(*types.PointerType); isPtr {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		// For any other type, we're not sure what to do
		panic(fmt.Sprintf("equality comparison not implemented for type: %v", leftValue.Type()))
	}

	// Add panic for unhandled operators
	panic(fmt.Sprintf("unsupported binary operator: %s", binExpr.Operator))
}

// generateDotCallExpression creates LLVM IR for method calls on objects
func (cg *CodeGenerator) generateDotCallExpression(dotCall *ast.DotCallExpression) value.Value {
	// Generate code for the object on which the method is called
	objectValue := cg.generateExpression(dotCall.Object)

	// Get the current block
	block := cg.CurrentBlock

	// Handle case expression results specially when used as method receivers
	// Check if the object value is a result of a case expression
	caseResult, isCaseResult := cg.CaseResults[objectValue]
	if isCaseResult {
		// Special handling for case results - use the stored type information
		// TODO: Use the branch types information to ensure correct dispatch
		objectValue = caseResult.Value
	}

	// Generate LLVM values for all arguments
	argValues := make([]value.Value, 0, len(dotCall.Arguments)+1)

	// The first argument to a method call is always the object itself (self)
	argValues = append(argValues, objectValue)

	// Add the rest of the arguments
	for _, arg := range dotCall.Arguments {
		argValues = append(argValues, cg.generateExpression(arg))
	}

	// Special case for IO.out_string and IO.out_int to call runtime functions directly
	if dotCall.Method.Value == "out_string" || dotCall.Method.Value == "out_int" {
		// First check if the object is an IO object or can be cast to an IO object
		var isIOObject bool
		if ptrType, isPtr := objectValue.Type().(*types.PointerType); isPtr {
			if structType, isStruct := ptrType.ElemType.(*types.StructType); isStruct {
				isIOObject = structType.Name() == "IO"
			}
		}

		// Try to handle IO methods
		if isIOObject || cg.canCastTo(objectValue, "IO") {
			// Cast object to IO if needed
			ioObjValue := objectValue
			if !isIOObject {
				ioType, exists := cg.TypeMap["IO"]
				if exists {
					ioObjValue = block.NewBitCast(objectValue, types.NewPointer(ioType))
				}
			}

			if dotCall.Method.Value == "out_string" && len(argValues) > 1 {
				// Find or create the runtime function
				var outStringFunc *ir.Func
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_string" {
						outStringFunc = f
						break
					}
				}

				if outStringFunc == nil {
					// If the function isn't already declared, declare it
					outStringFunc = cg.Module.NewFunc("IO.out_string", types.NewPointer(cg.TypeMap["IO"]),
						ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						ir.NewParam("str", types.NewPointer(types.I8)))
				}

				// Ensure string argument has the right type
				strArg := argValues[1]
				convertedStrArg := cg.ensureCorrectArgumentType(strArg, types.NewPointer(types.I8))

				// Make the call with the properly typed arguments
				callArgs := []value.Value{ioObjValue, convertedStrArg}
				block.NewCall(outStringFunc, callArgs...)

				// Return the IO object itself
				return ioObjValue

			} else if dotCall.Method.Value == "out_int" && len(argValues) > 1 {
				// Find or create the runtime function
				var outIntFunc *ir.Func
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_int" {
						outIntFunc = f
						break
					}
				}

				if outIntFunc == nil {
					// If the function isn't already declared, declare it
					outIntFunc = cg.Module.NewFunc("IO.out_int", types.NewPointer(cg.TypeMap["IO"]),
						ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						ir.NewParam("n", types.I32))
				}

				// Ensure int argument has the right type
				intArg := argValues[1]
				convertedIntArg := cg.ensureCorrectArgumentType(intArg, types.I32)

				// Make the call with properly typed arguments
				callArgs := []value.Value{ioObjValue, convertedIntArg}
				block.NewCall(outIntFunc, callArgs...)

				// Return the IO object itself
				return ioObjValue
			}
		}
	}

	// Check if this is a static dispatch (explicitly specifying the type)
	if dotCall.Type != nil {
		// Get the target type name
		targetTypeName := dotCall.Type.Value

		// Use static dispatch when the target type is explicitly specified
		return cg.generateStaticDispatch(
			objectValue,
			targetTypeName,
			dotCall.Method.Value,
			argValues[1:], // Skip the 'self' argument which was already added
		)
	}

	// For all other cases, use dynamic dispatch based on the runtime type of the object
	return cg.generateDynamicDispatch(
		objectValue,
		dotCall.Method.Value,
		argValues[1:], // Skip the 'self' argument which was already added
	)
}

// generateStaticDispatch creates LLVM IR for static dispatch (obj@Type.method())
func (cg *CodeGenerator) generateStaticDispatch(object value.Value, typeName string, methodName string, args []value.Value) value.Value {
	block := cg.CurrentBlock

	// Find the actual method function by name
	methodFuncName := fmt.Sprintf("%s.%s", typeName, methodName)
	var methodFunc *ir.Func

	// Look for the method in the module
	for _, f := range cg.Module.Funcs {
		if f.Name() == methodFuncName {
			methodFunc = f
			break
		}
	}

	if methodFunc == nil {
		panic(fmt.Sprintf("method %s not found in class %s", methodName, typeName))
	}

	// Get the target class type to cast the object
	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("class type %s not found", typeName))
	}

	// Find the target method to get its signature for correct argument conversion
	targetMethod := cg.findMethodByName(typeName, methodName)
	if targetMethod == nil {
		panic(fmt.Sprintf("method %s not found in class %s for signature lookup", methodName, typeName))
	}

	// Cast the object to the target class type
	castedObject := block.NewBitCast(object, types.NewPointer(classType))

	// Create a new list of arguments starting with the properly typed object
	allArgs := make([]value.Value, 0, len(args)+1)
	allArgs = append(allArgs, castedObject) // Add properly cast 'self' as the first argument

	// Ensure all arguments have the correct type based on the formal parameter types
	for i, arg := range args {
		if i < len(targetMethod.Formals) {
			paramType := cg.getLLVMTypeForCOOLType(targetMethod.Formals[i].TypeDecl.Value)
			convertedArg := cg.ensureCorrectArgumentType(arg, paramType)
			allArgs = append(allArgs, convertedArg)
		} else {
			// If we have more arguments than formals (shouldn't happen in a well-typed program),
			// just add the argument as is
			allArgs = append(allArgs, arg)
		}
	}

	// In static dispatch, we call the method directly by name rather than through the vtable
	call := block.NewCall(methodFunc, allArgs...)

	return call
}

// generateLetExpression creates LLVM IR for let expressions
func (cg *CodeGenerator) generateLetExpression(letExpr *ast.LetExpression) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Save the old symbol table to restore after the let expression
	oldSymbols := make(map[string]value.Value)
	for k, v := range cg.Symbols {
		oldSymbols[k] = v
	}

	// Process each binding in the let expression
	for _, binding := range letExpr.Bindings {
		// Allocate space for the variable on the stack
		varName := binding.Identifier.Value

		// Determine the LLVM type for the variable based on the COOL type
		var varType types.Type

		// For basic types, map them directly
		switch binding.Type.Value {
		case "Int":
			varType = types.I32
		case "Bool":
			varType = types.I1
		case "String":
			varType = types.NewPointer(types.I8) // Strings are pointers to char arrays
		default:
			// For class types, use a pointer to the class struct
			classType, exists := cg.TypeMap[binding.Type.Value]
			if !exists {
				panic(fmt.Sprintf("unknown type in let binding: %s", binding.Type.Value))
			}
			varType = types.NewPointer(classType)
		}

		// Allocate space for the variable
		alloca := block.NewAlloca(varType)

		// Initialize the variable
		var initValue value.Value

		if binding.Init != nil {
			// If there's an initialization expression, evaluate it
			initValue = cg.generateExpression(binding.Init)

			// Make sure types match
			if !initValue.Type().Equal(varType) {
				// If types don't match, we might need to cast
				// For example, if assigning a subclass instance to a superclass variable

				// For simplicity, just assume we need a bitcast if types don't match
				initValue = block.NewBitCast(initValue, varType)
			}
		} else {
			// If there's no initialization, use a default value based on the type
			switch binding.Type.Value {
			case "Int":
				initValue = constant.NewInt(types.I32, 0)
			case "Bool":
				initValue = constant.NewInt(types.I1, 0)
			case "String":
				// Use the pre-defined empty string global instead of creating a new one
				initValue = constant.NewGetElementPtr(
					cg.EmptyStringGlobal.ContentType,
					cg.EmptyStringGlobal,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, 0),
				)
			default:
				// For objects, use null
				ptrType, ok := varType.(*types.PointerType)
				if !ok {
					panic(fmt.Sprintf("expected pointer type for class type variable, got: %v", varType))
				}
				initValue = constant.NewNull(ptrType)
			}
		}

		// Store the initial value in the allocated space
		block.NewStore(initValue, alloca)

		// Add the variable to the symbol table
		cg.Symbols[varName] = alloca
	}

	// Generate code for the body of the let expression
	bodyValue := cg.generateExpression(letExpr.In)

	// Restore the old symbol table
	cg.Symbols = oldSymbols

	// The value of the let expression is the value of its body
	return bodyValue
}

// generateNewExpression creates LLVM IR for object instantiation
func (cg *CodeGenerator) generateNewExpression(newExpr *ast.NewExpression) value.Value {
	// Get the type name from the NewExpression
	typeName := newExpr.Type.Value

	// Handle special case for basic COOL types
	switch typeName {
	case "Int":
		// Create a new Int object with default value 0
		return constant.NewInt(types.I32, 0)

	case "Bool":
		// Create a new Bool object with default value false
		return constant.NewInt(types.I1, 0)

	case "String":
		return constant.NewGetElementPtr(
			cg.EmptyStringGlobal.ContentType,
			cg.EmptyStringGlobal,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0),
		)

	default:
		// For user-defined classes, use the object allocation function
		return cg.generateObjectAllocation(typeName)
	}
}

// generateIsVoidExpression creates LLVM IR for checking if a reference is null
func (cg *CodeGenerator) generateIsVoidExpression(isVoidExpr *ast.IsVoidExpression) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Generate code for the expression to check
	exprValue := cg.generateExpression(isVoidExpr.Expression)

	// The expression type determines how we check for "void" (null)
	switch exprValue.Type() {
	case types.I32, types.I1:
		// For primitive types like Int and Bool, they can never be void
		// Always return false
		return constant.NewInt(types.I1, 0)

	default:
		// For reference types (objects, strings), check if the pointer is null
		if ptrType, isPtr := exprValue.Type().(*types.PointerType); isPtr {
			// Compare the pointer with null
			nullVal := constant.NewNull(ptrType)
			return block.NewICmp(enum.IPredEQ, exprValue, nullVal)
		} else {
			// If it's not a pointer type or a primitive type we know about,
			// we're not sure what to do, so panic
			panic(fmt.Sprintf("isvoid check not implemented for type: %v", exprValue.Type()))
		}
	}
}

// generateUnaryExpression creates LLVM IR for unary operations
func (cg *CodeGenerator) generateUnaryExpression(unaryExpr *ast.UnaryExpression) value.Value {
	block := cg.CurrentBlock

	// Generate code for the expression being operated on
	exprValue := cg.generateExpression(unaryExpr.Right)

	// Handle different operators
	switch unaryExpr.Operator {
	case "~": // Integer negation
		// Check if the operand is an integer
		if exprValue.Type() != types.I32 {
			panic(fmt.Sprintf("operand of integer negation (~) must be an integer"))
		}

		// Negate the integer value
		// In LLVM IR, negation is implemented as 0 - value
		zero := constant.NewInt(types.I32, 0)
		return block.NewSub(zero, exprValue)

	case "not": // Boolean NOT
		// Check if the operand is a boolean
		if exprValue.Type() != types.I1 {
			panic(fmt.Sprintf("operand of boolean NOT (not) must be a boolean"))
		}

		// Perform logical NOT
		// In LLVM IR, this is done with XOR with true (1)
		// Or we can simply use the LLVM 'not' instruction
		return block.NewXor(exprValue, constant.NewInt(types.I1, 1))

	default:
		panic(fmt.Sprintf("unsupported unary operator: %s", unaryExpr.Operator))
	}
}

// generateCaseExpression creates LLVM IR for COOL's case expressions
func (cg *CodeGenerator) generateCaseExpression(caseExpr *ast.CaseExpression) value.Value {
	// Increment the case counter for unique block names
	cg.CaseCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.CaseCounter)

	// Get the current function
	currentFunc := cg.CurrentFunc

	// Generate code for the expression being dispatched on
	exprValue := cg.generateExpression(caseExpr.Expression)

	// Create a basic block for the end of the case expression
	// All branches will merge to this block
	endBlock := currentFunc.NewBlock("case.end" + counterSuffix)

	// Create a basic block for each branch
	branchBlocks := make([]*ir.Block, len(caseExpr.Branches))
	for i := range caseExpr.Branches {
		branchBlocks[i] = currentFunc.NewBlock(fmt.Sprintf("case.branch.%d%s", i, counterSuffix))
	}

	// This should never happen in well-typed COOL code, but we need it for LLVM
	noMatchBlock := currentFunc.NewBlock("case.nomatch" + counterSuffix)

	// Get the current block - we'll be branching from here
	currentBlock := cg.CurrentBlock

	// First, we need to check if the object is null
	// In COOL, a case expression with a null object raises a runtime error
	// We'll add a null check if the expression is a reference type
	if ptrType, isPtr := exprValue.Type().(*types.PointerType); isPtr {
		// Create a block for the null check
		notNullBlock := currentFunc.NewBlock("case.notnull" + counterSuffix)

		// Compare the object with null
		nullVal := constant.NewNull(ptrType)
		isNull := currentBlock.NewICmp(enum.IPredEQ, exprValue, nullVal)

		// If the object is null, jump to an error handler
		// In a real compiler, we would call a runtime error function
		// For simplicity, we'll just use the nomatch block
		currentBlock.NewCondBr(isNull, noMatchBlock, notNullBlock)

		// Set the current block to the not-null block for further code generation
		currentBlock = notNullBlock
		cg.CurrentBlock = notNullBlock
	}

	// Get the runtime type of the object for branching based on type
	// Create a temporary block for getting the runtime type
	typeCheckBlock := currentFunc.NewBlock("case.typecheck" + counterSuffix)
	currentBlock.NewBr(typeCheckBlock)
	cg.CurrentBlock = typeCheckBlock

	// Get the runtime type of the object
	objectType := cg.getObjectRuntimeType(exprValue)

	// Keep track of the current branching block
	branchingBlock := typeCheckBlock

	// Create branch decision blocks
	decisionBlocks := make([]*ir.Block, len(caseExpr.Branches))
	for i := range caseExpr.Branches {
		decisionBlocks[i] = currentFunc.NewBlock(fmt.Sprintf("case.decision.%d%s", i, counterSuffix))
	}

	// Set up the branch chain for type checking
	// Starting with the first branch
	branchingBlock.NewBr(decisionBlocks[0])

	// Process each branch with actual type checking
	for i, branch := range caseExpr.Branches {
		// Set current block to this branch's decision block
		branchingBlock = decisionBlocks[i]

		// Get the type declared in this branch
		branchType := branch.Type.Value

		// In COOL, a case branch matches if the object's type conforms to the branch's type
		// We need to check if objectType is a subtype of branchType
		// This would involve checking the class hierarchy

		var matchesCondition value.Value

		// Check if the runtime type matches or is a subtype of the branch type
		if branchType == objectType {
			// Direct match
			matchesCondition = constant.NewInt(types.I1, 1)
		} else {
			// Need to check class hierarchy
			// Start with direct match check
			matchesCondition = constant.NewInt(types.I1, 0)

			// Check inheritance chain
			currentType := objectType
			for {
				parent, exists := cg.ClassHierarchy[currentType]
				if !exists || parent == "" {
					break // Reached Object or unknown class
				}

				if parent == branchType {
					// Found a match in the hierarchy
					matchesCondition = constant.NewInt(types.I1, 1)
					break
				}

				currentType = parent
			}
		}

		// If this is the last branch, the next block is the no-match block
		// Otherwise, it's the next decision block
		var nextBlock *ir.Block
		if i < len(caseExpr.Branches)-1 {
			nextBlock = decisionBlocks[i+1]
		} else {
			nextBlock = noMatchBlock
		}

		// Branch based on the type match condition
		branchingBlock.NewCondBr(matchesCondition, branchBlocks[i], nextBlock)
	}

	// Add code to handle the case where no branch matches (a runtime error in COOL)
	// This should never happen in well-typed COOL programs
	// In a real implementation, this would call a runtime error function
	noMatchBlock.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
	noMatchBlock.NewUnreachable()

	// Generate code for each branch
	branchValues := make([]value.Value, len(caseExpr.Branches))
	branchRealTypes := make([]string, len(caseExpr.Branches)) // Track actual COOL types for branches
	branchEndBlocks := make([]*ir.Block, len(caseExpr.Branches))

	for i, branch := range caseExpr.Branches {
		// Set the current block to the branch block
		cg.CurrentBlock = branchBlocks[i]

		// Save old symbol table
		oldSymbols := make(map[string]value.Value)
		for k, v := range cg.Symbols {
			oldSymbols[k] = v
		}

		// Cast the expression to the branch type
		var castedValue value.Value
		if branchType, exists := cg.TypeMap[branch.Type.Value]; exists {
			// Check if we're trying to cast from a primitive type to an object type
			_, isExprPtr := exprValue.Type().(*types.PointerType)
			if !isExprPtr && (exprValue.Type() == types.I32 || exprValue.Type() == types.I1) {
				// First convert to i8* using inttoptr
				tmpPtr := cg.CurrentBlock.NewIntToPtr(exprValue, types.NewPointer(types.I8))
				// Then bitcast to the target type
				castedValue = cg.CurrentBlock.NewBitCast(tmpPtr, types.NewPointer(branchType))
			} else {
				// Normal pointer-to-pointer cast
				castedValue = cg.CurrentBlock.NewBitCast(exprValue, types.NewPointer(branchType))
			}
		} else {
			castedValue = exprValue // Fallback if type doesn't exist (shouldn't happen)
		}

		// Add the branch variable to the symbol table with the properly typed value
		cg.Symbols[branch.Identifier.Value] = castedValue

		// Generate code for the branch expression
		branchValues[i] = cg.generateExpression(branch.Expression)
		branchRealTypes[i] = branch.Type.Value // Store the COOL type for this branch

		// Restore the old symbol table
		cg.Symbols = oldSymbols

		// Get the current block after generating the branch expression
		branchEndBlocks[i] = cg.CurrentBlock

		// Branch to the end block
		branchEndBlocks[i].NewBr(endBlock)
	}

	// Set the current block to the end block
	cg.CurrentBlock = endBlock

	// Create a PHI node to merge all branch values
	// First, determine the common type for the result
	// In COOL, this would be the least common ancestor of all branch types
	var resultType types.Type
	if len(branchValues) > 0 {
		resultType = branchValues[0].Type()

		// Check if all branches have the same type
		for _, val := range branchValues[1:] {
			if !val.Type().Equal(resultType) {
				// If types don't match, use a generic object pointer type
				// In a full implementation, you would calculate the least common ancestor type
				resultType = types.NewPointer(types.I8)
				break
			}
		}
	} else {
		// If there are no branches (shouldn't happen in valid COOL), use a generic object pointer
		resultType = types.NewPointer(types.I8)
	}

	// Create the PHI node with the result type
	phi := &ir.InstPhi{Typ: resultType}
	endBlock.Insts = append(endBlock.Insts, phi)

	// Add incoming values for the PHI node
	for i, val := range branchValues {
		// If the branch value type doesn't match the result type, cast it
		if !val.Type().Equal(resultType) {
			// Check if we need to convert from primitive to pointer type
			_, valIsPtr := val.Type().(*types.PointerType)
			_, resultIsPtr := resultType.(*types.PointerType)

			if !valIsPtr && resultIsPtr && (val.Type() == types.I1 || val.Type() == types.I32) {
				// First convert primitive to i8* with inttoptr
				tmpPtr := branchEndBlocks[i].NewIntToPtr(val, types.NewPointer(types.I8))
				// Then bitcast to the target type
				val = branchEndBlocks[i].NewBitCast(tmpPtr, resultType)
			} else {
				// Use ensureCorrectArgumentType for other cases
				val = cg.ensureCorrectArgumentType(val, resultType)
			}
		}

		phi.Incs = append(phi.Incs, &ir.Incoming{X: val, Pred: branchEndBlocks[i]})
	}

	// Create a wrapper for the result that includes both the value and its COOL type
	result := &CaseResult{
		Value:       phi,
		BranchTypes: branchRealTypes,
	}

	// Store the result in a map so it can be retrieved later when needed
	if cg.CaseResults == nil {
		cg.CaseResults = make(map[value.Value]*CaseResult)
	}
	cg.CaseResults[phi] = result

	return phi
}

// generateCallExpression creates LLVM IR for function calls
func (cg *CodeGenerator) generateCallExpression(callExpr *ast.CallExpression) value.Value {
	// In COOL, a direct call without a receiver (e.g., factorial(5)) is implicitly a self call
	// So we need to get the self reference from the current function
	if cg.CurrentFunc == nil {
		panic("Cannot call method without object outside of method context")
	}

	// Self is always the first parameter in method functions
	selfObj := cg.CurrentFunc.Params[0]

	// Extract the method name from the Function expression
	// For a direct method call, Function should be an ObjectIdentifier
	methodIdent, ok := callExpr.Function.(*ast.ObjectIdentifier)
	if !ok {
		// This is unexpected - in COOL syntax the function part of a call expression
		// should be an identifier
		panic(fmt.Sprintf("Unexpected Function type in CallExpression: %T", callExpr.Function))
	}

	methodName := methodIdent.Value

	// Generate all arguments
	args := make([]value.Value, 0, len(callExpr.Arguments))
	for _, arg := range callExpr.Arguments {
		args = append(args, cg.generateExpression(arg))
	}

	// For simple method calls, we use dynamic dispatch
	return cg.generateDynamicDispatch(selfObj, methodName, args)
}

// DefineBuiltInClasses defines the built-in COOL classes: Object, IO, Int, String, Bool
func (cg *CodeGenerator) DefineBuiltInClasses() {
	// Define Object class - the root of the inheritance hierarchy
	objectClass := &ast.Class{
		Name: &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "abort"},
				TypeDecl: &ast.TypeIdentifier{Value: "Object"},
				Formals:  []*ast.Formal{},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "type_name"},
				TypeDecl: &ast.TypeIdentifier{Value: "String"},
				Formals:  []*ast.Formal{},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "copy"},
				TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
				Formals:  []*ast.Formal{},
			},
		},
	}
	cg.declareClassType(objectClass)
	cg.defineClassStruct(objectClass)
	cg.BuiltInClasses = append(cg.BuiltInClasses, objectClass)

	// Define IO class - for input/output operations
	ioClass := &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "IO"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "out_string"},
				TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
				Formals: []*ast.Formal{
					{
						Name:     &ast.ObjectIdentifier{Value: "x"},
						TypeDecl: &ast.TypeIdentifier{Value: "String"},
					},
				},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "out_int"},
				TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
				Formals: []*ast.Formal{
					{
						Name:     &ast.ObjectIdentifier{Value: "x"},
						TypeDecl: &ast.TypeIdentifier{Value: "Int"},
					},
				},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "in_string"},
				TypeDecl: &ast.TypeIdentifier{Value: "String"},
				Formals:  []*ast.Formal{},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "in_int"},
				TypeDecl: &ast.TypeIdentifier{Value: "Int"},
				Formals:  []*ast.Formal{},
			},
		},
	}
	cg.declareClassType(ioClass)
	cg.defineClassStruct(ioClass)
	cg.BuiltInClasses = append(cg.BuiltInClasses, ioClass)

	// Define Int class
	intClass := &ast.Class{
		Name:     &ast.TypeIdentifier{Value: "Int"},
		Parent:   &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{},
	}
	cg.declareClassType(intClass)
	cg.defineClassStruct(intClass)
	cg.BuiltInClasses = append(cg.BuiltInClasses, intClass)

	// Define String class
	stringClass := &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "String"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "length"},
				TypeDecl: &ast.TypeIdentifier{Value: "Int"},
				Formals:  []*ast.Formal{},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "concat"},
				TypeDecl: &ast.TypeIdentifier{Value: "String"},
				Formals: []*ast.Formal{
					{
						Name:     &ast.ObjectIdentifier{Value: "s"},
						TypeDecl: &ast.TypeIdentifier{Value: "String"},
					},
				},
			},
			&ast.Method{
				Name:     &ast.ObjectIdentifier{Value: "substr"},
				TypeDecl: &ast.TypeIdentifier{Value: "String"},
				Formals: []*ast.Formal{
					{
						Name:     &ast.ObjectIdentifier{Value: "i"},
						TypeDecl: &ast.TypeIdentifier{Value: "Int"},
					},
					{
						Name:     &ast.ObjectIdentifier{Value: "l"},
						TypeDecl: &ast.TypeIdentifier{Value: "Int"},
					},
				},
			},
		},
	}
	cg.declareClassType(stringClass)
	cg.defineClassStruct(stringClass)
	cg.BuiltInClasses = append(cg.BuiltInClasses, stringClass)

	// Define Bool class
	boolClass := &ast.Class{
		Name:     &ast.TypeIdentifier{Value: "Bool"},
		Parent:   &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{},
	}
	cg.declareClassType(boolClass)
	cg.defineClassStruct(boolClass)
	cg.BuiltInClasses = append(cg.BuiltInClasses, boolClass)

	// Create vtables for all built-in classes
	// We need to do this in a separate step after all classes are defined
	// to handle inheritance properly
	program := &ast.Program{
		Classes: []*ast.Class{objectClass, ioClass, intClass, stringClass, boolClass},
	}

	// Create vtables for all the built-in classes
	for _, class := range program.Classes {
		cg.createVTableForClass(class, program)
	}

	// When defining the String class struct:
	stringStruct := types.NewStruct(
		types.NewPointer(types.I8), // vtable pointer (index 0)
		types.NewPointer(types.I8), // actual string data (index 1)
	)
	cg.TypeMap["String"] = stringStruct
}

// Add this new method to handle IO.out_string code generation
func (cg *CodeGenerator) generateIOOutStringMethod() {
	// Define the out_string method for IO class
	funcName := "IO.out_string"

	// Get the function if it already exists
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		panic("IO.out_string function declaration not found")
	}

	// Create entry block
	entry := funcDecl.NewBlock("entry")

	// Get the string parameter (the second parameter, after self)
	strParam := funcDecl.Params[1]

	// Create format string constant for printf
	fmtStr := cg.Module.NewGlobalDef(".str.fmt", constant.NewCharArrayFromString("%s\x00"))
	fmtPtr := constant.NewGetElementPtr(fmtStr.ContentType, fmtStr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	// Call printf directly, with the format string and the string parameter
	printfFunc := cg.StdlibFuncs["printf"]
	entry.NewCall(printfFunc, fmtPtr, strParam)

	// Return self
	entry.NewRet(funcDecl.Params[0])
}

// Add this new method to handle IO.out_int code generation
func (cg *CodeGenerator) generateIOOutIntMethod() {
	funcName := "IO.out_int"
	ioType := cg.TypeMap["IO"]

	// Get or create the function
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(ioType),
			ir.NewParam("self", types.NewPointer(ioType)),
			ir.NewParam("x", types.I32),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Create integer format string
	fmtStr := cg.Module.NewGlobalDef(".str.fmt.int", constant.NewCharArrayFromString("%d\x00"))
	fmtStr.Immutable = true
	fmtPtr := constant.NewGetElementPtr(fmtStr.ContentType, fmtStr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	// Call printf
	entry.NewCall(cg.StdlibFuncs["printf"], fmtPtr, funcDecl.Params[1])

	// Return self
	entry.NewRet(funcDecl.Params[0])
}

// Add these new methods to handle Object built-in methods
func (cg *CodeGenerator) generateObjectAbortMethod() {
	funcName := "Object.abort"
	objType := cg.TypeMap["Object"]

	// Get or create the function
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(objType),
			ir.NewParam("self", types.NewPointer(objType)),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Just exit with status code 1 without printing any message
	entry.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
	entry.NewUnreachable() // exit doesn't return
}

func (cg *CodeGenerator) generateTypeNameMethod() {
	funcName := "Object.type_name"
	objType := cg.TypeMap["Object"]

	// Get or create the function
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(types.I8), // Returns String (i8*)
			ir.NewParam("self", types.NewPointer(objType)),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Create class name string
	className := "Object"
	strConst := constant.NewCharArrayFromString(className + "\x00")
	global := cg.Module.NewGlobalDef(fmt.Sprintf(".str.%s", className), strConst)
	global.Immutable = true
	gep := constant.NewGetElementPtr(strConst.Type(), global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	entry.NewRet(gep)
}

func (cg *CodeGenerator) generateCopyMethod() {
	funcName := "Object.copy"
	objType := cg.TypeMap["Object"]

	// Get or create the function
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(objType), // Returns SELF_TYPE
			ir.NewParam("self", types.NewPointer(objType)),
		)
	}

	entry := funcDecl.NewBlock("entry")
	// In a real implementation this would perform a shallow copy
	// For now just return self
	entry.NewRet(funcDecl.Params[0])
}

// Add this new method to handle IO.in_int
func (cg *CodeGenerator) generateIOInIntMethod() {
	funcName := "IO.in_int"
	ioType := cg.TypeMap["IO"]

	// Get the function if it already exists
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		// Create function with signature: Int (IO* self)
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.I32, // Returns Int
			ir.NewParam("self", types.NewPointer(ioType)),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Use scanf instead of custom in_int function
	// Allocate space for the integer result
	resultPtr := entry.NewAlloca(types.I32)

	// Create a global constant for the format string
	formatStrGlobal := cg.Module.NewGlobalDef(".str.scanf_d_fmt", constant.NewCharArrayFromString("%d\x00"))
	formatStrGlobal.Immutable = true

	// Get a pointer to the format string - using a simpler approach
	formatStr := entry.NewBitCast(formatStrGlobal, types.NewPointer(types.I8))

	// Call scanf to read the integer
	entry.NewCall(cg.StdlibFuncs["scanf"], formatStr, resultPtr)

	// Load the result
	result := entry.NewLoad(types.I32, resultPtr)

	// Return the integer
	entry.NewRet(result)
}

// Add this new method to handle IO.in_string
func (cg *CodeGenerator) generateIOInStringMethod() {
	funcName := "IO.in_string"
	ioType := cg.TypeMap["IO"]

	// Get the function if it already exists
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		// Create function with signature: String (IO* self)
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(types.I8), // Returns String (i8*)
			ir.NewParam("self", types.NewPointer(ioType)),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Use scanf instead of custom in_string function
	// Allocate space for the input string (using a fixed buffer size)
	buffer := entry.NewAlloca(types.NewArray(1024, types.I8))

	// Create a global constant for the format string
	formatStrGlobal := cg.Module.NewGlobalDef(".str.scanf_s_fmt", constant.NewCharArrayFromString("%s\x00"))
	formatStrGlobal.Immutable = true

	// Get a pointer to the format string - using a simpler approach
	formatStr := entry.NewBitCast(formatStrGlobal, types.NewPointer(types.I8))

	// Call scanf to read into the buffer
	entry.NewCall(cg.StdlibFuncs["scanf"], formatStr,
		entry.NewBitCast(buffer, types.NewPointer(types.I8)))

	// Allocate heap memory for the string and copy from the buffer
	// First determine the length of the input
	bufferPtr := entry.NewBitCast(buffer, types.NewPointer(types.I8))
	strLen := entry.NewCall(cg.StdlibFuncs["strlen"], bufferPtr)

	// Allocate memory for the string (length + 1 for null terminator)
	allocSize := entry.NewAdd(strLen, constant.NewInt(types.I32, 1))
	mallocResult := entry.NewCall(cg.StdlibFuncs["malloc"],
		entry.NewZExt(allocSize, types.I64))

	// Copy the string to the heap
	entry.NewCall(cg.StdlibFuncs["strcpy"], mallocResult, bufferPtr)

	// Return the string pointer
	entry.NewRet(mallocResult)
}

// Add string method implementations
func (cg *CodeGenerator) generateStringLengthMethod() {
	funcName := "String.length"
	stringType := cg.TypeMap["String"]

	// Get the function if it already exists
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.I32, // Returns Int
			ir.NewParam("self", types.NewPointer(stringType)),
		)
	}

	entry := funcDecl.NewBlock("entry")

	// Check if the self parameter is already a raw string pointer (i8*)
	isRawStringPtr := entry.NewICmp(enum.IPredEQ,
		entry.NewPtrToInt(funcDecl.Params[0], types.I64),
		entry.NewPtrToInt(entry.NewBitCast(funcDecl.Params[0], types.NewPointer(types.I8)), types.I64),
	)

	rawStringBlock := funcDecl.NewBlock("raw_string")
	structStringBlock := funcDecl.NewBlock("struct_string")

	entry.NewCondBr(isRawStringPtr, rawStringBlock, structStringBlock)

	// If it's a raw string pointer, use it directly
	rawStrPtr := rawStringBlock.NewBitCast(funcDecl.Params[0], types.NewPointer(types.I8))
	rawLength := rawStringBlock.NewCall(cg.StdlibFuncs["strlen"], rawStrPtr)
	rawStringBlock.NewRet(rawLength)

	// If it's a proper String struct
	// Instead of using getelementptr, use simple bitcast for compatibility
	// First cast the String* to i8** (pointer to pointer)
	selfAsI8PtrPtr := structStringBlock.NewBitCast(funcDecl.Params[0], types.NewPointer(types.NewPointer(types.I8)))
	// Then load the i8* stored in the struct
	stringPtr := structStringBlock.NewLoad(types.NewPointer(types.I8), selfAsI8PtrPtr)
	// Get the length
	structLength := structStringBlock.NewCall(cg.StdlibFuncs["strlen"], stringPtr)
	structStringBlock.NewRet(structLength)
}

// generateStringConcatMethod generates LLVM IR for String.concat
func (cg *CodeGenerator) generateStringConcatMethod() {
	// Find the function (it should already be declared)
	funcName := "String.concat"
	var concatFunc *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			concatFunc = f
			break
		}
	}

	if concatFunc == nil {
		// If function not found, create it
		stringType := cg.TypeMap["String"]
		concatFunc = cg.Module.NewFunc(
			funcName,
			types.NewPointer(types.I8), // Return type is String
			ir.NewParam("self", types.NewPointer(stringType)),
			ir.NewParam("s", types.NewPointer(types.I8)), // String to concatenate
		)
	}

	// Create entry block
	entry := concatFunc.NewBlock("entry")

	// Check if the self parameter is already a raw string pointer (i8*)
	isRawStringPtr := entry.NewICmp(enum.IPredEQ,
		entry.NewPtrToInt(concatFunc.Params[0], types.I64),
		entry.NewPtrToInt(entry.NewBitCast(concatFunc.Params[0], types.NewPointer(types.I8)), types.I64),
	)

	rawStringBlock := concatFunc.NewBlock("raw_string")
	structStringBlock := concatFunc.NewBlock("struct_string")

	entry.NewCondBr(isRawStringPtr, rawStringBlock, structStringBlock)

	// If it's a raw string pointer, use it directly
	// ---- Raw String Block ----
	// Convert self to a proper i8* if it's not already
	selfStr := rawStringBlock.NewBitCast(concatFunc.Params[0], types.NewPointer(types.I8))
	otherStr := concatFunc.Params[1] // Second parameter is the string to concatenate

	// Get lengths of both strings
	selfLen := rawStringBlock.NewCall(cg.StdlibFuncs["strlen"], selfStr)
	otherLen := rawStringBlock.NewCall(cg.StdlibFuncs["strlen"], otherStr)

	// Calculate the total size needed (selfLen + otherLen + 1 for null terminator)
	totalSize := rawStringBlock.NewAdd(rawStringBlock.NewAdd(selfLen, otherLen), constant.NewInt(types.I32, 1))

	// Allocate memory for the concatenated string
	allocSize := rawStringBlock.NewZExt(totalSize, types.I64)
	resultPtr := rawStringBlock.NewCall(cg.StdlibFuncs["malloc"], allocSize)

	// Use strcpy to copy the first string
	rawStringBlock.NewCall(cg.StdlibFuncs["strcpy"], resultPtr, selfStr)

	// Use strcat to append the second string
	rawStringBlock.NewCall(cg.StdlibFuncs["strcat"], resultPtr, otherStr)

	// Return the concatenated string
	rawStringBlock.NewRet(resultPtr)

	// ---- Struct String Block ----
	// Use bitcast approach for compatibility
	selfAsI8PtrPtr := structStringBlock.NewBitCast(concatFunc.Params[0], types.NewPointer(types.NewPointer(types.I8)))
	selfStrLoad := structStringBlock.NewLoad(types.NewPointer(types.I8), selfAsI8PtrPtr)
	otherStrLoad := concatFunc.Params[1] // Already an i8*

	// Get lengths of both strings
	selfLen = structStringBlock.NewCall(cg.StdlibFuncs["strlen"], selfStrLoad)
	otherLen = structStringBlock.NewCall(cg.StdlibFuncs["strlen"], otherStrLoad)

	// Calculate the total size needed (selfLen + otherLen + 1 for null terminator)
	totalSize = structStringBlock.NewAdd(structStringBlock.NewAdd(selfLen, otherLen), constant.NewInt(types.I32, 1))

	// Allocate memory for the concatenated string
	allocSize = structStringBlock.NewZExt(totalSize, types.I64)
	resultPtr = structStringBlock.NewCall(cg.StdlibFuncs["malloc"], allocSize)

	// Use strcpy to copy the first string
	structStringBlock.NewCall(cg.StdlibFuncs["strcpy"], resultPtr, selfStrLoad)

	// Use strcat to append the second string
	structStringBlock.NewCall(cg.StdlibFuncs["strcat"], resultPtr, otherStrLoad)

	// Return the concatenated string
	structStringBlock.NewRet(resultPtr)
}

// generateStringSubstrMethod generates LLVM IR for String.substr
func (cg *CodeGenerator) generateStringSubstrMethod() {
	// Find the function (it should already be declared)
	funcName := "String.substr"
	var substrFunc *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			substrFunc = f
			break
		}
	}

	if substrFunc == nil {
		// If function not found, create it
		stringType := cg.TypeMap["String"]
		substrFunc = cg.Module.NewFunc(
			funcName,
			types.NewPointer(types.I8), // Return type is String
			ir.NewParam("self", types.NewPointer(stringType)),
			ir.NewParam("i", types.I32), // Starting index
			ir.NewParam("l", types.I32), // Length
		)
	}

	// Create the entry block
	entry := substrFunc.NewBlock("entry")

	// Check if the self parameter is already a raw string pointer (i8*)
	isRawStringPtr := entry.NewICmp(enum.IPredEQ,
		entry.NewPtrToInt(substrFunc.Params[0], types.I64),
		entry.NewPtrToInt(entry.NewBitCast(substrFunc.Params[0], types.NewPointer(types.I8)), types.I64),
	)

	rawStringBlock := substrFunc.NewBlock("raw_string")
	structStringBlock := substrFunc.NewBlock("struct_string")

	entry.NewCondBr(isRawStringPtr, rawStringBlock, structStringBlock)

	// ---- Raw String Block (for i8* input) ----
	selfStr := rawStringBlock.NewBitCast(substrFunc.Params[0], types.NewPointer(types.I8))
	startIdx := substrFunc.Params[1] // Starting index
	length := substrFunc.Params[2]   // Length to extract

	// Get the length of the original string
	strLen := rawStringBlock.NewCall(cg.StdlibFuncs["strlen"], selfStr)

	// Split into bounds_check block for raw string
	boundsCheckRawBlock := substrFunc.NewBlock("bounds_check_raw")
	rawStringBlock.NewBr(boundsCheckRawBlock)

	// Error block shared by both paths
	errorBlock := substrFunc.NewBlock("error")

	// Allocation block for raw string
	allocRawBlock := substrFunc.NewBlock("alloc_raw")

	// Check if start < 0 or start >= strLen or length < 0
	startOutOfBoundsRaw := boundsCheckRawBlock.NewOr(
		boundsCheckRawBlock.NewICmp(enum.IPredSLT, startIdx, constant.NewInt(types.I32, 0)),
		boundsCheckRawBlock.NewICmp(enum.IPredSGE, startIdx, strLen),
	)
	lengthNegativeRaw := boundsCheckRawBlock.NewICmp(enum.IPredSLT, length, constant.NewInt(types.I32, 0))
	invalidInputRaw := boundsCheckRawBlock.NewOr(startOutOfBoundsRaw, lengthNegativeRaw)

	// Branch based on the bounds check
	boundsCheckRawBlock.NewCondBr(invalidInputRaw, errorBlock, allocRawBlock)

	// Allocate memory for the new substring
	allocSizeRaw := allocRawBlock.NewAdd(length, constant.NewInt(types.I32, 1))
	mallocCallRaw := allocRawBlock.NewCall(cg.StdlibFuncs["malloc"],
		allocRawBlock.NewZExt(allocSizeRaw, types.I64))

	// Get pointer to the start of the substring in the original string
	substrStartRaw := allocRawBlock.NewGetElementPtr(types.I8, selfStr, startIdx)

	// Use strncpy to copy exactly 'length' characters
	allocRawBlock.NewCall(cg.StdlibFuncs["strncpy"], mallocCallRaw, substrStartRaw, length)

	// Add null terminator at the end of the new substring
	nullTermPtrRaw := allocRawBlock.NewGetElementPtr(types.I8, mallocCallRaw, length)
	allocRawBlock.NewStore(constant.NewInt(types.I8, 0), nullTermPtrRaw)

	// Return the raw substring pointer
	allocRawBlock.NewRet(mallocCallRaw)

	// ---- Struct String Block (for String* input) ----
	// Use bitcast approach for compatibility
	selfAsI8PtrPtr := structStringBlock.NewBitCast(substrFunc.Params[0], types.NewPointer(types.NewPointer(types.I8)))
	loadedStr := structStringBlock.NewLoad(types.NewPointer(types.I8), selfAsI8PtrPtr)

	// Get the parameters
	startIdxStruct := substrFunc.Params[1] // Starting index
	lengthStruct := substrFunc.Params[2]   // Length to extract

	// Get the length of the original string
	strLenStruct := structStringBlock.NewCall(cg.StdlibFuncs["strlen"], loadedStr)

	// Split into bounds_check block for struct
	boundsCheckStructBlock := substrFunc.NewBlock("bounds_check_struct")
	structStringBlock.NewBr(boundsCheckStructBlock)

	// Allocation block for struct string
	allocStructBlock := substrFunc.NewBlock("alloc_struct")

	// Check if start < 0 or start >= strLen or length < 0
	startOutOfBoundsStruct := boundsCheckStructBlock.NewOr(
		boundsCheckStructBlock.NewICmp(enum.IPredSLT, startIdxStruct, constant.NewInt(types.I32, 0)),
		boundsCheckStructBlock.NewICmp(enum.IPredSGE, startIdxStruct, strLenStruct),
	)
	lengthNegativeStruct := boundsCheckStructBlock.NewICmp(enum.IPredSLT, lengthStruct, constant.NewInt(types.I32, 0))
	invalidInputStruct := boundsCheckStructBlock.NewOr(startOutOfBoundsStruct, lengthNegativeStruct)

	// Branch based on the bounds check
	boundsCheckStructBlock.NewCondBr(invalidInputStruct, errorBlock, allocStructBlock)

	// Setup error block
	errorMsgGlobal := cg.Module.NewGlobalDef(".str.substr_error",
		constant.NewCharArrayFromString("Runtime error: substring out of range\n\x00"))
	errorMsgGlobal.Immutable = true

	errorMsgPtr := errorBlock.NewBitCast(errorMsgGlobal, types.NewPointer(types.I8))

	errorBlock.NewCall(cg.StdlibFuncs["printf"], errorMsgPtr)
	errorBlock.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
	errorBlock.NewUnreachable()

	// Allocate memory for the new substring
	allocSizeStruct := allocStructBlock.NewAdd(lengthStruct, constant.NewInt(types.I32, 1))
	mallocCallStruct := allocStructBlock.NewCall(cg.StdlibFuncs["malloc"],
		allocStructBlock.NewZExt(allocSizeStruct, types.I64))

	// Get pointer to the start of the substring in the original string
	substrStartStruct := allocStructBlock.NewGetElementPtr(types.I8, loadedStr, startIdxStruct)

	// Use strncpy to copy exactly 'length' characters
	allocStructBlock.NewCall(cg.StdlibFuncs["strncpy"], mallocCallStruct, substrStartStruct, lengthStruct)

	// Add null terminator at the end of the new substring
	nullTermPtrStruct := allocStructBlock.NewGetElementPtr(types.I8, mallocCallStruct, lengthStruct)
	allocStructBlock.NewStore(constant.NewInt(types.I8, 0), nullTermPtrStruct)

	// Return the raw substring pointer
	allocStructBlock.NewRet(mallocCallStruct)
}

// isStringPointer checks if a type is a pointer to i8 (string type in LLVM) or a %String* type
func isStringPointer(t types.Type) bool {
	// Check if it's directly an i8 pointer
	if t == types.NewPointer(types.I8) {
		return true
	}

	// Check if it's a pointer type pointing to i8
	if ptrType, ok := t.(*types.PointerType); ok {
		if ptrType.ElemType == types.I8 {
			return true
		}

		// Check if it's a %String* type (pointer to String struct)
		if structType, ok := ptrType.ElemType.(*types.StructType); ok {
			if structType.Name() == "String" {
				return true
			}
		}
	}

	return false
}

// isStringObject checks if a type is specifically a %String* (not just an i8*)
func isStringObject(t types.Type, cg *CodeGenerator) bool {
	if ptrType, ok := t.(*types.PointerType); ok {
		stringType, exists := cg.TypeMap["String"]
		if !exists {
			return false
		}
		return ptrType.ElemType == stringType
	}
	return false
}

// isIntType checks if a type is an i32 (primitive int) or %Int* (Int object)
func isIntType(t types.Type, cg *CodeGenerator) bool {
	if t == types.I32 {
		return true
	}

	if ptrType, ok := t.(*types.PointerType); ok {
		intType, exists := cg.TypeMap["Int"]
		if !exists {
			return false
		}
		return ptrType.ElemType == intType
	}

	return false
}

// isBoolType checks if a type is an i1 (primitive bool) or %Bool* (Bool object)
func isBoolType(t types.Type, cg *CodeGenerator) bool {
	if t == types.I1 {
		return true
	}

	if ptrType, ok := t.(*types.PointerType); ok {
		boolType, exists := cg.TypeMap["Bool"]
		if !exists {
			return false
		}
		return ptrType.ElemType == boolType
	}

	return false
}

// ensureCorrectArgumentType checks if the argument type matches the expected parameter type
// and performs necessary type conversions if needed
func (cg *CodeGenerator) ensureCorrectArgumentType(arg value.Value, paramType types.Type) value.Value {
	block := cg.CurrentBlock

	// If types already match, no conversion needed
	if arg.Type() == paramType {
		return arg
	}

	// Handle conversions between primitive types and their object representations

	// Check if we're trying to convert between a primitive type and a pointer type
	_, argIsPtr := arg.Type().(*types.PointerType)
	_, paramIsPtr := paramType.(*types.PointerType)

	// Converting from a primitive type to a pointer type
	if !argIsPtr && paramIsPtr {
		// For integer or boolean to pointer conversion
		if arg.Type() == types.I32 || arg.Type() == types.I1 {
			// Use inttoptr for primitive to pointer conversion
			tmpPtr := block.NewIntToPtr(arg, types.NewPointer(types.I8))

			// If we need a specific pointer type, cast to it
			if paramType != types.NewPointer(types.I8) {
				return block.NewBitCast(tmpPtr, paramType)
			}
			return tmpPtr
		}

		// For other primitive types, create a boxed object
		if arg.Type() == types.I32 {
			// Create an Int object to wrap the int value
			intObj := cg.generateObjectAllocation("Int")
			// In a real implementation, we would store the value in the object
			return intObj
		} else if arg.Type() == types.I1 {
			// Create a Bool object to wrap the boolean value
			boolObj := cg.generateObjectAllocation("Bool")
			// In a real implementation, we would store the value in the object
			return boolObj
		}
	}

	// Converting from a pointer type to a primitive type
	if argIsPtr && !paramIsPtr {
		// For pointer to integer or boolean conversion
		if paramType == types.I32 || paramType == types.I1 {
			// Use ptrtoint for pointer to primitive conversion
			return block.NewPtrToInt(arg, paramType)
		}
	}

	// String conversion cases
	if isStringPointer(paramType) {
		// For String methods that expect %String* but getting i8*
		if arg.Type() == types.NewPointer(types.I8) && isStringObject(paramType, cg) {
			// Create a String object or wrap the string pointer if needed
			// This is a simplified approach - ideally we'd create a proper String object
			return block.NewBitCast(arg, paramType)
		} else if isStringObject(arg.Type(), cg) && paramType == types.NewPointer(types.I8) {
			// For String methods that expect i8* but getting %String*
			// Extract the string pointer from the String object
			// This is a simplified approach - ideally we'd access the string buffer inside the String object
			return block.NewBitCast(arg, paramType)
		} else if isStringPointer(arg.Type()) {
			// Handle any other string pointer casting needs
			return block.NewBitCast(arg, paramType)
		}
	}

	// Generic pointer casting for object types
	if paramIsPtr && argIsPtr {
		// Cast between pointer types
		return block.NewBitCast(arg, paramType)
	}

	// If we can't figure out a better conversion, try a final approach based on type
	if paramIsPtr {
		// If parameter type is a pointer but arg is not, and we didn't handle it above
		// Try to box the primitive value into an object
		if arg.Type() == types.I32 {
			intObj := cg.generateObjectAllocation("Int")
			return block.NewBitCast(intObj, paramType)
		} else if arg.Type() == types.I1 {
			boolObj := cg.generateObjectAllocation("Bool")
			return block.NewBitCast(boolObj, paramType)
		}
		// Last resort for primitives to pointers - use inttoptr as a generic approach
		if !argIsPtr && (arg.Type() == types.I32 || arg.Type() == types.I1) {
			tmpPtr := block.NewIntToPtr(arg, types.NewPointer(types.I8))
			return block.NewBitCast(tmpPtr, paramType)
		}
	}

	// We've tried everything, just attempt a bitcast and hope for the best
	// (this may fail during LLVM validation if types are incompatible)
	return block.NewBitCast(arg, paramType)
}

// canCastTo checks if an object can be cast to the specified class type
// This checks the class hierarchy to determine if the cast is valid
func (cg *CodeGenerator) canCastTo(obj value.Value, targetClassName string) bool {
	// Get the runtime type of the object
	runtimeType := cg.getObjectRuntimeType(obj)
	if runtimeType == targetClassName {
		// Direct match, can cast
		return true
	}

	// Check class hierarchy
	currentType := runtimeType
	for {
		parent, exists := cg.ClassHierarchy[currentType]
		if !exists || parent == "" {
			// Reached Object or unknown class
			return false
		}

		if parent == targetClassName {
			// Found target class in ancestry, can cast
			return true
		}

		currentType = parent
	}
}
