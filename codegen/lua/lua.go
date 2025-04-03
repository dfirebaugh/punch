package lua

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

const (
	LuaFunction       = "function"
	LuaLocal          = "local"
	LuaReturn         = "return"
	LuaIf             = "if"
	LuaElse           = "else"
	LuaEnd            = "end"
	LuaThen           = "then"
	LuaNil            = "nil"
	LuaUnsupported    = "-- Unsupported"
	LuaFileComment    = "-- File: %s\n"
	LuaPackageComment = "-- Package: %s\n\n"
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

	out.WriteString(fmt.Sprintf(LuaFileComment, file.Filename))
	out.WriteString(fmt.Sprintf(LuaPackageComment, file.PackageName))

	for _, stmt := range file.Statements {
		out.WriteString(t.transpileStatement(stmt))
		out.WriteString("\n")
	}

	return out.String()
}

func (t *Transpiler) transpileStatement(stmt ast.Statement) string {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		return t.transpileExpression(stmt.Expression)

	case *ast.LetStatement:
		return fmt.Sprintf("%s %s = %s", LuaLocal, stmt.Name.String(), t.transpileExpression(stmt.Value))

	case *ast.ReturnStatement:
		return LuaReturn + " " + t.transpileExpressions(stmt.ReturnValues)

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
		return LuaUnsupported + " statement"
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

	case *ast.StructFieldAssignment:
		return t.transpileStructFieldAssignment(expr)

	case *ast.ListLiteral:
		return t.transpileListLiteral(expr)

	default:
		return LuaUnsupported + " expression"
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
	out.WriteString(LuaFunction + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString("(")

	params := []string{}
	for _, param := range stmt.Parameters {
		params = append(params, param.Identifier.Token.Literal)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")

	out.WriteString(t.transpileBlockStatement(stmt.Body))
	out.WriteString(LuaEnd)
	return out.String()
}

func (t *Transpiler) transpileIfStatement(stmt *ast.IfStatement) string {
	var out bytes.Buffer

	out.WriteString(LuaIf + " ")
	out.WriteString(t.transpileExpression(stmt.Condition))
	out.WriteString(" " + LuaThen + "\n")
	out.WriteString(t.transpileBlockStatement(stmt.Consequence))

	if stmt.Alternative != nil {
		out.WriteString("\n" + LuaElse + "\n")
		out.WriteString(t.transpileBlockStatement(stmt.Alternative))
	}

	out.WriteString("\n" + LuaEnd)
	return out.String()
}

func (t *Transpiler) transpileBlockStatement(stmt *ast.BlockStatement) string {
	var out bytes.Buffer

	for _, s := range stmt.Statements {
		out.WriteString(t.transpileStatement(s))
		out.WriteString("\n")
	}

	return out.String()
}

func (t *Transpiler) transpileFunctionCall(expr *ast.FunctionCall) string {
	var out bytes.Buffer

	if expr.Function.String() == "println" {
		out.WriteString("print(")
	} else if expr.Function.String() == "len" && len(expr.Arguments) == 1 {
		out.WriteString("#" + t.transpileExpression(expr.Arguments[0]))
		return out.String()
	} else if expr.Function.String() == "append" && len(expr.Arguments) == 2 {
		out.WriteString(t.transpileExpression(expr.Arguments[0]) + "[#" + t.transpileExpression(expr.Arguments[0]) + "+1] = ")
		out.WriteString(t.transpileExpression(expr.Arguments[1]))
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

	out.WriteString(LuaFunction + "(")
	params := []string{}
	for _, param := range expr.Parameters {
		params = append(params, param.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")
	out.WriteString(t.transpileBlockStatement(expr.Body))
	out.WriteString(LuaEnd)

	return out.String()
}

func (t *Transpiler) transpileArrayLiteral(expr *ast.ArrayLiteral) string {
	var out bytes.Buffer

	out.WriteString("{")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, t.transpileExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (t *Transpiler) transpileStructDefinition(stmt *ast.StructDefinition) string {
	var out bytes.Buffer

	if t.definedStructs == nil {
		t.definedStructs = make(map[string]bool)
	}
	t.definedStructs[stmt.Name.String()] = true

	out.WriteString(fmt.Sprintf("%s = {}\n", stmt.Name.String()))
	out.WriteString(fmt.Sprintf("%s.__index = %s\n", stmt.Name.String(), stmt.Name.String()))
	out.WriteString(fmt.Sprintf("function %s:new(o)\n", stmt.Name.String()))
	out.WriteString("o = o or {}\n")
	out.WriteString("setmetatable(o, self)\n")
	out.WriteString("return o\n")
	out.WriteString(LuaEnd + "\n")

	for _, field := range stmt.Fields {
		out.WriteString(fmt.Sprintf("function %s:get_%s()\n", stmt.Name.String(), field.Name.String()))
		out.WriteString(fmt.Sprintf("return self.%s\n", field.Name.String()))
		out.WriteString(LuaEnd + "\n")

		out.WriteString(fmt.Sprintf("function %s:set_%s(value)\n", stmt.Name.String(), field.Name.String()))
		out.WriteString(fmt.Sprintf("self.%s = value\n", field.Name.String()))
		out.WriteString(LuaEnd + "\n")
	}

	return out.String()
}

func (t *Transpiler) transpileStructLiteral(expr *ast.StructLiteral) string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s:new({", expr.StructName.String()))
	fields := []string{}
	for name, value := range expr.Fields {
		fields = append(fields, fmt.Sprintf("%s = %s", name, t.transpileExpression(value)))
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

func (t *Transpiler) transpileStructFieldAssignment(expr *ast.StructFieldAssignment) string {
	return fmt.Sprintf("%s.%s = %s",
		t.transpileExpression(expr.Left.Left),
		expr.Left.Field.String(),
		t.transpileExpression(expr.Right),
	)
}

func (t *Transpiler) transpileVariableDeclaration(stmt *ast.VariableDeclaration) string {
	if stmt.Type.Type == token.IDENTIFIER && t.definedStructs[stmt.Type.Literal] {
		return fmt.Sprintf("%s %s = %s:new(%s)",
			LuaLocal,
			stmt.Name.String(),
			stmt.Type.Literal,
			t.transpileExpression(stmt.Value),
		)
	}
	if stmt.Type.Type == token.IDENTIFIER {
		return fmt.Sprintf("%s = %s",
			stmt.Name.String(),
			t.transpileExpression(stmt.Value),
		)
	}

	return fmt.Sprintf("%s %s = %s",
		LuaLocal,
		stmt.Name.String(),
		t.transpileExpression(stmt.Value),
	)
}

func (t *Transpiler) transpileForStatement(stmt *ast.ForStatement) string {
	var out bytes.Buffer

	out.WriteString("for ")

	if stmt.Init != nil {
		if letStmt, ok := stmt.Init.(*ast.VariableDeclaration); ok {
			out.WriteString(fmt.Sprintf("%s = %s, ",
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
		out.WriteString(t.transpileExpression(stmt.Condition) + ", ")
	} else {
		out.WriteString(";")
	}

	if stmt.Post != nil {
		if letStmt, ok := stmt.Post.(*ast.VariableDeclaration); ok {
			out.WriteString(letStmt.Name.String() + " = " + t.transpileExpression(letStmt.Value))
		} else if exprStmt, ok := stmt.Post.(*ast.ExpressionStatement); ok {
			out.WriteString(t.transpileExpression(exprStmt.Expression))
		} else {
			out.WriteString(t.transpileStatement(stmt.Post))
		}
	}

	out.WriteString("\n")
	out.WriteString(t.transpileBlockStatement(stmt.Body))
	out.WriteString(LuaEnd)
	return out.String()
}

func (t *Transpiler) transpileListLiteral(expr *ast.ListLiteral) string {
	var out bytes.Buffer

	out.WriteString("{")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, t.transpileExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (t *Transpiler) transpileListDeclaration(stmt *ast.ListDeclaration) string {
	var out bytes.Buffer

	out.WriteString(LuaLocal + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" = ")
	out.WriteString(t.transpileListLiteral(stmt.Value))
	out.WriteString(";")

	return out.String()
}
