package wat

import (
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

const (
	MemoryAllocateFunc   = "memory_allocate"
	MemoryDeallocateFunc = "memory_deallocate"
)

// Global or accessible map for function declarations
var functionDeclarations map[string]*ast.FunctionDeclaration

func findFunctionDeclarations(node ast.Node) {
	switch n := node.(type) {
	case *ast.FunctionDeclaration:
		functionDeclarations[n.Name.Value] = n
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			findFunctionDeclarations(stmt)
		}
	case *ast.Program:
		for _, stmt := range n.Statements {
			findFunctionDeclarations(stmt)
		}
	}
}

// GenerateWAT generates WAT IR code for a given AST.
func GenerateWAT(node ast.Node, withMemoryManagement bool) string {
	functionDeclarations = make(map[string]*ast.FunctionDeclaration)
	findFunctionDeclarations(node)
	switch n := node.(type) {
	case *ast.Program:
		return generateStatements(n.Statements, withMemoryManagement)
	}
	return ""
}

func generateMemoryManagementFunctions() string {
	return `
;; Declare a memory section with 1 page (64KB)
(memory 1)

;; Global variable to track the current memory allocation position
(global $mem_alloc_ptr (mut i32) (i32.const 0))

(func $memory_allocate (param $size i32) (result i32)
	(local $ptr i32)
	;; Call the function to find a free block
	(local.set $ptr (call $find_free_block (local.get $size)))
	;; Return the pointer
	(local.get $ptr)
)
(func $memory_deallocate (param $ptr i32)
    ;; Call the function to mark the block as free
    (call $mark_block_free (local.get $ptr))
)

;; Function to find a free memory block
(func $find_free_block (param $size i32) (result i32)
  (local $ptr i32)

  ;; Get the current allocation pointer
  (local.set $ptr (global.get $mem_alloc_ptr))

  ;; Update the allocation pointer (naive implementation, no checks for memory limits)
  (global.set $mem_alloc_ptr (i32.add (global.get $mem_alloc_ptr) (local.get $size)))

  ;; Return the starting address of the allocated block
  (local.get $ptr)
)

(data (i32.const 0) "\00\00\00\00") ;; Memory allocation bitmap, initialized to all free

;; mark_block_free -- may need to be revisited
(func $mark_block_free
  (param $ptr i32)  ;; Pointer to the memory block to free
)
`
}

func generateStatements(stmts []ast.Statement, withMemoryManagement bool) string {
	var out strings.Builder
	out.WriteString("(module\n")

	if withMemoryManagement {
		out.WriteString(generateMemoryManagementFunctions())
	}

	for _, stmt := range stmts {
		if stmt == nil {
			continue
		}
		out.WriteString(generateStatement(stmt))
	}
	out.WriteString(")\n")
	return out.String()
}

func generateStringLiteral(str *ast.StringLiteral) string {
	length := len(str.Value) + 1 // Including null terminator
	return fmt.Sprintf(
		`(call $%s (i32.const %d)) ;; allocate memory for string\n`+
			`(data (i32.const 0) "%s\u0000") ;; store string with null terminator\n`,
		MemoryAllocateFunc, length, str.Value)
}

func generateFunctionCall(call *ast.FunctionCall) string {
	var out strings.Builder
	for _, arg := range call.Arguments {
		out.WriteString(generateExpression(arg) + "\n")
	}
	out.WriteString(fmt.Sprintf("(call $%s)\n", call.FunctionName))

	if funcDecl, exists := functionDeclarations[call.FunctionName]; exists {
		if funcDecl.ReturnType.Type == token.STRING {
			out.WriteString(fmt.Sprintf("(call $%s (local.get 0)) ;; deallocate string memory\n", MemoryDeallocateFunc))
		}
	}

	return out.String()
}

func generateInfixExpression(infix *ast.InfixExpression) string {
	left := generateExpression(infix.Left)
	right := generateExpression(infix.Right)
	operator := infix.Operator.Type

	switch operator {
	case token.PLUS:
		return fmt.Sprintf("(i32.add %s %s)", left, right)
	case token.MINUS:
		return fmt.Sprintf("(i32.sub %s %s)", left, right)
	case token.ASTERISK:
		return fmt.Sprintf("(i32.mul %s %s)", left, right)
	case token.SLASH:
		return fmt.Sprintf("(i32.div_s %s %s)", left, right)
	case token.LT:
		return fmt.Sprintf("(i32.lt_s %s %s)", left, right)
	case token.GT:
		return fmt.Sprintf("(i32.gt_s %s %s)", left, right)
	case token.EQ:
		return fmt.Sprintf("(i32.eq %s %s)", left, right)
	case token.NOT_EQ:
		return fmt.Sprintf("(i32.ne %s %s)", left, right)
	default:
		println("\t\t operator:", operator)
		return fmt.Sprintf(";; unhandled operator: %s\n", operator)
	}
}

func generatePrefixExpression(prefix *ast.PrefixExpression) string {
	operand := generateExpression(prefix.Right)
	operator := prefix.Operator.Type
	switch operator {
	case token.BANG:
		return fmt.Sprintf("(i32.eqz %s)", operand)
	case token.MINUS:
		return fmt.Sprintf("(i32.sub (i32.const 0) %s)", operand)
	default:
		return fmt.Sprintf(";; unhandled operator: %s\n", operator)
	}
}

func generateBlockStatement(block *ast.BlockStatement) string {
	var out strings.Builder
	for _, stmt := range block.Statements {
		out.WriteString("\t")
		out.WriteString(generateStatement(stmt))
		out.WriteString("\n")
	}
	return out.String()
}

