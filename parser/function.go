package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *Parser) parseFunctionStatement() (*ast.FunctionStatement, error) {
	var isExported bool
	var returnType *ast.Identifier

	if p.curToken.Type == token.PUB {
		isExported = true
		p.nextToken()
	}

	if p.isTypeToken(p.curToken) {
		returnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else if p.curToken.Type == token.FUNCTION {
		p.nextToken()
	} else {
		return nil, p.errorf("expected return type or 'fn', got %s instead", p.curToken.Type)
	}

	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if ident == nil {
		return nil, p.error("expected function name")
	}
	p.nextToken()

	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	if params == nil {
		return nil, p.error("failed to parse function parameters")
	}

	if !p.expectCurrentTokenIs(token.LBRACE) {
		return nil, p.error("expected '{' to start function body")
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	stmt := &ast.FunctionStatement{
		IsExported: isExported,
		ReturnType: returnType,
		Name:       ident.(*ast.Identifier),
		Parameters: params,
		Body:       body,
	}

	return stmt, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.Parameter, error) {
	parameters := []*ast.Parameter{}

	if !p.expectCurrentTokenIs(token.LPAREN) {
		return parameters, nil
	}

	p.nextToken()

	for !p.curTokenIs(token.RPAREN) {
		param, err := p.parseFunctionParameter()
		if err != nil {
			return nil, err
		}
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

	return parameters, nil
}

func (p *Parser) parseFunctionParameter() (*ast.Parameter, error) {
	if !p.isTypeToken(p.curToken) {
		return nil, p.errorf("expected type token, got %s instead", p.curToken.Type)
	}
	paramType := p.curToken

	p.nextToken()

	if !p.curTokenIs(token.IDENTIFIER) {
		return nil, p.errorf("expected identifier token, got %s instead", p.curToken.Type)
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
	}, nil
}

func (p *Parser) parseFunctionCall(function ast.Expression) (ast.Expression, error) {
	var err error
	p.trace("parsing function call", function.TokenLiteral(), p.peekToken.Literal)
	exp := &ast.FunctionCall{
		FunctionName: function.TokenLiteral(),
		Token:        p.curToken,
		Function:     function,
	}
	p.nextToken() // consume (
	exp.Arguments, err = p.parseFunctionCallArguments()

	p.trace("parsed function call after args", p.curToken.Literal, p.peekToken.Literal)
	return exp, err
}

func (p *Parser) parseFunctionCallArguments() ([]ast.Expression, error) {
	args := []ast.Expression{}
	p.trace("parsing func call args beginning", p.curToken.Literal, p.peekToken.Literal)
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		p.trace("consume LPAREN", p.curToken.Literal, p.peekToken.Literal)
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken() // consume the closing parenthesis
		return args, nil
	}

	firstArg, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if firstArg == nil {
		return nil, p.error("could not parse first argument in function call")
	}
	p.trace("parseFunctionCallArguments: first arg", firstArg.String(), p.curToken.Literal, p.peekToken.Literal)
	args = append(args, firstArg)

	for p.curTokenIs(token.COMMA) {
		p.trace("parseFunctionCallArguments: consume COMMA", p.curToken.Literal, p.peekToken.Literal)
		// consume the comma
		p.nextToken()

		nextArg, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		if nextArg == nil {
			return nil, p.error("could not parse argument after comma in function call")
		}
		p.trace("parseFunctionCallArguments: next arg", nextArg.String(), p.curToken.Literal, p.peekToken.Literal)
		args = append(args, nextArg)
	}

	p.trace("parseFunctionCallArguments: end", p.curToken.Literal, p.peekToken.Literal)
	return args, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	p.trace("parsing return statement", p.curToken.Literal)
	stmt := &ast.ReturnStatement{Token: p.curToken}
	if p.curTokenIs(token.RETURN) {
		p.nextToken()
	}
	p.trace("after return keyword", p.curToken.Literal, string(p.curToken.Type))

	if p.peekTokenAfter(token.COMMA) || p.peekTokenIs(token.COMMA) {
		if p.peekTokenAfter(token.COMMA) && p.curTokenIs(token.LPAREN) {
			p.nextToken()
		}

		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			expr, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.error("expected expression in return statement")
			}

			stmt.ReturnValues = append(stmt.ReturnValues, expr)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}

	} else {
		p.trace("parse return value", p.curToken.Literal, p.peekToken.Literal)
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		p.trace("after parse return value", expr.String())
		p.trace("return expression", expr.String())
		if expr != nil {
			stmt.ReturnValues = append(stmt.ReturnValues, expr)
		} else {
			return nil, p.error("expected expression after 'return'")
		}
	}
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return stmt, nil
}
