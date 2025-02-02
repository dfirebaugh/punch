package compiler

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/dfirebaugh/punch/emitters/wat"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/parser"
	"github.com/sirupsen/logrus"
)

const (
	wasmDisabled = false
	astDisabled  = false
)

func Compile(filename string, source string) (string, []byte, string) {
	l := lexer.New(filename, source)
	p := parser.New(l)
	program, err := p.ParseProgram(filename)
	if err != nil {
		logrus.Error(err)
		return "", nil, ""
	}
	var ast string
	wat := wat.GenerateWAT(program, true)
	if !astDisabled {
		ast, _ = program.JSONPretty()
	}
	if wasmDisabled {
		return wat, nil, ast
	}
	wasm, err := wasmtime.Wat2Wasm(wat)
	if err != nil {
		panic(err)
	}

	return wat, wasm, ast
}
