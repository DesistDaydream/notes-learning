---
title: Conditionals 条件判断
---

# 概述

> 参考：
> - [官方文档,用户指南-传统目录-条件](https://docs.ansible.com/ansible/latest/user_guide/playbooks_conditionals.html)

通常，play 的结果可能取决于 variable，fact（有关远程系统的知识）或先前的任务结果。在某些情况下，变量的值可能取决于其他变量。可以基于主机是否符合其他条件来创建其他组来管理主机。

Ansible 在条件中使用 Jinja 的 [测试](https://docs.ansible.com/ansible/latest/user_guide/playbooks_tests.html) 和 [过滤器 ](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html)来实现条件判断。详见 [Ansible Template 文章中《Ansible 扩展测试函数》章节](✏IT 学习笔记/🛠️9.运维/Ansible/Playbook/Templates%20 模板(Jinja2).md 模板(Jinja2).md)

## 条件判断的简单样例

下面的样例表示：当 ansible_facts\['os_family'] 变量的值为 Debian 的时候，则执行上面的任务，任务内容是执行 shutdown 命令

    tasks:
    - name: "shut down Debian flavored systems"
      command: /sbin/shutdown -t now
      when: ansible_facts['os_family'] == "Debian"

判断主机是否在某个组中
when: inventory_hostname is search("master") # 当 inventory_hostname 变量的值含有 master 字符串时。
when: inventory_hostname == groups\['kube_master']\[0] # 当当前主机的 inventory_hostname 变量值等于 kube_master 主机组中的第一台主机时
when: inventory_hostname in groups\['kube_master'] # 当当前主机的 inventory_hostname 变量值在 kube_master 主机组中时。

when: testvar1 is none # 当变量 testvar1 已定义，但是值为空时。Note：值为空表示 key 后面的值什么都不写，双引号都不能有
when: ((groups\['kube_master'] | length) > 1) # 当 kube_master 主机组的主机数量大于 1 时。
