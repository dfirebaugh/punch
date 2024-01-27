package compiler

import "testing"

func TestRunCompiler(t *testing.T) {
	program, _ := Compile(`pub i8 addTwo(x i8, y i8) {
		return (x + y);
	}`)
	println(program)
}

func TestTwoFunctions(t *testing.T) {
	program, _ := Compile(`
pub i8 addTwo(x i8, y i8) {
	return (x + y);
}

pub i8 addFour(a i8, b i8, c i8, d i8) {
	return (a + b + c + d);
}
	`)

	println(program)
}
