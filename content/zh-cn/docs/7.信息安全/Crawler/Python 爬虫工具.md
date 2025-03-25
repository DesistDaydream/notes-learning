---
title: Python 爬虫工具
linkTitle: Python 爬虫工具
weight: 20
---

# 概述

> 参考：
>
> - [2022 Python Web 爬虫武器库盘点](https://zhuanlan.zhihu.com/p/461875098)

[curl_cffi: 支持原生模拟浏览器 TLS/JA3 指纹的 Python 库](https://zhuanlan.zhihu.com/p/601474166)

# HTTP 请求

## Requests

Requests 是从 Python 2 就开始流行的 HTTP 请求库，至今已有十一年历史。当初 Python 内置的 http 库过于底层，用起来非常麻烦，Requests 的出现像救世主一样，slogan 就是“http for humans”，标版自己简单好用人性化。Requests 也确实对得起自己的 solgan，至少提供了以下便捷功能：

1. 对常用的 HTTP 请求方法进行封装，可以非常方便的设置 url 参数，表单参数。
2. 可以自动探测响应内容，自动对响应进行解码和解压缩。
3. 自动维护长连接和连接池。
4. 多文件上传。
5. 自动处理分块响应。
6. 会话和 cookie 持久化。
7. 支持 HTTP/HTTPS 代理。

特别是会话维持、 cookie 持久化、设置代理，是写爬虫必须的。有段时间，知乎上 Python 话题非常火，Python 话题下最火的子话题就是 “Python 爬虫”，而 “Python 爬虫”最火的子话题就是 Requests。

但随着 Python 3.6、3.7、3.8、3.9、3.10 的发布，Python 社区对协程、Type Hint 的支持越来越成熟，同时有部分网站开始使用 HTTP/2 协议，Requests 对这些都不支持。从 2011 年开始，Requests 就显得不那么现代了。

另外，Requests 满足不了特定场景的需求，只能依靠第三方包 requests_toolbelt 的帮助。本人曾经写爬虫时，需要将爬取到的图片上传到另一个对象存储服务，图片爬虫下来后是二进制的内容，而 Requests 不支持直接将二进制内容当成文件上传，只能将图片内容先保存到磁盘文件再用 `open` 读取，将 `open` 返回的句柄对象传递给 Requests。最后不得不使用另一个包 requests_toolbelt 。

> 说一个新手不知道冷知识，Requests 最初的 logo 是一只海龟。

### 优点

1. 在 Python 2 时代提供了简单好用、人性化的接口。
2. 封装了常见的 HTTP 任务。
3. 用户基数大，能踩的坑都填上了，文档丰富。

### 缺点

1. 不支持协程，没有异步接口。
2. 不支持 Type Hint。（Type Hint 接口可以增加可读性，减少查看文档的次数）。
3. 只支持 HTTP/1.1，不支持 HTTP/2。
4. 不支持一些特定场景，需要 requests_toolbelt 作为辅助。

简单来说，Requests 的缺点是生的太早，没跟着社区一起进步。

### 参考

[Requests: HTTP for Humans™ (python-requests.org)](https://link.zhihu.com/?target=https%3A//docs.python-requests.org/en/latest/)

## HTTPX

HTTPX 的 slogan 是 “Python 的下一代 HTTP 客户端”，从出生开始就只支持 Python 3.6 及更高版本。使用了 Type Hint，同时支持同步和异步接口，同时支持 HTTP/1.1 和 HTTP/2，还提供了命令行工具，可以在命令行中直接发送 HTTP 请求。HTTPX 站在 Requests 的肩膀上，Requests 支持的功能它都支持，Requests 不支持的功能它也支持，比 Requests 更现代，没有历史包袱。

但 HTTPX 至今还没发布 1.0.0 版本，截至 2022-1-23 日，最新版本为 0.21.3。但本人在 pypi 上观察到 HTTPX 在 2021 年 9 月 14 发布了一个预发布版本 1.0.0b0，期待在 2022 年正式发布 1.0.0。本人从 2019 年开始从 Requests 转向 HTTPX，以后会一直使用 HTTPX。

### 优点

- 支持异步（协程）接口。
- 支持 HTTP/2。
- 支持 Type Hint。
- 功能更丰富。

简单来说，HTTPX 的优点在于更加现代化。

### 缺点

- 还未发布 1.0.0 版本，不够成熟，需要时间踩坑。
- 只支持 Python 3.6+。（这算不算缺点呢？）

### 参考

[encode/httpx: A next generation HTTP client for Python. (github.com)](https://link.zhihu.com/?target=https%3A//github.com/encode/httpx/)

# 数据解析

## Beautiful Soup

Beautiful Soup 也是从 Python2 时代就开始流行的解析库，用于从 HTML 或 XML 文档中提取数据。Beautiful Soup 会将文档解析成树形文档结构，树中每个节点都是一个 Python 对象，并将节点分为 4 种类型：`Tag` 、 `NavigableString` 、`BeautifulSoup` 、 `Comment` ，并提供了遍历、搜索、修改文档树的接口。支持 CSS 搜索，可以美化 HTML 文档，快速提取文本。Beautiful Soup 还可以指定解析文档时使用的解析器，有的解析器效率高，有的解析器兼容性好，可以根据需要选择。

但是 Beautiful Soup 毕竟出生于 Python 2 时代，没有 Type Hint，再加上接口比较多，每次写代码都要参考官方文档。

### 优点

- 提供了遍历、搜索、修改文档树的结构。
- 可以快速提取所有 HTML 标签中的文本。
- 可以美化 HTML 文档。
- 可以根据需要选择不同的解析器。

虽然后面出现的 Parsel 解析库很好用，但需要一次性提取所有标签文本或编辑 DOM 树结构时还是 Beautiful Soup 好用。

### 缺点

- 接口复杂。
- 没有 Type Hint。

这两个缺点导致写代码时非常依赖官方文档。本人只有在需要提取所有页面文本或编辑 DOM 树时才会用 Beautiful Soup。

### 参考

[Beautiful Soup Documentation — Beautiful Soup 4.4.0 documentation (beautiful-soup-4.readthedocs.io)](https://link.zhihu.com/?target=https%3A//beautiful-soup-4.readthedocs.io/en/latest/)

## Parsel

Parsel 是新一代的 HTML/XML 文档解析库，也是知名爬虫框架 Scrapy 内置的解析器，属于 Scrapy 项目。Parsel 支持 XPath 和 CSS 两种风格的解析器，并集成了正则表达式。

### 优点

- 初步使用了 Type Hint。
- 接口简单，习惯后不依赖文档。
- 提供了 XPath 和 CSS 两个风格的选择器。

### 缺点

- 只能查找和删除文档树中的 HTML 标签，不能插入新标签。编辑能力弱于 Beautiful Soup。

### 参考

[scrapy/parsel: Parsel lets you extract data from XML/HTML documents using XPath or CSS selectors (github.com)](https://link.zhihu.com/?target=https%3A//github.com/scrapy/parsel)

## JSONPath

JSONPath 是查询和解析复杂 JSON 文档的一门语言，JSONPath 的 slogan 是“XPath to JSON”，可以像用 XPath 解析 XML 文档一样用 JSONPath 解析 JSON 文档。JSONPath 由 Stefan Goessner 于 2007 年在[一篇博客](https://link.zhihu.com/?target=https%3A//goessner.net/articles/JsonPath/)中提出，Stefan Goessner 认为 JSON 是 C 系列编程语言中数据的自然表示，所以 JSONPath 表达式的语法也是 C 系列语言风格。例如，访问嵌套的结构可以用 `.` 也可以用 `[]`。

JSONPath 的不完整语法如下：

| JSONPath                                                                                                           | 描述                                 |
| ------------------------------------------------------------------------------------------------------------------ | ------------------------------------ |
| $                                                                                                                  | 根元素或对象                         |
| @                                                                                                                  | 当前对象或元素                       |
| . 或 \[]                                                                                                           | 访问子元素或对象                     |
| ..                                                                                                                 | 递归下降                             |
| \*                                                                                                                 | 通配符，匹配所有元素或对象           |
| \[]                                                                                                                | 下标运算符                           |
| \[,]                                                                                                               | 用名称或数组下标索引将元素提取为一组 |
| \[start![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b088ad0b-52ae-4bbb-84fb-4198566eb0c6/1f51a.svg)step] | 切片操作                             |
| ?()                                                                                                                | 应用过滤表达式                       |
| ()                                                                                                                 | 表达式                               |

有一些网站可以提供 JSONPath 表达式的在线验证。例如，下面从一个表示热搜的复杂 JSON 文档中提取出热搜词列表。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b088ad0b-52ae-4bbb-84fb-4198566eb0c6/v2-052f5331263a0505554b2705887fbe2d_b.jpg)

JSONPath 在线验证

JSONPath 有多种实现，在主流的编程语言中都有对应的库。但由于 JSONPath 没有严格规范的语法，各个实现都有自己的理解，导致同一个 JSONPath 表达式在不同的编程语言中有不同结果。

### 优点

- 可以方便的解析复杂 JSON 文档。
- 语法类 C，相对好记。
- 主流的编程语言都有库支持。
- 代码简洁，占用内存低，高效。

### 缺点

- 没有完整严格的语法规范，导致不同的实现各有差异。

### 参考

[JSONPath - XPath for JSON (goessner.net)](https://link.zhihu.com/?target=https%3A//goessner.net/articles/JsonPath/)

[JSONPath 在线验证](https://link.zhihu.com/?target=https%3A//www.jsonpath.cn/)

## JMESPath

JMESPath 是最近几年最流行的 JSON 文档查询、提取、转换的语言。克服了 JSONPath 的缺点，并提供更加强大功能，例如投影和管道。

JMESPath 的基本表达式包括：标识符、子表达式（用于访问嵌套结构）、下标访问表达式。高级表达式包括：切片、投影、管道、多选、函数。特意没有提供递归下降的功能。

### 优点

JSONPath 的优点它都有，此外还具有 JSONPath 没有的优点：

- 功能更强大。
- 完整而规范的语法，用 ABNF 描述。
- 有一套完整的数据驱动的测试用例，确保不同的 JMESPath 实现库都提供相同的语法支持。
- 针对不同语言都有通过官方认证的实现库。

### 缺点

- 语法相对复杂。

### 参考

[JMESPath — JMESPath](https://link.zhihu.com/?target=https%3A//jmespath.org/)

[RFC 4234 - Augmented BNF for Syntax Specifications: ABNF (ietf.org)](https://link.zhihu.com/?target=https%3A//datatracker.ietf.org/doc/html/rfc4234)

[jmespath vs json-query vs jsonata vs jsonpath vs jsonpath-plus | npm trends](https://link.zhihu.com/?target=https%3A//www.npmtrends.com/jmespath-vs-json-query-vs-jsonata-vs-jsonpath-vs-jsonpath-plus)

# 框架

Scrapy 是爬虫领域至今未被取代的框架，所以本节只介绍了 Scrapy。

## Scrapy

Scrapy 是一款流行了 11 年的 web 爬虫框架，也许是唯一广泛流行的 Python 爬虫框架。Scrapy 框架集成了 HTTP 请求、数据解析、数据持久化、请求调度、并发爬取。

Scrapy 架构包括七个组件，包括：引擎、spider、spider 中间件、调度器、item 管道、下载器、下载器中间件。最基本的爬虫只需要用户继承 `scrapy.Spider`，并实现解析方法 `parse`。

Scrapy 的一般用法是，自定义 spider 类实现请求和解析逻辑，通过 item 管道处理爬取到的数据，通过配置定制爬虫行为，添加 spider 中间件和下载器中间件扩展 Scrapy 的功能。

Scrapy 还提供了 Scrapy-Redis 用于支持分布式爬取。

但 Scrapy 发布时 Python 还没发布标准库 asyncio，也没有 type hint，HTTP/2 协议规范还未发布，所以 Scrapy 并不支持它们。现在只是对 HTTP/2 和 asyncio 有了实验性支持。

### 优点

- 对爬虫任务提供了相对完整的解决方案。
- 提供了一定的扩展性。

### 缺点

- 不支持动态渲染的 Web 页面。
- 不支持 Type Hint。
- 基于 Twisted 实现的异步模型不兼容 Python 标准库 asyncio。
- 对 HTTP/2 的支持不够成熟。

### 参考

[Scrapy | A Fast and Powerful Scraping and Web Crawling Framework](https://link.zhihu.com/?target=https%3A//scrapy.org/)

[Scrapy 2.5 documentation — Scrapy 2.5.1 documentation](https://link.zhihu.com/?target=https%3A//docs.scrapy.org/en/latest/)

[scrapy/scrapy: Scrapy, a fast high-level web crawling & scraping framework for Python. (github.com)](https://link.zhihu.com/?target=https%3A//github.com/scrapy/scrapy)

# 模拟/自动化工具

用自动化测试工具模拟真人爬取网页可以绕过大多数反爬策略，而且不用担心页面动态渲染的问题。

下面介绍的自动化测试工具，原本都是为 Web 自动化测试而生，并不是为爬虫而设计的。本人是从爬虫角度了解它们，所以对它们的介绍肯定不全面，也可能某些地方不准确。

## Selenium

Selenium 是用于支持 web 浏览器自动化的综合项目，包括一系列工具和库的集合。其中最核心的是 WebDriver。我们用 Selenium 模拟爬取就是通过特定于浏览器**驱动**控制浏览器，模拟真人使用浏览器的过程。驱动特定了浏览器，一般由浏览器厂商自己提供，例如 Chrome 浏览器的驱动为 ChromeDriver，Firefox 浏览器的驱动为 GeckoDriver。

爬虫代码通过 WebDrvier 发送指令给驱动，驱动翻译指令后发给浏览器，浏览器的响应会返回给驱动，驱动再返回给 WebDrvier。

Selenium 的目标是提供自动化测试套件，并没为爬取数据做优化，有时候爬取数据需要 hook 请求和返回，而 Selenium 并没有提供这样的功能。本人曾经在做某个项目时非常想要 hook 请求和返回，就去 Selenium 的 Github 仓库搜索相关 issue，发现很早就有人在 issue 里建议加上 hook 请求/响应的功能，但是官方回复说没有这种场景，所以不会提供，提 issue 的人也支支吾吾说不清应用场景，我猜他就是做爬虫的，不好啥意思说出来。

Selenium 3.x 完全没有网络拦截功能，只能通过浏览器插件的形式阻止某个网络请求。

2021 年 10 月 13，经过近三年的迭代，Selenium 终于发布了 4.0.0 版本。Selenium 4.0 引入了一个新特性，支持 Chrome DevTools 协议。利用 Chrome DevTools 协议可以根据 URL 规则阻止请求、mock 请求、记录请求、记录响应，但不支持修改请求或修改响应内容。

另外，Selenium 的固有缺陷导致在驱动浏览器时很容易产生浏览器僵尸进程。当然，僵尸进程有解决办法。

### 优点

- 社区庞大，文档丰富，积累了大量实践经验。
- 至少支持四种主流浏览器，包括：Google Chrome、Microsoft Edge、Mozilla Firefox、
- 针对主流的编程语言都有对应的库。

### 缺点

- 对网络请求/响应的 hook 支持不完善。
- 不支持异步。
- 容易产生浏览器僵尸进程。

### 参考

[Selenium](https://link.zhihu.com/?target=https%3A//www.selenium.dev/)

[https://github.com/SeleniumHQ/Selenium](https://link.zhihu.com/?target=https%3A//github.com/SeleniumHQ/Selenium)

## Puppeteer

Puppeteer 是 Google Chrome 官方团队于 2017 年发布的一个 Node 库，通过 DevTools 协议控制浏览器。能控制的浏览器包括 Google Chrome、Microsoft Edge，不包括 Mozilla Firefox。默认为无头模式，也可以为有头模式。

> Puppeteer 的出现直接导致了另外一款无头浏览器 PhantomJS 于 2018 年宣布停止维护。参见 [PhantomJS - Scriptable Headless Browser](https://link.zhihu.com/?target=https%3A//phantomjs.org/)

Puppeteer 是 Node 库，接口自然都是异步的，async/await 随处可见。

Puppeteer 除了可以模拟用户操作外，还可以拦截请求（修改请求、中止请求、定义返回）。

Puppeteer 有个 Python 迁移版，叫 Pyppeteer，用 Python 的协程语法一比一实现了 Puppeteer 的接口，但由于是个人主导作品，更新、维护以及代码质量都比不上 Puppeteer，不建议使用。

### 优点

- 由 Gooolge Chrome 官方出品，不用担心更新和维护问题。
- 可以拦截请求。

### 缺点

- 官方没有提供相应的 Python 库。第三方提供的 Python 迁移版本质量不高。

### 参考

[puppeteer/puppeteer: Headless Chrome Node.js API (github.com)](https://link.zhihu.com/?target=https%3A//github.com/puppeteer/puppeteer)

[Puppeteer v13.1.1 (pptr.dev)](https://link.zhihu.com/?target=https%3A//pptr.dev/)

[pyppeteer/pyppeteer: Headless chrome/chromium automation library (unofficial port of puppeteer) (github.com)](https://link.zhihu.com/?target=https%3A//github.com/pyppeteer/pyppeteer)

## Playwright

Playwright 是微软官方于 2020 年发布的 Web 自动化测试工具。

微软出手，果然不一般，直接支持了五大浏览器：Google Chrome、Microsoft Edge、Mozilla Firefox、Opera、Safari。支持跨平台，包括：Windows、Linux、macOS。并且还有官方维护的多种编程语言库，包括：TypeScript、JavaScript、Python、.NET、Java。

Playwright 提供的网络拦截功能很适合写爬虫，包括修改请求、中止请求、mock 请求、修改响应。有了这个网络拦截功能，就不需要 mitmproxy 了。

### 优点

- 微软官方出品，不用担心更新和维护问题。
- 兼容多种主流浏览器。
- 官方支持多种编程语言对应的库。
- 可以 hook 请求和响应。
- 更加现代。

### 缺点

- 暂无

参考

[microsoft/playwright: Playwright is a framework for Web Testing and Automation. It allows testing Chromium, Firefox and WebKit with a single API. (github.com)](https://link.zhihu.com/?target=https%3A//github.com/microsoft/playwright)

[Fast and reliable end-to-end testing for modern web apps | Playwright](https://link.zhihu.com/?target=https%3A//playwright.dev/)

# JS 引擎 for Python

写爬虫经常会遇到一些场景，某个请求需要一些参数，而这些参数是经过执行一段复杂的 JavaScript 代码生成的。在用 Python 代码模拟发出该请求时，你可以选择用 Python 重写对于的 js 逻辑，这样会很麻烦，因为你要保证你的 Python 代码逻辑和对应的 JS 代码一模一样，通常还有考虑两种编程语言的异同。如果这段 JS 代码被混淆过，读懂它的逻辑就需要很大的工作量，更便捷的方法是直接将这段 JS 代码拷贝下来，放到一个 JavaScript 环境中执行，再拿到其结果。这就需要一个能在 Python 中使用的 JS 引擎库。

另外在用 Python 脚本分析 JavaScript 代码时，也会需要一个 JavaScript 引擎 for Python。

在 PyV8、PyExecJS 等库都停止维护的情况下，目前唯一比较好的选择是 PyMiniRacer。

## PyMiniRacer

PyMiniRacer 是适用于 Python 的最小的现代嵌入式 V8。虽然维护团队很小，但却是目前的唯一选择。PyMiniRacer 支持最新的 ECMAScript 标准，支持 Assembly，并提供可重用的上下文。

### 优点

- 支持最新的 ECMAScript。
- 内存消耗小。
- 不需要额外安装 JavaScript 环境。
- 提供可重用上下文。

### 缺点

- 维护团队较小，至今未发布 1.0.0 版本。

### 参考

[py-mini-racer · PyPI](https://link.zhihu.com/?target=https%3A//pypi.org/project/py-mini-racer/)

# 虚拟显示器

服务器一般不会配置显式设备，服务器上的 Linux 系统一般也不会安装桌面环境。但如果我们要将一个使用有头浏览器的爬虫部署到服务器上，甚至是 Docker 容器里，该怎么办呢？这时候就可以用虚拟显示器了。目前只有 Xvfb 一个选择。

## Xvfb

Xvfb 是一个 X server，可以在没有显式硬件和物理输入设备的机器上运行，用虚拟内存模拟帧缓冲。

Xvfb 可以设置模拟屏幕的分辨率、像素深度。

我们可以使用 xvfb 启动需要界面的爬虫，或者使用 python 的 xvfbwrapper 在启动需要界面的浏览器之前启动一个虚拟显示器。

### 参考

[xvfb(1): virtual framebuffer X server for X - Linux man page (die.net)](https://link.zhihu.com/?target=https%3A//linux.die.net/man/1/xvfb)

[cgoldberg/xvfbwrapper: Manage headless displays with Xvfb (X virtual framebuffer) (github.com)](https://link.zhihu.com/?target=https%3A//github.com/cgoldberg/xvfbwrapper)

[X](https://link.zhihu.com/?target=https%3A//www.x.org/releases/X11R7.7/doc/man/man7/X.7.xhtml)

# 抓包工具

## Fiddler Classic

Fiddler Classic 是 Windows 平台上一款非常好用的抓包工具。网络上提到的 Fiddler 通常就是指的 Fiddler Classic。

Fiddler Classic 用创建代理的方式抓取浏览器或其他程序的 HTTP/HTTPS 流量。Fiddler Classic 是商业公司推出的网络调试工具之一，该司目前停止了对免费产品 Fiddler Classic 的更新，主力推荐它们的商业收费产品，所以 Fiddler Classic 最后版本为 2020 年 11 月 3 日构建的 v5.0.20204.45441 for .NET 4.6.1。不够 Fiddler Classic 也够用了。如果有跨平台的需求，可以购买该司的 Fiddler Every。

功能特性：

- 可以多种角度查看 HTTP/HTTPS 流量。
- 可以非常方便的构建 HTTP(S) 请求。
- 可以对请求和返回进行断点，并修改继续发送。
- 可以修改请求、重放请求。
- 可以 mock 请求。
- 可以将抓包结果持久化保存。

Fiddler 的抓包原理本质上属于一种“中间人攻击”，有的 app 会对请求进行中间人攻击检测，导致 Fiddler 不能顺利抓到流量。

### 优点

- 免费。
- 可视化操作非常方便。

### 缺点

- 只适用 Windows 平台。

### 参考

[Fiddler Classic | Original Web Capturing Tool for Windows (telerik.com)](https://link.zhihu.com/?target=https%3A//www.telerik.com/fiddler/fiddler-classic)

## 浏览器

浏览器自带的调试工具可以监控网络流量，而且不需要创建代理。但在浏览器的调试页面分析流量不怎么方便，对于复杂的网络交互过程可以将浏览器抓包到的结果保存为 HAR 文件，再用其他工具（例如 Fiddler Classic）分析。

### 优点

- 不用创建代理。
- 不用担心 HTTPS 流量的解密问题。
- 跨平台，免费。

### 缺点

- 分析网络交互过程没有专门的工具方便。

# 个人推荐的最佳组合

- 请求库： HTTPX。
- 框架：Scrapy
- 模拟爬取： Playwright
- JS 引擎 for Python：PyMiniRacer。
- 抓包工具：Fiddler 和浏览器。
