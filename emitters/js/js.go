package js

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

const (
	JSFunction       = "function"
	JSLet            = "let"
	JSReturn         = "return"
	JSIf             = "if"
	JSElse           = "else"
	JSClass          = "class"
	JSConstructor    = "constructor"
	JSNew            = "new"
	JSExport         = "export"
	JSConsoleLog     = "console.log"
	JSUnsupported    = "// Unsupported"
	JSFileComment    = "// File: %s\n"
	JSPackageComment = "// Package: %s\n\n"
)

type Transpiler struct {
	definedStructs map[string]bool
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		definedStructs: make(map[string]bool),
	}
}

func (t *Transpiler) Transpile(program *ast.Program) (string, error) {
	var out bytes.Buffer

	for _, file := range program.Files {
		out.WriteString(t.transpileFile(file))
		out.WriteString("\n")
	}

	return out.String(), nil
}

func (t *Transpiler) transpileFile(file *ast.File) string {
	var out bytes.Buffer
	var exports []string

	out.WriteString(fmt.Sprintf(JSFileComment, file.Filename))
	out.WriteString(fmt.Sprintf(JSPackageComment, file.PackageName))

	for _, stmt := range file.Statements {
		if functionStmt, ok := stmt.(*ast.FunctionStatement); ok && functionStmt.IsExported {
			exports = append(exports, functionStmt.Name.String())
		}
		out.WriteString(t.transpileStatement(stmt))
		out.WriteString("\n")
	}

	if len(exports) > 0 {
		out.WriteString("\n" + JSExport + " {\n")
		for _, export := range exports {
			out.WriteString(fmt.Sprintf("  %s,\n", export))
		}
		out.WriteString("};\n")
	}

	return out.String()
}

func (t *Transpiler) transpileStatement(stmt ast.Statement) string {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		return t.transpileExpression(stmt.Expression) + ";"

	case *ast.LetStatement:
		return fmt.Sprintf("%s %s = %s;", JSLet, stmt.Name.String(), t.transpileExpression(stmt.Value))

	case *ast.ReturnStatement:
		return JSReturn + " " + t.transpileExpressions(stmt.ReturnValues) + ";"

	case *ast.FunctionStatement:
		return t.transpileFunctionStatement(stmt)

	case *ast.IfStatement:
		return t.transpileIfStatement(stmt)

	case *ast.BlockStatement:
		return t.transpileBlockStatement(stmt)

	case *ast.StructDefinition:
		return t.transpileStructDefinition(stmt)

	case *ast.VariableDeclaration:
		return t.transpileVariableDeclaration(stmt)

	case *ast.ForStatement:
		return t.transpileForStatement(stmt)

	case *ast.ListDeclaration:
		return t.transpileListDeclaration(stmt)

	default:
		return JSUnsupported + " statement"
	}
}

func (t *Transpiler) transpileExpression(expr ast.Expression) string {
	switch expr := expr.(type) {
	case *ast.Identifier:
		return expr.String()

	case *ast.IntegerLiteral:
		return expr.String()

	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, expr.String())

	case *ast.BooleanLiteral:
		return expr.String()

	case *ast.BinaryExpression:
		return fmt.Sprintf("(%s %s %s)",
			t.transpileExpression(expr.Left),
			expr.Operator.Literal,
			t.transpileExpression(expr.Right),
		)

	case *ast.CallExpression:
		return t.transpileCallExpression(expr)

	case *ast.FunctionCall:
		return t.transpileFunctionCall(expr)

	case *ast.IndexExpression:
		return fmt.Sprintf("%s[%s]",
			t.transpileExpression(expr.Left),
			t.transpileExpression(expr.Index),
		)

	case *ast.PrefixExpression:
		return fmt.Sprintf("(%s%s)",
			expr.Operator.Literal,
			t.transpileExpression(expr.Right),
		)

	case *ast.InfixExpression:
		return fmt.Sprintf("(%s %s %s)",
			t.transpileExpression(expr.Left),
			expr.Operator.Literal,
			t.transpileExpression(expr.Right),
		)

	case *ast.AssignmentExpression:
		return t.transpileAssignmentExpression(expr)

	case *ast.FunctionLiteral:
		return t.transpileFunctionLiteral(expr)

	case *ast.ArrayLiteral:
		return t.transpileArrayLiteral(expr)

	case *ast.StructLiteral:
		return t.transpileStructLiteral(expr)

	case *ast.StructFieldAccess:
		return t.transpileStructFieldAccess(expr)

	case *ast.ListLiteral:
		return t.transpileListLiteral(expr)

	default:
		return JSUnsupported + " expression"
	}
}

