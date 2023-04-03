package wat

import (
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/internal/ast"
	"github.com/dfirebaugh/punch/internal/token"
)

// GenerateWAT generates WAT IR code for a given AST.
func GenerateWAT(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Program:
		return generateStatements(n.Statements)
	}
	return ""
}

func generateStatements(stmts []ast.Statement) string {
	var out strings.Builder
	out.WriteString("(module\n")
	for _, stmt := range stmts {
		out.WriteString(generateStatement(stmt))
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
					out.WriteString(fmt.Sprintf("(param $%s i32)\n", param.Value))
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
	}
	return ""
}

func generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return fmt.Sprintf("(i32.const %s)", e.Value)
	case *ast.Boolean:
		if e.Value {
			return "(i32.const 1)"
		} else {
			return "(i32.const 0)"
		}
	case *ast.StringLiteral:
		return "(i32.const 0)" // return pointer to string constant
	case *ast.PrefixExpression:
		operand := generateExpression(e.Right)
		switch e.Operator {
		case token.BANG:
			return fmt.Sprintf("(i32.eqz %s)", operand)
		case token.MINUS:
			return fmt.Sprintf("(i32.sub (i32.const 0) %s)", operand)
		}
	case *ast.InfixExpression:
		left, right := generateExpression(e.Left), generateExpression(e.Right)
		switch e.Operator.Type {
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
		}
	case *ast.Identifier:
		return fmt.Sprintf("(local.get $%s)", e.Value)
	case *ast.IfExpression:
		var out strings.Builder
		out.WriteString(fmt.Sprintf("%s\n", generateExpression(e.Condition)))
		out.WriteString("(if\n")
		out.WriteString(fmt.Sprintf("(result i32)\n"))
		out.WriteString(fmt.Sprintf("(param i32 i32)\n"))
		out.WriteString(fmt.Sprintf("(export \"add\")\n"))
		out.WriteString(fmt.Sprintf("(local.get 0)\n"))
		out.WriteString(fmt.Sprintf("(local.get 1)\n"))
		out.WriteString(fmt.Sprintf("(i32.add)\n"))
		out.WriteString(")\n")
		return out.String()
	case *ast.FunctionCall:
		var out strings.Builder
		for _, arg := range e.Arguments {
			out.WriteString(fmt.Sprintf("%s\n", generateExpression(arg)))
		}
		out.WriteString(fmt.Sprintf("(call $%s)\n", e.FunctionName))
		return out.String()
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
