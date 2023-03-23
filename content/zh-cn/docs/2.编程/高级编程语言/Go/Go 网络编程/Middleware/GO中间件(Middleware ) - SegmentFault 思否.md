---
title: GO中间件(Middleware ) - SegmentFault 思否
---

原文链接：[GO 中间件 (Middleware) - SegmentFault 思否](https://segmentfault.com/a/1190000018819804)

中间件是一种计算机[软件](https://link.segmentfault.com/?url=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FSoftware)，可为[操作系统](https://link.segmentfault.com/?url=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FOperating_system)提供的[软件应用程序](https://link.segmentfault.com/?url=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FSoftware_application)提供服务，以便于各个软件之间的沟通，特别是系统软件和应用软件。广泛用于 web 应用和面向服务的体系结构等。

纵观 GO 语言，中间件应用比较普遍，主要应用：

- 记录对服务器发送的请求（request）
- 处理服务器响应（response ）
- 请求和处理之间做一个权限认证工作
- 远程调用
- 安全
- 等等

**中间件处理程序**是简单的`http.Handler`，它包装另一个`http.Handler`做请求的一些预处理和 / 或后处理。它被称为 “中间件”，因为它位于 Go Web 服务器和实际处理程序之间的中间位置。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ed41e30b-8c8b-4238-90f6-663f162aa592/1460000018819807)

下面是一些中间件例子

### 记录日志中间件

    package main

    import (
       "fmt"
       "log"
       "net/http"
    )

    func logging(f http.HandlerFunc) http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
          log.Println(r.URL.Path)
          f(w, r)
       }
    }
    func foo(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "foo")
    }

    func bar(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "bar")
    }

    func main() {
       http.HandleFunc("/foo", logging(foo))
       http.HandleFunc("/bar", logging(bar))
       http.ListenAndServe(":8080", nil)
    }

访问 [http://localhost](https://link.segmentfault.com/?url=http%3A%2F%2Flocalhost):8080/foo

返回结果

foo

将上面示例修改下，也可以实现相同的功能。

    package main

    import (
       "fmt"
       "log"
       "net/http"
    )

    func foo(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "foo")
    }
    func bar(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "bar")
    }

    func loggingMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          log.Println(r.URL.Path)
          next.ServeHTTP(w, r)
       })
    }

    func main() {

       http.Handle("/foo", loggingMiddleware(http.HandlerFunc(foo)))
       http.Handle("/bar", loggingMiddleware(http.HandlerFunc(bar)))
       http.ListenAndServe(":8080", nil)
    }

访问 [http://localhost](https://link.segmentfault.com/?url=http%3A%2F%2Flocalhost):8080/foo

返回结果

foo

### 多中间件例子

    package main

    import (
       "fmt"
       "log"
       "net/http"
       "time"
    )

    type Middleware func(http.HandlerFunc) http.HandlerFunc

    func Logging() Middleware {


       return func(f http.HandlerFunc) http.HandlerFunc {


          return func(w http.ResponseWriter, r *http.Request) {


             start := time.Now()
             defer func() { log.Println(r.URL.Path, time.Since(start)) }()


             f(w, r)
          }
       }
    }

    func Method(m string) Middleware {


       return func(f http.HandlerFunc) http.HandlerFunc {


          return func(w http.ResponseWriter, r *http.Request) {


             if r.Method != m {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
             }


             f(w, r)
          }
       }
    }

    func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
       for \_, m := range middlewares {
          f = m(f)
       }
       return f
    }

    func Hello(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "hello world")
    }

    func main() {
       http.HandleFunc("/", Chain(Hello, Method("GET"), Logging()))
       http.ListenAndServe(":8080", nil)
    }

中间件本身只是将其`http.HandlerFunc`作为其参数之一，包装它并返回一个新`http.HandlerFunc`的服务器来调用。在这里，我们定义了一种新类型`Middleware`，最终可以更容易地将多个中间件链接在一起。

当然我们也可以改成如下形式

    package main

    import (
       "fmt"
       "log"
       "net/http"
       "time"
    )

    type Middleware func(http.Handler) http.Handler

    func Hello(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "hello world")
    }

    func Chain(f http.Handler, mmap ...Middleware) http.Handler {
       for \_, m := range mmap {
          f = m(f)
       }
       return f
    }
    func Method(m string) Middleware {
       return func(f http.Handler) http.Handler {
          return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
             log.Println(r.URL.Path)
             if r.Method != m {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
             }
             f.ServeHTTP(w, r)
          })
       }

    }
    func Logging() Middleware {
       return func(f http.Handler) http.Handler {
          return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


             start := time.Now()
             defer func() { log.Println(r.URL.Path, time.Since(start)) }()
             f.ServeHTTP(w, r)
          })
       }
    }

    func main() {
       http.Handle("/", Chain(http.HandlerFunc(Hello), Method("GET"), Logging()))
       http.ListenAndServe(":8080", nil)
    }

### 在 gin 框架下实现中间件

    r := gin.Default() 创建带有默认中间件的路由，默认是包含logger和recovery中间件的
    r :=gin.new()      创建带有没有中间件的路由

示例

    package main

    import (
       "github.com/gin-gonic/gin"
       "log"
       "time"
    )

    func Logger() gin.HandlerFunc {
       return func(c *gin.Context) {
          t := time.Now()

          c.Set("example", "12345")

          c.Next()

          latency := time.Since(t)
          log.Print(latency)

          status := c.Writer.Status()
          log.Println(status)
       }
    }
    func main() {
       r := gin.New()
       r.Use(Logger())

       r.GET("/test", func(c *gin.Context) {
          example := c.MustGet("example").(string)


          log.Println(example)
       })


       r.Run(":8080")
    }

以上示例也可改为

    package main

    import (
       "github.com/gin-gonic/gin"
       "log"
       "time"
    )

    func Logger() gin.HandlerFunc {
       return func(c *gin.Context) {
          t := time.Now()

          c.Set("example", "12345")

          c.Next()

          latency := time.Since(t)
          log.Print(latency)

          status := c.Writer.Status()
          log.Println(status)
       }
    }

    func main() {
       r := gin.New()
       r.GET("/test", Logger(), func(c *gin.Context) {
          example := c.MustGet("example").(string)

          log.Println(example)
       })

       r.Run(":8080")
    }

即不用 r.use 添加中间件，直接将 Logger() 写到 r.GET 方法的参数里（"/test" 之后）。

更多 gin 中间件示例可参考 [https://github.com/gin-gonic/gin](https://link.segmentfault.com/?url=https%3A%2F%2Fgithub.com%2Fgin-gonic%2Fgin)

**参考资料**

[https://drstearns.github.io/t...](https://link.segmentfault.com/?url=https%3A%2F%2Fdrstearns.github.io%2Ftutorials%2Fgomiddleware%2F)

[https://gowebexamples.com/adv...](https://link.segmentfault.com/?url=https%3A%2F%2Fgowebexamples.com%2Fadvanced-middleware%2F)

- [更多 GO 学习资料](https://link.segmentfault.com/?url=https%3A%2F%2Fgithub.com%2Fguyan0319%2Fgolang_development_notes%2Fblob%2Fmaster%2Fzh%2Fpreface.md)
