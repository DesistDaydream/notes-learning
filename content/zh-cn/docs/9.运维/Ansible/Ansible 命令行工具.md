---
title: Ansible 命令行工具
weight: 7
---

# 概述

> 参考：
>
> - [官方文档，用户指南 - 命令行工具](https://docs.ansible.com/ansible/latest/user_guide/command_line_tools.html#command-line-tools)
> - [官方文档，用户指南 - 传统目录 - 临时命令简介](https://docs.ansible.com/ansible/latest/user_guide/intro_adhoc.html)

由于 Ansible 是基于 SSH 远程管理主机，所以让 Ansible 的控制节点需要对受管理节点进行 ssh 的基于密钥的认证(方法详见 ssh 命令)或者在 inventory 文件中指定认证方式。

Note：Ansible 的控制节点和受管理节点的 Python 版本需要一致，否则 ansible 命令有时候会执行失败

# ansible

> 参考：
>
> - [官方文档，用户指南 - 使用命令行工具工作 - ansible](https://docs.ansible.com/ansible/latest/cli/ansible.html)

ansible 是 Ansible 的一个 ad-hoc(临时) 命令，可以在一个或多个受管理节点上自动执行单个任务。ansible 命令即简单又快速，但并不具备可重复性，通过 ansible 命令的使用，我们可以了解到 Ansible 的简单性和强大的功能。**并且，可以将类似的操作，直接移植到 Playbooks 中**。

临时命令非常适合很少重复，但是又需要大批量执行的任务，比如重启服务器、复制文件、管理服务、管理用户 等等。我们可以在临时任务中使用所有 Ansible 模块

## Syntax(语法)

**ansible \<HostPattern> \[OPTIONS]**

- **HostPattern** # 主机模式,可以是主机名，主机 IP，组名，还有一个 all(所有 hosts 里定义的主机)

### OPTIONS

- **--list-hosts** # 列出所有 HostPattern 定义的被管理 host 并统计数量，一般用于查看组内的主机有多少
- **-i, --inventory, --inventory-file INVENTORY** # 指定具体的 INVENTORY 路径或文件，而不使用配置中默认的。
  - INVENTORY 可以是 目录、文件、逗号分割的 IP 列表。
  - 可以使用 `-i 192.168.0.1,` 这种方式手动指定主机而不使用文件。注意：末尾的逗号必须存在，表示这是一个主机列表。多个主机以逗号分割，但是末尾必须有个逗号

**Modules Options(模块选项)**

- **-a, --args \<MODULE_ARGS>** # 以空格分割的模块参数。格式为 `ARG1=VAL1 ARG2=VAL2......`，注意使用引号，有的 VAL 也需要使用引号引起来
- 只要使用 -m 选项，就必须要是使用 -a 选项
- **-m, --module-name \<ModuleName>** # 执行任务要使用的模块，默认模块为 command。

**Privilege Escalation Options(权限提升选项)**

- **-b, --become** # 使用 become 模块执行所有操作。即开启权限提升功能

## EXAMPLE

- ansible all --list-hosts # 列出所有主机
- 测试 k8s_test_master 组的所有主机的连通性
  - ansible k8s_test_master -m ping
- 对所有管理主机使用默认模块 command 使用 date 命令
  - ansible all -a 'date'
    - 注：也可以使用 `-m 'shell'` 显式得指定 shell 模块。
- 将 resolv.conf 文件中的 nameserver 127.0.0.1 替换成 nameserver 10.8.8.8
  - ansible all -m lineinfile -a "dest=/etc/resolv.conf regexp='nameserver 127.0.1.1' line='nameserver 10.8.8.8'" #

临时测试

`ansible all -i '100.64.0.12:1070,' -u root -k -m ping`

### 常见模块命令示例

> 参考：
>
> - <https://docs.ansible.com/ansible/latest/user_guide/intro_adhoc.html#use-cases-for-ad-hoc-tasks>

- **文件管理**
  - 拷贝文件
    - ansible all -m copy -a "src=/etc/hosts dest=/tmp/hosts"
  - 创建目录，类似 mkdir -p 命令
    - ansible all -m file -a "dest=/tmp/hosts mode=755 owner=desistdaydream group=desistdaydream state=directory"
  - 删除文件
    - ansible all -m file -a "dest=/opt/nginx/config/stream.d/wireguard.conf state=absent"
- **包管理**
  - 安装最新的 net-snmp-utils 包
    - ansible -i inventory/ssc-pool-unicom-ha all -m yum -a "name=net-snmp-utils state=latest"
- **用户和组管理**
  - 创建一个名为 sudo 的组，设置 gid 为 27
    - ansible -i inventory/ssc-pool-datalake-ha -jxgz -m group -a "name=sudo gid=27"
- **cron** # 添加定时任务
  - ansible all -m cron -a 'minute=\*/10 job="/bin/echo hello" name="test1"'
- **script** # 脚本模块，为远程机器执行本地脚本
  - ansible -i ./inventory/ssc-pool-unicom-ha all -m script -a 'scripts.sh'
- **setup** # 收集远程主机的 facts
  - ansible all -m setup # 显示所有被管理节点的相关信息，每个被管理节点，在运行管理命令之前通常会将自己主机相关的信息如，OS 版本，IP 等报告给远程的 ansible 主机
- **unarchive** # 解包
  - ansible all -m unarchive -a "src=/root/downloads/docker-ehualu-20.10.9.tar.gz dest=/"

**检查变量**

检查 root_dir 变量的值

```bash
ansible -i inventory/ my_host -m debug -a "var=root_dir" -C
```

# ansible-doc

显示有关 Ansible 库中安装的模块的信息。 它显示了简短的插件清单及其简短描述，提供了其 DOCUMENTATION 字符串的打印输出，并且可以创建一个简短的“片段”，可以粘贴到 playbook 中。

OPTIONS

- **-l** # 显示所有可用模块
- **-s** # 显示指定的模块中在 playbook 可以定义的参数，可以当作该模块在命令行界面的使用方法

EXAMPLE

# ansible-vault

通过一个密码来加密或解密一个文件或字符串。可以加密 Ansible 使用的任何结构化数据文件。 这可以包括 group_vars/ 或 host_vars/ 清单变量、由 include_vars 或 vars_files 加载的变量，或在 ansible- playbook 命令行上使用 -e @file.yml 或 -e @file.json 传递的变量文件。 还包括角色变量和默认值！

由于 Ansible 任务、处理程序和其他对象都是数据，因此它们也可以使用 Vault 进行加密。 如果您不想公开您正在使用的变量，您可以将单个任务文件完全加密。

该程序加密文件的逻辑是通过管理员输入的密码，对文件进行加密，若想解密，也需要相同的密码才可以解密。与传统意义上通过算法加密的逻辑不太一样。

## Syntax(语法)

**ansible-vault COMMAND**

**COMMAND**

- create
- decrypt
- edit
- view
- encrypt
- encrypt_string
- rekey

# ansible-inventory

显示、转存 Ansible Inventory 信息。默认输出 JSON 格式信息

## Syntax(语法)

**ansible-inventory OPTIONS <--list] | --host HOST | --graph>**

- **--list** # 输出所有主机信息
- **--host \<HOST>** # 输出指定主机信息
- **--graph** #

**OPTIONS**

- **--export** # 使用 --list 时，优化输出内容，以便进行导出；而不是准确表示 Ansible 如何处理清单文件
  - 用人话说：如果不加 --export，那么当多个主机共享了 1 个变量时，输出的主机信息，会给每个主机都添加这个变量。不利于人类阅读与维护
- **-i, --inventory** # 使用 --list 时，将输出的内容保存到指定文件中，而不是输出到标准输出
- **-y, --yaml** # 使用 YAML 格式输出信息

## EXAMPLE

将 Inventory 文件以 YAML 格式输出并保存到文件

- ansible-inventory -i inventory --yaml --list --output dest_inventory.yaml --export
