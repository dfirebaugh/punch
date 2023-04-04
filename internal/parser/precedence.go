package parser

import "github.com/dfirebaugh/punch/internal/token"

type precedence int

// precedence order
const (
	_ int = iota
	LOWEST
	COND         // OR or AND
	ASSIGN       // =
	TERNARY      // ? :
	EQUALS       // == or !=
	REGEXP_MATCH // !~ ~=
	LESSGREATER  // > or <
	SUM          // + or -
	PRODUCT      // * or /
	POWER        // **
	MOD          // %
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	// DOTDOT       // ..
	INDEX // array[index], map[key]
	HIGHEST
)

// each token precedence
var precedences = map[token.Type]int{
	token.QUESTION:  TERNARY,
	token.ASSIGN:    ASSIGN,
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESSGREATER,
	token.LT_EQUALS: LESSGREATER,
	token.GT:        LESSGREATER,
	token.GT_EQUALS: LESSGREATER,

	token.PLUS:            SUM,
	token.PLUS_EQUALS:     SUM,
	token.MINUS:           SUM,
	token.MINUS_EQUALS:    SUM,
	token.SLASH:           PRODUCT,
	token.SLASH_EQUALS:    PRODUCT,
	token.ASTERISK:        PRODUCT,
	token.ASTERISK_EQUALS: PRODUCT,
	token.MOD:             MOD,
	token.AND:             COND,
	token.OR:              COND,
	token.LPAREN:          CALL,
	token.DOT:             CALL,
	token.LBRACKET:        INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
