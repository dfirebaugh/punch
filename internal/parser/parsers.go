package parser

import (
	"fmt"
	"strconv"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

func (p *Parser) ParseProgram(filename string) *ast.Program {
	program := &ast.Program{}
	program.Files = []*ast.File{}

	for !p.curTokenIs(token.EOF) {
		file := p.parseFile(filename)
		if file != nil {
			program.Files = append(program.Files, file)
		}
	}

	return program
}

func (p *Parser) parseFile(filename string) *ast.File {
	file := &ast.File{
		Filename: filename,
	}

	if !p.expectCurrentTokenIs(token.PACKAGE) {
		p.error("expected 'pkg' keyword")
		return nil
	}
	p.nextToken()
	if !p.expectCurrentTokenIs(token.IDENTIFIER) {
		p.error("expected package name")
		return nil
	}
	file.PackageName = p.curToken.Literal
	p.nextToken()

	if p.curTokenIs(token.IMPORT) {
		p.nextToken()
		if p.curTokenIs(token.LPAREN) {
			p.nextToken()
			for !p.curTokenIs(token.RPAREN) {
				if !p.expectCurrentTokenIs(token.STRING) {
					p.error("expected import path")
					return nil
				}
				file.Imports = append(file.Imports, p.curToken.Literal)
				p.nextToken()
			}
			p.nextToken()
		} else {
			if !p.expectCurrentTokenIs(token.STRING) {
				p.error("expected import path")
				return nil
			}
			file.Imports = append(file.Imports, p.curToken.Literal)
			p.nextToken()
		}
	}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			file.Statements = append(file.Statements, stmt)
		}
		p.nextToken()
	}

	return file
}

func (p *Parser) parseStatement() ast.Statement {
	p.trace("parsing statement", p.curToken.Literal, p.peekToken.Literal)
	if p.curTokenIs(token.PUB) && p.isTypeToken(p.peekToken) && p.peekTokenAfter(token.IDENTIFIER) || p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.LPAREN) {
		return p.parseFunctionStatement()
	}
	if p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.ASSIGN) {
		p.trace("parsing variable declaration", p.curToken.Literal, p.peekToken.Literal)
		s := p.parseTypeBasedVariableDeclaration()
		p.trace("after parsing variable declaration", p.curToken.Literal, p.peekToken.Literal)
		return s
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
		token.F32,
		token.F64:
		return true
	}

	_, exists := p.definedTypes[t.Literal]
	return exists
}

func (p *Parser) parseVariableDeclarationOrAssignment() ast.Statement {
	typeToken := p.curToken
	p.trace("parsing variable declaration or assignment", typeToken.Literal, p.curToken.Literal, string(p.curToken.Type))
	if !p.curTokenIs(token.IDENTIFIER) {
		p.error("expected identifier after type")
		return nil
	}
	name := p.parseIdentifier()

	p.nextToken()
	if !p.curTokenIs(token.ASSIGN) {
		p.error("expected '=' after identifier")
		return nil
	}

	p.nextToken()
	value := p.parseExpression(LOWEST)

	return &ast.VariableDeclaration{
		Type:  typeToken,
		Name:  name.(*ast.Identifier),
		Value: value,
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression != nil {
		p.trace("after parsing expression statement", stmt.Expression.String())
	}

	return stmt
}

func (p *Parser) parseNumberType() ast.Expression {
	d, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.error("could not parse number")
		return nil
	}
	return &ast.IntegerLiteral{
		Token: token.Token{
			Type:     token.I64,
			Literal:  p.curToken.Literal,
			Position: p.curToken.Position,
		},
		Value: int64(d),
	}
}

func (p *Parser) parseFloatType() ast.Expression {
	d, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.error("could not parse float")
		return nil
	}

	return &ast.FloatLiteral{
		Token: p.curToken,
		Value: d,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.DOT) {
		return p.parseStructFieldAccess(ident)
	}

	return ident
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

	expression.Right = p.parseExpression(PREFIX)

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

func (p *Parser) parseIfStatement() *ast.IfStatement {
	p.enterControlStatement()
	defer p.exitControlStatement()

	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken() // consume if
	stmt.Condition = p.parseExpression(LOWEST)
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectCurrentTokenIs(token.LBRACE) {
		p.error("expecting current token is a '{'", p.curToken.Literal, p.peekToken.Literal)
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()
	p.trace("parsing if statement", p.curToken.Literal, p.peekToken.Literal)
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume else

		if !p.expectCurrentTokenIs(token.LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return stmt
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken() // consume LBRACE
	}
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
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
	default:
		return
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
	p.trace("parse assignment", p.curToken.Literal, p.peekToken.Literal)

	println(p.curToken.Literal, p.curTokenIs(token.INFER))
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
	if !p.isTypeToken(p.curToken) && p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.INFER) {
		decl := &ast.VariableDeclaration{
			Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			Value: nil,
		}
		p.nextToken() // consume identifier
		p.nextToken() // consume INFER
		decl.Value = p.parseExpression(LOWEST)
		t := p.inferType(decl.Value)
		decl.Type = token.Token{
			Type:    t,
			Literal: string(t),
		}
		return decl
	}
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

func (p *Parser) parseForStatement() *ast.ForStatement {
	p.trace("parsing for statement", p.curToken.Literal, p.peekToken.Literal)

	// Create a ForStatement node in your AST
	stmt := &ast.ForStatement{Token: p.curToken}

	// Consume the "for" token so p.curToken is now what's after "for"
	p.nextToken()

	// -- Step 1: Parse the init statement --------------------------------
	//    e.g., "i32 i = 1;"
	stmt.Init = p.parseStatement()
	if stmt.Init == nil {
		p.error("expected initialization statement in for loop")
		return nil
	}

	p.trace("current token", p.curToken.Literal)

	// After parsing the init statement, we must see a semicolon.
	// For example, "i32 i = 1;" must end with ';'
	if !p.curTokenIs(token.SEMICOLON) {
		p.error("expected ';' after for-init statement")
		return nil
	}
	// Consume the ';'
	p.nextToken()

	p.trace("current token", p.curToken.Literal)
	// -- Step 2: Parse the condition -------------------------------------
	//    e.g., "i <= n;"
	// Condition can be empty in Go-like syntax, but your grammar might require it.
	// We'll parse an expression if we don't see another semicolon right away.
	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
		if stmt.Condition == nil {
			p.error("expected condition expression in for loop")
			return nil
		}
	}

	p.trace("current token", p.curToken.Literal)
	if !p.curTokenIs(token.SEMICOLON) {
		p.error("expected ';' after for-loop condition")
		return nil
	}
	// consume the second ';'
	p.nextToken()

	// -- Step 3: Parse the post statement --------------------------------
	//    e.g., "i = i + 1"
	// This can also be empty in many “Go-like” syntaxes, but typically you have something.
	if !p.curTokenIs(token.LBRACE) {
		stmt.Post = p.parseStatement()
	}

	// -- Step 4: Finally, we expect '{' to parse the for-block -----------
	if !p.curTokenIs(token.LBRACE) {
		p.error("expected '{' after for loop post statement")
		return nil
	}

	// parseBlockStatement parses everything in { ... } until matching '}'
	stmt.Body = p.parseBlockStatement()

	if !p.curTokenIs(token.LBRACE) {
		p.nextToken()
	}

	return stmt
}
