package compiler

import "testing"

func TestRunCompiler(t *testing.T) {
	program := Compile(`pub function addTwo(x, y) {
		return (x + y);
	}`)
	println(program)
}
