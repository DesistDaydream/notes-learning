---
title: "OAuth"
linkTitle: "OAuth"
weight: 20
---

# 概述

> 参考：
>
> - [RFC 6749, The OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
> - [Wiki, Oauth](https://en.wikipedia.org/wiki/Oauth)

**Open Authorization(简称 OAuth)** 是一种用于 access delegation(访问授权) 的开放标准。通常作为互联网用户授权网站或应用程序访问其在其它网站上的信息的一种方式，而无需提供密码。亚马逊、谷歌、Meta Platforms、微软、推特、etc. 公司采用此机制，允许用户与第三方应用程序或网站共享其账户信息。

OAuth 的出现主要是解决这么一个问题：第三方应用程序，如何安全地获得用户授权，以访问该用户在另一个服务上的资源。

一个简单的场景是：我开发了一个照片打印程序（程序 A），用户想要打印自己存储在网盘中的照片，那么用户如何在不告诉 A 账号密码或者任何登录用 Token 等认证信息的情况下，让 A 访问到网盘中的照片呢？

程序 A 应该先找网盘要登录方式，提供给用户。网盘验证用户登录成功后，告诉 A 用户已登录并同意了 A 访问照片，同时提供一个 A 用的 Token。此时 A 可以使用这个 Token 来访问允许的内容。这中间，用户并不用把自己的认证信息提供给 A。

> [!Note]
> 在上述场景中，为了避免将用户的敏感资源泄露给程序 A，这种程序一般都是**本地部署**，并不是通过互联网为多个用户提供服务的。比如 NAS，部署在本地，通过程序编写好的 OAuth 能力获取到权限，即可获取各种云盘中的资源
>
> 另外一个场景是非敏感信息的，比如 用户名、头像 之类的。这种就可以不用本地部署，而是第三方程序在互联网提供服务，比如各种网站在登陆时通过 Google 登录，本质也是使用了 Google 的 OAuth 能力获取 邮箱、头像 等非敏感信息。

![600](Excalidraw/information_security/oauth.excalidraw.md)

