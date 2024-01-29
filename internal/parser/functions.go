package parser

import (
	"fmt"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	var parameters []*ast.Parameter

	if !p.curTokenIs(token.LPAREN) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		param := p.parseFunctionParameter()
		if param == nil {
			return nil
		}

		parameters = append(parameters, param)

		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RPAREN) {
	} else {
		return nil
	}

	return parameters
}

func (p *Parser) parseFunctionParameter() *ast.Parameter {
	paramType := p.curToken.Type
	if !p.isTypeToken(p.curToken) {
		p.peekError(token.IDENTIFIER)
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(token.IDENTIFIER) {
		fmt.Printf("Current Token: %s (Type: %s)\n", p.curToken.Literal, p.curToken.Type)
		fmt.Printf("Expected IDENTIFIER, got %s (Type: %s)\n", p.peekToken.Literal, p.peekToken.Type)
		return nil
	}
	paramName := p.curToken.Literal

	return &ast.Parameter{
		Identifier: &ast.Identifier{Value: paramName},
		Type:       paramType,
	}
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	var isExported bool
	var returnType ast.Expression

	if p.curToken.Type == token.PUB {
		isExported = true
		p.nextToken()
	}

	if p.isTypeToken(p.curToken) {
		returnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else {
		p.errors = append(p.errors, "expected return type before function name")
		return nil
	}

	ident := p.parseIdentifier()
	if ident == nil {
		return nil
	}
	p.nextToken()

	params := p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()

	body := p.parseBlockStatement()

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	stmt := &ast.FunctionStatement{
		IsExported: isExported,
		ReturnType: returnType,
		Name:       ident.(*ast.Identifier),
		Parameters: params,
		Body:       body,
	}

	return stmt
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	exp := &ast.FunctionCall{
		FunctionName: function.TokenLiteral(),
		Token:        p.curToken,
		Function:     function,
	}
	exp.Arguments = p.parseFunctionCallArguments()
	return exp
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	return args
}
