---
title: HAProxy
---

# HAProxy

可以实现四层以及七层负载均衡

多用于七层负载均衡

http 层代理

tcp 层的负载均衡

目前，有两种主流的代理模式：tcp 代理(即所谓的 4 层代理)和 http 代理(即所谓的 7 层代理)。在 4 层代理模式下，haproxy 简单的在两端进行双向转发。在 7 层代理模式下，haproxy 会对协议进行分析，可以根据协议来允许、阻塞、切换、增加、修改和移除 request 或 response 中的属性内容。

## haproxy 工作逻辑

比如 client 为 114.114.114.114，haproyx 为 192.168.1.2，Server 为 192.168.1.3

- client 发送数据包给 haproxy 所在服务器 192.168.1.2
- 192.168.1.2 发现这个数据包是给自己的 haproxy 的，则剥离 IP 与 PORT，并把数据包发送给用户空间的 haproxy
- haproxy 由于在用户空间，所以收到的数据包已经被内核剥离了 IP 与 PORT，此时 haproxy 会根据自身的配置以及数据包内的相关信息来进行匹配选择一个合适的 Server，然后发送给内核，告诉内核这个数据包要发送给某 Server。这是 haproxy 与 Server 建立的一个新 TCP 连接。
- 内核根据 Server 这个目的 IP，再封装上 mac 地址从网卡中发送出去。
- 这时候 Server 就会收到请求，处理完成后把响应报文发送给 haproxy
- 由于 haproxy 与 Client 和 Server 分别建立的两个 TCP 连接，这会生成两个 Socket，所以发送给 client 的响应数据以及之后的数据交互就直接通过两个相连的 socket 来进行。
   - socket 介绍详解 TCPIP，UDP，端口 Port，Socket，API.note 的 socket 章节
   - haproxy 所在的设备就相当于创建了两个 socket，这两个 socket 又可以直接相连，这时候数据就不用再次经过用户空间而可以直接交互。client 作为客户端随机使用一个端口与 haproxy 监听的端口交互，之后 haproxy 作为客户端随机使用一个端口与 server 监听的端口交互。效果如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zdm35o/1616132304027-ba04f126-b601-4db8-b244-3b361691fcbd.jpeg)

## haproxy 建立链接的方式

- 比如此时有 5W 个用户请求，haproxy 会与后端服务器建立指定个数的连接，比如 500，这 500 个连接可以让每个连接处理 100 个请求。(并不是一个 Client 与 haproxy 连接，haproxy 就也要与 Server 建立链接，haproxy 两端的连接没有绝对关系。当然，也可以单独进行配置以便让一个 Client 单独享有一个与 Servier 链接的关系)
- 每一个客户端需要与 haproxy 要建立链接，但是只占用 haproxy 所监听的端口，haproxy 与后端服务器建立多少连接，就占用 haproxy 所在服务器多少端口。如果每次后端服务器响应完成之后就关闭连接(比如 http-server-close 配置)，那么如果客户端请求数据大，建立的链接数量也就变大大，很有可能会在端口关闭还没释放前，就提前把端口占满了。
- 所以，haproxy 最佳的性能不一定是一个客户端请求就让 haproxy 建立一个与后端的连接，当设备性能足够的时候，可以在尽量少建立链接的情况下，让单个链接处理多个请求。这个最佳的连接数，可以通过压测工具来测出来峰值，比如 haproxy 与后端连接数为 1、2、3....N 的时候，服务性能在哪个数值可以达到峰值，就可以配置 haproxy 与后端的连接数为多少。

### 注意 lvs 的四层代理与 haproxy 的 7 层代理的区别

