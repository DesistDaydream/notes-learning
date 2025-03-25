---
title: Vmess
linkTitle: Vmess
weight: 20
---

# 概述

> 参考：
>

https://www.chengxiaobai.com/essays/v2ray-trojan-xray.html

Vmess 协议是 V2Ray 项目的主流加密协议

# Xray

> 参考：
>
> - 原文: https://qoant.com/2021/04/vps-with-xray/

V2Fly 与 V2Ray

V2Fly 是 V2Ray 的延伸，因为 V2Ray 仓库的作者失踪，所以建立了 V2Fly，本质上两者没区别。

## Xray 介绍

由于 Debian 包维护人员发现 [XTLS库](https://github.com/XTLS/Go) 的 LICENSE 不是 BSD 许可，提了一个 issue 希望作者 [@rprx](https://github.com/rprx) 能修改 LISENCE 许可方便打包，详见 [https://github.com/XTLS/Go/issues/9](https://github.com/XTLS/Go/issues/9)。由这个 issue 引发了广泛讨论，rprx 认为目前许可不是问题，也有不少人认为协议是立场的体现，各执一词。

最终 V2ray(V2fly社区) 维护团队经过投票确认 XTLS 不符合 V2ray 的 MIT 协议，并在 [V2ray-core 4.33.0版本](https://github.com/v2fly/v2ray-core/releases/tag/v4.33.0)移除了 XTLS。rprx 和其拥护者行动起来，很快就创建了 [Project X](https://github.com/XTLS) 项目和其核心 [Xray](https://github.com/XTLS/Xray-core)（Xray取名来自XTLS和V2ray的结合），并以 XTLS 为核心协议陆续发布了 Xray-core 的多个版本，于是 Xray 诞生了。

XTLS 和 Xray 离不开作者 [@rprx](https://github.com/rprx) 的辛勤付出，因此也简要介绍一下 [@rprx](https://github.com/rprx)：

1. [@rprx](https://github.com/rprx) 是VLESS协议的设计者，在介绍VLESS协议时写下了“性能至上、可扩展性空前，目标是全场景终极协议”的宏壮愿景；
2. [@rprx](https://github.com/rprx) 是XTLS的作者，在[XTLS库](https://github.com/XTLS/Go)中写下了“THE FUTURE”的霸气描述。将内外两条 TLS 连接结合，rprx 可能不是第一个有这想法的人，但却是第一个将其实现、并成熟应用到实际中的作者。从使用表现上看，XTLS 无愧于 rprx 对其的评价：“划时代的革命性概念与技术：XTLS”，以及社区给出的“黑科技”称谓；
3. [@rprx](https://github.com/rprx) 是 Project X 和 Xray 项目的创始人。由于 LICENSE 理念之争，rprx 创建了对标 Project V 和 V2ray-core 的 Project X 和 Xray-core 项目，广受欢迎。

## Xray 和 V2ray 的区别

在说明 Xray 和 V2ray 区别之前，先说一下三个相近但不同的概念：

- **V2ray：** Project V是用于构建基础通信网络的工具合集，其核心工具称为V2Ray。V2ray主要负责网络协议和功能的实现，既可以单独运行，也可以和其它工具配合。V2ray官网是：[https://v2ray.com/](https://v2ray.com/)，Github项目主页是：[https://github.com/v2ray](https://github.com/v2ray)，TG讨论组是：[@projectv2ray](https://t.me/projectv2ray)；
- **V2fly：** 出现一些科学上网作者被喝茶事件后，V2ray原开发者长期不上线，其他维护者没有完整权限，导致V2ray项目维护困难。因此社区在2019年组建了V2fly组织，继续维护V2ray，也是目前V2ray发展的主力。V2fly官网是：[https://www.v2fly.org](https://www.v2fly.org/)，Github项目主页是：[https://github.com/v2fly](https://github.com/v2fly)，TG通知频道：[@v2fly](https://t.me/v2fly)，TG交流群为：[@v2fly_chat](https://t.me/v2fly_chat)；
- **Xray：** 因许可理念之争，VLESS 和 XTLS 的作者单独创建了 Xray 项目，目前是 V2ray 的超集，后续可能有不同的发展路线。Xray 文档官网：[https://xtls.github.io/](https://xtls.github.io/)，Github 项目主页：[https://github.com/XTLS](https://github.com/XTLS)，TG 交流群：[@projectXray](https://t.me/projectXray)。

从上面可以看到，先有 V2ray(Project V)，然后是 V2fly，最后才出来 Xray(Project X)。其中 V2fly 是 V2ray 的社区，可以认为两者是同一个组织。

Xray 和 V2ray 区别如下：

1. **Xray是V2ray的一个分支(Fork)**。Xray 项目基于V2ray而来，其支持并且兼容V2ray的配置；
2. **Xray是V2ray的超集**。虽然最新版 V2ray 删除了XTLS，但仍保留VLESS协议。Xray提供完整的VLESS和XTLS支持，目前是V2ray的超集，但后续Xray可能会有会有自己的发展方向；
3. **如果使用 XTLS，强烈推荐使用 Xray**，不使用 XTLS 的情况下，使用 V2ray 和 Xray 均可。

简而言之，Xray 是 V2Ray 的项目分支，Xray 是 V2Ray 的超集，就跟 Trojan-Go 和 Trojan-GFW 的关系类似，而且 Xray 性能更好、速度更快，更新迭代也更频繁。由于自 V2ray-core 4.33.0 版本起，删除了 XTLS 黑科技，但仍然支持 VLESS，所以是否原生支持 XTLS 是 Xray 和 V2Ray 最大的区别之一。

# 概述

> 参考：
>
> - [GItHub 项目，v2fly/v2ray-core](https://github.com/v2fly/v2ray-core)
> - [路由规则](https://github.com/Loyalsoldier/v2ray-rules-dat)

V2Ray 是 Project V 下的一个工具。Project V 是一个包含一系列构建特定网络环境工具的项目，而 V2Ray 属于最核心的一个。 官方中介绍 `Project V 提供了单一的内核和多种界面操作方式。内核（V2Ray）用于实际的网络交互、路由等针对网络数据的处理，而外围的用户界面程序提供了方便直接的操作流程。` 不过从时间上来说，先有 V2Ray 才有 Project V。 如果还是不理解，那么简单地说，V2Ray 是一个与 Shadowsocks 类似的代理软件，可以用来科学上网（翻墙）学习国外先进科学技术。
