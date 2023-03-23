---
title: CDN
---

# 概述

http://3ms.huawei.com/km/groups/1002549/home?l=zh-cn#category=5402776 学习材料

**Content Delivery Network(内容分发网络，简称 CDN)**

内容的定义：内容就是资源，人们浏览的网页，下载的数据，观看的视频等等都属于内容范畴

CDN 产生的原因以及 CDN 的基本概念

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yiw4sk/1616130888007-9f69ee92-f49a-47b2-981d-de386691f4b4.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yiw4sk/1616130887987-dbdd175d-b52c-42e6-93b4-dfb190aa09d2.jpeg)

非签约模式：即不通过与内容提供方签约的方式来获取资源镜像

1. DNS 引流
2. 流量镜像

签约模式：即与内容提供方签约后，获取对方的资源景象

1. 通过 CNAME，域名的别名方式来重定向用户请求

|                |                        |                |                        |                |     |     |
| -------------- | ---------------------- | -------------- | ---------------------- | -------------- | --- | --- |
| 模式           | 非签约模式             |                |                        | 签约模式       |     |     |
| 文件类型       | 大文件(视频类，下载类) | 小文件(网页类) | 大文件(视频类，下载类) | 小文件(网页类) |     |     |
| 调度模式       | 流量镜像               | 本地 DNS       |
| DNS 引流，转发 | 全局 DNS+本地 HTTP     | 全局 DNS       |                        |                |

CDN 流程

Cache 结构

用户请求资源的缓存状态信息表：

1. TCP_HIT：内网用户请求的资源是 HCS 已缓存资源，内网用户获取的资源来自于 HCS 中的已缓存资源。
2. TCP_MISS：内网用户请求的资源不是 HCS 已缓存资源，内网用户获取的资源来自于外网 Web 服务器。
3. TCP_CNC_MISS：内网用户请求的头部规定不缓存这个资源，HCS 不缓存这个资源，内网用户获取的资源来自于外网 Web 服务器。
4. TCP_SNC_MISS：外网 Web 服务器返回的头部规定不缓存这个资源，HCS 不缓存这个资源，内网用户获取的资源来自于外网 Web 服务器。
5. TCP_REFRESH(刷新)\_HIT：内网用户请求的资源命中了 HCS 已缓存资源，但 HCS 需要检查这个资源是否已更新，外网 Web 服务器通知 HCS 这个资源未修改，HCS 将这个资源发送给内网用户。
6. TCP_REFRESH_MISS_METADATA：外网 Web 服务器返回一个对应请求资源的 304 报文，表示这个资源已临时被移走。
7. TCP_REFRESH_MISS：内网用户请求的资源命中了 HCS 已缓存资源，但 HCS 需要检查这个资源是否已更新，外网 Web 服务器通知 HCS 这个资源已经过期，HCS 重新从 Web 服务器获取这个资源后再发送给内网用户。
8. TCP_PARTIAL(部分)\_HIT：客户端分段请求文件的时候，命中请求资源。
9. TCP_PARTIAL_MISS：客户端分段请求文件的时候，未命中请求资源。
10. TCP_REFRESH_UKN_MISS：内网用户请求的资源命中了 HCS 已缓存资源，但 HCS 需要检查这个资源是否已更新，但未能判断出是否更新，代理访问。
11. TCP_REFRESH_NC_MISS：内网用户请求的资源命中了 HCS 已缓存资源，但 HCS 需要检查这个资源是否已更新，，外网 Web 服务器通知 HCS 这个资源未修改，但未从本地吐出，代理访问。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yiw4sk/1616130887966-257146b0-91a9-4008-a239-adfa46ce58b9.jpeg)
