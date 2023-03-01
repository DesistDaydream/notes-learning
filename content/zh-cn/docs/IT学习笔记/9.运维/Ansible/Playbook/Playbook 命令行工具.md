---
title: Playbook 命令行工具
---

# 概述

> 参考：
> - [官方文档，用户指南-传统目录-使用命令行工具-ansible-playbook](https://docs.ansible.com/ansible/latest/cli/ansible-playbook.html)

ansible-playbook 用来运行运行 Ansible playbook，以便在目标主机上执行定义的任务。

# ansible-playbook

运行 Ansible playbooks，并在目标主机上执行剧本中定义的任务

## Syntax(语法)

**ansible-playbook \[OPTIONS] PLAYBOOK**

### OPTIONS

- \--ask-vault-pass # ask for vault password
- \--become-method # privilege escalation method to use (default=%(default)s), use ansible-doc -t become -l to list valid choices.
- \--become-user # run operations as this user (default=root)
- **-C, --check** # 不在目标主机上执行任务，仅检查任务是否可以完成
- **-C, --connection \<CONNECTION>** # 要使用的连接插件。`默认值：smart`
	- 可以设置为 local 以便让 playbook 在本地执行而不用去远程机器上运行
- \--flush-cache # clear the fact cache for every host in inventory
- \--force-handlers # run handlers even if a task fails
- **-i, --inventory, --inventory-file** # 指定 inventory 文件路径或者以逗号分隔的主机列表。(不推荐使用该选项)
- **-l , --limit \<SUBSET>** # 限定执行的主机范围。可以对一批主机的其中一台执行操作，但是依然可以使用其他主机的变量。further limit selected hosts to an additional pattern
- **--list-hosts** # 列出执行该剧本所能匹配到的主机，但并不会执行
- **--list-tags** # 列出所有可用的 tags
- **--list-tasks** # 列出所有即将被执行的任务
- \--private-key , --key-file # use this file to authenticate the connection
- \--scp-extra-args # specify extra arguments to pass to scp only (e.g. -l)
- \--sftp-extra-args # specify extra arguments to pass to sftp only (e.g. -f, -l)
- \--skip-tags # only run plays and tasks whose tags do not match these values
- \--ssh-common-args # specify common arguments to pass to sftp/scp/ssh (e.g. ProxyCommand)
- \--ssh-extra-args # specify extra arguments to pass to ssh only (e.g. -R)
- **--start-at-task** # start the playbook at the task matching this name
- **--step** # 一步一步运行,也就是说在运行每个任务和之前，会弹出确认信息，确认执行，才会执行该任务。有 3 个选项，执行、不执行、继续(继续就是指后续任务不再确认，从当前任务开始执行完剩下的任务)
- \--syntax-check # perform a syntax check on the playbook, but do not execute it
- **-t, --tags** # 仅运行带有名为 TAG 标签的 tasks 或者 plays
- \--vault-id # the vault identity to use
- \--vault-password-file # vault password file
- \--version # show program’s version number, config file location, configured module search path, module location, executable location and exit
- **-D, --diff** # 当使用 template、file、等指令更改文件时，显示这些文件更改前后的差异。通常与 --check 选项一起使用。
- -K, --ask-become-pass # ask for privilege escalation password
- -M, --module-path # prepend colon-separated path(s) to module library (default=~/.ansible/plugins/modules:/usr/share/ansible/plugins/modules)
- -T , --timeout # override the connection timeout in seconds (default=10)
- -b, --become # run operations with become (does not imply password prompting)
- -c , --connection # connection type to use (default=smart)
- **-e, --extra-vars <@FILE | KEY=VALUE>** # 添加额外的变量，可以是 `KEY=VALUE` 格式(若是 yaml 的话则是 `KEY: VALUE` 格式)，也可以直接指定 yaml 或 json 格式的文件，如果指定文件，以 `@` 开头，比如：
  - `--extra-vars @~/ansible/defaults/main.yaml`
- -f , --forks # specify number of parallel processes to use (default=5)
- -k, --ask-pass # ask for connection password
- -u , --user # connect as this user (default=None)
- **-v, --verbose** # 详细模式(-vvv 会输出更多信息, -vvvv 将会启用 DEBUG 模式)

## EXAMPLE

- 从 install packages 这个任务开始执行 playbook
  - ansible-playbook playbook.yml --start-at="install packages"
- 只对 HLJHEB-PSC-SCORE-PM-OS04-EBRS-HA02 主机执行 playbook
  - ansible-playbook -i inventory/ssc-pool-unicom-ha --limit HLJHEB-PSC-SCORE-PM-OS04-EBRS-HA02 ha-gdas-proxy.yaml

### 常见的本地调试

提前检查渲染的模板。通过 --connectoin=local 以在本地运行，使用 --diff 展示渲染后差异。

- ansible-playbook -i inventory/all.yaml  deploy-mysql.yaml --tag config-mysql --check --diff --connection=local --limit tj-test-spst-node-2