- 虽然从逻辑上看似都是修改了请求报文中的某些信息(比如目的 MAC、目的 IP、源 IP 等等)，但是这种逻辑上的修改其实有区别的
   - lvs 在内核空间操作数据包。内核空间的数据包外层还有 mac、ip 地址、port 等，可以由内核模块直接进行修改，然后让数据包不进入用户空间就直接通过网络栈发给网卡。实际上算是一种转发并且是真正的修改
   - haproxy 在用户空间操作数据。而一般情况数据本身并没有 ip、mac 等信息，所以，haproxy 就根据数据内容，找到匹配的 RS，把改数据包当做一个新包发送出去。发送的时候其实就是由 haproxy 发起的新请求。只不过在发送请求的时候，haproxy 还会记录这个请求。这样在 RS 响应该报文时，haproxy 就可以把响应数据再发送给 client。如果从 client 来看，确实是修改了目的 IP(虽然 client 并不知道，但是作为操作者纵观全局来说是修改了)，但是实际是 haproxy 把数据当做自己主动发送的数据再发送给别人了。这不是修改了原始数据包的信息而是代替 client 发送数据与 RS 建立一个新的 TCP 连接，是属于代理的类型，而 lvs 更像是转发。
- 总结起来就是 lvs 是在内核处理数据包、而 haproxy 在用户空间处理数据。所以 haproxy 还有有性能的损耗，因为数据会经由内核空间进到用户空间再进到内核空间发送出去

# HAProxy 关联文件与配置

**/etc/haproxy/haproxy.cfg** # haproxy 程序运行所需基本配置

haproxy.cfg 文件中各个参数解释

官方参数详解见：<http://cbonte.github.io/haproxy-dconv/1.9/configuration.html#4.1>

## 全局配置段

### global # 全局配置

## 代理配置段

- frontend 和 backend 主要用来实现 7 层代理。两个配置段相互配置。用来匹配代理规则的，用户访问哪一类页面，就把这类请求转发到某后端服务器上。可以与 nginx 的配置互相参考，acl 就是 naginx 的 location 里的匹配规则；backend 里的 server 就是 nginx 里 upstream 的 server；use_backend 就是 nginx 里的 proxy_pass。
- listen 主要用来实现 4 层代理。用户访问 haproxy 设备的哪个端口，就会把请求代理到指定的服务器的 IP:PORT 上

frontend、backend、listen 三个配置段都可以有多个，e.g.根据 acl1 的规则匹配到 backend1，符合 acl2 的规则匹配到 backend2，以此类推

### defaults \[DefaultsName] - 为 frontend、backend、listen 三段提供默认配置

defaults 中的配置，也可以单独配置在 frontend、backend、listen 配置段中，该段中的配置主要是为了防止后面面的三段中有重复的配置，比如 option http-keep-alive 该选项是保持连接，需要在所有前端配置中配置，这时候，可以在 defaults 配置段段配置一次，即可在所有的 frontend、backend、listen 中生效。

### frontend FrontendName \*:PORT - 前端配置段

用于让 haproxy 监听在某个 IP:PORT 上，用来接收客户端请求，然后根据匹配规则把改请求转发给指定 backend

### backend BackendName - 后端配置段

负责接收前段配置段转发的请求，后端里可以包含多台服务器来均衡处理前端转发过来的请求。BackendName 与 frontend 中的 default_backend 或 use_backend 所关联。

### listen ListenName - 指定一个四层代理名称，该配置段无法配置详细的 7 层转发规则

## 配置段的关键字(keywords)

### global 全局配置段 keywords

- maxconn NUM # 设置最大连 z 接数
- log \[Len ] \[Format ] \[ \[]] # 让 haproxy 程序记录日志，并指定记录方式
   - ADDRESS # 指定要把日志发送到哪台设备上的哪个 PORT，默认 PORT 为 514
   - FACILITY # 指定要使用的日志设施，一般为 local0-7 其中一个
   - LEVEL # 指定哪个级别的日志会被记录
- user USER # haproxy 以指定的 USER 用户运行
- group GROUP # haproxy 以指定的 GROUP 组运行
- daemon # haproxy 以守护进程运行，不加该参数，则 haproxy 则会运行在前台

### defaults、frontend、backend、listen 代理配置段 keywords

defaults、frontend、backend、listen 可用的关键字及其对应的值的简要说明，其余的配置详见官网参数说明

目前，有两种主流的代理模式：tcp 代理(即所谓的 4 层代理)和 http 代理(即所谓的 7 层代理)。在 4 层代理模式下，haproxy 简单的在两端进行双向转发。在 7 层代理模式下，haproxy 会对协议进行分析，可以根据协议来允许、阻塞、切换、增加、修改和移除 request 或 response 中的属性内容。

