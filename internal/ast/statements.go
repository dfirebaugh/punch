package ast

import (
	"bytes"
	"strings"

	"github.com/dfirebaugh/punch/internal/token"
)

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return string(ls.Token.Literal) }

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Token.Literal + " ")

	if ls.Name != nil {
		out.WriteString(ls.Name.TokenLiteral() + " " + token.ASSIGN + " ")
	}

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(token.SEMICOLON)
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(token.SEMICOLON)
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}
func (bs *BlockStatement) statementNode()  {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

type FunctionStatement struct {
	IsExported bool
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionStatement) expressionNode() {}
func (f *FunctionStatement) statementNode()  {}

func (f *FunctionStatement) TokenLiteral() string {
	return "function"
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
	out.WriteString("function ")
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
