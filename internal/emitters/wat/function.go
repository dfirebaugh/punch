package wat

import (
	"fmt"
	"log"
	"strings"

	"github.com/dfirebaugh/punch/internal/ast"
)

var (
	scopeStack       []map[string]string
	stringLiteralMap map[string]string
)

func pushScope() {
	scopeStack = append(scopeStack, make(map[string]string))
}

func popScope() {
	scopeStack = scopeStack[:len(scopeStack)-1]
}

func declareVariable(name, watType string) string {
	currentScope := scopeStack[len(scopeStack)-1]
	if _, exists := currentScope[name]; exists {
		log.Fatalf("Variable %s already declared in current scope", name)
	}
	declaration := fmt.Sprintf("(local $%s %s)\n", name, watType)
	currentScope[name] = declaration
	return declaration
}

func lookupVariable(name string) string {
	for i := len(scopeStack) - 1; i >= 0; i-- {
		if declaration, exists := scopeStack[i][name]; exists {
			return declaration
		}
	}
	log.Fatalf("Variable %s not found in any scope", name)
	return ""
}

func generateFunctionStatement(s *ast.FunctionStatement) string {
	var out strings.Builder

	out.WriteString(fmt.Sprintf("(func $%s ", s.Name.Value))
	if s.IsExported {
		out.WriteString(fmt.Sprintf("(export \"%s\") ", s.Name.Value))
	}

	pushScope()
	for _, param := range s.Parameters {
		declaration := fmt.Sprintf("(param $%s %s) ", param.Identifier.Value, mapTypeToWAT(string(param.Type)))
		scopeStack[len(scopeStack)-1][param.Identifier.Value] = declaration
		out.WriteString(declaration)
	}

	if s.ReturnType != nil {
		out.WriteString(fmt.Sprintf("(result %s) ", mapTypeToWAT(s.ReturnType.Value)))
	}
	out.WriteString("\n")

	stringLiteralMap = make(map[string]string)

	declaredLocals := make(map[string]bool)
	for _, param := range s.Parameters {
		declaredLocals[param.Identifier.Value] = true
	}
	var locals []string
	var initializations []string
	var stringLiterals []string
	for _, stmt := range s.Body.Statements {
		collectLocalsAndInitializations(stmt, declaredLocals, &locals, &initializations, &stringLiterals)
	}

	for _, local := range locals {
		out.WriteString(local)
	}

	for _, strInit := range stringLiterals {
		out.WriteString(strInit)
	}

	for _, init := range initializations {
		out.WriteString(init)
	}

	for _, stmt := range s.Body.Statements {
		out.WriteString(generateStatement(stmt))
	}

	popScope()
	out.WriteString(")\n")
	return out.String()
}

func collectLocalsAndInitializations(stmt ast.Statement, declaredLocals map[string]bool, locals *[]string, initializations *[]string, stringLiterals *[]string) {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		if !declaredLocals[s.Name.Value] {
			*locals = append(*locals, fmt.Sprintf("(local $%s %s)\n", s.Name.Value, mapTypeToWAT(s.Type.Literal)))
			declaredLocals[s.Name.Value] = true
		}
		if s.Value != nil {
			*initializations = append(*initializations, fmt.Sprintf("(local.set $%s %s)\n", s.Name.Value, generateExpression(s.Value)))
		}
	case *ast.BlockStatement:
		pushScope()
		for _, stmt := range s.Statements {
			collectLocalsAndInitializations(stmt, declaredLocals, locals, initializations, stringLiterals)
		}
		popScope()
	case *ast.IfStatement:
		collectLocalsAndInitializations(s.Consequence, declaredLocals, locals, initializations, stringLiterals)
		if s.Alternative != nil {
			collectLocalsAndInitializations(s.Alternative, declaredLocals, locals, initializations, stringLiterals)
		}
	case *ast.ExpressionStatement:
		collectExpressionLocals(s.Expression, declaredLocals, locals, initializations, stringLiterals)
	}
}

func collectExpressionLocals(
	expr ast.Expression,
	declaredLocals map[string]bool,
	locals *[]string,
	initializations *[]string,
	stringLiterals *[]string,
) {
	switch e := expr.(type) {
	case *ast.InfixExpression:
		collectExpressionLocals(e.Left, declaredLocals, locals, initializations, stringLiterals)
		collectExpressionLocals(e.Right, declaredLocals, locals, initializations, stringLiterals)
	case *ast.FunctionCall:
		for _, arg := range e.Arguments {
			collectExpressionLocals(arg, declaredLocals, locals, initializations, stringLiterals)
		}
	case *ast.Identifier:
		if !declaredLocals[e.Value] {
			log.Fatalf("Undeclared identifier: %s", e.Value)
		}
	case *ast.StringLiteral:
		if localVarName, exists := stringLiteralMap[e.Value]; exists {
			*initializations = append(*initializations, fmt.Sprintf("(local.get $%s)\n", localVarName))
		} else {
			length := len(e.Value) + 1
			localVarName := generateUniqueLocalVarName("str_ptr")
			stringLiteralMap[e.Value] = localVarName
			*locals = append(*locals, fmt.Sprintf("(local $%s i32)\n", localVarName))
			var strInit strings.Builder
			strInit.WriteString(fmt.Sprintf("(local.set $%s (call $%s (i32.const %d)))\n", localVarName, MemoryAllocateFunc, length))
			for i, char := range e.Value {
				strInit.WriteString(fmt.Sprintf("(i32.store8 offset=%d (local.get $%s) (i32.const %d)) ;; '%c'\n", i, localVarName, char, char))
			}
			strInit.WriteString(fmt.Sprintf("(i32.store8 offset=%d (local.get $%s) (i32.const 0))\n", length-1, localVarName))
			*stringLiterals = append(*stringLiterals, strInit.String())
			// *initializations = append(*initializations, fmt.Sprintf("(local.get $%s)\n", localVarName))
		}
	case *ast.StructLiteral:
		for _, fieldValue := range e.Fields {
			collectExpressionLocals(fieldValue, declaredLocals, locals, initializations, stringLiterals)
		}
	case *ast.StructFieldAccess:
		collectExpressionLocals(e.Left, declaredLocals, locals, initializations, stringLiterals)
	}
}

func generateFunctionCall(call *ast.FunctionCall) string {
	var out strings.Builder

	if call.FunctionName == "println" && len(call.Arguments) > 0 {
		arg := call.Arguments[0]
		if strLiteral, ok := arg.(*ast.StringLiteral); ok {
			localVarName, exists := stringLiteralMap[strLiteral.Value]
			if !exists {
				log.Fatalf("String literal not found: %s", strLiteral.Value)
			}
			out.WriteString(fmt.Sprintf("(call $println (local.get $%s))\n", localVarName))
		} else if fieldAccess, ok := arg.(*ast.StructFieldAccess); ok {
			out.WriteString(fmt.Sprintf("(call $println %s)\n", generateStructFieldAccess(fieldAccess)))
		} else if ident, ok := arg.(*ast.Identifier); ok {
			out.WriteString(fmt.Sprintf("(call $println (local.get $%s))\n", ident.Value))
		} else {
			log.Fatal("println expects a string argument or struct field access")
		}
	} else {
		out.WriteString(fmt.Sprintf("(call $%s ", call.FunctionName))
		for _, arg := range call.Arguments {
			out.WriteString(" ")
			out.WriteString(generateExpression(arg))
		}
		out.WriteString(")\n")
	}

	return out.String()
}