- **log global** # 让代理的日志使用 global 中的日志配置
- **mode {tcp|http}** # 设置实例的运行模式或者协议,默认为 tcp
   - 用于 defaults、frontend、listen、backend
   - tcp 模式：该实例将在纯 TCP 模式下工作。 将在客户端和服务器之间建立全双工连接，并且不会执行第 7 层检查。 这是默认模式。 它应该用于 SSL，SSH，SMTP，......
   - http 模式：该实例将在 HTTP 模式下工作。 在连接到任何服务器之前，将对客户端请求进行深入分析。 任何不符合 RFC 的请求都将被拒绝。 可以进行第 7 层过滤，处理和切换。这种模式为 HAProxy 带来了最大的价值。
- **acl CRITERION \[FLAG] \[OPERATOR] VALUE** # 定义一个名为 AclName 的规则，匹配规则为 MatchRule。匹配规则主要是针对访问的内容，i.e.用户访问哪一类 URL。e.g.用户访问某个路径下的资源，acl 匹配到后，把请求代理到 use_backend 字段中定义的 backend 上去.acl 官方使用文档：<http://cbonte.github.io/haproxy-dconv/1.9/configuration.html#7.1>
   - 用于 frontend、backend、listen。
   - EXAMPLE
      - acl invalid_src src 0.0.0.0/7 # CERITERION(规范)为 src，VALUE(值)为 0.0.0.0/7。规则名为 invalid_src
- **bind IP:PORT** # 指定该前端会监听在哪个 IP:PORT 上
   - 用于 frontend、listen
- **use_backend BackendName if ACLName** # 当满足 ACLName 策略的请求代理到 BackendName 这个后端上
   - 用于 frontend、listen。
- **default_backend BackendName** # 定义默认把请求转发到 backend 所定义的一组以 NAME 命名的后端服务器上
   - 用于 defaults、frontend、listen。
- **balance SCHEDULER** # 指明该组后端服务器接收请求的 SCHEDULER(scheduler 调度算法,也可以翻译为调度器)，调度算法可以使用以下的几种机型定义，有的 scheduler 还有子配置，定义在该 scheduler 之下
  - 用于 defaults、listen、backend。
  - 动态：权重可动态调整
  - 静态：调整权重不会实时生效
    - roundrobin # 动态轮询，基于权重，动态调整权重，每个后端设备最多支持 4128 个连接
    - static-rr # 静态轮询，基于权重，该组中的每个设备轮流接收请求，无接收上线
    - leastconn # 根据后端设备的负载数量进行调度，仅适用于长连接的会话
    - source # 对源 IP 地址进行哈希，用可用服务器的权重总数除以哈希值，根据结果进行分配。只要服务器正常，同一个客户端 IP 地址总是访问同一台服务器。如果哈希的结果随可用服务器数量而变化，那么有的客户端会定向到不同的服务器。该算法一般用于不能插入 cookie 的 TCP 模式。它还可以用于广域网上，为拒绝使用会话 cookie 的客户端提供最有效的粘连。该算法默认是静态的，所以运行时修改服务器的权重是无效的，但是算法会根据"hash-type"的变化做调整。将请求的源地址进行 hash 运算，然后进行调度，使用该 scheduler 后可以使用 hash-type 参数来定义 hash 类型
    - uri # 将请求的 uri 左边的部分(左边的部分就是 uri 的语法的？之前的内容，相当于资源的位置)进行 hash 运算后进行调度，常用于缓存服务器(可以根据每个资源的 hash 值，来让用户直接访问到该资源上，第一个用户请求到该资源在某台的时候，后续的用户再请求该资源还是会到该设备上，比如 CDN 中，)；使用该 scheduler 后可以使用 hash-type 参数来定义 hash 类型
    - url_param # 根据 url 中指定的 param 参数的值做 hash 运算后进行调度，可以根据参数中的内容进行会话保持；比如根据 param 中定义的用户名，来规定有哪一部分调度到哪台设备上；使用该 scheduler 后可以使用 hash-type 参数来定义 hash 类型
    - hdr(NAME) # 根据每个请求报文的 header 首部报文的(字段)进行调度
    - hash-type # 动态调整权重，一致性哈希算法 | 静态，基础映射哈希
