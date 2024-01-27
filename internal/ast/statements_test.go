package ast_test

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
)

func TestLetStatement(t *testing.T) {
	input := "i8 x = 5;"
	p, program := parse(input, t)
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Statements))
	}

	stmt := program.Statements[0]
	if stmt.TokenLiteral() != "i8" {
		t.Fatalf("stmt.TokenLiteral not 'i8'. got=%q", stmt.TokenLiteral())
	}

	letStmt, ok := stmt.(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("stmt not *ast.VariableDeclaration. got=%T", stmt)
	}

	if letStmt.Name.Value != "x" {
		t.Fatalf("letStmt.Name.Value not 'x'. got=%s", letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != "x" {
		t.Fatalf("letStmt.Name.TokenLiteral() not 'x'. got=%s", letStmt.Name.TokenLiteral())
	}
}

func TestReturnStatement(t *testing.T) {
	input := "return 5;"
	p, program := parse(input, t)
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Statements))
	}

	stmt := program.Statements[0]
	if stmt.TokenLiteral() != "return" {
		t.Fatalf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
	}

	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
	}

	if returnStmt.ReturnValue.String() != "5" {
		t.Fatalf("returnStmt.ReturnValue.String() not '5'. got=%s", returnStmt.ReturnValue.String())
	}
}

func TestFunctionStatement(t *testing.T) {
	input := "i8 add(x i8, y i8) { return x + y; }"
	p, program := parse(input, t)
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Statements))
	}

	stmt := program.Statements[0]
	if stmt.TokenLiteral() != "i8" {
		t.Fatalf("stmt.TokenLiteral not 'i8'. got=%q", stmt.TokenLiteral())
	}

	funcStmt, ok := stmt.(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionDeclaration. got=%T", stmt)
	}

	if funcStmt.Name.Value != "add" {
		t.Fatalf("funcStmt.Name.Value not 'add'. got=%s", funcStmt.Name.Value)
	}

	if len(funcStmt.Parameters) != 2 {
		t.Fatalf("len(funcStmt.Parameters) not 2. got=%d", len(funcStmt.Parameters))
	}

	if funcStmt.Parameters[0].Identifier.Value != "x" {
		t.Fatalf("funcStmt.Parameters[0].Value not 'x'. got=%s", funcStmt.Parameters[0].Identifier.Value)
	}

	if funcStmt.Parameters[1].Identifier.Value != "y" {
		t.Fatalf("funcStmt.Parameters[1].Value not 'y'. got=%s", funcStmt.Parameters[1].Identifier.Value)
	}

	if funcStmt.Body.String() != "{ return (x + y); }" {
		t.Fatalf("funcStmt.Body.String() not '{ return (x + y); }'. got=%s", funcStmt.Body.String())
	}
}

func parse(input string, t *testing.T) (*parser.Parser, *ast.Program) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return p, program
}
