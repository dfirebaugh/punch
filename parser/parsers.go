package parser

import (
	"fmt"
	"strconv"

	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *Parser) ParseProgram(filename string) (*ast.Program, error) {
	program := &ast.Program{}
	program.Files = []*ast.File{}

	for !p.curTokenIs(token.EOF) {
		file, err := p.parseFile(filename)
		if err != nil {
			return nil, err
		}
		if file != nil {
			program.Files = append(program.Files, file)
		}
	}

	return program, nil
}

func (p *Parser) parseFile(filename string) (*ast.File, error) {
	file := &ast.File{
		Filename: filename,
	}

	if !p.expectCurrentTokenIs(token.PACKAGE) {
		return nil, p.error("expected 'pkg' keyword")
	}
	p.nextToken()
	if !p.expectCurrentTokenIs(token.IDENTIFIER) {
		return nil, p.error("expected package name")
	}
	file.PackageName = p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(token.IMPORT) {
		if err := p.parseImports(file); err != nil {
			return nil, p.errorf("expected import path: %s", err)
		}
	}

	for !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			file.Statements = append(file.Statements, stmt)
		}
		p.nextToken()
	}

	return file, nil
}

func (p *Parser) parseImports(file *ast.File) error {
	p.nextToken()
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		for !p.curTokenIs(token.RPAREN) {
			if !p.expectCurrentTokenIs(token.STRING) {
				return p.error("expected import path")
			}
			file.Imports = append(file.Imports, p.curToken.Literal)
			p.nextToken()
		}
		p.nextToken()
	} else {
		if !p.expectCurrentTokenIs(token.STRING) {
			return p.error("expected import path")
		}
		file.Imports = append(file.Imports, p.curToken.Literal)
		p.nextToken()
	}

	return nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	p.trace("parsing statement", p.curToken.Literal, p.peekToken.Literal)

	if p.curTokenIs(token.PUB) && p.isTypeToken(p.peekToken) && p.peekTokenAfter(token.IDENTIFIER) || p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.LPAREN) {
		return p.parseFunctionStatement()
	}
	if p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.ASSIGN) {
		p.trace("parsing variable declaration", p.curToken.Literal, p.peekToken.Literal)
		s, err := p.parseTypeBasedVariableDeclaration()
		if err != nil {
			return nil, err
		}
		p.trace("after parsing variable declaration", p.curToken.Literal, p.peekToken.Literal)
		return s, nil
	}

	switch p.curToken.Type {
	case token.TRUE:
		return p.parseExpressionStatement()
	case token.FALSE:
		return p.parseExpressionStatement()
	case token.STRUCT:
		return p.parseStructDefinition()
	case token.SLASH_SLASH:
		p.parseComment()
		return nil, nil
	case token.SLASH_ASTERISK:
		p.parseComment()
		return nil, nil
	case token.PUB:
		return p.parseFunctionStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.DEFER:
		return p.parseDeferStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.IDENTIFIER:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseVariableDeclarationOrAssignment()
		}
		return p.parseExpressionStatement()
	default:
		p.trace("parsing statement - default - expression statement", p.curToken.Literal, p.peekToken.Literal)
		if p.curTokenIs(token.RBRACE) && p.peekTokenIs(token.EOF) {
			return nil, nil
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
		token.F32,
		token.F64:
		return true
	}

	_, exists := p.definedTypes[t.Literal]
	return exists
}

func (p *Parser) parseVariableDeclarationOrAssignment() (ast.Statement, error) {
	typeToken := p.curToken
	p.trace("parsing variable declaration or assignment", typeToken.Literal, p.curToken.Literal, string(p.curToken.Type))
	if !p.curTokenIs(token.IDENTIFIER) {
		return nil, p.error("expected identifier after type")
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.nextToken()
	if !p.curTokenIs(token.ASSIGN) {
		return nil, p.error("expected '=' after identifier")
	}

	p.nextToken()
	value, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	return &ast.VariableDeclaration{
		Type:  typeToken,
		Name:  name.(*ast.Identifier),
		Value: value,
	}, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	var err error
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if stmt.Expression != nil {
		p.trace("after parsing expression statement", stmt.Expression.String())
	}

	return stmt, nil
}

func (p *Parser) parseNumberType() (ast.Expression, error) {
	d, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil, p.error("could not parse number")
	}
	return &ast.IntegerLiteral{
		Token: token.Token{
			Type:     token.I64,
			Literal:  p.curToken.Literal,
			Position: p.curToken.Position,
		},
		Value: int64(d),
	}, nil
}

func (p *Parser) parseFloatType() (ast.Expression, error) {
	d, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, p.error("could not parse float")
	}

	return &ast.FloatLiteral{
		Token: p.curToken,
		Value: d,
	}, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.DOT) {
		return p.parseStructFieldAccess(ident)
	}

	return ident, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		return nil, p.error(msg)
	}

	lit.Value = value

	return lit, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	if !p.curTokenIs(token.STRING) {
		return nil, p.error("expected string literal")
	}

	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()

	return lit, nil
}

func (p *Parser) parseBooleanLiteral() (ast.Expression, error) {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	var err error
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken,
	}

	p.nextToken()

	expression.Right, err = p.parseExpression(PREFIX)

	return expression, err
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	var err error
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: p.curToken,
		Right:    nil,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right, err = p.parseExpression(precedence)
	return expression, err
}

