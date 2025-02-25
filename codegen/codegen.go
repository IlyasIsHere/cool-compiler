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

	// Symbol table for variables in current scope
	Symbols map[string]value.Value

	// Standard library functions
	StdlibFuncs map[string]*ir.Func
}

// Generate is the main entry point for code generation
func Generate(program *ast.Program) (*ir.Module, error) {
	// Initialize code generator
	cg := NewCodeGenerator()

	// Define built-in COOL classes first
	cg.DefineBuiltInClasses()

	// Generate class structures
	cg.GenerateClassStructs(program)

	// Generate virtual method tables for all classes
	cg.GenerateVTables(program)

	// Generate method implementations
	cg.GenerateMethods(program)

	// Generate the main function, which initializes the Main class and calls Main.main()
	cg.GenerateMain(program)

	return cg.Module, nil
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	cg := &CodeGenerator{
		Module:         ir.NewModule(),
		TypeMap:        make(map[string]*types.StructType),
		VTables:        make(map[string]*ir.Global),
		ClassHierarchy: make(map[string]string),
		Symbols:        make(map[string]value.Value),
		StdlibFuncs:    make(map[string]*ir.Func),
	}

	// Initialize standard library functions
	cg.initStdlib()

	return cg
}

// initStdlib initializes standard library functions
func (cg *CodeGenerator) initStdlib() {
	// Memory management functions

	// malloc for object allocation
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

	// IO functions

	// out_string - print a string to stdout
	cg.StdlibFuncs["out_string"] = cg.Module.NewFunc(
		"out_string",
		types.I32,
		ir.NewParam("str", types.NewPointer(types.I8)),
	)

	// out_int - print an integer to stdout
	cg.StdlibFuncs["out_int"] = cg.Module.NewFunc(
		"out_int",
		types.I32,
		ir.NewParam("num", types.I32),
	)

	// in_string - read a string from stdin
	cg.StdlibFuncs["in_string"] = cg.Module.NewFunc(
		"in_string",
		types.NewPointer(types.I8),
	)

	// in_int - read an integer from stdin
	cg.StdlibFuncs["in_int"] = cg.Module.NewFunc(
		"in_int",
		types.I32,
	)

	// String manipulation functions

	// string_length - get the length of a string
	cg.StdlibFuncs["string_length"] = cg.Module.NewFunc(
		"string_length",
		types.I32,
		ir.NewParam("str", types.NewPointer(types.I8)),
	)

	// string_concat - concatenate two strings
	cg.StdlibFuncs["string_concat"] = cg.Module.NewFunc(
		"string_concat",
		types.NewPointer(types.I8),
		ir.NewParam("str1", types.NewPointer(types.I8)),
		ir.NewParam("str2", types.NewPointer(types.I8)),
	)

	// string_substr - get a substring
	cg.StdlibFuncs["string_substr"] = cg.Module.NewFunc(
		"string_substr",
		types.NewPointer(types.I8),
		ir.NewParam("str", types.NewPointer(types.I8)),
		ir.NewParam("start", types.I32),
		ir.NewParam("length", types.I32),
	)

	// Runtime support functions

	// abort - terminate the program with an error message
	cg.StdlibFuncs["abort"] = cg.Module.NewFunc(
		"abort",
		types.Void,
	)

	// type_name - get the name of an object's type as a string
	cg.StdlibFuncs["type_name"] = cg.Module.NewFunc(
		"type_name",
		types.NewPointer(types.I8),
		ir.NewParam("obj", types.NewPointer(types.I8)),
	)

	// copy - create a shallow copy of an object
	cg.StdlibFuncs["object_copy"] = cg.Module.NewFunc(
		"object_copy",
		types.NewPointer(types.I8),
		ir.NewParam("obj", types.NewPointer(types.I8)),
	)

	// Runtime type checking for case expressions
	cg.StdlibFuncs["case_abort"] = cg.Module.NewFunc(
		"case_abort",
		types.Void,
	)

	// Dispatch on void check
	cg.StdlibFuncs["dispatch_abort"] = cg.Module.NewFunc(
		"dispatch_abort",
		types.Void,
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

	// Collect all fields for this class (including inherited fields)
	var fields []types.Type
	var parentFields []types.Type

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

			fields = append(fields, attrType)
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

	// Create global constant for vtable
	vtableName := fmt.Sprintf("vtable.%s", className)
	vtable := cg.Module.NewGlobalDef(vtableName, constant.NewZeroInitializer(vtableType))

	// Store in the VTables map
	cg.VTables[className] = vtable
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
	// Generate implementation for each class method
	for _, class := range program.Classes {
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

	// The method body should compute a value of the correct return type
	if bodyValue.Type() != methodFunc.Sig.RetType {
		// If types don't match, we may need to cast or handle special cases
		// This is a simplified version - real implementation would need more robust type handling
		if ptr, isPtrType := methodFunc.Sig.RetType.(*types.PointerType); isPtrType {
			if bodyValue.Type() == types.I32 || bodyValue.Type() == types.I1 {
				// Boxing primitive values when returning as Object
				// In a real implementation, would create proper boxed objects
				bodyValue = entryBlock.NewIntToPtr(bodyValue, ptr)
			} else if _, isOtherPtr := bodyValue.Type().(*types.PointerType); isOtherPtr {
				// Cast between pointer types
				bodyValue = entryBlock.NewBitCast(bodyValue, ptr)
			}
		}
	}

	// Return the method result
	entryBlock.NewRet(bodyValue)
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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

	// Get the LLVM struct type for the class
	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("attempt to create an instance of unknown type: %s", typeName))
	}

	// Allocate memory for the object using malloc
	// For this, we need the size of the class struct
	sizeOfClass := constant.NewPtrToInt(
		constant.NewGetElementPtr(
			classType,
			constant.NewNull(types.NewPointer(classType)),
			constant.NewInt(types.I32, 1),
		),
		types.I64,
	)

	// Call malloc to allocate memory for the object
	mallocFunc, exists := cg.StdlibFuncs["malloc"]
	if !exists {
		// Define malloc if it doesn't exist
		mallocFunc = cg.Module.NewFunc(
			"malloc",
			types.NewPointer(types.I8),
			ir.NewParam("size", types.I64),
		)
		cg.StdlibFuncs["malloc"] = mallocFunc
	}

	// Call malloc with the size of the class
	mallocCall := block.NewCall(mallocFunc, sizeOfClass)

	// Cast the result from i8* to the class pointer type
	objectPtr := block.NewBitCast(mallocCall, types.NewPointer(classType))

	// Get the vtable for this class
	vtable, exists := cg.VTables[typeName]
	if !exists {
		panic(fmt.Sprintf("vtable not found for type: %s", typeName))
	}

	// Get pointer to the vtable field in the object (first field)
	vtableFieldPtr := block.NewGetElementPtr(
		classType,
		objectPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	// We need to cast the vtable to the correct type for storage
	// The first field is a pointer to i8 (void*) which is where we store the vtable pointer
	// Cast the vtable to i8* for storage
	vtablePtr := block.NewBitCast(vtable, types.NewPointer(types.I8))

	// Store the vtable pointer in the object
	block.NewStore(vtablePtr, vtableFieldPtr)

	// Initialize other fields with default values
	// For a real implementation, we would iterate over all fields and initialize them

	// Return the newly allocated and initialized object
	return objectPtr
}

// generateDynamicDispatch generates code for method dispatch
func (cg *CodeGenerator) generateDynamicDispatch(object value.Value, methodName string, args []value.Value) value.Value {
	// Get the current block
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
		panic(fmt.Sprintf("expected object to point to a struct type, got: %v", objPtrType.ElemType))
	}

	// Get pointer to the vtable field in the object (first field)
	vtablePtrPtr := block.NewGetElementPtr(
		objStructType,
		object,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	// Load the vtable pointer
	vtablePtr := block.NewLoad(types.NewPointer(types.I8), vtablePtrPtr) // Adjust type as needed

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
	// We need to cast vtablePtr to the appropriate type first
	vtableType := types.NewPointer(types.NewArray(0, types.NewPointer(types.I8))) // Example vtable type
	castedVTablePtr := block.NewBitCast(vtablePtr, vtableType)

	// Get pointer to the method slot
	methodSlotPtr := block.NewGetElementPtr(vtableType.ElemType, castedVTablePtr, gepIndices...)

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

// GenerateMain generates the LLVM main function
func (cg *CodeGenerator) GenerateMain(program *ast.Program) {
	// Create the main function with signature: int main()
	mainFunc := cg.Module.NewFunc("main", types.I32)
	entryBlock := mainFunc.NewBlock("entry")

	// Set the current function for code generation
	cg.CurrentFunc = mainFunc

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

	// Try to find the identifier in the symbol table
	val, exists := cg.Symbols[identifier.Value]
	if !exists {
		// This should not happen if semantic analysis was successful
		panic(fmt.Sprintf("undefined identifier encountered during code generation: %s", identifier.Value))
	}

	// Check if the identifier refers to a local variable or parameter (already stored in register)
	if _, isLocalVar := val.(*ir.InstAlloca); isLocalVar {
		// For local variables (alloca instructions), we need to load the value
		block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1] // Get current block
		load := block.NewLoad(val.Type().(*types.PointerType).ElemType, val)
		return load
	}

	// For non-local variables (global, parameters, etc.)
	return val
}

// generateAssignmentExpression creates LLVM IR for variable assignment
func (cg *CodeGenerator) generateAssignmentExpression(assign *ast.AssignmentExpression) value.Value {
	// First, generate code for the right-hand side expression
	rhsValue := cg.generateExpression(assign.Expression)

	// Get the target variable from the symbol table
	target, exists := cg.Symbols[assign.Identifier.Value]
	if !exists {
		panic(fmt.Sprintf("undefined identifier in assignment: %s", assign.Identifier.Value))
	}

	// Get the current basic block
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

	// Handle different target types:
	if allocaInst, isLocalVar := target.(*ir.InstAlloca); isLocalVar {
		// For local variables (created with alloca), we use a store instruction
		block.NewStore(rhsValue, allocaInst)
	} else if _, isParam := target.(*ir.Param); isParam {
		// For parameters, we need to create a store if they've been copied to local storage
		// (In SSA form, parameters are immutable, so we typically need to have created an alloca for them)
		panic("assignment to parameter not properly handled - parameters should have local storage")
	} else {
		// Other cases (e.g., global variables) would be handled here
		panic(fmt.Sprintf("unsupported assignment target type: %T", target))
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

	// Create the basic blocks for the true, false, and merge paths
	currentFunc := cg.CurrentFunc
	trueBlock := currentFunc.NewBlock("if.then")
	falseBlock := currentFunc.NewBlock("if.else")
	mergeBlock := currentFunc.NewBlock("if.end")

	// Create the conditional branch
	currentBlock := currentFunc.Blocks[len(currentFunc.Blocks)-3] // Get the block before our new ones
	currentBlock.NewCondBr(condValue, trueBlock, falseBlock)

	// Generate code for the true branch
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-3] = trueBlock // Set current block to true block
	trueValue := cg.generateExpression(ifExpr.Consequence)
	trueBlock = cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-3] // Get updated true block
	trueBlock.NewBr(mergeBlock)

	// Generate code for the false branch
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-2] = falseBlock // Set current block to false block
	falseValue := cg.generateExpression(ifExpr.Alternative)
	falseBlock = cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-2] // Get updated false block
	falseBlock.NewBr(mergeBlock)

	// Set current block to merge block
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1] = mergeBlock

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

	// Create a PHI node to merge the values
	phi := mergeBlock.NewPhi()
	phi.Typ = resultType

	// Add the incoming values to the PHI node
	phi.Incs = append(phi.Incs, &ir.Incoming{X: trueValue, Pred: trueBlock})
	phi.Incs = append(phi.Incs, &ir.Incoming{X: falseValue, Pred: falseBlock})

	return phi
}

