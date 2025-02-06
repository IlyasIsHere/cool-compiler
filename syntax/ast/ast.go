package ast

import (
	"cool-compiler/lexer"
)

type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type Feature interface {
	Node
	featureNode()
}

type TypeIdentifier struct {
	Token lexer.Token
	Value string
}

func (ti *TypeIdentifier) TokenLiteral() string { return ti.Token.Literal }

type ObjectIdentifier struct {
	Token lexer.Token
	Value string
}

func (oi *ObjectIdentifier) TokenLiteral() string { return oi.Token.Literal }
func (oi *ObjectIdentifier) expressionNode()      {}

type Program struct {
	Classes []*Class
}

func (p *Program) TokenLiteral() string { return "" }

type Class struct {
	Token    lexer.Token
	Name     *TypeIdentifier
	Features []Feature
	Parent   *TypeIdentifier
}

func (c *Class) TokenLiteral() string { return c.Token.Literal }

type Attribute struct {
	Token      lexer.Token
	Name       *ObjectIdentifier
	TypeDecl   *TypeIdentifier
	Expression Expression
}

func (a *Attribute) TokenLiteral() string { return a.Token.Literal }
func (a *Attribute) featureNode()         {}

type Method struct {
	Token      lexer.Token
	Name       *ObjectIdentifier
	TypeDecl   *TypeIdentifier
	Formals    []Formal
	Expression Expression
}

func (m *Method) TokenLiteral() string { return m.Token.Literal }
func (m *Method) featureNode()         {}

type Formal struct {
	Token    lexer.Token
	Name     *ObjectIdentifier
	TypeDecl *TypeIdentifier
}

func (f *Formal) TokenLiteral() string { return f.Token.Literal }

type IntegerLiteral struct {
	Token lexer.Token
	Value int
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) expressionNode()      {}

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) expressionNode()      {}

type IfExpression struct {
	Token       lexer.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ifexp *IfExpression) TokenLiteral() string { return ifexp.Token.Literal }
func (ifexp *IfExpression) expressionNode()      {}

type WhileExpression struct {
	Token     lexer.Token
	Condition Expression
	Body      Expression
}

func (wexp *WhileExpression) TokenLiteral() string { return wexp.Token.Literal }
func (wexp *WhileExpression) expressionNode()      {}

type BlockExpression struct {
	Token       lexer.Token
	Expressions []Expression
}

func (bexp *BlockExpression) TokenLiteral() string { return bexp.Token.Literal }
func (bexp *BlockExpression) expressionNode()      {}

type Binding struct {
	Token      lexer.Token
	Identifier *ObjectIdentifier
	Type       *TypeIdentifier
	Init       Expression
}

func (b *Binding) TokenLiteral() string { return b.Token.Literal }

type LetExpression struct {
	Token    lexer.Token
	Bindings []*Binding
	In       Expression
}

func (lexp *LetExpression) TokenLiteral() string { return lexp.Token.Literal }
func (lexp *LetExpression) expressionNode()      {}

type NewExpression struct {
	Token lexer.Token
	Type  *TypeIdentifier
}

func (nexp *NewExpression) TokenLiteral() string { return nexp.Token.Literal }
func (nexp *NewExpression) expressionNode()      {}

type IsVoidExpression struct {
	Token      lexer.Token
	Expression Expression
}

func (ivexpr *IsVoidExpression) TokenLiteral() string { return ivexpr.Token.Literal }
func (ivexpr *IsVoidExpression) expressionNode()      {}

type UnaryExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (unaryexpr *UnaryExpression) TokenLiteral() string { return unaryexpr.Token.Literal }
func (unaryexpr *UnaryExpression) expressionNode()      {}

type BinaryExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (binexp *BinaryExpression) TokenLiteral() string { return binexp.Token.Literal }
func (binexp *BinaryExpression) expressionNode()      {}

type CaseExpression struct {
	Token      lexer.Token
	Expression Expression
	Branches   []CaseBranch
}

func (caseexp *CaseExpression) TokenLiteral() string { return caseexp.Token.Literal }
func (caseexp *CaseExpression) expressionNode()      {}

type CaseBranch struct {
	Token      lexer.Token
	Identifier *ObjectIdentifier
	Type       *TypeIdentifier
	Expression Expression
}

type Assignment struct {
	Token      lexer.Token
	Identifier *ObjectIdentifier
	Expression Expression
}

func (a *Assignment) TokenLiteral() string { return a.Token.Literal }
func (a *Assignment) expressionNode()      {}
