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
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Register prefix parsing functions for identifiers, integers,
	// booleans, and the negation operator.
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	numberTypes := []token.Type{token.U8, token.U16, token.U32, token.U64,
		token.I8, token.I16, token.I32, token.I64,
		token.F8, token.F16, token.F32, token.F64}
	for _, numberType := range numberTypes {
		p.registerPrefix(numberType, func() ast.Expression {
			return p.parseNumberType(numberType)
		})
	}
	// Register infix parsing functions for arithmetic operators and
	// comparison operators.
	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	// Register infix parsing function for logical operators
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	// Register parsing function for if expressions
	p.registerPrefix(token.IF, p.parseIfExpression)

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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken.Type == token.ASSIGN {
		p.errors = append(p.errors, "unexpected '='")
		return nil
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
}

// expectPeek checks if the next token is of the expected type.
// If it is, it consumes the token and returns true.
// If it isn't, it generates an error message and returns false.
func (p *Parser) expectPeek(expectedType token.Type) bool {
	if p.peekTokenIs(expectedType) {
		// p.nextToken()
		return true
	} else {
		fmt.Printf("Expected next token to be %s, got %s instead\n", expectedType, p.peekToken.Type)
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
