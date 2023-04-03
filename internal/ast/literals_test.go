package ast

import (
	"testing"

	"github.com/dfirebaugh/punch/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestNumberLiteral(t *testing.T) {
	nl := &NumberLiteral{Token: token.Token{Type: token.INT, Literal: "123"}, Value: "123"}
	assert.Equal(t, "123", nl.String())
	assert.Equal(t, "123", nl.TokenLiteral())
}

func TestIntegerLiteral(t *testing.T) {
	il := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "123"}, Value: 123}
	assert.Equal(t, "123", il.String())
	assert.Equal(t, "123", il.TokenLiteral())
}

func TestStringLiteral(t *testing.T) {
	sl := &StringLiteral{Token: token.Token{Type: token.STRING, Literal: "hello"}, Value: "hello world"}
	assert.Equal(t, "hello world", sl.String())
	assert.Equal(t, "hello", sl.TokenLiteral())
}

func TestBooleanLiteral(t *testing.T) {
	bl := &BooleanLiteral{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true}
	assert.Equal(t, "true", bl.String())
	assert.Equal(t, "true", bl.TokenLiteral())
}

func TestArrayLiteral(t *testing.T) {
	elts := []Expression{
		&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
		&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
		&IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "3"}, Value: 3},
	}
	al := &ArrayLiteral{Elements: elts}
	assert.Equal(t, "[1, 2, 3]", al.String())
	assert.Equal(t, "[ARRAYLITERAL]: [1, 2, 3]", al.TokenLiteral())
}
