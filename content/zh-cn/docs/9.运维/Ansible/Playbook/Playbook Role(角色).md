---
title: Playbook Role(角色)
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，Playbook 指南 - Roles](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_reuse_roles.html)
> - [Ansible 中文权威指南，Playbooks-Playbook 角色和 Incluede 语句](https://ansible-tran.readthedocs.io/en/latest/docs/playbooks_roles.html)

刚开始运用 Playbook 时，可能会把 Playbook 写成包含内容很多的单一文件，这种情况不利于扩展和复用。这时候可以使用一种方式，将这个复杂的 Playbook [模块化](/docs/8.通用技术/Modular.md)，i.e. 拆分一个复杂的 Playbook 文件成多个零散的小文件，将其组合成一个一个具有不同功能的 playbook。

实现 Ansible Playbook 模块化能力的是一种称为 **Roles(角色)** 的概念。Roles 实际上是对 playbook 进行逻辑上的划分，主要依赖于目录的命名和摆放，一个 Role 就是一个目录，Role 名与目录名相同。

当我们开始思考这些概念：tasks, handlers, variables 等等，是否可以将它们抽象为一个更大的概念呢。我们考虑的不再是”将这些 tasks，handlers，variables 等等应用到这些 hosts 中”，而是有了更抽象的概念，比如：”这些 hosts 是 dbservers” 或者 “那些 hosts 是 webservers”（注：dbserver，webservers 即是”角色”）。这种思考方式在编程中被称为”封装”，将其中具体的功能封装了起来。举个例子，你会开车但并不需要知道引擎的工作原理（注：同样的道理，我们只需要知道”这些 hosts 是 dbservers”，而不需要知道其中有哪些 task，handlers 等）。

# roles 目录结构

下面是一个最基本的 roles 目录结构。在这个目录结构里，有两个 roles，一个名为 common，另一个名为 webservers

```bash
site.yml
webservers.yml
fooservers.yml
roles/
   common/
     tasks/
     handlers/
     files/
     templates/
     vars/
     defaults/
     meta/
   webservers/
     tasks/
     defaults/
     meta/
```

每个目录的作用如下

- **tasks** # 包含角色要执行任务的主要列表
- **handlers** # 包含 handlers，该 role 甚至该 role 之外的任何地方都可以使用这些 handlers
- **defaults** # 包含该 role 的变量的默认值 (see Using Variables for more information).
- **vars** # 包含该 role 的变量的自定义值 (see Using Variables for more information).
- **files** # 包含可以通过该 role 部署的文件。比如通过 script 模块使用的脚本。
- **templates** # 包含可以通过该 role 部署的模板
- **meta** # 为该 role 定义的一些元数据

Note：如果想让这些目录生效，需要在 tasks、handlers、vars、defaults、meta 目录中保存名为 main.yml 的文件，main.yml 的作用详见下文。如果没有文件或目录不存在，则会忽略

## 使用 roles

在 playbook.yaml 文件中，使用关键字 roles 即可调用指定的 roles 内的工作内容

```yaml
- hosts: webservers
  roles:
    - common
    - webservers
```

roles 下指定的每个值(i.e. roles 名称)，ansbile 都会去与该值同名的目录中获取其中所有文件，这其中遵循如下规则：

- 如果 roles/X/tasks/main.yml 存在, 则其中列出的 task 将添加到 play 中。
- 如果 roles/X/handlers/main.yml 存在, 则其中列出的 handler 都将添加到 play 中
- 如果 roles/X/vars/main.yml 存在, 则其中列出的 variables 都将添加到 play 中
- 如果 roles/X/defaults/main.yml 存在, 则其中列出的默认变量值会被添加到 play 中，如果在其他地方没有指定其中列出的变量的值，则会用到默认值
- 如果 roles/X/meta/main.yml 存在,则其中列出的所有角色依赖项都将添加到角色列表中
- task 中 copy，script，template 或 include task 模块都会自动引用 role/X/{files,templates,tasks} 目录中文件，而不必使用绝对路径设置。

Note：

- 其中 X 为 Role 名字
- ansible 会从以下几个目录中寻找与 roles 同名的目录来获取其中的内容
  - ./roles # playbook.yaml 文件所在的目录寻找 roles 目录
  - /etc/ansbile/roles # 默认的系统级别的 roles 目录
  - /root/.ansible/roles
  - /usr/share/ansible/roles
- 也可以在 ansible 的配置文件 ansbile.cfg 中修改 roles_path 字段来指定默认系统级别 role 的位置
