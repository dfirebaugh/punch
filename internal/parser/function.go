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
				p.error(fmt.Sprintf("expected return type, got %s instead", p.curToken.Type))
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
			p.error("expected ')' after return types")
			return nil
		}
		p.nextToken()
	} else if p.isTypeToken(p.curToken) {
		returnType := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		returnTypes = append(returnTypes, returnType)
		p.nextToken()
	} else {
		p.error("expected return type or '(' for multiple return types")
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
		p.error(fmt.Sprintf("expected type token, got %s instead", p.curToken.Type))
		return nil
	}
	paramType := p.curToken

	p.nextToken()

	if !p.curTokenIs(token.IDENTIFIER) {
		p.error(fmt.Sprintf("expected identifier token, got %s instead", p.curToken.Type))
		return nil
	}
	paramName := p.curToken.Literal

	return &ast.Parameter{
		Identifier: &ast.Identifier{
			Token: token.Token{
				Type:     token.IDENTIFIER,
				Literal:  paramName,
				Position: p.curToken.Position,
			},
			Value: paramName,
		},
		Type: paramType.Type,
	}
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	p.trace("parsing function call", function.TokenLiteral(), p.peekToken.Literal)
	exp := &ast.FunctionCall{
		FunctionName: function.TokenLiteral(),
		Token:        p.curToken,
		Function:     function,
	}
	p.nextToken()
	exp.Arguments = p.parseFunctionCallArguments()

	if len(exp.Arguments) == 1 && p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}
	if !p.expectCurrentTokenIs(token.RPAREN) {
		p.error("expected ')' after function call arguments")
		return nil
	}
	p.trace("parsed function call after args", p.curToken.Literal, p.peekToken.Literal)
	return exp
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	args := []ast.Expression{}
	p.trace("parsing func call args beginning", p.curToken.Literal, p.peekToken.Literal)
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		p.trace("consume LPAREN", p.curToken.Literal, p.peekToken.Literal)
	}

	args = append(args, p.parseExpression(LOWEST))
	p.trace("after parsing first arg", args[0].String(), p.curToken.Literal, p.peekToken.Literal)

	if p.curTokenIs(token.RPAREN) {
		return args
	}

	for p.curTokenIs(token.COMMA) && !p.curTokenIs(token.RPAREN) {
		p.nextToken()
		p.trace("func args 1", p.curToken.Literal, p.peekToken.Literal)
		args = append(args, p.parseExpression(LOWEST))
		p.trace("func args 2", p.curToken.Literal, p.peekToken.Literal)
	}
	p.trace("after parsing func args", p.curToken.Literal, p.peekToken.Literal)
	return args
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.trace("parsing return statement", p.curToken.Literal)
	stmt := &ast.ReturnStatement{Token: p.curToken}
	if p.curTokenIs(token.RETURN) {
		p.nextToken() // consume return keyword
	}
	p.trace("after return keyword", p.curToken.Literal, string(p.curToken.Type))

	if p.peekTokenAfter(token.COMMA) || p.peekTokenIs(token.COMMA) {
		if p.peekTokenAfter(token.COMMA) && p.curTokenIs(token.LPAREN) {
			p.nextToken()
		}

		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			expr := p.parseExpression(LOWEST)
			if expr == nil {
				p.error("expected expression in return statement")
				return nil
			}

			stmt.ReturnValues = append(stmt.ReturnValues, expr)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}

	} else {
		p.trace("parse return value", p.curToken.Literal, p.peekToken.Literal)
		expr := p.parseExpression(LOWEST)
		p.trace("after parse return value", expr.String())
		p.trace("return expression", expr.String())
		if expr != nil {
			stmt.ReturnValues = append(stmt.ReturnValues, expr)
		} else {
			p.error("expected expression after 'return'")
			return nil
		}
	}
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return stmt
}
