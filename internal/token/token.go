package token

import (
	"text/scanner"
)

// Type represents the type of a token.
type Type string

// Token struct holds the information about a token.
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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	DOT       = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "if"

	IDENTIFIER = "IDENTIFIER"
)

// Keywords is a map of keyword token types to their language representation.
var Keywords = map[Type]string{
	RETURN:   "return",
	FUNCTION: "function",
	LET:      "let",
	IF:       "if",
}