func (p *Parser) parseIfStatement() (*ast.IfStatement, error) {
	var err error
	p.enterControlStatement()
	defer p.exitControlStatement()

	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken() // consume if
	stmt.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectCurrentTokenIs(token.LBRACE) {
		return nil, p.error("expecting current token is a '{'", p.curToken.Literal, p.peekToken.Literal)
	}
	stmt.Consequence, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	p.trace("parsing if statement", p.curToken.Literal, p.peekToken.Literal)
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume else

		if !p.expectCurrentTokenIs(token.LBRACE) {
			return nil, p.error("expected left brace")
		}

		stmt.Alternative, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) isInControlStatement() bool {
	return p.controlDepth > 0
}

func (p *Parser) enterControlStatement() {
	p.controlDepth++
}

func (p *Parser) exitControlStatement() {
	p.controlDepth--
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken() // consume LBRACE
	}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block, nil
}

// parseComment skips over a comment token and moves the parser to the end of the comment.
func (p *Parser) parseComment() {
	switch p.curToken.Type {
	case token.SLASH_SLASH:
		currentLine := p.curToken.Position.Line
		for p.curToken.Position.Line == currentLine {
			p.nextToken()
		}
	case token.SLASH_ASTERISK:
		p.nextToken()

		for !(p.curTokenIs(token.ASTERISK) && p.peekTokenIs(token.SLASH)) {
			if p.curTokenIs(token.EOF) {
				return
			}
			p.nextToken()
		}

		p.nextToken()
	default:
		return
	}
}

func (p *Parser) parseDeferStatement() (*ast.DeferStatement, error) {
	var err error
	deferStmt := &ast.DeferStatement{Token: p.curToken}

	p.nextToken()
	deferStmt.Statement, err = p.parseStatement()
	if err != nil {
		return nil, err
	}

	p.deferStack = append(p.deferStack, deferStmt.Statement)

	return deferStmt, nil
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) (ast.Expression, error) {
	var err error
	p.trace("parse assignment", p.curToken.Literal, p.peekToken.Literal)

	println(p.curToken.Literal, p.curTokenIs(token.INFER))
	if _, ok := left.(*ast.Identifier); !ok {
		return nil, p.error("left-hand side of assignment must be an identifier")
	}

	expression := &ast.AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Right, err = p.parseExpression(LOWEST)

	return expression, err
}

func (p *Parser) parseTypeBasedVariableDeclaration() (ast.Statement, error) {
	var err error
	if !p.isTypeToken(p.curToken) && p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.INFER) {
		decl := &ast.VariableDeclaration{
			Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			Value: nil,
		}
		p.nextToken() // consume identifier
		p.nextToken() // consume INFER
		decl.Value, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		t := p.inferType(decl.Value)
		decl.Type = token.Token{
			Type:    t,
			Literal: string(t),
		}
		return decl, nil
	}
	varType := p.curToken
	if !p.expectPeek(token.IDENTIFIER) {
		return nil, p.error("expected identifier")
	}
	p.nextToken()
	varName := p.curToken

	if !p.expectPeek(token.ASSIGN) {
		return nil, p.error("expected assigment operator")
	}
	varDecl := &ast.VariableDeclaration{
		Type:  varType,
		Name:  &ast.Identifier{Token: varName, Value: varName.Literal},
		Value: nil,
	}

	p.nextToken()
	p.nextToken()
	varDecl.Value, err = p.parseExpression(LOWEST)

	return varDecl, err
}

func (p *Parser) inferType(value ast.Expression) token.Type {
	switch value.(type) {
	case *ast.IntegerLiteral:
		return token.I32
	case *ast.FloatLiteral:
		return token.F32
	case *ast.StringLiteral:
		return token.STRING
	case *ast.BooleanLiteral:
		return token.BOOL
	default:
		return token.ILLEGAL
	}
}

func (p *Parser) parseForStatement() (*ast.ForStatement, error) {
	var err error
	p.trace("parsing for statement", p.curToken.Literal, p.peekToken.Literal)

	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	stmt.Init, err = p.parseStatement()
	if err != nil {
		return nil, err
	}
	if stmt.Init == nil {
		return nil, p.error("expected initialization statement in for loop")
	}

	p.trace("current token", p.curToken.Literal)

	if !p.curTokenIs(token.SEMICOLON) {
		return nil, p.error("expected ';' after for-init statement")
	}
	// Consume the ';'
	p.nextToken()

	p.trace("current token", p.curToken.Literal)
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Condition, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		if stmt.Condition == nil {
			return nil, p.error("expected condition expression in for loop")
		}
	}

	p.trace("current token", p.curToken.Literal)
	if !p.curTokenIs(token.SEMICOLON) {
		return nil, p.error("expected ';' after for-loop condition")
	}
	// consume the second ';'
	p.nextToken()

	if !p.curTokenIs(token.LBRACE) {
		stmt.Post, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	if !p.curTokenIs(token.LBRACE) {
		return nil, p.error("expected '{' after for loop post statement")
	}

	stmt.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	if !p.curTokenIs(token.LBRACE) {
		p.nextToken()
	}

	return stmt, nil
}
