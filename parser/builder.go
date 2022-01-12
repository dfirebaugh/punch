package parser

type Builder interface {
	SyntaxError(line int, pos int)
}
