---
title: Ansible
weith: 1
---

# 概述

> 参考：
> - [官网](https://www.ansible.com/)
> - [GitHub 项目](https://github.com/ansible/ansible)
> - [官方文档](https://docs.ansible.com/ansible/latest/index.html)
> - [公众号,程序员面试吧-快速入门 Ansible 自动化运维工具](https://mp.weixin.qq.com/s/qu0vPgyNBbRlTtf1pMtq7Q)
> - <https://www.zsythink.net/archives/tag/ansible/>

Ansible is a radically simple IT automation platform that makes your applications and systems easier to deploy and maintain. Automate everything from code deployment to network configuration to cloud management, in a language that approaches plain English, using SSH, with no agents to install on remote systems
Ansible 是一个非常简单的 IT 自动化系统。它处理配置管理、应用程序部署、云供应、临时任务执行、网络自动化和多节点编排。Ansible 可以轻松得批量进行复杂的更改，例如使用负载均衡器进行零停机滚动更新。而这一过程使用 SSH 实现，无需在远程系统上安装代理程序。
Ansible 的主要目标是简单易用。它还非常注重安全性和可靠性，具有最少的移动部件，使用 OpenSSH 进行传输（使用其他传输和拉模式作为替代），以及一种围绕人类可审计性设计的语言 - 即使是那些不熟悉的人该程序。

## 概念

**Control node(控制节点)** # 任何安装了 Ansible 的机器都可以称为控制节点。在控制节点中可以通过调用 `ansible` 或 `ansible-playbook` 命令来运行 Ansible 相关命令和 Playbooks。

**Managed nodes(受管理节点)** # 使用 Ansible 管理的 网络设备 或 服务器。受管理节点有时候也称为 **hosts**。

**Inventory(库存)** # 库存是一个受管理节点的列表。库存有时候也称为 **hostfile**。Inventory 还可以用来组织受管理节点，将每个节点进行分组，以便于扩展

**Collections** # 2.10 版本之后的新概念。Collections 是 Ansible 内容的分发格式，可以包括 Playbooks、Role、Modules、Plugins。新版中，Modules 就被托管于 Collections 中。

- 随着 Ansible 的发展，越来越多的模块、插件被开发并加入到 Ansible 的大家庭，这时候难免会出现命名上的冲突，或者调用上的重复。所以，从 2.10 版本之后，提出了 Collections 的概念。
  - Collections 最大的一个功能就是将模块分类，比如以前 核心模块 command，现在的全名就叫 ansible.builtin.command，前面的 ansible.builtin 就是 command 的 Collections。这种全名称为 **Full Qualified Class Name(完全限定类名，简称 FQCN)**。

**Tasks(任务)** # Ansible 工作的最小单元，Ansible 对受管理节点执行的操作，称为任务。

**Modules(模块)** # 模块就是 Ansible 用来执行 Tasks 的代码。

**Playbooks(剧本)** # 一个被保存起来的有序的 Tasks 列表，通过重复运行 Playbooks，可以方便得重复一组任务。Playbooks 中还可以包含变量、模板、条件语句、控制循环，从本质上来说，编写一个 Playbooks，就好像编写一个脚本代码一样。

Playbooks 是 Ansible 的精髓，如果把 Ansible 当做一门语言，那么就成可以称为 Playbooks 脚本编程语言。

# Ansible 的核心组件

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ot4g6f/1616125280904-828340be-8634-4a3f-a97b-d4600818bd6e.jpeg)

- **ansible core** # 核心组件，ansible 本身
- **host inventory** # 主机库存，Ansible 所管理的主机清单，一个文本文件
- **Modules **# 模块。ansible 的 modules 是实现 ansible 的核心，
  - **core modules** # 核心模块
    - ansible 执行任何命令，都是通过 module 来完成；比如 ansible 让被管理机创建一个用户，那么就会去 core modules 中调用一个能实现创建用户功能的模块，来执行这个操作。
  - **custom modules** # 自定义模块
    - 可以使用任何编程语言来编写模块，只要符合 ansible 的标准即可，可以实现 ansible 本身不具备的功能

# Ansible 关联文件与配置

**/etc/ansible/ansible.cfg** # ansible 使用时调用的配置文件
**/etc/ansible/hosts** # Inventory 的默认配置文件。该文件可以定义被管理主机的 IP，port 等，都可以定义在该文件中，具体格式如下

- 单独 host，任何未分组的主机，需要在定义主机组之前定义各单独的 host，可以是 IP 地址或者主机名
- 主机组，定义一个主机组，组名用\[]括起来，可以定义多个主机组；当使用 ansible 命令的时候，可以使用组名来对该组内所有主机进行操作
- 配置文件说明：详见：[inventory 配置文件详解](/docs/IT学习笔记/9.运维/Ansible/Inventory%20配置文件详解.md)
