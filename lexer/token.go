package lexer

import (
	"fmt"
	"text/scanner"
)

type Token struct {
	Text     string
	Position scanner.Position
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%d: %s", t.Position.Line, t.Position.Offset, t.Text)
}
