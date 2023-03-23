---
title: 别看 DNS 污染闹得欢，现在我用 CoreDNS 将它拉清单
---

## 配置 CoreDNS

既然 CoreDNS 如此优秀，我用它来抵御伟大的防火长城岂不美哉？研究了一圈，发现技术上还是可行的，唯一的一个缺点是不支持使用代理，不过你可以通过 proxychians-ng\[5] 或 proxifier\[6] 来强制使用代理。下面开始折腾。

具体的思路其实非常简单，就是将国内的域名查询请求转发到 114 等国内的公共 DNS 服务器，将国外的域名查询请求转发到 8.8.8.8 等国外的公共 DNS 服务器。然而 CoreDNS 的插件链有点反直觉，同一个插件链上的每一个插件只能出现一次，如果只使用 `forward` 插件是满足不了需求的。CoreDNS 原来还有个插件叫 `proxy`，功能和 `forward` 类似，目测好像同时利用 `proxy` 和 `forward` 插件就可以实现咱的需求了。但理想与现实的差距总是很大，不知道从什么时候开始，CoreDNS 官方编译的二进制文件已经没有 `proxy` 插件了，真是气人。

### dnsredir

偶然间发现了一个第三方插件 dnsredir\[7]，目测可以解决我的所有问题。该插件综合了 `proxy` 和 `forward` 插件的所有优点，支持 UDP、TCP、DNS-over-TLS 和 DNS-over-HTTPS，也支持多个后端，还具备健康检查和故障转移的功能，真是太香了！

它的语法是这样的：

    dnsredir FROM... {
        to TO...
    }

1
2
3
Plain Text

- `FROM...` 是一个文件列表，包含了匹配的域名和解析该域名的服务器，说白了就是 dnsmasq 所使用的格式，直接看例子：


    server=/0-100.com/114.114.114.114
    server=/0-100.com/114.114.114.114

1
2
Plain Text

为什么要用这种格式呢？当然是为了方便啦。

为什么这样会方便呢？当然是为了可以直接用上 FelixOnMars 的大陆区域名列表\[8]了。。。FelixOnMars 同时还提供了 `Google` 和 `Apple` 的域名列表，这在某些地区某些 ISP 可以得到国内镜像的 IP，从而加速访问，想想就刺激。

- 当然，除了使用文件列表外，还可以使用 `.`，类似于上面所说的根域。**这个插件最大的亮点是可以在插件链中重复使用 dnsredir 插件**，只要 `FROM...` 不重复就行。

- `to TO...` 用来将 DNS 解析请求发给上游 DNS 服务器。支持几乎所有 DNS 协议，例如：


    dns://1.1.1.1
    8.8.8.8
    tcp://9.9.9.9
    udp://2606:4700:4700::1111
    tls://1.1.1.1@one.one.one.one
    tls://8.8.8.8
    tls://dns.quad9.net
    doh://cloudflare-dns.com/dns-query
    json-doh://1.1.1.1/dns-query
    json-doh://dns.google/resolve
    ietf-doh://dns.quad9.net/dns-query

1
2
3
4
5
6
7
8
9
10
11
12
13
Plain Text

### 增强版 CoreDNS

dnsredir 虽香，但大家别忘了，它是第三方插件，官方默认的二进制文件是不包含该插件的。你可以选择自己编译，但如果经常需要升级怎么办？总不能每次都手动编译吧，也太累了。

好在有位大佬已经通过 `CI/CD` 流程将所需的第三方插件都集成编译进去了，并定期更新，简直就是我等的福音。大佬的项目地址为：

