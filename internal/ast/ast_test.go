package ast_test

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/token"
)

func TestLetStatementString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.NumberLiteral{
					Token: token.Token{Type: token.INT, Literal: "42"},
					Value: "42",
				},
			},
		},
	}

	expected := "let myVar = 42;\n"
	if program.String() != expected {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestReturnStatementString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &ast.NumberLiteral{
					Token: token.Token{Type: token.INT, Literal: "42"},
					Value: "42",
				},
			},
		},
	}

	expected := "return 42;\n"
	if program.String() != expected {
		t.Errorf("program.String() wrong. got=%q want=%q", program.String(), expected)
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + 5", "(5 + 5)"},
		{"5 - 5", "(5 - 5)"},
		{"5 * 5", "(5 * 5)"},
		{"5 / 5", "(5 / 5)"},
		{"5 > 5", "(5 > 5)"},
		{"5 < 5", "(5 < 5)"},
		{"5 == 5", "(5 == 5)"},
		{"5 != 5", "(5 != 5)"},
		{"true == true", "(true == true)"},
		{"true != false", "(true != false)"},
		{"false == false", "(false == false)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.String() != tt.expected {
			t.Errorf("exp.String() wrong. got=%q, expected=%q", exp.String(), tt.expected)
		}
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
