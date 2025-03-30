package parser

import (
	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *Parser) parseStructDefinition() (*ast.StructDefinition, error) {
	var err error
	structDef := &ast.StructDefinition{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil, p.error("expected identifier")
	}
	p.nextToken()
	structDef.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LBRACE) {
		return nil, p.error("expected '{' after struct name")
	}
	p.nextToken()

	structDef.Fields, err = p.parseStructFields()
	if err != nil {
		return nil, err
	}
	if !p.expectPeek(token.RBRACE) {
		return nil, p.error("expected '}' after struct fields")
	}
	p.nextToken()

	p.definedTypes[structDef.Name.Value] = true
	p.structDefinitions[structDef.Name.Value] = structDef

	return structDef, nil
}

func (p *Parser) parseStructField() (*ast.StructField, error) {
	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
	}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil, p.error("expected identifier")
	}
	field := &ast.StructField{
		Token: p.peekToken,
		Name:  &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal},
		Type:  p.curToken.Type,
	}
	p.nextToken()
	return field, nil
}

func (p *Parser) parseStructFields() ([]*ast.StructField, error) {
	fields := []*ast.StructField{}

	for !p.peekTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		p.nextToken()
		field, err := p.parseStructField()
		if err != nil {
			return nil, err
		}
		if field == nil {
			return nil, p.error("field is nil")
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func (p *Parser) parseStructLiteral() (ast.Expression, error) {
	structName := p.curToken.Literal
	structDef, ok := p.structDefinitions[structName]
	if !ok {
		return nil, p.errorf("undefined struct '%s'", structName)
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
		return nil, p.error("expected '{' after struct name")
	}
	p.nextToken()

	for i := 0; !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF); i++ {
		if p.peekTokenIs(token.COLON) {
			fieldName := p.curToken.Literal
			p.nextToken()
			p.nextToken()

			fieldValue, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			if fieldValue == nil {
				return nil, p.error("expected expression as struct value")
			}

			structLit.Fields[fieldName] = fieldValue
		} else {
			if i >= len(structDef.Fields) {
				return nil, p.errorf("too many values in struct literal '%s'", structName)
			}

			fieldName := structDef.Fields[i].Name.Value
			fieldValue, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			if fieldValue == nil {
				return nil, p.error("expected expression as struct value")
			}

			structLit.Fields[fieldName] = fieldValue
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}

		p.nextToken()
	}
	p.nextToken()
	return structLit, nil
}

func (p *Parser) parseStructFieldAccess(left ast.Expression) (ast.Expression, error) {
	if !p.expectCurrentTokenIs(token.DOT) {
		p.error("expected a dot, but got: ", p.curToken.Literal)
	}
	p.nextToken() // consume the dot

	if !p.expectPeek(token.IDENTIFIER) {
		return nil, p.error("expected identifier after dot operator")
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
	if p.isBinaryOperator(p.peekToken) {
		p.nextToken()
		return p.parseInfixExpression(fieldAccess)
	}

	return fieldAccess, nil
}

func (p *Parser) parseStructFieldAssignment(left ast.Expression) (ast.Expression, error) {
	tok := p.curToken
	if !p.curTokenIs(token.ASSIGN) {
		p.error("we were expecting an assignment operator, but got: ", p.curToken.Literal)
	}
	p.nextToken()

	right, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, p.error("expected expression after assignment operator")
	}

	fieldAccess, ok := left.(*ast.StructFieldAccess)
	if !ok {
		return nil, p.error("left-hand side of assignment must be a struct field access")
	}

	return &ast.StructFieldAssignment{
		Token: tok,
		Left:  fieldAccess,
		Right: right,
	}, nil
}
