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
	ASSIGN
	EQUALS
	SUM
	PRODUCT
	ISVOID
	PREFIX
	AT
	CALL
)

var precedences = map[lexer.TokenType]Precedence{
	lexer.ASSIGN: ASSIGN,
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

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.ASSIGN, p.parseInfixExpression)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.TIMES, p.parseInfixExpression)
	p.registerInfix(lexer.DIVIDE, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.LE, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)

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

	for !p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		// TODO: parse formals
	}

	if !p.expectAndPeek(lexer.RPAREN) && !p.expectAndPeek(lexer.COLON) {
		return nil
	}

	if !p.peekTokenIs(lexer.TYPEID) {
		p.peekError(lexer.TYPEID)
		return nil
	}
	typeid := &ast.TypeIdentifier{
		Token: p.peekToken,
		Value: p.peekToken.Literal,
	}
	method.TypeDecl = typeid

	p.nextToken()
	p.nextToken()

	if !p.expectCurrent(lexer.LBRACE) {
		return nil
	}
	// TODO: method.body = p.parseExpression()

	if !p.expectAndPeek(lexer.RBRACE) {
		return nil
	}
	return method
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
		attr.Expr = p.parseExpression()
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
	return &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
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
	boolLiteral := &ast.BooleanLiteral{Token: p.curToken}
	if p.curToken.Literal == "true" {
		boolLiteral.Value = true
	} else {
		boolLiteral.Value = false
	}

	return boolLiteral
}

func (p *Parser) parseParenthesisExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectAndPeek(lexer.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
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
	exp := &ast.IsVoidExpression{Token: p.curToken}
	exp.Expression = p.parseExpression(ISVOID)
	return exp
}

func (p *Parser) parseWhileExpression() ast.Expression {
	exp := &ast.WhileExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if !p.expectAndPeek(lexer.LOOP) {
		return nil
	}

	exp.Body = p.parseExpression(LOWEST)
	if !p.expectAndPeek(lexer.POOL) {
		return nil
	}

	return exp
}

// func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
// 	exp := &ast.CallExpression{Token: p.curToken, Function: function}
// }

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

func (p *Parser) parseExpression(minPrecedence Precedence) ast.Expression {
	// TODO
	return nil
}
