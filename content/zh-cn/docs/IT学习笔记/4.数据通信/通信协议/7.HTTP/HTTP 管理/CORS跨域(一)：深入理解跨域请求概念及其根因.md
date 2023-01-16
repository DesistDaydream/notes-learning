---
title: CORS跨域(一)：深入理解跨域请求概念及其根因
---

**原文链接：**<https://mp.weixin.qq.com/s/dynx7wrSINYFKZgGPcD3zQ>

> **❝Talk is cheap. Show me the money.**> ❞

# 前言

你好，我是**YourBatman**。
做 Web 开发的小伙伴对“跨域”定并不陌生，像狗皮膏药一样粘着几乎每位同学，对它可谓既爱又恨。跨域请求之于创业、小型公司来讲是个头疼的问题，因为这类企业还未沉淀出一套行之有效的、统一的解决方案。
让人担忧的是，据我了解不少程序员同学（不乏有高级开发）碰到跨域问题大都一头雾水：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182917-3c20dcea-66db-4950-bd08-2ca996afbd43.webp)然后很自然的 用谷歌去百度一下搜索答案，但相关文章可能参差不齐、鱼龙混杂。短则半天长则一天（包含改代码、部署等流程）此问题才得以解决，一个“小小跨域”问题成功偷走你的宝贵时间。
既然跨域是个如此常见（特别是当下前后端分离的开发模式），因此深入理解 CORS 变得就异常的重要了（反倒前端工程师不用太了解），因此早在 2019 年我刚开始写博客那会就有过较为详细的系列文章：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182890-838bfbb6-192c-4666-8af4-a89e80a8aad6.png)现在把它搬到公众号形成技术专栏，并且加点料，让它更深、更全面、更系统的帮助到你，希望可以助你从此不再怕 Cors 跨域资源共享问题。

## 本文提纲

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182696-51d2c077-47b0-4137-b399-a8bc3183bb21.png)

## 版本约定

- JDK：8
- Servlet：4.x

# 正文

文章遵循一贯的风格，本文将采用概念 + 代码示例的方式，层层递进的进行展开叙述。那么上菜，先来个示例预览，模拟一下**跨域请求**，后面的一些的概念示例将以此作为抓手。

## 模拟跨域请求

要模拟跨域请求的根本是需要**两个源**：让请求的来源和目标源不一样。这里我就使用 IDEA 作为静态 Web 服务器（63342），Tomcat 作为后端动态 Servlet 服务器（8080）。

> **❝**> 说明：服务器都在本机，端口不一样即可
> ❞

### 前端代码

index.html

    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>CORS跨域请求</title>
        <!--导入Jquery-->
        <script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.js"></script>
    </head>
    <body>
    <button id="btn">跨域从服务端获取内容</button>
    <div id="content"></div>
    <script>
        $("#btn").click(function () {
            // 跨域请求
            $.get("http://localhost:8080/cors", function (result) {
                $("#content").append(result).append("<br/>");
            });
            // 同域请求
            $.get("http://localhost:63342");
            $.post("http://localhost:63342");
        });
    </script>
    </body>
    </html>

使用 IDEA 作为静态 web 服务器，浏览器输入地址即可访问（注：端口号为 63342）：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183029-7a7f9241-b8f1-43b1-830a-55e244ea8427.png)

### 后端代码

后端写个 Servlet 来接收 cors 请求

    /**
     * 在此处添加备注信息
     *
     * @author YourBatman. <a href=mailto:yourbatman@aliyun.com>Send email to me</a>
     * @site https://yourbatman.cn
     * @date 2021/6/9 10:36
     * @since 0.0.1
     */
    @Slf4j
    @WebServlet(urlPatterns = "/cors")
    public class CorsServlet extends HttpServlet {
        @Override
        protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
            String requestURI = req.getRequestURI();
            String method = req.getMethod();
            String originHeader = req.getHeader("Origin");
            log.info("收到请求：{}，方法：{}， Origin头：{}", requestURI, method, originHeader);
            resp.getWriter().write("hello cors...");
        }
    }

