package parser

import (
	"fmt"
	"testing"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/token"
)

func TestParseIdentifier(t *testing.T) {
	input := "foobar;"
	lexer := lexer.New(input)
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
		lexer := lexer.New(tt.input)
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

func TestParseLetStatement(t *testing.T) {
	input := "let x = 5;"
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
	}

	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%q", stmt.Name.Value)
	}

	if stmt.Name.TokenLiteral() != "x" {
		t.Errorf("stmt.Name.TokenLiteral() not 'x'. got=%q", stmt.Name.TokenLiteral())
	}

	if _, ok := stmt.Value.(*ast.IntegerLiteral); !ok {
		t.Errorf("stmt.Value not *ast.IntegerLiteral. got=%T", stmt.Value)
	}

	if stmt.Value.TokenLiteral() != "5" {
		t.Errorf("stmt.Value.TokenLiteral() not '5'. got=%q", stmt.Value.TokenLiteral())
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := "return 5;"
	lexer := lexer.New(input)
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

	if _, ok := stmt.ReturnValue.(*ast.IntegerLiteral); !ok {
		t.Errorf("stmt.Value not *ast.IntegerLiteral. got=%T", stmt.ReturnValue)
	}

	if stmt.ReturnValue.TokenLiteral() != "5" {
		t.Errorf("stmt.Value.TokenLiteral() not '5'. got=%q", stmt.ReturnValue.TokenLiteral())
	}
}

func TestParsePrefixExpression(t *testing.T) {
	input := "!5;"
	lexer := lexer.New(input)
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

	if exp.Operator != "!" {
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
		lexer := lexer.New(tt.input)
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

func TestParseGroupedExpression(t *testing.T) {
	input := "(5 + 6);"
	lexer := lexer.New(input)
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

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if !testLiteralExpression(t, exp.Left, 5) {
		return
	}

	if exp.Operator.Literal != "+" {
		t.Errorf("exp.Operator is not '+'. got=%q", exp.Operator)
	}

	if !testLiteralExpression(t, exp.Right, 6) {
		return
	}
}

func TestParseIfExpression(t *testing.T) {
	input := "if (x < y) { x };"
	lexer := lexer.New(input)
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

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
		return
	}

	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative was not nil. got=%+v", exp.Alternative)
	}
}

func TestParseBlockStatement(t *testing.T) {
	input := `{
		let x = 5;
		let y = 10;
		let result = x + y;
	}`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	if len(parser.Errors()) != 0 {
		for _, err := range parser.Errors() {
			t.Errorf("parser error: %s", err)
		}
		t.Fatalf("parser has %d errors", len(parser.Errors()))
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong number of statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("expected *ast.BlockStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Statements) != 3 {
		t.Fatalf("expected 3 statements in block statement. got=%d", len(stmt.Statements))
	}

	expectedTokens := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.LET, Literal: "let"},
		{Type: token.LET, Literal: "let"},
	}
	for i, s := range stmt.Statements {
		if s.TokenLiteral() != string(expectedTokens[i].Literal) {
			t.Errorf("expected token literal=%q, got=%q", expectedTokens[i], s.TokenLiteral())
		}
	}
}

func testBlockStatement(t *testing.T, stmt ast.Statement, expectedLen int) bool {
	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Errorf("stmt is not ast.BlockStatement. got=%T", stmt)
		return false
	}

	if len(blockStmt.Statements) != expectedLen {
		t.Errorf("block does not contain %d statements. got=%d\n", expectedLen, len(blockStmt.Statements))
		return false
	}

	return true
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

func TestParseFunctionStatement(t *testing.T) {
	input := "function add(x, y) { return x + y; }"

	l := lexer.New(input)
	p := New(l)

	stmt := p.parseFunctionStatement()

	if stmt == nil {
		t.Fatalf("Expected function statement, got nil")
	}

	if stmt.Name.Value != "add" {
		t.Errorf("Expected function name to be 'add', got '%s'", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(stmt.Parameters))
	}

	expectedParams := []string{"x", "y"}

	for i, ident := range stmt.Parameters {
		if ident.Value != expectedParams[i] {
			t.Errorf("Parameter %d should have value %q, got %q", i+1, expectedParams[i], ident.Value)
		}
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("Expected function body to have 1 statement, got %d", len(stmt.Body.Statements))
	}

	blockStmt, ok := stmt.Body.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("Expected function body to contain a block statement, got %T", blockStmt.Statements[0])
	}
	returnStmt, ok := blockStmt.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected function body to contain a return statement, got %T", stmt.Body.Statements[0])
	}

	if !testInfixExpression(t, returnStmt.ReturnValue, "x", "+", "y") {
		t.Errorf("Unexpected return value: %v", returnStmt.ReturnValue)
	}
}

func TestParseFunctionParameters(t *testing.T) {
	input := "(foo, bar, baz)"

	l := lexer.New(input)
	p := New(l)

	params := p.parseFunctionParameters()

	if len(params) != 3 {
		t.Fatalf("Expected 3 parameters, got %d", len(params))
	}

	expectedParams := []string{"foo", "bar", "baz"}

	for i, ident := range params {
		if ident.Value != expectedParams[i] {
			t.Errorf("Parameter %d should have value %q, got %q", i+1, expectedParams[i], ident.Value)
		}
	}
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
