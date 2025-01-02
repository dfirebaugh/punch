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

type ReturnStatement struct {
	Token        token.Token
	ReturnValues []Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	for i, expr := range rs.ReturnValues {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(expr.String())
	}

	out.WriteString(";")
	return out.String()
}

type FunctionStatement struct {
	IsExported bool
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
	ReturnType *Identifier
}

func (f *FunctionStatement) expressionNode() {}
func (f *FunctionStatement) statementNode()  {}

func (f *FunctionStatement) TokenLiteral() string {
	if f.ReturnType != nil {
		return f.ReturnType.TokenLiteral()
	}
	return ""
}

func (f *FunctionStatement) String() string {
	if f == nil {
		return ""
	}

	var params []string
	if f.Parameters != nil {
		params = make([]string, len(f.Parameters))
		for i, p := range f.Parameters {
			if p != nil {
				params[i] = p.String()
			} else {
				params[i] = "nil"
			}
		}
	}

	var out bytes.Buffer
	if f.ReturnType != nil {
		out.WriteString(f.ReturnType.String() + " ")
	}
	if f.Name != nil {
		out.WriteString(f.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if f.Body != nil {
		out.WriteString(f.Body.String())
	}
	return out.String()
}

type FunctionDeclaration struct {
	ReturnType *Identifier
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
}

func (*FunctionDeclaration) statementNode() {}

func (fd *FunctionDeclaration) TokenLiteral() string {
	if fd.ReturnType != nil {
		return fd.ReturnType.TokenLiteral()
	}
	return ""
}

func (fd *FunctionDeclaration) String() string {
	var out strings.Builder

	if fd.ReturnType != nil {
		out.WriteString(fd.ReturnType.String() + " ")
	}
	out.WriteString(fd.Name.String())

	out.WriteString("(")
	for i, p := range fd.Parameters {
		out.WriteString(p.String())
		if i < len(fd.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")

	if fd.Body != nil {
		out.WriteString(fd.Body.String())
	}

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
