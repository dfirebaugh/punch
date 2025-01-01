package parser

import (
	"fmt"
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/token"
)

func TestParseMultipleReturnTypes(t *testing.T) {
	input := `pub (i8, string) myFunction(i32 x, bool y) { return 5, "hello" }`
	l := lexer.New("", input)
	p := New(l)

	program := p.ParseProgram("")
	checkParserErrors(t, p)

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Files[0].Statements))
	}

	stmt, ok := program.Files[0].Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Files[0].Statements[0])
	}

	if len(stmt.ReturnTypes) != 2 {
		t.Fatalf("function statement does not have 2 return types. got=%d", len(stmt.ReturnTypes))
	}
}

func TestParseFunctionWithControlFlow(t *testing.T) {
	input := `pub i8 myFunction(i32 x, bool y) { if true { return 0; } return 5 }`
	l := lexer.New("", input)
	p := New(l)

	program := p.ParseProgram("")
	checkParserErrors(t, p)

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Files[0].Statements))
	}

	stmt, ok := program.Files[0].Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Files[0].Statements[0])
	}

	if len(stmt.ReturnTypes) != 1 {
		t.Fatalf("function statement does not have 2 return types. got=%d", len(stmt.ReturnTypes))
	}
}

func TestParseFunctionStatement(t *testing.T) {
	input := "i8 add(i8 x, i8 y) { return x + y }"
	l := lexer.New("", input)
	p := New(l)

	program := p.ParseProgram("")
	checkParserErrors(t, p)

	if len(program.Files[0].Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Files[0].Statements))
	}

	stmt, ok := program.Files[0].Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T", program.Files[0].Statements[0])
	}

	if len(stmt.ReturnTypes) != 1 {
		t.Fatalf("function statement does not have 2 return types. got=%d", len(stmt.ReturnTypes))
	}
}

func TestParseFunctionParametersLexing(t *testing.T) {
	input := "pub i8 addTwo(x i8, y string) {return x + y; }"
	l := lexer.New("", input)
	p := New(l)

	expectedTokens := []token.Token{
		{Type: token.PUB, Literal: "pub"},
		{Type: token.I8, Literal: "i8"},
		{Type: token.IDENTIFIER, Literal: "addTwo"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.I8, Literal: "i8"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.STRING, Literal: "string"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.EOF, Literal: ""},
	}

	for _, expected := range expectedTokens {
		fmt.Printf("Token: Type=%s, Literal=%s\n", p.curToken.Type, p.curToken.Literal)
		if p.curToken.Type != expected.Type || p.curToken.Literal != expected.Literal {
			t.Fatalf("unexpected token - expected %q (%q), got %q (%q)", expected.Type, expected.Literal, p.curToken.Type, p.curToken.Literal)
		}
		p.nextToken()
	}
}

func TestParseFunctionParameter(t *testing.T) {
	input := "i8 x"

	l := lexer.New("", input)
	p := New(l)

	param := p.parseFunctionParameter()

	if param == nil {
		t.Fatalf("Expected a valid parameter, got nil")
	}

	expectedName := "x"
	expectedType := token.I8

	if param.Identifier.Value != expectedName {
		t.Errorf("Expected parameter name to be %q, got %q", expectedName, param.Identifier.Value)
	}

	if string(param.Type) != expectedType {
		t.Errorf("Expected parameter type to be %s, got %s", expectedType, param.Type)
	}
}

func TestParseFunctionParameters(t *testing.T) {
	input := "i8 addTwo(i8 x, string y) {return x + y;}"
	l := lexer.New("", input)
	p := New(l)
	p.nextToken()
	p.nextToken()

	parameters := p.parseFunctionParameters()

	expected := []*ast.Parameter{
		{Identifier: &ast.Identifier{Value: "x"}, Type: token.I8},
		{Identifier: &ast.Identifier{Value: "y"}, Type: token.STRING},
	}

	if len(parameters) != len(expected) {
		t.Fatalf("wrong number of parameters. want=%d, got=%d\n", len(expected), len(parameters))
	}

	for i, param := range parameters {
		if param.Identifier.Value != expected[i].Identifier.Value {
			t.Errorf("parameter %d identifier wrong. want=%s, got=%s", i, expected[i].Identifier.Value, param.Identifier.Value)
		}
		if param.Type != expected[i].Type {
			t.Errorf("parameter %d type wrong. want=%s, got=%s", i, expected[i].Type, param.Type)
		}
	}
}

