---
title: Playbook Role(角色)
---

# 概述

> 参考：
> - [官方文档，用户指南-Roles](https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html)
> - [Ansible 中文权威指南，Playbooks-Playbook 角色和 Incluede 语句](https://ansible-tran.readthedocs.io/en/latest/docs/playbooks_roles.html)

刚开始学习运用 playbook 时，可能会把 playbook 写成一个内容很多的文件，这种情况不利于扩展和复用。这时候可以使用一种方式，将这个复杂的 playbook 模块化，即拆分一个复杂的 playbook 文件成多个零散的小文件，将其组合成一个一个具有不同功能的 playbook。

这时候就需要用到 ansible playbook 的 roles 概念了。roles 实际上是对 playbook 进行逻辑上的划分，主要依赖于目录的命名和摆放，一个 Role 就是一个目录，Role 名与目录名相同。

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

    - hosts: webservers
      roles:
        - common
        - webservers

roles 下指定的每个值(roles 名)，ansbile 都会去与该值同名的目录中获取其中所有文件，这其中遵循如下规则：

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

# roles 目录结构的最佳示例

> 参考：
> - [官方文档，用户指南-配置样例](https://docs.ansible.com/ansible/latest/user_guide/sample_setup.html)

```bash
production                # 适用于 production 的 Inventory 文件
staging                   # 适用于 staging 的 Inventory 文件

group_vars/               # 在这里定义组的变量
   group1.yml             # 文件名以组名命名，group1.yml 是适用于 group1 组的变量
   group2.yml
host_vars/								# 在这里定义主机变量
   hostname1.yml          # 文件名以主机名命名，hostname1.yml 是适用于 hostname1 主机的变量
   hostname2.yml

library/                  # if any custom modules, put them here (optional)
module_utils/             # if any custom module_utils to support modules, put them here (optional)
filter_plugins/           # if any custom filter plugins, put them here (optional)

site.yml                  # master playbook
webservers.yml            # playbook for webserver tier
dbservers.yml             # playbook for dbserver tier

# 当需要管理多个 Role 时，可以在 roles/ 目录中
roles/
    common/               # 名为 common 的角色
        tasks/            #
            main.yml      #  <-- tasks file can include smaller files if warranted
        handlers/         #
            main.yml      #  <-- handlers file
        templates/        #  <-- files for use with the template resource
            ntp.conf.j2   #  <------- templates end in .j2
        files/            #
            bar.txt       #  <-- files for use with the copy resource
            foo.sh        #  <-- script files for use with the script resource
        vars/             #
            main.yml      #  <-- variables associated with this role
        defaults/         #
            main.yml      #  <-- default lower priority variables for this role
        meta/             #
            main.yml      #  <-- role dependencies
        library/          # roles can also include custom modules
        module_utils/     # roles can also include custom module_utils
        lookup_plugins/   # or other types of plugins, like lookup in this case

    webtier/              # 名为 webtier 的角色，其内的机构与 common 相同
        ......
    monitoring/           # 同上
        ......
    fooapp/               # 同上
        ......
```

## group_vars 与 host_vars 目录

组变量与主机变量的文件除了可以放在 Palybook 的根目录，还可以放在存放 Inventory 文件的目录中，比如：

```bash
inventories/							# 这里存放 Inventory 目录，通过在命令行中使用 -i 选项以指定 Inventory 文件
   production/
      hosts               # inventory file for production servers
      group_vars/
         group1.yml       # here we assign variables to particular groups
         group2.yml
      host_vars/
         hostname1.yml    # here we assign variables to particular systems
         hostname2.yml

   staging/
      hosts               # inventory file for staging environment
      group_vars/
         group1.yml       # here we assign variables to particular groups
         group2.yml
      host_vars/
         stagehost1.yml   # here we assign variables to particular systems
         stagehost2.yml

library/
module_utils/
filter_plugins/

site.yml
webservers.yml
dbservers.yml

roles/
    common/
    webtier/
    monitoring/
    fooapp/
```

至于 group_vars 与 host_vars 在不同目录的优先级可以参考 [Ansible 变量-优先级](✏IT 学习笔记/🛠️9.运维/Ansible/Ansible%20Variables(变量).md Variables(变量).md) 部分
