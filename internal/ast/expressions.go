package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/internal/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	if i == nil || i.Token.Type == token.EOF {
		return ""
	}
	return i.Token.Literal
}

func (i *Identifier) String() string {
	if i == nil {
		return ""
	}
	return i.Value
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}

func (be *BinaryExpression) TokenLiteral() string {
	if be == nil || be.Operator.Type == token.EOF {
		return ""
	}
	return be.Operator.Literal
}

func (be *BinaryExpression) String() string {
	if be == nil {
		return ""
	}
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator.Literal, be.Right.String())
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	if ie == nil || ie.Token.Type == token.EOF {
		return ""
	}
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	if ie == nil {
		return ""
	}
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	out.WriteString(" ")

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type WhileExpression struct {
	Condition Expression
	Body      Statement
	ID        int
}

func (w *WhileExpression) expressionNode() {}

func (w *WhileExpression) TokenLiteral() string {
	if w == nil {
		return ""
	}
	return "while"
}

func (w *WhileExpression) String() string {
	if w == nil {
		return ""
	}
	var out bytes.Buffer
	out.WriteString("while (")
	out.WriteString(w.Condition.String())
	out.WriteString(") {\n")
	out.WriteString(w.Body.String())
	out.WriteString("}\n")
	return out.String()
}

type CallExpression struct {
	Function  *Identifier
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) TokenLiteral() string {
	if c == nil || c.Function == nil {
		return ""
	}
	return "call"
}

func (c *CallExpression) String() string {
	if c == nil || c.Function == nil {
		return ""
	}
	args := make([]string, len(c.Arguments))
	for i, a := range c.Arguments {
		args[i] = a.String()
	}
	var out bytes.Buffer
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}

func (i *IndexExpression) TokenLiteral() string {
	if i == nil {
		return ""
	}
	return "index"
}

func (i *IndexExpression) String() string {
	if i == nil {
		return ""
	}
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator token.Token
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	if pe == nil || pe.Token.Type == token.EOF {
		return ""
	}
	return pe.TokenLiteral()
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(" " + pe.Operator.Literal + " ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string { return ie.Operator.Literal }

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator.Literal + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token // the token.Boolean token, either true or false
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string { return b.Token.Literal }

type Integer struct {
	Token token.Token
	Value int64
}

func (i *Integer) expressionNode() {}

func (i *Integer) TokenLiteral() string {
	if i == nil || i.Token.Type == token.EOF {
		return ""
	}
	return i.Token.Literal
}

func (i *Integer) String() string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", i.Value)
}

type AssignmentExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ae *AssignmentExpression) expressionNode() {}

func (ae *AssignmentExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Left.String())
	out.WriteString(" ")
	out.WriteString(ae.Token.Literal)
	out.WriteString(" ")
	out.WriteString(ae.Right.String())
	return out.String()
}
