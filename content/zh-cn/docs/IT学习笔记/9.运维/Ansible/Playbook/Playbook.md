---
title: Playbook
---

# 概述

> 参考：
> - [Ansible Galaxy](https://galaxy.ansible.com/)

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
> - [官方文档，参考与附录-Playbook 关键字](https://docs.ansible.com/ansible/latest/reference_appendices/playbooks_keywords.html)

关键字是配置 Ansible 行为的几个来源之一。有关每个源的相对优先级的详细信息，请参阅控制 Ansible 的行为方式：优先级规则。
hosts #&#x20;
tasks # 要在 Play 中执行的主要任务列表，这些任务在 `roles 关键字定义的任务之后`，以及 `post_tasks 关键字定义的任务之前` 执行
roles #&#x20;
name #

# Playbook 语法详解

```yaml
- hosts: STRING #指定该playbook要操作的主机
  tasks:
    - name: STRING #(可省略)指定该任务名称
      MODULES: #指定该任务所要使用的模块名称
        PARAMETER: #指定该模块参数
```

或

```yaml
- hosts:
  roles:
    - RoleNameOne
    - RoleNameTwo
```

## block 将多个 task 合并为一个进行统一处理

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

## handler 任务处理器，用于在执行任务时附加额外的任务

官方文档：<https://docs.ansible.com/ansible/latest/user_guide/playbooks_intro.html?highlight=handlers#handlers-running-operations-on-change>

ansible 执行的每一个 task 都会报告该任务是否改变了目标，即 changed=true 或 changed=false。当 ansible 捕捉到 changed 为 true 的时候，则会触发一个 notify(通知)组件，该组件的作用就是用来调用指定的 handler。

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
