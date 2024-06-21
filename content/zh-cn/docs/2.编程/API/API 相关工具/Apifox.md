---
title: Apifox
linkTitle: Apifox
date: 2024-01-16T11:21
weight: 20
---

# 概述

> 参考：
>
> -

内嵌 [PostMan](/docs/2.编程/API/API%20相关工具/PostMan.md) 的 Postman Collection SDK，可以直接使用 pm 对象对请求和响应进行控制。

在 https://apifox.com/help/pre-post-processors-and-scripts/scripts/api-references/pm-reference#pm 可以看到 Apifox 对 pm 对象的描述文档

# 自带的动态变量

> 参考：
>
> - [Postman 官方文档](https://learning.postman.com/docs/writing-scripts/script-references/variables-list/)
> - [ApiFox 官方文档](https://www.apifox.cn/help/app/api-manage/dynamic-variables/)

- {{$timestamp}} # 当前时间戳

程序脚本

> 参考：
>
> - [Postman 官方文档](https://learning.postman.com/docs/writing-scripts/intro-to-scripts/)
> - [ApiFox 官方文档](https://www.apifox.cn/help/app/scripts/)

## Gdas 签名

```javascript
用于 Apifox 进行 Gdas 签名
// 随机数
var nonce = pm.variables.replaceIn('{{$randomPassword}}');
// 随机数反序
var nonceReverse = nonce.split('').reverse().join('');
// 接入渠道标识
var appkey = "wo-obs";
var secretKey = "obs123456";
// 毫秒时间戳
var stimestamp =Date.parse(new Date());
// 组合签名
var signOriginal = secretKey + nonce + stimestamp + nonceReverse;
var cryptoJs = require("crypto-js");
// 使用 sha256 加密签名并转换为字符串
var signature = cryptoJs.SHA256(signOriginal).toString();
// 输出上面生成的变量的值
console.log("标识符：",appkey);
console.log("时间戳：",stimestamp);
console.log("随机数：",nonce);
console.log("随机数反序：",nonceReverse);
console.log("签名：",signature);
pm.variables.set("appkey", appkey);
pm.variables.set("stimestamp", stimestamp);
pm.variables.set("nonce", nonce);
pm.variables.set("signature", signature);
```

调用其他编程语言的脚本

### go 代码调用示例

```go
package main

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "math/rand"
    "strconv"

    "time"
)

type header struct {
    Appkey     string `json:"appkey"`
    Stimestamp string `json:"stimestamp"`
    Nonce      string `json:"nonce"`
    Signature  string `json:"signature"`
}

func main() {
    // 接入渠道标识
    appkey := "wo-obs"
    secretKey := "obs123456"
    // 毫秒时间戳
    stimestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
    // 随机字符串
    const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.NewSource(time.Now().UnixNano()) // 产生随机种子
    var s bytes.Buffer
    for i := 0; i < 20; i++ {
        s.WriteByte(char[rand.Int63()%int64(len(char))])
    }
    nonce := s.String()
    //nonce 逆序
    var bytes []byte = []byte(nonce)
    for i := 0; i < len(nonce)/2; i++ {
        tmp := bytes[len(nonce)-i-1]
        bytes[len(nonce)-i-1] = bytes[i]
        bytes[i] = tmp
    }
    nonceReverse := string(bytes)
    // 签名 secretKey+nonce+stime+nonce的倒序拼接再SHA256加密
    signOriginal := secretKey + nonce + stimestamp + nonceReverse
    // SHA256加密
    h := sha256.New()
    h.Write([]byte(signOriginal))
    signEncrypt := h.Sum(nil)
    signature := hex.EncodeToString(signEncrypt)

    hd := header{appkey, stimestamp, nonce, signature}
    j, _ := json.Marshal(hd)
    fmt.Println(string(j))
}
```

对应的 Apifox 中的 JS 代码

```javascript
try {
    // 执行 go 代码，代码输出的内容作为 test 变量的值
    const header = fox.execute("consoler_proxy_gen_header.go")
    console.log(header)
    // 解析 JSON
    const json = JSON.parse(header)
    // 设置变量 stime 和 enctyptSing。可以在参数中调用这些变量
    pm.environment.set("appkey", json.appkey)
    pm.environment.set("stimestamp", json.stimestamp)
    pm.environment.set("nonce", json.nonce)
    pm.environment.set("signature", json.signature)
    console.log(json)
} catch (e) {
    console.error(e.message)
}
```

# 当响应数据经过编码或加密，如何在 Apifox 中处理？

https://mp.weixin.qq.com/s/xKqvbtLgTnKz1OXaey7kMQ