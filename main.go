package main

import (
	"punch/compiler"
)

func main() {
	// 🥊
	compiler := compiler.New(`function helloWorld() {
			return "hello, world!";
		}

		helloWorld();
		print(helloWorld());
		`)
	compiler.Run()
}
