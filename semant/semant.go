package semant

import (
	"cool-compiler/ast"
	"cool-compiler/lexer"
	"fmt"
)

type SymbolTable struct {
	symbols map[string]*SymbolEntry
	parent  *SymbolTable
}

type SymbolEntry struct {
	Type     string // This could be "Class", "Method", or "Attribute". Otherwise, it is the type of a variable (e.g. if it's a method parameter)
	Token    lexer.Token
	AttrType *ast.TypeIdentifier // If it's an attribute
	Method   *ast.Method         // If it's a method
	Scope    *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*SymbolEntry),
		parent:  parent,
	}
}

func (st *SymbolTable) AddEntry(name string, entry *SymbolEntry) {
	st.symbols[name] = entry
}

func (st *SymbolTable) Lookup(name string) (*SymbolEntry, bool) {
	entry, ok := st.symbols[name]
	if !ok && st.parent != nil {
		return st.parent.Lookup(name)
	}
	return entry, ok
}

type SemanticAnalyser struct {
	globalSymbolTable   *SymbolTable
	errors              []string
	inheritanceGraph    *InheritanceGraph
	currentClass        string
	objectClassSymEntry *SymbolEntry
}

func NewSemanticAnalyser() *SemanticAnalyser {
	return &SemanticAnalyser{
		globalSymbolTable: NewSymbolTable(nil),
		errors:            []string{},
		inheritanceGraph:  NewInheritanceGraph(),
	}
}

func (sa *SemanticAnalyser) Errors() []string {
	return sa.errors
}

func (sa *SemanticAnalyser) GetSymbolTable() *SymbolTable {
	return sa.globalSymbolTable
}
func (sa *SemanticAnalyser) GetInheritanceGraph() *InheritanceGraph {
	return sa.inheritanceGraph
}

func (sa *SemanticAnalyser) Analyze(program *ast.Program) {
	sa.buildClassesSymboltables(program)
	sa.buildInheritanceGraph(program)
	sa.buildSymbolTables(program)
	sa.validateMainClass()
	sa.typeCheck(program)
}

func (sa *SemanticAnalyser) typeCheck(program *ast.Program) {
	for _, class := range program.Classes {
		st := sa.globalSymbolTable.symbols[class.Name.Value].Scope
		sa.typeCheckClass(class, st)
	}
}

func (sa *SemanticAnalyser) typeCheckClass(cls *ast.Class, st *SymbolTable) {
	sa.currentClass = cls.Name.Value
	defer func() { sa.currentClass = "" }()

	for _, feature := range cls.Features {
		switch f := feature.(type) {
		case *ast.Attribute:
			sa.typeCheckAttribute(f, st)
		case *ast.Method:
			sa.typeCheckMethod(f, st)
		}
	}
}

func (sa *SemanticAnalyser) typeCheckAttribute(attribute *ast.Attribute, st *SymbolTable) {
	// Check if attribute is named 'self'
	if attribute.Name.Value == "self" {
		sa.errors = append(sa.errors, "cannot have attribute named 'self'")
		return
	}

	// Check if attribute type exists
	if _, ok := sa.globalSymbolTable.Lookup(attribute.TypeDecl.Value); !ok {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s for attribute %s", attribute.TypeDecl.Value, attribute.Name.Value))
		return
	}

	if attribute.Expression != nil {
		expressionType := sa.getExpressionType(attribute.Expression, st)
		if !sa.isTypeConformant(expressionType, attribute.TypeDecl.Value) {
			sa.errors = append(sa.errors, fmt.Sprintf("attribute %s cannot be of type %s, expected %s", attribute.Name.Value, expressionType, attribute.TypeDecl.Value))
		}
	}
}

