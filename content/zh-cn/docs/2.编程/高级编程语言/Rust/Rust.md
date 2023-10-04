---
title: Rust
---

# 概述

> 参考：
> 
> - [GitHub 项目，rust-lang/rust](https://github.com/rust-lang/rust)
> - [官网](https://www.rust-lang.org/)、
> - [官方文档](https://doc.rust-lang.org/book/)

# 学习资料

[B 站-软件工艺师，Rust 编程语言入门教程（Rust 语言/Rust 权威指南配套）](https://www.bilibili.com/video/BV1hp4y1k7SV/)

# Hello World

代码：`hello_world.rs`

```rust
fn main() {
    println!("Hello, world!");
}
```

```bash
$ rustc main.rs
$ ./main
Hello, world!
```

# Rust 语言关键字

> 参考：
> 
> - [官方文档](https://doc.rust-lang.org/book/appendix-01-keywords.html)

- `as` - perform primitive casting, disambiguate the specific trait containing an item, or rename items in `use` and `extern crate` statements
- `async` - return a `Future` instead of blocking the current thread
- `await` - suspend execution until the result of a `Future` is ready
- `break` - exit a loop immediately
- `const` - define constant items or constant raw pointers
- `continue` - continue to the next loop iteration
- `crate` - link an external crate or a macro variable representing the crate in which the macro is defined
- `dyn` - dynamic dispatch to a trait object
- `else` - fallback for `if` and `if let` control flow constructs
- `enum` - define an enumeration
- `extern` - link an external crate, function, or variable
- `false` - Boolean false literal
- `fn` - define a function or the function pointer type
- `for` - loop over items from an iterator, implement a trait, or specify a higher-ranked lifetime
- `if` - branch based on the result of a conditional expression
- `impl` - implement inherent or trait functionality
- `in` - part of `for` loop syntax
- `let` - bind a variable
- `loop` - loop unconditionally
- `match` - match a value to patterns
- `mod` - define a module
- `move` - make a closure take ownership of all its captures
- `mut` - denote mutability in references, raw pointers, or pattern bindings
- `pub` - denote public visibility in struct fields, `impl` blocks, or modules
- `ref` - bind by reference
- `return` - return from function
- `Self` - a type alias for the type we are defining or implementing
- `self` - method subject or current module
- `static` - global variable or lifetime lasting the entire program execution
- `struct` - define a structure
- `super` - parent module of the current module
- `trait` - define a trait
- `true` - Boolean true literal
- `type` - define a type alias or associated type
- `union` - define a [union](https://doc.rust-lang.org/reference/items/unions.html) and is only a keyword when used in a union declaration
- `unsafe` - denote unsafe code, functions, traits, or implementations
- `use` - bring symbols into scope
- `where` - denote clauses that constrain a type
- `while` - loop conditionally based on the result of an expression

# Rust 语言规范

> 参考：

