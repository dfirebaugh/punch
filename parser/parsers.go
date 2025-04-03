package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *parser) ParseProgram(filename string) (*ast.Program, error) {
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

func (p *parser) parseFile(filename string) (*ast.File, error) {
	file := &ast.File{
		Filename:    filename,
		PackageName: "punch_default",
	}

	// if !p.expectCurrentTokenIs(token.PACKAGE) {
	// 	return nil, p.error("expected 'pkg' keyword")
	// }
	// if !p.expectCurrentTokenIs(token.IDENTIFIER) {
	// 	return nil, p.error("expected package name")
	// }
	if p.curTokenIs(token.PACKAGE) {
		p.nextToken()
		if p.curTokenIs(token.IDENTIFIER) {
			file.PackageName = p.curToken.Literal
			p.nextToken()
		}
	}

	if p.curTokenIs(token.IMPORT) {
		if err := p.parseImports(file); err != nil {
			return nil, p.errorf("expected import path: %s", err)
		}
	}

	for !p.curTokenIs(token.EOF) {
		// if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.TAB) {
		// 	p.nextToken()
		// 	continue
		// }
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
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

func (p *parser) parseImports(file *ast.File) error {
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

func (p *parser) parseStatement() (ast.Statement, error) {
	p.trace("parsing statement", p.curToken.Literal, p.peekToken.Literal)
	switch p.curToken.Type {
	case token.PUB:
		return p.parseFunctionStatement()
	case token.FUNCTION:
		println("\t\t\tparsing function statement")
		return p.parseFunctionStatement()
	case token.DEFER:
		return p.parseDeferStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LET:
		return p.parseVariableDeclarationOrAssignment()
	default:
		p.trace("-----------parsing statement - default - expression statement", p.curToken.Literal, p.peekToken.Literal)
		if p.curTokenIs(token.RBRACE) && p.peekTokenIs(token.EOF) {
			return nil, nil
		}
		return p.parseExpressionStatement()
		// return nil, nil
	}
}

func (p *parser) parseVariableDeclarationOrAssignment() (ast.Statement, error) {
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

	if p.curTokenIs(token.ASSIGN) {
		p.nextToken()
	}

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

func (p *parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	var err error
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	p.trace("parseExpressionStatement: before parsing expression", p.curToken.Literal, p.peekToken.Literal)
	stmt.Expression, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	// p.nextToken()

	return stmt, nil
}

func (p *parser) parseStringLiteral() (ast.Expression, error) {
	if !p.curTokenIs(token.STRING) {
		return nil, p.error("expected string literal")
	}

	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	if p.curTokenIs(token.STRING) {
		p.nextToken()
	}
	if p.curTokenIs(token.COMMA) {
		p.trace("after string literal -- consuming comma: ", p.curToken.Literal, p.peekToken.Literal)
		p.nextToken()
	}

	return lit, nil
}

func (p *parser) parseBooleanLiteral() (ast.Expression, error) {
	lit := &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	p.nextToken()
	return lit, nil
}

func (p *parser) parseIfExpression() (ast.Expression, error) {
	p.trace("parsing if statement-start:", p.curToken.Literal, p.peekToken.Literal)
	defer p.trace("parsing if statement-end:", p.curToken.Literal, p.peekToken.Literal)
	var err error
	p.enterControlStatement()
	defer p.exitControlStatement()

	stmt := &ast.IfStatement{Token: p.curToken}

	if !p.expectAndConsumeCurrentTokenIs(token.IF) {
		return nil, p.error("expected 'if'")
	}

	stmt.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.trace("parsing if statement consequence")
	// if !p.expectAndConsumeCurrentTokenIs(token.LBRACE) {
	// 	return nil, p.error("expected '{' after condition")
	// }
	stmt.Consequence, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	println("\t\t\t got out of block can now parse else")
	p.trace("parsing else block")
	if p.curTokenIs(token.ELSE) {
		if !p.expectAndConsumeCurrentTokenIs(token.ELSE) {
			p.nextToken()
		}
		if p.curTokenIs(token.IF) {
			alternative, err := p.parseIfExpression()
			if err != nil {
				return nil, err
			}
			stmt.Alternative = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Expression: alternative},
				},
			}
		} else {
			// if !p.expectAndConsumeCurrentTokenIs(token.LBRACE) {
			// 	return nil, p.error("expected '{' after 'else'")
			// }
			stmt.Alternative, err = p.parseBlockStatement()
			if err != nil {
				return nil, err
			}
		}
	}

	return stmt, nil
}

func (p *parser) enterControlStatement() {
	p.controlDepth++
}

func (p *parser) exitControlStatement() {
	p.controlDepth--
}

func (p *parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.trace("parsing block statement-start:", p.curToken.Literal, p.peekToken.Literal)
	if !p.expectAndConsumeCurrentTokenIs(token.LBRACE) {
		return nil, p.error("expected '{'")
	}

	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.EOF) {
			return nil, p.error("reached EOF while parsing block statement")
		}
		println("in block statement:", p.curToken.Literal, p.peekToken.Literal)
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.trace("----parsing block statement-mid:", p.curToken.Literal, p.peekToken.Literal)
		if !p.curTokenIs(token.RBRACE) {
			// p.nextToken()
		}
	}

	if !p.expectAndConsumeCurrentTokenIs(token.RBRACE) {
		return nil, p.error("expected '}'")
	}

	p.trace("parsing block statement-end:", p.curToken.Literal, p.peekToken.Literal)
	return block, nil
}

