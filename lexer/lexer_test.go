package lexer

import (
	"testing"
)

func testLexerSetup(source string) Lexer {
	parser := newTestParser()
	l := New(parser)
	l.lineNumber = 0
	l.position = 0
	l.source = source
	l.splitLines()

	return l
}

func TestFindName(t *testing.T) {

	l := testLexerSetup(` func( some other name )`)

	l.findWhiteSpace()
	if !l.findName() {
		t.Errorf("should identify %s as a name", l.source)
	}

	if l.position != 5 {
		t.Error("position should be after the name")
	}
}

func TestFindWhiteSpace(t *testing.T) {
	const sourceWithTwoSpaces = `  d`
	l := testLexerSetup(sourceWithTwoSpaces)

	if !l.findWhiteSpace() {
		t.Error("should recognize whitespace")
	}

	if l.position != 2 {
		t.Error("should advance the number of spaces found")
	}
}

func TestFindSingleChar(t *testing.T) {
	l := testLexerSetup(`(`)

	if !l.findSingleCharacterToken() {
		t.Error("should identify a single char")
	}
}