- [https://github.com/missdeer/coredns_custom_build\[9\]](https://github.com/missdeer/coredns_custom_build%5B9%5D)

现在只需要下载对应操作系统的二进制文件，到处拷贝，就可以运行了。

下面统统以 MacOS 为例作讲解。`Openwrt` 的玩法也一样，参考本文的方法论即可，具体本文就不展开了。

直接下载二进制文件：

    $ wget 'https://appveyorcidatav2.blob.core.windows.net/missdeer-15199/coredns-custom-build/1-7-1-514/idbodwxwywg1xgdg/distrib/coredns-linux-amd64.zip?sv=2015-12-11&sr=c&sig=BhMWcOVtDuaETyz2DcjpOr9GdvkpNVOqoIa7iWFpFNQ%3D&st=2020-12-23T15%3A26%3A19Z&se=2020-12-23T15%3A32%3A19Z&sp=r'
    $ $ tar zxf coredns-linux-amd64.zip
    $ mv coredns-linux-amd64/coredns /usr/local/bin/

1
2
3
Plain Text

### 配置

要深入了解 CoreDNS，请查看其文档\[10]，及 plugins 的介绍\[11]。下面是我的配置文件：

    cat > /usr/local/etc/Corefile <<EOF
    # https://coredns.io/plugins/cache/
    (global_cache) {
        cache {
            # [5, 60]
            success 65536 3600 300
            # [1, 10]
            denial 8192 600 60
            prefetch 1 60m 10%
        }
    }
    .:7913  {
      ads {
          default-lists
          blacklist https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt
          whitelist https://files.krnl.eu/whitelist.txt
          log
          auto-update-interval 24h
          list-store ads-cache
      }
      errors
      hosts {
        fallthrough
      }
      health
      prometheus :9153
      import global_cache
      template ANY AAAA {
          rcode NXDOMAIN
      }
      dnsredir accelerated-domains.china.conf google.china.conf apple.china.conf mydns.conf {
          expire 15s
          max_fails 3
          health_check 3s
          policy round_robin
          path_reload 2s
          to 114.114.114.114 223.5.5.5 119.29.29.29
      }
      dnsredir . {
          expire 60s
          max_fails 5
          health_check 5s
          policy random
          spray
          to tls://8.8.8.8@dns.google tls://8.8.4.4@dns.google
          to tls://1.1.1.1@1dot1dot1dot1.cloudflare-dns.com tls://1.0.0.1@1dot1dot1dot1.cloudflare-dns.com
          # Global TLS server name
          # tls_servername cloudflare-dns.com
      }
      log
      loop
      reload 6s
    }
    EOF

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
Plain Text

- **hosts** : `hosts` 是 CoreDNS 的一个 plugin，这一节的意思是加载 `/etc/hosts` 文件里面的解析信息。hosts 在最前面，则如果一个域名在 hosts 文件中存在，则优先使用这个信息返回；

- **fallthrough** : 如果 `hosts` 中找不到，则进入下一个 plugin 继续。缺少这一个指令，后面的 plugins 配置就无意义了；

- **cache** : 溯源得到的结果，缓存指定时间。类似 TTL 的概念；

- **reload** : 多久扫描配置文件一次。如有变更，自动加载；

- **errors** : 打印/存储错误日志；

- **dnsredir** : 这是重点插件。第一段 dnsredir 配置使用了 4 个文件列表，均是 FelixOnMars 的大陆区域名列表\[12]，这里我还加了一个自定义的文件列表 `mydns.conf`。第二段 dnsredir 配置表示默认的解析配置，可以理解为故障转移，如果某个域名没有匹配到任何一个文件列表，就使用第二段 dnsredir 的上游 DNS 服务器进行解析。通过这样的配置方式，就实现了将国内的域名查询请求转发到 114 等国内的公共 DNS 服务器，将国外的域名查询请求转发到 8.8.8.8 等国外的公共 DNS 服务器。

讲一下我自己的理解：

1. 配置文件类似于 nginx 配置文件的格式；

2. 最外面一级的大括号，对应‘服务’的概念。多个服务可以共用一个端口；

3. 往里面一级的大括号，对应 plugins 的概念，每一个大括号都是一个 plugin。这里可以看出，plugins 是 CoreDNS 的一等公民；

4. 服务之间顺序有无关联没有感觉，但 plugins 之间是严重顺序相关的。某些 plugin 必须用 `fallthrough` 关键字流向下一个 plugin；

5. plugin 内部的配置选项是顺序无关的；

6. 从 plugins\[13] 页面的介绍看，CoreDNS 的功能还是很强的，既能轻松从 bind 迁移，还能兼容 old-style dns server 的运维习惯；

7. 从 CoreDNS 的性能指标看，适合做大型服务。

**注意：该方案的前提是能够强制让 CoreDNS 使用代理，或者更精确一点，让 8.8.8.8 和 8.8.4.4 使用代理。这里的方法比较复杂一点，本文就不介绍了。如果你实在不知道怎么办，可以将 8.8.8.8 这一行删除，直接使用 Cloudflare 提供的 DNS 服务，虽然响应有点慢，但好在可以访问。**

如果你无法忍受 Cloudflare 的响应速度，可以考虑使用国内的无污染 DNS：红鱼 DNS\[14]。然后直接一劳永逸：

    cat > /usr/local/etc/Corefile <<EOF
    # https://coredns.io/plugins/cache/
    (global_cache) {
        cache {
            # [5, 60]
            success 65536 3600 300
            # [1, 10]
            denial 8192 600 60
            prefetch 1 60m 10%
        }
    }
    .:7913  {
      ads {
          default-lists
          blacklist https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt
          whitelist https://files.krnl.eu/whitelist.txt
          log
          auto-update-interval 24h
          list-store ads-cache
      }
      errors
      hosts {
        fallthrough
      }
      health
      prometheus :9153
      import global_cache
      template ANY AAAA {
          rcode NXDOMAIN
      }
      dnsredir accelerated-domains.china.conf google.china.conf apple.china.conf mydns.conf {
          expire 15s
          max_fails 3
          health_check 3s
          policy round_robin
          path_reload 2s
          to 114.114.114.114 223.5.5.5 119.29.29.29
      }
      dnsredir . {
          expire 60s
          max_fails 5
          health_check 5s
          policy random
          spray
          to doh://13800000000.rubyfish.cn
      }
      log
      loop
      reload 6s
    }
    EOF

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
Plain Text

这样 CoreDNS 就不用担心走代理的问题了。

### 定时更新国内域名列表

大陆域名列表每天都会更新，所以还需要写个脚本来更新文件列表。不用检查文件是否存在了，直接简单粗暴无脑更新：

    $ cat > /usr/local/bin/update_coredns.sh <<EOF
    #!/bin/bash
    rm accelerated-domains.china.conf
    wget https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list/accelerated-domains.china.conf -O /usr/local/etc/accelerated-domains.china.conf
    rm apple.china.conf
    wget https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list/apple.china.conf -O /usr/local/etc/apple.china.conf
    rm google.china.conf
    wget https://cdn.jsdelivr.net/gh/felixonmars/dnsmasq-china-list/google.china.conf -O /usr/local/etc/google.china.conf
    EOF
    $ sudo chmod +x /usr/local/bin/update_coredns.sh

1
2
3
4
5
6
7
8
9
10
11
Plain Text

先执行一遍该脚本，更新 Corefile 的配置：

    $ /usr/local/bin/update_coredns.sh

1
Plain Text

然后通过 `Crontab` 制作定时任务，每隔两天下午两点更新域名列表：

    $ crontab -l
    0 14 */2 * * /usr/local/bin/update_coredns.sh

1
2
Plain Text

### 开机自启

MacOS 可以使用 launchctl 来管理服务，它可以控制启动计算机时需要开启的服务，也可以设置定时执行特定任务的脚本，就像 Linux crontab 一样, 通过加装 `*.plist` 文件执行相应命令。Launchd 脚本存储在以下位置, 默认需要自己创建个人的 `LaunchAgents` 目录：

- `~/Library/LaunchAgents` : 由用户自己定义的任务项

- `/Library/LaunchAgents` : 由管理员为用户定义的任务项

- `/Library/LaunchDaemons` : 由管理员定义的守护进程任务项

- `/System/Library/LaunchAgents` : 由 MacOS 为用户定义的任务项

- `/System/Library/LaunchDaemons` : 由 MacOS 定义的守护进程任务项

我们选择在 `/Library/LaunchAgents/` 目录下创建 `coredns.plist` 文件，内容如下：

    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
    <dict>
    <key>Label</key>
    <string>coredns</string>
    <key>ProgramArguments</key>
    <array>
    <string>/usr/local/bin/coredns</string>
    <string>-conf</string>
    <string>/usr/local/etc/Corefile</string>
    </array>
    <key>StandardOutPath</key>
    <string>/var/log/coredns.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/coredns.stderr.log</string>
    <key>KeepAlive</key>
    <true/>
    <key>RunAtLoad</key>
    <true/>
    </dict>
    </plist>

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
XML

设置开机自动启动 coredns：

    $ sudo launchctl load -w /Library/LaunchAgents/coredns.plist

1
Plain Text

查看服务：

    $ sudo launchctl list|grep coredns
    61676 0 coredns

1
2
Plain Text

    $ sudo launchctl list coredns
    {
     "StandardOutPath" = "/var/log/coredns.stdout.log";
     "LimitLoadToSessionType" = "System";
     "StandardErrorPath" = "/var/log/coredns.stderr.log";
     "Label" = "coredns";
     "TimeOut" = 30;
     "OnDemand" = false;
     "LastExitStatus" = 0;
     "PID" = 61676;
     "Program" = "/usr/local/bin/coredns";
     "ProgramArguments" = (
      "/usr/local/bin/coredns";
      "-conf";
      "/usr/local/etc/Corefile";
     );
    };

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
Plain Text

查看端口号：

    $ sudo ps -ef|egrep -v grep|grep coredns
        0 81819     1   0  2:54下午 ??         0:04.70 /usr/local/bin/coredns -conf /usr/local/etc/Corefile
    $ sudo lsof -P -p 81819|egrep "TCP|UDP"
    coredns 81819 root    5u    IPv6 0x1509853aadbdf853      0t0     TCP *:5302 (LISTEN)
    coredns 81819 root    6u    IPv6 0x1509853acd2f39ab      0t0     UDP *:5302
    coredns 81819 root    7u    IPv6 0x1509853aadbdc493      0t0     TCP *:53 (LISTEN)
    coredns 81819 root    8u    IPv6 0x1509853acd2f5a4b      0t0     UDP *:53
    coredns 81819 root    9u    IPv6 0x1509853ac63bfed3      0t0     TCP *:5301 (LISTEN)
    coredns 81819 root   10u    IPv6 0x1509853acd2f5d03      0t0     UDP *:5301

1
2
3
4
5
6
7
8
9
10
11
12
Plain Text

大功告成，现在你只需要将系统的 DNS IP 设置为 `127.0.0.1` 就可以了。

### 验证

    $ doggo www.youtube.com @udp://127.0.0.1
    NAME                     TYPE  CLASS TTL  ADDRESS                  NAMESERVER
    www.youtube.com.         CNAME IN    293s youtube-ui.l.google.com. 127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.14.110           127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.11.174           127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.5.206            127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.5.78             127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.14.78            127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 142.250.72.238           127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 216.58.193.206           127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 142.250.68.110           127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 142.250.68.78            127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 172.217.4.142            127.0.0.1:53
    youtube-ui.l.google.com. A     IN    293s 142.250.68.14            127.0.0.1:53

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
Plain Text

搞定。

## 参考

- CoreDNS 使用与架构分析\[15]

- CoreDNS 搭建无污染 DNS\[16]

### 参考资料

\[1]CoreDNS: _<https://github.com/coredns/coredns>_

\[2]DSL: _[https://www.wikiwand.com/zh/%E9%A2%86%E5%9F%9F%E7%89%B9%E5%AE%9A%E8%AF%AD%E8%A8%80](https://www.wikiwand.com/zh/%25E9%25A2%2586%25E5%259F%259F%25E7%2589%25B9%25E5%25AE%259A%25E8%25AF%25AD%25E8%25A8%2580)_

\[3]DNS-over-TLS: _<https://www.wikiwand.com/zh/DNS_over_TLS>_

\[4]plugin.cfg: _<https://github.com/coredns/coredns/blob/master/plugin.cfg>_

\[5]proxychians-ng: _<https://github.com/rofl0r/proxychains-ng>_

\[6]proxifier: _[https://github.com/yangchuansheng/love-gfw#%E7%95%AA%E5%A4%96%E7%AF%87](https://github.com/yangchuansheng/love-gfw#%25E7%2595%25AA%25E5%25A4%2596%25E7%25AF%2587)_

\[7]dnsredir: _<https://github.com/leiless/dnsredir>_

\[8]FelixOnMars 的大陆区域名列表: _<https://github.com/felixonmars/dnsmasq-china-list>_

\[9]https://github.com/missdeer/coredns\_custom\_build: _<https://github.com/missdeer/coredns_custom_build>_

\[10]文档: _<https://coredns.io/manual/toc>_

\[11]plugins 的介绍: _<https://coredns.io/plugins/>_

\[12]FelixOnMars 的大陆区域名列表: _<https://github.com/felixonmars/dnsmasq-china-list>_

\[13]plugins: _<https://coredns.io/plugins/>_

\[14]红鱼 DNS: _<https://www.rubyfish.cn/dns/solutions/>_

\[15]CoreDNS 使用与架构分析: _<https://zhengyinyong.com/coredns-basis.html>_

\[16]CoreDNS 搭建无污染 DNS: _<https://blog.minidump.info/2019/07/coredns-no-dns-poisoning/>_
