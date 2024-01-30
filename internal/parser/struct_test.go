package parser

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/token"
)

func TestParseStructLiteral(t *testing.T) {
	input := `struct message {
		i32 		sender
		i8 			recipient
		string 	body
	}
	message msg = message{
		sender: 5,
		recipient: 10,
		body: "hello, world!"
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for _, s := range program.Statements {
		println(s.String())
	}

	for _, e := range p.Errors() {
		t.Errorf("parser has errors: %s", e)
	}
}

func TestParseStructLiteralWithoutFieldNames(t *testing.T) {
	input := `struct message {
		i32 		sender
		i8 			recipient
		string 	body
	}
	message msg = message{ 5, 10, "hello, world!" }
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	for _, s := range program.Statements {
		println(s.String())
	}

	for _, e := range p.Errors() {
		t.Errorf("parser has errors: %s", e)
	}
}

func TestParseStructDefinition(t *testing.T) {
	input := `
	struct message {
		i32 sender
		i8 recipient
		string body
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Expecting only one statement, the struct declaration
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	structDecl, ok := program.Statements[0].(*ast.StructDefinition)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.StructDefinition. got=%T", program.Statements[0])
	}

	if structDecl.Name.Value != "message" {
		t.Errorf("structDecl.Name.Value not 'message'. got=%s", structDecl.Name.Value)
	}

	expectedFields := []struct {
		typ  token.Type
		name string
	}{
		{token.I32, "sender"},
		{token.I8, "recipient"},
		{token.STRING, "body"},
	}

	if len(structDecl.Fields) != len(expectedFields) {
		t.Errorf("structDecl.Fields does not contain %d fields. got=%d", len(expectedFields), len(structDecl.Fields))
	}

	for i, field := range structDecl.Fields {
		if field.Type != expectedFields[i].typ {
			t.Errorf("structDecl.Fields[%d].Type.Value not '%s'. got=%s", i, expectedFields[i].typ, field.Type)
		}
		if field.Name.Value != expectedFields[i].name {
			t.Errorf("structDecl.Fields[%d].Name.Value not '%s'. got=%s", i, expectedFields[i].name, field.Name.Value)
		}
	}
}
