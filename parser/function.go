package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *parser) parseFunction() (ast.Expression, error) {
	if !p.isFunctionDeclaration() {
		p.trace("expected a function definition")
		return nil, nil
	}

	return p.parseFunctionStatement()
}

func (p *parser) parseFunctionName() (*ast.Identifier, error) {
	if !p.curTokenIs(token.IDENTIFIER) {
		return nil, p.error("expected function name")
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return ident, nil
}

func (p *parser) parseFunctionStatement() (*ast.FunctionStatement, error) {
	var isExported bool
	var returnType *ast.Identifier
	defer p.trace("end of parsing function statement")

	if p.curToken.Type == token.PUB {
		isExported = true
		p.nextToken()
	}

	if p.isTypeToken(p.curToken) || p.isStructType(p.curToken) {
		returnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else if p.curToken.Type == token.FUNCTION {
		p.nextToken()
	} else {
		return nil, p.errorf("expected return type or 'fn', got %s instead", p.curToken.Type)
	}

	ident, err := p.parseFunctionName()
	if err != nil {
		return nil, err
	}
	if ident == nil {
		return nil, p.error("expected function name")
	}

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
		Name:       ident,
		Parameters: params,
		Body:       body,
	}

	return stmt, nil
}

func (p *parser) parseFunctionParameters() ([]*ast.Parameter, error) {
	parameters := []*ast.Parameter{}

	if !p.expectCurrentTokenIs(token.LPAREN) {
		return parameters, nil
	}

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
	}

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

func (p *parser) parseFunctionParameter() (*ast.Parameter, error) {
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

func (p *parser) parseFunctionCall(function ast.Expression) (ast.Expression, error) {
	var err error
	p.trace("parsing function call", function.TokenLiteral(), p.peekToken.Literal)
	if function == nil {
		println("function is nil before creating FunctionCall")
	}
	exp := &ast.FunctionCall{
		FunctionName: function.TokenLiteral(),
		Token:        p.curToken,
		Function:     function,
	}
	if p.curTokenIs(token.IDENTIFIER) {
		p.nextToken()
	}
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
	}
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return exp, nil
	}

	exp.Arguments, err = p.parseFunctionCallArguments()

	p.trace("parsed function call after args", p.curToken.Literal, p.peekToken.Literal)

	return exp, err
}

func (p *parser) parseFunctionCallArguments() ([]ast.Expression, error) {
	args := []ast.Expression{}
	p.trace("parsing func call args beginning", p.curToken.Literal, p.peekToken.Literal)
	defer p.trace("parseFunctionCallArguments-end:", p.curToken.Literal, p.peekToken.Literal)
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		p.trace("consume LPAREN", p.curToken.Literal, p.peekToken.Literal)
	}
	defer func() {
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}()

	p.trace("parseFunctionCallArguments before parse expression", p.curToken.Literal, p.peekToken.Literal)
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return args, nil
	}
	firstArg, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if firstArg == nil {
		return nil, p.error("could not parse first argument in function call")
	}
	p.trace("parseFunctionCallArguments: first arg", firstArg.String(), p.peekToken.Literal)
	args = append(args, firstArg)

	if !p.curTokenIs(token.COMMA) && !p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	for p.curTokenIs(token.COMMA) {
		p.trace("parseFunctionCallArguments: consume COMMA", p.curToken.Literal, p.peekToken.Literal)
		if p.curTokenIs(token.RPAREN) {
			return args, nil
		}
		if p.curTokenIs(token.COMMA) && !p.peekTokenIs(token.RPAREN) {
			p.nextToken()
		}

		nextArg, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		if nextArg == nil {
			return nil, p.error("could not parse argument after comma in function call")
		}
		p.trace("parseFunctionCallArguments: next arg", nextArg.String(), p.curToken.Literal, p.peekToken.Literal)
		args = append(args, nextArg)
		// p.nextToken()
	}

	return args, nil
}

func (p *parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	p.trace("parsing return statement", p.curToken.Literal)
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // consume 'return' keyword

	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.RBRACE) {
		p.nextToken() // consume semicolon if present
		return stmt, nil
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if expr != nil {
		p.trace("return expression parsed:", expr.String())
		stmt.ReturnValues = append(stmt.ReturnValues, expr)
	} else {
		return nil, p.error("expected expression after 'return'")
	}

	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}
