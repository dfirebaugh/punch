package compiler

import (
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/wat"
)

func Compile(source string) string {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	return wat.GenerateWAT(program)
}
