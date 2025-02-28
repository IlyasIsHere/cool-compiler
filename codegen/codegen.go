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

	// Collect methods from parent classes first (to maintain correct override order)
	var parentMethods map[string]*ir.Func
	if parent, exists := cg.ClassHierarchy[className]; exists && parent != "" {
		// Make sure parent vtable is created first
		parentClass := findClass(parent, program)
		if parentClass != nil {
			cg.createVTableForClass(parentClass, program)
		}

		// Copy methods from parent vtable
		parentMethods = make(map[string]*ir.Func)
		// In a full implementation, you would copy methods from parent vtable here

		// Copy the parent method indices to this class
		for methodName, methodIndex := range cg.MethodIndices[parent] {
			cg.MethodIndices[className][methodName] = methodIndex
		}
	}

	// Create a map of method names to functions for this class
	methods := make(map[string]*ir.Func)
	if parentMethods != nil {
		// Copy parent methods (which will be overridden by this class's methods if same name)
		for name, method := range parentMethods {
			methods[name] = method
		}
	}

	// Process methods of this class
	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod {
			// Check if this is a built-in method that's already handled
			if isBuiltInMethod(className, method.Name.Value) {
				// Use existing function declaration
				funcName := fmt.Sprintf("%s.%s", className, method.Name.Value)
				var found bool
				for _, f := range cg.Module.Funcs {
					if f.Name() == funcName {
						methods[method.Name.Value] = f
						found = true
						break
					}
				}

				// If we didn't find the built-in method, create it on demand
				if !found && className == "IO" {
					switch method.Name.Value {
					case "out_string":
						// Create IO.out_string function
						outStringFunc := cg.Module.NewFunc(
							"IO.out_string",
							types.NewPointer(cg.TypeMap["IO"]),
							ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
							ir.NewParam("str", types.NewPointer(types.I8)),
						)
						methods[method.Name.Value] = outStringFunc
					case "out_int":
						// Create IO.out_int function
						outIntFunc := cg.Module.NewFunc(
							"IO.out_int",
							types.NewPointer(cg.TypeMap["IO"]),
							ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
							ir.NewParam("i", types.I32),
						)
						methods[method.Name.Value] = outIntFunc
					case "in_string":
						// Create IO.in_string function
						inStringFunc := cg.Module.NewFunc(
							"IO.in_string",
							types.NewPointer(types.I8),
							ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						)
						methods[method.Name.Value] = inStringFunc
					case "in_int":
						// Create IO.in_int function
						inIntFunc := cg.Module.NewFunc(
							"IO.in_int",
							types.I32,
							ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						)
						methods[method.Name.Value] = inIntFunc
					}
				}

				continue // Skip creating a new declaration
			}

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
			funcName := fmt.Sprintf("%s.%s", className, method.Name.Value)

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

			// Add to methods map
			methods[method.Name.Value] = function
		}
	}

	// Sort method names for consistent ordering
	var methodNames []string
	for name := range methods {
		methodNames = append(methodNames, name)
	}
	sort.Strings(methodNames)

	// Create an array of function pointers for the vtable
	methodCount := len(methods)
	vtableType := types.NewArray(uint64(methodCount), types.NewPointer(types.I8))

	// Create global array with proper initialization
	vtableName := fmt.Sprintf("vtable.%s", className)

	// Create the initializers for the vtable
	initializers := make([]constant.Constant, methodCount)
	for i, name := range methodNames {
		funcPtr := methods[name]
		// Cast the function pointer to i8*
		initializers[i] = constant.NewBitCast(funcPtr, types.NewPointer(types.I8))

		// Store the method index in the MethodIndices map
		cg.MethodIndices[className][name] = i
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
			cg.generateObjectAbortMethod(class, method)
			return
		case "type_name":
			cg.generateTypeNameMethod(class, method)
			return
		case "copy":
			cg.generateCopyMethod(class, method)
			return
		}
	}

	if className == "IO" {
		switch methodName {
		case "out_string":
			cg.generateIOOutStringMethod(class, method)
			return
		case "out_int":
			cg.generateIOOutIntMethod(class, method)
			return
		case "in_string":
			cg.generateIOInStringMethod(class, method)
			return
		case "in_int":
			cg.generateIOInIntMethod(class, method)
			return
		}
	}

	// Handle String class methods
	if className == "String" {
		switch methodName {
		case "length":
			cg.generateStringLengthMethod(class, method)
			return
		case "concat":
			cg.generateStringConcatMethod(class, method)
			return
		case "substr":
			cg.generateStringSubstrMethod(class, method)
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

	// The method body should compute a value of the correct return type
	if bodyValue.Type() != methodFunc.Sig.RetType {
		// If types don't match, we may need to cast or handle special cases
		if ptr, isPtrType := methodFunc.Sig.RetType.(*types.PointerType); isPtrType {
			if bodyValue.Type() == types.I32 || bodyValue.Type() == types.I1 {
				// Boxing primitive values when returning as Object
				// In a real implementation, would create proper boxed objects
				bodyValue = cg.CurrentBlock.NewIntToPtr(bodyValue, ptr)
			} else if bodyValue.Type() == types.Void {
				// Handle void return values when pointer expected (common with IO operations)
				// Return 'self' as a sensible default
				bodyValue = cg.CurrentBlock.NewBitCast(selfParam, ptr)
			} else if _, isOtherPtr := bodyValue.Type().(*types.PointerType); isOtherPtr {
				// Cast between pointer types
				bodyValue = cg.CurrentBlock.NewBitCast(bodyValue, ptr)
			}
		} else if methodFunc.Sig.RetType == types.Void && bodyValue.Type() != types.Void {
			// If we need to return void but have a value, just ignore the value
			cg.CurrentBlock.NewRet(nil)
			return
		} else {
			// added this to handle the case where the return type is not a pointer
			bodyValue = cg.CurrentBlock.NewBitCast(bodyValue, methodFunc.Sig.RetType)
		}
	}

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

// generateObjectAllocation creates a new instance of a class
func (cg *CodeGenerator) generateObjectAllocation(typeName string) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Get the LLVM struct type for the class
	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("attempt to create an instance of unknown type: %s", typeName))
	}

	// Calculate the size of the class struct using getelementptr
	sizeGEP := constant.NewGetElementPtr(
		classType,
		constant.NewNull(types.NewPointer(classType)),
		constant.NewInt(types.I32, 1),
	)

	// Call malloc with the size of the class
	mallocFunc, exists := cg.StdlibFuncs["malloc"]
	if !exists {
		panic("malloc function not found")
	}

	// malloc returns i8* which we'll cast to the appropriate type
	mallocCall := block.NewCall(mallocFunc, sizeGEP)

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
				if attr.Expression != nil {
					// Get the attribute index
					attrIndex, exists := cg.AttributeIndices[className][attr.Name.Value]
					if !exists {
						continue // Skip if the attribute index isn't found
					}

					// Generate the init expression
					initValue := cg.generateExpression(attr.Expression)

					// Get the attribute type
					attrType := classType.Fields[attrIndex]

					// Make sure types match
					if !initValue.Type().Equal(attrType) {
						// Need to cast if types don't match
						initValue = block.NewBitCast(initValue, attrType)
					}

					// Get a pointer to the attribute field and store the initial value
					attrPtr := block.NewGetElementPtr(
						classType,
						objectPtr,
						constant.NewInt(types.I32, 0),
						constant.NewInt(types.I32, int64(attrIndex)),
					)
					block.NewStore(initValue, attrPtr)
				}
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

