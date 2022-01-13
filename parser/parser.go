package parser

import (
	"punch/lexer"
	"unicode"
)

type Parser struct {
	Builder Builder
}

func New(builder Builder) Parser {
	return Parser{Builder: builder}
}

func (p Parser) Collect(token lexer.Token) {
	p.evaluateToken(token)
}

func (p Parser) evaluateToken(token lexer.Token) {
	if p.isIdentifier(token) {
		p.Name(token)
	} else if p.isString(token) {
		p.String(token)
	} else if p.isSingleChar(token) {
		p.evaluateSingleChar(token)
	} else {
		// p.HandleEventError(event, token)
	}
}

func (p Parser) evaluateSingleChar(token lexer.Token) {
	switch token.Text[0] {
	case '{':
		p.OpenBrace(token)
	case '}':
		p.ClosedBrace(token)
	case '(':
		p.OpenParen(token)
	case ')':
		p.ClosedParen(token)
	case '.':
		p.Dot(token)
	case ';':
	default:
	}
}

func (p Parser) isIdentRune(ch rune, i int) bool {
	return ch == '%' && i == 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}

func (p Parser) isIdentifier(token lexer.Token) bool {
	for i, k := range token.Text {
		if !p.isIdentRune(k, i) {
			return false
		}
	}

	return true
}

func (p Parser) isString(token lexer.Token) bool {
	return token.Text[0] == '"' && token.Text[len(token.Text)-1] == '"'
}

func (p Parser) isSingleChar(token lexer.Token) bool {
	return len(token.Text) == 1
}

func (p Parser) OpenBrace(token lexer.Token) {
	p.HandleEvent(OPEN_BRACE, token)
	println(token.String(), "-> openBrace")
}
func (p Parser) ClosedBrace(token lexer.Token) {
	p.HandleEvent(CLOSED_BRACE, token)
	println(token.String(), "-> closedBrace")
}
func (p Parser) OpenParen(token lexer.Token) {
	p.HandleEvent(OPEN_PAREN, token)
	println(token.String(), "-> openParen")
}
func (p Parser) ClosedParen(token lexer.Token) {
	p.HandleEvent(CLOSED_PAREN, token)
	println(token.String(), "-> closeParen")
}
func (p Parser) OpenAngle(token lexer.Token) {
	p.HandleEvent(OPEN_ANGLE, token)
	println(token.String(), "-> openangel")
}
func (p Parser) ClosedAngel(token lexer.Token) {
	p.HandleEvent(CLOSED_ANGLE, token)
	println(token.String(), "-> closedangel")
}
func (p Parser) Dot(token lexer.Token) {
	p.HandleEvent(DOT, token)
	println(token.String(), "-> dot")
}
func (p Parser) Dash(token lexer.Token) {
	p.HandleEvent(DASH, token)
	println(token.String(), "-> dash")
}
func (p Parser) Colon(token lexer.Token) {
	p.HandleEvent(COLON, token)
	println(token.String(), "-> colon")
}
func (p Parser) Name(token lexer.Token) {
	println(token.String(), "-> identifier")
	p.HandleEvent(NAME, token)
}
func (p Parser) String(token lexer.Token) {
	println(token.String(), "-> string")
	p.HandleEvent(STRING, token)
}
func (p Parser) Error(token lexer.Token) {
	// p.Builder.SyntaxError(token.Position)
}
func (p Parser) HandleEvent(event ParserEvent, token lexer.Token) {
	// println(event, token.String())
}
func (p Parser) HandleEventError(event ParserEvent, token lexer.Token) {}
