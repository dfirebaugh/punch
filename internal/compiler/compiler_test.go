package compiler

import "testing"

func TestRunCompiler(t *testing.T) {
	program, _, _ := Compile(`pub i8 addTwo(i8 x, i8 y) {
		return (x + y);
	}`)
	println(program)
}

func TestTwoFunctions(t *testing.T) {
	program, _, _ := Compile(`
pub i8 addTwo(i8 x, i8 y) {
	return (x + y);
}

pub i8 addFour(i8 a, i8 b, i8 c, i8 d) {
	return (a + b + c + d);
}
	`)

	println(program)
}
