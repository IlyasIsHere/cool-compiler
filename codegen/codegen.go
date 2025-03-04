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
	Value       value.Value
	BranchTypes []string
}

// CodeGenerator handles the code generation process
type CodeGenerator struct {
	Module            *ir.Module
	TypeMap           map[string]*types.StructType
	VTables           map[string]*ir.Global
	ClassHierarchy    map[string]string
	CurrentFunc       *ir.Func
	CurrentBlock      *ir.Block
	Symbols           map[string]value.Value
	StdlibFuncs       map[string]*ir.Func
	BuiltInClasses    []*ast.Class
	ProgramClasses    []*ast.Class
	ClassNameToAST    map[string]*ast.Class
	AttributeIndices  map[string]map[string]int
	MethodIndices     map[string]map[string]int
	IfCounter         int
	WhileCounter      int
	CaseCounter       int
	EmptyStringGlobal *ir.Global
	CaseResults       map[value.Value]*CaseResult
}

// Generate is the main entry point for code generation
func Generate(program *ast.Program) (*ir.Module, error) {
	cg := NewCodeGenerator()
	cg.DefineBuiltInClasses()
	cg.initStdlib()
	cg.ProgramClasses = program.Classes

	for _, class := range program.Classes {
		cg.ClassNameToAST[class.Name.Value] = class
	}

	for _, builtInClass := range cg.BuiltInClasses {
		cg.ClassNameToAST[builtInClass.Name.Value] = builtInClass
	}

	cg.GenerateClassStructs(program)
	cg.GenerateVTables(program)
	cg.GenerateMethods(program)

	if _, exists := cg.TypeMap["Main"]; !exists {
		panic("Program must have a Main class")
	}

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
	cg.StdlibFuncs["malloc"] = cg.Module.NewFunc(
		"malloc",
		types.NewPointer(types.I8),
		ir.NewParam("size", types.I64),
	)

	cg.StdlibFuncs["free"] = cg.Module.NewFunc(
		"free",
		types.Void,
		ir.NewParam("ptr", types.NewPointer(types.I8)),
	)

	cg.StdlibFuncs["exit"] = cg.Module.NewFunc(
		"exit",
		types.Void,
		ir.NewParam("status", types.I32),
	)

	// IO functions
	printfFunc := cg.Module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true
	cg.StdlibFuncs["printf"] = printfFunc

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
	// First declare all class types, then define their structures
	for _, class := range program.Classes {
		cg.declareClassType(class)
	}

	for _, class := range program.Classes {
		cg.defineClassStruct(class)
	}
}

// declareClassType creates an LLVM struct type for a class
func (cg *CodeGenerator) declareClassType(class *ast.Class) {
	if _, exists := cg.TypeMap[class.Name.Value]; exists {
		return
	}

	if class.Parent != nil {
		cg.ClassHierarchy[class.Name.Value] = class.Parent.Value
	} else if class.Name.Value != "Object" {
		cg.ClassHierarchy[class.Name.Value] = "Object"
	}

	structType := types.NewStruct()
	cg.Module.NewTypeDef(class.Name.Value, structType)
	cg.TypeMap[class.Name.Value] = structType
}

// defineClassStruct defines the fields of a class struct
func (cg *CodeGenerator) defineClassStruct(class *ast.Class) {
	className := class.Name.Value
	classType := cg.TypeMap[className]

	if _, exists := cg.AttributeIndices[className]; !exists {
		cg.AttributeIndices[className] = make(map[string]int)
	}

	var fields []types.Type
	var parentFields []types.Type
	fieldIndex := 1 // Start at 1 because index 0 is vtable pointer

	// First field is always a pointer to the vtable
	fields = append(fields, types.NewPointer(types.I8))

	// Include parent fields if this class inherits
	if parent, exists := cg.ClassHierarchy[className]; exists && parent != "" {
		parentType, exists := cg.TypeMap[parent]
		if !exists {
			panic(fmt.Sprintf("parent class %s not found for class %s", parent, className))
		}

		for i := 1; i < len(parentType.Fields); i++ {
			parentFields = append(parentFields, parentType.Fields[i])
		}

		fields = append(fields, parentFields...)

		for attrName, attrIndex := range cg.AttributeIndices[parent] {
			cg.AttributeIndices[className][attrName] = attrIndex
		}

		fieldIndex += len(parentFields)
	}

	// Add class's own fields
	for _, feature := range class.Features {
		if attr, isAttr := feature.(*ast.Attribute); isAttr {
			var attrType types.Type

			switch attr.TypeDecl.Value {
			case "Int":
				attrType = types.I32
			case "Bool":
				attrType = types.I1
			case "String":
				attrType = types.NewPointer(types.I8)
			case "SELF_TYPE":
				attrType = types.NewPointer(classType)
			default:
				referencedType, exists := cg.TypeMap[attr.TypeDecl.Value]
				if !exists {
					panic(fmt.Sprintf("undefined type %s in attribute %s of class %s",
						attr.TypeDecl.Value, attr.Name.Value, className))
				}
				attrType = types.NewPointer(referencedType)
			}

			fields = append(fields, attrType)
			cg.AttributeIndices[className][attr.Name.Value] = fieldIndex
			fieldIndex++
		}
	}

	classType.Fields = fields
}

// GenerateVTables creates virtual method tables for all classes
func (cg *CodeGenerator) GenerateVTables(program *ast.Program) {
	for _, class := range program.Classes {
		cg.createVTableForClass(class, program)
	}
}

