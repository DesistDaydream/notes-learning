---
title: Packet analyzer
linkTitle: Packet analyzer
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Packet analyzer(包分析器)](https://en.wikipedia.org/wiki/Packet_analyzer)

**Packet analyzer(包分析器)** 是一种计算器程序或计算机硬件，可以拦截和记录通过计算机网络的流量，有的地方也称之为 **Packet sniffer(包嗅探器)**。数据包捕获是拦截和记录流量的过程。随着数据流跨网络流流，分析器捕获每个数据包，如果需要，可以解码分组的原始数据，显示分组中的各种字段的值，并根据适当的 [RFC](/docs/Standard/Internet/IETF.md) 或其他规范分析其内容。

## Packet Analyzer 的实现

各种实现的对比: https://en.wikipedia.org/wiki/Comparison_of_packet_analyzers

- [TCPDump](/docs/7.信息安全/Packet%20analyzer/TCPDump/TCPDump.md)
- [WireShark](/docs/7.信息安全/Packet%20analyzer/WireShark/WireShark.md)
- ......等等

# 抓包程序

Reqable

- https://github.com/reqable/reqable-app # 非开源，只是有个仓库
- 官网 https://reqable.com/
- 图标是 小黄鸟，有 移动端  和 PC 端。宣传自己是 Fiddler + Charles + Postman

[Fiddler](/docs/7.信息安全/Packet%20analyzer/Fiddler.md)

[Charles](/docs/7.信息安全/Packet%20analyzer/Charles.md)

ProxyPin

- [GitHub 项目，wanghongenpin/proxypin](https://github.com/wanghongenpin/proxypin)
- 开源的 HTTP(S) 流量捕获软件，支持 Windows、Mac、Android、iOS、Linux 全平台系统

mitmproxy

- [GitHub 项目，mitmproxy/mitmproxy](github.com/mitmproxy/mitmproxy)
- Python 编写，为渗透测试人员和软件开发人员提供的交互式、支持 TLS 的拦截 HTTP 代理。

HTTP Debugger

- https://www.httpdebugger.com/
- 可以抓进程的包，而不是通过代理的方式抓包

openQPA

- https://github.com/l7dpi/openQPA, https://gitee.com/l7dpi/openQPA
- http://www.l7dpi.com/
- 基于进程抓包

[SunnyNetTools](https://github.com/qtgolang/SunnyNetTools)

- Sunny 网络中间件-抓包工具

# Android 抓包

抓包方式

- PC 上运行抓包程序，Wiki 设置代理到 PC
    - 有些人选择 fiddler 的 CA 导入到安卓手机内，然后电脑上运行 fiddler，安卓 WIFI 设置代理为电脑的fddler 端口，然后发现有些 apk 会检测 WIFI是否设置代理，存在就不走代理导致无法抓到手机 apk 的https 报文。所以还是要安卓上安装软件+VPNService 指定 app 来搞中间人拦截。
- Android 上安装抓包程序

常见问题：

- 证书安装成功，但是抓到的包都是 unknow，可能的原因：
    - Android7.0 之后默认不信任用户级别 CA 证书
    - 此时开启抓包后，很多 APP 都是无网络的情况；但是 chrome 打开网页是可以抓到 https 的包
    - 需要想办法安装在系统级别下的 CA 证书。
        - 有些 APP 内嵌证书，需要修改程序内部逻辑
    - 可能的方法
      - 平行空间
      - 获取系统 Root 权限

解决方式：

- 用苹果的 IOS 系统~~

# 最佳实践

[馆长 Blog，安卓 HTTPS 抓包那些事](https://zhangguanzhang.github.io/2026/02/23/android-https-capture/)
