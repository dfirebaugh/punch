package lexer

import (
	"fmt"
	"text/scanner"
)

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	Text     string
	Position scanner.Position
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// literals
	INT    = "INT"
	STRING = "STRING"

	// operators
	ASSIGN = "="
	PLUS   = "+"

	// delimeters
	COMMA     = ","
	SEMICOLON = ";"
	DOT       = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	PRINT    = "PRINT"
	RETURN   = "RETURN"
)

var keywords = map[TokenType]string{
	RETURN:   "return",
	FUNCTION: "function",
	LET:      "let",
	PRINT:    "print",
}

// String will format a token as a String so that it can pretty print.
func (t Token) String() string {
	return fmt.Sprintf("%d:%d: %s", t.Position.Line, t.Position.Offset, t.Text)
}
