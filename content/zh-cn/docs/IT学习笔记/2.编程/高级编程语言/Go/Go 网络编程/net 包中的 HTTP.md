---
title: net 包中的 HTTP
---

# 概述

> 参考：
> - [GitHub 项目，DesistDaydream/go-net](https://github.com/DesistDaydream/go-net)(学习代码)
> - [GoWeb 编程](https://github.com/astaxie/build-web-application-with-golang)
> - [看云 GoWeb 编程](https://www.kancloud.cn/kancloud/web-application-with-golang)

go 使用 net/http 标准库来实现基本的 web 功能

- **form(表单)** # 描述网页表单的处理
- **middleware(中间件)** # 常用来处理认证等行为

## 一般的上网过程概述

浏览器本身是一个客户端，当你输入 URL 的时候，首先浏览器会去请求 DNS 服务器，通过 DNS 获取相应的域名对应的 IP，然后通过 IP 地址找到 IP 对应的服务器后，要求建立 TCP 连接，等浏览器发送完 HTTP Request（请求）包后，服务器接收到请求包之后才开始处理请求包，服务器调用自身服务，返回 HTTP Response（响应）包；客户端收到来自服务器的响应后开始渲染这个 Response 包里的主体（body），等收到全部的内容随后断开与该服务器之间的 TCP 连接

# Hello World

```go
package main

import (
	"fmt"
	"net/http"
)

// HelloWorld 处理客户端请求 /hello 时的具体逻辑
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	// 将 Hello DesistDaydream! 这一串字符写入到 Response 中，并响应给客户端
	fmt.Fprintf(w, "Hello DesistDaydream!")
}

func main() {
	// 设置访问的路由,一般也称为 Handler(处理器)，用来处理 http 请求。比如这里就是处理一个访问 /hello 的 hettp 请求。
	// 当客户端发起 http 请求，访问 http://IP:8080/hello ，由 HelloWorld 函数处理该请求。
	http.HandleFunc("/hello", HelloWorld)

	// 设置监听的端口
	http.ListenAndServe(":8080", nil)
}
```

# net/http 包解析

以 1.16 版本为例
[http.Client{}](https://github.com/golang/go/blob/release-branch.go1.16/src/net/http/client.go#L57) 结构体 # HTTP 客户端，作用在该结构体的方法，就是用来发起 HTTP 请求的方法，比如 GET、POST 等
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xi368g/1626314106732-1e5da8f4-21b9-485f-988c-c99837e21513.png)

[do()](https://cs.opensource.google/go/go/+/refs/tags/go1.16.6:src/net/http/client.go;l=590) 方法 # do 是用来发送 HTTP 请求并返回 HTTP 响应的最本质方法。像 Get()、Post() 等方法，最终还是调用的 do()

# 使用 Go 发起 HTTP Request 并处理 Response Body 的基本示例

```go
package main

import (
	"bufio"
	"fmt"
	"net/http"
)

// 设置一些会用到的全局变量，省的每次都要重新初始化
var (
	req  *http.Request
	resp *http.Response
	err  error
)


// Client1 直接使用 http.Get() 来发起请求
func Client1() {
	// net/http 标准库中还可以实现作为客户端发送 http 请求
	// Get() 向指定的服务器发送一个 HTTP GET 请求，并返回一个 Response
	resp, err := http.Get("http://172.38.40.250:8080/index")
	if err != nil {
		panic(err)
	}
	// 关闭连接
	defer resp.Body.Close()

	// 输出服务端响应的的状态码
	fmt.Println("Response status:", resp.Status)

	// 处理 Response 中的 Body，并输出响应体的字符串格式内容
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// Client2 先构建一个 Request，再根据这个 Request 发起请求，这种方式常用来自定义请求内容
func Client2() {
	// 构建 Request
	req, _ := http.NewRequest("GET", "http://172.38.40.250:8080/index", nil)
	// 为构建的 Request 设定请求头信息，可以多次使用 Set() 来设定多个 Header 信息
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	// 查看一下将要发起的请求内容
	fmt.Printf("本次 HTTP Request 为：%v\n请求方法为：%v\n请求头为：%v\n", req, req.Method, req.Header)

	// 根据新构建的 req 来发起请求，并获取响应信息
	// 这里的 http.Client{} 中可以设置一些发起 HTTP 请求时的一些信息，比如 TLS 等
    client := &http.Client{}
	if resp, err = client.Do(req); err != nil {
		panic(err)
	}
	// 关闭连接
	defer resp.Body.Close()

	// 处理响应，并输出 Response Body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	Client1()
	Client2()
}
```

# [Go HTTP GET/POST JSON 的服务端、客户端示例，包含序列化、反序列化](https://www.cnblogs.com/junneyang/p/6211190.html)

## 服务端代码示例：

```go
package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Message 是一条消息应该具有的基本属性
type Message struct {
	Name string `json:"name"`
	Body string `json:"body"`
	Time string `json:"time"`
}

// NewMessage 实例化 Message
func NewMessage() *Message {
	return &Message{
		Name: "DesistDaydream",
		Body: "Hello World",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// ResponseJSON 将会响应 JSON 格式数据
func ResponseJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("当前客户端的请求 %v 页面的 Method 为：%v\n", r.RequestURI, r.Method)
	// 初始化结构体，用于存储和响应 JSON 数据
	m := NewMessage()
	// 声明两个常用的变量
	var err error
	var jsonData []byte
	// 根据不同请求方法，执行不同的行为
	switch r.Method {
	case "GET":
		// 将 struct 中的数据转换为 JSON 格式
		if jsonData, err = json.Marshal(m); err != nil {
			fmt.Println(err)
			return
		}
		// 响应 JSON 格式的默认值
		fmt.Fprintf(w, string(jsonData))
	default:
		// 模拟下面这样的 curl 请求，程序将会根据 Request Body 中的内容替换 Message 结构体数据中的值，并返回结构体中的数据
		// 这就好比请求一个需要 TOKEN 的 API，我们只有使用正确的 TOKEN，才可以获取想要的信息
		// curl -XPOST http://172.38.40.250:8080/json -d '{"name":"lichenhao"}'
		//
		// 读取 Request 的 Body
		RequestBody, _ := io.ReadAll(r.Body)
		fmt.Printf("请求体为：%v\n", string(RequestBody))
		// 将 Request Body 的 JSON 格式转换为 struct 类型，并将 struct 中的值替换为 JSON 中的值
		// 注意，struct 中仅传入一个 key 的值，则 struct 中也只有一个属性的值被替代，其他属性的值保持不变
		if err = json.Unmarshal(RequestBody, m); err != nil {
			fmt.Fprintf(w, "请检查 Body，格式不正确或数据类型不对")
			return
		}
		fmt.Printf("请求体转换为 struct 后的值为：%v\n", m)

		// 根据传入的 请求体 的值，判断认证是否成功
		// 比如现在假设，只有传入 {"name":"DesistDaydream"} 这个请求体时，才会响应结构体的数据给给客户端。
		switch m.Name {
		case "DesistDaydream":
			// 认证正确，将 struct 类型数据转换为 JSON 格式数据并响应给客户端
			if jsonData, err = json.Marshal(m); err != nil {
				fmt.Fprintf(w, "序列化出错，请始终其他数据格式的 Body")
			}
			fmt.Fprint(w, string(jsonData))
		default:
			fmt.Fprintf(w, "你好 %v,认证失败，请重试\n", m.Name)
		}
	}
}

func main() {
	// 设置访问的路由
	http.HandleFunc("/json", handler.ResponseJSON)

	// 设置监听的端口
	fmt.Println("开始监听8080端口")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

## 客户端代码示例：

```go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Message 是一条消息应该具有的基本属性
type Message struct {
	Name string `json:"name"`
	Body string `json:"body"`
	Time string `json:"time"`
}

func main() {
	// 第一种请求
	// 模拟从外部读取 json 格式文件，将 json 与 struct 绑定，然后再发送请求
	m := NewMessage()
	// 这里假定 struct 的值时从外部文件获取的
	m.Name = "DesistDaydream"
	m.Body = "你好"
	m.Time = time.Now().Format("2006-01-02 15:04:05")
	// 将 struct 转换为 json
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	// 构建 Request
	req, _ := http.NewRequest("POST", "http://172.38.40.250:8080/json", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	// 处理响应信息并输出
	resp, _ := (&http.Client{}).Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	// 第二种请求
	// 手动指定 json 数据，并发起请求
	// 下面的代码等效于使用 crul 命令发起请求
	// curl -XPOST http://172.38.40.250:8080/json -d '{"name":"lichenhao"}'
	// 创建一个 json 数据
	jsonReqBody := []byte(`{"name":"lichenhao"}`)
	// 构建 Request
	req, _ = http.NewRequest("POST", "http://172.38.40.250:8080/json", bytes.NewBuffer(jsonReqBody))
	req.Header.Set("Content-type", "application/json")
	// 处理响应信息并输出
	resp, _ = (&http.Client{}).Do(req)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```
