package compiler

import "testing"

func TestRunCompiler(t *testing.T) {
	program, _ := Compile(`pub function addTwo(x, y) {
		return (x + y);
	}`)
	println(program)
}

func TestTwoFunctions(t *testing.T) {
	program, _ := Compile(`
pub fn addTwo(x, y) {
	return (x + y);
}

pub fn addFour(a, b, c, d) {
	return (a + b + c + d);
}
	`)

	println(program)
}
