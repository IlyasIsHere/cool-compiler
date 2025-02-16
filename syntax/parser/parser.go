package parser

import (
	"cool-compiler/ast"
	"cool-compiler/lexer"
	"fmt"
	"strconv"
)

type Precedence int

const (
	LOWEST Precedence = iota
	EQUALS
	SUM
	PRODUCT
	ISVOID
	PREFIX
	AT
	CALL
)

var precedences = map[lexer.TokenType]Precedence{
	lexer.ASSIGN: LOWEST,
	lexer.PLUS:   SUM,
	lexer.MINUS:  SUM,
	lexer.TIMES:  PRODUCT,
	lexer.DIVIDE: PRODUCT,
	lexer.LT:     EQUALS,
	lexer.LE:     EQUALS,
	lexer.EQ:     EQUALS,
	lexer.LPAREN: CALL,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.OBJECTID, p.parseIdentifier)
	p.registerPrefix(lexer.INT_CONST, p.parseInteger)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.BOOL_CONST, p.parseBoolean)
	p.registerPrefix(lexer.LPAREN, p.parseParenthesisExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.WHILE, p.parseWhileExpression)
	p.registerPrefix(lexer.LET, p.parseLetExpression)
	p.registerPrefix(lexer.CASE, p.parseCaseExpression)
	p.registerPrefix(lexer.LBRACE, p.parseBlockExpression)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.ASSIGN, p.parseAssignment)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.TIMES, p.parseInfixExpression)
	p.registerInfix(lexer.DIVIDE, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.LE, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	// p.registerInfix(lexer.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectAndPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) expectCurrent(t lexer.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.currentError(t)
	return false
}

func (p *Parser) peekError(t lexer.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected next token to be %v, got %v line %d col %d", t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column))
}

func (p *Parser) currentError(t lexer.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("Expected current token to be %v, got %v line %d col %d", t, p.curToken.Type, p.curToken.Line, p.curToken.Column))
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}
	for p.curToken.Type != lexer.EOF && p.curToken.Type != lexer.ERROR {
		c := p.ParseClass()

		if !p.expectAndPeek(lexer.SEMI) {
			continue
		}
		p.nextToken()
		prog.Classes = append(prog.Classes, c)
	}
	return prog
}

func (p *Parser) ParseClass() *ast.Class {

	c := &ast.Class{Token: p.curToken}
	if !p.curTokenIs(lexer.CLASS) {
		p.currentError(lexer.CLASS)
		return nil
	}
	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}

	c.Name = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Doing: handle inheritance
	if p.peekTokenIs(lexer.INHERITS) {
		p.nextToken()

		if !p.expectAndPeek(lexer.TYPEID) {
			return nil
		}
		c.Parent = &ast.TypeIdentifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}

	if !p.expectAndPeek(lexer.LBRACE) {
		return nil
	}
	for !p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		c.Features = append(c.Features, p.parseFeature())
		if !p.expectAndPeek(lexer.SEMI) {
			return nil
		}
	}

	if !p.expectAndPeek(lexer.RBRACE) {
		return nil
	}

	return c
}

func (p *Parser) parseFeature() ast.Feature {
	if p.peekTokenIs(lexer.LPAREN) {
		return p.parseMethod()
	}
	return p.parseAttribute()
}

func (p *Parser) parseMethod() *ast.Method {
	method := &ast.Method{Token: p.curToken}

	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}
	name := &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	method.Name = name

	if !p.expectAndPeek(lexer.LPAREN) {
		return nil
	}

	method.Formals = p.parseFormals()
	if method.Formals == nil {
		return nil
	}

	if !p.expectAndPeek(lexer.RPAREN) && !p.expectAndPeek(lexer.COLON) {
		return nil
	}

	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}
	method.TypeDecl = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.LBRACE) {
		return nil
	}
	p.nextToken()
	method.Expression = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.RBRACE) {
		return nil
	}
	return method
}

func (p *Parser) parseFormals() []*ast.Formal {
	var formals []*ast.Formal

	if p.peekTokenIs(lexer.RPAREN) {
		return formals
	}

	p.nextToken()
	formal := p.parseFormal()
	if formal == nil {
		return nil
	}

	formals = append(formals, formal)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move the next formal
		formal := p.parseFormal()
		if formal == nil {
			return nil
		}
		formals = append(formals, formal)
	}

	return formals
}

func (p *Parser) parseFormal() *ast.Formal {
	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}

	formal := &ast.Formal{Token: p.curToken}
	formal.Name = &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.COLON) {
		return nil
	}

	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}

	formal.TypeDecl = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return formal
}

func (p *Parser) parseAttribute() *ast.Attribute {
	attr := &ast.Attribute{
		Token: p.curToken,
	}

	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}
	attr.Name = &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.COLON) {
		return nil
	}

	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}
	attr.TypeDecl = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()
		p.nextToken()
		attr.Expression = p.parseExpression(LOWEST)
	}

	return attr
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tt lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}
	return &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	if !p.curTokenIs(lexer.BOOL_CONST) {
		p.currentError(lexer.BOOL_CONST)
		return nil
	}
	intLiteral := &ast.IntegerLiteral{Token: p.curToken}

	num, err := strconv.Atoi(p.curToken.Literal)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	intLiteral.Value = num

	return intLiteral
}

func (p *Parser) parseBoolean() ast.Expression {
	if !p.curTokenIs(lexer.BOOL_CONST) {
		p.currentError(lexer.BOOL_CONST)
		return nil
	}
	boolLiteral := &ast.BooleanLiteral{Token: p.curToken}
	if p.curToken.Literal == "true" {
		boolLiteral.Value = true
	} else {
		boolLiteral.Value = false
	}

	return boolLiteral
}

