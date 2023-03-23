---
title: systemd-resolved.service
---

# 概述

> 参考：[金步国-system 中文手册](http://www.jinbuguo.com/systemd/systemd-resolved.service.html)

systemd-resolved.service 是一个类似于 DNSmasq 的域名解析服务，只不过这个服务只适用于 Linux 中，且被 systemd 所管理。

# 配置

**/run/systemd/resolve/resolv.conf **# 具体的解析配置
**/usr/lib/systemd/resolv.conf** # 顶替 /etc/resolv.conf 文件