func (sa *SemanticAnalyser) typeCheckMethod(method *ast.Method, st *SymbolTable) {

	methodSt := st.symbols[method.Name.Value].Scope
	for _, formal := range method.Formals {
		if formal.Name.Value == "self" {
			sa.errors = append(sa.errors, "cannot use 'self' as formal parameter")
			continue
		}

		if formal.TypeDecl.Value == "SELF_TYPE" {
			sa.errors = append(sa.errors, "SELF_TYPE is not allowed as formal parameter type")
			continue
		}

		// Check if formal type exists
		if _, ok := sa.globalSymbolTable.Lookup(formal.TypeDecl.Value); !ok {
			sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s in formal parameter of method %s", formal.TypeDecl.Value, method.Name.Value))
			continue
		}

		if _, ok := methodSt.Lookup(formal.Name.Value); ok {
			sa.errors = append(sa.errors, fmt.Sprintf("argument %s in method %s is already defined", formal.Name.Value, method.Name.Value))
			continue
		}

		methodSt.AddEntry(formal.Name.Value, &SymbolEntry{Token: formal.Token, Type: formal.TypeDecl.Value})
	}

	// Check if return type exists
	if method.TypeDecl.Value != "SELF_TYPE" {
		if _, ok := sa.globalSymbolTable.Lookup(method.TypeDecl.Value); !ok {
			sa.errors = append(sa.errors, fmt.Sprintf("undefined return type %s for method %s", method.TypeDecl.Value, method.Name.Value))
			return
		}
	}

	methodExpressionType := sa.getExpressionType(method.Expression, methodSt)
	expectedReturnType := method.TypeDecl.Value
	if expectedReturnType == "SELF_TYPE" {
		expectedReturnType = sa.currentClass
	}

	if !sa.isTypeConformant(methodExpressionType, expectedReturnType) {
		sa.errors = append(sa.errors, fmt.Sprintf("method %s is expected to return %s, found %s", method.Name.Value, method.TypeDecl.Value, methodExpressionType))
	}
}

func (sa *SemanticAnalyser) isTypeConformant(type1, type2 string) bool {
	return sa.inheritanceGraph.IsConformant(type1, type2)
}

func (sa *SemanticAnalyser) getExpressionType(expression ast.Expression, st *SymbolTable) string {
	switch e := expression.(type) {
	case *ast.IntegerLiteral:
		return "Int"
	case *ast.StringLiteral:
		return "String"
	case *ast.BooleanLiteral:
		return "Bool"
	case *ast.ObjectIdentifier:
		return sa.getObjectIdentifierType(e, st)
	case *ast.BlockExpression:
		return sa.getBlockExpressionType(e, st)
	case *ast.IfExpression:
		return sa.getIfExpressionType(e, st)
	case *ast.WhileExpression:
		return sa.getWhileExpressionType(e, st)
	case *ast.NewExpression:
		return sa.GetNewExpressionType(e, st)
	case *ast.LetExpression:
		return sa.GetLetExpressionType(e, st)
	case *ast.AssignmentExpression:
		return sa.GetAssignmentExpressionType(e, st)
	case *ast.UnaryExpression:
		return sa.GetUnaryExpressionType(e, st)
	case *ast.BinaryExpression:
		return sa.GetBinaryExpressionType(e, st)
	case *ast.CaseExpression:
		return sa.GetCaseExpressionType(e, st)
	case *ast.CallExpression:
		return sa.getCallExpressionType(e, st)
	case *ast.DotCallExpression:
		return sa.getDotCallExpressionType(e, st)
	case *ast.IsVoidExpression:
		return "Bool"
	default:
		return "Object"
	}
}

func (sa *SemanticAnalyser) getObjectIdentifierType(identifier *ast.ObjectIdentifier, st *SymbolTable) string {
	entry, ok := st.Lookup(identifier.Value)
	if !ok {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined identifier %s", identifier.Value))
		return "Object"
	}
	if entry.Type == "SELF_TYPE" {
		return sa.currentClass
	}
	return entry.Type
}

func (sa *SemanticAnalyser) getBlockExpressionType(bexpr *ast.BlockExpression, st *SymbolTable) string {
	lastType := ""

	// Go recursively to check all inner expressions
	for _, expression := range bexpr.Expressions {
		lastType = sa.getExpressionType(expression, st)
	}

	return lastType
}

func (sa *SemanticAnalyser) getIfExpressionType(ifexpr *ast.IfExpression, st *SymbolTable) string {
	conditionType := sa.getExpressionType(ifexpr.Condition, st)
	if conditionType != "Bool" {
		sa.errors = append(sa.errors, fmt.Sprintf("condition of if statement is of type %s, expected Bool", conditionType))
		return "Object"
	}

	constype := sa.getExpressionType(ifexpr.Consequence, st)
	alttype := sa.getExpressionType(ifexpr.Alternative, st)

	return sa.inheritanceGraph.FindLCA(constype, alttype)
}