// getObjectRuntimeType gets the runtime type name of an object
// In a real compiler, this would use runtime type information
// For our simple implementation, we'll extract it from the object's type
func (cg *CodeGenerator) getObjectRuntimeType(object value.Value, block *ir.Block) string {
	// Handle primitive types first
	switch object.Type() {
	case types.I1: // Boolean type
		return "Bool"
	case types.I32: // Integer type
		return "Int"
	case types.I8Ptr: // String type (i8*)
		return "String"
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

	// First, we need to load the vtable pointer from the object
	// In our implementation, the vtable pointer is always the first field of any object

	// Get the object's type
	objPtrType, ok := object.Type().(*types.PointerType)
	if !ok {
		panic(fmt.Sprintf("expected object to be a pointer type, got: %v", object.Type()))
	}

	// Get the underlying struct type
	objStructType, ok := objPtrType.ElemType.(*types.StructType)
	if !ok {
		// Handle i8* pointer type - we need to bitcast it to the appropriate type first
		// Get the object's runtime type name
		objectTypeName := cg.getObjectRuntimeType(object, block)

		// Look up the struct type for this class
		structType, exists := cg.TypeMap[objectTypeName]
		if !exists {
			panic(fmt.Sprintf("cannot find struct type for class: %s", objectTypeName))
		}

		// Bitcast the i8* pointer to the appropriate struct pointer type
		structPtrType := types.NewPointer(structType)
		object = block.NewBitCast(object, structPtrType)

		// Update the struct type for further operations
		objStructType = structType
	}

	// Get pointer to the vtable field in the object (first field)
	vtablePtrPtr := block.NewGetElementPtr(
		objStructType,
		object,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	// Load the vtable pointer
	vtablePtr := block.NewLoad(types.NewPointer(types.I8), vtablePtrPtr)

	// Find the method's index in the vtable
	// Get the object's runtime type
	objectTypeName := cg.getObjectRuntimeType(object, block)
	methodIndex := -1

	// Look up the method index in the MethodIndices map
	if indices, exists := cg.MethodIndices[objectTypeName]; exists {
		if idx, exists := indices[methodName]; exists {
			methodIndex = idx
		}

	}

	if methodIndex == -1 {
		panic(fmt.Sprintf("method %s not found in class %s", methodName, objectTypeName))
	}

	// Load the function pointer from the vtable
	// In LLVM IR, this involves a getelementptr to get to the right slot, then a load
	gepIndices := []value.Value{
		constant.NewInt(types.I32, 0),                  // First index is always 0 for struct GEP
		constant.NewInt(types.I32, int64(methodIndex)), // Second index is the method index
	}

	// We need to cast vtablePtr to the appropriate type first
	vtableType := types.NewPointer(types.NewArray(0, types.NewPointer(types.I8)))
	castedVTablePtr := block.NewBitCast(vtablePtr, vtableType)

	// Get pointer to the method slot in the vtable
	methodSlotPtr := block.NewGetElementPtr(
		types.NewArray(0, types.NewPointer(types.I8)),
		castedVTablePtr,
		gepIndices...,
	)

	// Load the method function pointer
	methodPtr := block.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	// For a real implementation, we should look up the method's signature
	// Determine appropriate return type for the method
	var returnType types.Type
	// For special IO methods like out_string that return SELF_TYPE
	if methodName == "out_string" || methodName == "out_int" {
		// These methods return the object itself (SELF_TYPE)
		returnType = object.Type()
	} else {
		// Look up method's actual return type by finding the method declaration
		// First, determine object's class type from its LLVM type
		className := objectTypeName

		if className != "" {
			// Look for the function declaration in the module
			funcName := fmt.Sprintf("%s.%s", className, methodName)
			var methodFunc *ir.Func

			// Search for method in the module
			for _, f := range cg.Module.Funcs {
				if f.Name() == funcName {
					methodFunc = f
					break
				}
			}

			if methodFunc != nil {
				// Use the actual return type from the method's signature
				returnType = methodFunc.Sig.RetType
			} else {
				// If method not found, use the generic Object pointer type as fallback
				returnType = types.NewPointer(types.I8)
			}
		} else {
			// If we can't determine the class name, use generic Object pointer
			returnType = types.NewPointer(types.I8)
		}
	}

	// Create a function type with the appropriate parameters and return type
	paramTypes := make([]types.Type, len(args))
	for i, arg := range args {
		paramTypes[i] = arg.Type()
	}
	funcType := types.NewPointer(types.NewFunc(returnType, paramTypes...))

	// Cast the i8* function pointer to the correct function type
	castedMethodPtr := block.NewBitCast(methodPtr, funcType)

	// Call the method with the provided arguments
	call := block.NewCall(castedMethodPtr, args...)

	// For methods that return SELF_TYPE (like out_string), the result should be the object itself
	if methodName == "out_string" || methodName == "out_int" {
		// If the method returns SELF_TYPE, return the object itself
		return object
	}

	return call
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

	// Call the Main.main() method
	// First, find the main method in the vtable of Main
	vtable, exists := cg.VTables["Main"]
	if !exists {
		panic("Main class must have a vtable")
	}

	// We need to find the index of the main method in the vtable
	// In a real implementation, we would have a mapping from method names to indices
	// For simplicity, we'll assume we can look up the method directly
	mainMethodIndex := 0 // This should be the actual index of main in the vtable

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
	cg.CurrentBlock.NewBr(mergeBlock)

	// Generate code for the false branch
	cg.CurrentBlock = falseBlock
	falseValue := cg.generateExpression(ifExpr.Alternative)
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
		// For simplicity, use i8* as a generic object pointer type
		// In a full implementation, you would calculate the least common ancestor type
		resultType = types.NewPointer(types.I8)
	}

	// Create a PHI node with incoming values right away
	phi := cg.CurrentBlock.NewPhi(
		&ir.Incoming{X: trueValue, Pred: trueBlock},
		&ir.Incoming{X: falseValue, Pred: falseBlock},
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

	// Get the current block and create a branch to the condition block
	currentBlock := cg.CurrentBlock
	currentBlock.NewBr(condBlock)

	// Set current block to condition block
	cg.CurrentBlock = condBlock

	// Generate code for the condition
	condValue := cg.generateExpression(whileExpr.Condition)

	// Check if condition is a boolean
	if condValue.Type() != types.I1 {
		panic(fmt.Sprintf("condition in while expression must be of boolean type"))
	}

	// Create conditional branch: if condition is true, enter body, otherwise exit
	condBlock.NewCondBr(condValue, bodyBlock, exitBlock)

	// Generate code for the loop body
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-2] = bodyBlock

	// Generate the loop body expression
	cg.generateExpression(whileExpr.Body)

	// After executing the body, jump back to the condition block
	bodyBlock = cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-2] // Get updated body block
	bodyBlock.NewBr(condBlock)

	// Set the current block to the exit block
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1] = exitBlock

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

	// Generate LLVM values for all arguments
	argValues := make([]value.Value, 0, len(dotCall.Arguments)+1)

	// The first argument to a method call is always the object itself (self)
	argValues = append(argValues, objectValue)

	// Add the rest of the arguments
	for _, arg := range dotCall.Arguments {
		argValues = append(argValues, cg.generateExpression(arg))
	}

	// Special case for IO.out_string and IO.out_int to call runtime functions directly
	// This avoids the vtable dispatch which is causing issues
	if dotCall.Method.Value == "out_string" || dotCall.Method.Value == "out_int" {
		// First check if the object is an IO object
		objType := objectValue.Type().(*types.PointerType).ElemType
		if structType, isStruct := objType.(*types.StructType); isStruct && structType.Name() == "IO" {
			if dotCall.Method.Value == "out_string" && len(argValues) > 1 {
				// Find or create the runtime function
				var outStringFunc *ir.Func
				// Look for an existing declaration
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_string" { // Corrected function name
						outStringFunc = f
						break
					}
				}

				if outStringFunc == nil {
					// If the function isn't already declared, declare it
					outStringFunc = cg.Module.NewFunc("IO.out_string", types.NewPointer(cg.TypeMap["IO"]), ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])), ir.NewParam("str", types.NewPointer(types.I8))) // Corrected signature
				}

				// Make the call with the self and string argument
				block.NewCall(outStringFunc, argValues...)

				// Return the IO object itself (this is what COOL's out_string does)
				return objectValue

			} else if dotCall.Method.Value == "out_int" && len(argValues) > 1 {
				// Get the int argument (skip self which is at index 0)

				// Find or create the runtime function
				var outIntFunc *ir.Func
				// Look for an existing declaration
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_int" { // Corrected function name
						outIntFunc = f
						break
					}
				}

				if outIntFunc == nil {
					// If the function isn't already declared, declare it
					outIntFunc = cg.Module.NewFunc("IO.out_int", types.NewPointer(cg.TypeMap["IO"]), ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])), ir.NewParam("n", types.I32)) // Corrected signature
				}

				// Make the call with self and the int argument
				block.NewCall(outIntFunc, argValues...)

				// Return the IO object itself (this is what COOL's out_int does)
				return objectValue
			}
		}
	}

	// Determine if this is a dynamic or static dispatch
	var result value.Value
	if dotCall.Type != nil {
		// This is a static dispatch (e.g., obj@Type.method())
		result = cg.generateStaticDispatch(objectValue, dotCall.Type.Value, dotCall.Method.Value, argValues)
	} else {
		// This is a dynamic dispatch (e.g., obj.method())
		result = cg.generateDynamicDispatch(objectValue, dotCall.Method.Value, argValues)
	}

	// Special case for IO methods like out_string that return SELF_TYPE
	// In COOL, a method with return type SELF_TYPE returns the object itself
	methodName := dotCall.Method.Value
	if methodName == "out_string" || methodName == "out_int" {
		// For IO methods that return the object itself (SELF_TYPE)
		// Return the object (self) instead of the void result from the call

		// Make sure the return value has the expected type based on the current function's return type
		if cg.CurrentFunc != nil && result.Type() != cg.CurrentFunc.Sig.RetType {
			if _, isPtrType := cg.CurrentFunc.Sig.RetType.(*types.PointerType); isPtrType {
				// If we need to return a pointer type (like Object), convert void or other types
				if result.Type() == types.Void {
					// Convert object value to the expected return type if it's a pointer
					result = block.NewBitCast(objectValue, cg.CurrentFunc.Sig.RetType)
				} else if _, isOtherPtr := result.Type().(*types.PointerType); isOtherPtr {
					// Cast between pointer types
					result = block.NewBitCast(result, cg.CurrentFunc.Sig.RetType)
				}
			}
		}
	}

	return result
}

