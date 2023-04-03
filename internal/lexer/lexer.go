package lexer

import (
	"strings"
	"text/scanner"

	"github.com/dfirebaugh/punch/internal/token"
)

type Lexer struct {
	Collector token.TokenCollector
	scanner   scanner.Scanner
}

func New(source string) *Lexer {
	var s scanner.Scanner

	s.Init(strings.NewReader(source))

	lexer := &Lexer{
		Collector: &Collector{},
		scanner:   s,
	}

	return lexer
}

func (l *Lexer) Run() []token.Token {
	var tokens []token.Token
	for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
		tokens = append(tokens, t)
	}

	tokens = append(tokens, token.Token{Type: token.EOF})
	return tokens
}

func (l *Lexer) NextToken() token.Token {
	tok := l.scanner.Scan()
	if tok == scanner.EOF {
		return token.Token{
			Literal:  "",
			Position: scanner.Position{Line: -1, Offset: -1},
			Type:     token.EOF,
		}
	}

	t := token.Token{
		Literal:  l.scanner.TokenText(),
		Position: l.scanner.Position,
	}
	t.Type = l.evaluateType(t)
	if l.isMultiCharOperator(t.Type) {
		t.Literal = string(t.Type)
	}
	l.Collector.Collect(t)

	return t
}

func (l *Lexer) evaluateType(t token.Token) token.Type {
	switch {
	case l.isSpecialCharacter(t.Literal):
		if m := l.evaluateMultiCharOperators(t.Literal); m != token.UNKNOWN {
			return m
		}
		return l.evaluateSpecialCharacter(t.Literal)
	case t.IsInt():
		return token.INT
	case t.IsSingleCharIdentifier():
		return token.IDENTIFIER
	case t.IsFloatInt():
		return token.FLOAT
	case t.IsString():
		return token.STRING
	case t.IsIdentifier():
		return l.evaluateKeyword(t.Literal)
	default:
		return token.ILLEGAL
	}
}

func (l Lexer) isMultiCharOperator(t token.Type) bool {
	return t == token.PLUS_EQUALS || t == token.MINUS_EQUALS || t == token.ASTERISK_EQUALS || t == token.SLASH_EQUALS || t == token.AND || t == token.OR || t == token.EQ || t == token.NOT_EQ || t == token.LT_EQUALS || t == token.GT_EQUALS
}

func (l Lexer) isSpecialCharacter(literal string) bool {
	return l.evaluateSpecialCharacter(literal) != token.ILLEGAL
}

func (l *Lexer) evaluateMultiCharOperators(literal string) token.Type {
	switch literal {
	case token.ASSIGN:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.EQ
		}
		return token.ASSIGN
	case token.BANG:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.NOT_EQ
		}
		return token.UNKNOWN
	case token.GT:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.GT_EQUALS
		}
		return token.GT
	case token.LT:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.LT_EQUALS
		}
		return token.LT
	case token.PIPE:
		if l.scanner.Peek() == rune('|') {
			l.scanner.Scan()
			return token.OR
		}
		return token.UNKNOWN
	case token.AMPERSAND:
		if l.scanner.Peek() == rune('&') {
			l.scanner.Scan()
			return token.AND
		}
		return token.UNKNOWN
	case token.SLASH:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.SLASH_EQUALS
		}
		if l.scanner.Peek() == rune('/') {
			return token.SLASH_SLASH
		}
		if l.scanner.Peek() == rune('*') {
			return token.SLASH_ASTERISK
		}
		return token.SLASH
	case token.ASTERISK:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.ASTERISK_EQUALS
		}
		if l.scanner.Peek() == rune('/') {
			l.scanner.Scan()
			return token.ASTERISK_SLASH
		}
		return token.ASTERISK
	case token.PLUS:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.PLUS_EQUALS
		}
		return token.PLUS
	case token.MINUS:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.MINUS_EQUALS
		}
		return token.MINUS
	default:
		return token.UNKNOWN
	}
}

func (l Lexer) evaluateSpecialCharacter(literal string) token.Type {
	switch literal {
	case token.PLUS:
		return token.PLUS
	case token.MINUS:
		return token.MINUS
	case token.ASTERISK:
		return token.ASTERISK
	case token.SLASH:
		return token.SLASH
	case token.GT:
		return token.GT
	case token.LT:
		return token.LT
	case token.ASSIGN:
		return token.ASSIGN
	case token.COMMA:
		return token.COMMA
	case token.DOT:
		return token.DOT
	case token.SEMICOLON:
		return token.SEMICOLON
	case token.LPAREN:
		return token.LPAREN
	case token.RPAREN:
		return token.RPAREN
	case token.LBRACE:
		return token.LBRACE
	case token.RBRACE:
		return token.RBRACE
	case token.BANG:
		return token.BANG
	default:
		return token.ILLEGAL
	}
}

func (l *Lexer) evaluateKeyword(literal string) token.Type {
	switch literal {
	case token.Keywords[token.FUNCTION]:
		return token.FUNCTION
	case token.Keywords[token.LET]:
		return token.LET
	case token.Keywords[token.RETURN]:
		return token.RETURN
	case token.Keywords[token.IF]:
		return token.IF
	case token.Keywords[token.TRUE]:
		return token.TRUE
	case token.Keywords[token.FALSE]:
		return token.FALSE
	case token.Keywords[token.PUB]:
		return token.PUB
	default:
		return token.IDENTIFIER
	}
}
