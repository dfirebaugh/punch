package parser

import (
	"github.com/dfirebaugh/punch/ast"
)

func (p *parser) parseIdentifier() (ast.Expression, error) {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.isStructAccess() {
		return p.parseStructFieldAccess(ident)
	}

	if p.isFunctionCall() {
		p.trace("parsing identifier function call expression", p.curToken.Literal, p.peekToken.Literal)
		fnCall, err := p.parseFunctionCall(ident)
		if err != nil {
			return nil, err
		}
		p.trace("parsed identifier function call expression", p.curToken.Literal, p.peekToken.Literal)
		return fnCall, nil
	}

	return ident, nil
}
