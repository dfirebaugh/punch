package ast

import (
	"bytes"
	"strings"

	"github.com/dfirebaugh/punch/token"
)

type ListDeclaration struct {
	Token token.Token  // The 'let' or 'const' token
	Type  token.Type   // The type of the list elements
	Name  *Identifier  // The name of the list
	Value *ListLiteral // The list literal
}

func (ld *ListDeclaration) statementNode() {}

func (ld *ListDeclaration) TokenLiteral() string {
	return ld.Token.Literal
}

func (ld *ListDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(ld.TokenLiteral() + " ")
	out.WriteString(string(ld.Type) + " ")
	out.WriteString(ld.Name.String() + " = ")
	out.WriteString(ld.Value.String())
	out.WriteString(";")
	return out.String()
}

type ListLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode() {}

func (ll *ListLiteral) TokenLiteral() string {
	return ll.Token.Literal
}

func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type ListOperation struct {
	Token    token.Token // The operation token (e.g., 'append', 'len')
	Operator string
	List     Expression
	Element  Expression // Used for operations like append
}

func (lo *ListOperation) expressionNode() {}

func (lo *ListOperation) TokenLiteral() string {
	return lo.Token.Literal
}

func (lo *ListOperation) String() string {
	var out bytes.Buffer
	out.WriteString(lo.Operator)
	out.WriteString("(")
	out.WriteString(lo.List.String())
	if lo.Element != nil {
		out.WriteString(", ")
		out.WriteString(lo.Element.String())
	}
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression  // The expression being indexed
	Index Expression  // The index expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
