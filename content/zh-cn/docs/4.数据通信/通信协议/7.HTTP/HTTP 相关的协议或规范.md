---
title: HTTP 相关的协议或规范
---

# 概述

HTTP 无法单独存在，要想让它生效，必须依赖其他的协议或者规范

## URI 与 URL

> 参考：
> 
> - [Wiki，URI](https://en.wikipedia.org/wiki/Uniform_Resource_Identifier)
> - [Wiki，IRI](https://en.wikipedia.org/wiki/Internationalized_Resource_Identifier)
> - [Wiki，URL](https://en.wikipedia.org/wiki/URL)
> - [Wiki，CleanURL-slug](https://en.wikipedia.org/wiki/Clean_URL)
> - [Wiki，URL encoding](https://en.wikipedia.org/wiki/Percent-encoding)
> - [RFC 3986，Uniform Resource Identifier(URI): Generic Syntax](https://www.rfc-editor.org/rfc/rfc3986.html)
> - [RFC 1738，Uniform Resource Locators (URL)](https://www.rfc-editor.org/rfc/rfc1738)
> - https://www.ruanyifeng.com/blog/2010/02/url_encoding.html

既然 HTTP 的本质是是在两点之间传输超文本，那么这个超文本又该如何表示呢？我们应该如何正确得找到这个超文本呢？所以，人们将超文本描述为 **Resource(资源)**，互联网上如此之多得资源，就需要一个唯一标识符来标识每一个资源。URI 就是这么一个用来标识资源的规范。

**Uniform Resource Identifier(统一资源标识符，简称 URI)** 是 Web 技术使用的唯一标识符。URI 可以用于标识任何东西，包括现实世界中的对象，例如人和地方，概念或信息资源，例如网页和书籍。某些 URI 提供了一种在网络上(在 Internet 上或在另一个专用网络上，例如在计算机文件系统或 Intranet 上)定位和检索信息资源的方法，它们是 **Uniform Resource Locator(统一资源定位符，简称 URL)**。而其他 URI 仅提供一个唯一名称，而没有找到或检索该资源的信息，这类 URI 被称为 **Uniform Resource Name(统一资源名称，简称 URN)**。

> 尽管 URI 仍然是常用术语，但定义 URI 的规范已经被 Internationalized Resource Identifiers(国际化资源标识符，简称 IRI) 的规范所取代。IRI 扩展了 URI 的定义，以便 IRI 可以处理诸如 Kanji(汉字) 之类的字符集，而不是仅限于 ASCII。

**Uniform Resource Locator(统一资源定位符，简称 URL)** 是 URI 的一种子集。
要强制区分 URL 和 URI/IRI 这两种标准化术语是很难的。实际上，两者都使用同一种算法，因此没必要强行区分二者的区别，而且 URL 这个词也更具有人气。所以没有必要强制区分 URI 与 URL。

### URL Syntax(URL 语法)

URL 主要由四个部分组成：协议、主机、端口、路径

**`**Scheme:[//Authority]/Path[?query][#fragment]**`**

- **Scheme://** # URL 方案，即访问协议,指定低层使用的协议(例如：http, https, ftp)
- **Authority** # 分为三个部分 `[UserInfo@]Host[:Port]`
- **UserInfo** # 认证信息。由于安全原因，现在这个年代，都省略这部分，而通过其他方式传递认证信息。
- **Host:PORT** # 服务器 IP 地址或者域名:要访问的端口号
- **/PATH** # 要访问的资源路径。即资源在主机上的存放相对路径。
- **QUERY** # 其实就是参数。用于给动态网页或接口传递参数，可有多个参数，用“&”符号隔开，每个参数的名和值用“=”符号隔开。
- **Fragment** # 片段，主要用于浏览器中。当前页面的其中一段的位置，比如一篇小说有 N 个章节在统一页面显示，开头有目录，点击某一个章节会跳到该页面的某一段，该段的开头就是整个页面的片段，类似于一个位置锚定的作用，该字段即实现这个功能

### URL slug

URL slug 是位于域扩展名之后的 URL 或链接的一部分。

- 它们可用于网站：
  - www.rebrandly.com/links # 'links' 是 URL Slug。
- 或者它们可以用于您的自定义短链接：
  - rebrandly.rocks/content-curation # 'content-curation' 是 URL slug。

### URL Encoding(URL 编码)

通常来说，URL 只能使用英文字母、阿拉伯数字和某些标点符号。那么如果在 URL 中想使用其他字符，比如中文或某些特殊字符怎么办呢？
上述问题就是 URL 编码的由来。在初期，并没有 RFC 规定具体如何进行编码，而是由应用程序自行处理。这就导致 URL 编码称为一个混乱的领域。
在浏览器中，如果输入 `http://www.google.com/你好`，则会被编码为 `http://www.google.com/%E4%BD%A0%E5%A5%BD`。这里面的规则非常简单：

- 你好 两个汉子的 UTF-8 编码结果为 E4 BD A0 E5 A5 BD，每个字节前面加上个 `%`，就得到了 URL 编码。

### EXAMPLE

```bash
          userinfo       host      port
          ┌──┴───┐ ┌──────┴──────┐ ┌┴┐
  https://john.doe@www.example.com:123/forum/questions/?tag=networking&order=newest#top
  └─┬─┘   └───────────┬──────────────┘└───────┬───────┘ └───────────┬─────────────┘ └┬┘
  scheme          authority                  path                 query           fragment
```

- `unix:///run/containerd/containerd.sock` # 这也是 URI 的一种，这就是不同于网络定位符的地方，没有域名之类的东西。Scheme 后面直接接的是 PATH

## TCP/IP

详见：[TCP/IP 相关文章](https://www.yuque.com/go/doc/33218376)

## DNS

详见：[DNS 相关文章](https://www.yuque.com/go/doc/33218346)
