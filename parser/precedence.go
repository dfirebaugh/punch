package parser

import "github.com/dfirebaugh/punch/token"

const (
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LOGICAL     // && and ||
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	MOD         // %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index], map[key]
	HIGHEST
)

var precedences = map[token.Type]int{
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESSGREATER,
	token.LT_EQUALS: LESSGREATER,
	token.GT:        LESSGREATER,
	token.GT_EQUALS: LESSGREATER,
	token.AND:       LOGICAL,
	token.OR:        LOGICAL,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.MOD:       MOD,
	token.LPAREN:    CALL,
	token.LBRACKET:  INDEX,
}

func (p *parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
