package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/token"
)

type (
	prefixParseFn  func() (ast.Expression, error)
	infixParseFn   func(ast.Expression) (ast.Expression, error)
	postfixParseFn func(ast.Expression) (ast.Expression, error)
)

type parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns  map[token.Type]prefixParseFn
	infixParseFns   map[token.Type]infixParseFn
	postfixParseFns map[token.Type]postfixParseFn

	deferStack []ast.Statement

	definedTypes      map[string]bool
	structDefinitions map[string]*ast.StructDefinition

	controlDepth int
}

type parseRule struct {
	prefixFn  prefixParseFn
	infixFn   infixParseFn
	postfixFn postfixParseFn
}

func New(l *lexer.Lexer) *parser {
	p := &parser{
		l:                 l,
		errors:            []string{},
		definedTypes:      make(map[string]bool),
		prefixParseFns:    make(map[token.Type]prefixParseFn),
		infixParseFns:     make(map[token.Type]infixParseFn),
		postfixParseFns:   make(map[token.Type]postfixParseFn),
		structDefinitions: make(map[string]*ast.StructDefinition),
	}

	p.registerParseRules()

	// Read two tokens to initialize curToken and peekToken.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *parser) registerPostfix(tokenType token.Type, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *parser) registerParseRules() {
	parseRules := map[token.Type]parseRule{
		token.FUNCTION:   {prefixFn: p.parseFunction},
		token.PUB:        {prefixFn: p.parseFunction},
		token.IDENTIFIER: {prefixFn: p.parseIdentifier},
		token.IF:         {prefixFn: p.parseIfExpression},
		token.TRUE:       {prefixFn: p.parseBooleanLiteral},
		token.FALSE:      {prefixFn: p.parseBooleanLiteral},
		token.BANG:       {prefixFn: p.parsePrefixExpression},
		token.APPEND:     {prefixFn: p.parseListOperation},
		token.LEN:        {prefixFn: p.parseListOperation},
		token.STRUCT:     {prefixFn: p.parseStructLiteral},
		token.NUMBER:     {prefixFn: p.parseNumberExpression},
		token.STRING:     {prefixFn: p.parseStringLiteral},
		// token.LPAREN:     {prefixFn: p.parseGroupedExpression},

		token.ASSIGN:    {infixFn: p.parseAssignmentExpression},
		token.MINUS:     {infixFn: p.parseInfixExpression, prefixFn: p.parsePrefixExpression},
		token.PLUS:      {infixFn: p.parseInfixExpression, prefixFn: p.parsePrefixExpression},
		token.ASTERISK:  {infixFn: p.parseInfixExpression},
		token.MOD:       {infixFn: p.parseInfixExpression},
		token.SLASH:     {infixFn: p.parseInfixExpression},
		token.EQ:        {infixFn: p.parseInfixExpression},
		token.NOT_EQ:    {infixFn: p.parseInfixExpression},
		token.LT_EQUALS: {infixFn: p.parseInfixExpression},
		token.GT_EQUALS: {infixFn: p.parseInfixExpression},
		token.LT:        {infixFn: p.parseInfixExpression},
		token.GT:        {infixFn: p.parseInfixExpression},
		token.AND:       {infixFn: p.parseInfixExpression},
		token.OR:        {infixFn: p.parseInfixExpression},
		// token.LPAREN:      {infixFn: p.parseFunctionCall},
		token.LBRACKET:    {infixFn: p.parseIndexExpression},
		token.MINUS_MINUS: {postfixFn: p.parsePostfixExpression},
		token.PLUS_PLUS:   {postfixFn: p.parsePostfixExpression},
	}
	for tokenType, rule := range parseRules {
		if rule.prefixFn != nil {
			p.registerPrefix(tokenType, rule.prefixFn)
		}
		if rule.infixFn != nil {
			p.registerInfix(tokenType, rule.infixFn)
		}
		if rule.postfixFn != nil {
			p.registerPostfix(tokenType, rule.postfixFn)
		}
	}
}

func (p *parser) parseExpression(precedence int) (ast.Expression, error) {
	var err error
	p.trace("parseExpression-start:", p.curToken.Literal, p.peekToken.Literal)
	defer p.trace("parseExpression-end:", p.curToken.Literal, p.peekToken.Literal)

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

	if postfix := p.postfixParseFns[p.curToken.Type]; postfix != nil {
		leftExp, err = postfix(leftExp)
		if err != nil {
			return nil, err
		}
	}

  p.trace("parsed expression: [", leftExp.String(), "], curtokenis:", p.curToken.Literal, "peektokenis:", p.peekToken.Literal)

	return leftExp, nil
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *parser) parsePrefixExpression() (ast.Expression, error) {
	p.trace("parsing prefix expression", p.curToken.Literal, p.peekToken.Literal)
	var err error
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken,
	}
	p.nextToken()

	expression.Right, err = p.parseExpression(PREFIX)

	return expression, err
}

func (p *parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	p.trace("parsing infix expression", p.curToken.Literal, p.peekToken.Literal)
	var err error
	expression := &ast.InfixExpression{
		Left:     left,
		Operator: p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right, err = p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *parser) parsePostfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	p.nextToken()
	return expression, nil
}
