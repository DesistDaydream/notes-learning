---
title: Python 规范与标准库
weight: 1
---

# 概述

> 参考：
> - [官方文档-3，标准库](https://docs.python.org/3/library/index.html)


- **Python 语言参考**描述了 Python 语言的具体语法和语义
- **Python 标准库则**是与 Python 语言一起发行的一些可选功能，以便人们可以从一开始就轻松得使用 Python 进行编程。

# Python 标识符与关键字

> 参考：
> - [官方文档，参考-2.3.标识符和关键字](https://docs.python.org/3/reference/lexical_analysis.html#identifiers)

and
as
assert
async
await
break
class
continue
def
del
elif
else
except
False
finally
for
from
global
if
import
in
is
lambda
None
nonlocal
not
or
pass
raise
return
True
try
while
with
yield

# Python 语言规范

> 参考：
> 
> - [官方文档-3，参考](https://docs.python.org/3/reference/index.html)

# Python 标准库

> 参考：
> 
> - [官方文档-3，标准库](https://docs.python.org/3/library/index.html)
> - [官方文档-3，标准库参考-内置函数](https://docs.python.org/3/library/functions.html)
> - [官方文档-3，标准库参考-内置常量](https://docs.python.org/3/library/constants.html)

**Python Standard Library(Python 标准库)** 是所有 Python 内置 **Package(包)** 或 **Module(模块)** 的集合，每个 package 都可以实现一类功能。

[pypi.org](https://pypi.org/) 是 Python 的配套网站，可以查找所有可以通过 pip 命令安装内置的或第三方的的 Package

Python 标准库非常庞大，所提供的组件涉及范围十分广泛，正如以下内容目录所显示的。这个库包含了多个内置模块 (以 C 编写)，Python 程序员必须依靠它们来实现系统级功能，例如文件 I/O，此外还有大量以 Python 编写的模块，提供了日常编程中许多问题的标准解决方案。其中有些模块经过专门设计，通过将特定平台功能抽象化为平台中立的 API 来鼓励和加强 Python 程序的可移植性。

## 内置函数

## 内置常量

## 互联网协议和支持

- [webbrowser--- 方便的 Web 浏览器控制工具](https://docs.python.org/zh-cn/3/library/webbrowser.html)
- [wsgiref--- WSGI 工具和参考实现](https://docs.python.org/zh-cn/3/library/wsgiref.html)
- [urllib--- URL 处理模块](https://docs.python.org/zh-cn/3/library/urllib.html)
- [urllib.request--- 用于打开 URL 的可扩展库](https://docs.python.org/zh-cn/3/library/urllib.request.html)
- [urllib.response--- urllib 使用的 Response 类](https://docs.python.org/zh-cn/3/library/urllib.request.html#module-urllib.response)
- [urllib.parse 用于解析 URL](https://docs.python.org/zh-cn/3/library/urllib.parse.html)
- [urllib.error--- urllib.request 引发的异常类](https://docs.python.org/zh-cn/3/library/urllib.error.html)
- [urllib.robotparser--- robots.txt 语法分析程序](https://docs.python.org/zh-cn/3/library/urllib.robotparser.html)
- [http--- HTTP 模块](https://docs.python.org/zh-cn/3/library/http.html)
- [http.client--- HTTP 协议客户端](https://docs.python.org/zh-cn/3/library/http.client.html)
- [ftplib--- FTP 协议客户端](https://docs.python.org/zh-cn/3/library/ftplib.html)
- [poplib--- POP3 协议客户端](https://docs.python.org/zh-cn/3/library/poplib.html)
- [imaplib--- IMAP4 协议客户端](https://docs.python.org/zh-cn/3/library/imaplib.html)
- [smtplib--- SMTP 协议客户端](https://docs.python.org/zh-cn/3/library/smtplib.html)
- [uuid---RFC 4122 定义的 UUID 对象](https://docs.python.org/zh-cn/3/library/uuid.html)
- [socketserver--- 用于网络服务器的框架](https://docs.python.org/zh-cn/3/library/socketserver.html)
- [http.server--- HTTP 服务器](https://docs.python.org/zh-cn/3/library/http.server.html)
- [http.cookies--- HTTP 状态管理](https://docs.python.org/zh-cn/3/library/http.cookies.html)
- [http.cookiejar—— HTTP 客户端的 Cookie 处理](https://docs.python.org/zh-cn/3/library/http.cookiejar.html)
- [xmlrpc--- XMLRPC 服务端与客户端模块](https://docs.python.org/zh-cn/3/library/xmlrpc.html)
- [xmlrpc.client--- XML-RPC 客户端访问](https://docs.python.org/zh-cn/3/library/xmlrpc.client.html)
- [xmlrpc.server--- 基本 XML-RPC 服务器](https://docs.python.org/zh-cn/3/library/xmlrpc.server.html)
- [ipaddress--- IPv4/IPv6 操作库](https://docs.python.org/zh-cn/3/library/ipaddress.html)

## 内嵌模块

### Python 运行时服务

Python 运行时服务

[`sys`](https://docs.python.org/3/library/sys.html) # System-specific parameters and functions