// createVTableForClass creates a vtable for a specific class
func (cg *CodeGenerator) createVTableForClass(class *ast.Class, program *ast.Program) {
	className := class.Name.Value

	if _, exists := cg.VTables[className]; exists {
		return
	}

	if _, exists := cg.MethodIndices[className]; !exists {
		cg.MethodIndices[className] = make(map[string]int)
	}

	methods := make(map[string]*ir.Func)
	methodIndices := make(map[string]int)
	var methodNames []string

	// Process parent class first to inherit its methods
	if parent, exists := cg.ClassHierarchy[className]; exists && parent != "" {
		parentClass := findClass(parent, program)
		if parentClass != nil {
			if _, parentVTableExists := cg.VTables[parent]; !parentVTableExists {
				cg.createVTableForClass(parentClass, program)
			}
		}

		// Copy ALL methods from parent vtable, maintaining the same indices
		currentParent := parent
		for currentParent != "" {
			for methodName, methodIndex := range cg.MethodIndices[currentParent] {
				if _, exists := methods[methodName]; !exists {
					methodIndices[methodName] = methodIndex

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

			currentParent = cg.ClassHierarchy[currentParent]
		}
	}

	// Process methods of this class
	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod {
			methodName := method.Name.Value

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

			function := cg.Module.NewFunc(funcName, returnType, params...)

			methods[methodName] = function

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
			methodIndices[name] = len(finalMethods) - 1
		}
	}

	// Create an array of function pointers for the vtable
	methodCount := len(finalMethods)
	vtableType := types.NewArray(uint64(methodCount), types.NewPointer(types.I8))

	vtableName := fmt.Sprintf("vtable.%s", className)

	initializers := make([]constant.Constant, methodCount)
	for i, method := range finalMethods {
		initializers[i] = constant.NewBitCast(method, types.NewPointer(types.I8))
	}

	var vtable *ir.Global
	if methodCount > 0 {
		arrayConst := &constant.Array{
			Typ:   vtableType,
			Elems: initializers,
		}
		vtable = cg.Module.NewGlobalDef(vtableName, arrayConst)
	} else {
		vtable = cg.Module.NewGlobalDef(vtableName, constant.NewZeroInitializer(vtableType))
	}

	cg.VTables[className] = vtable
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

	// Find the function declaration
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
	cg.CurrentFunc = methodFunc
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
		cg.CurrentBlock.NewRet(selfParam)
		return
	}

	// Convert the bodyValue to the correct return type if needed
	returnType := cg.getLLVMTypeForCOOLTypeWithContext(method.TypeDecl.Value, className)

	// Check if we need to convert from primitive type to object pointer
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
	block := cg.CurrentBlock

	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("attempt to create an instance of unknown type: %s", typeName))
	}

	sizeValue := cg.getSizeOf(classType, block)

	mallocFunc, exists := cg.StdlibFuncs["malloc"]
	if !exists {
		panic("malloc function not found")
	}

	mallocCall := block.NewCall(mallocFunc, sizeValue)
	objectPtr := block.NewBitCast(mallocCall, types.NewPointer(classType))

	vtable, exists := cg.VTables[typeName]
	if !exists {
		panic(fmt.Sprintf("vtable not found for type: %s", typeName))
	}

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

	cg.initializeAttributes(typeName, objectPtr)

	return objectPtr
}

// initializeAttributes initializes all attributes of a class with their default values
func (cg *CodeGenerator) initializeAttributes(className string, objectPtr value.Value) {
	block := cg.CurrentBlock

	_, exists := cg.ClassNameToAST[className]
	if !exists {
		return
	}

	classType := cg.TypeMap[className]
	oldSelf, hasSelf := cg.Symbols["self"]
	cg.Symbols["self"] = objectPtr

	ancestors := []string{className}
	current := className
	for {
		parent, exists := cg.ClassHierarchy[current]
		if !exists || parent == "" {
			break
		}
		ancestors = append([]string{parent}, ancestors...)
		current = parent
	}

	for _, ancestorName := range ancestors {
		ancestor, exists := cg.ClassNameToAST[ancestorName]
		if !exists {
			continue
		}

		for _, feature := range ancestor.Features {
			if attr, isAttr := feature.(*ast.Attribute); isAttr {

				attrIndex, exists := cg.AttributeIndices[className][attr.Name.Value]
				if !exists {
					continue
				}

				attrPtr := block.NewGetElementPtr(
					classType,
					objectPtr,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, int64(attrIndex)),
				)

				var initValue value.Value

				if attr.Expression != nil {
					initValue = cg.generateExpression(attr.Expression)
				} else {
					switch attr.TypeDecl.Value {
					case "Int":
						initValue = constant.NewInt(types.I32, 0)
					case "Bool":
						initValue = constant.NewInt(types.I1, 0)
					case "String":
						initValue = cg.EmptyStringGlobal
					default:
						objType := cg.getLLVMTypeForCOOLType(attr.TypeDecl.Value)
						if ptrType, ok := objType.(*types.PointerType); ok {
							initValue = constant.NewNull(ptrType)
						} else {
							initValue = constant.NewZeroInitializer(objType)
						}
					}
				}

				attrType := classType.Fields[attrIndex]

				if !initValue.Type().Equal(attrType) {
					initValue = block.NewBitCast(initValue, attrType)
				}

				block.NewStore(initValue, attrPtr)
			}
		}
	}

	if hasSelf {
		cg.Symbols["self"] = oldSelf
	} else {
		delete(cg.Symbols, "self")
	}
}

