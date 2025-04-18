---
title: Playbook
weight: 1
---

# 概述

> 参考：
> 
> - [官方文档，Playbook指南 - Playbook 介绍](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_intro.html)
> - [官方文档，Playbook 指南 - 使用 Playbook](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks.html)
> - [Ansible Galaxy](https://galaxy.ansible.com/) 类似于 playbook 仓库的地方
> - [公众号，任务中心之Ansible进阶篇](https://mp.weixin.qq.com/s/HA0vKnuKwKOaB5kdcYX9rg)

与通过命令行来执行 Ansbile 任务模式相比，Playbook 是 Ansible 执行任务的另一种方式，而且功能非常强大。

playbook 可以通过定义一个或多个文件，然后让 ansible 使用这些文件来完成一系列复杂的任务。如果说通过命令行是对多台设备执行一个任务，那么 Playbook 则是可以对多台设备按顺序执行不同任务。

playbook 称为"剧本"。每个 playbook 都包含一个或多个 plays(戏剧)。拿拍电影举例，一部电影会有一部“剧本 playbook”来描述电影情节，而整部电影都是由一场一场的“戏剧 play”拼凑起来的。每一场戏剧又需要执行多种“任务 task”(比如亲嘴、打架、聊天、上床~~~)

首先，下面是一个 playbook 的样例。这个 playbook 中包含两个 play，一个叫 webservers，另一个叫 databases。其中 webservers 中包含两个 tasks，一个要使用 yum 模块执行动作，另一个要使用 template 模块，向文件中写入内容

```yaml
- hosts: webservers
  remote_user: root
  tasks:
    - name: ensure apache is at the latest version
      yum:
        name: httpd
        state: latest
    - name: write the apache config file
      template:
        src: /srv/httpd.j2
        dest: /etc/httpd.conf

- hosts: databases
  remote_user: root
  tasks:
    - name: ensure postgresql is at the latest version
      yum:
        name: postgresql
        state: latest
    - name: ensure that postgresql is started
      service:
        name: postgresql
        state: { { item } }
      with_items:
        - started
```

在 [Ansible Galaxy](https://galaxy.ansible.com/) 网站上，我们可以找到大量的社区已经编写好的 Playbook。

# Playbook 关键字

> 参考：
>
> - [官方文档，参考与附录 - Playbook 关键字](https://docs.ansible.com/ansible/latest/reference_appendices/playbooks_keywords.html)

关键字是配置 Ansible 行为的几个来源之一。有关每个源的相对优先级的详细信息，请参阅控制 Ansible 的行为方式：优先级规则。

> 这些关键字有各自的适用场景，通常分为 Play、Role、Block、Task 这几种

**hosts** #

**tasks** # 要在 Play 中执行的主要任务列表，这些任务在 `roles 关键字定义的任务之后`，以及 `post_tasks 关键字定义的任务之前` 执行

**roles** #

**name** #

**check_mode** # 控制任务是否在检查模式下运行。若命令行中使用了 -C 则所有任务得 check_mode 为 true；若 check_mode 为 false，即使使用了 -C 该任务也会真实执行。

- 适用场景：任务
- 常用在通过 Shell 类型任务注册变量的场景。

## Task 任务中的关键字

register

# Playbook 语法详解

```yaml
- hosts: STRING # 指定该playbook要操作的主机
  tasks:
    - name: STRING #(可省略)指定该任务名称
      MODULES: # 指定该任务所要使用的模块名称
        PARAMETER: # 指定该模块参数
```

或

```yaml
- hosts:
  roles:
    - RoleNameOne
    - RoleNameTwo
```

## block - 将多个 task 合并为一个进行统一处理

块允许对任务进行逻辑分组以及进行中的错误处理。您可以应用于单个任务的大多数内容（循环除外）都可以应用于块级，这也使设置任务通用的数据或指令变得更加容易。这并不意味着该指令会影响块本身，而是被块所包含的任务继承。即何时将应用于任务，而不是块本身。

```yaml
tasks:
  - name: Install, configure, and start Apache
    block:
      - name: install httpd and memcached
        yum:
          name:
            - httpd
            - memcached
          state: present

      - name: apply the foo config template
        template:
          src: templates/src.j2
          dest: /etc/foo.conf
      - name: start service bar and enable it
        service:
          name: bar
          state: started
          enabled: True
    when: ansible_facts['distribution'] == 'CentOS'
    become: true
    become_user: root
    ignore_errors: yes
```

## handler - 任务处理器，用于在执行任务时附加额外的任务

> 参考：
>
> - [官方文档，使用 Ansible playbooks - Handlers: 任务 change 时运行操作](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_handlers.html)

ansible 执行的每一个 task 都会报告该任务是否改变了目标，即 changed=true 或 changed=false。当 ansible 捕捉到 changed 为 true 的时候，则会触发一个 notify(通知)组件，该组件的作用就是用来调用指定的 handler。

> 注意：通过 notify 调用的 handler 任务，只有**在所有任务全部完成后，才会执行**。

handlers 示例：在 task 下定义 notify 来指定要调用的 handlers，需要与后面定义的 handlers 的 name 相同

```yaml
tasks:
  - name: copy template file to remote host
    template: src=/etc/ansible/nginx.conf.j2 dest=/etc/nginx/nginx.conf
    notify:
      - restart nginx
      - test web page
  - copy: src=nginx/index.html.j2 dest=/usr/share/nginx/html/index.html
    notify:
      - restart nginx # 与handlers中的name相同才会调用成功
```

handlers 示例：定义 handlers，与定义 task 类似。

```yaml
handlers:
  - name: restart nginx
    service: name=nginx state=restarted
  - name: test web page
    shell: curl -I http://192.168.100.10/index.html | grep 200 || /bin/false
```

上面的示例表示：当执行 template 模块的任务时，如果捕捉到 changed=true，那么就会触发 notify，如果分发的 index.html 改变了，那么也重启 nginx(当然这是没必要的，仅做示例演示)

Note：notify 是在执行完一个 play 中所有 task 后被触发的，在一个 play 中也只会被触发一次。意味着如果一个 play 中有多个 task 出现了 changed=true，它也只会触发一次。例如上面的示例中，向 nginx 复制配置文件和复制 index.html 时如果都发生了改变，都会触发重启 apache 操作。但是只会在执行完 play 后重启一次，以避免多余的重启。

## playbook 中的错误处理

通常情况下, 当出现失败时 Ansible 会停止在宿主机上执行.有时候,你会想要继续执行下去.为此 你需要像这样编写任务:

```yaml
- name: this will not be counted as a failure
  command: /bin/false
  ignore_errors: yes # 使用ignore_errors来忽略该任务失败后终止ansbile的效果
```

# Ansible Artifacts

在 Ansible Playbook 中，可以将 task、playbook、role、var、各种文件 等等统一抽象为 **Artifacts(工件)**

## 复用 Ansible 工件

> 参考：
>
> - [官方文档，Playbook 指南 - 复用 Ansible 工件](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_reuse.html)

我们可以在一个非常大的文件中编写简单的 Playbook，大多数用户首先学习单个文件的方法。然而，将自动化工作分解为较小的文件是组织复杂任务集并重复使用它们的绝佳方法。更小，更分散的 Artifacts 让您在多个 Playbook 中多次使用相同的变量、任务和操作以解决不同的用例。您可以在多个父 Playbook 中使用分布式 Artifacts，甚至可以在一个 Playbook 中多次使用。例如，您可能希望在几个不同的 Playbook 中更新客户数据库。如果您将与更新数据库相关的所有任务放在一个任务文件或角色中，您可以在许多 Playbook 中重复使用它们，同时只需在一个地方维护它们。

Ansible 提供四种可分发、可重复使用的 Artifacts：

- Variables files
- Task files
- palybooks
- roles

### include 与 import 的区别

> 参考：
>
> - https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_reuse.html#comparing-includes-and-imports-dynamic-and-static-re-use
> - https://lework.github.io/2018/01/25/Ansible-xiao-shou-ce-xi-lie-er-shi-san-(-dong-tai-he-jing-tai-bao-han-)/

重用分布式 Ansible 工件的每种方法都有优点和局限性。您可以为某些剧本选择动态重用，为其他剧本选择静态重用。尽管您可以在单个剧本中同时使用动态和静态重用，但最好为每个剧本选择一种方法。混合静态和动态重用可能会在您的剧本中引入难以诊断的错误。此表总结了主要差异，因此您可以为您创建的每个剧本选择最佳方法。

|                                                    | Include_*                               | Import_*                                 |
| -------------------------------------------------- | --------------------------------------- | ---------------------------------------- |
| Type of re-use                                     | Dynamic(动态)                             | Static(静态)                               |
| When processed                                     | At runtime, when encountered            | Pre-processed during playbook parsing    |
| Task or play                                       | All includes are tasks                  | `import_playbook` cannot be a task       |
| Task options                                       | Apply only to include task itself       | Apply to all child tasks in import       |
| Calling from loops                                 | Executed once for each loop item        | Cannot be used in a loop                 |
| Using `--list-tags`<br>i.e. 是否可以通过 -t 选项指定执行包含中的任务 | 无法列出 includes 中的标签                      | `--list-tags` 命令可以列出包含所有导入的任务的标签         |
| Using `--list-tasks`                               | Tasks within includes not listed        | All tasks appear with `--list-tasks`     |
| Notifying handlers                                 | Cannot trigger handlers within includes | Can trigger individual imported handlers |
| Using --start-at-task                              | Cannot start at tasks within includes   | Can start at imported tasks              |
| Using inventory variables                          | Can `include_*: {{ inventory_var }}`    | Cannot `import_*: {{ inventory_var }}`   |
| With playbooks                                     | No `include_playbook`                   | Can import full playbooks                |
| With variables files                               | Can include variables files             | Use `vars_files:` to import variables    |


