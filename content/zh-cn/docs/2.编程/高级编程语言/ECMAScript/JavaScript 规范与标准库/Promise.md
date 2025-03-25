---
title: Promise
linkTitle: Promise
weight: 20
---

# 概述

> 参考：
>
> - [B 站 up-思学堂，5 分钟彻底学会使用 Promise，你真的懂 Promise 吗？](https://www.bilibili.com/video/BV1TG411L7KM)
> - <https://www.runoob.com/w3cnote/javascript-promise-object.html>

Promise 是一个对象，它代表了一个异步操作的最终完成或者失败。

一个 `Promise` 必然处于以下几种状态之一：

- **pending(待定)** # 初始状态，既没有被兑现，也没有被拒绝。
- **fulfilled(已兑现)** # 意味着操作成功完成。
- **rejected(已拒绝)** # 意味着操作失败。

## Promise 的链式调用

[`Promise.prototype.then()`](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Promise/then)、[`Promise.prototype.catch()`](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Promise/catch) 和 [`Promise.prototype.finally()`](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Promise/finally) 方法用于将进一步的操作与已敲定的 Promise 相关联。由于这些方法返回 Promise，因此它们可以被链式调用。

`.then()` 方法最多接受两个参数；第一个参数是 Promise 兑现时的回调函数，第二个参数是 Promise 拒绝时的回调函数。每个 `.then()` 返回一个新生成的 Promise 对象，这个对象可被用于链式调用，例如：

```js
const myPromise = new Promise((resolve, reject) => {
  setTimeout(() => {
    resolve("foo");
  }, 300);
});

myPromise
  .then(handleFulfilledA, handleRejectedA)
  .then(handleFulfilledB, handleRejectedB)
  .then(handleFulfilledC, handleRejectedC);
```

# 方法

Promise.resolve() 方法

Promise.reject() 方法

## Promise 的异常处理

promise.catch