// getObjectRuntimeType determines the runtime type of an object
func (cg *CodeGenerator) getObjectRuntimeType(object value.Value) string {
	switch object.Type() {
	case types.I1:
		return "Bool"
	case types.I32:
		return "Int"
	case types.I8Ptr:
		return "String"
	case types.I8:
		return "String"
	}

	if ptrType, ok := object.Type().(*types.PointerType); ok {
		if ptrType.ElemType == types.I8 {
			return "String"
		}
	}

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

	methodSlotPtr := block.NewGetElementPtr(
		vtable.ContentType,
		vtable,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(methodIndex)),
	)

	methodPtr := block.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	targetMethod := cg.findMethodByName(runtimeType, methodName)
	if targetMethod == nil {
		panic(fmt.Sprintf("method %s not found in class %s for signature lookup", methodName, runtimeType))
	}

	methodDeclaringClass := cg.findMethodDeclaringClass(runtimeType, methodName)
	if methodDeclaringClass == "" {
		panic(fmt.Sprintf("declaring class for method %s not found", methodName))
	}

	classType, exists := cg.TypeMap[methodDeclaringClass]
	if !exists {
		panic(fmt.Sprintf("class type %s not found", methodDeclaringClass))
	}

	paramTypes := []types.Type{types.NewPointer(classType)} // Self parameter first
	for _, formal := range targetMethod.Formals {
		paramType := cg.getLLVMTypeForCOOLType(formal.TypeDecl.Value)
		paramTypes = append(paramTypes, paramType)
	}

	returnType := cg.getLLVMTypeForCOOLTypeWithContext(targetMethod.TypeDecl.Value, methodDeclaringClass)
	funcType := types.NewPointer(types.NewFunc(returnType, paramTypes...))

	castedMethodPtr := block.NewBitCast(methodPtr, funcType)

	selfArg := block.NewBitCast(object, types.NewPointer(classType))

	allArgs := make([]value.Value, 0, len(args)+1)
	allArgs = append(allArgs, selfArg)

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

	return block.NewCall(castedMethodPtr, allArgs...)
}

// findMethodDeclaringClass returns the class that defines a method
func (cg *CodeGenerator) findMethodDeclaringClass(className, methodName string) string {
	class, exists := cg.ClassNameToAST[className]
	if !exists {
		return ""
	}

	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod && method.Name.Value == methodName {
			return className
		}
	}

	parent, exists := cg.ClassHierarchy[className]
	if exists && parent != "" {
		return cg.findMethodDeclaringClass(parent, methodName)
	}

	return ""
}

// findMethodByName returns a method's AST node
func (cg *CodeGenerator) findMethodByName(className, methodName string) *ast.Method {
	class, exists := cg.ClassNameToAST[className]
	if !exists {
		return nil
	}

	for _, feature := range class.Features {
		if method, isMethod := feature.(*ast.Method); isMethod && method.Name.Value == methodName {
			return method
		}
	}

	parent, exists := cg.ClassHierarchy[className]
	if exists && parent != "" {
		return cg.findMethodByName(parent, methodName)
	}

	return nil
}

// getLLVMTypeForCOOLType converts COOL types to LLVM types
func (cg *CodeGenerator) getLLVMTypeForCOOLType(coolType string) types.Type {
	switch coolType {
	case "Int":
		return types.I32
	case "Bool":
		return types.I1
	case "String":
		return types.NewPointer(types.I8)
	case "SELF_TYPE":
		// Default to Object if we can't determine the current context
		classType, exists := cg.TypeMap["Object"]
		if !exists {
			panic("Object type not found when processing SELF_TYPE")
		}
		return types.NewPointer(classType)
	case "Object", "IO":
		fallthrough
	default:
		classType, exists := cg.TypeMap[coolType]
		if !exists {
			panic(fmt.Sprintf("unknown type: %s", coolType))
		}
		return types.NewPointer(classType)
	}
}

// getLLVMTypeForCOOLTypeWithContext converts COOL types to LLVM types using class context
func (cg *CodeGenerator) getLLVMTypeForCOOLTypeWithContext(coolType string, currentClass string) types.Type {
	if coolType == "SELF_TYPE" {
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
	mainFunc := cg.Module.NewFunc("main", types.I32)
	entryBlock := mainFunc.NewBlock("entry")

	cg.CurrentFunc = mainFunc
	cg.CurrentBlock = entryBlock

	mainClass, exists := cg.TypeMap["Main"]
	if !exists {
		panic("Program must have a Main class")
	}

	cg.Symbols = make(map[string]value.Value)

	mainObj := cg.generateObjectAllocation("Main")

	mainObjAlloca := entryBlock.NewAlloca(types.NewPointer(mainClass))
	entryBlock.NewStore(mainObj, mainObjAlloca)
	cg.Symbols["self"] = mainObjAlloca

	cg.initializeAttributes("Main", mainObj)

	vtable, exists := cg.VTables["Main"]
	if !exists {
		panic("Main class must have a vtable")
	}

	mainMethodIndex, exists := cg.MethodIndices["Main"]["main"]
	if !exists {
		panic("Main class must have a main method")
	}

	methodSlotPtr := entryBlock.NewGetElementPtr(
		vtable.ContentType,
		vtable,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(mainMethodIndex)),
	)

	methodPtr := entryBlock.NewLoad(types.NewPointer(types.I8), methodSlotPtr)

	funcType := types.NewPointer(types.NewFunc(types.NewPointer(types.I8), types.NewPointer(mainClass)))
	castedMethodPtr := entryBlock.NewBitCast(methodPtr, funcType)

	entryBlock.NewCall(castedMethodPtr, mainObj)

	entryBlock.NewRet(constant.NewInt(types.I32, 0))
}

