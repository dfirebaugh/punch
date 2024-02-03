package parser

import (
	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn

	deferStack []ast.Statement

	definedTypes      map[string]bool
	structDefinitions map[string]*ast.StructDefinition

	controlDepth int
}

type parseRule struct {
	prefixFn prefixParseFn
	infixFn  infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type ParseFn interface {
	Parse(p *Parser) ast.Expression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:            l,
		errors:       []string{},
		definedTypes: make(map[string]bool),
	}
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.structDefinitions = make(map[string]*ast.StructDefinition)

	p.registerParseRules()

	// Read two tokens to initialize curToken and peekToken.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerParseRules() {
	parseRules := map[token.Type]parseRule{
		token.IDENTIFIER: {prefixFn: p.parseIdentifier},
		token.STRING:     {prefixFn: p.parseStringLiteral},
		token.NUMBER:     {prefixFn: p.parseIntegerLiteral},
		token.TRUE:       {prefixFn: p.parseBooleanLiteral},
		token.FALSE:      {prefixFn: p.parseBooleanLiteral},
		token.BANG:       {prefixFn: p.parsePrefixExpression},
		token.ASSIGN:     {infixFn: p.parseAssignmentExpression},
		token.MINUS:      {infixFn: p.parseInfixExpression},
		token.PLUS:       {infixFn: p.parseInfixExpression},
		token.ASTERISK:   {infixFn: p.parseInfixExpression},
		token.SLASH:      {infixFn: p.parseInfixExpression},
		token.EQ:         {infixFn: p.parseInfixExpression},
		token.NOT_EQ:     {infixFn: p.parseInfixExpression},
		token.LT:         {infixFn: p.parseInfixExpression},
		token.GT:         {infixFn: p.parseInfixExpression},
		token.AND:        {infixFn: p.parseInfixExpression},
		token.OR:         {infixFn: p.parseInfixExpression},
	}
	numberTypes := []token.Type{
		token.U8, token.U16, token.U32, token.U64,
		token.I8, token.I16, token.I32, token.I64,
		token.F32, token.F64}
	for _, numberType := range numberTypes {
		p.registerPrefix(numberType, func() ast.Expression {
			return p.parseNumberType(numberType)
		})
	}
	for tokenType, rule := range parseRules {
		if rule.prefixFn != nil {
			p.registerPrefix(tokenType, rule.prefixFn)
		}
		if rule.infixFn != nil {
			p.registerInfix(tokenType, rule.infixFn)
		}
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken.Literal == token.TRUE || p.curToken.Literal == token.FALSE {
		p.trace("parsing bool", p.curToken.Literal, p.peekToken.Literal)
		b := p.parseBooleanLiteral()
		if p.peekToken.Type == token.MINUS ||
			p.peekToken.Type == token.PLUS ||
			p.peekToken.Type == token.ASTERISK ||
			p.peekToken.Type == token.SLASH ||
			p.peekToken.Type == token.EQ ||
			p.peekToken.Type == token.NOT_EQ ||
			p.peekToken.Type == token.LT ||
			p.peekToken.Type == token.GT ||
			p.peekToken.Type == token.AND ||
			p.peekToken.Type == token.OR {
			p.trace("parsing infixed expression", p.curToken.Literal, p.peekToken.Literal)
			p.nextToken()
			return p.parseInfixExpression(b)
		}
		p.nextToken()
		return b
	}

	if p.curTokenIs(token.STRING) {
		return p.parseStringLiteral()
	}

	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.FLOAT) {
		p.trace("parsing number", p.curToken.Literal, p.peekToken.Literal)
		var n ast.Expression
		if p.curTokenIs(token.NUMBER) {
			n = p.parseNumberType(p.curToken.Type)
		}
		if p.curTokenIs(token.FLOAT) {
			n = p.parseFloatType()
		}

		if n == nil {
			p.errors = append(p.errors, "could not parse number")
			return nil
		}
		if p.peekToken.Type == token.MINUS ||
			p.peekToken.Type == token.PLUS ||
			p.peekToken.Type == token.ASTERISK ||
			p.peekToken.Type == token.SLASH ||
			p.peekToken.Type == token.EQ ||
			p.peekToken.Type == token.NOT_EQ ||
			p.peekToken.Type == token.LT ||
			p.peekToken.Type == token.GT ||
			p.peekToken.Type == token.AND ||
			p.peekToken.Type == token.OR {
			p.trace("parsing infixed expression", p.curToken.Literal, p.peekToken.Literal)
			p.nextToken()
			return p.parseInfixExpression(n)
		}
		p.nextToken()
		return n
	}

	if p.curToken.Type == token.IDENTIFIER {
		if p.peekToken.Type == token.MINUS ||
			p.peekToken.Type == token.PLUS ||
			p.peekToken.Type == token.ASTERISK ||
			p.peekToken.Type == token.SLASH ||
			p.peekToken.Type == token.EQ ||
			p.peekToken.Type == token.NOT_EQ ||
			p.peekToken.Type == token.LT ||
			p.peekToken.Type == token.GT ||
			p.peekToken.Type == token.AND ||
			p.peekToken.Type == token.OR {
			p.trace("parsing identifier infix expression", p.curToken.Literal, p.peekToken.Literal)
			ident := p.parseIdentifier()
			p.nextToken()
			return p.parseInfixExpression(ident)
		}
		if p.peekToken.Type == token.ASSIGN || p.peekToken.Type == token.INFER {
			ident := p.parseIdentifier()
			p.nextToken()
			return p.parseAssignmentExpression(ident)
		}

		if p.peekTokenIs(token.LBRACE) && !p.isInControlStatement() {
			return p.parseStructLiteral()
		}
		if p.curTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.LPAREN) {
			p.trace("parsing identifier functioncall expression", p.curToken.Literal, p.peekToken.Literal)
			ident := p.parseIdentifier()
			fnCall := p.parseFunctionCall(ident)
			p.trace("parsed identifier functioncall expression", p.curToken.Literal, p.peekToken.Literal)
			if p.curTokenIs(token.RPAREN) {
				p.nextToken()
			}
			return fnCall
		}
		ident := p.parseIdentifier()
		p.nextToken()
		return ident
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		panic("no prefix parse function for " + p.curToken.Literal)
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
}

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
