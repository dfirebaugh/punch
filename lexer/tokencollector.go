package lexer

type TokenCollector interface {
	OpenBrace(line int, pos int)
	ClosedBrace(line int, pos int)
	OpenParen(line int, pos int)
	ClosedParen(line int, pos int)
	OpenAngle(line int, pos int)
	ClosedAngel(line int, pos int)
	Dash(line int, pos int)
	Dot(line int, pos int)
	Colon(line int, pos int)
	Keyword(name string, line int, pos int)
	Name(name string, line int, pos int)
	String(name string, line int, pos int)
	Error(line int, pos int)
}