// generateObjectIdentifier creates LLVM IR to access a variable by its identifier
func (cg *CodeGenerator) generateObjectIdentifier(identifier *ast.ObjectIdentifier) value.Value {
	if identifier.Value == "self" {
		return cg.Symbols["self"]
	}

	block := cg.CurrentBlock

	val, exists := cg.Symbols[identifier.Value]
	if exists {
		if _, isLocalVar := val.(*ir.InstAlloca); isLocalVar {
			load := block.NewLoad(val.Type().(*types.PointerType).ElemType, val)
			return load
		}
		return val
	}

	selfPtr, exists := cg.Symbols["self"]
	if !exists {
		panic("'self' not found in symbol table")
	}

	selfPtrType := selfPtr.Type().(*types.PointerType)
	structType := selfPtrType.ElemType.(*types.StructType)
	className := ""

	for name, typ := range cg.TypeMap {
		if typ == structType {
			className = name
			break
		}
	}

	if className == "" {
		panic("couldn't determine class name for self")
	}

	attrIndex, exists := cg.AttributeIndices[className][identifier.Value]
	if !exists {
		// Check if this identifier is a method in the class
		if _, methodExists := cg.MethodIndices[className][identifier.Value]; methodExists {
			return selfPtr
		}

		panic(fmt.Sprintf("undefined attribute in class %s: %s", className, identifier.Value))
	}

	attributeType := structType.Fields[attrIndex]

	attrPtr := block.NewGetElementPtr(structType, selfPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(attrIndex)))

	load := block.NewLoad(attributeType, attrPtr)
	return load
}

// generateAssignmentExpression creates LLVM IR for variable assignment
func (cg *CodeGenerator) generateAssignmentExpression(assign *ast.AssignmentExpression) value.Value {
	rhsValue := cg.generateExpression(assign.Expression)
	block := cg.CurrentBlock

	if target, exists := cg.Symbols[assign.Identifier.Value]; exists {
		if allocaInst, isLocalVar := target.(*ir.InstAlloca); isLocalVar {
			targetType := allocaInst.Type().(*types.PointerType).ElemType
			if !targetType.Equal(rhsValue.Type()) {
				rhsValue = block.NewBitCast(rhsValue, targetType)
			}
			block.NewStore(rhsValue, allocaInst)
		} else if _, isParam := target.(*ir.Param); isParam {
			panic("assignment to parameter not properly handled - parameters should have local storage")
		} else {
			panic(fmt.Sprintf("unsupported assignment target type: %T", target))
		}
	} else {
		selfPtr, exists := cg.Symbols["self"]
		if !exists {
			panic("'self' not found in symbol table")
		}

		selfPtrType := selfPtr.Type().(*types.PointerType)
		structType := selfPtrType.ElemType.(*types.StructType)
		className := ""

		for name, typ := range cg.TypeMap {
			if typ == structType {
				className = name
				break
			}
		}

		if className == "" {
			panic("couldn't determine class name for self")
		}

		attrIndex, exists := cg.AttributeIndices[className][assign.Identifier.Value]
		if !exists {
			panic(fmt.Sprintf("undefined attribute in class %s: %s", className, assign.Identifier.Value))
		}

		attrPtr := block.NewGetElementPtr(structType, selfPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(attrIndex)))

		attributeType := structType.Fields[attrIndex]
		if !attributeType.Equal(rhsValue.Type()) {
			rhsValue = block.NewBitCast(rhsValue, attributeType)
		}

		block.NewStore(rhsValue, attrPtr)
	}

	return rhsValue
}

// generateIfExpression creates LLVM IR for a conditional expression
func (cg *CodeGenerator) generateIfExpression(ifExpr *ast.IfExpression) value.Value {
	condValue := cg.generateExpression(ifExpr.Condition)

	if condValue.Type() != types.I1 {
		panic("condition in if expression must be of boolean type")
	}

	cg.IfCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.IfCounter)

	currentFunc := cg.CurrentFunc
	trueBlock := currentFunc.NewBlock("if.then" + counterSuffix)
	falseBlock := currentFunc.NewBlock("if.else" + counterSuffix)
	mergeBlock := currentFunc.NewBlock("if.end" + counterSuffix)

	cg.CurrentBlock.NewCondBr(condValue, trueBlock, falseBlock)

	cg.CurrentBlock = trueBlock
	trueValue := cg.generateExpression(ifExpr.Consequence)
	trueEndBlock := cg.CurrentBlock
	cg.CurrentBlock.NewBr(mergeBlock)

	cg.CurrentBlock = falseBlock
	falseValue := cg.generateExpression(ifExpr.Alternative)
	falseEndBlock := cg.CurrentBlock
	cg.CurrentBlock.NewBr(mergeBlock)

	cg.CurrentBlock = mergeBlock

	var resultType types.Type

	if trueValue.Type().Equal(falseValue.Type()) {
		resultType = trueValue.Type()
	} else {
		_, trueIsPtr := trueValue.Type().(*types.PointerType)
		_, falseIsPtr := falseValue.Type().(*types.PointerType)

		if trueValue.Type() == types.I1 && falseIsPtr {
			resultType = falseValue.Type()
			trueValue = cg.ensureProperValue(trueValue, resultType, trueEndBlock)
		} else if falseValue.Type() == types.I1 && trueIsPtr {
			resultType = trueValue.Type()
			falseValue = cg.ensureProperValue(falseValue, resultType, falseEndBlock)
		} else {
			resultType = types.NewPointer(types.I8)

			trueValue = cg.ensureProperValue(trueValue, resultType, trueEndBlock)
			falseValue = cg.ensureProperValue(falseValue, resultType, falseEndBlock)
		}
	}

	phi := cg.CurrentBlock.NewPhi(
		&ir.Incoming{X: trueValue, Pred: trueEndBlock},
		&ir.Incoming{X: falseValue, Pred: falseEndBlock},
	)

	phi.Typ = resultType

	return phi
}

