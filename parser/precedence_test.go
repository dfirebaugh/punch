package parser_test

import (
	"testing"

	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/parser"
	"github.com/sirupsen/logrus"
)

// func TestOperatorPrecedenceParsing(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected string
// 	}{
// 		// {
// 		// 	input:    "a + b * c",
// 		// 	expected: "package test\n\n(a + (b * c))\n\n",
// 		// },
// 		// {
// 		// 	input:    "a * b + c",
// 		// 	expected: "package test\n\n((a * b) + c)\n\n",
// 		// },
// 		// {
// 		// 	input:    "a + b - c",
// 		// 	expected: "package test\n\n((a + b) - c)\n\n",
// 		// },
// 		// {
// 		// 	input:    "a * b / c",
// 		// 	expected: "package test\n\n((a * b) / c)\n\n",
// 		// },
// 		// {
// 		// 	input:    "a + b * c - d / e",
// 		// 	expected: "package test\n\n((a + (b * c)) - (d / e))\n\n",
// 		// },
// 		{
// 			input:    "a + (b + c) * d",
// 			expected: "package test\n\n(a + ((b + c) * d))\n\n",
// 		},
// 		// {
// 		// 	input:    "-a * b",
// 		// 	expected: "package test\n\n((-a) * b)\n\n",
// 		// },
// 		// {
// 		// 	input:    "!-a",
// 		// 	expected: "package test\n\n(!(-a))\n\n",
// 		// },
// 		// {
// 		// 	input:    "a + b * c + d / e - f",
// 		// 	expected: "package test\n\n(((a + (b * c)) + (d / e)) - f)\n\n",
// 		// },
// 		// {
// 		// 	input:    "3 + 4; -5 * 5",
// 		// 	expected: "package test\n\n(3 + 4)((-5) * 5)\n\n",
// 		// },
// 		// {
// 		// 	input:    "5 > 4 == 3 < 4",
// 		// 	expected: "package test\n\n((5 > 4) == (3 < 4))\n\n",
// 		// },
// 		// {
// 		// 	input:    "5 < 4 != 3 > 4",
// 		// 	expected: "package test\n\n((5 < 4) != (3 > 4))\n\n",
// 		// },
// 		// {
// 		// 	input:    "3 + 4 * 5 == 3 * 1 + 4 * 5",
// 		// 	expected: "package test\n\n((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))\n\n",
// 		// },
// 	}
//
// 	for i, tt := range tests {
// 		l := lexer.New("test_file", "pkg test\n" + tt.input)
// 		p := parser.New(l)
//
//     logrus.SetLevel(logrus.TraceLevel)
// 		program, err := p.ParseProgram("test_file")
// 		if err != nil {
// 			t.Fatalf("ParseProgram() returned error: %v %d", err, i)
// 		}
//
// 		actual := program.String()
// 		if actual != tt.expected {
// 			t.Errorf("expected=%q, got=%q", tt.expected, actual)
// 		}
// 	}
// }


func TestLogicalPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "true || false",
			expected: "package test\n\n(true || false)\n\n",
		},
		{
			input:    "true && false",
			expected: "package test\n\n(true && false)\n\n",
		},
		{
			input:    "true && false || true",
			expected: "package test\n\n((true && false) || true)\n\n",
		},
		{
			input:    "true || false && true",
			expected: "package test\n\n(true || (false && true))\n\n",
		},
	}

	for i, tt := range tests {
		l := lexer.New("test_file", "pkg test\n" + tt.input)
		p := parser.New(l)

		logrus.SetLevel(logrus.TraceLevel)
		program, err := p.ParseProgram("test_file")
		if err != nil {
			t.Fatalf("ParseProgram() returned error: %v %d", err, i)
		}

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
