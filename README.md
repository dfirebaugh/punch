> work in progress
> Work in progress

# PUNCH 🥊
The goal is to build a simple language that compiles directly to WebAssembly.

Also, I'm just kind of toying around with making a language.

## Goals
* Target WebAssembly Text format (WAT) as an intermediate representation.
* Strict types.
* Enums.
* User-defined types.
* Built-in tools for testing.
* Easy-to-use build tools.
* Module/package management.
* Runtimes for different use cases.
* Functions are private by default and can easily be exported with the `pub` keyword.

### Example

The ideal syntax will look similar to below.

```punch
// addTwo is an exported function that adds two ints together and returns the result.
pub int addTwo(a int, b int) {
    return a + b;
}
```

The above function should build to WAT with the following command:

```
punch -o addTwo.wat ./examples/addTwo.pn
```

output:
```wat
(module
    (func $addTwo (export "addTwo")(param $x i32)(param $y i32)(result i32)
        (return (i32.add (local.get $x) (local.get $y)))
    )
)
```
Ideally, we would also be able to compile to the WebAssembly binary format with the build tool.

## Reference
- [WebAssembly Text Format (WAT)](https://webassembly.github.io/spec/core/text/index.html)
- [WebAssembly Binary Format Specification](https://webassembly.github.io/spec/core/binary/index.html)
- [WABT - The WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt)