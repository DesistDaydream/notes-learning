---
title: AJAX
linkTitle: AJAX
weight: 2
---

# 概述

> 参考：
>
> - [Wiki, Ajax](https://en.wikipedia.org/wiki/Ajax_(programming))

在 Web 应用程序中（比如浏览器），用户可能需要与服务器进行数据交互，但传统的同步方式在发起 HTTP 请求数据时，会让浏览器在接收到响应后刷新整个页面，导致用户体验变差且不利于更多功能的实现。

所以，如何在不刷新整个页面的情况下更新页面的部分内容呢？

**Asynchronous JavaScript and XML(简称 AJAX)** 就是为了解决上述问题而提出的编程概念，也可称为 Web 开发技术、Web 标准。最早，实现 AJAX 技术的是 [XMLHttpRequest](#xmlhttprequest)

# XMLHttpRequest

> 参考：
>
> - [MDN，术语 - XMLHttpRequest](https://developer.mozilla.org/zh-CN/docs/Glossary/XMLHttpRequest)
> - [MDN，Web API - XMLHttpRequest](https://developer.mozilla.org/zh-CN/docs/Web/API/XMLHttpRequest)
> - [Wiki, XMLHttpRequest](https://en.wikipedia.org/wiki/XMLHttpRequest)

**XMLHttpRequest(简称 XHR)** 是一种创建 AJAX 请求的 [JavaScript](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20规范与标准库/JavaScript%20规范与标准库.md) API，通过 XHR 可以方便得让浏览器发起异步请求到服务器以更新 Web 页面的部分内容。

使用 XHR 主要依赖于 **XMLHttpRequest 对象**，该对象下有很多属性和方法用来定义请求、发送请求、处理响应、等等。尽管名称为 `XMLHttpRequest`，但其可以用于获取任何类型的数据，而不仅仅是 XML。它甚至支持 [HTTP](https://developer.mozilla.org/zh-CN/docs/Web/HTTP) 以外的协议（包括 file:// 和 FTP），尽管可能受到更多出于安全等原因的限制。

```js
// 实例化 XMLHttpRequest 对象
let xhr = new XMLHttpRequest()

// 配置请求信息
xhr.open("GET", "https://api.github.com/users/DesistDaydream", true)

// 绑定 onload 事件，以便在执行 xhr.send() 后处理响应
xhr.onload = function () {
  // 输出一些响应信息
  console.log(xhr.responseText)
  console.log(xhr.status)
}

// 使用配置好的信息发起 HTTP 请求
xhr.send()
```

# Fetch

> 参考：
>
> - [MDN，Web API - Fetch API - Fetch 基本概念](https://developer.mozilla.org/zh-CN/docs/Web/API/Fetch_API/Basic_concepts)
> - [阮一峰，Fetch API 教程](https://www.ruanyifeng.com/blog/2020/12/fetch-tutorial.html)

**Fetch** 是一个现代的概念，等同于 [XMLHttpRequest](https://developer.mozilla.org/zh-CN/docs/Web/API/XMLHttpRequest)。它提供了许多与 XMLHttpRequest 相同的功能，但被设计成更具可扩展性和高效性。Fetch 基于 [Promise](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20规范与标准库/Promise.md) 还原生就拥有异步的特性。

Fetch 的代码写起来比 XMLHttpRequest 更简洁。直接使用 `fetch()` 函数即可，该函数接收两个参数第一个是要获取资源的路径(e.g. URL)，第二个参数可以用来配置请求(e.g. 请求体、请求方法、等等)。fetch 返回一个被泛型约束为 [Response](https://developer.mozilla.org/zh-CN/docs/Web/API/Response) 的 [Promise](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20规范与标准库/Promise.md) 对象(i.e.`Promise<Response>`)

```js
fetch('https://api.github.com/users/DesistDaydream')
  .then(resp => resp.json()) // Response.json() 的返回值也是 Promise
  .then(json => console.log(json))
  .catch(err => console.log('Request Failed', err));
```

Promise 可以使用 await 语法改写，使得语义更清晰。

```js
async function getJSON() {
  let url = 'https://api.github.com/users/DesistDaydream';
  try {
    let resp = await fetch(url);
    return await resp.json();
  } catch (error) {
    console.log('Request Failed', error);
  }
}
```

> 要中止未完成的 `fetch()`，甚至 `XMLHttpRequest` 操作，请使用 [`AbortController`](https://developer.mozilla.org/zh-CN/docs/Web/API/AbortController) 和 [`AbortSignal`](https://developer.mozilla.org/zh-CN/docs/Web/API/AbortSignal) 接口。

`fetch` 规范与 `jQuery.ajax()` 有很多不同的地方。
