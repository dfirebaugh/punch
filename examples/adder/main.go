package main

import (
	"fmt"
	"os"

	wt "github.com/bytecodealliance/wasmtime-go"
)

func main() {
	data, err := os.ReadFile("./examples/adder/adder.wasm")
	if err != nil {
		fmt.Println("Error reading the WASM file:", err)
		return
	}

	engine := wt.NewEngine()
	store := wt.NewStore(engine)

	module, err := wt.NewModule(store.Engine, data)
	if err != nil {
		fmt.Println("Error loading the module:", err)
		return
	}

	instance, err := wt.NewInstance(store, module, nil)
	if err != nil {
		fmt.Println("Error instantiating the module:", err)
		return
	}

	addTwo := instance.GetFunc(store, "addTwo")

	result, err := addTwo.Call(store, 5, 3)
	if err != nil {
		fmt.Println("Error calling the 'addTwo' function:", err)
		return
	}
	fmt.Println("Result of addTwo(5, 3):", result)

	addFour := instance.GetFunc(store, "addFour")

	result, err = addFour.Call(store, 1, 2, 3, 4)
	if err != nil {
		fmt.Println("Error calling the 'addFour' function:", err)
		return
	}
	fmt.Println("Result of addFour(1, 2, 3, 4):", result)
}
