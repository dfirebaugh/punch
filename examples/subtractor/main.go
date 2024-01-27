package main

import (
	"fmt"
	"os"

	wt "github.com/bytecodealliance/wasmtime-go"
)

func main() {
	data, err := os.ReadFile("./examples/subtractor/subtractor.wasm")
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

	subTwo := instance.GetFunc(store, "subTwo")

	result, err := subTwo.Call(store, 5, 3)
	if err != nil {
		fmt.Println("Error calling the 'subTwo' function:", err)
		return
	}
	fmt.Println("Result of subTwo(5, 3):", result)

	subFour := instance.GetFunc(store, "subFour")

	result, err = subFour.Call(store, 1, 2, 3, 4)
	if err != nil {
		fmt.Println("Error calling the 'subFour' function:", err)
		return
	}
	fmt.Println("Result of subFour(1, 2, 3, 4):", result)
}
