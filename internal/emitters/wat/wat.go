package wat

import (
	"fmt"
	"log"
	"strings"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

const (
	MemoryAllocateFunc   = "memory_allocate"
	MemoryDeallocateFunc = "memory_deallocate"
)

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
		for _, stmt := range n.Files[0].Statements {
			findFunctionDeclarations(stmt)
		}
	}
}

func GenerateWAT(node ast.Node, withMemoryManagement bool) string {
	functionDeclarations = make(map[string]*ast.FunctionDeclaration)
	findFunctionDeclarations(node)
	switch n := node.(type) {
	case *ast.Program:
		return generateStatements(n.Files[0].Statements, withMemoryManagement)
	}
	return ""
}

func generateImports() string {
	return `
(import "imports" "println" (func $println (param i32)))
`
}

func generateMemoryManagementFunctions() string {
	return `
;; Declare a memory section with 1 page (64KB)
(memory 1)
(export "memory" (memory 0))

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

func mapTypeToWAT(t string) string {
	switch t {
	case "u8", "i8", "u16", "i16", "u32", "i32", "bool":
		return "i32"
	case "u64", "i64":
		return "i64"
	case "f8", "f16", "f32", "float":
		return "f32"
	case "f64":
		return "f64"
	default:
		log.Fatalf("Unsupported type: %s", t)
		return ""
	}
}

func generateStatements(stmts []ast.Statement, withMemoryManagement bool) string {
	var out strings.Builder
	out.WriteString("(module\n")

	out.WriteString(generateImports())
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

var localVarCounter int

func generateUniqueLocalVarName(base string) string {
	localVarCounter++
	return fmt.Sprintf("%s_%d", base, localVarCounter)
}

func generateStringLiteral(str *ast.StringLiteral) string {
	length := len(str.Value) + 1
	var out strings.Builder

	localVarName := generateUniqueLocalVarName("str_ptr")

	out.WriteString(fmt.Sprintf("(local $%s i32)\n", localVarName))

	out.WriteString(fmt.Sprintf("(local.set $%s (call $%s (i32.const %d)))\n", localVarName, MemoryAllocateFunc, length))

	for i, char := range str.Value {
		if char == '\\' && i+1 < len(str.Value) && str.Value[i+1] == '"' {
			out.WriteString(fmt.Sprintf("(i32.store8 offset=%d (local.get $%s) (i32.const %d)) ;; '\"'\n", i, localVarName, '"'))
			i++
		} else {
			out.WriteString(fmt.Sprintf("(i32.store8 offset=%d (local.get $%s) (i32.const %d)) ;; '%c'\n", i, localVarName, char, char))
		}
	}

	out.WriteString(fmt.Sprintf("(i32.store8 offset=%d (local.get $%s) (i32.const 0))\n", length-1, localVarName))

	// out.WriteString(fmt.Sprintf("(local.get $%s)\n", localVarName))

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
	if e == nil {
		return ""
	}
	var out strings.Builder
	out.WriteString("\t\t(if ")
	out.WriteString(generateExpression(e.Condition))
	out.WriteString("\n\t\t\t(then\n")
	out.WriteString(generateBlockStatement(e.Consequence))
	out.WriteString("\n\t\t\t)")
	if e.Alternative != nil {
		out.WriteString("\n\t\t\t(else\n")
		out.WriteString(generateBlockStatement(e.Alternative))
		out.WriteString("\n\t\t\t)")
	}
	out.WriteString("\n\t\t)")
	return out.String()
}

func collectLocals(stmt ast.Statement, declaredLocals map[string]bool, locals *[]string) {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		if !declaredLocals[s.Name.Value] {
			*locals = append(*locals, fmt.Sprintf("(local $%s %s)\n", s.Name.Value, mapTypeToWAT(s.Type.Literal)))
			declaredLocals[s.Name.Value] = true
		}
	case *ast.BlockStatement:
		for _, stmt := range s.Statements {
			collectLocals(stmt, declaredLocals, locals)
		}
	case *ast.IfStatement:
		collectLocals(s.Consequence, declaredLocals, locals)
		if s.Alternative != nil {
			collectLocals(s.Alternative, declaredLocals, locals)
		}
	case *ast.FunctionStatement:
		for _, stmt := range s.Body.Statements {
			collectLocals(stmt, declaredLocals, locals)
		}
	case *ast.ExpressionStatement:
		exp := stmt.(*ast.ExpressionStatement)
		fn, ok := exp.Expression.(*ast.FunctionCall)
		if !ok {
			break
		}
		collectLocals(stmt, declaredLocals, locals)
		for _, arg := range fn.Arguments {
			println(arg.String())
		}
	}
}

func generateVariableDeclaration(decl *ast.VariableDeclaration) string {
	var out strings.Builder
	watType := mapTypeToWAT(decl.Type.Literal)

	if decl.Value != nil {
		out.WriteString(fmt.Sprintf("(local.set $%s ", decl.Name.Value))
		out.WriteString(generateExpression(decl.Value))
		out.WriteString(") \n")
	} else {
		out.WriteString(fmt.Sprintf("(local $%s %s)\n", decl.Name.Value, watType))
	}

	return out.String()
}

func generateReturnStatement(s *ast.ReturnStatement) string {
	if s == nil {
		log.Println("Encountered nil *ast.ReturnStatement")
		return ""
	}

	if len(s.ReturnValues) == 0 {
		return "\t\t(return (i32.const 0)) ;; No return values, return null pointer\n"
	} else if len(s.ReturnValues) == 1 {
		return fmt.Sprintf("\t\t(return %s)\n", generateExpression(s.ReturnValues[0]))
	} else {
		var out strings.Builder
		out.WriteString("\t\t(local $retPtr i32)\n")
		out.WriteString(fmt.Sprintf("\t\t(local.set $retPtr (call $%s (i32.const %d)))\n", MemoryAllocateFunc, len(s.ReturnValues)*4))

		for i, retVal := range s.ReturnValues {
			expr := generateExpression(retVal)
			out.WriteString(fmt.Sprintf("\t\t(i32.store offset=%d (local.get $retPtr) %s)\n", i*4, expr))
		}

		out.WriteString("\t\t(return (local.get $retPtr))\n")
		return out.String()
	}
}

func generateStatement(stmt ast.Statement) string {
	if stmt == nil {
		return ""
	}

	switch s := stmt.(type) {
	case *ast.LetStatement:
		if str, ok := s.Value.(*ast.StringLiteral); ok {
			return fmt.Sprintf(
				`(data (i32.const 0) "%s")`+"\n"+
					"(local $%s i32)\n"+
					`(i32.const 0)`+"\n"+
					`(i32.load)`+"\n"+
					`(local.set $%s)`+"\n",
				str.Value, s.Name.Value, s.Name.Value,
			)
		}
		return fmt.Sprintf(
			"(local $%s i32)\n"+
				"%s\n"+
				"(local.set $%s %s)\n",
			s.Name.Value,
			generateExpression(s.Value),
			s.Name.Value,
			s.Name.Value,
		)
	case *ast.VariableDeclaration:
		return generateVariableDeclaration(s)
	case *ast.ReturnStatement:
		return generateReturnStatement(s)
	case *ast.IfStatement:
		return generateIfStatement(s)
	case *ast.FunctionStatement:
		return generateFunctionStatement(s)
	case *ast.ExpressionStatement:
		return generateExpression(s.Expression)
	}
	return ""
}

func generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return fmt.Sprintf("(%s.const %d)", mapTypeToWAT(string(e.Token.Type)), e.Value)
	case *ast.FloatLiteral:
		return fmt.Sprintf("(%s.const %f)", mapTypeToWAT(string(e.Token.Type)), e.Value)
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
