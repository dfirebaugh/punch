package token

import "text/scanner"

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
	STRING = "string"
	NUMBER = "number"
	BOOL   = "bool"

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
	FUNCTION  = "FUNCTION"
	CONST     = "CONST"
	LET       = "LET"
	RETURN    = "RETURN"
	IF        = "if"
	ELSE      = "else"
	PUB       = "pub"
	PACKAGE   = "pkg"
	IMPORT    = "import"
	INTERFACE = "interface"
	STRUCT    = "struct"
	TEST      = "test"
	ENUM      = "enum"
	DEFER     = "defer"

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
	F8  = "f8"
	F16 = "f16"
	F32 = "f32"
	F64 = "f64"
)

var Keywords = map[Type]string{
	PACKAGE:   "pkg",
	IMPORT:    "import",
	INTERFACE: "interface",
	STRUCT:    "struct",
	TEST:      "test",
	ENUM:      "enum",
	RETURN:    "return",
	CONST:     "const",
	LET:       "let",
	IF:        "if",
	ELSE:      "else",
	TRUE:      "true",
	FALSE:     "false",
	PUB:       "pub",
}
