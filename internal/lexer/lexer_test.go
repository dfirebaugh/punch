package lexer

import (
	"punch/internal/token"
	"testing"
	"text/scanner"
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

	expectedTokens := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "42"},
		{Type: token.FLOAT, Literal: "3.14"},
		{Type: token.STRING, Literal: "\"hello world\""},
		{Type: token.IDENTIFIER, Literal: "foo"},
		{Type: token.IF, Literal: "if"},
	}

	for i, expectedToken := range expectedTokens {
		t.Run(expectedToken.Literal, func(t *testing.T) {
			token := l.evaluateType(expectedToken)
			if token != expectedToken.Type {
				t.Errorf("test %d: expected %+v but got %+v", i+1, expectedToken, token)
			}
		})
	}
}