- **option httpchk HEAD /PATH/TO/FILE** # 用于配置健康检查所使用的文件，HEAD 是关键字
   - 用于 defaults、listen、backend。
- **option http-keep-alive** # 启用客户端和服务端与 haproxy 之间的长连接。haproxy 将处理所有请求和响应报文，请求完后 haproxy 两端的连接都处于空闲状态。
   - 用于 defaults、frontend、backend、listen
- **option http-server-close** # 启用在 haproxy 处理完第一次响应之后关闭 haproxy 到服务端之间长连接的功能，但客户端的长连接还保持，后续的每次请求都重新建立和后端的连接，每次响应后都关闭和后端的连接。启用该选项时，haproxy 将会在发送给后端 server 的 request 数据包中添加一个"Connection:Close"标记，后端 Server 看到此标记就会在响应后关闭 tcp 连接。
  - 用于 defaults、frontend、backend、listen
  - 一般来说，后端是静态内容缓存服务器时，或者就是静态服务器时，首选使用 http-keep-alive 模式，后端是动态应用程序服务器时，首选使用 http-server-close 模式。
- **option httpchk \[METHOD] URI \[VERSION]** # 开启 http 协议以检查 server 字段定义的各个服务器的健康状态
- **server ServerName IP:PORT check inter 5s rise 2 fall 2 weight 2** # 指定被代理的服务器的 IP 与地址还有其余信息
  - 用于 backend 和 listen。
  - 在 backend 里声明 server，这些 server 用来处理前端代理过来的请求。可以定义多个 server 来负载均衡这些请求，每个 server 都可以自定义一个 NAME，这些名字用来在 haproxy 的日志中查看是哪个服务器出问题或者接收请求等

实际应用举例

