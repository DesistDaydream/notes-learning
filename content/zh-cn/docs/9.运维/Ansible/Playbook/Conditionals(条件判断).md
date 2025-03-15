---
title: Conditionals(条件判断)
weight: 5
---

# 概述

> 参考：
>
> - [官方文档，用户指南 - 传统目录 - 条件](https://docs.ansible.com/ansible/latest/user_guide/playbooks_conditionals.html)

通常，play 的结果可能取决于 variable，fact（有关远程系统的知识）或先前的任务结果。在某些情况下，变量的值可能取决于其他变量。可以基于主机是否符合其他条件来创建其他组来管理主机。

Ansible 在条件中使用 Jinja 的 [测试](https://docs.ansible.com/ansible/latest/user_guide/playbooks_tests.html) 和 [过滤器](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html)来实现条件判断。详见 [Ansible Template 文章中《Ansible 扩展测试函数》章节](/docs/9.运维/Ansible/Playbook/Templates%20模板(Jinja2).md#Ansible%20扩展的测试函数)


# 基于变量的条件

https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_conditionals.html#conditionals-based-on-variables

您还可以根据剧本或库存中定义的变量创建条件。由于条件需要布尔值输入（必须评估测试以触发条件），因此您必须应用| Bool过滤到非树状变量，例如带有“是”，“ ON”，“ 1”或“ TRUE”的内容的字符串变量。您可以定义这样的变量：

```yaml
vars:
  epic: true
  monumental: "yes"
```

对于上面的变量，Ansible将运行其中一个任务并跳过另一个任务：

```yaml
tasks:
    - name: Run the command if "epic" or "monumental" is true
      ansible.builtin.shell: echo "This certainly is epic!"
      when: epic or monumental | bool

    - name: Run the command if "epic" is false
      ansible.builtin.shell: echo "This certainly isn't epic!"
      when: not epic
```

如果尚未设置所需的变量，则可以使用Jinja2的定义测试跳过或失败。例如：

```yaml
tasks:
    - name: Run the command if "foo" is defined
      ansible.builtin.shell: echo "I've got '{{ foo }}' and am not afraid to use it!"
      when: foo is defined

    - name: Fail if "bar" is undefined
      ansible.builtin.fail: msg="Bailing out. This play requires 'bar'"
      when: bar is undefined
```

# 最佳实践

## 简单样例

下面的样例表示：当 `ansible_facts['os_family']` 变量的值为 Debian 的时候，则执行上面的任务，任务内容是执行 shutdown 命令

```yaml
tasks:
- name: "shut down Debian flavored systems"
  command: /sbin/shutdown -t now
  when: ansible_facts['os_family'] == "Debian"
```

判断主机是否在某个组中

- when: inventory_hostname is search("master") # 当 inventory_hostname 变量的值含有 master 字符串时。
- when: inventory_hostname == groups\['kube_master']\[0] # 当当前主机的 inventory_hostname 变量值等于 kube_master 主机组中的第一台主机时
- when: inventory_hostname in groups\['kube_master'] # 当当前主机的 inventory_hostname 变量值在 kube_master 主机组中时。

- when: testvar1 is none # 当变量 testvar1 已定义，但是值为空时。Note：值为空表示 key 后面的值什么都不写，双引号都不能有
- when: ((groups\['kube_master'] | length) > 1) # 当 kube_master 主机组的主机数量大于 1 时。

## 每次处理变量是否要判断该变量是否已被定义？

在对变量的值进行判断时，尽量先判断变量是否定义，再判断其值是否为某个值，<font color="#ff0000">避免当变量没定义时报错</font>

比如在模板渲染的场景中

- `{% if hostvars[target]['ipmi_ip'] is defined and hostvars[target]['ipmi_ip'] != "" %}`
- 和
- `{% if hostvars[target]['ipmi_ip'] != "" %}` 

比如在任务的 when 条件中

```yaml
- name: 任务1
  when:
    - scrape_hds is defined
    - scrape_hds | bool
```

和

```yaml
- name: 任务1
  when:
    - scrape_hds | bool
```

上面这几个例子中的前者不会因为变量没有定义而报错

