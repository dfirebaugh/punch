package lexer

type parser struct {
}

func newTestParser() parser {
	return parser{}
}

func (p parser) OpenBrace(line int, pos int)            {}
func (p parser) ClosedBrace(line int, pos int)          {}
func (p parser) OpenParen(line int, pos int)            {}
func (p parser) ClosedParen(line int, pos int)          {}
func (p parser) OpenAngle(line int, pos int)            {}
func (p parser) ClosedAngel(line int, pos int)          {}
func (p parser) Dash(line int, pos int)                 {}
func (p parser) Colon(line int, pos int)                {}
func (p parser) Keyword(name string, line int, pos int) {}
func (p parser) Name(name string, line int, pos int)    {}
func (p parser) String(name string, line int, pos int)  {}
func (p parser) Error(line int, pos int)                {}
func (p parser) Quote(line int, pos int)                {}
func (p parser) Dot(line int, pos int)                  {}
