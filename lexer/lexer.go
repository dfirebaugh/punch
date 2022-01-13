package lexer

import (
	"strings"
	"text/scanner"
)

type lexer struct {
	Collector TokenCollector
	scanner   scanner.Scanner
}

func New(collector TokenCollector, source string) lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(source))

	lexer := lexer{
		Collector: collector,
		scanner:   s,
	}

	lexer.run()

	return lexer
}

func (l *lexer) run() {
	for tok := l.scanner.Scan(); tok != scanner.EOF; tok = l.scanner.Scan() {
		l.Collector.Collect(Token{Text: l.scanner.TokenText(), Position: l.scanner.Position})
	}
}
