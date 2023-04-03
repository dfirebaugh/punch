package token

import (
	"testing"
	"text/scanner"
)

func TestTokenString(t *testing.T) {
	tok := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	want := "1:2:3: | type: IDENTIFIER | literal: \"myVar\""
	got := tok.String()
	if got != want {
		t.Errorf("Token.String(): got %q, want %q", got, want)
	}
}

func TestTokenIsKeyword(t *testing.T) {
	kwTok := Token{Type: IF, Literal: "if"}
	idTok := Token{Type: IDENTIFIER, Literal: "myVar"}
	if !kwTok.IsKeyword() {
		t.Errorf("Token.IsKeyword(): expected true for keyword token, got false")
	}
	if idTok.IsKeyword() {
		t.Errorf("Token.IsKeyword(): expected false for non-keyword token, got true")
	}
}

func TestTokenIsString(t *testing.T) {
	strTok := Token{Type: STRING, Literal: "\"hello world\""}
	notStrTok := Token{Type: IDENTIFIER, Literal: "myVar"}
	if !strTok.IsString() {
		t.Errorf("Token.IsString(): expected true for string token, got false")
	}
	if notStrTok.IsString() {
		t.Errorf("Token.IsString(): expected false for non-string token, got true")
	}
}

func TestTokenIsIdentifier(t *testing.T) {
	identTok := Token{Type: IDENTIFIER, Literal: "myVar"}
	numTok := Token{Type: INT, Literal: "1234"}
	strTok := Token{Type: STRING, Literal: "\"hello\""}
	symTok := Token{Type: PLUS, Literal: "+"}
	if !identTok.IsIdentifier() {
		t.Errorf("Token.IsIdentifier(): expected true for identifier token, got false")
	}
	if numTok.IsIdentifier() {
		t.Errorf("Token.IsIdentifier(): expected false for integer token, got true")
	}
	if strTok.IsIdentifier() {
		t.Errorf("Token.IsIdentifier(): expected false for string token, got true")
	}
	if symTok.IsIdentifier() {
		t.Errorf("Token.IsIdentifier(): expected false for symbol token, got true")
	}
}

func TestTokenIsFloatInt(t *testing.T) {
	floatTok := Token{Type: FLOAT, Literal: "3.14"}
	intTok := Token{Type: INT, Literal: "42"}
	notNumTok := Token{Type: IDENTIFIER, Literal: "myVar"}
	if !floatTok.IsFloatInt() {
		t.Errorf("Token.IsFloatInt(): expected true for floating point number token, got false")
	}
	if !intTok.IsFloatInt() {
		t.Errorf("Token.IsFloatInt(): expected true for integer token, got false")
	}
	if notNumTok.IsFloatInt() {
		t.Errorf("Token.IsFloatInt(): expected false for non-number token, got true")
	}
}

func TestTokenIsInt(t *testing.T) {
	intTok := Token{Type: INT, Literal: "42"}
	floatTok := Token{Type: FLOAT, Literal: "3.14"}
	notNumTok := Token{Type: IDENTIFIER, Literal: "myVar"}
	if !intTok.IsInt() {
		t.Errorf("Token.IsInt(): expected true for integer token, got false")
	}
	if floatTok.IsInt() {
		t.Errorf("Token.IsInt(): expected false for floating point number token, got true")
	}
	if notNumTok.IsInt() {
		t.Errorf("Token.IsInt(): expected false for non-number token, got true")
	}
}

func TestTokenIsIdentRune(t *testing.T) {
	identChar := 'a'
	numChar := '1'
	symChar := '+'
	underChar := '_'
	tok := Token{}
	if !tok.IsIdentRune(identChar, 0) {
		t.Errorf("IsIdentRune(): expected true for alphabetical character at position 0, got false")
	}
	if !tok.IsIdentRune(numChar, 1) {
		t.Errorf("IsIdentRune(): expected true for numeric character after position 0, got false")
	}
	if tok.IsIdentRune(symChar, 0) {
		t.Errorf("IsIdentRune(): expected false for symbol character, got true")
	}
	if !tok.IsIdentRune(underChar, 2) {
		t.Errorf("IsIdentRune(): expected true for underscore character at the middle of identifier, got false")
	}
}

func TestTokenIsSingleCharIdentifier(t *testing.T) {
	singleCharIdent := Token{Type: IDENTIFIER, Literal: "x"}
	multipleCharIdent := Token{Type: IDENTIFIER, Literal: "var"}
	numTok := Token{Type: INT, Literal: "1234"}
	if !singleCharIdent.IsSingleCharIdentifier() {
		t.Errorf("Token.IsSingleCharIdentifier(): expected true for single character identifier token, got false")
	}
	if multipleCharIdent.IsSingleCharIdentifier() {
		t.Errorf("Token.IsSingleCharIdentifier(): expected false for multi-character identifier token, got true")
	}
	if numTok.IsSingleCharIdentifier() {
		t.Errorf("Token.IsSingleCharIdentifier(): expected false for integer token, got true")
	}
}

func TestMatch(t *testing.T) {
	tok1 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	tok2 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	if !Match(tok1, tok2) {
		t.Errorf("Match(): expected true for equal tokens, got false")
	}
	tok3 := Token{
		Type:     IDENTIFIER,
		Literal:  "differentVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	tok4 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	if Match(tok3, tok4) {
		t.Errorf("Match(): expected false for tokens with different literals, got true")
	}

	tok5 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	tok6 := Token{
		Type:     INT,
		Literal:  "1234",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	if Match(tok5, tok6) {
		t.Errorf("Match(): expected false for tokens with different types, got true")
	}

	tok7 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 1, Column: 2, Offset: 3},
	}
	tok8 := Token{
		Type:     IDENTIFIER,
		Literal:  "myVar",
		Position: scanner.Position{Line: 2, Column: 3, Offset: 4},
	}
	if Match(tok7, tok8) {
		t.Errorf("Match(): expected false for tokens with different positions, got true")
	}
}

func TestIsOperator(t *testing.T) {
	tests := []struct {
		token Token
		want  bool
	}{
		{Token{Type: PLUS}, true},
		{Token{Type: MINUS}, true},
		{Token{Type: ASTERISK}, true},
		{Token{Type: SLASH}, true},
		{Token{Type: EQ}, true},
		{Token{Type: NOT_EQ}, true},
		{Token{Type: LT}, true},
		{Token{Type: GT}, true},
		{Token{Type: ASSIGN}, true},
		{Token{Type: INT, Literal: "42"}, false},
		{Token{Type: FLOAT, Literal: "3.14"}, false},
		{Token{Type: IDENTIFIER, Literal: "foo"}, false},
	}

	for _, tc := range tests {
		got := tc.token.IsOperator()
		if got != tc.want {
			t.Errorf("IsOperator(%v) = %v, want %v", tc.token, got, tc.want)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		token Token
		want  bool
	}{
		{Token{Type: INT, Literal: "42"}, true},
		{Token{Type: FLOAT, Literal: "3.14"}, true},
		{Token{Type: PLUS}, false},
		{Token{Type: MINUS}, false},
		{Token{Type: IDENTIFIER, Literal: "foo"}, false},
		{Token{Type: ASSIGN}, false},
	}

	for _, tc := range tests {
		got := tc.token.IsNumber()
		if got != tc.want {
			t.Errorf("IsNumber(%v) = %v, want %v", tc.token, got, tc.want)
		}
	}
}
