---
title: Libvirt 客户端库
---

# 概述

> 参考：
> - [官方文档，binding](https://libvirt.org/bindings.html)

Libvirt 直接支持 C 和 C++，并且具有可用于其他语言的绑定：

- **C#**：Arnaud Champion 开发[C# 绑定](https://libvirt.org/csharp.html)。
- **Go**：Daniel Berrange 开发了 [Go 绑定](https://pkg.go.dev/libvirt.org/go/libvirt)。
- **Java**：Daniel Veillard 开发 [Java 绑定](https://libvirt.org/java.html)。
- **OCaml**：Richard Jones 开发 [OCaml 绑定](https://libvirt.org/ocaml/)。
- **Perl**：Daniel Berrange 开发 [Perl 绑定](https://search.cpan.org/dist/Sys-Virt/)。
- **PHP**：Radek Hladik 于 2010 年开始开发 [PHP 绑定](https://libvirt.org/php)。2011 年 2 月，绑定开发已作为 libvirt-php 项目移至 libvirt.org 网站。该项目现在由 Michal Novotny 维护，并且很大程度上基于 Radek 的版本。有关更多信息，包括发布补丁到 libvirt-php 的信息，请参阅[PHP 绑定](https://libvirt.org/php)站点。
- **Python**：Libvirt 的 python 绑定从 1.2.0 版本开始被拆分为一个单独的 [包](https://gitlab.com/libvirt/libvirt-python)，旧版本直接支持 Python 语言。如果您的 libvirt 是作为软件包安装的，而不是由您从源代码编译的，请确保您安装了适当的软件包。这在 RHEL/Fedora 上被命名为 **libvirt-python** ，在 Ubuntu 上被命名为 [**python-libvirt**](https://packages.ubuntu.com/search?keywords=python-libvirt) ，并且在其他人上可能有不同的命名。有关使用信息，请参阅[Python API 绑定](https://libvirt.org/python.html) 页面。
- **Ruby**：Chris Lalancette 开发[Ruby 绑定](https://libvirt.org/ruby/)。

集成 API 模块：

- **D-Bus**：Pavel Hrdina 开发[D-Bus API](https://libvirt.org/dbus.html)。

有关在 **Windows 上使用 libvirt 的信息，** [请参阅 Windows 支持页面](https://libvirt.org/windows.html)。

# Go 库

> 参考：
> 
> - [GitHub 项目，libvirt/libvirt-go-module](https://github.com/libvirt/libvirt-go-module)
> - [官方文档，go-libvirt](https://libvirt.org/go/libvirt.html)

Go 语言的 `libvirt.org/go/libvrit` 包可以提供来自 OS 原生的 Libvirt API 的 CGO 绑定。该软件包替换了过时的 libvirt.org/libvirt-go 软件包，以便切换到使用 semver 和 Go 模块。除了更改的导入路径和版本控制方案之外，该 API 与旧包完全兼容。

一般来说，Go 表示是从本机 API 概念到 Go 的直接 1-1 映射，因此本机 API 文档应该作为大多数行为的参考。

有关 Go 特定行为的详细信息，请参阅 [Go Pagkage 中的文档](https://pkg.go.dev/libvirt.org/go/libvirt)。

# Python 库

> 参考：
> - [libvirt 官方文档，使用 Python 开发 Libvirt 应用程序指南](https://libvirt.org/docs/libvirt-appdev-guide-python/en-US/html/index.html)

<https://blog.51cto.com/u_10616534/1878609>
<https://cloud.tencent.com/developer/article/1603833>
