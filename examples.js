export const snippets = {
  example: `
pkg main

bool is_eq(i32 a, i32 b) {
  return a == b
}

pub i32 add_two(i32 x, i32 y) {
  println("x =", x, "y =", y)
  println("Hello, World!")
  return x + y
}

println(add_two(2, 5))
  `.trim(),
  multiply: `
pkg main

pub i32 multiply(i32 a, i32 b) {
  return a * b
}

println(multiply(3, 4))
  `.trim(),
  greet: `
pkg main

pub fn greet(str name) {
    println("Hello,", name)
}

greet("World!")
  `.trim(),
  math: `
pkg main

pub fn math_operations() {
    i32 a = 10
    i32 b = 20
    println("Addition: ", a + b)
    println("Subtraction: ", a - b)
    println("Multiplication: ", a * b)
    println("Division: ", b / a)
    println("Modulus: ", b % a)
}

math_operations()
  `.trim(),
  fib: `
pkg main

i32 fibonacci(i32 n) {
    if n == 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    return fibonacci(n - 1) + fibonacci(n - 2)
}

println("Fibonacci(10) is {}", fibonacci(10))
`,
  loop: `
pkg main

pub fn count_to(i32 n) {
    for i32 i = 1; i <= n; i = i + 1 {
        println(i)
    }
}

count_to(5)
  `.trim(),
  return: `
pkg main

i32 square(i32 n) {
    return n * n
}

println("Square of 4 is", square(4))
  `.trim(),
  types: `
pkg main

pub fn log_types() {
    i32 c = 42
    i64 d = 42
    u32 g = 42
    u64 h = 42
    //f32 k = 42.0
    //f64 l = 42.0
    bool m = true
    str n = "hello"

    println("i32:", c)
    println("i64:", d)
    println("u32:", g)
    println("u64:", h)
    //println("f32:", k)
    //println("f64:", l)
    println("bool:", m)
    println("str:", n)
}

log_types()
  `.trim(),
  struct: `
pkg main

struct extra {
  str note
}

struct other {
  str message
  extra extra
}

struct message {
  i32 sender
  i32 receiver
  str body
  other other
}

fn send_message() {
  message msg = message {
    sender: 2,
    receiver: 4,
    body: "hello, world",
    other: other {
      message: "hello",
      extra: extra {
        note: "this is extra info",
      },
    },
  }

  println(msg)
  println(msg.sender)
  println(msg.receiver)
  println(msg.body)
  println(msg.other)
  println(msg.other.message)
  println(msg.other.extra)
  println(msg.other.extra.note)
}

send_message()

  `.trim(),
  list: `
pkg main

fn greet(str name) {
    println("Hello,", name)
}

fn main() {
    []str names = {"Alice", "Bob", "Charlie"}
    append(names, "Alf")
    for i = 0; i < len(names); i = i + 1 {
        greet(names[i])
    }
}

main()
  `.trim(),
};
