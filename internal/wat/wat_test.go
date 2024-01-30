package wat_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/token"
	"github.com/dfirebaugh/punch/internal/wat"
)

func _TestGenerateWAT(t *testing.T) {
	t.Run("generate WAT code for a program with one let statement", func(t *testing.T) {
		ast := &ast.Program{
			Statements: []ast.Statement{
				&ast.LetStatement{
					Name: &ast.Identifier{Value: "x"},
					Value: &ast.IntegerLiteral{
						Value: 10,
					},
				},
			},
		}
		expected := `(module
(local $x i32)
(i32.const 10)
(local.set $x x)
)
`
		assert.Equal(t, expected, wat.GenerateWAT(ast, false))
	})

	t.Run("generate WAT code for a program with one return statement", func(t *testing.T) {
		ast := &ast.Program{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					ReturnValues: []ast.Expression{&ast.Boolean{
						Value: true,
					},
					},
				},
			},
		}
		expected := `(module
(return (i32.const 1))
)
`
		assert.Equal(t, expected, wat.GenerateWAT(ast, false))
	})

	t.Run("generate WAT code for a program with one block statement", func(t *testing.T) {
		ast := &ast.Program{
			Statements: []ast.Statement{
				&ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.LetStatement{
							Name: &ast.Identifier{Value: "x"},
							Value: &ast.IntegerLiteral{
								Value: 10,
							},
						},
						&ast.ReturnStatement{
							ReturnValues: []ast.Expression{&ast.Identifier{Value: "x"}}},
					},
				},
			},
		}
		expected := `(module

(local $x i32)
(i32.const 10)
(local.set $x x)
(return (local.get $x))
)
`
		assert.Equal(t, expected, wat.GenerateWAT(ast, false))
	})

	t.Run("generate WAT code for a program with one function declaration", func(t *testing.T) {
		ast := &ast.Program{
			Statements: []ast.Statement{
				&ast.FunctionStatement{
					Name: &ast.Identifier{Value: "add"},
					Parameters: []*ast.Parameter{
						{Identifier: &ast.Identifier{Value: "x"}},
						{Identifier: &ast.Identifier{Value: "y"}},
					},
					Body: &ast.BlockStatement{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								ReturnValues: []ast.Expression{&ast.InfixExpression{
									Left:     &ast.Identifier{Value: "x"},
									Right:    &ast.Identifier{Value: "y"},
									Operator: token.Token{Type: token.PLUS},
								},
								}},
						},
					},
				},
			},
		}
		expected := `(module
(func $add (param $x i32)
(param $y i32)
(result i32)

(return (i32.add (local.get $x) (local.get $y)))

)
)
`
		assert.Equal(t, expected, wat.GenerateWAT(ast, false))
	})
}

func TestFunctionDeclaration(t *testing.T) {
	input := `
bool is_eq(i32 a, i32 b) {
	return (a == b)
}

pub i32 add_two(i32 x, i32 y) {
	return x + y
}

pub i32 add_four(i32 a, i32 b, i32 c, i32 d) {
	if !is_eq(a, c) {
		return (a - b - c - d)
	}
	return a + b + c + d
}
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	expected := "(module\n\t(func $is_eq (param $a i32) (param $b i32) (result i32)\n\t\t\t(return (i32.eq (local.get $a) (local.get $b)))\n\n)\n\t(func $add_two (export \"add_two\") (param $x i32) (param $y i32) (result i32)\n\t\t\t(return (i32.add (local.get $x) (local.get $y)))\n\n)\n\t(func $add_four (export \"add_four\") (param $a i32) (param $b i32) (param $c i32) (param $d i32) (result i32)\n\t\t\t(if (i32.eqz \n\t\t(call $is_eq (local.get $a) (local.get $c)))\n\t\t\t(then\n\t\t\t(return (i32.sub (i32.sub (i32.sub (local.get $a) (local.get $b)) (local.get $c)) (local.get $d)))\n\n\n\t\t\t)\n\t\t)\n\t\t\t(return (i32.add (i32.add (i32.add (local.get $a) (local.get $b)) (local.get $c)) (local.get $d)))\n\n)\n)"
	w := wat.GenerateWAT(program, false)
	t.Run("generate WAT code for a program with one function declaration", func(t *testing.T) {
		assert.Equal(
			t,
			strings.Trim(strings.Trim(expected, "\n"), "\t"),
			strings.Trim(strings.Trim(w, "\n"), "\t"))
	})
}