// generateWhileExpression creates LLVM IR for a while loop expression
func (cg *CodeGenerator) generateWhileExpression(whileExpr *ast.WhileExpression) value.Value {
	cg.WhileCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.WhileCounter)

	currentFunc := cg.CurrentFunc
	condBlock := currentFunc.NewBlock("while.cond" + counterSuffix)
	bodyBlock := currentFunc.NewBlock("while.body" + counterSuffix)
	exitBlock := currentFunc.NewBlock("while.exit" + counterSuffix)

	cg.CurrentBlock.NewBr(condBlock)

	cg.CurrentBlock = condBlock

	condValue := cg.generateExpression(whileExpr.Condition)

	if condValue.Type() != types.I1 {
		panic("condition in while expression must be of boolean type")
	}

	cg.CurrentBlock.NewCondBr(condValue, bodyBlock, exitBlock)

	cg.CurrentBlock = bodyBlock
	cg.generateExpression(whileExpr.Body)

	cg.CurrentBlock.NewBr(condBlock)

	cg.CurrentBlock = exitBlock

	nullValue := constant.NewNull(types.NewPointer(types.I8))

	return nullValue
}

// generateBlockExpression creates LLVM IR for a block of expressions
func (cg *CodeGenerator) generateBlockExpression(blockExpr *ast.BlockExpression) value.Value {
	var lastValue value.Value

	for _, expr := range blockExpr.Expressions {
		lastValue = cg.generateExpression(expr)
	}

	if lastValue == nil {
		return constant.NewNull(types.NewPointer(types.I8))
	}

	return lastValue
}

// generateBinaryExpression creates LLVM IR for binary operations
func (cg *CodeGenerator) generateBinaryExpression(binExpr *ast.BinaryExpression) value.Value {
	block := cg.CurrentBlock

	leftValue := cg.generateExpression(binExpr.Left)
	rightValue := cg.generateExpression(binExpr.Right)

	switch binExpr.Operator {
	case "+":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewAdd(leftValue, rightValue)
		} else {
			panic("operands of '+' must be integers")
		}

	case "-":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewSub(leftValue, rightValue)
		} else {
			panic("operands of '-' must be integers")
		}

	case "*":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewMul(leftValue, rightValue)
		} else {
			panic("operands of '*' must be integers")
		}

	case "/":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewSDiv(leftValue, rightValue)
		} else {
			panic("operands of '/' must be integers")
		}

	case "<":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredSLT, leftValue, rightValue)
		} else {
			panic("operands of '<' must be integers")
		}

	case "<=":
		if leftValue.Type() == types.I32 && rightValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredSLE, leftValue, rightValue)
		} else {
			panic("operands of '<=' must be integers")
		}

	case "=":
		if !leftValue.Type().Equal(rightValue.Type()) {
			return constant.NewInt(types.I1, 0)
		}

		if leftValue.Type() == types.I32 {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		if leftValue.Type() == types.I1 {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		if _, isPtr := leftValue.Type().(*types.PointerType); isPtr {
			return block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}

		panic(fmt.Sprintf("equality comparison not implemented for type: %v", leftValue.Type()))
	}

	panic(fmt.Sprintf("unsupported binary operator: %s", binExpr.Operator))
}

// generateDotCallExpression creates LLVM IR for method calls on objects
func (cg *CodeGenerator) generateDotCallExpression(dotCall *ast.DotCallExpression) value.Value {
	objectValue := cg.generateExpression(dotCall.Object)

	block := cg.CurrentBlock

	caseResult, isCaseResult := cg.CaseResults[objectValue]
	if isCaseResult {
		objectValue = caseResult.Value
	}

	argValues := make([]value.Value, 0, len(dotCall.Arguments)+1)

	argValues = append(argValues, objectValue)

	for _, arg := range dotCall.Arguments {
		argValues = append(argValues, cg.generateExpression(arg))
	}

	if dotCall.Method.Value == "out_string" || dotCall.Method.Value == "out_int" {
		var isIOObject bool
		if ptrType, isPtr := objectValue.Type().(*types.PointerType); isPtr {
			if structType, isStruct := ptrType.ElemType.(*types.StructType); isStruct {
				isIOObject = structType.Name() == "IO"
			}
		}

		if isIOObject || cg.canCastTo(objectValue, "IO") {
			ioObjValue := objectValue
			if !isIOObject {
				ioType, exists := cg.TypeMap["IO"]
				if exists {
					ioObjValue = block.NewBitCast(objectValue, types.NewPointer(ioType))
				}
			}

			if dotCall.Method.Value == "out_string" && len(argValues) > 1 {
				var outStringFunc *ir.Func
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_string" {
						outStringFunc = f
						break
					}
				}

				if outStringFunc == nil {
					outStringFunc = cg.Module.NewFunc("IO.out_string", types.NewPointer(cg.TypeMap["IO"]),
						ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						ir.NewParam("str", types.NewPointer(types.I8)))
				}

				strArg := argValues[1]
				convertedStrArg := cg.ensureCorrectArgumentType(strArg, types.NewPointer(types.I8))

				callArgs := []value.Value{ioObjValue, convertedStrArg}
				block.NewCall(outStringFunc, callArgs...)

				return ioObjValue

			} else if dotCall.Method.Value == "out_int" && len(argValues) > 1 {
				var outIntFunc *ir.Func
				for _, f := range cg.Module.Funcs {
					if f.Name() == "IO.out_int" {
						outIntFunc = f
						break
					}
				}

				if outIntFunc == nil {
					outIntFunc = cg.Module.NewFunc("IO.out_int", types.NewPointer(cg.TypeMap["IO"]),
						ir.NewParam("self", types.NewPointer(cg.TypeMap["IO"])),
						ir.NewParam("n", types.I32))
				}

				intArg := argValues[1]
				convertedIntArg := cg.ensureCorrectArgumentType(intArg, types.I32)

				callArgs := []value.Value{ioObjValue, convertedIntArg}
				block.NewCall(outIntFunc, callArgs...)

				return ioObjValue
			}
		}
	}

	if dotCall.Type != nil {
		targetTypeName := dotCall.Type.Value

		return cg.generateStaticDispatch(
			objectValue,
			targetTypeName,
			dotCall.Method.Value,
			argValues[1:],
		)
	}

	return cg.generateDynamicDispatch(
		objectValue,
		dotCall.Method.Value,
		argValues[1:],
	)
}

