---
title: cURL
---

# 概述

> 参考：
>
> - [官网](https://curl.se/)
> - [manual(手册)](https://man.cx/curl)

cURL 是一个用于 URLs 传输的命令行工具和库，始于 1998 年。

早在 20 世纪 90 年代中期，当时互联网还处于萌芽阶段，瑞典程序员 Daniel Stenberg 开始了一个项目，这个项目最终发展成了我们今天所知道的 curl 程序。

最初，他的目标是开发一种机器人，可以定期从网页上下载货币汇率，并向 IRC 用户提供等值的瑞典克朗美元。这个项目蓬勃发展，在这个过程中添加了几个协议和特性——剩下的就是历史了。

curl 是常用的命令行工具，用来请求 Web 服务器。它的名字就是“客户端(client)的 URL ”工具的意思。

> 注意：curl 最常见是通过网络 URL 来传输数据。但是，curl 还可以通过 Socket 的 URL 来传输数据，只需要使用 --unix-socket 选项指定 Socket 路径即可。

# Syntax(语法)

**curl \[OPTIONS] URL....**
如果没有另外说明，curl 将接收到的数据写入标准输出。可以使用 -o、--output 或 -O、--remote-name 选项将该数据保存到本地文件中。如果 curl 被赋予多个 URL 以在命令行上传输，它同样需要多个选项来保存它们。curl 不会解析或以其他方式“理解”它作为输出获取或写入的内容。它不进行编码或解码，除非使用专用命令行选项明确要求。

## OPTIONS

由于 curl 程序支持多种协议，可以使用各种不同的协议向指定的 URL 发起请求，所以，并不是所有选项都适用于所有协议。在下面的笔记中，每个选项后面会添加一个 `()`，括号中说明此选项支持的协议，多个协议以空格分割；没有 `()` 的表示该选项适用于所有协议。若括号内为 TLS 则表示使用安全的各种协议，比如 https、ftps、imaps 等等

- **--compressed** # (HTTP)使用 curl 支持的算法之一请求压缩响应，并自动解压缩响应体。
  - 在有 Sever 端会将响应体压缩，若不使用该选项，则响应体是无法输出到标准输出的，且会出现如下报错：

```bash
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.
```

- **-d,--data DATA>** # (HTTP MQTT)使用 POST 请求将指定数据作为请求体。
  - 与 `Content-Type: application/x-www-form-urlencoded` 头信息配合，发送的 DATA 是 x-www-form-urlencoded 类型的请求体数据。
- **--data-urlencode DATA>** # (HTTP)与 -d 选项类似，发起 POST 请求，但是它执行 URL 编码。(urlencode 就是 URL Encode)
- **-f,--fail** # (HTTP)连接失败时不显示 HTTP 错误信息
- **-g, --globoff** # 关闭 `URL Globbing Parser(URL全局解析器)`。设置此选项，则 URL 中可以包含 `{}` 和 `[]` 符号，这些符号将被当做字符。
  - 该选项常用来配合 IPv6 使用
- **-H,--header STRING>**# (HTTP)使用指定的 STRING 作为请求 header 发送给服务器
  - STRING 可以使用 @FILE 格式来通过文件传递请求头信息。
- **-I,--head**# (HTTP FTP FILE)只显示本次请求的 Header 信息。当用于 FTP 或 FILE 时，则只显示文件大小和最后修改时间。
- **-k,--insecure** # (TLS)此选项表示此次 curl 请求允许"不安全"的 SSL 连接和传输。也就是说对于 https 请求，可以允许私有证书。如果使用 curl 进行 https 请求的时候，不使用该参数的话，服务端使用的私有证书或自建 CA 的证书，则有可能产生如下报错
  - curl: (60) Peer's certificate issuer has been marked as not trusted by the user.
  - curl: (60) Peer's Certificate issuer is not recognized.
- **-L, --location** # (HTTP)如果服务器报告所请求的页面已移动到其他位置（用 Location：标题和 3XX 响应代码表示），则此选项将使 curl 重做新位置的请求。
  - Note:如果下载文件出错之后，发现文件大小异常，则说明该文件被移动到其他链接下了，需要使用-L 与-O 配合使用才能正确下载
- **--limit-rate NUM>** # 限制现在时的速率，NMU 为每秒下载速度，单位可以使 K、M、G
- **-m, --max-time TIME>**# 指定 curl 不管访问成功还是失败，最大消耗时间为 TIME。TIME 时间后服务端未响应，则视为无法连接。
- **-O, --remote-name**# 将输入写入的一个文件中，默认的文件名与请求的资源的名称一样。i.e.下载文件
  - curl -O <https://www.example.com/foo/bar.html> # 将服务器回应保存成文件，文件名为 bar.html。
  - 可以在一条命令中多次使用 -O 来下载多个文件
- **-o, --output FileName>** # 与 -O 一样，下载文件，只不过可以自己制定下载到本地后的文件名。可以重定向到 /dev/null，以便隐藏输出。
- **--resolve DN:PORT:IP,IP...** # 指定将 DN(域名)解析成哪个 IP。DN 可以使用通配符
- **-s, --silent** # 静默模式。将不输出错误和进度信息,但是会正常显示运行结果。
- **--trace**# 与-v 类似也可以用于调试，还会输出原始的二进制数据。
- **-u, --user <USER:[PASSWORD]>** # 指定发起请求时，所使用的基本认证信息。若省略 PASSWORD，则会以交互方式，在执行命令之后输入。
- **--unix-socket PATH>** # (HTTP)通过 Unix 套接字连接，而不使用网络。
- **-v, --verbose** # 输出通信的整个过程，用于调试。
- **-w, --write-out FORMAT>** # 指定在 curl 完成后，输出的信息，详细介绍可以参考[样例](/docs/1.操作系统/X.Linux%20管理/Linux%20网络管理工具/cURL.md)
- **-X, --request METHOD>** # (HTTP)指定 HTTP 请求的方法。

### -w, --write-out 选项详解

-w, --write-out 选项可以根据指定格式输出本次请求的一些统计信息。简单的效果如下所示：

```bash
~# curl -o /dev/null -s -w "DNS解析时间："%{time_namelookup}"\n"\
"TCP建立时间:"%{time_connect}"\n"\
"响应第一个字节时间:"%{time_starttransfer}"\n"\
"总时间:"%{time_total}"\n"\
"下载速度:"%{speed_download}"\n" \
"http://www.taobao.com"
# 下面是命令结果
DNS解析时间：0.534294
TCP建立时间:0.551090
响应第一个字节时间:0.569936
总时间:0.570050
下载速度:487.000
```

FORMAT 中可用字段：

- time_namelookup # DNS 解析域名\[www.taobao.com]的时间
- time_commect # client 和 server 端建立 TCP 连接的时间
- time_starttransfer # 从 client 发出请求；到 web 的 server 响应第一个字节的时间
- time_total # client 发出请求；到 web 的 server 发送会所有的相应数据的时间
- speed_download # 下载速度 单位 byte/s
- http_code # 本次请求的 http 响应状态码。
- content_type # 显示在 Header 里面使用 Content-Type 来表示的具体请求中的媒体类型信息；
- time_namelookup # 从请求开始，到 DNS 解析完成所使用的时间，经常用来排除 DNS 解析的问题；
- time_redirect # 重定向的时间；
- time_pretransfer # 从开始到准备传输的时间；
- time_connect # 连接时间，从开始到 TCP 三次握手完成时间，这里面包括 DNS 解析的时候，如果想求连接时间，需要减去上面的解析时间；
- time_starttransfer # 开始传输时间，从发起请求开始，到服务器返回第一个字段的时间；
- time_total # 总时间；
- speed_download # 经常使用它来测试网速度，下载速度，单位是字节每秒；
- size_request # 请求头的大小；
- size_header # 下载的 header 的大小；
- 等等

# EXAMPLE

- 基本示例。不带有任何参数时，curl 就是发出 GET 请求。命令向 www.baidu.com 发出 GET 请求，服务器返回的内容会在命令行输出。
  - curl https://www.baidu.com
- 使用 curl 访问 IPv6
  - curl -g -6 'http://\[2408:8210:3c3c:35e0:7df1:783c:ce23:e958]:8080'
  - curl --ipv6 'http://2408:8210:3c3c:35e0:7df1:783c:ce23:e958:8080'
- 下载指定的件
  - curl -LO https://github.com/goharbor/harbor/releases/download/v1.9.3/harbor-online-installer-v1.9.3.tgz
- 请求一个域名时，指定要解析的 IP
  - curl http://myapp.example.com/myapp --resolv 'myapp.example.com:80:172.19.42.217'
- 修改请求头的 Head 信息来发送请求
  - curl -v -H"Host: gw-test.wisetv.com.cn" http://10.10.100.116/app-node/monitor
- 通过文件下载多个 URL
  - xargs -n 1 curl -O < wenjianlisturls.txt # 从 wenjianlisturls.txt 中的 url 列表下载文件
- 不去验证目标证书直接获取 /healthz
  - curl --insecure https://localhost:6443/healthz
- 通过 docker 的 socket 文件获取容器信息
  - **curl --unix-socket /var/run/docker.sock http://localhost/containers/json**
- 使用或不使用身份验证将文件上载到 FTP 服务器。要使用 curl 将名为 wodewenjian.tar.gz 的本地文件上载到 ftp://ftpserver，请执行以下操作：
  - **curl -u username:password -T wodewenjian.tar.gz ftp://ftpserver**
- 存储 Cookie。使用以下命令将它们保存到 linuxidccookies.txt。然后，您可以使用 cat 命令查看该文件。
  - **curl --cookie-jar linuxidcookies.txt https://www.linuxidc.com/index.htm -O**
- 使用 Cookie 发起请求。
  - **curl --cookie cnncookies.txt https://www.linuxidc.com**
- 使用 -d 参数以后，HTTP 请求会自动加上标头 Content-Type : application/x-www-form-urlencoded。并且会自动将请求转为 POST 方法，因此可以省略 -X POST。
  - curl -d'login=emma＆password=123'-X POST https://google.com/login
  - curl -d 'login=emma' -d 'password=123' -X POST https://google.com/login

## -d, --data 选项示例

- 读取 data.txt 文件的内容，作为请求体向服务器发送。
  - curl -d '@data.txt' https://google.com/login
- --data-urlencode 参数等同于 -d，发送 POST 请求的数据体，区别在于会自动将发送的数据进行 URL 编码。下面代码中，发送的数据 hello world 之间有一个空格，需要进行 URL 编码。
  - curl --data-urlencode 'comment=hello world' https://google.com/login

## -w, --write-out 选项示例

- 获取请求 www.baidu.com 总共花费的时间
  - curl -o /dev/null -s -w '%{time_total}\n' https://www.baidu.com
- 只显示响应的状态码。
  - curl -s -o /dev/null -w %{http_code}"\n" http://www.baidu.com

# 复杂应用实例

## 从 Chrome 浏览器的请求中，获取 curl 参数

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fwy2as/1616165084827-164a7bfc-f35f-43a4-8580-9a8eb3856784.png)

# 可能是东半球最好的 Curl 学习指南

-A

-A 参数指定客户端的用户代理标头，即 User-Agent。curl 的默认用户代理字符串是 curl/\[version]。

$ curl -A 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36' <https://google.com>

上面命令将 User-Agent 改成 Chrome 浏览器。

$ curl -A '' <https://google.com>

上面命令会移除 User-Agent 标头。你也可以通过 -H 参数直接指定标头，更改 User-Agent。

$ curl -H 'User-Agent: php/1.0' <https://google.com>

-b

-b 参数用来向服务器发送 Cookie。

$ curl -b 'foo=bar' <https://google.com>

上面命令会生成一个标头 Cookie: foo=bar，向服务器发送一个名为 foo、值为 bar 的 Cookie。

$ curl -b 'foo1=bar' -b 'foo2=baz' <https://google.com>

上面命令发送两个 Cookie。

$ curl -b cookies.txt <https://www.google.com>

上面命令读取本地文件 cookies.txt，里面是服务器设置的 Cookie（参见 -c 参数），将其发送到服务器。

-c

-c 参数将服务器设置的 Cookie 写入一个文件。

$ curl -c cookies.txt <https://www.google.com>

上面命令将服务器的 HTTP 回应所设置 Cookie 写入文本文件 cookies.txt。

-e

-e 参数用来设置 HTTP 的标头 Referer，表示请求的来源。

$ curl -e 'https://google.com?q=example' <https://www.example.com>

上面命令将 Referer 标头设为 <https://google.com?q=example>。

-H 参数可以通过直接添加标头 Referer，达到同样效果。

$ curl -H 'Referer: <https://google.com?q=example'> <https://www.example.com>

-F

-F 参数用来向服务器上传二进制文件。

$ curl -F 'file=@photo.png' <https://google.com/profile>

上面命令会给 HTTP 请求加上标头 Content-Type: multipart/form-data，然后将文件 photo.png 作为 file 字段上传。

-F 参数可以指定 MIME 类型。

$ curl -F 'file=@photo.png;type=image/png' <https://google.com/profile>

上面命令指定 MIME 类型为 image/png，否则 curl 会把 MIME 类型设为 application/octet-stream。

-F 参数也可以指定文件名。

$ curl -F 'file=@photo.png;filename=me.png' <https://google.com/profile>

上面命令中，原始文件名为 photo.png，但是服务器接收到的文件名为 me.png。

-G

-G 参数用来构造 URL 的查询字符串。

$ curl -G -d 'q=kitties' -d 'count=20' <https://google.com/search>

上面命令会发出一个 GET 请求，实际请求的 URL 为 <https://google.com/search?q=kitties&count=20>。如果省略 --G，会发出一个 POST 请求。

如果数据需要 URL 编码，可以结合 --data--urlencode 参数。

$ curl -G --data-urlencode 'comment=hello world' <https://www.example.com>

-H

-H 参数添加 HTTP 请求的标头。

$ curl -H 'Accept-Language: en-US' <https://google.com>

上面命令添加 HTTP 标头 Accept-Language: en-US。

$ curl -H 'Accept-Language: en-US' -H 'Secret-Message: xyzzy' <https://google.com>

上面命令添加两个 HTTP 标头。

$ curl -d '{"login": "emma", "pass": "123"}' -H 'Content-Type: application/json' <https://google.com/login>

上面命令添加 HTTP 请求的标头是 Content-Type: application/json，然后用 -d 参数发送 JSON 数据。

-i

-i 参数打印出服务器回应的 HTTP 标头。

$ curl -i <https://www.example.com>

上面命令收到服务器回应后，先输出服务器回应的标头，然后空一行，再输出网页的源码。

-s

-s 参数将不输出错误和进度信息。

$ curl -s <https://www.example.com>

上面命令一旦发生错误，不会显示错误信息。不发生错误的话，会正常显示运行结果。

如果想让 curl 不产生任何输出，可以使用下面的命令。

$ curl -s -o /dev/null <https://google.com>

-S

-S 参数指定只输出错误信息，通常与 -s 一起使用。

$ curl -s -o /dev/null <https://google.com>

上面命令没有任何输出，除非发生错误。

-u

-u 参数用来设置服务器认证的用户名和密码。

$ curl -u 'bob:12345' <https://google.com/login>

上面命令设置用户名为 bob，密码为 12345，然后将其转为 HTTP 标头 Authorization: Basic Ym9iOjEyMzQ1。

curl 能够识别 URL 里面的用户名和密码。

$ curl <https://bob:12345@google.com/login>

上面命令能够识别 URL 里面的用户名和密码，将其转为上个例子里面的 HTTP 标头。

$ curl -u 'bob' <https://google.com/login>

上面命令只设置了用户名，执行后，curl 会提示用户输入密码。

-x

-x 参数指定 HTTP 请求的代理。

$ curl -x socks5://james:cats@myproxy.com:8080 <https://www.example.com>

上面命令指定 HTTP 请求通过 myproxy.com:8080 的 socks5 代理发出。

如果没有指定代理协议，默认为 HTTP。

$ curl -x james:cats@myproxy.com:8080 <https://www.example.com>

上面命令中，请求的代理使用 HTTP 协议。
