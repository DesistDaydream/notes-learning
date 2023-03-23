---
title: "Cockpit"
---


# 概述

> 参考：
> - 官网：<https://cockpit-project.org/>

cockpit 是一个基于 web 的 Linxu 服务器管理工具。可以通过 web 端管理服务器上的虚拟机、容器、服务、网络、存储等等。还可以提供一个 web 版的控制台。

Cockpit 配置

**/etc/cockpit/**

- **./ws-certs.d/** # https 证书保存目录。cockpit 第一次启动时，会在该目录生成 https 所需的 证书与私钥