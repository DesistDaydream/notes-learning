---
title: Rust 环境安装与使用
created: 2026-08-06T09:44
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，其他 Rust 安装方式 - 独立安装程序](https://forge.rust-lang.org/infra/other-installation-methods.html#standalone-installers)
> - [官方文档，历史版本](https://forge.rust-lang.org/infra/archive-stable-version-installers.html)

# 安装 Rust

gnu vs musl

|       | gnu                    | musl                                                      |
| ----- | ---------------------- | --------------------------------------------------------- |
| C 库   | 链接 glibc（动态链接）         | 链接 musl libc（可静态链接）                                       |
| 生态兼容性 | 最好，绝大多数 crate 无问题      | 部分依赖 glibc 特性或调用 C 库的 crate 可能有兼容性问题（如某些 DNS 解析、动态加载相关行为） |
| 产物    | 依赖系统的 glibc            | 可以生成完全静态、无外部依赖的二进制                                        |
| 典型场景  | 标准 Linux 发行版日常开发       | Alpine Linux、极简 Docker 镜像、需要单文件可移植二进制分发                   |
| 性能    | 一般更优（尤其多线程下 malloc 实现） | 静态链接方便部署，但某些高并发场景 malloc 性能不如 glibc                       |

## Linux 安装 - 使用 rustup

下载 rustup

```bash
export RustupVersion=1.29.0
export TargetTuple=x86_64-unknown-linux-gnu
# export TargetTuple=x86_64-unknown-linux-musl
wget https://static.rust-lang.org/rustup/archive/${RustupVersion}/${TargetTuple}/rustup-init
```

> [!Note] **target-tuple(目标三元组)**，格式一般是：`CPU 架构-厂商-操作系统-ABI`，用来表示某个二进制文件适用的平台

使用 rustup 安装

```bash
# sudo mkdir -p /usr/local/rust
# sudo chown $(whoami):$(whoami) /usr/local/rust
# export RUSTUP_HOME=/usr/local/rust/rustup
# export CARGO_HOME=/usr/local/rust/cargo

export RUSTUP_DIST_SERVER="https://rsproxy.cn"
export RUSTUP_UPDATE_ROOT="https://rsproxy.cn/rustup"
export RustVersion=1.97.1
./rustup-init -y --no-modify-path --default-toolchain ${RustVersion} --profile minimal
```

配置环境变量

```bash
sudo tee /etc/profile.d/rust.sh > /dev/null <<-"EOF"
# Rustup 镜像
export RUSTUP_DIST_SERVER="https://rsproxy.cn"
export RUSTUP_UPDATE_ROOT="https://rsproxy.cn/rustup"
# export RUSTUP_HOME=/usr/local/rust/rustup
# export CARGO_HOME=/usr/local/rust/cargo
# export PATH=${CARGO_HOME}/bin:${PATH}
export PATH=~/.cargo/bin:${PATH}
EOF

source /etc/profile.d/rust.sh
```

安装 IDE 中必备的组件

```bash
rustup component add rust-src
rustup component add rustfmt
rustup component add clippy
```

设置 crates.io 的镜像

```bash
tee ~/.cargo/config.toml > /dev/null <<-"EOF"
[source.crates-io]
replace-with = 'rsproxy-sparse'
[source.rsproxy]
registry = "https://rsproxy.cn/crates.io-index"
[source.rsproxy-sparse]
registry = "sparse+https://rsproxy.cn/index/"
[registries.rsproxy]
index = "https://rsproxy.cn/crates.io-index"
[net]
git-fetch-with-cli = true
EOF
```

## Linux 安装 - 手工

```bash
export RustVersion=1.97.1
export TargetTuple=x86_64-unknown-linux-gnu
wget https://static.rust-lang.org/dist/rust-${RustVersion}-${TargetTuple}.tar.xz

tar -xvf rust-${RustVersion}-${TargetTuple}.tar.xz
cd rust-${RustVersion}-${TargetTuple}

sudo mkdir -p /usr/local/rust
sudo chown $(whoami):$(whoami) /usr/local/rust
export RUSTUP_HOME=/usr/local/rust/rustup
export CARGO_HOME=/usr/local/rust/cargo
sudo ./install.sh --prefix=/usr/local/rust --without=rust-docs,rust-docs-json-preview
```

> [!Attention] Rust 为什么要把[文档](https://doc.rust-lang.org/std/)（有好多小文件）打包到离线安装包里？安装的时候还会占用大量空间
> 我使用 --without=rust-docs,rust-docs-json-preview 跳过了文档的安装，看看后续有什么不好的影响。可以在安装前，使用 `./install.sh --list-components` 命令查看一下都会安装哪些组件。

配置环境变量

```bash
sudo tee /etc/profile.d/rust.sh > /dev/null <<-"EOF"
export RUSTUP_HOME=/usr/local/rust/rustup
export CARGO_HOME=/usr/local/rust/cargo
export PATH=${CARGO_HOME}/bin:${PATH}
EOF

source /etc/profile.d/rust.sh
```

> [!TODO] 还没装完，unknown-linux-gnu 是动态编译，系统环境里的 glibc 相关文件变了就没法运行二进制文件了。需要再手动安装 musl 与 gnu 互相配合。

> [!TODO] 还没有配置包管理的国内镜像，待验证。好像可以用 https://rsproxy.cn/ ？

## Windows 安装


# 初始化项目

Rust 项目通常由 cargo 工具管理

```bash
cargo init
```

# Rust 编译


# Rust 关联文件与配置

> [!Tip] 这部分笔记中的目录，如果记录了环境变量，则 Rust 生态相关程序可以通过读取环境变量来改变其工作使用的目录。环境变量为空时使用默认目录，环境变量不为空时。

**~/.rustup/** # 存放 rustup 管理的工具链和配置

- 环境变量: `RUSTUP_HOME`

**~/.cargo/** # cargo 程序使用的目录。包括 依赖库、etc.

- 环境变量: `CARGO_HOME`
- 实际要放进 PATH 的那批命令：`bin/` 下的 `rustup`、`rustc`、`cargo`、`rustfmt`、`clippy-driver` 等（这些其实都是指向 rustup 自身二进制的硬链接/代理），还有 `registry/`（crates.io 索引和下载的包缓存）、`git/`（git 依赖缓存）、`env`（PATH 追加脚本）、可选的 `config.toml`

# 最佳实践

# Rust 库

[crates.io](https://crates.io/)

Rust 包或库称为 **crate**。

# Rust 工具

- rustup # Rust 工具链管理工具
- rustc # 编译器
- cargo # 包管理工具

## rustup

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

## cargo

> 参考：
>
> - [GitHub 项目，rust-lang/cargo](https://github.com/rust-lang/cargo)

cargo 是 Rust 的 package(包) 管理器

### cargo 关联文件与配置

**${CARGO_HOME}/** # cargo 工作目录。环境变量为空时，默认为: `~/.cargo/`

- **./config** # TODO: 有啥用? 好像可以配置 crates.io 的国内镜像代理。
