> work in progress

# PUNCH 🥊
The goal is to build a simple language that compiles directly to WebAssembly.

Also, I'm just kind of toying around with making a language.

The ideal syntax will look similar to below.
```rust
// addTwo is an exported function that adds two ints together and returns the result.
pub i8 addTwo(i8 a, i8 b) {
    return a + b;
}
```

which should output something like the following:
```wat
(module
    (func $addTwo (export "addTwo")(param $x i32)(param $y i32)(result i32)
        (return (i32.add (local.get $x) (local.get $y)))
    )
)
```

### Example

Compile a file to wasm:

```bash
punch -o ./examples/adder/adder ./examples/adder/adder.pn
```

This will output a `adder.wat` file and an `adder.wasm` file.

To execute the `.wasm` file, you can run `go run ./examples/adder/`.
`./examples/adder/main.go` uses wasmtime to load in the wasm file and execute functions that it exports.

## Reference
- [WebAssembly Text Format (WAT)](https://webassembly.github.io/spec/core/text/index.html)
- [WebAssembly Binary Format Specification](https://webassembly.github.io/spec/core/binary/index.html)
- [WABT - The WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt)
    - [wasm-validate](https://webassembly.github.io/wabt/doc/wasm-validate.1.html) (included in wabt)
- [wat2wasm in browser](https://webassembly.github.io/wabt/demo/wat2wasm/)