func (sa *SemanticAnalyser) getWhileExpressionType(wexpr *ast.WhileExpression, st *SymbolTable) string {
	conditionType := sa.getExpressionType(wexpr.Condition, st)
	if conditionType != "Bool" {
		sa.errors = append(sa.errors, fmt.Sprintf("condition of if statement is of type %s, expected Bool", conditionType))
		return "Object"
	}

	return sa.getExpressionType(wexpr.Body, st)
}

func (sa *SemanticAnalyser) buildClassesSymboltables(program *ast.Program) {
	// Add built-in classes with their methods
	sa.initializeObjectClass()
	sa.initializeIOClass()
	sa.initializeStringClass()
	sa.initializeIntClass()
	sa.initializeBoolClass()

	for _, class := range program.Classes {
		// Only check for duplicate classes, not parent relationships
		if _, ok := sa.globalSymbolTable.Lookup(class.Name.Value); ok {
			sa.errors = append(sa.errors, fmt.Sprintf("class %s is already defined", class.Name.Value))
			continue
		}

		// Create scope with nil parent temporarily
		classScope := NewSymbolTable(nil)
		sa.globalSymbolTable.AddEntry(class.Name.Value, &SymbolEntry{
			Type:  "Class",
			Token: class.Name.Token,
			Scope: classScope,
		})
		sa.inheritanceGraph.AddNode(class.Name.Value, class)
	}
}

func (sa *SemanticAnalyser) initializeObjectClass() {
	objectScope := NewSymbolTable(nil)

	// Add Object methods
	objectScope.AddEntry("abort", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "Object"},
		},
	})
	objectScope.AddEntry("type_name", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "String"},
		},
	})
	objectScope.AddEntry("copy", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
		},
	})

	objectSymbolEntry := &SymbolEntry{
		Type:  "Class",
		Token: lexer.Token{Literal: "Object"},
		Scope: objectScope,
	}

	sa.globalSymbolTable.AddEntry("Object", objectSymbolEntry)
	sa.objectClassSymEntry = objectSymbolEntry

	sa.inheritanceGraph.AddNode("Object", &ast.Class{Name: &ast.TypeIdentifier{Value: "Object"}})
}

func (sa *SemanticAnalyser) initializeIOClass() {
	ioScope := NewSymbolTable(sa.objectClassSymEntry.Scope) // Set parent to Object's scope

	// Add IO methods
	ioScope.AddEntry("out_string", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			Formals: []*ast.Formal{
				{Name: &ast.ObjectIdentifier{Value: "x"}, TypeDecl: &ast.TypeIdentifier{Value: "String"}},
			},
			TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
		},
	})
	ioScope.AddEntry("out_int", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			Formals: []*ast.Formal{
				{Name: &ast.ObjectIdentifier{Value: "x"}, TypeDecl: &ast.TypeIdentifier{Value: "Int"}},
			},
			TypeDecl: &ast.TypeIdentifier{Value: "SELF_TYPE"},
		},
	})
	ioScope.AddEntry("in_string", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "String"},
		},
	})
	ioScope.AddEntry("in_int", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "Int"},
		},
	})

	sa.globalSymbolTable.AddEntry("IO", &SymbolEntry{
		Type:  "Class",
		Token: lexer.Token{Literal: "IO"},
		Scope: ioScope,
	})

	sa.inheritanceGraph.AddNode("IO", &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "IO"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
	})
	sa.inheritanceGraph.AddEdge("IO", "Object")
}

