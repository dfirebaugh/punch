package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/token"
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
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
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
		token.TRUE:       {prefixFn: p.parseBooleanLiteral},
		token.FALSE:      {prefixFn: p.parseBooleanLiteral},
		token.BANG:       {prefixFn: p.parsePrefixExpression},
		token.ASSIGN:     {infixFn: p.parseAssignmentExpression},
		token.MINUS:      {infixFn: p.parseInfixExpression, prefixFn: p.parsePrefixExpression},
		token.PLUS:       {infixFn: p.parseInfixExpression, prefixFn: p.parsePrefixExpression},
		token.ASTERISK:   {infixFn: p.parseInfixExpression},
		token.MOD:        {infixFn: p.parseInfixExpression},
		token.SLASH:      {infixFn: p.parseInfixExpression},
		token.EQ:         {infixFn: p.parseInfixExpression},
		token.NOT_EQ:     {infixFn: p.parseInfixExpression},
		token.LT_EQUALS:  {infixFn: p.parseInfixExpression},
		token.GT_EQUALS:  {infixFn: p.parseInfixExpression},
		token.LT:         {infixFn: p.parseInfixExpression},
		token.GT:         {infixFn: p.parseInfixExpression},
		token.AND:        {infixFn: p.parseInfixExpression},
		token.OR:         {infixFn: p.parseInfixExpression},
		// token.LBRACKET:   {prefixFn: p.parseListLiteral},
		token.APPEND: {prefixFn: p.parseListOperation},
		token.LEN:    {prefixFn: p.parseListOperation},
		token.LPAREN: {infixFn: p.parseFunctionCall},
	}
	numberTypes := []token.Type{
		token.U8, token.U16, token.U32, token.U64,
		token.I8, token.I16, token.I32, token.I64,
		token.F32, token.F64,
	}
	for _, numberType := range numberTypes {
		p.registerPrefix(numberType, func() (ast.Expression, error) {
			return p.parseNumberType()
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

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	var err error
	p.trace("parseExpression", p.curToken.Literal, p.peekToken.Literal)
	if p.isBooleanLiteral() {
		p.trace("parsing bool", p.curToken.Literal, p.peekToken.Literal)
		b, err := p.parseBooleanLiteral()
		if err != nil {
			return nil, err
		}
		if p.isBinaryOperator(p.peekToken) {
			p.trace("parsing infixed expression", p.curToken.Literal, p.peekToken.Literal)
			p.nextToken()
			return p.parseInfixExpression(b)
		}
		p.nextToken()
		return b, nil
	}

	if p.curTokenIs(token.STRING) {
		p.trace("parseExpression - parsing string literal:", p.curToken.Literal)
		return p.parseStringLiteral()
	}

	if p.isNumber() {
		p.trace("parsing number", p.curToken.Literal, p.peekToken.Literal)
		var n ast.Expression
		if p.curTokenIs(token.NUMBER) {
			n, err = p.parseNumberType()
			if err != nil {
				return nil, err
			}
		}
		if p.curTokenIs(token.FLOAT) {
			n, err = p.parseFloatType()
			if err != nil {
				return nil, err
			}
		}

		if n == nil {
			return nil, p.error("could not parse number")
		}
		if p.isBinaryOperator(p.peekToken) {
			p.trace("parsing infixed expression", p.curToken.Literal, p.peekToken.Literal)
			p.nextToken()
			return p.parseInfixExpression(n)
		}
		if p.curTokenIs(token.NUMBER) {
			p.nextToken()
		}
		return n, nil
	}

	if p.isIndexExpression() {
		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if ident == nil {
			return nil, p.error("identifier is nil")
		}
		return p.parseIndexExpression(ident)
	}
	if p.curToken.Type == token.IDENTIFIER {
		if p.isBinaryOperator(p.peekToken) {
			p.trace("parsing identifier infix expression", p.curToken.Literal, p.peekToken.Literal)
			ident, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if ident == nil {
				return nil, p.error("identifier is nil")
			}
			p.nextToken()
			return p.parseInfixExpression(ident)
		}
		if p.isAssignmentExpression() {
			ident, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if ident == nil {
				return nil, p.error("identifier is nil")
			}
			p.nextToken()
			return p.parseAssignmentExpression(ident)
		}

		if p.isStructLiteral() {
			return p.parseStructLiteral()
		}
		if p.isFunctionCall() {
			p.trace("parsing identifier functioncall expression", p.curToken.Literal, p.peekToken.Literal)
			ident, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			if ident == nil {
				return nil, p.error("identifier is nil")
			}
			fnCall, err := p.parseFunctionCall(ident)
			if err != nil {
				return nil, err
			}
			p.trace("parsed identifier functioncall expression", p.curToken.Literal, p.peekToken.Literal)

			if p.isBinaryOperator(p.curToken) {
				p.trace("parsed identifier functioncall expression - is binary operator", p.curToken.Literal, p.peekToken.Literal)
				return p.parseInfixExpression(fnCall)
			}
			return fnCall, nil
		}

		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if ident == nil {
			return nil, p.error("identifier is nil")
		}
		if p.isStructAccess() {
			return p.parseStructFieldAccess(ident.(*ast.Identifier))
		}

		p.nextToken()
		return ident, nil
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return nil, nil
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil, p.errorf("no prefix parse function for %s", p.curToken.Literal)
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp, nil
		}
		p.nextToken()
		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