```
       ####################全局配置信息######################## 
       #######参数是进程级的，通常和操作系统（OS）相关######### 
global 
       maxconn 20480                   #默认最大连接数 
       log 127.0.0.1 local3            #[err warning info debug] 
       chroot /var/haproxy             #chroot运行的路径 
       uid 99                          #所属运行的用户uid，也可以user后边接用户名
       gid 99                          #所属运行的用户组 ，也可以改成group后边接组名
       daemon                          #以后台形式运行haproxy 
       nbproc 1                        #进程数量(可以设置多个进程提高性能) 
       pidfile /var/run/haproxy.pid    #haproxy的pid存放路径,启动进程的用户必须有权限访问此文件 
       ulimit-n 65535                  #ulimit的数量限制 
 
 
       #####################默认的全局设置###################### 
       ##这些参数可以被利用配置到frontend，backend，listen组件## 
defaults 
       log global 
       mode http                       #所处理的类别 (#7层 http;4层tcp  ) 
       maxconn 20480                   #最大连接数 
       option httplog                  #日志类别http日志格式 
       option httpclose                #每次请求完毕后主动关闭http通道 
       option dontlognull              #不记录健康检查的日志信息 
       option forwardfor               #如果后端服务器需要获得客户端真实ip需要配置的参数，可以从Http Header中获得客户端ip  
       option redispatch               #serverId对应的服务器挂掉后,强制定向到其他健康的服务器  
       option abortonclose             #当服务器负载很高的时候，自动结束掉当前队列处理比较久的连接 
       stats refresh 30                #统计页面刷新间隔 
       retries 3                       #3次连接失败就认为服务不可用，也可以通过后面设置 
       balance roundrobin              #默认的负载均衡的方式,轮询方式 
      #balance source                  #默认的负载均衡的方式,类似nginx的ip_hash 
      #balance leastconn               #默认的负载均衡的方式,最小连接 
       contimeout 5000                 #连接超时 
       clitimeout 50000                #客户端超时 
       srvtimeout 50000                #服务器超时 
       timeout check 2000              #心跳检测超时 
 
       ####################监控页面的设置####################### 
listen admin_status                    #Frontend和Backend的组合体,监控组的名称，按需自定义名称 
        bind 0.0.0.0:65532             #监听端口 
        mode http                      #http的7层模式 
        log 127.0.0.1 local3 err       #错误日志记录 
        stats refresh 5s               #每隔5秒自动刷新监控页面 
        stats uri /admin?stats         #监控页面的url 
        stats realm itnihao\ itnihao   #监控页面的提示信息 
        stats auth admin:admin         #监控页面的用户和密码admin,可以设置多个用户名 
        stats auth admin1:admin1       #监控页面的用户和密码admin1 
        stats hide-version             #隐藏统计页面上的HAproxy版本信息  
        stats admin if TRUE            #手工启用/禁用,后端服务器(haproxy-1.4.9以后版本) 
 
 
       errorfile 403 /etc/haproxy/errorfiles/403.http 
       errorfile 500 /etc/haproxy/errorfiles/500.http 
       errorfile 502 /etc/haproxy/errorfiles/502.http 
       errorfile 503 /etc/haproxy/errorfiles/503.http 
       errorfile 504 /etc/haproxy/errorfiles/504.http 
 
       #################HAProxy的日志记录内容设置################### 
       capture request  header Host           len 40 
       capture request  header Content-Length len 10 
       capture request  header Referer        len 200 
       capture response header Server         len 40 
       capture response header Content-Length len 10 
       capture response header Cache-Control  len 8 
     
       #######################网站监测listen配置##################### 
       ###########此用法主要是监控haproxy后端服务器的监控状态############ 
listen site_status 
       bind 0.0.0.0:1081                    #监听端口 
       mode http                            #http的7层模式 
       log 127.0.0.1 local3 err             #[err warning info debug] 
       monitor-uri /site_status             #网站健康检测URL，用来检测HAProxy管理的网站是否可以用，正常返回200，不正常返回503 
       acl site_dead nbsrv(server_web) lt 2 #定义网站down时的策略当挂在负载均衡上的指定backend的中有效机器数小于1台时返回true 
       acl site_dead nbsrv(server_blog) lt 2 
       acl site_dead nbsrv(server_bbs)  lt 2  
       monitor fail if site_dead            #当满足策略的时候返回503，网上文档说的是500，实际测试为503 
       monitor-net 192.168.16.2/32          #来自192.168.16.2的日志信息不会被记录和转发 
       monitor-net 192.168.16.3/32 
 
       ########frontend配置############ 
       #####注意，frontend配置里面可以定义多个acl进行匹配操作######## 
frontend http_80_in 
       bind 0.0.0.0:80      #监听端口，即haproxy提供web服务的端口，和lvs的vip端口类似 
       mode http            #http的7层模式 
       log global           #应用全局的日志配置 
       option httplog       #启用http的log 
       option httpclose     #每次请求完毕后主动关闭http通道，HA-Proxy不支持keep-alive模式 
       option forwardfor    #如果后端服务器需要获得客户端的真实IP需要配置次参数，将可以从Http Header中获得客户端IP 
       ########acl策略配置############# 
       acl itnihao_web hdr_reg(host) -i ^(www.itnihao.cn|ww1.itnihao.cn)$    
       #如果请求的域名满足正则表达式中的2个域名返回true -i是忽略大小写 
       acl itnihao_blog hdr_dom(host) -i blog.itnihao.cn 
       #如果请求的域名满足www.itnihao.cn返回true -i是忽略大小写 
       acl itnihao    hdr(host) -i itnihao.cn 
       #如果请求的域名满足itnihao.cn返回true -i是忽略大小写 
       acl file_req url_sub -i  killall= 
       #在请求url中包含killall=，则此控制策略返回true,否则为false 
       acl dir_req url_dir -i allow 
       #在请求url中存在allow作为部分地址路径，则此控制策略返回true,否则返回false 
       acl missing_cl hdr_cnt(Content-length) eq 0 
       #当请求的header中Content-length等于0时返回true 
       
       ########acl策略匹配相应############# 
       block if missing_cl 
       #当请求中header中Content-length等于0阻止请求返回403 
       block if !file_req || dir_req 
       #block表示阻止请求，返回403错误，当前表示如果不满足策略file_req，或者满足策略dir_req，则阻止请求 
       use_backend  server_web  if itnihao_web 
       #当满足itnihao_web的策略时使用server_web的backend 
       use_backend  server_blog if itnihao_blog 
       #当满足itnihao_blog的策略时使用server_blog的backend 
       redirect prefix http://blog.itniaho.cn code 301 if itnihao 
       #当访问itnihao.cn的时候，用http的301挑转到http://192.168.16.3 
       default_backend server_bbs 
       #以上都不满足的时候使用默认server_bbs的backend  
 
       ##########backend的设置############## 
       #下面我将设置三组服务器 server_web，server_blog，server_bbs
###########################backend server_web############################# 
backend server_web 
       mode http            #http的7层模式 
       balance roundrobin   #负载均衡的方式，roundrobin平均方式 
       cookie SERVERID      #允许插入serverid到cookie中，serverid后面可以定义 
       option httpchk GET /index.html #心跳检测的文件 
       server web1 192.168.16.2:80 cookie web1 check inter 1500 rise 3 fall 3 weight 1  
       #服务器定义，cookie 1表示serverid为web1，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用， 
       #fall 3是3次失败认为服务器不可用，weight代表权重 
       server web2 192.168.16.3:80 cookie web2 check inter 1500 rise 3 fall 3 weight 2 
       #服务器定义，cookie 1表示serverid为web2，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用， 
       #fall 3是3次失败认为服务器不可用，weight代表权重 
 
###################################backend server_blog############################################### 
backend server_blog 
       mode http            #http的7层模式 
       balance roundrobin   #负载均衡的方式，roundrobin平均方式 
       cookie SERVERID      #允许插入serverid到cookie中，serverid后面可以定义 
       option httpchk GET /index.html #心跳检测的文件 
       server blog1 192.168.16.2:80 cookie blog1 check inter 1500 rise 3 fall 3 weight 1  
       #服务器定义，cookie 1表示serverid为web1，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用，fall 3是3次失败认为服务器不可用，weight代表权重 
       server blog2 192.168.16.3:80 cookie blog2 check inter 1500 rise 3 fall 3 weight 2 
        #服务器定义，cookie 1表示serverid为web2，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用，fall 3是3次失败认为服务器不可用，weight代表权重 
 
###################################backend server_bbs############################################### 
 
backend server_bbs 
       mode http            #http的7层模式 
       balance roundrobin   #负载均衡的方式，roundrobin平均方式 
       cookie SERVERID      #允许插入serverid到cookie中，serverid后面可以定义 
       option httpchk GET /index.html #心跳检测的文件 
       server bbs1 192.168.16.2:80 cookie bbs1 check inter 1500 rise 3 fall 3 weight 1  
       #服务器定义，cookie 1表示serverid为web1，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用，fall 3是3次失败认为服务器不可用，weight代表权重 
       server bbs2 192.168.16.3:80 cookie bbs2 check inter 1500 rise 3 fall 3 weight 2 
        #服务器定义，cookie 1表示serverid为web2，check inter 1500是检测心跳频率rise 3是3次正确认为服务器可用，fall 3是3次失败认为服务器不可用，weight代表权重
```

# haproxy 日志配置

开启 haproxy 日志，修改 haproxy 配置文件中，log 关键字

配置 rsyslgo 中 haproxy 的日志保存路径，其中需要让 rsyslog 开启 udp 的 514 端口

```bash
cat > /etc/rsyslog.d/haproxy.conf << \EOF
$ModLoad imudp
$UDPServerRun 514
local2.\* /var/log/haproxy/haproxy.log
& stop
EOF
```

配置 haproxy 的日志轮替

```bash
cat > /etc/logrotate.d/haproxy << \EOF
/var/log/haproxy/haproxy.log {
    daily
    copytruncate
    rotate 10
    missingok
    dateext
    notifempty
    compress
    sharedscripts
    postrotate
        /bin/kill -HUP `cat /var/run/syslogd.pid 2> /dev/null` 2> /dev/null || true
        /bin/kill -HUP `cat /var/run/rsyslogd.pid 2> /dev/null` 2> /dev/null || true
    endscript
}
EOF
```
