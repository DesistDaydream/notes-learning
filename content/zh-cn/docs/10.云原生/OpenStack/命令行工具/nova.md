---
title: nova
linkTitle: nova
weight: 20
---

# 概述

> 参考：
>
> -

# Syntax(语法)

**nova  \[OPTIONS]  \[SubCommand  \[OPTIONS]]**

nova list \[OPTIONS] # 列出SERVER相关信息

nova show  # 显示指定SERVER的详细信息，非常详细

nova instacne-action-list  # 列出指定SERVER的操作，创建、启动、停止、删除时间等，

注意：语法中的SERVER指的都是已经创建的虚拟服务器，SERVER可以用实例的NAME(实例名)或者UUID(实例的ID)来表示，SERVER的ID和NAME可以用过nova list命令查到

可以使用nova help SubCommand命令查看相关子命令的使用方法

nova list \[OPTIONS] # 列出SERVER相关信息

OPTIONS

- --all-tenants # 显示所有租户的SERVER信息，可简写为--all-t
- --tenant  \[] # 显示指定租户的SERVER信息

EXAMPLE

- nova list --all-t # 显示所有正在运行的实例，可以查看实例以及ID和主机名
- nova list --all-t --host `cat /etc/uuid` # 显示`cat /etc/uuid`命令输出的主机名的节点运行的实例信息

nova show  # 显示指定SERVER的详细信息，非常详细

EXAMPLE

- nova show ID # 以实例ID展示该实例的详细信息
- nova show ID | grep host # 以实例ID查看所在节点的主机名

nova instacne-action-list  # 列出指定SERVER的操作，创建、启动、停止、删除时间等，

EXAMPLE

- nova instance-action-list ID|NAME # 以实例ID显示该实例的活动信息，包括启动、停止、创建时间等