启动后端服务器，点击页面上的**按钮**，结果如下：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183029-6d0cd299-4e43-4c40-b451-8a036ce30019.png)服务端控制台输出：

    ... INFO  c.y.cors.servlet.CorsServlet - 收到请求：/cors，方法：GET， Origin头：http://localhost:63342

> **❝**> 服务端输出日志，说明即使前端的 Http Status 是 error，但服务端还是收到**并处理了**这个请求的
> ❞
> 下面以此代码示例为基础，普及一下和 Cors 跨域资源共享相关的概念。

## Host、Referer、Origin 的区别

这哥三看起来很是相似，下面对概念作出区分。![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182992-b5ba405f-82dd-4071-a854-b6d62a458711.png)

- **Host**：去哪里。**域名+端口**。值为客户端将要访问的远程主机，浏览器在发送 Http 请求时会带有此 Header
- **Referer**：来自哪里。**协议+域名+端口+路径+参数**。当前请求的**来源页面**的地址，服务端一般使用 Referer 首部识别访问来源，可能会以此进行统计分析、日志记录以及缓存优化等
  - 来源页面协议为 File 或者 Data URI（如页面从本地打开的）
  - 来源页面是 Https，而目标 URL 是 http
  - 浏览器地址栏直接输入网址访问，或者通过浏览器的书签直接访问
  - 使用 JS 的 location.href 跳转
  - ...
  - 常见应用场景：百度的搜索广告就会分析 Referer 来判断打开站点是从百度搜索跳转的，还是直接 URL 输入地址的
  - 一般情况下浏览器会带有此 Header，但这些 case 不会带有 Referer 这个头
- **Origin**：来自哪里（跨域）。**协议+域名+端口**。它用于 Cors 请求和同域 POST 请求

可以看到 Referer 与 Origin 功能相似，前者一般用于统计和阻止盗链，后者用于 CORS 请求。但是还是有几点不同：

1. 只有**跨域请求**，或者**同域时发送 post**请求，才会携带 Origin 请求头；而 Referer 只要浏览器能获取到都会携带（除了上面说明的几种 case 外）![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183074-01c3b097-a969-4628-8afc-1b048c8cab82.png)

2. 若浏览器不能获取到请求源页面地址（如上面的几种 case），Referer 头不会发送，但 Origin 依旧会发送，只是值是 null 而已（注：虽然值为 null，但此请求依旧属于 Cors 请求哦），如下图所示：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183177-3c056e2d-4e67-4242-995e-a946a5f4692e.png)

3. Origin 的值只包括**协议、域名和端口**，而 Rerferer 不但包括协议、域名、端口还包括路径，参数，注意不包括 hash 值

## 浏览器的同源策略

浏览器的职责是展示/渲染 document、css、script 脚本等，但是这些资源（将 document、css、script 统一称为资源）可能来自不同的地方，如本地、远程服务器、甚至黑客的服务器......浏览器作为万维网的入口，是我们接入互联网最重要的软件之一（甚至没有之一），因此它的安全性显得尤为重要，这就出现了浏览器的同源策略。
**同源策略**是浏览器一个重要的安全策略，它用于限制一个 origin 源的 document 或者它加载的脚本如何能与另一个 origin 源的资源进行交互。它能帮助阻隔恶意文档，**减少**（并不是杜绝）可能被攻击的媒介。

> **❝**> 方便和安全往往是相悖的：安全性增高了，方便性就会有所降低
> ❞
> 那么问题来了，什么才算同源？

### 同源的定义

URL 被称作：统一资源定位符，同源是针对 URL 而言的。一个完整的 URL 各部分如下图所示：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182773-6d283640-a36c-47f5-891b-35544f3429f1.png)