// generateStaticDispatch creates LLVM IR for static dispatch (obj@Type.method())
func (cg *CodeGenerator) generateStaticDispatch(object value.Value, typeName string, methodName string, args []value.Value) value.Value {
	block := cg.CurrentBlock

	// Look up the method in the specified class's vtable
	// In a static dispatch, we bypass dynamic dispatch and directly call the method of the specified type

	// First, get the vtable for the specified type
	vtable, exists := cg.VTables[typeName]
	if !exists {
		panic(fmt.Sprintf("unknown type in static dispatch: %s", typeName))
	}

	// Get the global vtable value
	vtablePtr := vtable

	// Find the method's index in the vtable
	// This would require a mapping from method names to vtable indices
	// For simplicity, we'll assume there's a helper function that finds the method's index
	methodIndex := 0 // Placeholder - would need actual implementation

	// Load the function pointer from the vtable
	// In LLVM IR, this involves a getelementptr to get to the right slot, then a load
	gepIndices := []value.Value{
		constant.NewInt(types.I32, 0),                  // First index is always 0 for struct GEP
		constant.NewInt(types.I32, int64(methodIndex)), // Second index is the method index
	}

	// Get pointer to the method slot in the vtable
	methodSlotPtr := block.NewGetElementPtr(vtablePtr.ContentType, vtablePtr, gepIndices...)

	// Load the method function pointer
	methodPtr := block.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	// Cast the i8* function pointer to the correct function type
	// For a real implementation, we would need to know the function's signature
	funcType := types.NewPointer(types.NewFunc(types.Void)) // Placeholder - need actual signature
	castedMethodPtr := block.NewBitCast(methodPtr, funcType)

	// Call the method with the provided arguments
	call := block.NewCall(castedMethodPtr, args...)

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
				// For strings, use an empty string
				emptyStr := constant.NewCharArrayFromString("\x00")
				global := cg.Module.NewGlobalDef(".str.empty", emptyStr)
				global.Immutable = true
				// Get a pointer to the first character
				initValue = constant.NewGetElementPtr(
					emptyStr.Type(),
					global,
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
		// Create a new empty string
		emptyStr := constant.NewCharArrayFromString("\x00")
		global := cg.Module.NewGlobalDef(".str.empty", emptyStr)
		global.Immutable = true
		// Get a pointer to the first character
		return constant.NewGetElementPtr(
			emptyStr.Type(),
			global,
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
	currentBlock := currentFunc.Blocks[len(currentFunc.Blocks)-(len(caseExpr.Branches)+2)]

	// In a real implementation, we would check the actual runtime type of the object
	// This involves using RTTI (Runtime Type Information) stored in the object

	// For simplicity in this implementation, we'll assume we can extract the type information
	// and branch based on it. In a real compiler, we would likely have a helper function to do this.

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
	}

	// For simplicity, we'll use a simple implementation that just checks each type in sequence
	// In a real compiler, we might use a more efficient dispatch mechanism like a jump table

	// Generate code to branch to the appropriate block based on the runtime type
	// Since we don't have actual runtime type checking here, we'll simulate it with a simple chain

	// Keep track of the current branching block
	branchingBlock := currentBlock

	// Process each branch in reverse order (COOL case semantics check branches in order listed)
	for i := range caseExpr.Branches {
		// For each branch, we check if the object's type matches the branch's type
		// For simplicity, we'll just branch to each block in sequence
		// In a real compiler, we would do actual runtime type checking here

		if i < len(caseExpr.Branches)-1 {
			// Create a chained comparison for all but the last branch
			// We would check the actual type here, but for simplicity we'll just always go to the branch
			// This is a placeholder for real type checking logic
			branchingBlock.NewBr(branchBlocks[i])

			// Move to the next branch block
			branchingBlock = branchBlocks[i]
		} else {
			// The last branch is taken if none of the previous branches matched
			// In a well-typed COOL program, one branch must match, so this is a catch-all
			branchingBlock.NewBr(branchBlocks[i])
		}
	}

	// Add code to handle the case where no branch matches (a runtime error in COOL)
	// This should never happen in well-typed COOL programs
	noMatchBlock.NewUnreachable()

	// Generate code for each branch
	branchValues := make([]value.Value, len(caseExpr.Branches))
	branchEndBlocks := make([]*ir.Block, len(caseExpr.Branches))

	for i, branch := range caseExpr.Branches {
		// Set the current block to the branch block
		currentFunc.Blocks[len(currentFunc.Blocks)-(len(caseExpr.Branches)+2-i)] = branchBlocks[i]

		// In a real implementation, we would create a new variable bound to the object
		// For the branch's scope, with the branch's declared type
		// For simplicity, we'll skip this and just generate code for the branch expression

		// To simulate this, we would save the old symbol table, add the branch variable,
		// generate the branch code, and then restore the old symbol table

		// Save old symbol table
		oldSymbols := make(map[string]value.Value)
		for k, v := range cg.Symbols {
			oldSymbols[k] = v
		}

		// Add the branch variable to the symbol table
		// In a real implementation, we would cast the object to the branch's type
		cg.Symbols[branch.Identifier.Value] = exprValue

		// Generate code for the branch expression
		branchValues[i] = cg.generateExpression(branch.Expression)

		// Restore the old symbol table
		cg.Symbols = oldSymbols

		// Get the current block after generating the branch expression
		branchEndBlocks[i] = currentFunc.Blocks[len(currentFunc.Blocks)-(len(caseExpr.Branches)+2-i)]

		// Branch to the end block
		branchEndBlocks[i].NewBr(endBlock)
	}

	// Set the current block to the end block
	currentFunc.Blocks[len(currentFunc.Blocks)-1] = endBlock

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

	// Create the PHI node
	phi := endBlock.NewPhi()
	phi.Typ = resultType

	// Add incoming values for the PHI node
	for i, val := range branchValues {
		// If the branch value type doesn't match the result type, cast it
		if !val.Type().Equal(resultType) {
			// Cast to the common type
			// This would require custom casting logic in a real implementation
			// For simplicity, we'll use a bitcast for pointer types
			if _, isResultPtr := resultType.(*types.PointerType); isResultPtr {
				if _, isValPtr := val.Type().(*types.PointerType); isValPtr {
					// Bitcast between pointer types
					val = branchEndBlocks[i].NewBitCast(val, resultType)
				} else {
					// Cast from non-pointer to pointer not handled here
					panic(fmt.Sprintf("cannot cast from %v to %v", val.Type(), resultType))
				}
			} else {
				// Other casts not handled here
				panic(fmt.Sprintf("cannot cast from %v to %v", val.Type(), resultType))
			}
		}

		phi.Incs = append(phi.Incs, &ir.Incoming{X: val, Pred: branchEndBlocks[i]})
	}

	// Return the PHI node as the result of the case expression
	return phi
}

// generateCallExpression creates LLVM IR for function calls
func (cg *CodeGenerator) generateCallExpression(callExpr *ast.CallExpression) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Generate code for each argument
	args := make([]value.Value, 0, len(callExpr.Arguments))
	for _, arg := range callExpr.Arguments {
		args = append(args, cg.generateExpression(arg))
	}

	// Generate code for the function expression
	functionExpr := cg.generateExpression(callExpr.Function)

	// This is a simplified implementation that assumes function calls are
	// either to object identifiers (direct calls) or dot calls that have
	// already been processed.

	// In a real implementation, you would need to:
	// 1. Check if this is a method call (via a dot expression)
	// 2. Determine if it's a static or dynamic dispatch
	// 3. Handle the case where the function itself is an expression

	// For now, we'll assume this is a direct call to a function
	// whose value we've evaluated to functionExpr

	// We need to determine what kind of call this is
	if _, isPtr := functionExpr.Type().(*types.PointerType); isPtr {
		// If it's a pointer, assume it's an object/method call
		// and use dynamic dispatch

		// In a full implementation, you would extract the method name
		// from the original call expression
		methodName := "method" // Placeholder - you'd need actual method name

		return cg.generateDynamicDispatch(functionExpr, methodName, args)
	} else {
		// Otherwise assume it's a direct function call
		return block.NewCall(functionExpr, args...)
	}
}

