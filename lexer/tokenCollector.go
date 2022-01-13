package lexer

type TokenCollector interface {
	Collect(token Token)
}
