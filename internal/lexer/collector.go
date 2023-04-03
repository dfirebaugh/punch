package lexer

import "github.com/dfirebaugh/punch/internal/token"

type Collector struct {
	result []token.Token
}

func (c *Collector) Collect(t token.Token) {
	c.result = append(c.result, t)
}
func (c *Collector) Result() []token.Token {
	return c.result
}
