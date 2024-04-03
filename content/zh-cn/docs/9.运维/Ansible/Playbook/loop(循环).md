---
title: loop(循环)
weight: 6
---

# 概述

> 参考：
>
> - [官方文档，用户指南-传统目录-Loops](https://docs.ansible.com/ansible/latest/user_guide/playbooks_loops.html)

有时需要重复执行多次任务。在计算机编程中，这称为循环。常见的 Ansible 循环包括使用文件模块更改多个文件和/或目录的所有权，使用用户模块创建多个用户以及重复轮询步骤直到达到特定结果。Ansible 提供了两个用于创建循环的关键字：`loop` 和 `with_XX`

- `with_XX` 关键字依赖于 [Lookup Plugins(Lookup 插件)](https://docs.ansible.com/ansible/latest/plugins/inventory.html)。其中 根据插件的不同功能，使用不同的字符串。e.g.with_items 也是 Lookup 插件。
  - 插件的介绍详见：[Ansible Plugins](docs/9.运维/Ansible/Ansible%20Plugins/Ansible%20Plugins.md)
- `loop` 关键字与 with_list 等效，是简单循环的最佳选择

## 循环的简单样例

下面展示了循环的基本功能：通过 loop 或者 with\_\*来对一个列表中的值逐一操作

```yaml
- name: 添加几个用户。循环的基本使用方式
  user:
    name: "{{ item }}"
    state: present
  loop:
    - testuser1
    - testuser2
```

等同于

```yaml
- name: 添加几个用户。先赋值给一个变量，然后在loop关键字中引用变量。
  vars:
    users: [testuser1, testuser2]
  user:
    name: "{{ item }}"
    state: present
  loop: "{{ users}}"
```

上述示例与下面的任务相同。这就相当于将两个任务模块相同但是操作内容不同的任务合并成为一个任务

```yaml
- name: 添加testuser1用户
  user:
    name: "testuser1"
    state: present

- name: 添加testuser2用户
  user:
    name: "testuser2"
    state: present
```

从示例可以看出，循环是使用两个部分来组成整个循环的功能

1. {{ item }}变量来引用 loop 关键字定义的内容
   1. Note：在定义任何变量时，一定要避免使用 item 为变量名
2. loop 关键字来定义列表内容

# 循环样例

## 逐行读取文件

参考链接：<https://stackoverflow.com/questions/48403508/ansible-read-local-file-to-var-and-then-loop-read-line-by-line>

通过 lookup 插件中的 file 功能读取 ip_vs.conf 文件，并逐行读取文件赋值给变量 info，然后执行 modprobe 命令调用变量，逐一加载模块

Note：一定要加 .splitlines() ，否则会将文件里的换行符都赋值给变量，这样，变量只有一个值，而不是一个列表了

```yaml
- name: 加载模块
  vars:
    info: "{{ lookup('file', 'ip_vs.conf').splitlines() }}"
  loop: "{{ info }}"
  shell: modprobe {{ item }}
```

## 匹配目录下所有文件

这个任务通过 with_fieglob 插件来获取该剧本的 file 目录下 rpm 目录下的所有文件，将这些文件复制到远程设备的 /root/downloads/kubeadm-let-ctl/ 目录下。

```yaml
- name: 拷贝rpm安装包
    copy:
      src: "{{ item }}"
      dest: /root/downloads/kubeadm-let-ctl/
    with_fileglob:
    - rpm/*
```
