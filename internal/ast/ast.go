package ast

import (
	"fmt"
	"punch/internal/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) ExpressionNode()      {}
func (i Identifier) String() string       { return i.Value }
func (i Identifier) TokenLiteral() string { return i.Token.Literal }

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l LetStatement) StatementNode()       {}
func (l LetStatement) String() string       { return l.Value.String() }
func (l LetStatement) TokenLiteral() string { return l.Name.Token.Literal }

type ReturnStatement struct {
	Token token.Token
	Name  Identifier
	Value Expression
}

func (l ReturnStatement) StatementNode()       {}
func (l ReturnStatement) String() string       { return l.Value.String() }
func (l ReturnStatement) TokenLiteral() string { return l.Name.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value string
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression  // the expression itself
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
