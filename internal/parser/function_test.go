package parser

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
)

func TestParseMultipleReturnTypes(t *testing.T) {
	input := `pub (i8, string) myFunction(i32 x, bool y) { return 5, "hello" }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Statements[0])
	}

	if len(stmt.ReturnTypes) != 2 {
		t.Fatalf("function statement does not have 2 return types. got=%d", len(stmt.ReturnTypes))
	}
}

func TestParseFunctionWithControlFlow(t *testing.T) {
	input := `pub i8 myFunction(i32 x, bool y) { if true { return 0; } return 5 }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Statements[0])
	}

	if len(stmt.ReturnTypes) != 1 {
		t.Fatalf("function statement does not have 2 return types. got=%d", len(stmt.ReturnTypes))
	}
}
