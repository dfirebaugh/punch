package ast_test

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestLetStatement(t *testing.T) {
	input := "i8 x = 5;"
	p, program := parse(input, t)
	checkParserErrors(t, p)

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Files[0].Statements))
	}

	stmt := program.Files[0].Statements[0]
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

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Files[0].Statements))
	}

	stmt := program.Files[0].Statements[0]
	if stmt.TokenLiteral() != "return" {
		t.Fatalf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
	}

	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
	}

	if returnStmt.ReturnValues[0].String() != "5" {
		t.Fatalf("returnStmt.ReturnValue.String() not '5'. got=%s", returnStmt.ReturnValues[0].String())
	}
}

func TestFunctionStatement(t *testing.T) {
	input := "i8 add(i8 x, i8 y) { return x + y }"
	p, program := parse(input, t)
	checkParserErrors(t, p)

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program does not contain 1 statements. got=%d\n", len(program.Files[0].Statements))
	}

	stmt := program.Files[0].Statements[0]
	if stmt.TokenLiteral() != "i8" {
		t.Fatalf("stmt.TokenLiteral not 'i8'. got=%q", stmt.TokenLiteral())
	}

	funcStmt, ok := stmt.(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionStatement. got=%T", stmt)
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

func TestIfStatement(t *testing.T) {
	cond := &ast.Boolean{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true}
	cons := &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "foo"}, Value: "foo"}}}}
	alt := &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "bar"}, Value: "bar"}}}}
	ie := &ast.IfStatement{Token: token.Token{Type: token.IF, Literal: "if"}, Condition: cond, Consequence: cons, Alternative: alt}
	assert.Equal(t, "if true { foo } else { bar }", ie.String())
	assert.Equal(t, "if", ie.TokenLiteral())
}

func parse(input string, t *testing.T) (*parser.Parser, *ast.Program) {
	l := lexer.New("", input)
	p := parser.New(l)
	program := p.ParseProgram("")

	return p, program
}
