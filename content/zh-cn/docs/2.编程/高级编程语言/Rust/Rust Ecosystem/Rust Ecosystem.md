---
title: Rust Ecosystem
created: 2026-08-07T13:00
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Rust - Ecosystem](https://en.wikipedia.org/wiki/Rust_(programming_language)#Ecosystem)

**Rust Ecosystem(Rust 生态)** 包括其 [Compiler](/docs/2.编程/Programming%20tools/Compiler.md)、标准库以及用于软件开发的其他组件。组件安装通常由 Rust 项目开发的 Rust工具链安装程序 管理。

- rustup # Rust 工具链管理工具
- rustc # 编译器
- cargo # 包管理工具

# rustup

> 参考：
>
> - [GitHub 项目，rust-lang/rustup](https://github.com/rust-lang/rustup)

rustup 是 Rust 的 toolchain(工具链) 管理器

> [!Tip] rustup 与 rustup-init 是同一个二进制文件
> 以 rustup-init 作为文件名运行时，会进入"setup 模式"：把自己拷贝到 `$CARGO_HOME/bin/rustup`，再给 rustc、cargo、etc. 命令建一批指向自己的硬链接/代理（这些代理靠"调用时的文件名"来判断该表现成哪个命令），然后修改 PATH，跑一遍 rustup default stable

`rustup`/`rustup-init`/`~/.cargo/bin/` 下的 `rustc`、`cargo`、`rustfmt`、`clippy-driver` 等，在装好之后其实都是**同一个 rustup 二进制的硬链接/拷贝**，它靠自己被调用时的文件名（`argv[0]`）来决定"这次我该表现成哪个工具"，然后转发给 `${RUSTUP_HOME}/toolchains/<当前 toolchain>/bin/` 下真正的可执行文件去执行。这个设计模式更接近 **busybox**（一个二进制，靠文件名分发多种行为）而不是典型的 Go 工具链风格，但"一个核心二进制驱动一大堆看起来独立的命令"这个感觉是对的。

**rustup-init OPTS**

- **-y** # 跳过交互确认
- **--no-modify-path** # 不让 rustup-init 修改用户的 shell rc 文件（我安装的时候自己配置 `/etc/profile.d/` 下的文件管理 `$PATH`）
- **--default-toolchain**(STRING) # 要安装的工具链版本，详见 [概念 - Toolchains](https://rust-lang.github.io/rustup/concepts/toolchains.html)
- **--profile**(STRING) # 安装工具链中的哪些组件。`默认值: default`。可用的值有: minimal, default, complete，详见 [概念 - Profiles](https://rust-lang.github.io/rustup/concepts/profiles.html)。
    - 最小安装用 minimal（包含 rustc/cargo/rust-std）

# cargo

> 参考：
>
> - [GitHub 项目，rust-lang/cargo](https://github.com/rust-lang/cargo)

cargo 是 Rust 的 package(包) 管理器

## cargo 关联文件与配置

**${CARGO_HOME}/** # cargo 工作目录。环境变量为空时，默认为: `~/.cargo/`

- **./config** # TODO: 有啥用? 好像可以配置 crates.io 的国内镜像代理。

## Syntax

# Rust 库

[crates.io](https://crates.io/)

Rust 包或库称为 **crate**。
