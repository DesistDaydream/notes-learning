---
title: jq 命令行工具
---

# 概述

> 参考：
> 
> - 官方文档：<https://stedolan.github.io/jq/>

jq 是轻量级且灵活的处理 JSON 数据的 shell 命令行工具

[这里是官方提供的 jq 命令在线测试工具](https://jqplay.org/)，提供原始 JSON 内容，会自动根据 表达式 输出结果。

# jq 用法详解

官方文档：<https://stedolan.github.io/jq/manual/>

jq 程序是一个`过滤器`，接收一个输入，并产生一个输出。

# 基础过滤

官方文档：<https://stedolan.github.io/jq/manual/#Basicfilters>

下面的 jq 用法，都是用下面这个 json 文件作为演示

    {"favorite":{"drink":"water","food":"sushi","game":"WOW & PAL"},"sushiKinds":["sashimi",{"name":"hot"},{"name":"handRoll","rice":"more"},{"name":null}],"arrayBrowser":[{"name":"360","url":"http://www.so.com"},{"name":"bing","url":"http://www.bing.com"}]}

格式化后的内容如下，格式化内容仅作参考对照，因为 jq 命令本身就可以实现格式化的 json 的作用。

```json
{
  "favorite": {
    "drink": "water",
    "food": "sushi",
    "game": "WOW & PAL"
  },
  "sushiKinds": [
    "sashimi",
    {
      "name": "hot"
    },
    {
      "name": "handRoll",
      "rice": "more"
    },
    {
      "name": null
    }
  ],
  "arrayBrowser": [
    {
      "name": "360",
      "url": "http://www.so.com"
    },
    {
      "name": "bing",
      "url": "http://www.bing.com"
    }
  ]
}
```

## `.` 符号

点`.`符号与 go 模板中的点作用一样，表示**当前作用域**的**对象**。对于 jq 来说，所有给 jq 输入的内容，都是当前作用域的对象。比如

```json
~]# cat demo.json
{"favorite":{"drink":"water","food":"sushi","game":"WOW & PAL"},"sushiKinds":["sashimi",{"name":"hot"},{"name":"handRoll","rice":"more"},{"name":null}],"arrayBrowser":[{"name":"360","url":"http://www.so.com"},{"name":"bing","url":"http://www.bing.com"}]}
~]# cat demo.json | jq .
{
  "favorite": {
    "drink": "water",
    "food": "sushi",
    "game": "WOW & PAL"
  },
  "sushiKinds": [
    "sashimi",
    {
      "name": "hot"
    },
    {
      "name": "handRoll",
      "rice": "more"
    },
    {
      "name": null
    }
  ],
  "arrayBrowser": [
    {
      "name": "360",
      "url": "http://www.so.com"
    },
    {
      "name": "bing",
      "url": "http://www.bing.com"
    }
  ]
}
```

## 获取 map 的值

给定 map 的名称，获取其值。`.foo.bar`与`.foo|.bar`作用相同。如果 map 名称中包含特殊字符或以数字开头，则需要适用双引号将其括起来，例如

```bash
~]# cat demo.json | jq '.favorite.food'
"sushi"
```

## 获取 array 的值

```json
~]# cat demo.json | jq .arrayBrowser
[
  {
    "name": "360",
    "url": "http://www.so.com"
  },
  {
    "name": "bing",
    "url": "http://www.bing.com"
  }
]
~]# cat demo.json | jq .arrayBrowser[]
{
  "name": "360",
  "url": "http://www.so.com"
}
{
  "name": "bing",
  "url": "http://www.bing.com"
}
~]# cat demo.json | jq .arrayBrowser[].name
"360"
"bing"
~]# cat demo.json | jq .arrayBrowser[1]
{
  "name": "bing",
  "url": "http://www.bing.com"
}
```
