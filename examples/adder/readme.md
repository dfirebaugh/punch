# adder

build adder wasm

(ran from the root of this project)

```bash
punch build -o ./examples/adder/adder ./examples/adder/adder.pn
```

This should output a `adder.wat` and an `adder.wasm` file.

`./examples/adder/main.go` is a go program that can load up the `adder.wasm` file and run the functions it contains.

```bash
go run ./examples/adder/main.go
```

```bash
cd ./examples/adders/
node adder.js
```