// generateWhileExpression creates LLVM IR for a while loop expression
func (cg *CodeGenerator) generateWhileExpression(whileExpr *ast.WhileExpression) value.Value {
	// Create the basic blocks for the loop
	currentFunc := cg.CurrentFunc
	condBlock := currentFunc.NewBlock("while.cond")
	bodyBlock := currentFunc.NewBlock("while.body")
	exitBlock := currentFunc.NewBlock("while.exit")

	// Get the current block and create a branch to the condition block
	currentBlock := currentFunc.Blocks[len(currentFunc.Blocks)-3] // Get the block before our new ones
	currentBlock.NewBr(condBlock)

	// Set current block to condition block
	cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-3] = condBlock

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

		// Each expression's result is discarded except for the last one
		// We still generate code for all expressions because they might have side effects
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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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

	default:
		panic(fmt.Sprintf("unsupported binary operator: %s", binExpr.Operator))
	}
}

// generateDotCallExpression creates LLVM IR for method calls on objects
func (cg *CodeGenerator) generateDotCallExpression(dotCall *ast.DotCallExpression) value.Value {
	// Generate code for the object on which the method is called
	objectValue := cg.generateExpression(dotCall.Object)

	// Generate LLVM values for all arguments
	argValues := make([]value.Value, 0, len(dotCall.Arguments)+1)

	// The first argument to a method call is always the object itself (self)
	argValues = append(argValues, objectValue)

	// Add the rest of the arguments
	for _, arg := range dotCall.Arguments {
		argValues = append(argValues, cg.generateExpression(arg))
	}

	// Determine if this is a dynamic or static dispatch
	if dotCall.Type != nil {
		// This is a static dispatch (e.g., obj@Type.method())
		return cg.generateStaticDispatch(objectValue, dotCall.Type.Value, dotCall.Method.Value, argValues)
	} else {
		// This is a dynamic dispatch (e.g., obj.method())
		return cg.generateDynamicDispatch(objectValue, dotCall.Method.Value, argValues)
	}
}

