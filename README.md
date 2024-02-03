# PUNCH 🥊
`punch` is a hobby programming language.  At the moment, `punch` targets wasm.
> I'm mainly working on this as a learning experience.

### Build
Compile a punch program to wasm:
```bash
# build the wasm file
punch -o ./examples/adder/adder ./examples/adder/adder.p
# execute the wasm file
cd ./examples/adder
node adder.js
```

> a `.wat` file and `.ast` file will also be output for debug purposes

#### Functions

```rust
// function declaration
bool is_best(i8 a, i8 b)

// simple function
i8 add(i8 a, i8 b) {
    return a + b
}

// exported function
pub i8 add_two(i8 a, i8 b) {
    return a + b
}

// multiple return types
(i8, bool) add_eq(i8 a, i8 b) {
    return a + b, a == b
}

// no return
main() {
    println("hello world")
}
```

#### Conditions

```rust
if a && b {
    println("abc")
}
```

#### Assignment

```rust
i8 a     = 42
i16 b    = 42
i32 c    = 42
i64 d    = 42
u8 e     = 42
u16 f    = 42
u32 g    = 42
u64 h    = 42
f32 k    = 42.0
f64 l    = 42.0
bool m   = true
str n    = "hello"
```

#### Structs

```rust
struct message {
    i8  sender
    i8 	recipient
    str body
}
message msg = {
    sender: 5,
    recipient: 10,
    body: "hello"
}

println(msg.sender, msg.recipient, msg.body)
```

#### Loops

```go
// traditional for loop
for i := 0; i < 10 ; i = i + 1 {

}

// loop while true
for true {

}

// loop forever
for {

}
```

#### Simple Program

```rust
pkg main

import (
    "fmt"
)

main() {
    fmt.Println("hello, world!")
}
```

#### Status
> work in progress

| Feature | ast | wasm |
| - | - | - |
| function declaration | ✅ | ✅ |
| function calls | ✅ | ✅ |
| function multiple returns | ✅ | ❌ |
| if/else | ✅ | ✅ |
| strings | ✅ | ✅ |
| integers | ✅ | ✅ |
| floats | ✅ |  ✅ |
| structs | ✅ | ✅ |
| struct access | ❌ | ❌ |
| loops | ❌ | ❌ |
| lists | ❌ | ❌ |
| maps | ❌ | ❌ |
| pointers | ❌ | ❌ |
| enums | ❌ | ❌ |
| modules | ❌ | ❌ |
| type inference | ✅ | ✅ |
| interfaces | ❌ | ❌ |

## Reference
- [WebAssembly Text Format (WAT)](https://webassembly.github.io/spec/core/text/index.html)
- [WebAssembly Binary Format Specification](https://webassembly.github.io/spec/core/binary/index.html)
- [WABT - The WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt)
    - [wasm-validate](https://webassembly.github.io/wabt/doc/wasm-validate.1.html) (included in wabt)
- [wat2wasm in browser](https://webassembly.github.io/wabt/demo/wat2wasm/)
