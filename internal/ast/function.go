package ast

import (
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
