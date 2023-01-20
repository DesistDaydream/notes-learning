---
title: ECMAScript 第三方库
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
> - [GitHub 项目，axios/axios](https://github.com/axios/axios)
> - [官网](https://axios-http.com/)

Axios 是用于浏览器和 node.js 的基于 Promise 的 HTTP 客户端

# Element Plus

> 参考：
> - [GitHub 项目，element-plus/element-plus](https://github.com/element-plus/element-plus)
> - [官网](https://element-plus.org/)

Element Plus 是一个基于 Vue3 的 UI 框架。它是 Element UI 基于 Vue3 的重构版本。
