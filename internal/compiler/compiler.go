package compiler

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/wat"
)

func Compile(source string) (string, []byte) {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()
	wat := wat.GenerateWAT(program)
	wasm, err := wasmtime.Wat2Wasm(wat)
	if err != nil {
		panic(err)
	}
	return wat, wasm
}