> **❝**> Tips：域名和 host 是等同的概念，域名+端口号 = host+端口号（大部分情况下你看到域名并没有端口号，那是采用了默认端口号 80 而已）
> ❞
> **同源：只和上图的前两部分（protocol + domain）有关，规则为：全部相同则为同源**。这个定义不难理解，但有几点需要再强调一下：

- 两部分必须**完全一样**才算同源
- 这里的 domain 包含 port 端口号，所以总共是两部分而非三部分
  - 当然也有说三部分的（协议+host+port），理解其含义就成

下面通过举例来彻底了解下。譬如，我的源 URL 为：`[http://www.baidu.com/api/user](http://www.baidu.com/api/user)`，下面表格描述了不同 URL 的各类情况：

| URL                                                                          | 是否同源 | 原因说明                              |
| :--------------------------------------------------------------------------- | :------- | :------------------------------------ |
| `[http://www.baidu.com/account](http://www.baidu.com/account)`               | 是       | 前两部分相同，path 路径不一样而已     |
| `[http://www.baidu.com/account?name=a](http://www.baidu.com/account?name=a)` | 是       | 前两部分相同，path 路径、参数不同而已 |
| `[https://www.baidu.com/api/user](https://www.baidu.com/api/user)`           | 否       | 协议不同                              |
| `[http://www.baidu.com:8080/api/user](http://www.baidu.com:8080/api/user)`   | 否       | 端口不同(domain 不同)                 |
| `[http://map.baidu.com/api/user](http://map.baidu.com/api/user)`             | 否       | host 不同(domain 不同)                |

### 不同源的网络访问

浏览器同源策略的存在，限制了不同源之间的交互，实为不便。但是浏览器也开了一些“绿灯”，让其不受同源策略的约束。此种情况一般可分为如下三类：

1. 跨域写操作（Cross-origin writes）：一般是被允许的。如链接（如 a 标签）、重定向以及表单提交（如 form 表单的提交）
2. 跨域资源嵌入（Cross-origin embedding）：一般是允许的。比如下面这些例子：
   1. `<script src="..."></script>`标签嵌入 js 脚本
   2. `<link rel="stylesheet" href="...">`标签嵌入 CSS
   3. `<img>`展示的图片
   4. `<video>`和`<audio>`媒体资源
   5. `<object>、 <embed> 、<applet>`嵌入的插件
   6. CSS 中使用`@font-face`引入字体
   7. 通过`<iframe>`载入资源
3. 跨域读操作（Cross-origin reads）：一般是**不被允许**的。比如我们的 http 接口请求等都属于此范畴，也是本专栏关注的焦点

简单总结成一句话：**浏览器自己**是可以发起跨域请求的（比如 a 标签、img 标签、form 表单等），但是 Javascript 是不能去跨域获取资源（如 ajax）。

### 如何允许不同源的网络访问

上面说到的第三种情况：跨域读操作一般是不允许跨域访问的，而这种情况是我们开发过程中最关心、最常见的 case，因此必须解决。

> **❝**> Tips：这里的读指的是广义上的读，指的是从服务器获取资源（有 response）的都叫读操作，而和具体是什么 Http Method 无关。换句话讲，所有的 Http API 接口请求都在这里都指的是读操作
> ❞
> 可以使用 CORS 来允许跨源访问。**CORS 是 HTTP 的一部分**，它允许服务端来指定哪些主机可以从这个服务端加载资源。

## 什么是 Cors 跨域

Cors(Cross-origin resource sharing)：跨域资源共享，它是**浏览器**的一个技术**规范**，由 W3C 规定，规范的 wiki 地址在此：<https://www.w3.org/wiki/CORS_Enabled#What_is_CORS_about.3F>

