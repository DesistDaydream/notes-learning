---
title: network
linkTitle: network
weight: 20
---

# 概述

> 参考：
>
> -

network 命令用于管理容器的网络

## connect - Connect a container to a network

## create - Create a network

**docker network create \[OPTIONS] NETWORK**

EXAMPLE

- docker network create -d bridge --subnet "172.26.0.0/16" --gateway "172.26.0.1" mybr0 # 创建一个桥接的网络，网段是 172.26.0.0/16,网关是 172.26.0.1

## disconnect - Disconnect a container from a network

## inspect - Display detailed information on one or more networks

## ls - List networks

## prune - 移除所有未使用的网络

## rm - 移除一个或多个网络
