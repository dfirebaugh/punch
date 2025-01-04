# PUNCH 🥊
`punch` is a hobby programming language.  
> I'm mainly working on this as a learning experience.

[demo playground](https://dfirebaugh.github.io/punch/)

I have some aspirations of working on a backend for this.  To work out some issues with the AST, I added a js code generation step (to easily produce runnable code). I'm not sure if i'll fully commit to that.

### Build
To build you will need [golang installed](https://go.dev/doc/install).

To run code locally, you will need `node` or `bun` installed in your PATH.

```bash
go build ./cmd/punch/

./punch ./examples/simple.pun # output: Hello, World!
```

#### Functions

```rust
// function declaration
bool is_best(i32 a, i32 b)

// simple function
i8 add(i32 a, i32 b) {
    return a + b
}

// exported function
pub i32 add_two(i32 a, i32 b) {
    return a + b
}

// multiple return types
(i32, bool) add_eq(i32 a, i32 b) {
    return a + b, a == b
}

// no return
fn main() {
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
i32 c    = 42
i64 d    = 42
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

fn main() {
    fmt.Println("hello, world!")
}
```

#### Status
> work in progress

| Feature | ast | wasm | js |
| - | - | - | - |
| function declaration | ✅ | ✅ | ✅ |
| function calls | ✅ | ✅ | ✅ |
| function multiple returns | ❌ | ❌ | ❌ |
| if/else | ✅ | ✅ | ✅ |
| strings | ✅ | ✅ | ✅ |
| integers | ✅ | ✅ | ✅ |
| floats | ✅ |  ✅ | ❌ |
| structs | ✅ | ✅ | ✅ |
| struct access | ✅ | ❌ | ✅ |
| loops | ✅ | ❌ | ✅ |
| lists | ❌ | ❌ | ❌ |
| maps | ❌ | ❌ | ❌ |
| pointers | ❌ | ❌ | ❌ |
| enums | ❌ | ❌ | ❌ |
| modules | ❌ | ❌ | ❌ |
| type inference | ❌ | ❌ | ❌ |
| interfaces | ❌ | ❌ | ❌ |

## Reference
- [WebAssembly Text Format (WAT)](https://webassembly.github.io/spec/core/text/index.html)
- [WebAssembly Binary Format Specification](https://webassembly.github.io/spec/core/binary/index.html)
- [WABT - The WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt)
    - [wasm-validate](https://webassembly.github.io/wabt/doc/wasm-validate.1.html) (included in wabt)
- [wat2wasm in browser](https://webassembly.github.io/wabt/demo/wat2wasm/)
