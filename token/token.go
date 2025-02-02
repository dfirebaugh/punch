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
	STRING = "STRING"
	NUMBER = "NUMBER"
	FLOAT  = "FLOAT"
	BOOL   = "BOOL"

	// Operators
	ASSIGN          = "="
	INFER           = ":="
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
	SLASH_SLASH     = "//"
	SLASH_ASTERISK  = "/*"
	ASTERISK_SLASH  = "/*"

	// Delimiters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FN        = "FN"
	FUNCTION  = "FUNCTION"
	CONST     = "CONST"
	LET       = "LET"
	RETURN    = "RETURN"
	FOR       = "FOR"
	IF        = "IF"
	ELSE      = "ELSE"
	PUB       = "PUB"
	PACKAGE   = "PKG"
	IMPORT    = "IMPORT"
	INTERFACE = "INTERFACE"
	STRUCT    = "STRUCT"
	TEST      = "TEST"
	ENUM      = "ENUM"
	DEFER     = "DEFER"
	APPEND    = "APPEND"
	LEN       = "LEN"

	IDENTIFIER = "IDENTIFIER"

	TRUE  = "true"
	FALSE = "false"

	BANG = "!"
)

const (
	U8  = "u8"
	U16 = "u16"
	U32 = "u32"
	U64 = "u64"
	I8  = "i8"
	I16 = "i16"
	I32 = "i32"
	I64 = "i64"
	F32 = "f32"
	F64 = "f64"
)

var Keywords = map[Type]string{
	PACKAGE:   "pkg",
	IMPORT:    "import",
	INTERFACE: "interface",
	STRUCT:    "struct",
	TEST:      "test",
	FN:        "fn",
	ENUM:      "enum",
	RETURN:    "return",
	CONST:     "const",
	LET:       "let",
	IF:        "if",
	ELSE:      "else",
	FOR:       "for",
	TRUE:      "true",
	FALSE:     "false",
	PUB:       "pub",
	BOOL:      "bool",
	STRING:    "str",
	APPEND:    "append",
	LEN:       "len",
}
