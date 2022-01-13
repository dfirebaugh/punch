package compiler

import (
	"punch/lexer"
	"punch/parser"
)

type Compiler struct {
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
	lexer.New(p, source)
}

func (c Compiler) getSourceCode() string {
	return c.Source
}
func (c Compiler) reportSyntaxError() {}
func (c Compiler) optimize()          {}
func (c Compiler) generateCode()      {}
func (c Compiler) createGenerator()   {}