// generateRegularCall creates LLVM IR for direct function calls to known functions
func (cg *CodeGenerator) generateRegularCall(funcName string, args []value.Value) value.Value {
	// Get the current block
	block := cg.CurrentBlock

	// Look up the function in the current scope/module
	funcValue, exists := cg.Symbols[funcName]
	if !exists {
		// Try to find in the module's functions
		var foundFunc *ir.Func
		for _, f := range cg.Module.Funcs {
			if f.Name() == funcName {
				foundFunc = f
				break
			}
		}

		if foundFunc == nil {
			panic(fmt.Sprintf("undefined function: %s", funcName))
		}
		funcValue = foundFunc
	}

	// Call the function
	return block.NewCall(funcValue, args...)
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

	// Declare String.concat and String.substr functions
	stringType := cg.TypeMap["String"]
	cg.Module.NewFunc(
		"String.concat",
		types.NewPointer(stringType), // Returns a pointer to the new String object
		ir.NewParam("self", types.NewPointer(stringType)),
		ir.NewParam("other", types.NewPointer(stringType)),
	)

	cg.Module.NewFunc(
		"String.substr",
		types.NewPointer(stringType), // Returns a pointer to the new String object
		ir.NewParam("self", types.NewPointer(stringType)),
		ir.NewParam("start", types.I32),
		ir.NewParam("length", types.I32),
	)
}

