package token

import (
	"text/scanner"
)

type Type string

type Token struct {
	Type     Type
	Literal  string
	Position scanner.Position
}

type TokenCollector interface {
	Collect(token Token)
	Result() []Token
}

// Token types
const (
	ILLEGAL = "ILLEGAL"
	UNKNOWN = "UNKNOWN"
	EOF     = "EOF"

	// Literals
	FLOAT  = "FLOAT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN          = "="
	PLUS            = "+"
	PLUS_EQUALS     = "+="
	MINUS           = "-"
	MINUS_EQUALS    = "-="
	ASTERISK        = "*"
	ASTERISK_EQUALS = "*="
	SLASH           = "/"
	SLASH_EQUALS    = "/="
	MOD             = "%"
	AND             = "&&"
	AMPERSAND       = "&"
	OR              = "||"
	PIPE            = "|"
	EQ              = "=="
	NOT_EQ          = "!="
	LT              = "<"
	LT_EQUALS       = "<="
	GT              = ">"
	GT_EQUALS       = ">="
	QUESTION        = "?"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "if"
	ELSE     = "else"
	PUB      = "pub"

	IDENTIFIER = "IDENTIFIER"

	TRUE  = "true"
	FALSE = "false"

	BANG = "!"
)

var Keywords = map[Type]string{
	RETURN:   "return",
	FUNCTION: "function",
	LET:      "let",
	IF:       "if",
	TRUE:     "true",
	FALSE:    "false",
	PUB:      "pub",
}
