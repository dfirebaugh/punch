package parser

type Parser struct {
	Builder Builder
}

func New(builder Builder) Parser {
	return Parser{Builder: builder}
}
func (p Parser) OpenBrace(line int, pos int) {
	p.HandleEvent(OPEN_BRACE, line, pos)
}
func (p Parser) ClosedBrace(line int, pos int) {
	p.HandleEvent(CLOSED_BRACE, line, pos)
}
func (p Parser) OpenParen(line int, pos int) {
	p.HandleEvent(OPEN_PAREN, line, pos)
}
func (p Parser) ClosedParen(line int, pos int) {
	p.HandleEvent(CLOSED_PAREN, line, pos)
}
func (p Parser) OpenAngle(line int, pos int) {
	p.HandleEvent(OPEN_ANGLE, line, pos)
}
func (p Parser) ClosedAngel(line int, pos int) {
	p.HandleEvent(CLOSED_ANGLE, line, pos)
}
func (p Parser) Dot(line int, pos int) {
	p.HandleEvent(DOT, line, pos)
}
func (p Parser) Dash(line int, pos int) {
	p.HandleEvent(DASH, line, pos)
}
func (p Parser) Colon(line int, pos int) {
	p.HandleEvent(COLON, line, pos)
}
func (p Parser) Keyword(name string, line int, pos int) {
	p.HandleEvent(KEYWORD, line, pos)
}
func (p Parser) Name(name string, line int, pos int) {
	print("name:|", name, "|", line, ":", pos, "\n")
	p.HandleEvent(NAME, line, pos)
}
func (p Parser) String(name string, line int, pos int) {
	print("string:|", name, "|", line, ":", pos, "\n")
	p.HandleEvent(STRING, line, pos)
}
func (p Parser) Error(line int, pos int) {
	p.Builder.SyntaxError(line, pos)
}
func (p Parser) HandleEvent(event ParserEvent, line int, pos int) {
	println(event, line, pos)
}

func (p Parser) HandleEventError(event ParserEvent, line int, pos int) {}