// generateStaticDispatch creates LLVM IR for static dispatch (obj@Type.method())
func (cg *CodeGenerator) generateStaticDispatch(object value.Value, typeName string, methodName string, args []value.Value) value.Value {
	block := cg.CurrentBlock

	methodFuncName := fmt.Sprintf("%s.%s", typeName, methodName)
	var methodFunc *ir.Func

	for _, f := range cg.Module.Funcs {
		if f.Name() == methodFuncName {
			methodFunc = f
			break
		}
	}

	if methodFunc == nil {
		panic(fmt.Sprintf("method %s not found in class %s", methodName, typeName))
	}

	classType, exists := cg.TypeMap[typeName]
	if !exists {
		panic(fmt.Sprintf("class type %s not found", typeName))
	}

	targetMethod := cg.findMethodByName(typeName, methodName)
	if targetMethod == nil {
		panic(fmt.Sprintf("method %s not found in class %s for signature lookup", methodName, typeName))
	}

	castedObject := block.NewBitCast(object, types.NewPointer(classType))

	allArgs := make([]value.Value, 0, len(args)+1)
	allArgs = append(allArgs, castedObject)

	for i, arg := range args {
		if i < len(targetMethod.Formals) {
			paramType := cg.getLLVMTypeForCOOLType(targetMethod.Formals[i].TypeDecl.Value)
			convertedArg := cg.ensureCorrectArgumentType(arg, paramType)
			allArgs = append(allArgs, convertedArg)
		} else {
			allArgs = append(allArgs, arg)
		}
	}

	call := block.NewCall(methodFunc, allArgs...)

	return call
}

// generateLetExpression creates LLVM IR for let expressions
func (cg *CodeGenerator) generateLetExpression(letExpr *ast.LetExpression) value.Value {
	block := cg.CurrentBlock

	oldSymbols := make(map[string]value.Value)
	for k, v := range cg.Symbols {
		oldSymbols[k] = v
	}

	for _, binding := range letExpr.Bindings {
		varName := binding.Identifier.Value

		var varType types.Type

		switch binding.Type.Value {
		case "Int":
			varType = types.I32
		case "Bool":
			varType = types.I1
		case "String":
			varType = types.NewPointer(types.I8)
		default:
			classType, exists := cg.TypeMap[binding.Type.Value]
			if !exists {
				panic(fmt.Sprintf("unknown type in let binding: %s", binding.Type.Value))
			}
			varType = types.NewPointer(classType)
		}

		alloca := block.NewAlloca(varType)

		var initValue value.Value

		if binding.Init != nil {
			initValue = cg.generateExpression(binding.Init)

			if !initValue.Type().Equal(varType) {
				initValue = block.NewBitCast(initValue, varType)
			}
		} else {
			switch binding.Type.Value {
			case "Int":
				initValue = constant.NewInt(types.I32, 0)
			case "Bool":
				initValue = constant.NewInt(types.I1, 0)
			case "String":
				initValue = constant.NewGetElementPtr(
					cg.EmptyStringGlobal.ContentType,
					cg.EmptyStringGlobal,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, 0),
				)
			default:
				ptrType, ok := varType.(*types.PointerType)
				if !ok {
					panic(fmt.Sprintf("expected pointer type for class type variable, got: %v", varType))
				}
				initValue = constant.NewNull(ptrType)
			}
		}

		block.NewStore(initValue, alloca)

		cg.Symbols[varName] = alloca
	}

	bodyValue := cg.generateExpression(letExpr.In)

	cg.Symbols = oldSymbols

	return bodyValue
}

// generateNewExpression creates LLVM IR for object instantiation
func (cg *CodeGenerator) generateNewExpression(newExpr *ast.NewExpression) value.Value {
	typeName := newExpr.Type.Value

	switch typeName {
	case "Int":
		return constant.NewInt(types.I32, 0)

	case "Bool":
		return constant.NewInt(types.I1, 0)

	case "String":
		return constant.NewGetElementPtr(
			cg.EmptyStringGlobal.ContentType,
			cg.EmptyStringGlobal,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0),
		)

	default:
		return cg.generateObjectAllocation(typeName)
	}
}

