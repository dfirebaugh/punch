package parser

import (
	"fmt"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	var isExported bool
	returnTypes := []ast.Expression{}

	if p.curToken.Type == token.PUB {
		isExported = true
		p.nextToken()
	}

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()

		for !p.curTokenIs(token.RPAREN) {
			if !p.isTypeToken(p.curToken) {
				p.errors = append(p.errors, fmt.Sprintf("expected return type, got %s instead", p.curToken.Type))
				return nil
			}

			returnType := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			returnTypes = append(returnTypes, returnType)
			p.nextToken()

			if p.curTokenIs(token.COMMA) {
				p.nextToken()
			}
		}

		if !p.expectCurrentTokenIs(token.RPAREN) {
			p.errors = append(p.errors, "expected ')' after return types")
			return nil
		}
		p.nextToken()
	} else if p.isTypeToken(p.curToken) {
		returnType := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		returnTypes = append(returnTypes, returnType)
		p.nextToken()
	} else {
		p.errors = append(p.errors, "expected return type or '(' for multiple return types")
		return nil
	}

	ident := p.parseIdentifier()
	if ident == nil {
		return nil
	}
	p.nextToken()

	params := p.parseFunctionParameters()
	if !p.expectCurrentTokenIs(token.LBRACE) {
		return nil
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}
	body := p.parseBlockStatement()

	stmt := &ast.FunctionStatement{
		IsExported:  isExported,
		ReturnTypes: returnTypes,
		Name:        ident.(*ast.Identifier),
		Parameters:  params,
		Body:        body,
	}

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	if !p.expectCurrentTokenIs(token.LPAREN) {
		return parameters
	}

	p.nextToken()

	for !p.curTokenIs(token.RPAREN) {
		param := p.parseFunctionParameter()
		if param != nil {
			parameters = append(parameters, param)
		}

		p.nextToken()

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	return parameters
}

func (p *Parser) parseFunctionParameter() *ast.Parameter {
	if !p.isTypeToken(p.curToken) {
		p.errors = append(p.errors, fmt.Sprintf("expected type token, got %s instead", p.curToken.Type))
		return nil
	}
	paramType := p.curToken

	p.nextToken()

	if !p.curTokenIs(token.IDENTIFIER) {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier token, got %s instead", p.curToken.Type))
		return nil
	}
	paramName := p.curToken.Literal

	return &ast.Parameter{
		Identifier: &ast.Identifier{Token: paramType, Value: paramName},
		Type:       paramType.Type,
	}
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	exp := &ast.FunctionCall{
		FunctionName: function.TokenLiteral(),
		Token:        p.curToken,
		Function:     function,
	}
	p.nextToken()
	exp.Arguments = p.parseFunctionCallArguments()
	p.nextToken()
	return exp
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
	}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	return args
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	if p.peekTokenAfter(token.COMMA) || p.peekTokenIs(token.COMMA) {
		if p.peekTokenAfter(token.COMMA) && p.curTokenIs(token.LPAREN) {
			p.nextToken()
		}

		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			expr := p.parseExpression(LOWEST)
			if expr == nil {
				p.errors = append(p.errors, "expected expression in return statement")
				return nil
			}

			stmt.ReturnValues = append(stmt.ReturnValues, expr)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}

		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	} else {
		expr := p.parseExpression(LOWEST)
		if expr != nil {
			stmt.ReturnValues = append(stmt.ReturnValues, expr)
		} else {
			p.errors = append(p.errors, "expected expression after 'return'")
			return nil
		}
	}

	return stmt
}
