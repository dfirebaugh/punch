package parser

import (
	"fmt"

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

	definedTypes map[string]bool
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
		token.MINUS: {
			prefixFn: p.parsePrefixExpression,
			infixFn:  p.parseInfixExpression,
		},
		token.LPAREN:   {prefixFn: p.parseGroupedExpression},
		token.ASSIGN:   {infixFn: p.parseAssignmentExpression},
		token.PLUS:     {infixFn: p.parseInfixExpression},
		token.ASTERISK: {infixFn: p.parseInfixExpression},
		token.SLASH:    {infixFn: p.parseInfixExpression},
		token.EQ:       {infixFn: p.parseInfixExpression},
		token.NOT_EQ:   {infixFn: p.parseInfixExpression},
		token.LT:       {infixFn: p.parseInfixExpression},
		token.GT:       {infixFn: p.parseInfixExpression},
		token.AND:      {infixFn: p.parseInfixExpression},
		token.OR:       {infixFn: p.parseInfixExpression},
	}
	numberTypes := []token.Type{token.U8, token.U16, token.U32, token.U64,
		token.I8, token.I16, token.I32, token.I64,
		token.F8, token.F16, token.F32, token.F64}
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
	switch p.curToken.Type {
	case token.ASSIGN:
		p.errors = append(p.errors, "unexpected '='")
		return nil
	case token.IDENTIFIER:
		if p.peekTokenIs(token.LBRACE) {
			return p.parseStructLiteral()
		}
	}

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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
// If it is, it consumes the token and returns true.
// If it isn't, it generates an error message and returns false.
func (p *Parser) expectPeek(expectedType token.Type) bool {
	if p.peekTokenIs(expectedType) {
		return true
	} else {
		fmt.Printf("Expected next token to be %s, got %s instead\n", expectedType, p.peekToken.Type)
		p.peekError(expectedType)
		return false
	}
}

func (p *Parser) expectCurrentTokenIs(expectedType token.Type) bool {
	if p.curTokenIs(expectedType) {
		return true
	} else {
		fmt.Printf("Expected current token to be %s, got %s instead\n", expectedType, p.curToken.Type)
		p.peekError(expectedType)
		return false
	}
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

// peekError generates an error message for unexpected tokens.
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
