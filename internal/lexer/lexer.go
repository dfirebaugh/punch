package lexer

import (
	"punch/internal/token"
	"strings"
	"text/scanner"
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
	l.Collector.Collect(t)

	return t
}

func (l *Lexer) evaluateType(t token.Token) token.Type {
	switch {
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
	case l.isSpecialCharacter(t.Literal):
		return l.evaluateSpecialCharacter(t.Literal)
	default:
		return token.ILLEGAL
	}
}

func (l Lexer) isSpecialCharacter(literal string) bool {
	return l.evaluateSpecialCharacter(literal) != token.ILLEGAL
}

func (l Lexer) evaluateSpecialCharacter(literal string) token.Type {
	switch literal {
	case token.PLUS:
		return token.PLUS
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
	default:
		return token.IDENTIFIER
	}
}
