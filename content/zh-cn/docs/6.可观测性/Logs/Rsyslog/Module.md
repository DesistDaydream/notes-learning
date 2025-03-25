---
title: Module
linkTitle: Module
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，模块](https://www.rsyslog.com/doc/configuration/modules/index.html)

在配置模块时，可以会利用各种 Rsyslog [objects](https://www.rsyslog.com/doc/rainerscript/configuration_objects.html#objects)(对象) 对模块进行更多配置，比如 input、[action](https://www.rsyslog.com/doc/configuration/actions.html)、etc. 这些对象。

> objects 的概念是 Rsyslog 设计的 RainerScript 脚本语言中的。

- **module** 对象 # 用于加载模块
- **input** 对象 # 用于描述 Input 行为的主要手段，用于收集 rsyslog 处理的消息。
- **action** 对象 # TODO
- TODO

在阅读各种模块文档时，通常能看到这种标题，这部分内容就是在说利用 input 对象对该模块进行配置。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/observability/rsyslog/202411072158977.png)

比如我 `module(load="imuxsock")` 这样加载了 imuxsock 模块，然后可以这样 `input(type="imuxsock" Socket="/root/tmp/log")` 利用 input 对象控制 imuxsock 模块的行为。Socket 参数是告诉 imuxsock 模块要在 /root/tmp/log 文件处监听 Unix Socket。

还有类似这样的标题: `Action Parameters`，这就是在说可以利用 action 对象对该模块进行配置。

# Output 模块

https://www.rsyslog.com/doc/configuration/modules/idx_output.html

# Input 模块

https://www.rsyslog.com/doc/configuration/modules/idx_input.html

Input 模块用于从各种来源收集消息。它们与消息生成器交互。它们通常是通过输入配置对象定义的。如果要配置具体的输入参数，通常是利用 input 对象实现的。

## imuxsock

https://www.rsyslog.com/doc/configuration/modules/imuxsock.html

配置示例： https://www.rsyslog.com/doc/configuration/modules/imuxsock.html#examples

```bash
# 加载 imuxsock 模块，以便让 Rsyslog 可以监听 /dev/log 这个 Unix Socket 以接收日志消息
module(load="imuxsock")
# 为 imuxsock 模块添加 input 参数，多监听一个 /root/tmp/log 这个 Unix Socket。发送到这里的消息也会被 Rsyslog 处理
input(type="imuxsock" Socket="/root/tmp/log")
```

## imjournal

https://www.rsyslog.com/doc/configuration/modules/imjournal.html

## 远程输入相关模块

### imtcp

配置 TCP 协议的 syslog 接收，用于在日志服务器的时候配置

- **$ModLoad imtcp** # 使用 tcp 进行传输
- **$InputTCPServerRun 514** # 监听在 514 端口上