// generateStaticDispatch creates LLVM IR for static dispatch (obj@Type.method())
func (cg *CodeGenerator) generateStaticDispatch(object value.Value, typeName string, methodName string, args []value.Value) value.Value {
	// Get the current basic block
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
	// Get the current block
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
	// Get the current function
	currentFunc := cg.CurrentFunc

	// Generate code for the expression being dispatched on
	exprValue := cg.generateExpression(caseExpr.Expression)

	// Create a basic block for the end of the case expression
	// All branches will merge to this block
	endBlock := currentFunc.NewBlock("case.end")

	// Create a basic block for each branch
	branchBlocks := make([]*ir.Block, len(caseExpr.Branches))
	for i := range caseExpr.Branches {
		branchBlocks[i] = currentFunc.NewBlock(fmt.Sprintf("case.branch.%d", i))
	}

	// Create a default block for when no branch matches
	// This should never happen in well-typed COOL code, but we need it for LLVM
	noMatchBlock := currentFunc.NewBlock("case.nomatch")

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
		notNullBlock := currentFunc.NewBlock("case.notnull")

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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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
	block := cg.CurrentFunc.Blocks[len(cg.CurrentFunc.Blocks)-1]

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

	// Define Int class
	intClass := &ast.Class{
		Name:     &ast.TypeIdentifier{Value: "Int"},
		Parent:   &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{},
	}
	cg.declareClassType(intClass)
	cg.defineClassStruct(intClass)

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

	// Define Bool class
	boolClass := &ast.Class{
		Name:     &ast.TypeIdentifier{Value: "Bool"},
		Parent:   &ast.TypeIdentifier{Value: "Object"},
		Features: []ast.Feature{},
	}
	cg.declareClassType(boolClass)
	cg.defineClassStruct(boolClass)

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
}
