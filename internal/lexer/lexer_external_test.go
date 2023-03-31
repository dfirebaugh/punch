package lexer_test

import (
	"punch/internal/lexer"
	"punch/internal/token"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		output []token.Token
	}{
		{
			input: "let xx = 10;",
			output: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENTIFIER, Literal: "xx"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		tokens := l.Run()

		for i, tok := range tokens {
			expected := tt.output[i]
			if tok.Type != expected.Type || tok.Literal != expected.Literal {
				t.Errorf("unexpected token[%d]: got=%v, want=%v", i, tok.Type, expected.Type)
			}
		}
	}
}