func (p *Parser) parseParenthesisExpression() ast.Expression {
	if !p.curTokenIs(lexer.LPAREN) {
		p.currentError(lexer.LPAREN)
		return nil
	}
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	if !p.curTokenIs(lexer.IF) {
		p.currentError(lexer.IF)
		return nil
	}
	exp := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.THEN) {
		return nil
	}
	p.nextToken()
	exp.Consequence = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.ELSE) {
		return nil
	}
	p.nextToken()
	exp.Alternative = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.FI) {
		return nil
	}
	return exp
}

func (p *Parser) parseNewExpression() ast.Expression {
	if !p.curTokenIs(lexer.NEW) {
		p.currentError(lexer.NEW)
		return nil
	}
	exp := &ast.NewExpression{Token: p.curToken}
	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}

	exp.Type = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	return exp
}

func (p *Parser) parseIsVoidExpression() ast.Expression {
	if !p.curTokenIs(lexer.ISVOID) {
		p.currentError(lexer.ISVOID)
		return nil
	}
	exp := &ast.IsVoidExpression{Token: p.curToken}
	p.nextToken()
	exp.Expression = p.parseExpression(ISVOID)
	return exp
}

func (p *Parser) parseWhileExpression() ast.Expression {
	if !p.curTokenIs(lexer.WHILE) {
		p.currentError(lexer.WHILE)
		return nil
	}
	exp := &ast.WhileExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if !p.expectAndPeek(lexer.LOOP) {
		return nil
	}

	p.nextToken()
	exp.Body = p.parseExpression(LOWEST)
	if !p.expectAndPeek(lexer.POOL) {
		return nil
	}

	return exp
}

func (p *Parser) parseLetExpression() ast.Expression {
	if !p.curTokenIs(lexer.LET) {
		p.currentError(lexer.LET)
		return nil
	}

	exp := &ast.LetExpression{Token: p.curToken}
	var bindings []*ast.Binding

	p.nextToken()
	binding := p.parseBinding()
	if binding == nil {
		return nil
	}
	bindings = append(bindings, binding)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		binding := p.parseBinding()
		if binding == nil {
			return nil
		}
		bindings = append(bindings, binding)
	}

	exp.Bindings = bindings

	if !p.expectAndPeek(lexer.IN) {
		return nil
	}
	p.nextToken()
	exp.In = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parseBinding() *ast.Binding {
	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}

	binding := &ast.Binding{Token: p.curToken}
	binding.Identifier = &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.COLON) {
		return nil
	}
	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}

	binding.Type = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()
		p.nextToken()
		binding.Init = p.parseExpression(LOWEST)
	}

	return binding
}

func (p *Parser) parseCaseExpression() ast.Expression {
	if !p.curTokenIs(lexer.CASE) {
		p.currentError(lexer.CASE)
		return nil
	}

	exp := &ast.CaseExpression{Token: p.curToken}
	p.nextToken()
	exp.Expression = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.OF) {
		return nil
	}

	var branches []*ast.CaseBranch
	p.nextToken()
	// parsing the first branch
	branch := p.parseCaseBranch()
	if branch == nil {
		return nil
	}
	branches = append(branches, branch)

	// parsing the remaining branches if any
	for !p.peekTokenIs(lexer.ESAC) {
		p.nextToken()
		branch := p.parseCaseBranch()
		if branch == nil {
			return nil
		}
		branches = append(branches, branch)
	}

	exp.Branches = branches

	if !p.expectAndPeek(lexer.ESAC) {
		return nil
	}

	return exp
}

func (p *Parser) parseCaseBranch() *ast.CaseBranch {
	if !p.curTokenIs(lexer.OBJECTID) {
		p.currentError(lexer.OBJECTID)
		return nil
	}

	branch := &ast.CaseBranch{Token: p.curToken}
	branch.Identifier = &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.COLON) {
		return nil
	}
	if !p.expectAndPeek(lexer.TYPEID) {
		return nil
	}

	branch.Type = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectAndPeek(lexer.DARROW) {
		return nil
	}
	p.nextToken()

	branch.Expression = p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.SEMI) {
		return nil
	}

	return branch
}

func (p *Parser) parseBlockExpression() ast.Expression {
	if !p.curTokenIs(lexer.LBRACE) {
		p.currentError(lexer.LBRACE)
		return nil
	}

	blockExp := &ast.BlockExpression{Token: p.curToken}
	var expressions []ast.Expression

	p.nextToken()

	// parse the first expression
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	expressions = append(expressions, expr)
	if !p.expectAndPeek(lexer.SEMI) {
		return nil
	}

	// parse the remaining ones
	for !p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}
		expressions = append(expressions, expr)

		if !p.expectAndPeek(lexer.SEMI) {
			return nil
		}
	}

	blockExp.Expressions = expressions
	return blockExp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// TODO:
	return nil
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignment(left ast.Expression) ast.Expression {
	identifier, ok := left.(*ast.ObjectIdentifier)
	if !ok {
		msg := fmt.Sprintf("expected identifier on left side of assignment, got %T", left)
		p.errors = append(p.errors, msg)
		return nil
	}

	assignment := &ast.AssignmentExpression{
		Token:      p.curToken,
		Identifier: identifier,
	}

	p.nextToken()
	assignment.Expression = p.parseExpression(LOWEST)

	return assignment
}

func (p *Parser) curPrecedence() Precedence {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(minPrecedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMI) && minPrecedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
