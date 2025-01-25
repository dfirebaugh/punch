package ast

import (
	"bytes"
	"strings"

	"github.com/dfirebaugh/punch/token"
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

type StructDefinition struct {
	Token  token.Token
	Name   *Identifier
	Fields []*StructField
}

func (sd *StructDefinition) statementNode() {}

func (sd *StructDefinition) TokenLiteral() string {
	return sd.Token.Literal
}

func (sd *StructDefinition) String() string {
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

type StructLiteral struct {
	Token      token.Token
	Fields     map[string]Expression
	StructName *Identifier
}

func (sl *StructLiteral) expressionNode() {}

func (sl *StructLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StructLiteral) String() string {
	var out bytes.Buffer
	if sl.StructName != nil {
		out.WriteString(sl.StructName.String())
	}
	out.WriteString("{ ")
	fields := []string{}
	for name, expr := range sl.Fields {
		fields = append(fields, name+": "+expr.String())
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString(" }")
	return out.String()
}

type StructFieldAccess struct {
	Token token.Token // The '.' token
	Left  Expression
	Field *Identifier
}

func (s *StructFieldAccess) expressionNode() {}

func (s *StructFieldAccess) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StructFieldAccess) String() string {
	return s.Left.String() + "." + s.Field.String()
}

type StructFieldAssignment struct {
	Token token.Token
	Left  *StructFieldAccess
	Right Expression
}

func (sfa *StructFieldAssignment) expressionNode() {}

func (sfa *StructFieldAssignment) TokenLiteral() string {
	return sfa.Token.Literal
}

func (sfa *StructFieldAssignment) String() string {
	var out bytes.Buffer
	out.WriteString(sfa.Left.String())
	out.WriteString(" = ")
	out.WriteString(sfa.Right.String())
	return out.String()
}