func (sa *SemanticAnalyser) initializeStringClass() {
	stringScope := NewSymbolTable(sa.objectClassSymEntry.Scope) // Set parent to Object's scope

	// Add String methods
	stringScope.AddEntry("length", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			TypeDecl: &ast.TypeIdentifier{Value: "Int"},
		},
	})
	stringScope.AddEntry("concat", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			Formals: []*ast.Formal{
				{Name: &ast.ObjectIdentifier{Value: "s"}, TypeDecl: &ast.TypeIdentifier{Value: "String"}},
			},
			TypeDecl: &ast.TypeIdentifier{Value: "String"},
		},
	})
	stringScope.AddEntry("substr", &SymbolEntry{
		Type: "Method",
		Method: &ast.Method{
			Formals: []*ast.Formal{
				{Name: &ast.ObjectIdentifier{Value: "i"}, TypeDecl: &ast.TypeIdentifier{Value: "Int"}},
				{Name: &ast.ObjectIdentifier{Value: "l"}, TypeDecl: &ast.TypeIdentifier{Value: "Int"}},
			},
			TypeDecl: &ast.TypeIdentifier{Value: "String"},
		},
	})

	sa.globalSymbolTable.AddEntry("String", &SymbolEntry{
		Type:  "Class",
		Token: lexer.Token{Literal: "String"},
		Scope: stringScope,
	})

	sa.inheritanceGraph.AddNode("String", &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "String"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
	})
	sa.inheritanceGraph.AddEdge("String", "Object")
}

func (sa *SemanticAnalyser) initializeIntClass() {
	// Int has no methods but needs a scope
	sa.globalSymbolTable.AddEntry("Int", &SymbolEntry{
		Type:  "Class",
		Token: lexer.Token{Literal: "Int"},
		Scope: NewSymbolTable(sa.objectClassSymEntry.Scope), // Set parent to Object's scope
	})

	sa.inheritanceGraph.AddNode("Int", &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "Int"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
	})
	sa.inheritanceGraph.AddEdge("Int", "Object")
}

func (sa *SemanticAnalyser) initializeBoolClass() {
	// Bool has no methods but needs a scope
	sa.globalSymbolTable.AddEntry("Bool", &SymbolEntry{
		Type:  "Class",
		Token: lexer.Token{Literal: "Bool"},
		Scope: NewSymbolTable(sa.objectClassSymEntry.Scope), // Set parent to Object's scope
	})

	sa.inheritanceGraph.AddNode("Bool", &ast.Class{
		Name:   &ast.TypeIdentifier{Value: "Bool"},
		Parent: &ast.TypeIdentifier{Value: "Object"},
	})
	sa.inheritanceGraph.AddEdge("Bool", "Object")
}

func (sa *SemanticAnalyser) buildSymbolTables(program *ast.Program) {
	for _, class := range program.Classes {
		classEntry, _ := sa.globalSymbolTable.Lookup(class.Name.Value)
		for _, feature := range class.Features {
			switch f := feature.(type) {
			case *ast.Attribute:
				if sa.isAttributeRedefined(f.Name.Value, classEntry.Scope.parent) {
					sa.errors = append(sa.errors, fmt.Sprintf("attribute %s is already defined in a parent class of %s", f.Name.Value, class.Name.Value))
					continue
				}
				if _, ok := classEntry.Scope.Lookup(f.Name.Value); ok {
					sa.errors = append(sa.errors, fmt.Sprintf("attribute %s is already defined in class %s", f.Name.Value, class.Name.Value))
					continue
				}
				classEntry.Scope.AddEntry(f.Name.Value, &SymbolEntry{Type: "Attribute", Token: f.Name.Token, AttrType: f.TypeDecl})
			case *ast.Method:
				methodST := NewSymbolTable(classEntry.Scope)

				redefined, parentEntry := sa.isMethodRedefined(f.Name.Value, classEntry.Scope.parent)
				if redefined {
					sa.validateMethodOverride(f, parentEntry.Method)
					classEntry.Scope.AddEntry(f.Name.Value, &SymbolEntry{Type: "Method", Token: f.Name.Token, Scope: methodST, Method: f})
					continue
				}

				if _, ok := classEntry.Scope.Lookup(f.Name.Value); ok {
					sa.errors = append(sa.errors, fmt.Sprintf("method %s is already defined in class %s", f.Name.Value, class.Name.Value))
					continue
				}
				classEntry.Scope.AddEntry(f.Name.Value, &SymbolEntry{Type: "Method", Token: f.Name.Token, Scope: methodST, Method: f})
			}
		}
	}
}

func (sa *SemanticAnalyser) GetNewExpressionType(ne *ast.NewExpression, st *SymbolTable) string {
	typeName := ne.Type.Value
	if typeName == "SELF_TYPE" {
		if sa.currentClass == "" {
			sa.errors = append(sa.errors, "SELF_TYPE used outside of class context")
			return "Object"
		}
		typeName = sa.currentClass
	}

	if _, ok := sa.globalSymbolTable.Lookup(typeName); !ok {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s in new expression", typeName))
		return "Object"
	}
	return typeName
}

