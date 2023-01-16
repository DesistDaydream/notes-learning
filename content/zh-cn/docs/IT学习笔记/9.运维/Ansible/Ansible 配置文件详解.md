---
title: Ansible 配置文件详解
---

# 概述

> 参考：

Ansible 的配置文件使用 INI 语法书写

# \[defaults]

**deprecation_warnings = <BOOLEAN>** # 是否显示某些功能的弃用警告。`默认值：TRUE`
**host_key_checking = <BOOLEAN>** # 主机 SSH 密钥检查。`默认值：TRUE`。如果启用检查，则对从未 ssh 登录过的主机执行任务将会失败。
**inventory = /etc/ansible/hosts** # 指定 ansible 运行时所用的主机清单路径。默认路径为/etc/ansible/hosts

- Note：可以指定文件或者路径，当指定路径时，则会从该路径下所有文件中读取 host 信息

# \[inventory]

# \[privilege_escalation] 权限提升部分

**become=<BOOLEAN>** # 是否启用以指定用户执行命令。`默认值：False`
**become_method=<STRING>** # 提升权限的方式。`默认值：sudo`
**become_user=<STRING>** # 提升权限所使用的`默认值：root`
**become_ask_pass=<BOOLEAN>** # `默认值：True`

# \[paramiko_connection]

# \[ssh_connection]

**transfer_method = <STRING>** # 传输文件的机制。`默认值：smart`。该指令代替旧版的 scp_if_ssh 指令

- sftp # 仅使用 sftp
- scp # 仅使用 scp
- piped # 仅使用 piped
- smart # 按照顺序尝试 sftp、scp、piped

# \[persistent_connection]

# \[accelerate]

# \[selinux]

# \[colors]

# \[diff]
