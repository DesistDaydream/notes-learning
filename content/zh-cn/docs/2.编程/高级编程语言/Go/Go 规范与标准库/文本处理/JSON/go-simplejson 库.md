---
title: "go-simplejson 库"
linkTitle: "go-simplejson 库"
date: "2023-06-05T16:06"
weight: 20
---

# 概述

> 参考：
>
> [GitHub 项目，bitly/go-simplejson](https://github.com/bitly/go-simplejson)
> [笔记文档](https://github.com/qq1060656096/go-simplejson)

simplejson 是一个 golang 包， 提供快速和简单的方法来从 JSON 文件中获取值、设置值、删除值。

### 1. Json 解码编码?

```go
// json解码
j, err := simplejson.NewJson([]byte(jsonStr))
// json漂亮编码
s, err := j.EncodeJsonPretty()
// json编码
s, err1 := j.EncodeJSON()
```

### 2. Json 获取值?

```go
// 初始化
j, err := simplejson.NewJson([]byte(jsonStr))
// 获取json object对象值
j.Get(key string).String()
// 获取json 多级object对象值
j.Get(key string).Get(key string).Int()
// 获取json array数组索引值
j.GetIndex(index int).String()
// 获取json 多级array数组索引值
j.GetIndex(index int).GetIndex(index int).Int()
// 获取json 多级数组对象组合键索引值
j.GetIndex(index int).Get(key string).GetIndex(index int).Int()
// 获取数组
// 当我们想以一种简洁的方式对数组值进行迭代时，这很有用:
//  for i, v := range j.Get("results").MustArray() {
//   fmt.Println(i, v)
//  }
// 还可以通过 len(j.MustArray()) 获取数组长度
j.MustArray()
```

### 3. Json 获取类型值?

```go
j, err := simplejson.NewJson([]byte(jsonStr))
// 获取json object对象值
v, err := j.Object()
// 获取json Array数组值
v, err := j.Array()
// 获取json 字符串值
v, err := j.String()
// 获取json bool布尔值
v, err := j.Bool()
// 获取json int整型值
v, err := j.Int()
// 获取json int64整型值
v, err := j.Int64()
// 获取json uint64整型值
v, err := j.Uint64()
// 获取json flat32浮点类型值
v, err := j.Float32()
// 获取json flat64浮点类型值
v, err := j.Float64()
```

### 4. Json 设置值?

```go
j, err := simplejson.NewJson([]byte(jsonStr))
// 设置json字段值, 支持多级键, 支持多级键
j.MustSet(value interface{}, key string)
j.MustSet(value interface{}, index int)
j.MustSet(value interface{}, key1 string|index1 int, key2 string,index2 int, keyN string|indexN int)
j.MustSet(value interface{}, key1 string|index1 int, key2 string,index2 int, keyN string|indexN int).MustSet(value interface{}, key1 string|index1 int, key2 string,index2 int, keyN string|indexN int)
```

### 5. Json 删除值?

```go
j, err := simplejson.NewJson([]byte(jsonStr))
// 删除json值, 支持多级键删除最后一个键, 支持连贯操作
j.Del(key string)
j.Del(index int)
j.Del(key1 string|index1 int, key2 string,index2 int, keyN string|indexN int).Del(key1 string|index1 int, key2 string,index2 int, keyN string|indexN int)
// 删除json字段值, 支持多级键删除最后一个键, 支持连贯操作
j.MustDel(key string)
j.MustDel(index int)
j.MustDel(key1 string|index1 int, key2 string,index2 int, keyN string|indexN int).MustDel(key1 string|index1 int, key2 string,index2 int, keyN string|indexN int)
```

### 6. json 类型对应 golang 类型?

```text
boolean >> bool
number  >> float32,float64,int, int64, uint64
string  >> string
null    >> nil
array   >> []interface{}
object  >> map[string]interface{}
```

## 文档

1\. 字符串解析成 json, 并获取值

```go
j, err := simplejson.NewJson([]byte(jsonStr))
// 获取json object对象值
j.Get(key string).String()
// 获取json 多级object对象值
j.Get(key string).Get(key string).Int()
// 获取json array数组索引值
j.GetArrayIndex(index int).String()
// 获取json 多级array数组索引值
j.GetArrayIndex(index int).GetArrayIndex(index int).Int()
// 获取json 多级数组对象组合键索引值
j.GetArrayIndex(index int).Get(key string).GetArrayIndex(index int).Int()
```

### 示例

```go
package main

import (
 "fmt"
 "github.com/qq1060656096/go-simplejson"
)

func main() {
 jsonStr := `
{
 "uid": 1,
 "name": "tester1",
 "pass": "123456",
 "profile": {
  "age": 18,
  "weight": "75kg",
  "height": "1.71m",
  "mobile": [
   15400000001,
   15400000002
  ]
 }
}
`
 // 字符串解析成json对象
 j, err := simplejson.NewJson([]byte(jsonStr))

 // 简单获取值, 并转换成string类型
 nameValue, err := j.Get("name").String()
 fmt.Println(err)
 fmt.Println(nameValue) // 输出: tester1

 // 连贯操作获取子级键的值, 并转换成int类型
 ageValue, err := j.Get("profile").Get("age").Int()
 fmt.Println(ageValue) // 输出: 18

 // 连贯操作获取子级数组索引值, 并转换成int类型
 mobileIndex2Value, err := j.Get("profile").Get("mobile").GetArrayIndex(1).Int()
 fmt.Println(mobileIndex2Value) // 输出: 15400000002
}
```

2\. 设置 json 对象值

```go
j, err := simplejson.NewJson([]byte(jsonStr))
j.MustSet(value interface{}, key string)
j.MustSet(value interface{}, index int)
j.MustSet(value interface{}, key1 string|index1 index, key2 string,index2 index, keyN string|indexN index)
```

### 示例

```go
package main

import (
 "fmt"
 "github.com/qq1060656096/go-simplejson"
)

func main() {
 jsonStr := `
{
 "mobile": [
  15400000001,
  15400120302
 ],
 "uid": 1
}
`
 j, err := simplejson.NewJson([]byte(jsonStr))
 fmt.Println(err)
 s0, err := j.EncodeJsonPretty()
 fmt.Printf("%s\n", s0)
 // 设置uid值
 j.MustSet(2, "uid")// uid设置为2
 s, err := j.EncodeJsonPretty()
 fmt.Printf("%s\n", s)
 /*
{
  "mobile": [
    15400000001,
    15400120302
  ],
  "uid": 2
}
*/
 // 设置多层级键的值
 j.MustSet(25400000002, "mobile", 1)// 设置mobile索引1的值为25400000002
 s1, err := j.EncodeJsonPretty()
 fmt.Printf("%s\n", s1)
 /*
{
  "mobile": [
    15400000001,
    25400000002
  ],
  "uid": 2
}
*/
}
```
