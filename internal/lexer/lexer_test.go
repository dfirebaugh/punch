package lexer

import (
	"reflect"
	"testing"
	"text/scanner"

	"github.com/dfirebaugh/punch/internal/token"
)

type testCollector struct{}

func (testCollector) Collect(token token.Token) {}

func TestEvaluateToken(t *testing.T) {
	l := New(``)

	tokens := []struct {
		ExpectedType    token.Type
		ExpectedLiteral string
	}{
		{
			ExpectedType:    token.PLUS,
			ExpectedLiteral: "+",
		},
		{
			ExpectedType:    token.ASSIGN,
			ExpectedLiteral: "=",
		},
		{
			ExpectedType:    token.COMMA,
			ExpectedLiteral: ",",
		},
		{
			ExpectedType:    token.SEMICOLON,
			ExpectedLiteral: ";",
		},
		{
			ExpectedType:    token.LPAREN,
			ExpectedLiteral: "(",
		},
		{
			ExpectedType:    token.RPAREN,
			ExpectedLiteral: ")",
		},
		{
			ExpectedType:    token.LBRACE,
			ExpectedLiteral: "{",
		},
		{
			ExpectedType:    token.RBRACE,
			ExpectedLiteral: "}",
		},
		{
			ExpectedType:    token.FUNCTION,
			ExpectedLiteral: "function",
		},
		{
			ExpectedType:    token.LET,
			ExpectedLiteral: "let",
		},
		{
			ExpectedType:    token.RETURN,
			ExpectedLiteral: "return",
		},
		{
			ExpectedType:    token.STRING,
			ExpectedLiteral: `"some string - could be any kind of string of any length"`,
		},
		{
			ExpectedType:    token.INT,
			ExpectedLiteral: `42`,
		},
		{
			ExpectedType:    token.INT,
			ExpectedLiteral: `1`,
		},
		{
			ExpectedType:    token.FLOAT,
			ExpectedLiteral: `42.2`,
		},
	}

	for _, v := range tokens {
		received := l.evaluateType(token.Token{Literal: v.ExpectedLiteral})
		if received != v.ExpectedType {
			t.Errorf("could not evaluate token expected: %s | received: %s", v.ExpectedType, received)
		}
	}
}

func TestLexLetStatement(t *testing.T) {
	input := "let x = 5;"
	expectedTokens := []token.Token{
		{Type: token.LET, Literal: "let", Position: scanner.Position{Line: 1, Column: 1, Offset: 0}},
		{Type: token.IDENTIFIER, Literal: "x", Position: scanner.Position{Line: 1, Column: 5, Offset: 4}},
		{Type: token.ASSIGN, Literal: "=", Position: scanner.Position{Line: 1, Column: 7, Offset: 6}},
		{Type: token.INT, Literal: "5", Position: scanner.Position{Line: 1, Column: 9, Offset: 8}},
		{Type: token.SEMICOLON, Literal: ";", Position: scanner.Position{Line: 1, Column: 10, Offset: 9}},
	}

	l := New(input)
	for i, expectedToken := range expectedTokens {
		tok := l.NextToken()

		if !token.Match(expectedToken, tok) {
			t.Errorf("Test[%d] - Token wrong. expected=%d:%d:%d | type: %s | literal: %q |, got=%d:%d:%d | type: %s | literal: %q |",
				i, expectedToken.Position.Line, expectedToken.Position.Column, expectedToken.Position.Offset, expectedToken.Type, expectedToken.Literal,
				tok.Position.Line, tok.Position.Column, tok.Position.Offset, tok.Type, tok.Literal)
		}
	}
}
func Test_evaluateType(t *testing.T) {
	input := `= 42 3.14 "hello world" foo if`
	l := New(input)
	expectedTokens := []token.Type{
		token.ASSIGN,
		token.INT,
		token.FLOAT,
		token.STRING,
		token.IDENTIFIER,
		token.IF,
	}
	for i, expectedTokenType := range expectedTokens {
		t.Run(string(expectedTokenType), func(t *testing.T) {
			token := l.evaluateType(l.NextToken())
			if token != expectedTokenType {
				t.Errorf("test %d: expected %+v but got %+v", i+1, expectedTokenType, token)
			}
		})
	}
}

func Test_isMultiCharOperator(t *testing.T) {
	operators := []struct {
		input string
		want  bool
	}{
		{">=", true},
		{"<=", true},
		{"==", true},
		{"!=", true},
		{"+=", true},
		{"-=", true},
		{"*=", true},
		{"/=", true},
		{"|", false},
		{"&", false},
		{"+", false},
		{"-", false},
		{"*", false},
		{"/", false},
		{">", false},
		{"<", false},
		{"=", false},
		{"!", false},
	}

	for _, tt := range operators {
		got := Lexer{}.isMultiCharOperator(token.Type(tt.input))
		if got != tt.want {
			t.Errorf("isMultiCharOperator(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestBooleanLiterals(t *testing.T) {
	source := "let a = true; let b = false;"
	expectedTokens := []token.Token{
		{Type: token.LET, Literal: "let", Position: scanner.Position{Line: 1, Column: 1, Offset: 0}},
		{Type: token.IDENTIFIER, Literal: "a", Position: scanner.Position{Line: 1, Column: 5, Offset: 4}},
		{Type: token.ASSIGN, Literal: "=", Position: scanner.Position{Line: 1, Column: 7, Offset: 6}},
		{Type: token.TRUE, Literal: "true", Position: scanner.Position{Line: 1, Column: 9, Offset: 8}},
		{Type: token.SEMICOLON, Literal: ";", Position: scanner.Position{Line: 1, Column: 13, Offset: 12}},
		{Type: token.LET, Literal: "let", Position: scanner.Position{Line: 1, Column: 15, Offset: 14}},
		{Type: token.IDENTIFIER, Literal: "b", Position: scanner.Position{Line: 1, Column: 19, Offset: 18}},
		{Type: token.ASSIGN, Literal: "=", Position: scanner.Position{Line: 1, Column: 21, Offset: 20}},
		{Type: token.FALSE, Literal: "false", Position: scanner.Position{Line: 1, Column: 23, Offset: 22}},
		{Type: token.SEMICOLON, Literal: ";", Position: scanner.Position{Line: 1, Column: 28, Offset: 27}},
		{Type: token.EOF},
	}

	l := New(source)
	tokens := l.Run()

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Expected tokens %v, but got %v", expectedTokens, tokens)
	}
}

func TestEvaluateKeyword(t *testing.T) {
	l := New("")

	tests := []struct {
		input    string
		expected token.Type
	}{
		{"function", token.FUNCTION},
		{"let", token.LET},
		{"return", token.RETURN},
		{"if", token.IF},
		{"true", token.TRUE},
		{"false", token.FALSE},
		{"foo", token.IDENTIFIER},
	}

	for _, tt := range tests {
		actual := l.evaluateKeyword(tt.input)
		if actual != tt.expected {
			t.Errorf("For input '%s', expected %v but got %v", tt.input, tt.expected, actual)
		}
	}
}