func (sa *SemanticAnalyser) checkLetBinding(binding *ast.Binding, st *SymbolTable) string {
	declaredType := binding.Type.Value

	// Handle SELF_TYPE
	if declaredType == "SELF_TYPE" {
		if sa.currentClass == "" {
			sa.errors = append(sa.errors, "SELF_TYPE used outside of class context")
			return "Object"
		}
		declaredType = sa.currentClass
	}

	// Check if declaredType is a valid type
	if _, ok := sa.globalSymbolTable.Lookup(declaredType); !ok {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s in let binding", declaredType))
		declaredType = "Object" // Recover by assuming Object type
	}

	if binding.Init != nil {
		initExprType := sa.getExpressionType(binding.Init, st)
		if !sa.isTypeConformant(initExprType, declaredType) {
			sa.errors = append(sa.errors, fmt.Sprintf("Let binding with wrong type: variable %s of type %s initialized with type %s", binding.Identifier.Value, declaredType, initExprType))
		}
	}

	return declaredType
}

func (sa *SemanticAnalyser) GetLetExpressionType(letExpr *ast.LetExpression, st *SymbolTable) string {
	letScope := NewSymbolTable(st) // Create a new scope, child of the current scope 'st'

	for _, binding := range letExpr.Bindings {
		declaredType := sa.checkLetBinding(binding, st)

		// Add variable to the *letScope*
		letScope.AddEntry(binding.Identifier.Value, &SymbolEntry{
			Type:  declaredType,
			Token: binding.Identifier.Token,
		})
	}

	// Type-check the 'in' expression in the *letScope*
	return sa.getExpressionType(letExpr.In, letScope)
}

func (sa *SemanticAnalyser) GetAssignmentExpressionType(a *ast.AssignmentExpression, st *SymbolTable) string {
	// Prevent assignment to 'self'
	if a.Identifier.Value == "self" {
		sa.errors = append(sa.errors, "cannot assign to 'self'")
		return "Object"
	}

	// Check local scope first
	se, ok := st.Lookup(a.Identifier.Value)
	if !ok {
		// Check if it's a class attribute
		if classEntry, ok := sa.globalSymbolTable.Lookup(sa.currentClass); ok {
			if attrEntry, ok := classEntry.Scope.Lookup(a.Identifier.Value); ok && attrEntry.Type == "Attribute" {
				se = attrEntry
				ok = true
			}
		}
	}

	if !ok {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined identifier %s in assignment", a.Identifier.Value))
		return "Object"
	}

	// Determine declared type (attribute vs local variable)
	var declaredType string
	if se.Type == "Attribute" {
		declaredType = se.AttrType.Value
	} else {
		declaredType = se.Type
	}

	exprType := sa.getExpressionType(a.Expression, st)
	if !sa.isTypeConformant(exprType, declaredType) {
		sa.errors = append(sa.errors, fmt.Sprintf("type mismatch in assignment: variable '%s' has type %s but was assigned value of type %s",
			a.Identifier.Value, declaredType, exprType))
	}

	return exprType
}

func (sa *SemanticAnalyser) GetUnaryExpressionType(uexpr *ast.UnaryExpression, st *SymbolTable) string {
	rightType := sa.getExpressionType(uexpr.Right, st)
	switch uexpr.Operator {
	case "~":
		if rightType != "Int" {
			sa.errors = append(sa.errors, fmt.Sprintf("bitwise negation on non-Int type: %s", rightType))
		}
		return "Int"
	case "not":
		if rightType != "Bool" {
			sa.errors = append(sa.errors, fmt.Sprintf("logical negation on non-Bool type: %s", rightType))
		}
		return "Bool"
	default:
		sa.errors = append(sa.errors, fmt.Sprintf("unknown unary operator %s", uexpr.Operator))
		return "Object"
	}
}

func isComparable(t string) bool {
	return t == "Int" || t == "Bool" || t == "String"
}

