---
title: ECMAScript 第三方库
weight: 1
---

# Postman

```javascript
// 将响应体解析为 JSON 格式
respBodyJSON = JSON.parse(responseBody)

respBodyJSON.data.forEach(function (item) {
  console.log(item.title)
})
```

# Axios

> 参考：
>
> - [GitHub 项目，axios/axios](https://github.com/axios/axios)
> - [官网](https://axios-http.com/)

Axios 是用于浏览器和 node.js 的基于 [Promise](/docs/2.编程/高级编程语言/ECMAScript/JavaScript%20规范与标准库/Promise.md) 的 HTTP 客户端，它是基于 XHR 进行的二次封装，传统的 XHR 并没有使用到 Promise
