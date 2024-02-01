package parser

import (
	"fmt"
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
)

func TestParseIdentifier(t *testing.T) {
	input := "foobar;"
	lexer := lexer.New("", input)
	parser := New(lexer)

	program := parser.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression is not ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Value)
	}
}

func TestParseBoolean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		lexer := lexer.New("", tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("expression is not ast.Boolean. got=%T", stmt.Expression)
		}

		if boolean.Value != tt.want {
			t.Errorf("for input %q, boolean.Value is not %t. got=%t", tt.input, tt.want, boolean.Value)
		}
	}
}

func TestParseAssignment(t *testing.T) {
	input := "i8 x = 5;"
	lexer := lexer.New("", input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	expectedStatements := []string{
		"i8 x = 5;",
	}

	if len(expectedStatements) != len(program.Statements) {
		t.Errorf("program has wrong number of statements. got=%d expected=%d", len(program.Statements), len(expectedStatements))
	}

	for i, s := range program.Statements {
		if s.String() != expectedStatements[i] {
			t.Errorf("expected %q, got %q", expectedStatements[i], s.String())
		}
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := "return 5;"
	lexer := lexer.New("", input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStatement. got=%T", program.Statements[0])
	}

	if stmt.TokenLiteral() != "return" {
		t.Errorf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
	}

	if _, ok := stmt.ReturnValues[0].(*ast.IntegerLiteral); !ok {
		t.Errorf("stmt.Value not *ast.IntegerLiteral. got=%T", stmt.ReturnValues)
	}

	if stmt.ReturnValues[0].TokenLiteral() != "5" {
		t.Errorf("stmt.Value.TokenLiteral() not '5'. got=%q", stmt.ReturnValues[0].TokenLiteral())
	}
}

func TestParsePrefixExpression(t *testing.T) {
	input := "!5;"
	lexer := lexer.New("", input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
	}

	if exp.Operator.Literal != "!" {
		t.Errorf("exp.Operator is not '!'. got=%q", exp.Operator)
	}

	if _, ok := exp.Right.(*ast.IntegerLiteral); !ok {
		t.Errorf("exp.Right is not *ast.IntegerLiteral. got=%T", exp.Right)
	}

	if exp.Right.TokenLiteral() != "5" {
		t.Errorf("exp.Right.TokenLiteral() is not '5'. got=%q", exp.Right.TokenLiteral())
	}
}

func TestParseInfixExpression(t *testing.T) {
	tests := []struct {
		input string
		left  int64
		op    string
		right int64
	}{
		{"5 + 6;", 5, "+", 6},
		{"5 - 6;", 5, "-", 6},
		{"5 * 6;", 5, "*", 6},
		{"5 / 6;", 5, "/", 6},
		{"5 > 6;", 5, ">", 6},
		{"5 < 6;", 5, "<", 6},
		{"5 == 6;", 5, "==", 6},
		{"5 != 6;", 5, "!=", 6},
	}

	for _, tt := range tests {
		lexer := lexer.New("", tt.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		println(program.JSONPretty())

		_, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.left) {
			return
		}

		if exp.Operator.Literal != tt.op {
			t.Errorf("exp.Operator is not '%s'. got=%q", tt.op, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.right) {
			return
		}
	}
}

func TestParseInfixLT(t *testing.T) {
	input := "5 < 6;"
	l := lexer.New("", input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	println(program.JSONPretty())
	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}
}

func TestParseInfixIdentifiersLT(t *testing.T) {
	input := "z < x"
	l := lexer.New("", input)
	parser := New(l)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	println(program.JSONPretty())
	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
		return false
	}
}

func TestParseComment(t *testing.T) {
	input := `i8 commentFn() {
		i8 x = 5
		/*
		This is a multi-line comment
		that spans multiple lines.
		*/
		i8 y = 10
		return y + x
	}`
	expected := `i8 commentFn() { i8 x = 5; i8 y = 10; return (y + x); }`

	l := lexer.New("", input)

	p := New(l)
	stmt := p.parseFunctionStatement()
	if stmt == nil {
		t.Errorf("statement was nil")
		return
	}
	fmt.Printf("Parsed Statement: %s\n", stmt.String())

	if got := stmt.String(); got != expected {
		t.Errorf("expected comments to be stripped out:\nwant %q\ngot  %q", expected, got)
	}
}

func TestParseTypeBasedVariableDeclaration(t *testing.T) {
	input := `i32 myVar = 5;`
	l := lexer.New("", input)
	p := New(l)

	decl := p.parseTypeBasedVariableDeclaration()
	if decl == nil {
		t.Fatalf("parseTypeBasedVariableDeclaration() returned nil")
	}

	varDecl, ok := decl.(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("Statement is not *ast.VariableDeclaration. got=%T", decl)
	}

	if varDecl.Type.Literal != "i32" {
		t.Errorf("varDecl.Type.Literal not 'i32'. got=%s", varDecl.Type.Literal)
	}

	if varDecl.Name.Value != "myVar" {
		t.Errorf("varDecl.Name.Value not 'myVar'. got=%s", varDecl.Name.Value)
	}

	val, ok := varDecl.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("varDecl.Value not *ast.IntegerLiteral. got=%T", varDecl.Value)
	}

	if val.Value != 5 {
		t.Errorf("varDecl.Value.Value not 5. got=%d", val.Value)
	}
}

func TestParseExpressionAssignment(t *testing.T) {
	input := "i8 myVar = 5"
	l := lexer.New("", input)
	p := New(l)
	p.nextToken() // Initialize the parser to point to the first token.

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		t.Fatalf("parseExpression() returned nil")
	}

	assignExpr, ok := expr.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("Expression is not *ast.AssignmentExpression. got=%T", expr)
	}

	if !testIdentifier(t, assignExpr.Left, "myVar") {
		return
	}

	val, ok := assignExpr.Right.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("assignExpr.Right is not *ast.IntegerLiteral. got=%T", assignExpr.Right)
	}

	if val.Value != 5 {
		t.Errorf("assignExpr.Right.Value is not 5. got=%d", val.Value)
	}
}

