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
		program.Statements = append(program.Statements, p.parseStatement())
		p.nextToken()
		if p.curTokenIs(token.RBRACE) {
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.LPAREN) {
		return p.parseFunctionStatement()
	}
	if p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.ASSIGN) {
		return p.parseTypeBasedVariableDeclaration()
	}

	switch p.curToken.Type {
	case token.STRUCT:
		return p.parseStructDefinition()
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
	case token.DEFER:
		return p.parseDeferStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		if p.curTokenIs(token.RBRACE) && p.peekTokenIs(token.EOF) {
			return nil
		}
		return p.parseExpressionStatement()
	}
}

func (p *Parser) isTypeToken(t token.Token) bool {
	switch t.Type {
	case token.STRING,
		token.BOOL,
		token.U8,
		token.U16,
		token.U32,
		token.U64,
		token.I8,
		token.I16,
		token.I32,
		token.I64,
		token.F8,
		token.F16,
		token.F32,
		token.F64:
		return true
	}

	_, exists := p.definedTypes[t.Literal]
	return exists
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// Expect and parse the identifier (variable name)
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect and consume the '=' token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Move past '=' and parse the expression/value
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseNumberType(typeToken token.Type) ast.Expression {
	return &ast.NumberType{
		Token: p.curToken,
		Type:  typeToken,
	}
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

func (p *Parser) parseStringLiteral() ast.Expression {
	if !p.curTokenIs(token.STRING) {
		p.errors = append(p.errors, "expected string literal")
		return nil
	}

	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken,
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
		Right:    nil,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		fmt.Println("Failed to find closing parenthesis")
		return nil
	}
	p.nextToken()

	return exp
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	expression := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
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
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		switch stmt.(type) {
		case *ast.ReturnStatement:
			if p.curTokenIs(token.RPAREN) {
				p.nextToken()
			}
			return block
		}
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
			if p.curTokenIs(token.EOF) {
				return
			}
			p.nextToken()
		}

		// Advance past the closing token
		p.nextToken()
	}
}

func (p *Parser) parseDeferStatement() *ast.DeferStatement {
	deferStmt := &ast.DeferStatement{Token: p.curToken}

	p.nextToken()
	deferStmt.Statement = p.parseStatement()

	p.deferStack = append(p.deferStack, deferStmt.Statement)

	return deferStmt
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	if _, ok := left.(*ast.Identifier); !ok {
		p.errors = append(p.errors, "left-hand side of assignment must be an identifier")
		return nil
	}

	expression := &ast.AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Right = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseTypeBasedVariableDeclaration() ast.Statement {
	varType := p.curToken
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	p.nextToken()
	varName := p.curToken

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	varDecl := &ast.VariableDeclaration{
		Type:  varType,
		Name:  &ast.Identifier{Token: varName, Value: varName.Literal},
		Value: nil,
	}

	p.nextToken()
	p.nextToken()
	varDecl.Value = p.parseExpression(LOWEST)

	return varDecl
}
