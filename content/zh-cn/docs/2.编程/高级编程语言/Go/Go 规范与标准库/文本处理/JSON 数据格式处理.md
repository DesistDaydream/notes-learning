---
title: JSON 数据格式处理
---

# 概述

> 参考：
> - [Go 包，标准库-encoding/json](https://pkg.go.dev/encoding/json)
> - [Go 官方博客《JSON and Go》](https://blog.golang.org/json)
> - [骏马金龙](https://www.cnblogs.com/f-ck-need-u/p/10080793.html)
> - [在线 JSON 转 Go Struct](https://transform.tools/json-to-go)

**JavaScript Object Notation(简称 JSON)** 是一种简单的数据交换格式。从句法上讲，它类似于 JavaScript 的对象和列表。它最常用于 Web 后端与浏览器中运行的 JavaScript 程序之间的通信，但它也用于许多其他地方。它的主页 json.org 提供了一个清晰，简洁的标准定义。

**JSON 类型 与 Go 类型 对应关系**

    boolean >> bool
    number  >> float32,float64,int, int64, uint64
    string  >> string
    null    >> nil
    array   >> []interface{}
    object  >> map[string]interface{}

使用 json 包，可以轻松地从 Go 程序中读取和写入 JSON 数据。

# Encoding 与 Decoding

**Encoding(编码)** 与 **Decoding(解码)** 是 JSON 数据处理的基本操作

在 json 包中，使用`Marshal()`和 `Unmarshal()` 函数来执行最基本的 Encoding 与 Decoding 行为。

Marshal:直译为“编排、整理、排列、序列”，表示整理指定的内容，将内容整理成 json 数据。所以有时候也称此行为叫 **serializable(序列化)。**这种称呼是相对的。在计算机中特指将数据按某种描述格式编排出来，通常来说一般是从非文本格式到文本格式的数据转化。unmarshal 自然是指 marshal 的逆过程。

> 比如在 WebService 中，我们需要把 go 的 struct 以 JSON 方式表示并在网络间传输，把 go struct 转化成 JSON 的过程就是 marshal.

用白话说：

- **Encoding 就是将 struct、slice、array、map 等 转换为 JSON 格式**
- **Decoding 就是将 JSON 格式转换为 struct、slice、array、map。**

### sturct 结构与 JSON 结构的对应关系

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
  "comments": [
    {
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
	NextPost  *Post         `json:"nextPost"`
	Comments  []*Comment    `json:"comments"`
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
  - **struct、slice、array、map 都可以转换成 json**
  - **struct 转换成 json 的时候，struck 中只有字段首字母大写的属性才会被转换**
  - **map 转换的时候，key 必须为 string**
- **封装的时候，如果是指针，会追踪指针指向的对象进行封装**
- JSON 对象仅支持字符串作为键；要编码 Go map 类型，它必须采用以下形式`map[string]T`（`T` json 包支持的所有 Go 类型）。
- Channel、complex、function 类型无法编码。
- 不支持循环数据结构；它们将导致`Marshal`陷入无限循环。
- 指针将被编码为其所指向的值（如果指针为，则为“ null” `nil`）。
- json 包仅访问结构类型（以大写字母开头的结构类型）的导出字段。因此，仅 struct 的导出字段将出现在 JSON 输出中。

在 json 包中，可以使用`Marshal()`或者 `Mashallndent()` 函数来执行 Encoding 行为。

    func Marshal(v interface{}) ([]byte, error)

### 简单示例

假如现在有一个名为`Message`的 struct(结构体)，这个结构体表示一条消息中应该具有的属性。比如发送者、消息内容、发送时间，等等。

```go
type Message struct {
    Name string
    Body string
    Time int64
}
```

要想将这个 struct 中的数据转换为 JSON 格式，只需要使用 `Marshal()` 函数即可

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

    {"Name":"Alice","Body":"Hello","Time":1294706395881547000}

注意：
由于转换规范的原因导致 json 格式数据的 key 的首字母都是大写的，如果想要小写的，只需要给 struct 属性添加注释可，比如：

    type Message struct {
        Name string `json:"name"`
        Body string `json:"body"`
        Time int64  `json:"time"`
    }

那么输出结果就是这样的：

    {"name":"Alice","body":"Hello","time":1294706395881547000}

MarshalIndent() 函数，则是可以在 Encoding 成 JSON 的时候进行美化，将会自动添加前缀和缩进(前缀字符串一般设置为空)

    c,err := json.MarshalIndent(Message,"","\t")
    if err != nil {
    	fmt.Println(nil)
    }
    fmt.Println(string(c))

输出结果为：

    {
    	"Name": "Alice",
    	"Body": "Hello",
    	"Time": 1294706395881547000
    }

## Decoding(解码)

要解码 JSON 数据，我们使用`Unmarshal()`函数。

> Marshal 有整理、排列、序列的含义，表示整理指定的内容，将内容整理成 json 数据。那么 Unmarshal 就是 打散 这种含义。有时候也称为 **反序列化。**
> 比如可以这么描述：将 json 数据反序列化成指定的数据

    func Unmarshal(data []byte, v interface{}) error

我们首先必须创建一个存储解码数据的地方

    var m Message

并调用`json.Unmarshal`，将`[]byte`JSON 数据和一个指针传递给它`m`

    err := json.Unmarshal(b, &m)

如果`b`包含有效的 JSON，适合在`m`后电话`err`将`nil`与从数据`b`将被存储在结构`m`，仿佛像一个任务：

    m = Message{
        Name: "Alice",
        Body: "Hello",
        Time: 1294706395881547000,
    }

如何`Unmarshal`识别存储解码数据的字段？对于给定的 JSON 键`"Foo"`， `Unmarshal`将浏览目标结构的字段以查找(按优先顺序)：

- 标记为的导出字段`"Foo"`（ 有关 struct 标记的更多信息，请参见 Go 规范），
- 名为`"Foo"`或的导出字段
- 名为`"FOO"`或`"FoO"`或其他不区分大小写的匹配项的导出字段`"Foo"`。

当 JSON 数据的结构与 Go 类型不完全匹配时会发生什么？

    b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
    var m Message
    err := json.Unmarshal(b, &m)

`Unmarshal`只会解码在目标类型中可以找到的字段。在这种情况下，将仅填充 m 的 Name 字段，而 Food 字段将被忽略。当您希望从大型 JSON Blob 中仅选择几个特定字段时，此行为特别有用。这也意味着目标 struct 中所有未导出的字段都不会受到的影响`Unmarshal`。

但是，如果您事先不知道 JSON 数据的结构怎么办？

# 使用 interface{} 存放通用 JSON 数据

的`interface{}`（空接口）类型描述了具有零种方法的接口。每个 Go 类型至少实现零个方法，因此满足空接口。

空接口用作常规容器类型：

    var i interface{}
    i = "a string"
    i = 2011
    i = 2.777

类型断言访问基础的具体类型：

    r := i.(float64)
    fmt.Println("the circle's area", math.Pi*r*r)

或者，如果基础类型未知，则由类型开关确定类型：

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

json 包使用`map[string]interface{}`和 `[]interface{}`值来存储任意 JSON 对象和数组；它将很乐意将任何有效的 JSON Blob 解组为纯 `interface{}`值。默认的具体 Go 类型为：

- `bool` 对于 JSON 布尔值，
- `float64` 对于 JSON 数字，
- `string` 用于 JSON 字符串，以及
- `nil` JSON null。

# 解码任意数据

考虑以下存储在变量中的 JSON 数据`b`：

    b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

在不知道此数据结构的情况下，我们可以使用以下命令将其解码为一个`interface{}`值`Unmarshal`：

    var f interface{}
    err := json.Unmarshal(b, &f)

此时，Go 值`f`将是一个映射，其键为字符串，其值本身存储为空接口值：

    f = map[string]interface{}{
        "Name": "Wednesday",
        "Age":  6,
        "Parents": []interface{}{
            "Gomez",
            "Morticia",
        },
    }

要访问此数据，我们可以使用类型断言来访问`f`的底层`map[string]interface{}`：

    m := f.(map[string]interface{})

1
Plain Text

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

这样，您可以使用未知的 JSON 数据，同时仍然享有类型安全的好处。

# 参考类型

让我们定义一个 Go 类型以包含上一个示例中的数据：

    type FamilyMember struct {
        Name    string
        Age     int
        Parents []string
    }
    var m FamilyMember
    err := json.Unmarshal(b, &m)

将数据分解为一个`FamilyMember`值可以按预期工作，但是如果仔细观察，我们可以看到发生了一件了不起的事情。使用 var 语句，我们分配了一个`FamilyMember`结构，然后将指向该值的指针提供给`Unmarshal`，但那时该`Parents`字段是一个`nil`切片值。要填充该`Parents`字段，请`Unmarshal`在幕后分配一个新切片。这是`Unmarshal`与支持的参考类型（指针，切片和地图）一起使用的典型方式。
考虑拆封到此数据结构中：

    type Foo struct {
        Bar *Bar
    }

如果`Bar`JSON 对象中有一个字段，`Unmarshal`则将分配一个新字段 `Bar`并填充它。如果不是，`Bar`则将其留为`nil`指针。
由此产生一种有用的模式：如果您的应用程序接收一些不同的消息类型，则可以定义“接收器”结构，例如

    type IncomingMessage struct {
        Cmd *Command
        Msg *Message
    }

发送方可以根据他们想要传达的消息类型来填充顶级 JSON 对象的`Cmd`字段和/或`Msg`字段。 `Unmarshal`，当将 JSON 解码为`IncomingMessage`结构时，只会分配 JSON 数据中存在的数据结构。要知道这消息的过程中，程序员需要简单地测试，要么`Cmd`或`Msg`不是`nil`。

# 流编码器和解码器

json 包提供了`Decoder`和`Encoder`类型，以支持读写 JSON 数据流的通用操作。的`NewDecoder`和`NewEncoder`功能包裹`[io.Reader](https://golang.org/pkg/io/#Reader)` 和`[io.Writer](https://golang.org/pkg/io/#Writer)`接口类型。

    func NewDecoder(r io.Reader) *Decoder
    func NewEncoder(w io.Writer) *Encoder

这是一个示例程序，该程序从标准输入读取一系列 JSON 对象，`Name`从每个对象中删除除字段以外的所有内容，然后将这些对象写入标准输出：

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

由于读取器和编写的普及，这些`Encoder`和`Decoder`类型可以在宽范围内的情况下，如读出和写入 HTTP 连接，的 WebSockets，或文件中使用。

## 参考

有关更多信息，请参阅[json 包文档](https://golang.org/pkg/encoding/json/)。有关 json 的用法示例，请参阅[jsonrpc 包](https://golang.org/pkg/net/rpc/jsonrpc/)的源文件。

# 相关文章

- [用于协议缓冲区的新 Go API](https://blog.golang.org/protobuf-apiv2)
- [在 Go 1.13 中处理错误](https://blog.golang.org/go1.13-errors)
- [调试在 Go 1.12 中部署的内容](https://blog.golang.org/debug-opt)
- [HTTP / 2 服务器推送](https://blog.golang.org/h2push)
- [介绍 HTTP 跟踪](https://blog.golang.org/http-tracing)
- [产生程式码](https://blog.golang.org/generate)
- [隆重推出 Go Race Detector](https://blog.golang.org/race-detector)
- [行动地图](https://blog.golang.org/maps)
- [去你的代码](https://blog.golang.org/gofmt)
- [组织 Go 代码](https://blog.golang.org/organizing-go-code)
- [使用 GNU 调试器调试 Go 程序](https://blog.golang.org/debug-gdb)
- [Go 图片/绘图包](https://blog.golang.org/image-draw)
- [Go 图像包](https://blog.golang.org/image)
- [反射定律](https://blog.golang.org/laws-of-reflection)
- [错误处理和执行](https://blog.golang.org/error-handling-and-go)
- [Go 中的一流函数](https://blog.golang.org/functions-codewalk)
- [分析 Go 程序](https://blog.golang.org/pprof)
- [GIF 解码器：Go 接口中的练习](https://blog.golang.org/gif-decoder)
- [介绍 Gofix](https://blog.golang.org/introducing-gofix)
- [Godoc：记录 Go 代码](https://blog.golang.org/godoc)
- [数据块](https://blog.golang.org/gob)
- [C？走？go！](https://blog.golang.org/cgo)
- [切成薄片：用法和内部原理](https://blog.golang.org/slices-intro)
- [Go 并发模式：超时，继续前进](https://blog.golang.org/concurrency-timeouts)
- [推迟，恐慌和恢复](https://blog.golang.org/defer-panic-and-recover)
- [通过通信共享内存](https://blog.golang.org/codelab-share)
- [JSON-RPC：接口的故事](https://blog.golang.org/json-rpc)

# 其他文章

使用了太长时间的 python，对于强类型的 Golang 适应起来稍微有点费力，不过操作一次之后发现，只有这么严格的类型规定，才能让数据尽量减少在传输和解析过程中的错误。我尝试使用 Golang 创建了一个公司的 OpenAPI 的 demo，记录一下中间遇到的问题。

## 编码(Encode)Json：

首先来看下如何将字典编码成 Json：

    // 首先使用字面量来申明和初始化一个字典
    param := map[string]int{"page_no": 1, "page_size": 40}
    paramJson, err := json.Marshal(param)

1
2
3
Go

使用 json.Marshal 接收需要 json.encode 的变量。而 json.Marshal 接收的是 interface{}接口变量，该接口变量可以接收任何类型的数据。

## Http 包的 POST 请求来实践对 JSON 的序列化、反序列化：

当我们把 json 编码好之后我们需要将信息传递给服务器。所以用到了 http 包。

在使用了之后我觉得 go 的 http 包真的非常方便，的确如传言中描述的强大和人性化，方便实用。

    resp , err := http.PostForm(requestUrl, url.Values{"api_key": {ApiKey}, "api_sign": {apiSign},
    "param": {string(param)}, "time": {now_time}, "version": {version}})

1
2
Go

这里我使用 http.PostForm 方法使用带参数传递的 post 方法请求服务器。url.Values 后面可以跟 key\[string]\[]string 的形式传递参数。返回一个 http.response 结构体指针和一个 error 类型。

http.response 具体带有哪些属性可以详细查看一下包，这里我们会去解析他的 Body 字段，里面存储着返回的内容：

    // The Body is automatically dechunked if the server replied
    // with a "chunked" Transfer-Encoding.
    Body io.ReadCloser

1
2
3
Go

这里 Body 是一个有 io.ReadCloser 接口的值。io.ReadCloser 接口实现了 Read()和 Write()方法。

我会用 json 的 Decoder 去解析它：

    var response openApiResponse
    resp := request.RequestHeader(paramJson, version, SyncUrl)
    err1 := json.NewDecoder(resp.Body).Decode(&response)
    if err1 != nil {
        log.Println(err1)
    }
    return resp

1
2
3
4
5
6
7
Go

这里 json.NewDecoder 接收一个有 Reader 方法的变量，之后我们调用了 Decoder 的方法 decode 将里面的内容都存入事先申请好的 response 结构体变量中。这个变量初始化了我们通过文档了解到的返回的结构体字段类型。

    openApiResponse struct {
        Success    bool   `json:"success"`
        ResultCode int    `json:"result_code"`
        ResultMsg  string `json:"result_msg"`
    // 接收JSON字段
        Result GoodsSyncResult `json:"result"`
    }

1
2
3
4
5
6
7
Go

这样一级一级解析下去，在构造接收返回回来数据的结构体的时候，注意到后面的 json 字段。他是一个 tag，可以在解析 json 的时候将对应名字的 tag 解析到对应的变量中。

这样就相当于你做好了数据结构，然后将对应的数据放到对应的字段里面去。

当然还有一种办法，当你不知道你所接收数据的数据结构的时候，你是没有办法提前声明好这些数据结构然后来接收的。这时我们可以申明一个空接口 interface{}，让空接口的指针来接收这组数据，可以查看这组数据的数据结构。

    var hahaha interface{}
    resp := request.RequestHeader(paramJson, version, SyncUrl)
    err1 := json.NewDecoder(resp.Body).Decode(&hahaha)
    if err1 != nil {
        log.Println(err1)
    }

1
2
3
4
5
6
Go

上面的 hahaha 可以接收并 decodejson，来接收这组数据。并且可以直接使用 fmt.Print 之类函数直接打印接收到的数据。如果想直接使用，我们可以使用类型断言但是更推荐的方法是，我们可以根据这组数据来写对应的结构体，然后将数据接收到结构体上进行操作。就像上面一样。

同样的我们还可以使用一个 map\[string]interface{}来接收这个 Json 以方便对其进行后续操作，避免不需要的多余的反射。

    var hahaha map[string]interface{}
    resp := request.RequestHeader(paramJson, version, SyncUrl)
    err1 := json.NewDecoder(resp.Body).Decode(&hahaha)
    return hahaha

1
2
3
4
Plain Text

除了实现一个 decoder 来处理数据，我们往往有 Json 序列化之后就立即需要序列化的操作，这个同样很容易使用：

    json.Unmarshal([]byte, &xx)

1
Plain Text

来处理就好了。参数一是需要 decode 的 Json 数据, 参数二是用于接收这组数据的结构体字段。同样的我们也可以使用一个空接口来接收数据，也可以使用一一对应的结构体来放置数据。

看了上面的一堆介绍有一个感觉，就处理 Json 数据和类型转换来说。。python 真是简单到爆炸，一个 dumps 一个 loads 轻松搞定。但是 Golang 严格的参数类型缺可以保证解析过来的数据一定是对应的数据结构和数据类型。不会在类型上报错更为严谨。个人觉得这很有趣，也很喜欢。
