package ast

// Node is the interface that all AST nodes must implement.
type Node interface {
	TokenLiteral() string // Returns the literal value of the node's token
	String() string       // Returns a string representation of the node for debugging purposes
}

// Statement represents a statement in the code.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the code.
type Expression interface {
	Node
	expressionNode()
}

// Program represents a complete program in the code.
type Program struct {
	Statements []Statement // The statements in the program
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
		out += "\n"
	}
	return out
}
