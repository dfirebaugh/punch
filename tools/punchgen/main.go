package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	js_gen "github.com/dfirebaugh/punch/emitters/js"
	"github.com/dfirebaugh/punch/emitters/wat"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/parser"
	"github.com/dfirebaugh/punch/token"
)

func parse(this js.Value, p []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in parse:", r)
		}
	}()
	source := p[0].String()
	l := lexer.New("example", source)
	parser := parser.New(l)
	program := parser.ParseProgram("ast_explorer")

	astJSON, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to generate AST JSON: %v", err),
		}
	}

	return string(astJSON)
}

func lex(this js.Value, p []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in lex:", r)
		}
	}()
	source := p[0].String()
	l := lexer.New("example", source)
	var tokens []string
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokens = append(tokens, fmt.Sprintf("%s: %q", tok.Type, tok.Literal))
	}

	tokensJSON, err := json.Marshal(tokens)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to generate tokens JSON: %v", err),
		}
	}

	return string(tokensJSON)
}

func generateWAT(this js.Value, p []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in generateWAT:", r)
		}
	}()
	source := p[0].String()
	l := lexer.New("example", source)
	parser := parser.New(l)
	program := parser.ParseProgram("ast_explorer")

	watCode := wat.GenerateWAT(program, true)
	return watCode
}

func generateJS(this js.Value, p []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in generateJS:", r)
		}
	}()
	source := p[0].String()
	l := lexer.New("example", source)
	parser := parser.New(l)
	program := parser.ParseProgram("ast_explorer")

	t := js_gen.NewTranspiler()
	jsCode, err := t.Transpile(program)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to transpile to JS: %v", err),
		}
	}

	return jsCode
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in main:", r)
		}
	}()
	js.Global().Set("parse", js.FuncOf(parse))
	js.Global().Set("lex", js.FuncOf(lex))
	js.Global().Set("generateWAT", js.FuncOf(generateWAT))
	js.Global().Set("generateJS", js.FuncOf(generateJS))

	select {}
}