// parseComment skips over a comment token and moves the parser to the end of the comment.
func (p *parser) parseComment() {
	switch p.curToken.Type {
	case token.SLASH_SLASH:
		currentLine := p.curToken.Position.Line
		for p.curToken.Position.Line == currentLine {
			p.nextToken()
		}
	case token.SLASH_ASTERISK:
		p.nextToken()

		for !p.curTokenIs(token.ASTERISK) || !p.peekTokenIs(token.SLASH) {
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

func (p *parser) parseDeferStatement() (*ast.DeferStatement, error) {
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

func (p *parser) parseAssignmentExpression(left ast.Expression) (ast.Expression, error) {
	var err error
	p.trace("parse assignment", p.curToken.Literal, p.peekToken.Literal)

	if _, ok := left.(*ast.Identifier); !ok {
		return nil, p.error("left-hand side of assignment must be an identifier")
	}

	expression := &ast.AssignmentExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expression.Right, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	return expression, err
}

func (p *parser) parseTypeBasedVariableDeclaration() (ast.Statement, error) {
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
	// if not an identifier, it could be a struct member
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

func (p *parser) inferType(value ast.Expression) token.Type {
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

func (p *parser) parseForStatement() (*ast.ForStatement, error) {
	var err error
	p.trace("parsing for statement", p.curToken.Literal, p.peekToken.Literal)

	stmt := &ast.ForStatement{Token: p.curToken}

	if p.curTokenIs(token.FOR) {
		p.nextToken()
	}

	p.trace("parsing for loop init expression", p.curToken.Literal, p.peekToken.Literal)
	stmt.Init, err = p.parseStatement()
	if err != nil {
		return nil, err
	}
	if stmt.Init == nil {
		return nil, p.error("expected initialization statement in for loop")
	}

	if p.curTokenIs(token.SEMICOLON) {
		p.trace("consume semi", p.curToken.Literal)
		p.nextToken()
	}

	p.trace("current token", p.curToken.Literal)
	p.trace("parsing for loop condition expression", p.curToken.Literal, p.peekToken.Literal)
	stmt.Condition, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if stmt.Condition == nil {
		return nil, p.error("expected condition expression in for loop")
	}

	if p.curTokenIs(token.SEMICOLON) {
		p.trace("consume semi", p.curToken.Literal)
		p.nextToken()
	}

	p.trace("parsing for loop condition post expression", p.curToken.Literal, p.peekToken.Literal)
	if !p.curTokenIs(token.LBRACE) {
		stmt.Post, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	stmt.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	// if p.curTokenIs(token.RBRACE) {
	// 	p.trace("consuming right brace after for statement", p.curToken.Literal)
	// 	p.nextToken()
	// }

	return stmt, nil
}

func (p *parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	indexExp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken() // consume the identifier
	p.nextToken() // consume '['

	index, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	indexExp.Index = index

	if !p.expectCurrentTokenIs(token.RBRACKET) {
		return nil, p.error("expected ']' after index expression")
	}
	p.nextToken()

	return indexExp, nil
}