func generateIfStatement(e *ast.IfStatement) string {
	var out strings.Builder
	out.WriteString("(if\n")
	out.WriteString("(i32.eqz ")
	out.WriteString(generateExpression(e.Condition))
	out.WriteString(")\n")
	out.WriteString("  (then\n")
	out.WriteString(generateBlockStatement(e.Consequence))
	out.WriteString("  )\n")
	if e.Alternative != nil {
		out.WriteString("  (else\n")
		out.WriteString(generateBlockStatement(e.Alternative))
		out.WriteString("  )\n")
	}
	out.WriteString(")\n")
	return out.String()
}

func generateStatement(stmt ast.Statement) string {
	if stmt == nil {
		return ""
	}

	switch s := stmt.(type) {
	case *ast.LetStatement:
		// check if value is a string literal
		if str, ok := s.Value.(*ast.StringLiteral); ok {
			// declare string constant in data section
			return fmt.Sprintf(
				`(data (i32.const 0) "%s")`+"\n"+
					"(local $%s i32)\n"+
					`(i32.const 0)`+"\n"+ // load address of string
					`(i32.load)`+"\n"+ // load pointer to string value
					`(local.set $%s)`+"\n",
				str.Value, s.Name.Value, s.Name.Value,
			)
		}
		// handle other value types as before
		return fmt.Sprintf(
			"(local $%s i32)\n"+
				"%s\n"+
				"(local.set $%s %s)\n",
			s.Name.Value,
			generateExpression(s.Value),
			s.Name.Value,
			s.Name.Value,
		)
	case *ast.ReturnStatement:
		if s.ReturnValue == nil {
			return "(return)\n"
		}
		return fmt.Sprintf("(return %s)\n", generateExpression(s.ReturnValue))
	case *ast.BlockStatement:
		var out strings.Builder
		out.WriteString("\n")
		for _, stmt := range s.Statements {
			out.WriteString(generateStatement(stmt))
		}
		return out.String()
	case *ast.IfStatement:
		return generateIfStatement(s)
	case *ast.FunctionStatement:
		var out strings.Builder
		if s == nil {
			return ""
		}
		returnType := "i32" // default return type is i32
		// if s.ReturnType != nil {
		// 	returnType = s.ReturnType.Value // use specified return type, if any
		// }
		out.WriteString(fmt.Sprintf("(func $%s ", s.Name.Value))
		if s.IsExported {
			out.WriteString(fmt.Sprintf("(export \"%s\")", s.Name))
		}
		if s.Parameters != nil {
			for _, param := range s.Parameters {
				if param != nil {
					out.WriteString(fmt.Sprintf("(param $%s i32)\n", param.Identifier.Value))
				} else {
					out.WriteString("(param)\n")
				}
			}
		}
		out.WriteString(fmt.Sprintf("(result %s)\n", returnType))
		bodyStr := generateStatement(s.Body)
		if bodyStr != "" {
			out.WriteString(fmt.Sprintf("%s\n", bodyStr))
		}
		out.WriteString(")\n")
		return out.String()
	case *ast.ExpressionStatement:
		return generateExpression(s.Expression)
	}
	return ""
}

func generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return fmt.Sprintf("(i32.const %d)", e.Value)
	case *ast.Boolean:
		if e.Value {
			return "(i32.const 1)"
		} else {
			return "(i32.const 0)"
		}
	case *ast.StringLiteral:
		return generateStringLiteral(e)
	case *ast.PrefixExpression:
		return generatePrefixExpression(e)
	case *ast.InfixExpression:
		return generateInfixExpression(e)
	case *ast.Identifier:
		return fmt.Sprintf("(local.get $%s)", e.Value)
	case *ast.FunctionCall:
		return generateFunctionCall(e)
	case *ast.ArrayLiteral:
		var out strings.Builder
		out.WriteString(fmt.Sprintf("(i32.const %d)\n", len(e.Elements)))
		out.WriteString(fmt.Sprintf("(i32.const %d)\n", len(e.Elements)))
		out.WriteString("(i32.mul)\n")
		out.WriteString("(memory.grow 1)\n")
		out.WriteString("(i32.store8)\n")
		for _, elem := range e.Elements {
			out.WriteString(fmt.Sprintf("%s\n", generateExpression(elem)))
		}
		return out.String()
	case *ast.IndexExpression:
		var out strings.Builder
		out.WriteString(fmt.Sprintf("%s\n", generateExpression(e.Left)))
		out.WriteString(fmt.Sprintf("%s\n", generateExpression(e.Index)))
		out.WriteString("(i32.add)\n")
		out.WriteString("(i32.load8_s)\n")
		return out.String()
	case *ast.CallExpression:
		var out strings.Builder
		for _, arg := range e.Arguments {
			out.WriteString(fmt.Sprintf("%s\n", generateExpression(arg)))
		}
		out.WriteString(fmt.Sprintf("(call $%s)\n", e.Function.Value))
		return out.String()
	case *ast.WhileExpression:
		var out strings.Builder
		out.WriteString("(block\n")
		out.WriteString(fmt.Sprintf("(loop $%d\n", e.ID))
		out.WriteString(fmt.Sprintf("%s\n", generateExpression(e.Condition)))
		out.WriteString("(br_if 1\n")
		out.WriteString(fmt.Sprintf("%s\n", generateStatement(e.Body)))
		out.WriteString(")\n")
		out.WriteString("(br 0)\n")
		out.WriteString(")\n")
		out.WriteString(")\n")
		return out.String()
	}
	return ""
}
