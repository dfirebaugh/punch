package punchgen

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/dfirebaugh/punch/ast"
)

type Generator struct {
	indentLevel int
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) indent() {
	g.indentLevel++
}

func (g *Generator) unindent() {
	if g.indentLevel > 0 {
		g.indentLevel--
	}
}

func (g *Generator) currentIndent() string {
	return strings.Repeat("\t", g.indentLevel)
}

func (g *Generator) Generate(program *ast.Program) (string, error) {
	var out bytes.Buffer

	for _, file := range program.Files {
		out.WriteString(g.generateFile(file))
		out.WriteString("\n")
	}

	// generated := g.removeExcessiveBlankLines(out.String())

	// return generated, nil
	return out.String(), nil
}

func (g *Generator) removeExcessiveBlankLines(code string) string {
	re := regexp.MustCompile("\n{3,}")
	return re.ReplaceAllString(code, "\n\n")
}

func (g *Generator) generateFile(file *ast.File) string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("pkg %s\n\n", file.PackageName))

	for _, stmt := range file.Statements {
		out.WriteString(g.generateStatement(stmt))
		out.WriteString("\n\n")
	}

	return out.String()
}

func (g *Generator) generateStatement(stmt ast.Statement) string {
	switch stmt := stmt.(type) {
	case *ast.FunctionStatement:
		return g.generateFunctionStatement(stmt)

	case *ast.ReturnStatement:
		return g.generateReturnStatement(stmt)

	case *ast.LetStatement:
		return g.generateLetStatement(stmt)

	case *ast.VariableDeclaration:
		return g.generateVariableDeclaration(stmt)

	case *ast.IfStatement:
		return g.generateIfStatement(stmt)

	case *ast.StructDefinition:
		return g.generateStructDefinition(stmt)

	case *ast.ExpressionStatement:
		return g.currentIndent() + g.generateExpression(stmt.Expression)

	case *ast.BlockStatement:
		return g.generateBlockStatement(stmt)

	case *ast.ListDeclaration:
		return g.generateListDeclaration(stmt)

	case *ast.ForStatement:
		return g.generateForStatement(stmt)

	default:
		return g.currentIndent() + "// Not Supported statement"
	}
}

func (g *Generator) generateFunctionStatement(stmt *ast.FunctionStatement) string {
	var out bytes.Buffer

	if stmt.IsExported {
		out.WriteString("pub ")
	}

	if stmt.ReturnType != nil && stmt.ReturnType.Value != "void" {
		out.WriteString(fmt.Sprintf("%s %s(", stmt.ReturnType.Value, stmt.Name.Value))
	} else {
		out.WriteString(fmt.Sprintf("fn %s(", stmt.Name.Value))
	}

	params := []string{}
	for _, param := range stmt.Parameters {
		params = append(params, fmt.Sprintf("%s %s", param.Identifier.Value, param.Type))
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	out.WriteString(g.generateBlockStatement(stmt.Body))

	return out.String()
}

func (g *Generator) generateReturnStatement(stmt *ast.ReturnStatement) string {
	return g.currentIndent() + fmt.Sprintf("return %s", g.generateExpressions(stmt.ReturnValues))
}

func (g *Generator) generateLetStatement(stmt *ast.LetStatement) string {
	return g.currentIndent() + fmt.Sprintf("let %s = %s", stmt.Name.Value, g.generateExpression(stmt.Value))
}

func (g *Generator) generateVariableDeclaration(stmt *ast.VariableDeclaration) string {
	return g.currentIndent() + fmt.Sprintf("%s %s = %s", stmt.Type.Literal, stmt.Name.Token.Literal, g.generateExpression(stmt.Value))
}

func (g *Generator) generateIfStatement(stmt *ast.IfStatement) string {
	var out bytes.Buffer

	out.WriteString(g.currentIndent() + fmt.Sprintf("if %s ", g.generateExpression(stmt.Condition)))
	out.WriteString(g.generateBlockStatement(stmt.Consequence))

	if stmt.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(g.generateBlockStatement(stmt.Alternative))
	}

	return out.String()
}

func (g *Generator) generateStructDefinition(stmt *ast.StructDefinition) string {
	var out bytes.Buffer

	out.WriteString(g.currentIndent() + fmt.Sprintf("struct %s {\n", stmt.Name.Value))
	g.indent()
	for _, field := range stmt.Fields {
		out.WriteString(g.currentIndent() + fmt.Sprintf("%s %s\n", field.Type.Literal, field.Name.Value))
	}
	g.unindent()
	out.WriteString(g.currentIndent() + "}")

	return out.String()
}

func (g *Generator) generateStructFieldAssignment(stmt *ast.StructFieldAssignment) string {
	return g.currentIndent() + fmt.Sprintf("%s.%s = %s",
		g.generateExpression(stmt.Left.Left),
		stmt.Left.Field.Value,
		g.generateExpression(stmt.Right),
	)
}

