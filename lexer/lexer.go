package lexer

import (
	"strings"
	"text/scanner"

	"github.com/dfirebaugh/punch/token"
)

type Lexer struct {
	Collector    token.TokenCollector
	scanner      scanner.Scanner
	savedScanner scanner.Scanner
}

func New(filename string, source string) *Lexer {
	var s scanner.Scanner

	s.Init(strings.NewReader(source))

	s.Filename = filename

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
	if t.IsString() {
		t.Literal = l.processStringLiteral(t.Literal)
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
	case t.Literal == token.BOOL:
		return token.BOOL
	case t.IsSingleCharIdentifier():
		return token.IDENTIFIER
	case t.IsString():
		return token.STRING
	case t.IsIdentifier():
		return l.evaluateKeyword(t.Literal)
	case t.IsNumber():
		return token.NUMBER
	case t.IsFloat():
		return token.FLOAT
	default:
		return token.ILLEGAL
	}
}

func (l *Lexer) processStringLiteral(literal string) string {
	var result strings.Builder
	escaped := false

	for i := 1; i < len(literal)-1; i++ {
		char := literal[i]

		if escaped {
			switch char {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte(char)
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else {
			result.WriteByte(char)
		}
	}

	return result.String()
}

func (l Lexer) isMultiCharOperator(t token.Type) bool {
	return t == token.PLUS_EQUALS || t == token.MINUS_EQUALS || t == token.ASTERISK_EQUALS || t == token.SLASH_EQUALS || t == token.AND || t == token.OR || t == token.EQ || t == token.NOT_EQ || t == token.LT_EQUALS || t == token.GT_EQUALS
}

func (l Lexer) isSpecialCharacter(literal string) bool {
	return l.evaluateSpecialCharacter(literal) != token.ILLEGAL
}

func (l *Lexer) evaluateMultiCharOperators(literal string) token.Type {
	switch literal {
	case token.COLON:
		if l.scanner.Peek() == rune('=') {
			l.scanner.Scan()
			return token.INFER
		}
		return token.COLON
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
	case token.INFER:
		return token.INFER
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
	case token.COLON:
		return token.COLON
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
	case token.MOD:
		return token.MOD
	case token.LBRACKET:
		return token.LBRACKET
	case token.RBRACKET:
		return token.RBRACKET
	default:
		return token.ILLEGAL
	}
}

func (l *Lexer) evaluateKeyword(literal string) token.Type {
	switch literal {
	case token.U8:
		return token.U8
	case token.U16:
		return token.U16
	case token.U32:
		return token.U32
	case token.U64:
		return token.U64
	case token.I8:
		return token.I8
	case token.I16:
		return token.I16
	case token.I32:
		return token.I32
	case token.I64:
		return token.I64
	case token.F32:
		return token.F32
	case token.F64:
		return token.F64
	case token.STRING:
		return token.STRING
	case token.Keywords[token.STRING]:
		return token.STRING
	case token.Keywords[token.STRUCT]:
		return token.STRUCT
	case token.Keywords[token.INTERFACE]:
		return token.INTERFACE
	case token.Keywords[token.DEFER]:
		return token.DEFER
	case token.Keywords[token.PACKAGE]:
		return token.PACKAGE
	case token.Keywords[token.IMPORT]:
		return token.IMPORT
	case token.Keywords[token.FUNCTION]:
		return token.FUNCTION
	case token.Keywords[token.FN]:
		return token.FUNCTION
	case token.Keywords[token.LET]:
		return token.LET
	case token.Keywords[token.RETURN]:
		return token.RETURN
	case token.Keywords[token.IF]:
		return token.IF
	case token.Keywords[token.TRUE]:
		return token.TRUE
	case token.Keywords[token.FOR]:
		return token.FOR
	case token.Keywords[token.FALSE]:
		return token.FALSE
	case token.Keywords[token.PUB]:
		return token.PUB
	default:
		return token.IDENTIFIER
	}
}

func (l *Lexer) SaveState() {
	l.savedScanner = l.scanner
}

func (l *Lexer) RestoreState() {
	l.scanner = l.savedScanner
}