> **❝**> 话外音：它是浏览器的一种（自我保护）行为，并且已形成规范。也就是说：backend 请求 backend 是不存在此现象的喽
> ❞
> 若想实现 Cors 机制的跨域请求，是需要浏览器和服务器同时支持的。关于浏览器对 CORS 的支持情况：现在都 2021 年了，so 可以认为 100%的浏览器都是支持的，再加上 CORS 的整个过程都由浏览器自动完成，**前端无需做任何设置**，所以前端工程师的 ajax 原来怎么用现在还是怎么用，它对前段开发人员是完全透明的。

### 为何需要 Cors 跨域访问？

浏览器费尽心思的搞个同源策略来保护我们的安全，但为何又需要跨域来打破这种安全策略呢？其实啊，这一切都和互联网的快速发展有关~
随着 Web 开放的程度越来越高，页面的内容也是越来越丰富。因此页面上出现的元素也就越来越多：图片、视频、各种文字内容等。为了分而治之，一个页面的内容可能来自不同地方，也就是不同的 domain 域，因此通过 API 跨域访问成了**必然**。
浏览器作为进入 Internet 最大的入口，很长时间它是个大互联公司的必争之地，因此市面上并存的浏览器种类繁多且鱼龙混扎：IE 7、8、9、10，Chrome、Safari、火狐，每个浏览器对跨域的实现可能都不一样。因此对开发者而言亟待需要一个规范的、统一方案，它就是`Cors`。
**CORS（Cross-Origin Resource Sharing）由 W3C 组织于 2009-03-17 编写工作草案，直到 2014-01-16 才正式毕业成为行业规范，所有浏览器得以遵守**。至此，程序员同学们在解决跨域问题上，只需按照 Cors 规范实施即可。

### Cors 的工作原理

Web 资源涉及到两个角色：浏览器（消费者）和服务器（提供者），面向这两个角色来了解 Cors 的原理非常简单，如下图所示：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183135-6470a44e-03d9-49d7-af38-910b998f158c.png)

1. 若浏览器发送的是个跨域请求，http 请求中就会携带一个名为**Origin**的头表明自己的“位置”，如`Origin: http://localhost:5432`
2. 服务端接到请求后，就可以根据传过来的 Origin 头做逻辑，决定是否要将资源共享给这个源喽。而这个决定通过响应头**Access-Control-Allow-Origin**来承载，它的 value 值可以是任意值，有如下情况：
   1. 值为`*`，通配符，允许所有的 Origin 共享此资源
   2. 值为 http://localhost:5432（也就是和 Origin 相同），共享给此 Origin
   3. 值为非 http://localhost:5432（也就是和 Origin 不相同），不共享给此 Origin
   4. 无此头：不共享给此 origin
   5. 有此头：值有如下可能情况
3. 浏览器接收到 Response 响应后，会去提取 Access-Control-Allow-Origin 这个头。然后根据上述规则来决定要接收此响应内容还是拒绝

> **❝**> Tips：Access-Control-Allow-Origin 响应头只能有 1 个，且 value 值就是个字符串。另外，value 值即使写为`[http://aa.com,http://bb.com](http://aa.com,http://bb.com)`这种也属于一个而非两个值
> ❞

### Cors 细粒度控制：授权响应头

在 Cors 规范中，除了可以通过**Access-Control-Allow-Origin**响应头来对主体资源（URL 级别）进行授权外，还提供了针对于具体响应头更细粒度的控制，这个响应头就是：**Access-Control-Expose-Headers**。换句话讲，该头用于规定哪些响应头（们）可以暴露给前端，默认情况下这 6 个响应头无需特别的显示指定就支持：

- Cache-Control
- Content-Language
- Content-Type
- Expires
- Last-Modified
- Pragma

若不在此值里面的头将**不会返回给前端**（其实返回了，只是浏览器让其对前端不可见了而已，对 JavaScript 也不可见哦）。
**但是，但是，但是**，这种细粒度控制 header 的机制对简单请求是无效的，只针对于非简单请求（也叫复杂请求）。由此可见，将哪些类型的跨域资源请求划分为简单请求的范畴就显得特备重要了。

