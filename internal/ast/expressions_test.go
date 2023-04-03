package ast

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	ident := &Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "foo"}, Value: "foo"}
	assert.Equal(t, "foo", ident.String())
	assert.Equal(t, "foo", ident.TokenLiteral())
}

func TestBinaryExpression(t *testing.T) {
	left := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5}
	op := token.Token{Type: token.PLUS, Literal: "+"}
	right := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10}
	be := &BinaryExpression{Left: left, Operator: op, Right: right}
	assert.Equal(t, "(5 + 10)", be.String())
	assert.Equal(t, "+", be.TokenLiteral())
}

func TestIfExpression(t *testing.T) {
	cond := &Boolean{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true}
	cons := &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: &Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "foo"}, Value: "foo"}}}}
	alt := &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: &Identifier{Token: token.Token{Type: token.IDENTIFIER, Literal: "bar"}, Value: "bar"}}}}
	ie := &IfExpression{Token: token.Token{Type: token.IF, Literal: "if"}, Condition: cond, Consequence: cons, Alternative: alt}
	assert.Equal(t, "if true { foo } else { bar }", ie.String())
	assert.Equal(t, "if", ie.TokenLiteral())
}
