package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/dfirebaugh/punch/token"
)

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string {
	if il != nil {
		return strconv.FormatInt(il.Value, 10)
	}
	return ""
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (il *FloatLiteral) expressionNode() {}

func (il *FloatLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *FloatLiteral) String() string {
	if il != nil {
		return strconv.FormatFloat(il.Value, 10, 64, 64)
	}
	return ""
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string {
	if sl != nil {
		return sl.Value
	}
	return ""
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) String() string {
	if bl != nil {
		return bl.Token.Literal
	}
	return ""
}

type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string {
	return "[ARRAYLITERAL]: " + a.String()
}

func (a *ArrayLiteral) String() string {
	if a == nil {
		return ""
	}
	elts := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elts[i] = e.String()
	}
	var out bytes.Buffer
	out.WriteString("[")
	out.WriteString(strings.Join(elts, ", "))
	out.WriteString("]")
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	if fl == nil || fl.Token.Type == token.EOF {
		return ""
	}
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	if fl == nil {
		return ""
	}
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string {
	if hl == nil || hl.Token.Type == token.EOF {
		return ""
	}
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	if hl == nil {
		return ""
	}
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type NumberType struct {
	Token token.Token
	Type  token.Type
}

func (nt *NumberType) expressionNode() {}

func (nt *NumberType) TokenLiteral() string {
	return nt.Token.Literal
}

func (nt *NumberType) String() string {
	return nt.Token.Literal
}
