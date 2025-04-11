---
title: "Ansible Collection"
linkTitle: "Ansible Collection"
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，使用 Ansible collections](https://docs.ansible.com/ansible/latest/collections_guide/index.html)
> - [官方文档，参考 - Collection 索引](https://docs.ansible.com/ansible/latest/collections/index.html)

**Collection(集合)** 是对 [Ansible Plugins](docs/9.运维/Ansible/Ansible%20Plugins/Ansible%20Plugins.md) 和 [Ansible Modules](docs/9.运维/Ansible/Ansible%20Modules/Ansible%20Modules.md) 的高层抽象。还可以包括 Playbooks、Role、Modules、Plugins。

> 随着 Ansible 的发展，越来越多的模块、插件被开发并加入到 Ansible 的大家庭，这时候难免会出现命名上的冲突，或者调用上的重复。所以，从 2.10 版本之后，提出了 Collections 的概念。
>
> Collections 最大的一个功能就是将模块分类，比如以前 核心模块 command，现在的全名就叫 ansible.builtin.command，前面的 ansible.builtin 就是 command 的 Collections。这种全名称为 **Full Qualified Class Name(完全限定类名，简称 FQCN)**。

在 [Ansible Galaxy](https://galaxy.ansible.com/) 中可以找到非常多用户公开的 Collection。

# 关联文件与配置


# ansible.builtin

> 参考：
> 
> - [官方文档，参考 - Collection 索引 - Ansible.Builtin](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/index.html)

Ansible 的内置 Collection，包括了全部的内置模块与内置插件（比如最常用的连接插件、文件模块等等）

