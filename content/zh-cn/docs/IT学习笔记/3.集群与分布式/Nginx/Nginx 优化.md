---
title: Nginx 优化
---

1. 性能优化的相关配置
   1. work_processes NUM; #常用，指定 work 线程个数，通常应该少于 cpu 物理核心数，设为 auto 为自动判断
   2. work_cpu_affinity CpuMask; #常用，设定 cpu 掩码，用于绑定给 nginx 专用的 cpu 数
   3. timer_resolution Num; #计时器解析度，降低此值，可提高性能
   4. worker_priority NUM; #设定优先级，即 worker 线程的 nice 值
2. 事件相关配置
   1. worker_connections NUM; #常用，指定每个 worker 线程所能处理的最大并发连接数
   2. accept_mutex on|off; #调度用户请求至 worker 线程时使用的负载均衡锁。on 是让多个 worker 轮流的，序列化地响应新请求
   3. lock_file /PATH/FILE; #指定 accept_mutex 开启后用到的锁文件路径
3. 用于调试、定位问题的配置
   1. daemon on|off； #是否以守护进程方式运行 nginx，调试时设置为 off
   2. master_process on|off； #是否以 master/worker 模式来运行 ngins，调试时可以设置为 off

###

### 优化 Nginx 数据包头缓存

1）优化前，使用脚本测试长头部请求是否能获得响应

    [root@proxy ~]# cat lnmp_soft/buffer.sh
    #!/bin/bash
    URL=http://192.168.4.5/index.html?
    for i in {1..5000}
    do
        URL=${URL}v$i=$i
    done
    curl $URL                                //经过5000次循环后，生成一个长的URL地址栏
    [root@proxy ~]# ./buffer.sh
    .. ..
    <center><h1>414 Request-URI Too Large</h1></center>        //提示头部信息过大

2）修改 Nginx 配置文件，增加数据包头部缓存大小

```nginx
# vim /usr/local/nginx/conf/nginx.conf
... ..
http {
    client_header_buffer_size    1k;          // 默认请求包头信息的缓存
    large_client_header_buffers  4 4k;        // 大请求包头部信息的缓存个数与容量
    .. ..
}
# /usr/local/nginx/sbin/nginx -s reload
```

3）优化后，使用脚本测试长头部请求是否能获得响应

    1.[root@proxy ~]# cat buffer.sh
    2.#!/bin/bash
    3.URL=http://192.168.4.5/index.html?
    4.for i in {1..5000}
    5.do
    6.    URL=${URL}v$i=$i
    7.done
    8.curl $URL
    9.[root@proxy ~]# ./buffer.sh

### 浏览器本地缓存静态数据

1）使用 Firefox 浏览器查看缓存
以 Firefox 浏览器为例，在 Firefox 地址栏内输入 about:cache 将显示 Firefox 浏览器的缓存信息，如图所示，点击 List Cache Entries 可以查看详细信息。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bzuh2e/1621231619142-0d5c79d5-689e-48b4-9b82-2e6fd6450a02.webp)
2）清空 firefox 本地缓存数据，如图所示。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/bzuh2e/1621231619179-e56a8fb4-8a7c-4e06-afe1-6bb66fe30d1b.webp)
3）改 Nginx 配置文件，定义对静态页面的缓存时间

    # vim /usr/local/nginx/conf/nginx.conf
    server {
            listen       80;
            server_name  localhost;
            location / {
                root   html;
                index  index.html index.htm;
            }
    location ~* \.(jpg|jpeg|gif|png|css|js|ico|xml)$ {
    expires        30d;            //定义客户端缓存时间为30天
    }
    }
    # cp /usr/share/backgrounds/day.jpg /usr/local/nginx/html
    # /usr/local/nginx/sbin/nginx -s reload
    #请先确保nginx是启动状态，否则运行该命令会报错,报错信息如下：16.#[error] open() "/usr/local/nginx/logs/nginx.pid" failed (2: No such file or directory)

4）优化后，使用 Firefox 浏览器访问图片，再次查看缓存信息

    # firefox http://192.168.4.5/day.jpg

在 firefox 地址栏内输入 about:cache，查看本地缓存数据，查看是否有图片以及过期时间是否正确。
