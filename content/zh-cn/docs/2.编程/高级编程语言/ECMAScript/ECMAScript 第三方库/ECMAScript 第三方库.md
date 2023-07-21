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

Axios 是用于浏览器和 node.js 的基于 Promise 的 HTTP 客户端