// Add this new method to handle IO.out_string code generation
func (cg *CodeGenerator) generateIOOutStringMethod(class *ast.Class, method *ast.Method) {
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
func (cg *CodeGenerator) generateIOOutIntMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
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
func (cg *CodeGenerator) generateObjectAbortMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
	objType := cg.TypeMap["Object"]

	funcDecl := cg.Module.NewFunc(
		funcName,
		types.NewPointer(objType),
		ir.NewParam("self", types.NewPointer(objType)),
	)

	entry := funcDecl.NewBlock("entry")

	// Print error message
	errorMsg := constant.NewCharArrayFromString("Program aborted\n\x00")
	global := cg.Module.NewGlobalDef(".str.abort_msg", errorMsg)
	global.Immutable = true
	msgPtr := constant.NewGetElementPtr(errorMsg.Type(), global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	entry.NewCall(cg.StdlibFuncs["printf"], msgPtr)

	// Exit with status code 1
	entry.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
	entry.NewUnreachable() // exit doesn't return
}

func (cg *CodeGenerator) generateTypeNameMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
	objType := cg.TypeMap["Object"]

	funcDecl := cg.Module.NewFunc(
		funcName,
		types.NewPointer(types.I8), // Returns String (i8*)
		ir.NewParam("self", types.NewPointer(objType)),
	)

	entry := funcDecl.NewBlock("entry")

	// Create class name string
	className := class.Name.Value
	strConst := constant.NewCharArrayFromString(className + "\x00")
	global := cg.Module.NewGlobalDef(fmt.Sprintf(".str.%s", className), strConst)
	global.Immutable = true
	gep := constant.NewGetElementPtr(strConst.Type(), global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	entry.NewRet(gep)
}

func (cg *CodeGenerator) generateCopyMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
	objType := cg.TypeMap["Object"]

	funcDecl := cg.Module.NewFunc(
		funcName,
		types.NewPointer(objType), // Returns SELF_TYPE
		ir.NewParam("self", types.NewPointer(objType)),
	)

	entry := funcDecl.NewBlock("entry")
	// In a real implementation this would perform a shallow copy
	// For now just return self
	entry.NewRet(funcDecl.Params[0])
}

