---
title: JSON
linkTitle: JSON
date: 2023-06-05T16:13
weight: 1
---

# 概述

> 参考：
>
> - [Go 包，标准库 - encoding/json](https://pkg.go.dev/encoding/json)
> - [Go 官方博客《JSON and Go》](https://blog.golang.org/json)
> - [骏马金龙](https://www.cnblogs.com/f-ck-need-u/p/10080793.html)
> - [在线 JSON 转 Go Struct](https://transform.tools/json-to-go)

Go 语言内置了 encoding 标准库，并使用其中的 json 包来实现 [JSON](/docs/2.编程/无法分类的语言/JSON.md) 的 **Encoding(编码)** 和 **Decoding(解码)** 逻辑，这里面的 JSON 符合 [RFC 7159](https://rfc-editor.org/rfc/rfc7159.html) 标准。

这里面说的编码则是指将 Go 自身可以识别的数据类型编码为 JSON 格式的数据；解码主要是指解码后数据编程 Go 可以理解的数据格式，比如 struct、slice、map 等。

> 说白了，Go 编完的码，是编成人类可读的，自己不认识；解码才是将数据解成 Go 自己认识的。

**JSON类型 与 Go类型 对应关系**

```text
boolean >> bool
number  >> float32,float64,int, int64, uint64
string  >> string
null    >> nil
array   >> []interface{}
object  >> map[string]interface{}
```

使用 json 包，可以轻松地从 Go 程序中读取和写入 JSON 数据。

# Encoding 与 Decoding

**Encoding(编码)** 与 **Decoding(解码)** 是 JSON 数据处理的基本操作

在 json 包中，使用 `Marshal()` 和 `Unmarshal()` 函数来执行最基本的 Encoding 与 Decoding 行为。

> [!Note]
> Marshal: 直译为“编排、整理、排列、序列”，表示整理指定的内容，将内容整理成 json 数据。所以有时候也称此行为叫 **serializable(序列化)**。这种称呼是相对的。在计算机中特指将数据按某种描述格式编排出来。unMarshal 自然是指 Marshal 的逆过程。
>
> 比如在 Web 服务中，我们需要把 go 的 struct 以 JSON 方式表示并在网络间传输，把 go struct 转化成 JSON 的过程就是marshal。

用白话说：

- **Encoding 就是将 struct、slice、array、map 等 转换为 JSON 格式**
- **Decoding 就是将 JSON 格式转换为 struct、slice、array、map。**

这是一个 JSON 结构的数据

```json
{
    "id": 1,
    "content": "hello world",
    "author": {
        "id": 2,
        "name": "userA"
    },
    "published": true,
    "label": [],
    "nextPost": null,
    "comments": [{
            "id": 3,
            "content": "good post1",
            "author": "userB"
        },
        {
            "id": 4,
            "content": "good post2",
            "author": "userC"
        }
    ]
}
```

如果想要让 struct 可以存储上述 JSON 格式数据，那么需要如下定义方式：

```go
type Post struct {
 ID        int64         `json:"id"`
 Content   string        `json:"content"`
 Author    Author        `json:"author"`
 Published bool          `json:"published"`
 Label     []string      `json:"label"`
 NextPost  Post         `json:"nextPost"`
 Comments  []Comment    `json:"comments"`
}

type Author struct {
 ID   int64  `json:"id"`
 Name string `json:"name"`
}

type Comment struct {
 ID      int64  `json:"id"`
 Content string `json:"content"`
 Author  string `json:"author"`
}
```

## Encoding(编码)

**Encoding(编码)** 就是指将其他类型数据封装成 JSON 格式的数据。编码编码，也就是将某些数据编排一下变成另外一种样子。

数据转换时，遵循着一定的规范：

- **只有可以表示为有效 JSON 的数据结构才会被编码：**
  - **struct、slice、array、map 都可以转换成json**
  - **struct 转换成 json的时候，struck 中只有字段首字母大写的属性才会被转换**
  - **map转换的时候，key必须为string**
- **封装的时候，如果是指针，会追踪指针指向的对象进行封装**
- JSON 对象仅支持字符串作为键；要编码Go map 类型，它必须采用以下形式`map[string]T`（`T` json包支持的所有Go类型）。
- Channel、complex、function 类型无法编码。
- 不支持循环数据结构；它们将导致`Marshal`陷入无限循环。
- 指针将被编码为其所指向的值（如果指针为，则为“ null” `nil`）。
- json 包仅访问结构类型（以大写字母开头的结构类型）的导出字段。因此，仅 struct 的导出字段将出现在 JSON 输出中。

在 json 包中，可以使用`Marshal()`或者 `Mashallndent()` 函数来执行 Encoding 行为。

```go
func Marshal(v interface{}) ([]byte, error)
```

### Example

假如现在有一个名为 `Message` 的 Struct，这个结构体表示一条消息中应该具有的属性。比如发送者、消息内容、发送时间，等等。

```go
type Message struct {
    Name string
    Body string
    Time int64
}
```

要想将这个 Struct 中的数据转换为 JSON 格式，只需要使用 `Marshal()` 函数即可

和一个实例 `Message`

```go
// 向结构体中写入数据
m := Message{"DesistDaydream", "Hello", 1294706395881547000}
// 使用 Marshal() 方法，将 m 编码为 b
b, err := json.Marshal(m)
```

`Marshal()` 返回的是一个 `[]byte` 类型，现在变量 b 就存储了一段 `[]byte` 类型的 JSONG 格式数据。可以使用 `string()` 将类型转换为人类可读的字符串类型：

```go
fmt.Println(string(b))
```

输出结果为：

```json
{"Name":"Alice","Body":"Hello","Time":1294706395881547000}
```

注意：

由于转换规范的原因导致 json 格式数据的 key 的首字母都是大写的，如果想要小写的，只需要给 struct 属性添加注释可，比如：

```go
type Message struct {
    Name string `json:"name"`
    Body string `json:"body"`
    Time int64  `json:"time"`
}
```

那么输出结果就是这样的：

```json
{"name":"Alice","body":"Hello","time":1294706395881547000}
```

MarshalIndent() 函数，则是可以在 Encoding 成 JSON 的时候进行美化，将会自动添加前缀和缩进(前缀字符串一般设置为空)

```go
c,err := json.MarshalIndent(Message,"","\t")
if err != nil {
 fmt.Println(nil)
}
fmt.Println(string(c))
```

输出结果为：

```json
{
 "Name": "Alice",
 "Body": "Hello",
 "Time": 1294706395881547000
}
```

## Decoding(解码)

要解码 JSON 数据，我们使用`Unmarshal()`函数。

> Marshal 有整理、排列、序列的含义，表示整理指定的内容，将内容整理成 JSON 数据。那么 Unmarshal 就是 打散 这种含义。有时候也称为 **反序列化。**
> 比如可以这么描述：将 JSON 数据反序列化成指定的数据

```go
func Unmarshal(data []byte, v interface{}) error
```

我们首先必须创建一个存储解码数据的地方

```go
var m Message
```

并调用`json.Unmarshal`，将`[]byte`JSON数据和一个指针传递给它`m`

```go
err := json.Unmarshal(b, &m)
```

如果`b`包含有效的JSON，适合在`m`后电话`err`将`nil`与从数据`b`将被存储在结构`m`，仿佛像一个任务：

```go
m = Message{
    Name: "Alice",
    Body: "Hello",
    Time: 1294706395881547000,
}
```

如何`Unmarshal`识别存储解码数据的字段？对于给定的JSON键`"Foo"`， `Unmarshal`将浏览目标结构的字段以查找(按优先顺序)：

- 标记为的导出字段`"Foo"`（ 有关struct标记的更多信息，请参见Go规范），
- 名为`"Foo"`或的导出字段
- 名为`"FOO"`或`"FoO"`或其他不区分大小写的匹配项的导出字段`"Foo"`。

当JSON数据的结构与Go类型不完全匹配时会发生什么？

```go
b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
var m Message
err := json.Unmarshal(b, &m)
```

`Unmarshal`只会解码在目标类型中可以找到的字段。在这种情况下，将仅填充m的Name字段，而Food字段将被忽略。当您希望从大型JSON Blob中仅选择几个特定字段时，此行为特别有用。这也意味着目标struct中所有未导出的字段都不会受到的影响`Unmarshal`。

但是，如果您事先不知道JSON数据的结构怎么办？

# 如何使用 struct Tag

https://pkg.go.dev/encoding/json#Marshal

json 库可以利用 [Struct](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Map%20AND%20Struct.md) 中的 Tag 控制编码与解码过程中的部分行为。

常见 Tag 关键字

- **json** # 定义 JSON 原始数据中字段的 key 的名称。避免 json 库根据 struct 字段命自动推导的 key 命与原始 JSON 中的 key 命无法对应，导致解析失败。
- **omitempty** # omitempty 标签指定如果字段具有空值（定义为 false、0、nil 指针、nil 接口值以及任何空数组、切片、映射或字符串），则应从编码中省略该字段。
- **-** # 作为一种特殊情况，如果字段标记为 “-”，则始终省略该字段。请注意，仍然可以使用标记 “-，” 生成名称为 “-” 的字段。

# 使用 interface{} 存放通用 JSON 数据

的`interface{}`（空接口）类型描述了具有零种方法的接口。每个Go类型至少实现零个方法，因此满足空接口。

空接口用作常规容器类型：

```go
var i interface{}
i = "a string"
i = 2011
i = 2.777
```

类型断言访问基础的具体类型：

```go
r := i.(float64)
fmt.Println("the circle's area", math.Pi*r*r)
```

或者，如果基础类型未知，则由类型开关确定类型：

```go
switch v := i.(type) {
case int:
    fmt.Println("twice i is", v*2)
case float64:
    fmt.Println("the reciprocal of i is", 1/v)
case string:
    h := len(v) / 2
    fmt.Println("i swapped by halves is", v[h:]+v[:h])
default:
    // i isn't one of the types above
}
```

json包使用`map[string]interface{}`和 `[]interface{}`值来存储任意JSON对象和数组；它将很乐意将任何有效的JSON Blob解组为纯 `interface{}`值。默认的具体Go类型为：

- `bool` 对于JSON布尔值，
- `float64` 对于JSON数字，
- `string` 用于JSON字符串，以及
- `nil` JSON null。

# 解码任意数据

考虑以下存储在变量中的JSON数据`b`：

```go
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
```

在不知道此数据结构的情况下，我们可以使用以下命令将其解码为一个`interface{}`值`Unmarshal`：

```go
var f interface{}
err := json.Unmarshal(b, &f)
```

此时，Go值`f`将是一个映射，其键为字符串，其值本身存储为空接口值：

```go
f = map[string]interface{}{
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{}{
        "Gomez",
        "Morticia",
    },
}
```

要访问此数据，我们可以使用类型断言来访问`f`的底层`map[string]interface{}`：

```go
m := f.(map[string]interface{})
```

然后，我们可以使用 range 语句遍历 map，并使用类型开关将其值作为其具体类型来访问：

```go
for k, v := range m {
    switch vv := v.(type) {
    case string:
        fmt.Println(k, "is string", vv)
    case float64:
        fmt.Println(k, "is float64", vv)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don't know how to handle")
    }
}
```

这样，您可以使用未知的JSON数据，同时仍然享有类型安全的好处。

# 参考类型

让我们定义一个Go类型以包含上一个示例中的数据：

```go
type FamilyMember struct {
    Name    string
    Age     int
    Parents []string
}
var m FamilyMember
err := json.Unmarshal(b, &m)
```

将数据分解为一个`FamilyMember`值可以按预期工作，但是如果仔细观察，我们可以看到发生了一件了不起的事情。使用var语句，我们分配了一个`FamilyMember`结构，然后将指向该值的指针提供给`Unmarshal`，但那时该`Parents`字段是一个`nil`切片值。要填充该`Parents`字段，请`Unmarshal`在幕后分配一个新切片。这是`Unmarshal`与支持的参考类型（指针，切片和地图）一起使用的典型方式。

考虑拆封到此数据结构中：

```go
type Foo struct {
    Bar *Bar
}
```

如果`Bar` JSON对象中有一个字段，`Unmarshal`则将分配一个新字段 `Bar`并填充它。如果不是，`Bar`则将其留为`nil`指针。

由此产生一种有用的模式：如果您的应用程序接收一些不同的消息类型，则可以定义“接收器”结构，例如

```go
type IncomingMessage struct {
    Cmd *Command
    Msg *Message
}
```

发送方可以根据他们想要传达的消息类型来填充顶级JSON对象的`Cmd`字段和/或`Msg`字段。 `Unmarshal`，当将JSON解码为`IncomingMessage`结构时，只会分配JSON数据中存在的数据结构。要知道这消息的过程中，程序员需要简单地测试，要么`Cmd`或`Msg`不是`nil`。

# 流编码器和流解码器

json 包提供了 `Decoder` 和 `Encoder` 类型，以支持读写 JSON 数据流的通用操作。`NewDecoder()` 和 `NewEncoder()` 函数包装了 [io.Reader](https://golang.org/pkg/io/#Reader) 和 [io.Writer](https://golang.org/pkg/io/#Writer)接口类型。

```go
func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
```

这是一个示例程序，该程序从标准输入读取一系列 JSON 对象，`Name` 从每个对象中删除除字段以外的所有内容，然后将这些对象写入标准输出：

```go
package main
import (
    "encoding/json"
    "log"
    "os"
)
func main() {
    dec := json.NewDecoder(os.Stdin)
    enc := json.NewEncoder(os.Stdout)
    for {
        var v map[string]interface{}
        if err := dec.Decode(&v); err != nil {
            log.Println(err)
            return
        }
        for k := range v {
            if k != "Name" {
                delete(v, k)
            }
        }
        if err := enc.Encode(&v); err != nil {
            log.Println(err)
        }
    }
}
```

由于读取器和编写的普及，这些 `Encoder` 和 `Decoder` 类型可以在宽范围内的情况下，如读出和写入 HTTP 连接，的 WebSockets，或文件中使用。

## 参考

有关更多信息，请参阅 [json包文档](https://golang.org/pkg/encoding/json/)。有关json的用法示例，请参阅 [jsonrpc包](https://golang.org/pkg/net/rpc/jsonrpc/) 的源文件。

## go 中的 Decoder 和 Encoder 有什么用？为什么有了 Marshal 和 Unmarshal 还要有流式编码和解码的逻辑？

下面回答来自 NewBing：

你好，这是必应。我会尽力用中文回答你的问题。

go 中的 Decoder 和 Encoder 是用来实现流式编码和解码的接口，它们可以将数据从一个 io.Reader 或 io.Writer 转换成 go 的数据类型，或者反过来。<sup>1</sup> <sup>4</sup> <sup>5</sup>

Marshal 和 Unmarshal 是一些方便的函数，它们可以将数据从一个字符串或字节切片转换成 go 的数据类型，或者反过来。<sup>1</sup> <sup>4</sup> <sup>5</sup>

有了 Marshal 和 Unmarshal，为什么还要有流式编码和解码的逻辑呢？主要有以下几个原因：

- 流式编码和解码可以节省内存，因为它们不需要一次性加载整个数据到内存中，而是可以逐渐处理数据。<sup>1</sup> <sup>4</sup>
- 流式编码和解码可以提高性能，因为它们可以在读写数据的同时进行转换，而不需要等待数据完全读写完毕。<sup>1</sup> <sup>4</sup>
- 流式编码和解码可以支持更多的数据格式，比如 XML，CSV，YAML 等，因为它们可以自定义如何读写数据。<sup>6</sup>

下面是一个简单的例子，演示了如何使用 Encoder 和 Decoder 来编码和解码一个结构体：

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// 创建一个缓冲区，用来存储编码后的数据
	var buf bytes.Buffer
	// 创建一个 Encoder，用来将 Person 结构体编码到缓冲区
	enc := gob.NewEncoder(&buf)
	// 创建一个 Person 实例
	p := Person{Name: "Alice", Age: 25}
	// 调用 Encode 方法，将 p 编码到缓冲区
	err := enc.Encode(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 打印缓冲区的内容（二进制格式）
	fmt.Println(buf.Bytes())

	// 创建一个 Decoder，用来将缓冲区的数据解码成 Person 结构体
	dec := gob.NewDecoder(&buf)
	// 创建一个空的 Person 实例
	var q Person
	// 调用 Decode 方法，将缓冲区的数据解码到 q 中
	err = dec.Decode(&q)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 打印 q 的内容（结构体格式）
	fmt.Println(q)
}
```

\-
\[1]: [go - json.Marshal vs Encoder.Encode - Stack Overflow](https://stackoverflow.com/questions/62379309/json-marshal-vs-encoder-encode)
\[2]: [go - in Golang, what is the difference between json encoding and ...](https://stackoverflow.com/questions/33061117/in-golang-what-is-the-difference-between-json-encoding-and-marshalling)
\[3]: [Golang JSON Marshal(encode) and Unmarshal(decode/parse) with examples](https://www.jsonformatter.io/golang-json-parse-marshal-unmarshal-encode-decode/)
\[4]: [encoding package - encoding - Go Packages](https://pkg.go.dev/encoding)
\[5]: [go - Wrapping gob decoder - Stack Overflow](https://stackoverflow.com/questions/70555069/wrapping-gob-decoder)
\[6]: [go - How to marshal array to binary and unmarshal binary to array in ...](https://stackoverflow.com/questions/44290639/how-to-marshal-array-to-binary-and-unmarshal-binary-to-array-in-golang)

# 流式与非流式的区别

来自 Claude 的回答

json.Unmarshal() 和 json.Decode() 的区别主要有:

- Unmarshal() 接收 []byte 类型；Decode() 接收 io.Reader 接口类型,比如 strings.NewReader。
- Unmarshal() 会把整个 JSON 字符串解析成 Go 结构体；Decode() 采用流式解析,逐行解析 JSON 字符串。
- 如果 JSON 字符串很大,Unmarshal() 会把整个字符串加载到内存中,然后解析,可能导致 OOM；而 Decode() 采用流式解析,不会出现 OOM 的问题。
- Unmarshal() 解析失败会返回错误；Decode() 解析失败可以从 err 中获取错误信息,并继续解析其他 JSON 数据。

所以总的来说:

- 如果 JSON 数据量不大,使用 Unmarshal() 更方便。
- 如果 JSON 数据量很大,要防止 OOM,用 Decode() 流式解析。
- 如果 JSON 字符串包含多条 JSON 数据,用 Decode() 可以逐条解析,Unmarshal() 只能解析第一条数据。

根据你的代码示例,两个方法的输出结果是一样的,因为 JSON 字符串数据量很小,所以两种方法都可以使用。但如果数据量很大,Unmarshal() 可能会导致 OOM,这时 Decode() 是更好的选择。

**总得来说，如果是从 HTTP 的响应体中读取 JSON 数据的情况，建议使用流式的 json.NewDecoder().Decode() 来处理 JSON 数据转换为 Go 结构体的逻辑。**

# 其他文章

使用了太长时间的python，对于强类型的Golang适应起来稍微有点费力，不过操作一次之后发现，只有这么严格的类型规定，才能让数据尽量减少在传输和解析过程中的错误。我尝试使用Golang创建了一个公司的OpenAPI的demo，记录一下中间遇到的问题。

## 编码(Encode)Json

首先来看下如何将字典编码成Json：

```go
// 首先使用字面量来申明和初始化一个字典
param := map[string]int{"page_no": 1, "page_size": 40}
paramJson, err := json.Marshal(param)
```

使用json.Marshal接收需要json.encode的变量。而json.Marshal接收的是interface{}接口变量，该接口变量可以接收任何类型的数据。

## Http 包的 POST 请求来实践对 JSON 的序列化、反序列化

当我们把json编码好之后我们需要将信息传递给服务器。所以用到了http包。

在使用了之后我觉得go的http包真的非常方便，的确如传言中描述的强大和人性化，方便实用。

```go
resp , err := http.PostForm(requestUrl, url.Values{"api_key": {ApiKey}, "api_sign": {apiSign},
"param": {string(param)}, "time": {now_time}, "version": {version}})
```

这里我使用 http.PostForm 方法使用带参数传递的 post 方法请求服务器。url.Values 后面可以跟 `key[string][]string` 的形式传递参数。返回一个 http.response 结构体指针和一个 error 类型。

http.response 具体带有哪些属性可以详细查看一下包，这里我们会去解析他的 Body 字段，里面存储着返回的内容：

```go
// The Body is automatically dechunked if the server replied
// with a "chunked" Transfer-Encoding.
Body io.ReadCloser
```

这里 Body 是一个有 io.ReadCloser 接口的值。io.ReadCloser 接口实现了 Read() 和 Write() 方法。

我会用json的Decoder去解析它：

```go
var response openApiResponse
resp := request.RequestHeader(paramJson, version, SyncUrl)
err1 := json.NewDecoder(resp.Body).Decode(&response)
if err1 != nil {
    log.Println(err1)
}
return resp
```

这里json.NewDecoder接收一个有Reader方法的变量，之后我们调用了Decoder的方法decode将里面的内容都存入事先申请好的response结构体变量中。这个变量初始化了我们通过文档了解到的返回的结构体字段类型。

```go
openApiResponse struct {
    Success    bool   `json:"success"`
    ResultCode int    `json:"result_code"`
    ResultMsg  string `json:"result_msg"`
// 接收JSON字段
    Result GoodsSyncResult `json:"result"`
}
```

这样一级一级解析下去，在构造接收返回回来数据的结构体的时候，注意到后面的json字段。他是一个tag，可以在解析json的时候将对应名字的tag解析到对应的变量中。

这样就相当于你做好了数据结构，然后将对应的数据放到对应的字段里面去。

当然还有一种办法，当你不知道你所接收数据的数据结构的时候，你是没有办法提前声明好这些数据结构然后来接收的。这时我们可以申明一个空接口interface{}，让空接口的指针来接收这组数据，可以查看这组数据的数据结构。

```go
var hahaha interface{}
resp := request.RequestHeader(paramJson, version, SyncUrl)
err1 := json.NewDecoder(resp.Body).Decode(&hahaha)
if err1 != nil {
    log.Println(err1)
}
```

上面的 hahaha 可以接收并 decodejson，来接收这组数据。并且可以直接使用fmt.Print之类函数直接打印接收到的数据。如果想直接使用，我们可以使用类型断言但是更推荐的方法是，我们可以根据这组数据来写对应的结构体，然后将数据接收到结构体上进行操作。就像上面一样。

同样的我们还可以使用一个map[string]interface{}来接收这个Json以方便对其进行后续操作，避免不需要的多余的反射。

```go
var hahaha map[string]interface{}
resp := request.RequestHeader(paramJson, version, SyncUrl)
err1 := json.NewDecoder(resp.Body).Decode(&hahaha)
return hahaha
```

除了实现一个decoder来处理数据，我们往往有Json序列化之后就立即需要序列化的操作，这个同样很容易使用：

```go
json.Unmarshal([]byte, &xx)
```

来处理就好了。参数一是需要decode的Json数据, 参数二是用于接收这组数据的结构体字段。同样的我们也可以使用一个空接口来接收数据，也可以使用一一对应的结构体来放置数据。

看了上面的一堆介绍有一个感觉，就处理Json数据和类型转换来说。。python真是简单到爆炸，一个dumps一个loads轻松搞定。但是Golang严格的参数类型缺可以保证解析过来的数据一定是对应的数据结构和数据类型。不会在类型上报错更为严谨。个人觉得这很有趣，也很喜欢。
