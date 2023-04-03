package token

import (
	"fmt"
	"strconv"
	"unicode"
)

// String formats a token as a string for pretty printing.
func (t Token) String() string {
	return fmt.Sprintf("%d:%d:%d: | type: %s | literal: %q", t.Position.Line, t.Position.Column, t.Position.Offset, t.Type, t.Literal)
}

func (t Token) IsKeyword() bool {
	_, ok := Keywords[t.Type]
	return ok
}

func (t Token) IsOperator() bool {
	switch t.Type {
	case ASSIGN, PLUS, MINUS, ASTERISK, SLASH, EQ, NOT_EQ, LT, GT:
		return true
	default:
		return false
	}
}

func (t Token) IsNumber() bool {
	return t.Type == INT || t.Type == FLOAT
}

func (t Token) IsString() bool {
	return len(t.Literal) != 0 && t.Literal[0] == '"' && t.Literal[len(t.Literal)-1] == '"'
}

func (t Token) IsIdentifier() bool {
	if t.IsString() || t.IsInt() || t.IsFloatInt() {
		return false
	}
	if len(t.Literal) == 1 && unicode.IsLetter([]rune(t.Literal)[0]) {
		return true
	}
	for i, k := range t.Literal {
		if !t.IsIdentRune(k, i) {
			return false
		}
	}
	return true
}

func (t Token) IsFloatInt() bool {
	if _, err := strconv.ParseFloat(t.Literal, 64); err != nil {
		return false
	}

	return true
}

func (t Token) IsInt() bool {
	if _, err := strconv.ParseInt(t.Literal, 10, 64); err != nil {
		return false
	}

	return true
}

func (t Token) IsIdentRune(ch rune, i int) bool {
	// Allow underscores at any position
	if ch == '_' {
		return true
	}

	// Allow alphabetic characters or digits after the first character
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}

func (t Token) IsSingleCharIdentifier() bool {
	return len(t.Literal) == 1 && unicode.IsLetter([]rune(t.Literal)[0])
}

func Match(t1, t2 Token) bool {
	if t1.Type != t2.Type {
		fmt.Printf("Token type doesn't match: expected=%s got=%s\n", t1.Type, t2.Type)
		return false
	}
	if t1.Literal != t2.Literal {
		fmt.Printf("Token literal doesn't match: expected=%s got=%s\n", t1.Literal, t2.Literal)
		return false
	}
	if t1.Position.Line != t2.Position.Line || t1.Position.Column != t2.Position.Column || t1.Position.Offset != t2.Position.Offset {
		fmt.Printf("Token position doesn't match: expected=%d:%d:%d got=%d:%d:%d\n", t1.Position.Line, t1.Position.Column, t1.Position.Offset, t2.Position.Line, t2.Position.Column, t2.Position.Offset)
		return false
	}
	return true
}
