---
title: DDNS
---

# 概述

> 参考：
> 
> - [Wiki，DDNS](https://en.wikipedia.org/wiki/Dynamic_DNS)

**Dynamic Domain Name System(动态域名系统，简称 DDNS)** 是一种方法、概念，这个方法用来动态更新 DNS 中名称对应的 IP。通常情况下，域名都是解析到一个固定的 IP，但 DDNS 系统为动态网域提供一个固定的[名称服务器](https://zh.wikipedia.org/wiki/%E5%90%8D%E7%A8%B1%E4%BC%BA%E6%9C%8D%E5%99%A8)（Name server），透过即时更新，使外界用户能够连上动态用户的网址。

比如家庭宽带，获取到的 IP 地址是实时变化的，要想通过域名访问当家庭宽带内部的服务，则必须使用 DDNS。


[GitHub 项目，jeessy2/ddns-go](https://github.com/jeessy2/ddns-go) 是一个使用 Go 写的，带有 Web 管理页面的 DDNS 工具