### 何为简单请求

Cors 规范定义简单请求的原则是：请求不是以更新（添加、修改和删除）资源为目的，服务端对请求的处理不会导致自身维护资源的改变。对于简单跨域资源请求来说，浏览器将两个步骤（取得授权和获取资源）合二为一，由于不涉及到资源的改变，所以不会带来任何副作用。
对于一个请求，必须**同时符合**如下要求才被划为简单请求：

1. Http Method 只能为其一：
   1. GET
   2. POST
   3. HEAD
2. 请求头只能在如下范围：
   1. application/x-www-form-urlencoded
   2. multipart/form-data
   3. text/plain
   4. Accept
   5. Accept-Language
   6. Content-Language
   7. **Content-Type**，其中它的值必须如下其一：

除此之外的请求都为非简单请求（也可称为复杂请求）。非简单请求可能对服务端资源改变，因此 Cors 规定浏览器在发出此类请求**之前**必须有一个“预检（Preflight）”机制，这也就是我们熟悉的`OPTIONS`请求。

### 什么是 Preflight 预检机制

顾名思义，它表示在浏览器发出**真正**请求**之前**，先发送一个预检请求，这个在 Http 里就是 OPTIONS 请求方式。这个请求很特殊，它不包含主体（无请求参数、请求体等），主要就是将一些凭证、授权相关的辅助信息放在请求头里交给服务器去做决策。因此它除了携带 Origin 请求头外，还会额外携带如下两个**请求头**：

- **Access-Control-Request-Method**：真正请求的方法
- **Access-Control-Request-Headers**：真正请求的**自定义**请求头（若没有自定义的就是空呗）

服务端在接收到此类请求后，就可以根据其值做逻辑决策啦。如果允许预检请求通过，返回个 200 即可，否则返回 400 或者 403 呗。
如果预检成功，在响应里应该包含上文提到的响应头**Access-Control-Allow-Origin**和**Access-Control-Expose-Headers**，除此之外，服务端还可以做更精细化的控制，这些精细化控制的响应头为：

- **Access-Control-Allow-Methods**：允许实际请求的 Http 方法（们）
- **Access-Control-Allow-Headers**：允许实际请求的请求头（们）
- **Access-Control-Max-Age**：允许浏览器缓存此结果多久，单位：**秒**。有了缓存，以后就不用每次请求都发送预检请求啦

> **❝**> 说明：以上响应头并不是必须的。若没有此响应头，代表接受所有
> ❞
> ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183047-fe0cff44-2fad-43ff-9b1d-418e6cc8bbd2.webp)预检请求完成后，有个关键点，便是浏览器拿到预检请求的响应后的处理逻辑，这里描述如下：

1. 先通过自己的 Origin 匹配预检响应中的**Access-Control-Allow-Origin**的值，若不匹配就结束请求，若匹配就继续下一步验证
   1. 关于 Access-Control-Allow-Origin 的验证逻辑，请参考文上描述
2. 拿到预检响应中的**Access-Control-Allow-Methods**头。若此头不存在，则进行下一步，若存在则校验预检请求头 Access-Control-Request-Method 的值是否在此列表中，在其内继续下一步，否则失败
3. 拿到预检响应中的**Access-Control-Request-Headers**头。同请求头中的**Access-Control-Allow-Headers**值记性比较，**全部包含在内**则匹配成功，否则失败

以上全部匹配成功，就代表预检成功，可以开始发送正式请求了。值得一提的事，Access-Control-Max-Age 控制预检结果的浏览器缓存，若缓存还生效的话，是不用单独再发送 OPTIONS 请求的，匹配成功直接发送目标真实即可。

#### Access-Control-Max-Age 使用细节

Access-Control-Max-Age 用于控制浏览器缓存预检请求结果的时间，这里存在一些使用细节你需要注意：

