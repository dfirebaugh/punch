package lexer

import "testing"

type testCollector struct{}

func (testCollector) Collect(token Token) {}

func TestEvaluateToken(t *testing.T) {
	l := New(testCollector{}, ``)

	tokens := []struct {
		ExpectedType    TokenType
		ExpectedLiteral string
	}{
		{
			ExpectedType:    PLUS,
			ExpectedLiteral: "+",
		},
		{
			ExpectedType:    ASSIGN,
			ExpectedLiteral: "=",
		},
		{
			ExpectedType:    COMMA,
			ExpectedLiteral: ",",
		},
		{
			ExpectedType:    SEMICOLON,
			ExpectedLiteral: ";",
		},
		{
			ExpectedType:    LPAREN,
			ExpectedLiteral: "(",
		},
		{
			ExpectedType:    RPAREN,
			ExpectedLiteral: ")",
		},
		{
			ExpectedType:    LBRACE,
			ExpectedLiteral: "{",
		},
		{
			ExpectedType:    RBRACE,
			ExpectedLiteral: "}",
		},
		{
			ExpectedType:    FUNCTION,
			ExpectedLiteral: "function",
		},
		{
			ExpectedType:    LET,
			ExpectedLiteral: "let",
		},
		{
			ExpectedType:    PRINT,
			ExpectedLiteral: "print",
		},
		{
			ExpectedType:    RETURN,
			ExpectedLiteral: "return",
		},
		{
			ExpectedType:    STRING,
			ExpectedLiteral: `"some string - could be any kind of string of any length"`,
		},
	}

	for _, v := range tokens {
		received := l.evaluateTokenType(Token{Literal: v.ExpectedLiteral})
		if received != v.ExpectedType {
			t.Errorf("could not evaluate token expected: %s | received: %s", v.ExpectedType, received)
		}
	}
}
