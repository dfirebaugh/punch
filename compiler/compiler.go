package compiler

import (
	"punch/lexer"
	"punch/parser"
)

type Compiler struct {
	Lexer  lexer.Lexer
	Source string
}

func New(source string) Compiler {
	return Compiler{
		Source: source,
	}
}

func (c Compiler) Run() {
	c.compile(c.getSourceCode())
}
func (c Compiler) extractCommandLineArguments() {}

func (c Compiler) compile(source string) {
	syntaxBuilder := parser.NewSyntaxBuilder()
	p := parser.New(syntaxBuilder)
	l := lexer.New(p)
	l.Lex(source)

	p.HandleEvent(parser.EOF, -1, -1)

}

func (c Compiler) getSourceCode() string {
	return c.Source
}
func (c Compiler) reportSyntaxError() {}
func (c Compiler) optimize()          {}
func (c Compiler) generateCode()      {}
func (c Compiler) createGenerator()   {}