func (sa *SemanticAnalyser) GetBinaryExpressionType(be *ast.BinaryExpression, st *SymbolTable) string {
	leftType := sa.getExpressionType(be.Left, st)
	rightType := sa.getExpressionType(be.Right, st)
	switch be.Operator {
	case "+", "*", "/", "-":
		if leftType != "Int" || rightType != "Int" {
			sa.errors = append(sa.errors, fmt.Sprintf("arithmetic operation on non-Int types: %s %s %s", leftType, be.Operator, rightType))
		}
		return "Int"
	case "<", "<=", "=":
		if leftType != rightType || !isComparable(leftType) {
			sa.errors = append(sa.errors, fmt.Sprintf("comparison between incompatible types: %s %s %s", leftType, be.Operator, rightType))
		}
		return "Bool"
	default:
		sa.errors = append(sa.errors, fmt.Sprintf("unknown binary operator %s", be.Operator))
		return "Object"
	}
}

func (sa *SemanticAnalyser) GetCaseExpressionType(ce *ast.CaseExpression, st *SymbolTable) string {
	// Check minimum branches requirement
	if len(ce.Branches) == 0 {
		sa.errors = append(sa.errors, "case expression must have at least one branch")
		return "Object"
	}

	// Get and resolve case expression type
	exprType := sa.getExpressionType(ce.Expression, st)
	if exprType == "SELF_TYPE" {
		exprType = sa.currentClass
	}

	var caseTypes []string
	seenTypes := make(map[string]bool)

	for _, branch := range ce.Branches {
		// Check for duplicate types in branches
		if seenTypes[branch.Type.Value] {
			sa.errors = append(sa.errors, fmt.Sprintf("duplicate branch type %s in case expression", branch.Type.Value))
			continue
		}
		seenTypes[branch.Type.Value] = true

		// Verify branch type exists
		if _, ok := sa.globalSymbolTable.Lookup(branch.Type.Value); !ok {
			sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s in case branch", branch.Type.Value))
			continue
		}

		// Create branch-specific scope
		branchSt := NewSymbolTable(st)
		branchSt.AddEntry(branch.Identifier.Value, &SymbolEntry{
			Type:  branch.Type.Value,
			Token: branch.Identifier.Token,
		})

		// Check branch expression with proper scope
		branchExprType := sa.getExpressionType(branch.Expression, branchSt)
		if branchExprType == "SELF_TYPE" {
			branchExprType = sa.currentClass
		}

		caseTypes = append(caseTypes, branchExprType)
	}

	if len(caseTypes) == 0 {
		return "Object"
	}

	// Find Least Common Ancestor of all case types
	lca := caseTypes[0]
	for _, ct := range caseTypes[1:] {
		lca = sa.inheritanceGraph.FindLCA(lca, ct)
	}
	return lca
}

func (sa *SemanticAnalyser) getCallExpressionType(call *ast.CallExpression, st *SymbolTable) string {
	// Dynamic dispatch:  ID(arg1, arg2, ...)

	// Function name is in call.Function, which is an Expression.
	// For simple dynamic dispatch, it will be an ObjectIdentifier

	methodIdentifier, ok := call.Function.(*ast.ObjectIdentifier)
	if !ok {
		sa.errors = append(sa.errors, "invalid method call syntax") // More specific error needed
		return "Object"
	}
	methodName := methodIdentifier.Value

	// 1. Get the type of 'self' - in dynamic dispatch, it's 'SELF_TYPE' which resolves to currentClass
	selfType := sa.currentClass // In dynamic dispatch, 'self' is implicit

	// 2. Lookup the method in the symbol table of the class 'selfType'
	classEntry, classFound := sa.globalSymbolTable.Lookup(selfType)
	if !classFound {
		// This should not happen in valid COOL programs as currentClass is always valid
		sa.errors = append(sa.errors, fmt.Sprintf("class %s not found in symbol table", selfType))
		return "Object"
	}
	methodScope := classEntry.Scope // Scope of the class
	methodEntry, methodFound := methodScope.Lookup(methodName)

	if !methodFound || methodEntry.Type != "Method" {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined method %s in class %s", methodName, selfType))
		return "Object"
	}

	methodDef := methodEntry.Method // Retrieve the method definition from SymbolEntry

	// 3. Type check arguments
	if len(call.Arguments) != len(methodDef.Formals) {
		sa.errors = append(sa.errors, fmt.Sprintf("method %s expects %d arguments, but got %d", methodName, len(methodDef.Formals), len(call.Arguments)))
		return methodDef.TypeDecl.Value // Return declared return type even with argument error for partial type info
	}
	for i, arg := range call.Arguments {
		argType := sa.getExpressionType(arg, st)
		formalType := methodDef.Formals[i].TypeDecl.Value

		if !sa.isTypeConformant(argType, formalType) {
			sa.errors = append(sa.errors, fmt.Sprintf("argument %d of method %s expects type %s, but got %s", i+1, methodName, formalType, argType))
		}
	}

	// 4. Return the method's return type (handle SELF_TYPE)
	returnType := methodDef.TypeDecl.Value
	if returnType == "SELF_TYPE" {
		return sa.currentClass // SELF_TYPE return resolves to the current class
	}
	return returnType
}

