---
title: Tags
linkTitle: Tags
weight: 20
date: 2025-01-06T22:09:00
---

# 概述

> 参考：
>
> - [官方文档，执行 Playbooks - 标签](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_tags.html#)

如果您有一个很大的剧本，则仅运行其中的特定部分而不是运行整个剧本可能会很有用。您可以使用 Ansible 标签来做到这一点。使用标签执行或跳过选定的任务是一个两步过程：

1. 将标签添加到您的任务中，可以单独添加标签，也可以使用继承自 block、play、role 或 import 的标签
2. 运行 playbook 时选择或跳过标签

# 继承 Tags

https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_tags.html#tag-inheritance-for-includes-blocks-and-the-apply-keyword
