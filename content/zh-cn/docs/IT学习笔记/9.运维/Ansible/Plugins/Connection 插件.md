---
title: "Connection 插件"
linkTitle: "Connection 插件"
weight: 20
---

# 概述
> 参考：
> - [官方文档，Connection 插件](https://docs.ansible.com/ansible/latest/plugins/connection.html)

连接插件允许 Ansible 连接到目标主机，以便它可以在它们上执行任务。 Ansible 附带了许多连接插件，但每个主机一次只能使用一个。

默认情况下，Ansible 附带了几个连接插件。 最常用的是 ssh 和 local 类型。 所有这些都可以在剧本中使用，并与 /usr/bin/ansible 一起决定你想如何与远程机器交谈。 如有必要，我们可以创建自定义连接插件。