// generateIsVoidExpression creates LLVM IR for checking if a reference is null
func (cg *CodeGenerator) generateIsVoidExpression(isVoidExpr *ast.IsVoidExpression) value.Value {
	block := cg.CurrentBlock

	exprValue := cg.generateExpression(isVoidExpr.Expression)

	switch exprValue.Type() {
	case types.I32, types.I1:
		return constant.NewInt(types.I1, 0)

	default:
		if ptrType, isPtr := exprValue.Type().(*types.PointerType); isPtr {
			nullVal := constant.NewNull(ptrType)
			return block.NewICmp(enum.IPredEQ, exprValue, nullVal)
		} else {
			panic(fmt.Sprintf("isvoid check not implemented for type: %v", exprValue.Type()))
		}
	}
}

// generateUnaryExpression creates LLVM IR for unary operations
func (cg *CodeGenerator) generateUnaryExpression(unaryExpr *ast.UnaryExpression) value.Value {
	block := cg.CurrentBlock

	exprValue := cg.generateExpression(unaryExpr.Right)

	switch unaryExpr.Operator {
	case "~":
		if exprValue.Type() != types.I32 {
			panic("operand of integer negation (~) must be an integer")
		}

		zero := constant.NewInt(types.I32, 0)
		return block.NewSub(zero, exprValue)

	case "not":
		if exprValue.Type() != types.I1 {
			panic("operand of boolean NOT (not) must be a boolean")
		}

		return block.NewXor(exprValue, constant.NewInt(types.I1, 1))

	default:
		panic(fmt.Sprintf("unsupported unary operator: %s", unaryExpr.Operator))
	}
}

// generateCaseExpression creates LLVM IR for COOL's case expressions
func (cg *CodeGenerator) generateCaseExpression(caseExpr *ast.CaseExpression) value.Value {
	cg.CaseCounter++
	counterSuffix := fmt.Sprintf(".%d", cg.CaseCounter)

	currentFunc := cg.CurrentFunc

	exprValue := cg.generateExpression(caseExpr.Expression)

	endBlock := currentFunc.NewBlock("case.end" + counterSuffix)

	branchBlocks := make([]*ir.Block, len(caseExpr.Branches))
	for i := range caseExpr.Branches {
		branchBlocks[i] = currentFunc.NewBlock(fmt.Sprintf("case.branch.%d%s", i, counterSuffix))
	}

	noMatchBlock := currentFunc.NewBlock("case.nomatch" + counterSuffix)

	currentBlock := cg.CurrentBlock

	if ptrType, isPtr := exprValue.Type().(*types.PointerType); isPtr {
		notNullBlock := currentFunc.NewBlock("case.notnull" + counterSuffix)

		nullVal := constant.NewNull(ptrType)
		isNull := currentBlock.NewICmp(enum.IPredEQ, exprValue, nullVal)

		currentBlock.NewCondBr(isNull, noMatchBlock, notNullBlock)

		currentBlock = notNullBlock
		cg.CurrentBlock = notNullBlock
	}

	typeCheckBlock := currentFunc.NewBlock("case.typecheck" + counterSuffix)
	currentBlock.NewBr(typeCheckBlock)
	cg.CurrentBlock = typeCheckBlock

	objectType := cg.getObjectRuntimeType(exprValue)

	branchingBlock := typeCheckBlock

	decisionBlocks := make([]*ir.Block, len(caseExpr.Branches))
	for i := range caseExpr.Branches {
		decisionBlocks[i] = currentFunc.NewBlock(fmt.Sprintf("case.decision.%d%s", i, counterSuffix))
	}

	branchingBlock.NewBr(decisionBlocks[0])

	for i, branch := range caseExpr.Branches {
		branchingBlock = decisionBlocks[i]

		branchType := branch.Type.Value

		var matchesCondition value.Value

		if branchType == objectType {
			matchesCondition = constant.NewInt(types.I1, 1)
		} else {
			matchesCondition = constant.NewInt(types.I1, 0)

			currentType := objectType
			for {
				parent, exists := cg.ClassHierarchy[currentType]
				if !exists || parent == "" {
					break
				}

				if parent == branchType {
					matchesCondition = constant.NewInt(types.I1, 1)
					break
				}

				currentType = parent
			}
		}

		var nextBlock *ir.Block
		if i < len(caseExpr.Branches)-1 {
			nextBlock = decisionBlocks[i+1]
		} else {
			nextBlock = noMatchBlock
		}

		branchingBlock.NewCondBr(matchesCondition, branchBlocks[i], nextBlock)
	}

	noMatchBlock.NewCall(cg.StdlibFuncs["exit"], constant.NewInt(types.I32, 1))
	noMatchBlock.NewUnreachable()

	branchValues := make([]value.Value, len(caseExpr.Branches))
	branchRealTypes := make([]string, len(caseExpr.Branches))
	branchEndBlocks := make([]*ir.Block, len(caseExpr.Branches))

	for i, branch := range caseExpr.Branches {
		cg.CurrentBlock = branchBlocks[i]

		oldSymbols := make(map[string]value.Value)
		for k, v := range cg.Symbols {
			oldSymbols[k] = v
		}

		var castedValue value.Value
		if branchType, exists := cg.TypeMap[branch.Type.Value]; exists {
			_, isExprPtr := exprValue.Type().(*types.PointerType)
			if !isExprPtr && (exprValue.Type() == types.I32 || exprValue.Type() == types.I1) {
				tmpPtr := cg.CurrentBlock.NewIntToPtr(exprValue, types.NewPointer(types.I8))
				castedValue = cg.CurrentBlock.NewBitCast(tmpPtr, types.NewPointer(branchType))
			} else {
				castedValue = cg.CurrentBlock.NewBitCast(exprValue, types.NewPointer(branchType))
			}
		} else {
			castedValue = exprValue
		}

		cg.Symbols[branch.Identifier.Value] = castedValue

		branchValues[i] = cg.generateExpression(branch.Expression)
		branchRealTypes[i] = branch.Type.Value

		cg.Symbols = oldSymbols

		branchEndBlocks[i] = cg.CurrentBlock

		branchEndBlocks[i].NewBr(endBlock)
	}

	cg.CurrentBlock = endBlock

	var resultType types.Type
	if len(branchValues) > 0 {
		resultType = branchValues[0].Type()

		for _, val := range branchValues[1:] {
			if !val.Type().Equal(resultType) {
				resultType = types.NewPointer(types.I8)
				break
			}
		}
	} else {
		resultType = types.NewPointer(types.I8)
	}

	phi := &ir.InstPhi{Typ: resultType}
	endBlock.Insts = append(endBlock.Insts, phi)

	for i, val := range branchValues {
		if !val.Type().Equal(resultType) {
			_, valIsPtr := val.Type().(*types.PointerType)
			_, resultIsPtr := resultType.(*types.PointerType)

			if !valIsPtr && resultIsPtr && (val.Type() == types.I1 || val.Type() == types.I32) {
				tmpPtr := branchEndBlocks[i].NewIntToPtr(val, types.NewPointer(types.I8))
				val = branchEndBlocks[i].NewBitCast(tmpPtr, resultType)
			} else {
				val = cg.ensureCorrectArgumentType(val, resultType)
			}
		}

		phi.Incs = append(phi.Incs, &ir.Incoming{X: val, Pred: branchEndBlocks[i]})
	}

	result := &CaseResult{
		Value:       phi,
		BranchTypes: branchRealTypes,
	}

	if cg.CaseResults == nil {
		cg.CaseResults = make(map[value.Value]*CaseResult)
	}
	cg.CaseResults[phi] = result

	return phi
}