func (sa *SemanticAnalyser) getDotCallExpressionType(dotCall *ast.DotCallExpression, st *SymbolTable) string {
	// Dispatch:  expr.method(arg1, arg2, ...) or expr@Type.method(arg1, arg2, ...)

	objectType := sa.getExpressionType(dotCall.Object, st) // Type of the object on which method is called
	if objectType == "SELF_TYPE" {
		objectType = sa.currentClass // Resolve SELF_TYPE if needed
	}

	var dispatchType string  // Type to lookup method in (for static dispatch)
	if dotCall.Type != nil { // Static Dispatch: expr@Type.method(...)
		dispatchType = dotCall.Type.Value

		// Check if the static dispatch type exists
		_, typeFound := sa.globalSymbolTable.Lookup(dispatchType)
		if !typeFound {
			sa.errors = append(sa.errors, fmt.Sprintf("undefined type %s in static dispatch", dispatchType))
			return "Object"
		}

		if dispatchType == "SELF_TYPE" {
			dispatchType = sa.currentClass // Resolve static SELF_TYPE if needed (though semantically less common)
		}

		if !sa.isTypeConformant(objectType, dispatchType) { // Check if object's type conforms to static dispatch type
			sa.errors = append(sa.errors, fmt.Sprintf("static dispatch on type %s but receiver is of type %s", dispatchType, objectType))
			return "Object" // Or perhaps dispatchType for better error recovery
		}

	} else { // Dynamic Dispatch: expr.method(...)
		dispatchType = objectType // Dispatch on the object's type
	}

	// 1. Lookup the class symbol table based on dispatchType
	classEntry, classFound := sa.globalSymbolTable.Lookup(dispatchType)
	if !classFound {
		sa.errors = append(sa.errors, fmt.Sprintf("class %s not found for dispatch", dispatchType))
		return "Object"
	}
	classScope := classEntry.Scope
	methodName := dotCall.Method.Value

	// 2. Lookup the method in the class's symbol table
	methodEntry, methodFound := classScope.Lookup(methodName)
	if !methodFound || methodEntry.Type != "Method" {
		sa.errors = append(sa.errors, fmt.Sprintf("undefined method %s in %s", methodName, dispatchType))
		return "Object"
	}
	methodDef := methodEntry.Method

	// 3. Type check arguments (same as in getCallExpressionType)
	if len(dotCall.Arguments) != len(methodDef.Formals) {
		sa.errors = append(sa.errors, fmt.Sprintf("method %s expects %d arguments, but got %d", methodName, len(methodDef.Formals), len(dotCall.Arguments)))
		return methodDef.TypeDecl.Value // Partial type info
	}
	for i, arg := range dotCall.Arguments {
		argType := sa.getExpressionType(arg, st)
		formalType := methodDef.Formals[i].TypeDecl.Value
		if !sa.isTypeConformant(argType, formalType) {
			sa.errors = append(sa.errors, fmt.Sprintf("argument %d of method %s expects type %s, but got %s", i+1, methodName, formalType, argType))
		}
	}

	// 4. Return method return type (handle SELF_TYPE)
	returnType := methodDef.TypeDecl.Value
	if returnType == "SELF_TYPE" {
		// For SELF_TYPE return in dispatch, the result type is the *object's type*, not necessarily currentClass
		return objectType // Crucial adjustment for SELF_TYPE in dispatch
	}
	return returnType
}