func TestParseExportedFunction(t *testing.T) {
	input := "pub i8 addTwo(i8 x, string y) {return x + y;}"
	l := lexer.New("", input)
	p := New(l)
	prog := p.ParseProgram("")

	f, ok := prog.Files[0].Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("expected *ast.FunctionStatement, got %T", prog.Files[0].Statements[0].(*ast.FunctionStatement))
	}

	if f.Name.Value != "addTwo" {
		t.Errorf("function name wrong. want=%s, got=%s", "addTwo", f.Name)
	}
	if f.IsExported == false {
		t.Fatalf("expected function to be exported, got %t", f.IsExported)
	}
	if f.Parameters[0].Identifier.Token.Literal != "x" {
		t.Fatalf("expected function parameter x, got %s", f.Parameters[0].Identifier.Token.String())
	}
	if f.Parameters[0].Type != "i8" {
		t.Fatalf("expected function parameter type to be i8 got %s", f.Parameters[0].Type)
	}
	if f.Parameters[1].Identifier.Token.Literal != "y" {
		t.Fatalf("expected function parameter y, got %s", f.Parameters[1].Identifier.Token.String())
	}
	if f.Parameters[1].Type != "string" {
		t.Fatalf("expected function parameter type to be string got %s", f.Parameters[1].Type)
	}
	ret, ok := f.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected *ast.ReturnStatement, got %T", f.Body.Statements[0].(*ast.ReturnStatement))
	}
	if ret.ReturnValues[0].String() != "(x + y)" {
		t.Errorf("return value wrong. want=%s, got=%s", "x + y", ret.ReturnValues[0].String())
	}
}

func TestParseFunction(t *testing.T) {
	input := `
pkg main
i8 addTwo(i8 x, string y) {
	if true {
		return 0;
	}
	return x + y;
}`
	l := lexer.New("", input)
	p := New(l)
	prog := p.ParseProgram("")

	f, ok := prog.Files[0].Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("expected *ast.FunctionStatement, got %T", prog.Files[0].Statements[0].(*ast.FunctionStatement))
	}

	if f.Name.Value != "addTwo" {
		t.Errorf("function name wrong. want=%s, got=%s", "addTwo", f.Name)
	}
	if f.IsExported != false {
		t.Fatalf("expected function to be exported, got %t", f.IsExported)
	}
	if f.Parameters[0].Identifier.Token.Literal != "x" {
		t.Fatalf("expected function parameter x, got %s", f.Parameters[0].Identifier.Token.String())
	}
	if f.Parameters[0].Type != "i8" {
		t.Fatalf("expected function parameter type to be i8 got %s", f.Parameters[0].Type)
	}
	if f.Parameters[1].Identifier.Token.Literal != "y" {
		t.Fatalf("expected function parameter y, got %s", f.Parameters[1].Identifier.Token.String())
	}
	if f.Parameters[1].Type != "string" {
		t.Fatalf("expected function parameter type to be string got %s", f.Parameters[1].Type)
	}
	ifBlock, ok := f.Body.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected *ast.ReturnStatement, got %T", f.Body.Statements[0].(*ast.ReturnStatement))
	}
	if ifBlock.Condition.String() != "true" {
		t.Errorf("if condition wrong. want=%s, got=%s", "true", ifBlock.Condition.String())
	}
	ret, ok := f.Body.Statements[1].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected *ast.ReturnStatement, got %T", f.Body.Statements[1].(*ast.ReturnStatement))
	}
	if ret.ReturnValues[0].String() != "(x + y)" {
		t.Errorf("return value wrong. want=%s, got=%s", "x + y", ret.ReturnValues[0].String())
	}
}
