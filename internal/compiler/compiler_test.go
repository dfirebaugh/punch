package compiler

import "testing"

func TestRunCompiler(t *testing.T) {
	c := New(`function helloWorld() {
		return "hello, world!";
	}

	let test = 1 + 1
	helloWorld();
	print(helloWorld());
	`)
	c.Run()
}
