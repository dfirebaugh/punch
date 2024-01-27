package ast

import (
	"bytes"

	"github.com/dfirebaugh/punch/internal/token"
)

type StructField struct {
	Token token.Token
	Name  *Identifier
	Type  token.Type
}

func (sf *StructField) statementNode() {}

func (sf *StructField) TokenLiteral() string {
	return sf.Token.Literal
}

func (sf *StructField) String() string {
	var out bytes.Buffer
	if sf.Name != nil {
		out.WriteString(sf.Name.String())
	}
	out.WriteString(" ")
	out.WriteString(string(sf.Type))
	return out.String()
}

type StructDeclaration struct {
	Token  token.Token
	Name   *Identifier
	Fields []*StructField
}

func (sd *StructDeclaration) statementNode() {}

func (sd *StructDeclaration) TokenLiteral() string {
	return sd.Token.Literal
}

func (sd *StructDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(sd.TokenLiteral() + " ")
	if sd.Name != nil {
		out.WriteString(sd.Name.String())
	}
	out.WriteString(" {")
	for _, field := range sd.Fields {
		out.WriteString("\n  ")
		out.WriteString(field.String())
	}
	if len(sd.Fields) > 0 {
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}
