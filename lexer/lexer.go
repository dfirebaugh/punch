package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	Collector  TokenCollector
	start      int // start position of this item
	position   int // current position in the input
	lineNumber int
	lines      []string
	source     string
}

func New(collector TokenCollector) Lexer {
	return Lexer{
		Collector:  collector,
		start:      0,
		position:   0,
		lineNumber: 0,
	}
}

func (l *Lexer) Lex(source string) {
	l.source = source
	l.splitLines()
	for _, line := range l.lines {
		l.lexLine(line)
		l.lineNumber++
	}
}

func (l *Lexer) splitLines() {
	l.lines = strings.Split(l.source, "\n")
}

func (l *Lexer) resetLinePosition() {
	l.position = 0
}

func (l *Lexer) lexLine(line string) {
	l.resetLinePosition()
	for _, char := range line {
		l.lexToken(char)
	}
}

func (l *Lexer) lexToken(char rune) {
	if !l.findToken(char) {
		l.Collector.Error(l.lineNumber, l.position)
		l.position++
		return
	}
}

func (l *Lexer) findToken(char rune) bool {
	if l.position >= len(l.getLine()) {
		return false
	}

	return l.findString() || l.findWhiteSpace() || l.findName() || l.findSingleCharacterToken()
}

func (l *Lexer) findString() bool {
	if l.getRune() != '"' {
		return false
	}

	i := 0
	for {
		i++
		if rune(l.getLine()[l.position+i]) == '"' {
			break
		}
	}

	str := l.getLine()[l.position : l.position+i+1]
	l.position += len(str)
	l.Collector.String(str, l.lineNumber, l.position)
	return true
}

func (l *Lexer) findName() bool {
	if !unicode.IsLetter(rune(l.getRune())) {
		return false
	}

	// if it is a letter, then we can examine whether or not we have a keyword
	i := 0
	for {
		peek := rune(l.getLine()[l.position+i])
		if !unicode.IsLetter(peek) && !unicode.IsNumber(peek) {
			break
		}
		i++
	}

	word := l.getLine()[l.position : l.position+i]
	l.position += len(word)
	l.Collector.Name(word, l.lineNumber, l.position)
	return true
}

func (l *Lexer) isWhiteSpace() bool {
	return unicode.IsSpace(rune(l.getRune()))
}

// findWhiteSpace - identify a whitespace and how many whitespace characters after
func (l *Lexer) findWhiteSpace() bool {
	if !l.isWhiteSpace() {
		return false
	}

	count := 0
	for {
		if len(l.getLine()) == 0 {
			return false
		}
		if count == len(l.getLine()) {
			break
		}
		if l.position+count >= len(l.getLine()) {
			break
		}
		if !unicode.IsSpace(rune(l.getLine()[l.position+count])) {
			break
		}
		count++
	}

	if count == 0 {
		return false
	}

	l.position += count

	return true
}

func (l *Lexer) findSingleCharacterToken() bool {
	switch l.getRune() {
	case '{':
		l.Collector.OpenBrace(l.lineNumber, l.position)
		l.position++
		return true
	case '}':
		l.position++
		l.Collector.ClosedBrace(l.lineNumber, l.position)
		return true
	case '(':
		l.position++
		l.Collector.OpenParen(l.lineNumber, l.position)
		return true
	case ')':
		l.position++
		l.Collector.ClosedParen(l.lineNumber, l.position)
		return true
	case '.':
		l.position++
		l.Collector.Dot(l.lineNumber, l.position)
		return true
	case ';':
		return true
	default:
		return false
	}
}

// getRune - returns Rune at the current position
func (l *Lexer) getRune() byte {
	return l.getLine()[l.position]
}

func (l *Lexer) getChar() string {
	return string(l.getRune())
}

func (l *Lexer) getLine() string {
	if len(l.lines) == 0 {
		return l.lines[0]
	}
	if l.lineNumber >= len(l.lines) {
		return ""
	}
	return l.lines[l.lineNumber]
}
