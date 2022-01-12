package parser

type SyntaxBuilder struct{}

func NewSyntaxBuilder() SyntaxBuilder {
	return SyntaxBuilder{}
}

func (s SyntaxBuilder) SetName(name string)           {}
func (s SyntaxBuilder) SyntaxError(line int, pos int) {}
