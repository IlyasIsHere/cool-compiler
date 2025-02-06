package parser

import (
	"cool-compiler/ast"
	"cool-compiler/lexer"
	"fmt"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

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
	if !p.expectCurrent(lexer.CLASS) {
		return nil
	}

	if !p.curTokenIs(lexer.TYPEID) {
		p.currentError(lexer.TYPEID)
		return nil
	}

	c.Name = &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// Doing: handle inheritance
	if p.peekTokenIs(lexer.INHERITS) {
		p.nextToken()

		if !p.peekTokenIs(lexer.TYPEID) {
			p.peekError(lexer.TYPEID)
			return nil
		}
		p.nextToken()

		// TODO:
		c.Parent = &ast.TypeIdentifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

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
	name := &ast.ObjectIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	attr.Name = name

	if !p.expectAndPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	if !p.curTokenIs(lexer.TYPEID) {
		p.currentError(lexer.TYPEID)
		return nil
	}
	typeid := &ast.TypeIdentifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	attr.TypeDecl = typeid

	p.nextToken()

	if p.curTokenIs(lexer.SEMI) {
		p.nextToken()
		return attr
	}

	if !p.expectCurrent(lexer.ASSIGN) {
		return nil
	}
	// TODO: attr.Expr = p.parseExpression()

}
