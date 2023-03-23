---
title: "Fact Variables"
linkTitle: "Fact Variables"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，Playbook 指南-facts 和 magic 变量](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_vars_facts.html#ansible-facts)

在 ansible 执行任务时，会默认执行名为 Gathering Facts 的任务，以获取目标主机的一些系统信息，如图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nsvz9y/1616125069706-0662e031-1bfe-478b-bb7d-09cb313f4fe0.jpeg)

这些信息以变量的形式体现，每个变量都有其对应的值。可以通过命令 ansible all -m setup 获取这些信息。即通过 [setup 模块](/docs/9.运维/Ansible/Ansible%20Modules/ansible.builtin(内置模块)/System%20模块.md#setup%20-%20收集受管理节点的信息)实现。如下所示

ansible_facts 字段下面的所有字段才是可以直接引用的变量

```json
~]# ansible all -m setup
10.10.100.249 | SUCCESS => {
    "ansible_facts": {
        "ansible_all_ipv4_addresses": [
            "10.10.100.249"
        ],
        "ansible_all_ipv6_addresses": [
            "fe80::47e1:ea44:cfc8:cad0"
        ],
        "ansible_devices": {
            "fd0": {
                "holders": [],
                "host": "",
                "model": null,
                "partitions": {},
                "removable": "1",
                "rotational": "1",
                "scheduler_mode": "deadline",
                "sectors": "0",
                "sectorsize": "512",
                "size": "0.00 Bytes",
                "support_discard": "0",
                "vendor": null
            },
            "sda": {
                "holders": [],
                "host": "SCSI storage controller: LSI Logic / Symbios Logic 53c1030 PCI-X Fusion-MPT Dual Ultra320 SCSI (rev 01)",
                "model": "VMware Virtual S",
                "partitions": {
                    "sda1": {
                        "sectors": "39843840",
                        "sectorsize": 512,
                        "size": "19.00 GB",
                        "start": "2048"
                    }
                },
                "removable": "0",
                "rotational": "1",
                "scheduler_mode": "deadline",
                "sectors": "41943040",
                "sectorsize": "512",
                "size": "20.00 GB",
                "support_discard": "0",
                "vendor": "VMware,"
            },
......后续数据省略
```

可以在 Playbook 中以 `{{ ansible_devices.sda.model }}` 这种方式引用 ansible_devices 下面的 sda 下的 model 变量的值

Note：当进行大规模设备使用 ansible 时，如果每台设备都要获取 fact 信息，ansible 的压力会非常大，这时候推荐关闭 fact 功能，可以在 playbook.yaml 文件中使用 gather_facts 字段即可。如下所示

```yaml
- hosts: WHAT EVER
  gather_facts: no
```

# Facts 关联文件与配置

> 参考：
>
> - 公众号，https://mp.weixin.qq.com/s/HA0vKnuKwKOaB5kdcYX9rg

**/etc/ansible/facts.d/** # 目录是 Ansible 在目标主机上使用的一个配置目录，用于收集系统及应用程序的事实信息（facts）

通过 /etc/ansible/facts.d/ 目录，我们可以自定义该主机的 Facts。这些自定义事实可以通过 Ansible 的 `setup` 模块获取，或者在 Playbook 中通过 `ansible_local` 变量访问。

例如，如果在 `/etc/ansible/facts.d/` 目录中创建了一个名为 `myfact.fact` 的文件，其中包含以下内容：

```ini
[myfact]
mykey=myvalue
```

然后可以在 Playbook 中使用以下代码段访问这个自定义事实：

```yaml
- name: Print custom fact
  debug:
    var: ansible_local.myfact.mykey
```

这将输出 `myvalue`。

## set_fact

在 playbook 中可以使用使用 set_fact 指令设置 facts 变量。