1. 若浏览器禁用了缓存，也就是勾选了`Disable cache`，那么此属性无效。也就说每次都还得发送 OPTIONS 请求
2. 判断此缓存结果的因素有两个：
   1. 必须是同一 URL（也就是 Origin 相同才会去找对应的缓存）
   2. header 变化了，也会重新去发 OPTIONS 请求（当然若去掉一些 header 编程简单请求了，就另当别论喽）

## 跨域请求代码示例

正所谓说再多，也抵不上跑几个 case，毕竟 show me your code 才是最重要。下面就针对跨域情况的简单请求、非简单请求（预检通过、预检不通过）等 case 分别用代码（基于文首代码）说明。

### 简单请求

简单请求正如其名，是最简单的请求方式。

    // 跨域请求
    $.get("http://localhost:8080/cors", function (result) {
        $("#content").append(result).append("<br/>");
    });

服务端结果：

    INFO ...CorsServlet - 收到请求：/cors，方法：GET， Origin头：http://localhost:63342

浏览器结果：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183182-8f80baa8-d029-42c4-8bce-a29833a2f28d.png)若想让请求正常，只需在服务端响应头里“加点料”就成：

    ...
    resp.setHeader("Access-Control-Allow-Origin","http://localhost:63342");
    resp.getWriter().write("hello cors...");
    ...

再次请求，结果成功：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182990-03d5ab7f-9e0f-42a9-852e-7f1d0f6d6430.png)
对于简单请求来讲，服务端只需要设置**Access-Control-Allow-Origin**这个一个头即可，一个即可。

### 非简单请求

非简单请求的模拟非常简单，随便打破一个简单请求的约束即可。比如我们先在上面 get 请求的基础上自定义个请求头：

    $.ajax({
        type: "get",
        url: "http://localhost:8080/cors",
        headers: {secret:"kkjtjnbgjlfrfgv",token: "abc123"}
    });

服务端代码：

    /**
     * 在此处添加备注信息
     *
     * @author YourBatman. <a href=mailto:yourbatman@aliyun.com>Send email to me</a>
     * @site https://yourbatman.cn
     * @date 2021/6/9 10:36
     * @since 0.0.1
     */
    @Slf4j
    @WebServlet(urlPatterns = "/cors")
    public class CorsServlet extends HttpServlet {
        @Override
        protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
            String requestURI = req.getRequestURI();
            String method = req.getMethod();
            String originHeader = req.getHeader("Origin");
            log.info("收到请求：{}，方法：{}， Origin头：{}", requestURI, method, originHeader);
            resp.setHeader("Access-Control-Allow-Origin","http://localhost:63342");
            resp.setHeader("Access-Control-Expose-Headers","token,secret");
            resp.setHeader("Access-Control-Allow-Headers","token,secret"); // 一般来讲，让此头的值是上面那个的【子集】（或相同）
            resp.getWriter().write("hello cors...");
        }
    }

点击按钮，浏览器发送请求，结果为：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260183014-9b0617b6-5b7b-4c38-b7e1-c7fe01a03cdd.png)根本原因为：OPTIONS 的响应头里并未含有任何跨域相关信息，虽然预检通过（注意：这个预检是通过的哟，预检不通过的场景就不用额外演示了吧~），但预检的结果经浏览器判断此跨域实际请求不能发出，所以给拦下来了。
从代码层面问题就出现在`resp.setHeader(xxx,xxx)`放在了处理实际方法的 Get 方法上，显然不对嘛，应该放在`doOptions()`方法里才行：

    @Override
    protected void doOptions(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        super.doOptions(req, resp);
        resp.setHeader("Access-Control-Allow-Origin","http://localhost:63342");
        resp.setHeader("Access-Control-Expose-Headers","token,secret");
        resp.setHeader("Access-Control-Allow-Headers","token,secret"); // 一般来讲，让此头的值是上面那个的【子集】（或相同）
    }

