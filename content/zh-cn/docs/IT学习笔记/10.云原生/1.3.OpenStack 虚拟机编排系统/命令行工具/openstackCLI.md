---
title: openstackCLI
---

#

- 常用命令

- 语法格式

- 子命令列表

- CLI 命令行控制命令组

- common 通用命令组

- compute 计算服务命令组

- identity 身份服务命令组

- image 镜像服务命令组

- network 网络服务命令组

- neutron 客户端命令组

- objectStore 对象存储服务命令组

- volume 卷服务命令组

# 常用命令

查看虚拟机的：server

查看网络的：

查看镜像的：

查看存储的：

# openstack 命令基础

Openstack Command Line Client 官方介绍：<https://docs.openstack.org/python-openstackclient/rocky/>

OpenStackClient(又名 OSC)是**openstack 的命令行客户端**，这个客户端将 compute、identity、image、object、storage 和 blockStorage 这些的 APIs 一起放在一个具有统一命令结构的 shell 中。e.g.nova、neutron、glance 等命令，都会集中在 openstack 的子命令中。

# 语法格式：openstack \[OPTIONS] Command \[CommandArguments]

启动一个 shell 来执行 openstack 客户端中的 Command，或者直接使用 openstack+Command 来对 openstack 进行管理与控制

## OPTIONS

## Command

所有 openstack 的可用的 Command 可以通过`openstack command list` 命令所列出的列表来查看。这些命令通过组来划分，每个命令组代表对一类服务的控制命令

### openstack.cli

**command list \[--group <GroupKeyword>]** #按组列出 openstack 可以支持的所有 Command，可以在选项中指定要查看的具体组名，只查看该组的命令。GroupKeyword 可以使组名中的关键字，不用使用完整的组名

**module list \[--all]** # 显示 OSC 程序已经安装的 python 模块

### openstack.common

availability

configuration

extension

limits

project

quota

versions

### openstack.compute.v2

aggregate

compute

console

flavor

host

hypervisor

ip

keypair

**server** #控制 openstack 上所开虚拟机的命令

openstack server SubCommand \[OPTIONS] \[ARGS]

- SubCommand

- list # 列出 openstack 上所开的虚拟机

- OPTIONS

- --long # 列出更多的关于虚拟机的信息

- EXAMPLE

- \*usage\*\*

### openstack.identity.v3

access

application

catalog

consumer

credential

**domain** #domain 域，是用户、组、项目的集合，每个组合项目仅有一个域

ec2

endpoint

federation

group

identity

implied

limit

mapping

policy

project

region

registered

request

role

service

token

trust

user

### openstack.image.v2

image

### openstack.network.v2

address

floating

ip

network

port

router

security

subnet

### openstack.neutronclient.v2

bgp

bgpvpn

firewall

network

sfc

vpn

### openstack.object_store.v1

container

object

### openstack.volume.v2

backup

consistency

snapshot

volume
