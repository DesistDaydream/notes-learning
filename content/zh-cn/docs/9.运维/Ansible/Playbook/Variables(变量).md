---
title: Variables(变量)
linkTitle: Variables(变量)
weight: 3
---

# 概述

> 参考：
> 
> - [官方文档，用户指南 - 目录 - 使用变量](https://docs.ansible.com/ansible/latest/user_guide/playbooks_variables.html)

# 变量基本的定义与引用方式

变量名应为字母、数字、下划线。并且始终应该以字母开头。

基础的定义方式是在一个 play 中使用 vars 关键字定义变量，示例如下

```yaml
- host: webservers
  vars:
    http_port: 80 # 定义一个名为http_port的变量，值为80
  tasks:
  - debug
```

Note: vars 关键字可以用在 host 环境中，也可以用在 tasks 环境中，用在 tasks 环境时，则变量仅对当前任务生效

下面是一个在角色中定义字典类型变量的样例：

```yaml
~]# cat roles/test/defaults/main.yml
docker:
  version: 18.06.2
  dataDir: /var/lib/docker
  registryMirrors:
  - https://ac1rmo5p.mirror.aliyuncs.com
  execOpts:
  - 'native.cgroupdriver=systemd'
  insecureRegistries:
  - 100.64.2.52:9999
  - 100.64.1.31:9999
```

### 变量的引用方式

一般情况使用 `{{ VarName }}` 来引用变量，特殊情况不用加 `{{ }}` 而可以直接引用，比如在某些控制结构(比如 when)的语句中。

变量可以通过两种方式引用字典内特定字段的变量

- 使用方括号 `[]` 进行引用
   - `docker['registryMirrors']` 变量的值为 <https://ac1rmo5p.mirror.aliyuncs.com>
- 使用点号 . 进行引用
   - `docker.registryMirrors` 变量的值为 <https://ac1rmo5p.mirror.aliyuncs.com>

Note：如果使用 点号 引用变量可能会引起问题，因为会与 python 字典的属性和方法冲突。所以，尽量使用方括号来引用变量

# 特殊变量

## fact 变量

关闭 fact 变量

```yaml
gather_facts: no
```