func (t *Transpiler) transpileAssignmentExpression(expr *ast.AssignmentExpression) string {
	return fmt.Sprintf("%s = %s",
		t.transpileExpression(expr.Left),
		t.transpileExpression(expr.Right),
	)
}

func (t *Transpiler) transpileExpressions(exprs []ast.Expression) string {
	var out []string
	for _, expr := range exprs {
		out = append(out, t.transpileExpression(expr))
	}
	return strings.Join(out, ", ")
}

func (t *Transpiler) transpileFunctionStatement(stmt *ast.FunctionStatement) string {
	var out bytes.Buffer
	out.WriteString(JSFunction + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString("(")

	params := []string{}
	for _, param := range stmt.Parameters {
		params = append(params, param.Identifier.Token.Literal)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	out.WriteString(t.transpileBlockStatement(stmt.Body))
	return out.String()
}

func (t *Transpiler) transpileIfStatement(stmt *ast.IfStatement) string {
	var out bytes.Buffer

	out.WriteString(JSIf + " (")
	out.WriteString(t.transpileExpression(stmt.Condition))
	out.WriteString(") ")
	out.WriteString(t.transpileBlockStatement(stmt.Consequence))

	if stmt.Alternative != nil {
		out.WriteString(" " + JSElse + " ")
		out.WriteString(t.transpileBlockStatement(stmt.Alternative))
	}

	return out.String()
}

func (t *Transpiler) transpileBlockStatement(stmt *ast.BlockStatement) string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, s := range stmt.Statements {
		out.WriteString(t.transpileStatement(s))
		out.WriteString("\n")
	}
	out.WriteString("}")

	return out.String()
}

func (t *Transpiler) transpileFunctionCall(expr *ast.FunctionCall) string {
	var out bytes.Buffer

	if expr.Function.String() == "println" {
		out.WriteString(JSConsoleLog + "(")
	} else if expr.Function.String() == "len" && len(expr.Arguments) == 1 {
		out.WriteString(t.transpileExpression(expr.Arguments[0]) + ".length")
		return out.String()
	} else if expr.Function.String() == "append" && len(expr.Arguments) == 2 {
		out.WriteString(t.transpileExpression(expr.Arguments[0]) + ".push(")
		out.WriteString(t.transpileExpression(expr.Arguments[1]))
		out.WriteString(")")
		return out.String()
	} else {
		out.WriteString(expr.Function.String() + "(")
	}

	args := []string{}
	for _, arg := range expr.Arguments {
		args = append(args, t.transpileExpression(arg))
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (t *Transpiler) transpileCallExpression(expr *ast.CallExpression) string {
	return t.transpileFunctionCall(&ast.FunctionCall{
		Function:  expr.Function,
		Arguments: expr.Arguments,
	})
}

func (t *Transpiler) transpileFunctionLiteral(expr *ast.FunctionLiteral) string {
	var out bytes.Buffer

	out.WriteString(JSFunction + "(")
	params := []string{}
	for _, param := range expr.Parameters {
		params = append(params, param.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(t.transpileBlockStatement(expr.Body))

	return out.String()
}

func (t *Transpiler) transpileArrayLiteral(expr *ast.ArrayLiteral) string {
	var out bytes.Buffer

	out.WriteString("[")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, t.transpileExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (t *Transpiler) transpileStructDefinition(stmt *ast.StructDefinition) string {
	var out bytes.Buffer

	if t.definedStructs == nil {
		t.definedStructs = make(map[string]bool)
	}
	t.definedStructs[stmt.Name.String()] = true

	out.WriteString(JSClass + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" {\n")
	out.WriteString(JSConstructor + "({")
	fields := []string{}
	for _, field := range stmt.Fields {
		fields = append(fields, field.Name.String())
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString("}) {\n")
	for _, field := range stmt.Fields {
		out.WriteString(fmt.Sprintf("this.%s = %s;\n", field.Name.String(), field.Name.String()))
	}
	out.WriteString("}\n")
	out.WriteString("}")

	return out.String()
}

func (t *Transpiler) transpileStructLiteral(expr *ast.StructLiteral) string {
	var out bytes.Buffer

	out.WriteString(JSNew + " ")
	out.WriteString(expr.StructName.String())
	out.WriteString("({")
	fields := []string{}
	for name, value := range expr.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", name, t.transpileExpression(value)))
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString("})")

	return out.String()
}

func (t *Transpiler) transpileStructFieldAccess(expr *ast.StructFieldAccess) string {
	return fmt.Sprintf("%s.%s",
		t.transpileExpression(expr.Left),
		expr.Field.String(),
	)
}

func (t *Transpiler) transpileVariableDeclaration(stmt *ast.VariableDeclaration) string {
	if stmt.Type.Type == token.IDENTIFIER && t.definedStructs[stmt.Type.Literal] {
		return fmt.Sprintf("const %s = new %s(%s);",
			stmt.Name.String(),
			stmt.Type.Literal,
			t.transpileExpression(stmt.Value),
		)
	}
	if stmt.Type.Type == token.IDENTIFIER {
		return fmt.Sprintf("%s = %s;",
			stmt.Name.String(),
			t.transpileExpression(stmt.Value),
		)
	}

	return fmt.Sprintf("%s %s = %s;",
		JSLet,
		stmt.Name.String(),
		t.transpileExpression(stmt.Value),
	)
}

func (t *Transpiler) transpileForStatement(stmt *ast.ForStatement) string {
	var out bytes.Buffer

	out.WriteString("for (")

	if stmt.Init != nil {
		if letStmt, ok := stmt.Init.(*ast.VariableDeclaration); ok {
			out.WriteString(fmt.Sprintf("%s %s = %s;",
				JSLet,
				letStmt.Name.String(),
				t.transpileExpression(letStmt.Value),
			))
		} else {
			out.WriteString(t.transpileStatement(stmt.Init))
		}
	} else {
		out.WriteString(";")
	}

	if stmt.Condition != nil {
		out.WriteString(" " + t.transpileExpression(stmt.Condition) + ";")
	} else {
		out.WriteString(";")
	}

	if stmt.Post != nil {
		if letStmt, ok := stmt.Post.(*ast.VariableDeclaration); ok {
			out.WriteString(" " + letStmt.Name.String() + " = " + t.transpileExpression(letStmt.Value))
		} else if exprStmt, ok := stmt.Post.(*ast.ExpressionStatement); ok {
			out.WriteString(" " + t.transpileExpression(exprStmt.Expression))
		} else {
			out.WriteString(" " + t.transpileStatement(stmt.Post))
		}
	}

	out.WriteString(") ")
	out.WriteString(t.transpileBlockStatement(stmt.Body))
	return out.String()
}

func (t *Transpiler) transpileListLiteral(expr *ast.ListLiteral) string {
	var out bytes.Buffer

	out.WriteString("[")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, t.transpileExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (t *Transpiler) transpileListDeclaration(stmt *ast.ListDeclaration) string {
	var out bytes.Buffer

	out.WriteString(JSLet + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" = ")
	out.WriteString(t.transpileListLiteral(stmt.Value))
	out.WriteString(";")

	return out.String()
}
