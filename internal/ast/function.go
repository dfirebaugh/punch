package ast

import (
	"bytes"
	"strings"

	"github.com/dfirebaugh/punch/internal/token"
)

type Parameter struct {
	Identifier *Identifier
	Type       token.Type
}

func (p *Parameter) String() string {
	if p == nil {
		return ""
	}
	return p.Identifier.String()
}

type FunctionDeclaration struct {
	ReturnType token.Token
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) TokenLiteral() string {
	return fd.ReturnType.Literal
}
func (fd *FunctionDeclaration) statementNode() {}

func (fd *FunctionDeclaration) String() string {
	var out strings.Builder

	out.WriteString(fd.ReturnType.Literal + " ")

	out.WriteString(fd.Name.String())

	out.WriteString("(")
	for i, p := range fd.Parameters {
		out.WriteString(p.String())
		if i < len(fd.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(fd.Body.String())

	return out.String()
}

type FunctionCall struct {
	FunctionName string
	Token        token.Token
	Function     Expression
	Arguments    []Expression
}

func (f *FunctionCall) expressionNode() {}

func (f *FunctionCall) TokenLiteral() string {
	if f == nil {
		return ""
	}
	return "function call"
}

func (f *FunctionCall) String() string {
	if f == nil {
		return ""
	}
	args := make([]string, len(f.Arguments))
	for i, a := range f.Arguments {
		args[i] = a.String()
	}
	var out bytes.Buffer
	out.WriteString(f.FunctionName)
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
