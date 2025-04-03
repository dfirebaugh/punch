package parser

import (
	"fmt"
	"strconv"

	"github.com/dfirebaugh/punch/ast"
	"github.com/dfirebaugh/punch/token"
)

func (p *parser) parseNumberExpression() (ast.Expression, error) {
	var err error
	p.trace("parsing number", p.curToken.Literal, p.peekToken.Literal)
	var n ast.Expression
	if p.curTokenIs(token.NUMBER) {
		n, err = p.parseNumberType()
		if err != nil {
			return nil, err
		}
	}
	if p.curTokenIs(token.FLOAT) {
		n, err = p.parseFloatType()
		if err != nil {
			return nil, err
		}
	}

	if n == nil {
		return nil, p.error("could not parse number")
	}
	// if p.isBinaryOperator(p.peekToken) {
	// 	p.trace("parsing infixed expression", p.curToken.Literal, p.peekToken.Literal)
	// 	p.nextToken()
	// 	return p.parseInfixExpression(n)
	// }

	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.FLOAT) {
		if p.curTokenIs(token.SEMICOLON) {
		}

		p.nextToken()
	}
	return n, nil
}

func (p *parser) parseNumberType() (ast.Expression, error) {
	d, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		return nil, p.error("could not parse number")
	}
	return &ast.IntegerLiteral{
		Token: token.Token{
			Type:     token.I32,
			Literal:  p.curToken.Literal,
			Position: p.curToken.Position,
		},
		Value: int64(d),
	}, nil
}

func (p *parser) parseFloatType() (ast.Expression, error) {
	d, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, p.error("could not parse float")
	}

	return &ast.FloatLiteral{
		Token: p.curToken,
		Value: d,
	}, nil
}

func (p *parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		return nil, p.error(msg)
	}

	lit.Value = value

	return lit, nil
}
