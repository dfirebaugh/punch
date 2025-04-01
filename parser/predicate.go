package parser

import "github.com/dfirebaugh/punch/token"

// expectPeek checks if the next token is of the expected type.
func (p *Parser) expectPeek(expectedType token.Type) bool {
	return p.peekTokenIs(expectedType)
}

func (p *Parser) expectCurrentTokenIs(expectedType token.Type) bool {
	return p.curTokenIs(expectedType)
}

// peekTokenAfter checks if the token after the expectedType token is of a specific type.
// It temporarily advances the parser to check the token and then restores the parser's state.
func (p *Parser) peekTokenAfter(expectedType token.Type) bool {
	curToken := p.curToken
	peekToken := p.peekToken
	p.l.SaveState()

	p.nextToken()
	result := p.peekToken.Type == expectedType

	p.curToken = curToken
	p.peekToken = peekToken

	p.l.RestoreState()
	return result
}

// curTokenIs checks if the current token is of the specified type.
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks if the peek token is of the specified type.
func (p *Parser) peekTokenIs(t token.Type) bool {
	return string(p.peekToken.Type) == string(t)
}

func (p *Parser) isFunctionDeclaration() bool {
	return p.curTokenIs(token.FN) || p.curTokenIs(token.PUB) && p.isTypeToken(p.peekToken) && p.peekTokenAfter(token.IDENTIFIER) || p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.LPAREN)
}

func (p *Parser) isVariableDeclaration() bool {
	return p.isTypeToken(p.curToken) && p.peekTokenIs(token.IDENTIFIER) && p.peekTokenAfter(token.ASSIGN)
}

func (p *Parser) isInControlStatement() bool {
	return p.controlDepth > 0
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

func (p *Parser) isBooleanLiteral() bool {
	return p.curToken.Literal == token.TRUE || p.curToken.Literal == token.FALSE
}

func (p *Parser) isBinaryOperator(t token.Token) bool {
	return t.Type == token.MINUS ||
		t.Type == token.PLUS ||
		t.Type == token.ASTERISK ||
		t.Type == token.SLASH ||
		t.Type == token.MOD ||
		t.Type == token.EQ ||
		t.Type == token.NOT_EQ ||
		t.Type == token.LT ||
		t.Type == token.GT ||
		t.Type == token.LT_EQUALS ||
		t.Type == token.GT_EQUALS ||
		t.Type == token.AND ||
		t.Type == token.OR
}

func (p *Parser) isStructLiteral() bool {
	return p.peekTokenIs(token.LBRACE) && !p.isInControlStatement()
}

func (p *Parser) isIdentifier(t token.Type) bool {
	return t == token.IDENTIFIER
}

func (p *Parser) isAssignmentExpression() bool {
	return p.peekToken.Type == token.ASSIGN || p.peekToken.Type == token.INFER
}

func (p *Parser) isFunctionCall() bool {
	return p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.LPAREN)
}

func (p *Parser) isNumber() bool {
	return p.curTokenIs(token.NUMBER) || p.curTokenIs(token.FLOAT)
}

func (p *Parser) isIndexExpression() bool {
	return p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.LBRACKET)
}

func (p *Parser) isStructType(t token.Token) bool {
	_, exists := p.structDefinitions[t.Literal]
	return exists
}

func (p *Parser) isStructAccess() bool {
	return p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.DOT) && p.peekTokenAfter(token.IDENTIFIER)
}
