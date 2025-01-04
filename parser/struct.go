package parser

import (
	"fmt"

	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *Parser) parseStructDefinition() *ast.StructDefinition {
	structDef := &ast.StructDefinition{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	p.nextToken()
	structDef.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LBRACE) {
		p.error("expected '{' after struct name")
		return nil
	}
	p.nextToken()

	structDef.Fields = p.parseStructFields()
	if !p.expectPeek(token.RBRACE) {
		p.error("expected '}' after struct fields")
		return nil
	}
	p.nextToken()

	p.definedTypes[structDef.Name.Value] = true
	p.structDefinitions[structDef.Name.Value] = structDef

	return structDef
}

func (p *Parser) parseStructField() *ast.StructField {
	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
	}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	field := &ast.StructField{
		Token: p.peekToken,
		Name:  &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal},
		Type:  p.curToken.Type,
	}
	p.nextToken()
	return field
}

func (p *Parser) parseStructFields() []*ast.StructField {
	fields := []*ast.StructField{}

	for !p.peekTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
		field := p.parseStructField()
		if field == nil {
			p.error("field is nil")
			return nil
		}
		fields = append(fields, field)
	}

	return fields
}

func (p *Parser) parseStructLiteral() ast.Expression {
	structName := p.curToken.Literal
	structDef, ok := p.structDefinitions[structName]
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("undefined struct '%s'", structName))
		return nil
	}

	structLit := &ast.StructLiteral{
		Token:  p.curToken,
		Fields: make(map[string]ast.Expression),
		StructName: &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		},
	}

	p.nextToken()

	if !p.expectCurrentTokenIs(token.LBRACE) {
		p.error("expected '{' after struct name")
		return nil
	}
	p.nextToken()

	for i := 0; !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF); i++ {
		if p.peekTokenIs(token.COLON) {
			fieldName := p.curToken.Literal
			p.nextToken()
			p.nextToken()

			fieldValue := p.parseExpression(LOWEST)
			if fieldValue == nil {
				p.errors = append(p.errors, "expected expression as struct value")
				return nil
			}

			structLit.Fields[fieldName] = fieldValue
		} else {
			if i >= len(structDef.Fields) {
				p.errors = append(p.errors, fmt.Sprintf("too many values in struct literal '%s'", structName))
				return nil
			}

			fieldName := structDef.Fields[i].Name.Value
			fieldValue := p.parseExpression(LOWEST)
			if fieldValue == nil {
				p.errors = append(p.errors, "expected expression as struct value")
				return nil
			}

			structLit.Fields[fieldName] = fieldValue
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}

		p.nextToken()
	}
	p.nextToken()
	return structLit
}

func (p *Parser) parseStructFieldAccess(left ast.Expression) ast.Expression {
	p.nextToken() // consume the dot

	if !p.expectPeek(token.IDENTIFIER) {
		p.error("expected identifier after dot operator")
		return nil
	}

	fieldAccess := &ast.StructFieldAccess{
		Token: p.curToken, // the dot
		Left:  left,
		Field: &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal},
	}

	p.nextToken() // consume the field identifier

	if p.peekTokenIs(token.DOT) {
		return p.parseStructFieldAccess(fieldAccess)
	}

	return fieldAccess
}
