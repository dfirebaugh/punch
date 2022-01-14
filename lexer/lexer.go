package lexer

import (
	"strings"
	"text/scanner"
	"unicode"
)

type lexer struct {
	Collector TokenCollector
	scanner   scanner.Scanner
}

func New(collector TokenCollector, source string) lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(source))

	lexer := lexer{
		Collector: collector,
		scanner:   s,
	}

	lexer.run()

	return lexer
}

func (l *lexer) run() {
	for tok := l.scanner.Scan(); tok != scanner.EOF; tok = l.scanner.Scan() {
		l.Collector.Collect(Token{
			Literal:  l.scanner.TokenText(),
			Text:     l.scanner.TokenText(),
			Position: l.scanner.Position},
		)
	}
	l.Collector.Collect(Token{
		Literal:  "",
		Text:     "",
		Position: scanner.Position{Line: -1, Offset: -1},
		Type:     EOF,
	},
	)
}

func (l lexer) evaluateTokenType(token Token) TokenType {
	if l.isSingleChar(token) {
		return l.evaluateSpecialCharacter(token.Literal)
	}
	if l.isString(token) {
		return STRING
	}
	if l.isIdentifier(token) {
		return l.evaluateIdentifier(token.Literal)
	}

	return ILLEGAL
}

func (lexer) evaluateSpecialCharacter(literal string) TokenType {
	switch literal {
	case PLUS:
		return PLUS
	case ASSIGN:
		return ASSIGN
	case COMMA:
		return COMMA
	case DOT:
		return DOT
	case SEMICOLON:
		return SEMICOLON
	case LPAREN:
		return LPAREN
	case RPAREN:
		return RPAREN
	case LBRACE:
		return LBRACE
	case RBRACE:
		return RBRACE
	default:
		return ILLEGAL
	}
}

func (lexer) evaluateIdentifier(literal string) TokenType {
	switch literal {
	case keywords[FUNCTION]:
		return FUNCTION
	case keywords[LET]:
		return LET
	case keywords[PRINT]:
		return PRINT
	case keywords[RETURN]:
		return RETURN
	default:
		return ILLEGAL
	}
}

func (lexer) isIdentRune(ch rune, i int) bool {
	return ch == '%' && i == 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}

func (l lexer) isIdentifier(token Token) bool {
	for i, k := range token.Text {
		if !l.isIdentRune(k, i) {
			return false
		}
	}

	return true
}

func (lexer) isSingleChar(token Token) bool {
	return len(token.Literal) == 1
}

func (lexer) isString(token Token) bool {
	return len(token.Literal) != 0 && token.Literal[0] == '"' && token.Literal[len(token.Literal)-1] == '"'
}