func (sa *SemanticAnalyser) buildInheritanceGraph(program *ast.Program) {
	for _, class := range program.Classes {
		className := class.Name.Value
		if className == "Object" {
			continue // error would have already been reported before
		}

		parentName := "Object"
		if class.Parent != nil {
			parentName = class.Parent.Value
		}

		// Always look up parent - will be Object for classes without explicit parent
		parentEntry, ok := sa.globalSymbolTable.Lookup(parentName)
		if !ok {
			sa.errors = append(sa.errors,
				fmt.Sprintf("class %s inherits from undefined class %s",
					class.Name.Value, parentName))
			continue
		}

		if class.Parent != nil { // Only check forbidden inheritance for explicit parents
			if parentName == "Int" || parentName == "String" || parentName == "Bool" {
				sa.errors = append(sa.errors, fmt.Sprintf("class %s cannot inherit from built-in class %s", class.Name.Value, parentName))
				continue
			}
		}

		// Set parent scope for ALL classes (including those inheriting from Object)
		sa.globalSymbolTable.symbols[class.Name.Value].Scope.parent = parentEntry.Scope
		sa.inheritanceGraph.AddEdge(class.Name.Value, parentName)

		// Check for cycles
		if hasCycle, cycle := sa.inheritanceGraph.detectCycle(class.Name.Value); hasCycle {
			sa.errors = append(sa.errors, fmt.Sprintf("inheritance cycle detected: %v", cycle))
		}
	}
}

func (sa *SemanticAnalyser) validateMethodOverride(method *ast.Method, parentMethod *ast.Method) bool {
	// Check number of arguments
	if len(method.Formals) != len(parentMethod.Formals) {
		sa.errors = append(sa.errors, fmt.Sprintf(
			"method %s overrides parent method but has different number of parameters (%d vs %d)",
			method.Name.Value, len(method.Formals), len(parentMethod.Formals)))
		return false
	}

	// Check parameter types
	for i, formal := range method.Formals {
		if formal.TypeDecl.Value != parentMethod.Formals[i].TypeDecl.Value {
			sa.errors = append(sa.errors, fmt.Sprintf(
				"method %s overrides parent method but parameter %d has different type (%s vs %s)",
				method.Name.Value, i+1, formal.TypeDecl.Value, parentMethod.Formals[i].TypeDecl.Value))
			return false
		}
	}

	// Check return type
	expectedReturn := parentMethod.TypeDecl.Value
	actualReturn := method.TypeDecl.Value

	if expectedReturn != actualReturn {
		sa.errors = append(sa.errors, fmt.Sprintf(
			"method %s overrides parent method but has different return type (%s vs %s)",
			method.Name.Value, actualReturn, expectedReturn))
		return false
	}

	return true
}

func (sa *SemanticAnalyser) validateMainClass() {
	mainEntry, exists := sa.globalSymbolTable.Lookup("Main")
	if !exists {
		sa.errors = append(sa.errors, "program must have a class Main")
		return
	}

	// Check main method exists directly in Main class (not inherited)
	mainScope := mainEntry.Scope
	mainMethod, methodExists := mainScope.symbols["main"]
	if !methodExists || mainMethod.Type != "Method" {
		sa.errors = append(sa.errors, "class Main must define method 'main' with 0 parameters")
		return
	}

	// Verify no parameters
	if len(mainMethod.Method.Formals) > 0 {
		sa.errors = append(sa.errors, "main method must have 0 parameters")
	}
}

func (sa *SemanticAnalyser) isAttributeRedefined(attrName string, parentScope *SymbolTable) bool {
	if parentScope == nil {
		return false
	}
	if entry, ok := parentScope.Lookup(attrName); ok && entry.Type == "Attribute" {
		return true
	}
	return sa.isAttributeRedefined(attrName, parentScope.parent)
}

func (sa *SemanticAnalyser) isMethodRedefined(methodName string, parentScope *SymbolTable) (bool, *SymbolEntry) {
	if parentScope == nil {
		return false, nil
	}
	if entry, ok := parentScope.Lookup(methodName); ok && entry.Type == "Method" {
		return true, entry
	}
	return sa.isMethodRedefined(methodName, parentScope.parent)
}
