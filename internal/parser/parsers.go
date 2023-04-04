package parser

import (
	"fmt"
	"strconv"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.SLASH_SLASH:
		p.parseComment()
		return nil
	case token.SLASH_ASTERISK:
		p.parseComment()
		return nil
	case token.PUB:
		return p.parseFunctionStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if stmt.ReturnValue == nil {
		p.errors = append(p.errors, "expected expression after 'return'")
		return nil
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	var isExported bool
	if p.curToken.Type != token.FUNCTION && p.curToken.Type != token.PUB {
		return nil
	}

	if p.curToken.Type == token.PUB {
		isExported = true
		p.nextToken()
	}

	p.nextToken()
	ident := p.parseIdentifier()
	if ident == nil {
		return nil
	}

	if p.peekToken.Type != token.LPAREN {
		return nil
	}

	p.nextToken()

	params := p.parseFunctionParameters()

	if p.curToken.Type != token.RPAREN {
		return nil
	}
	p.nextToken()

	body := p.parseBlockStatement()

	stmt := &ast.FunctionStatement{
		IsExported: isExported,
		Name:       ident.(*ast.Identifier),
		Parameters: params,
		Body:       body,
	}

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	p.nextToken()

	return identifiers
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	right := p.parseExpression(PREFIX)

	expression.Right = right

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	expression.Right = right

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}

// parseComment skips over a comment token and moves the parser to the end of the comment.
func (p *Parser) parseComment() {
	if p.curToken.Type != token.SLASH_SLASH && p.curToken.Type != token.SLASH_ASTERISK {
		return
	}

	if p.curToken.Type == token.SLASH_SLASH {
		currentLine := p.curToken.Position.Line
		for p.curToken.Position.Line == currentLine {
			p.nextToken()
		}
	} else if p.curToken.Type == token.SLASH_ASTERISK {
		// Advance past the opening token
		p.nextToken()

		// Loop until we find the closing token
		for !(p.curTokenIs(token.ASTERISK) && p.peekTokenIs(token.SLASH)) {
			println(p.curToken.Type)
			if p.curTokenIs(token.EOF) {
				return
			}
			println(p.curToken.Type)
			p.nextToken()
		}

		println(p.curToken.Type)
		// Advance past the closing token
		p.nextToken()
		println(p.curToken.Type)
	}
}