在此运行，一切正常：![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glml9g/1624260182853-033ea229-1f94-47fe-abf1-be70ab519582.webp)值得特别注意的是：设置跨域的响应头这块代码，在处理真实请求的 doGet 里也必须得有，否则服务端处理了，浏览器“不认”也是会出跨域错误的。
另外就是，Access-Control-Allow-Headers/Access-Control-Expose-Headers 这两个头里必须包含你的请求的自定义的 Header（标准的 header 不需要包含），否则依旧跨域失败哦~
在实际生产场景中，Http 请求的 Content-type 大都是`application/json`并非简单请求的头，所以有个现实情况是：**实际的跨域请求中，几乎 100%的情况下我们发的都是非简单请求**。

## Cros 跨域使用展望

如上代码示例，处理简单请求尚且简单，但对于非简单请求来说，我们在 doOptions 和 doGet 都写了一段 setHeader 的代码，是否觉得麻烦呢？
另外，对于**Access-Control-Allow-Origin**若我需要允许多个源怎么办呢？

> **❝**> Tips：Access-Control-Allow-Origin 头只允许一个，且 Access-Control-Allow-Origin:a.com,b.com 依旧算作一个源的，它没有逗号分隔的“特性”。从命名的艺术你也可看出，它并非是 xxx-Origins 而是 xxx-Origin
> ❞
> 既然实际场景中几乎 100%都是非简单请求，那么对于控制非简单请求的**Access-Control-Allow-Methods**、**Access-Control-Allow-Headers**、**Access-Control-Max-Age**这些都都改如何赋值？是否有最佳实践？
> 现在我们大都在 Spring Framework/Spring Boot 场景下开发应用，框架层面是否提供一些优雅的解决方案？
> 作为一名后端开发工程师（编程语言不限），也许你从未处理过跨域问题，那么到底是谁默默的帮你解决了这一切呢？是否想知其所以然？
> 如果这些问题也是你在使用过程中的疑问，或者希望了解的知识点，那么请关注专栏吧。

# 总结

本文用很长的篇幅介绍了 Cors 跨域资源共享的相关知识，并且用代码做了示范，希望能助你通关 Cors 这个狗皮膏药一样粘着我们的硬核知识点。本文文字叙述较多，介绍了同源、跨域、Cors 的几乎所有概念，虽然略显难啃，但这些是指导我们实践的说明书。
革命尚未统一，带着 👆🏻 给到的问题，一起开启通过 Cors 跨域之旅吧~

## 本文思考题

本文已被`[https://yourbatman.cn](https://yourbatman.cn)`收录。所属专栏：**点拨-Cors 跨域**，后台回复“专栏列表”即可查看详情。
**看完了不一定懂，看懂了不一定会**。来，3 个思考题帮你复盘：

1. 试想一下，如果浏览器没有同源策略，将有多大的风险？
2. Cors 共涉及到哪些请求头？哪些响应头？
3. 你所知道的解决 Cors 跨域问题最佳实践是什么？

## 推荐阅读

- [10. 原来是这么玩的，@DateTimeFormat 和@NumberFormat](https://mp.weixin.qq.com/s?__biz=MzI0MTUwOTgyOQ==&mid=2247493900&idx=1&sn=f903380f7d7065192959b07acd1683ae&scene=21#wechat_redirect)
- [9. 细节见真章，Formatter 注册中心的设计很讨巧](https://mp.weixin.qq.com/s?__biz=MzI0MTUwOTgyOQ==&mid=2247491641&idx=1&sn=794c1fadf5d5144df83079a549d3b15f&scene=21#wechat_redirect)
- [8. 格式化器大一统 -- Spring 的 Formatter 抽象](https://mp.weixin.qq.com/s?__biz=MzI0MTUwOTgyOQ==&mid=2247491457&idx=1&sn=dc16a2b7f069df4b0329e66365efd980&scene=21#wechat_redirect)