func TestStringLiteralReturn(t *testing.T) {
	input := `string main() {
		return "hello, world!";
	}`

	l := lexer.New("", input)
	p := New(l)
	p.ParseProgram()
	for _, e := range p.Errors() {
		t.Errorf("parser has errors: %s", e)
	}
}

func TestStructDeclaration(t *testing.T) {
	input := `struct message {
		i32 sender
		i8 recipient
		string body
	}`
	l := lexer.New("", input)
	p := New(l)
	program := p.ParseProgram()

	for _, s := range program.Statements {
		println(s.String())
	}

	for _, e := range p.Errors() {
		t.Errorf("parser has errors: %s", e)
	}
}

func TestIfCondition(t *testing.T) {
	input := `
i8 main() {
	if 1 == 1 {
		return 1
	}
	return 0
}`

	l := lexer.New("", input)
	p := New(l)
	program := p.ParseProgram()
	println(program.JSONPretty())
	for _, s := range program.Statements {
		fn, ok := s.(*ast.FunctionDeclaration)
		if !ok {
			break
		}
		for i, bs := range fn.Body.Statements {
			// println(bs.TokenLiteral())
			println(bs.String())
			println(i, fn.Body.Statements[i].String())
		}
		if s == nil {
			continue
		}
		println(s.String())
		// println(s.TokenLiteral())
	}

	for _, e := range p.Errors() {
		t.Errorf("parser has errors: %s", e)
	}
}

func TestParseBlockStatement(t *testing.T) {
	input := `{
			i8 x = 5
			return x
	}`

	l := lexer.New("", input)
	p := New(l)
	p.nextToken()

	blockStmt := p.parseBlockStatement()
	if blockStmt == nil {
		t.Fatalf("parseBlockStatement() returned nil")
	}
	println(blockStmt.Statements[0].String())
	println(blockStmt.Statements[1].String())
	if blockStmt.Statements[0].String() != "i8 x = 5;" {
		t.Errorf("blockStmt.Statements[0].String not 'i8 x = 5;'. got=%s", blockStmt.Statements[0].String())
	}
	_, ok := blockStmt.Statements[1].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("blockStmt.Statements[1] is not ast.ReturnStatement. got=%T", blockStmt.Statements[1])
	}
	if blockStmt.Statements[1].String() != "return x;" {
		t.Errorf("blockStmt.Statements[1].String not'return x;'. got=%s", blockStmt.Statements[1].String())
	}
}

func TestParseBoolReturn(t *testing.T) {
	input := `
bool isEq(i8 a, i8 b) {
		return a == b;
}
	`
	l := lexer.New("", input)
	p := New(l)

	program := p.ParseProgram()
	println(program.JSONPretty())
	checkParserErrors(t, p)
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	lit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if lit.Value != value {
		t.Errorf("lit.Value not %d. got=%d", value, lit.Value)
		return false
	}

	if lit.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("lit.TokenLiteral not %d. got=%q", value, lit.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator.Literal != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator.Literal)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
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
