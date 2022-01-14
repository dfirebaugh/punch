package parser

import (
	"punch/lexer"
)

type Parser struct {
	Builder Builder
}

func New(builder Builder) Parser {
	return Parser{Builder: builder}
}

func (p Parser) Collect(token lexer.Token) {
	switch token.Literal {
	case lexer.EOF:
		p.eof(token)
	case lexer.INT:
		p.int(token)
	case lexer.STRING:
		p.string(token)
	case lexer.ASSIGN:
		p.assign(token)
	case lexer.PLUS:
		p.plus(token)
	case lexer.COMMA:
		p.comma(token)
	case lexer.SEMICOLON:
		break
	case lexer.LPAREN:
		p.openParen(token)
	case lexer.RPAREN:
		p.closedParen(token)
	case lexer.LBRACE:
		p.openBrace(token)
	case lexer.RBRACE:
		p.closedBrace(token)
	case lexer.FUNCTION:
		p.name(token)
	case lexer.LET:
		p.name(token)
	case lexer.PRINT:
		p.name(token)
	case lexer.RETURN:
		p.name(token)
	case lexer.DOT:
		p.dot(token)
	case lexer.ILLEGAL:
		p.error(token)
	default:
		p.error(token)
	}
}

func (p Parser) eof(token lexer.Token) {
	p.HandleEvent(EOF, token)
}
func (p Parser) int(token lexer.Token) {
	p.HandleEvent(INT, token)
}
func (p Parser) assign(token lexer.Token) {
	p.HandleEvent(ASSIGN, token)
}
func (p Parser) comma(token lexer.Token) {
	p.HandleEvent(COMMA, token)
}
func (p Parser) openBrace(token lexer.Token) {
	p.HandleEvent(OPEN_BRACE, token)
}
func (p Parser) closedBrace(token lexer.Token) {
	p.HandleEvent(CLOSED_BRACE, token)
}
func (p Parser) openParen(token lexer.Token) {
	p.HandleEvent(OPEN_PAREN, token)
}
func (p Parser) closedParen(token lexer.Token) {
	p.HandleEvent(CLOSED_PAREN, token)
}
func (p Parser) dot(token lexer.Token) {
	p.HandleEvent(DOT, token)
}
func (p Parser) plus(token lexer.Token) {
	p.HandleEvent(PLUS, token)
}
func (p Parser) dash(token lexer.Token) {
	p.HandleEvent(DASH, token)
}
func (p Parser) colon(token lexer.Token) {
	p.HandleEvent(COLON, token)
}
func (p Parser) name(token lexer.Token) {
	p.HandleEvent(NAME, token)
}
func (p Parser) string(token lexer.Token) {
	p.HandleEvent(STRING, token)
}
func (p Parser) error(token lexer.Token) {
	println("ERROR -> ", token.Literal)
}
func (p Parser) HandleEvent(event ParserEvent, token lexer.Token) {
	println(event, token.String())
}
func (p Parser) HandleEventError(event ParserEvent, token lexer.Token) {}
