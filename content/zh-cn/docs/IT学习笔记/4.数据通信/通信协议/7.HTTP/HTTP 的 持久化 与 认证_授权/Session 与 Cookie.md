---
title: Session 与 Cookie
---

Cookie

背景

- HTTP 是无状态协议，服务器不能记录浏览器的访问状态，也就是说服务器不能区分两次请求是否由同一个客户端发出

- **Cookie(小甜饼) **就是解决 HTTP 协议无状态的方案之一

- Cookie 实际上就是服务器保存再浏览器上的一段信息。浏览器有了 Cookie 之后，每次向服务器发送请求时都会同时将该信息发送给服务器，服务器收到请求后，就可以根据该信息处理请求

- Cookie 由服务器创建，并发送给浏览器，最终由浏览器保存

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168269-5454b393-dbe3-4518-a3e0-480ab5358176.png)

Cookie 的用途

- 保持用户登陆状态，由于不够安全，现在有其他方式替代

- 京东未登录的状态下，使用 Cookie 存储购物车中的物品的

> 淘宝不是这么实现的，淘宝必须登录才能浏览详细商品

- 上一次连接时打开的页面

- 与某个账号关联

- 等等

## Cookie 的属性

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1623676019731-6e1d1e07-4c5a-4d22-a396-32d3e541c07f.png)

一个 cookie 将会具有如下字段：

- **name** # cookie 的名称

- **value **# cookie 的值

- **domain** # 可以访问此 cookie 的域名

- **path **# 可以访问此 cookie 的页面路径。比如 domain 是 desistdaydream.ltd，path 是 /cookie，那么只有访问 <http://desistdaydream.ltd/cookie> 路径下的页面时，才可以读取此 cookie

- **MaxAge **或 **Expires** # 设置 cookie 持久化时的过期时长

  - 注意：Expires 是老式的过期方法， 如果可以，应该使用 MaxAge 设置过期时间，但有些老版本的浏览器不支持 MaxAge。 如果要支持所有浏览器，要么使用 Expires，要么同时使用 MaxAge 和 Expires。

- **size** # cookie 的大小

- **httpOnly** # 是否允许别人通过 js 获取自己的 cookie
  - httpOnly 属性限制了 cookie 对 HTTP 请求的作用范围。特别的，该属性指示用户代理忽略那些通过"非 HTTP" 方式对 cookie 的访问（比如浏览器暴露给 js 的接口）。
- **secure** # 是否只能通过 https 访问

注意：

- HttpOnly 属性和 Secure 属性相互独立：一个 cookie 既可以是 HttpOnly 的也可以有 Secure 属性。 在前段时间的项目中我就用 js 去读取一个 cookie，结果怎么都取不到这个值，最后查证这个 cookie 是 httpOnly 的，花了近 2 个小时，悲剧了。

Cookie 的缺点

- 不安全，明文

- 增加带宽消耗

- 可以在客户端手动禁用

- Cookie 是有上限的，最大 4096 字节

Session

背景

Cookie 虽然在一定程度上解决了 "保持状态" 的需求，但是由于 Cookie 本身最大支持 4096 字节，以及 Cookie 本身保存在客户端，可能被拦截或窃取，因此就需要有一种新的东西，它能支持更多的字节，并且他保存在服务器，有较高的安全性，这，就是 **Session(会话)**。

但是这时，问题又来了，基于 HTTP 协议的无状态特征，服务器根本不知道访问者是谁。又如何保存呢？此时，Cookie 又来了，起到了一个桥接的作用。

用户登录成功之后，我们在服务端为每个用户创建一个特定的 **SessionData **和 **SessionID**，他们一一对应。其中：

- SessionData 是在服务端保存的一个数据结构，用来跟踪用户的状态，这个数据可以保存在集群、数据库、文件中

  - 这个 SessionData 的数据结构一般都是 KEY/VALUE 类型的结构，相当于一个 **大 map。**可以保存在内存中、关系型数据库、Rdis、文件、等等地方

- SessionID 作为唯一标识符，通常会写入到用户的 Cookie 中。

所以，**Session 必须依赖于 Cookie 才能使用**，生成一个 SessionID 放在 Cookie 里传给客户端即可。

在互联网早期，访问一个网站为了持久化，将用户名和密码放在 cookie 里，一条信息就对应一条 cookie。

比如访问 baidu.com 后，浏览器记录了这么几个 cookie：

- username: DesistDaydream

- password: mypassword

- XXX: XXXX

- YYY: YYY

- .....等等

而使用了 Session 之后，Cookie 不再记录这些敏感信息，只保存一个 ID 用来标识这个 Session：

- sessionID: DesistDaydream

剩下的信息一般都保存在服务器本地，根据 sessionID 找到对应信息即可

Session 逻辑

这是一个简单的 Session 处理请求的逻辑

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168291-de47fffb-605b-41a5-ae65-aa3004c78762.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168320-400fd569-42fc-41ba-a439-6556a456a72e.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168293-4d7ef534-5772-4f15-8eaf-ccf0aa43f1db.png)

## Session 设计思路

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168286-70f4c25f-49bd-4b84-8b50-f0bca0900dd1.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qpvk57/1616161168270-49a7f14e-1cf0-4e50-8ddc-3cc8d9537544.png)

Session 本质上是一个 K/V ，通过 key 进行增删改查。可以存储在 RAM 或者 Redis 中。

- 一般有多个 Session，所以需要一个管理器统一管理

- Session 一般作为中间件，所以需要暴露结构给其他代码使用

Session 接口设计

- Set()

- Get()

- Del()

- Save() # 持久存储

SessionMgr 接口设计

- Init() # 初始化，加载 RAM 或 Redis&#x20;

- CreateSession() # 创建一个新的 Session

- GetSession() # 通过 SessionID 获取对应的 Session 对象

RAM Session 设计

- 定义 RAM Session 对象

  - SessionID

  - 存 K/V 的 map

  - 读写锁

- 构造函数，为了获取对象

  - Set()

  - Get()

  - Del()

  - Save()

RAM SessionMgr 设计

- 定义 Memory SessionMgr 对象

  - 存放所有 Session 的 map

  - 读写锁

- 构造函数

  - Init()

  - CreateSession()

  - GetSession()

Redis Session 设计

- 定义 RedisSession 对象

  - SessionID

  - 存 K/V 的 map

  - 读写锁

  - Redis 连接池

  - 记录内存中 map 是否被修改的编辑

- 构造函数

  - Set()

  - Get()

  - Del()

  - Save()

Redis SessionMgr 设计

- 定义 Redis SessionMgr 对象

  - Redis 地址、密码、连接池、读写锁

  - 存放所有 Session 的 map

  - 读写锁

- 构造函数

  - Init()

  - CreateSession()

  - GetSession()