func (g *Generator) generateExpression(expr ast.Expression) string {
	switch expr := expr.(type) {
	case *ast.Identifier:
		return expr.Value

	case *ast.IntegerLiteral:
		return expr.TokenLiteral()

	case *ast.FloatLiteral:
		return expr.TokenLiteral()

	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, expr.Value)

	case *ast.BooleanLiteral:
		return fmt.Sprintf("%t", expr.Value)

	case *ast.BinaryExpression:
		return fmt.Sprintf("%s %s %s",
			g.generateExpression(expr.Left),
			expr.Operator.Literal,
			g.generateExpression(expr.Right),
		)

	case *ast.FunctionCall:
		return g.generateFunctionCall(expr)

	case *ast.StructLiteral:
		return g.generateStructLiteral(expr)

	case *ast.StructFieldAccess:
		return fmt.Sprintf("%s.%s",
			g.generateExpression(expr.Left),
			expr.Field.Value,
		)

	case *ast.StructFieldAssignment:
		return g.generateStructFieldAssignment(expr)

	case *ast.PrefixExpression:
		return fmt.Sprintf("%s%s",
			expr.Operator.Literal,
			g.generateExpression(expr.Right),
		)

	case *ast.InfixExpression:
		return fmt.Sprintf("%s %s %s",
			g.generateExpression(expr.Left),
			expr.Operator.Literal,
			g.generateExpression(expr.Right),
		)

	case *ast.IndexExpression:
		return fmt.Sprintf("%s[%s]",
			g.generateExpression(expr.Left),
			g.generateExpression(expr.Index),
		)

	case *ast.ArrayLiteral:
		return g.generateArrayLiteral(expr)

	case *ast.ListLiteral:
		return g.generateListLiteral(expr)

	default:
		return g.currentIndent() + "// Not Supported expression"
	}
}

func (g *Generator) generateExpressions(exprs []ast.Expression) string {
	var out []string
	for _, expr := range exprs {
		out = append(out, g.generateExpression(expr))
	}
	return strings.Join(out, ", ")
}

func (g *Generator) generateFunctionCall(expr *ast.FunctionCall) string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s(", expr.Function))

	args := []string{}
	for _, arg := range expr.Arguments {
		args = append(args, g.generateExpression(arg))
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (g *Generator) generateStructLiteral(expr *ast.StructLiteral) string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("%s {\n", expr.StructName.Value))
	g.indent()
	fields := []string{}
	for name, value := range expr.Fields {
		fieldStr := fmt.Sprintf("%s%s: %s,", g.currentIndent(), name, g.generateExpression(value))
		fields = append(fields, fieldStr)
	}
	out.WriteString(strings.Join(fields, "\n"))
	g.unindent()
	out.WriteString("\n" + g.currentIndent() + "}")

	return out.String()
}

func (g *Generator) generateBlockStatement(stmt *ast.BlockStatement) string {
	var out bytes.Buffer

	out.WriteString("{\n")
	g.indent()
	for i, s := range stmt.Statements {
		out.WriteString(g.generateStatement(s))
		if i < len(stmt.Statements)-1 {
			out.WriteString("\n")
		}
	}
	g.unindent()
	out.WriteString("\n" + g.currentIndent() + "}")

	return out.String()
}

func (g *Generator) generateArrayLiteral(expr *ast.ArrayLiteral) string {
	var out bytes.Buffer

	out.WriteString("{")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, g.generateExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (g *Generator) generateListLiteral(expr *ast.ListLiteral) string {
	var out bytes.Buffer

	out.WriteString("{")
	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, g.generateExpression(el))
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (g *Generator) generateListDeclaration(stmt *ast.ListDeclaration) string {
	var out bytes.Buffer

	out.WriteString(g.currentIndent() + fmt.Sprintf("let %s = ", stmt.Name.Value))
	out.WriteString(g.generateListLiteral(stmt.Value))
	out.WriteString(";")

	return out.String()
}

func (g *Generator) generateForStatement(stmt *ast.ForStatement) string {
	var out bytes.Buffer

	out.WriteString(g.currentIndent() + "for ")

	if stmt.Init != nil {
		switch init := stmt.Init.(type) {
		case *ast.LetStatement:
			out.WriteString(fmt.Sprintf("%s = %s", init.Name.Value, g.generateExpression(init.Value)))
		case *ast.VariableDeclaration:
			out.WriteString(fmt.Sprintf("%s %s = %s", init.Type.Literal, init.Name.Token.Literal, g.generateExpression(init.Value)))
		default:
			out.WriteString(g.generateStatement(stmt.Init))
		}
		out.WriteString(";")
	} else {
		out.WriteString(";")
	}

	if stmt.Condition != nil {
		out.WriteString(fmt.Sprintf(" %s;", g.generateExpression(stmt.Condition)))
	} else {
		out.WriteString(";")
	}

	if stmt.Post != nil {
		switch post := stmt.Post.(type) {
		case *ast.LetStatement:
			out.WriteString(fmt.Sprintf("%s = %s", post.Name.Value, g.generateExpression(post.Value)))
		case *ast.VariableDeclaration:
			out.WriteString(fmt.Sprintf("%s = %s", post.Name.Token.Literal, g.generateExpression(post.Value)))
		default:
			out.WriteString(g.generateStatement(stmt.Post))
		}
	}

	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))
	return out.String()
}
