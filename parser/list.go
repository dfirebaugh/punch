package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
	"github.com/sirupsen/logrus"
)

func (p *parser) parseListLiteral() (ast.Expression, error) {
	if p.curTokenIs(token.LBRACKET) && p.peekTokenIs(token.RBRACKET) && p.peekTokenAfter(token.IDENTIFIER) {
		p.nextToken()
		p.nextToken()
		p.nextToken()
	}
	list := &ast.ListLiteral{Token: p.curToken}
	list.Elements = p.parseExpressionList()
	return list, nil
}

func (p *parser) parseListOperation() (ast.Expression, error) {
	var err error
	op := &ast.ListOperation{Token: p.curToken}
	op.Operator = p.curToken.Literal

	p.nextToken()

	op.List, err = p.parseExpression(LOWEST)
	if err != nil {
		logrus.Error(err)
	}

	if p.curTokenIs(token.COMMA) {
		p.nextToken()
		op.Element, err = p.parseExpression(LOWEST)
		if err != nil {
			logrus.Error(err)
		}
	}

	return op, nil
}

func (p *parser) parseExpressionList() []ast.Expression {
	p.trace("parsing expression list")
	list := []ast.Expression{}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken() // consume '{'
	}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		logrus.Error(err)
	}
	list = append(list, expression)

	for p.curTokenIs(token.COMMA) {
		p.trace("parse expression in list")
		if p.curTokenIs(token.COMMA) {
			p.nextToken() // consume comma
		}
		expression, err = p.parseExpression(LOWEST)
		if err != nil {
			logrus.Error(err)
		}
		list = append(list, expression)
	}

	return list
}

func (p *parser) parseListDeclaration() (*ast.ListDeclaration, error) {
	decl := &ast.ListDeclaration{Token: p.curToken}
	if p.curTokenIs(token.LBRACKET) {
		p.nextToken()
	}
	if p.curTokenIs(token.RBRACKET) {
		p.nextToken()
	}
	// Expect the type token (e.g., "let" or "const")
	if !p.expectPeek(token.IDENTIFIER) {
		return nil, p.error("expected type token")
	}
	decl.Type = p.curToken.Type
	p.nextToken()

	// Expect the identifier token (name of the list)
	if !p.expectCurrentTokenIs(token.IDENTIFIER) {
		return nil, p.error("expected identifier")
	}
	decl.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect the assignment token
	if !p.expectPeek(token.ASSIGN) {
		return nil, p.error("expected '=' after identifier")
	}

	// Parse the list literal
	p.nextToken()
	p.nextToken() // consume assign operator

	listLiteral, err := p.parseListLiteral()
	if err != nil {
		return nil, err
	}
	decl.Value = listLiteral.(*ast.ListLiteral)

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}
	return decl, nil
}
