package main

import (
	"punch/internal/compiler"
)

func main() {
	// 🥊
	compiler := compiler.New(`function helloWorld() {
	let name = "world";
	return "hello, "+name+"!";
}

let test = 1 + 1;
helloWorld();
print(helloWorld());
		`)
	compiler.Run()
}