// Add this new method to handle IO.in_int
func (cg *CodeGenerator) generateIOInIntMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
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
func (cg *CodeGenerator) generateIOInStringMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
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
func (cg *CodeGenerator) generateStringLengthMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
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
	// Get the string pointer from the String object
	strPtr := entry.NewGetElementPtr(
		stringType,
		funcDecl.Params[0],
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1), // Access the string data field
	)
	loadedPtr := entry.NewLoad(types.NewPointer(types.I8), strPtr)

	// Call strlen instead of string_length
	length := entry.NewCall(cg.StdlibFuncs["strlen"], loadedPtr)
	entry.NewRet(length)
}

// generateStringConcatMethod generates LLVM IR for String.concat
func (cg *CodeGenerator) generateStringConcatMethod(class *ast.Class, method *ast.Method) {
	// Find the function (it should already be declared)
	mangledName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
	var concatFunc *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == mangledName {
			concatFunc = f
			break
		}
	}

	if concatFunc == nil {
		panic(fmt.Sprintf("String.concat function not found: %s", mangledName))
	}

	// Create entry block
	entryBlock := concatFunc.NewBlock("entry")
	cg.CurrentFunc = concatFunc

	// Get the 'self' and 'other' string parameters
	selfParam := concatFunc.Params[0]
	otherParam := concatFunc.Params[1]

	// 1. Allocate a new String object
	newStringObj := cg.generateObjectAllocation("String")

	// 2. Load the string data (i8*) from 'self' and 'other'
	block := entryBlock // For brevity
	selfStringPtr := block.NewGetElementPtr(
		cg.TypeMap["String"],
		selfParam,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	selfStringValue := block.NewLoad(types.NewPointer(types.I8), selfStringPtr)

	otherStringPtr := block.NewGetElementPtr(
		cg.TypeMap["String"],
		otherParam,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	otherStringValue := block.NewLoad(types.NewPointer(types.I8), otherStringPtr)

	// 3. Use standard C library functions for concatenation
	// Calculate the length of both strings
	selfLen := block.NewCall(cg.StdlibFuncs["strlen"], selfStringValue)
	otherLen := block.NewCall(cg.StdlibFuncs["strlen"], otherStringValue)

	// Calculate total length needed
	totalLen := block.NewAdd(selfLen, otherLen)
	// Add 1 for the null terminator
	allocSize := block.NewAdd(totalLen, constant.NewInt(types.I32, 1))

	// Allocate memory for the new string
	mallocCall := block.NewCall(cg.StdlibFuncs["malloc"],
		block.NewZExt(allocSize, types.I64))

	// Copy the first string
	block.NewCall(cg.StdlibFuncs["strcpy"], mallocCall, selfStringValue)

	// Concatenate the second string
	block.NewCall(cg.StdlibFuncs["strcat"], mallocCall, otherStringValue)

	// 4. Store the result (i8*) into the new String object's string field
	newStringDataPtr := block.NewGetElementPtr(
		cg.TypeMap["String"],
		newStringObj,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	block.NewStore(mallocCall, newStringDataPtr)

	// 5. Return the new String object
	block.NewRet(newStringObj)
}

// generateStringSubstrMethod generates LLVM IR for String.substr
func (cg *CodeGenerator) generateStringSubstrMethod(class *ast.Class, method *ast.Method) {
	funcName := fmt.Sprintf("%s.%s", class.Name.Value, method.Name.Value)
	stringType := cg.TypeMap["String"]

	// Find existing function declaration
	var funcDecl *ir.Func
	for _, f := range cg.Module.Funcs {
		if f.Name() == funcName {
			funcDecl = f
			break
		}
	}

	if funcDecl == nil {
		// Create function only if not found (shouldn't happen if vtables were generated first)
		funcDecl = cg.Module.NewFunc(
			funcName,
			types.NewPointer(stringType),
			ir.NewParam("self", types.NewPointer(stringType)),
			ir.NewParam("start", types.I32),
			ir.NewParam("length", types.I32),
		)
	}

	// Only generate implementation if function doesn't have a body
	if len(funcDecl.Blocks) == 0 {
		entry := funcDecl.NewBlock("entry")

		// Get the underlying C string from the String object
		strPtr := entry.NewGetElementPtr(
			stringType,
			funcDecl.Params[0],
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 1),
		)
		loadedStr := entry.NewLoad(types.NewPointer(types.I8), strPtr)
		startIdx := funcDecl.Params[1]
		length := funcDecl.Params[2]

		// Perform bounds checking
		strLen := entry.NewCall(cg.StdlibFuncs["strlen"], loadedStr)

		// Create a new block for the bounds check
		boundsCheckBlock := funcDecl.NewBlock("bounds_check")
		allocBlock := funcDecl.NewBlock("alloc")
		errorBlock := funcDecl.NewBlock("error")

		// Branch to the bounds check block
		entry.NewBr(boundsCheckBlock)

		// Check if start < 0 or start >= strLen or length < 0
		startOutOfBounds := boundsCheckBlock.NewOr(
			boundsCheckBlock.NewICmp(enum.IPredSLT, startIdx, constant.NewInt(types.I32, 0)),
			boundsCheckBlock.NewICmp(enum.IPredSGE, startIdx, strLen),
		)
		lengthNegative := boundsCheckBlock.NewICmp(enum.IPredSLT, length, constant.NewInt(types.I32, 0))
		invalidInput := boundsCheckBlock.NewOr(startOutOfBounds, lengthNegative)

		// Branch based on the bounds check
		boundsCheckBlock.NewCondBr(invalidInput, errorBlock, allocBlock)

		// Handle error case
		errorMsgGlobal := cg.Module.NewGlobalDef(".str.substr_error",
			constant.NewCharArrayFromString("Runtime error: substring out of range\n\x00"))
		errorMsgGlobal.Immutable = true

		errorMsgPtr := errorBlock.NewBitCast(errorMsgGlobal, types.NewPointer(types.I8))

		errorBlock.NewCall(cg.StdlibFuncs["printf"], errorMsgPtr)
		errorBlock.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
		errorBlock.NewUnreachable()

		// Allocate the new substring in the alloc block
		// Calculate how much memory to allocate (length + 1 for null terminator)
		allocSize := allocBlock.NewAdd(length, constant.NewInt(types.I32, 1))
		mallocCall := allocBlock.NewCall(cg.StdlibFuncs["malloc"],
			allocBlock.NewZExt(allocSize, types.I64))

		// Get pointer to the start of the substring in the original string
		substrStart := allocBlock.NewGetElementPtr(types.I8, loadedStr, startIdx)

		// Use strncpy to copy exactly 'length' characters
		allocBlock.NewCall(cg.StdlibFuncs["strncpy"], mallocCall, substrStart, length)

		// Add null terminator at the end of the new substring
		nullTermPtr := allocBlock.NewGetElementPtr(types.I8, mallocCall, length)
		allocBlock.NewStore(constant.NewInt(types.I8, 0), nullTermPtr)

		// Create new String object with proper allocation
		// Calculate the size for malloc
		sizeGEP := constant.NewGetElementPtr(
			stringType,
			constant.NewNull(types.NewPointer(stringType)),
			constant.NewInt(types.I32, 1),
		)

		// Call malloc with the size of the String struct
		newStringMalloc := allocBlock.NewCall(cg.StdlibFuncs["malloc"], sizeGEP)

		// Cast the malloc result to String pointer type
		newString := allocBlock.NewBitCast(newStringMalloc, types.NewPointer(stringType))

		// Set up the vtable pointer
		vtableFieldPtr := allocBlock.NewGetElementPtr(
			stringType,
			newString,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0),
		)
		vtable, exists := cg.VTables["String"]
		if !exists {
			panic("String vtable not found")
		}

		allocBlock.NewStore(
			allocBlock.NewBitCast(vtable, types.NewPointer(types.I8)),
			vtableFieldPtr,
		)

		// Store the substring in the new String object
		newStrPtr := allocBlock.NewGetElementPtr(
			stringType,
			newString,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 1),
		)
		allocBlock.NewStore(mallocCall, newStrPtr)

		allocBlock.NewRet(newString)
	}
}
