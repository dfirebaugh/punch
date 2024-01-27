package wat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dfirebaugh/punch/internal/ast"
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
		assert.Equal(t, expected, wat.GenerateWAT(ast))
	})

	t.Run("generate WAT code for a program with one return statement", func(t *testing.T) {
		ast := &ast.Program{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					ReturnValue: &ast.Boolean{
						Value: true,
					},
				},
			},
		}
		expected := `(module
(return (i32.const 1))
)
`
		assert.Equal(t, expected, wat.GenerateWAT(ast))
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
							ReturnValue: &ast.Identifier{Value: "x"},
						},
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
		assert.Equal(t, expected, wat.GenerateWAT(ast))
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
								ReturnValue: &ast.InfixExpression{
									Left:     &ast.Identifier{Value: "x"},
									Right:    &ast.Identifier{Value: "y"},
									Operator: token.Token{Type: token.PLUS},
								},
							},
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
		assert.Equal(t, expected, wat.GenerateWAT(ast))
	})
}