// generateCallExpression creates LLVM IR for function calls
func (cg *CodeGenerator) generateCallExpression(callExpr *ast.CallExpression) value.Value {
	if cg.CurrentFunc == nil {
		panic("Cannot call method without object outside of method context")
	}

	selfObj := cg.CurrentFunc.Params[0]

	methodIdent, ok := callExpr.Function.(*ast.ObjectIdentifier)
	if !ok {
		panic(fmt.Sprintf("Unexpected Function type in CallExpression: %T", callExpr.Function))
	}

	methodName := methodIdent.Value

	args := make([]value.Value, 0, len(callExpr.Arguments))
	for _, arg := range callExpr.Arguments {
		args = append(args, cg.generateExpression(arg))
	}

	return cg.generateDynamicDispatch(selfObj, methodName, args)
}

// DefineBuiltInClasses defines the built-in COOL classes: Object, IO, Int, String, Bool
func (cg *CodeGenerator) DefineBuiltInClasses() {
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

// ensureProperValue ensures that value is properly typed for the expected type
// particularly handling the case of integer 0 -> null pointer conversion
func (cg *CodeGenerator) ensureProperValue(val value.Value, expectedType types.Type, block *ir.Block) value.Value {
	// If the value is already of the expected type, return it
	if val.Type().Equal(expectedType) {
		return val
	}

	// Check if we're trying to use an integer in a pointer context
	_, expectedIsPtr := expectedType.(*types.PointerType)

	// Special handling for integer constant 0 -> null pointer
	if expectedIsPtr && val.Type() == types.I1 || val.Type() == types.I8 || val.Type() == types.I32 || val.Type() == types.I64 {
		if constVal, isConst := val.(constant.Constant); isConst {
			// For all integer types, check if the value is 0
			if constInt, isInt := constVal.(*constant.Int); isInt && constInt.X.Int64() == 0 {
				// Create a proper null pointer of the expected type
				if ptrType, ok := expectedType.(*types.PointerType); ok {
					return constant.NewNull(ptrType)
				}
			}
		}
	}

	// Handle pointer to pointer conversions (bitcast)
	_, valIsPtr := val.Type().(*types.PointerType)
	if expectedIsPtr && valIsPtr {
		return block.NewBitCast(val, expectedType)
	}

	// Handle integer to pointer conversions
	if expectedIsPtr && !valIsPtr {
		// First convert to i8* then bitcast if needed
		i8PtrType := types.NewPointer(types.I8)
		ptrVal := block.NewIntToPtr(val, i8PtrType)

		// If the expected type is not i8*, bitcast to the correct type
		if !expectedType.Equal(i8PtrType) {
			return block.NewBitCast(ptrVal, expectedType)
		}
		return ptrVal
	}

	// Handle pointer to integer conversions
	if !expectedIsPtr && valIsPtr {
		// If converting pointer to integer, use ptrtoint
		intType := types.I64 // Use a reasonable default
		if intTy, isInt := expectedType.(*types.IntType); isInt {
			intType = intTy
		}
		return block.NewPtrToInt(val, intType)
	}

	// For all other cases (like int to int conversions), use the appropriate cast
	if intTy, isInt := expectedType.(*types.IntType); isInt {
		if valIntTy, valIsInt := val.Type().(*types.IntType); valIsInt {
			// Integer size conversions
			if intTy.BitSize > valIntTy.BitSize {
				return block.NewZExt(val, expectedType)
			} else if intTy.BitSize < valIntTy.BitSize {
				return block.NewTrunc(val, expectedType)
			}
		}
	}

	fmt.Printf("Warning: Couldn't convert value of type %v to expected type %v\n",
		val.Type(), expectedType)
	return val
}
