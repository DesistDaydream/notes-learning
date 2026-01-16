---
title: Ansible Variables
---

# group_vars 概述

> 参考：
>
> - [官方文档，传统目录 - 使用变量](https://docs.ansible.com/ansible/latest/user_guide/playbooks_variables.html)
> - [官方文档，传统目录 - 使用变量 - 变量优先级](https://docs.ansible.com/ansible/latest/user_guide/playbooks_variables.html#variable-precedence-where-should-i-put-a-variable)

虽然通过自动化可以使事情更简单、更可重复，但是并非所有系统都完全相同。在某些情况下，观察到一个系统的行为或状态可能会影响到配置其他系统的方式。比如，我们可能需要找出一个系统的 IP 地址，并将这个 IP 地址作为另一个系统中配置的值。

基于上述目的，Ansible 可以通过 **Variables(变量)** 来管理各个系统之间的差异。

Ansible 的变量就跟编程语言中的变量概念一样，同样可以定义、引用。我们使用标准的 YAML 语法创建变量，包括列表和字典；可以这么说，YAML 中每个字段的 key 就是变量名，value 就是变量的值。我们可以在 Playbooks、Inventory、甚至命令行中定义与引用变量。我们还可以在 Playbooks 运行期间，将任务的返回值注册为变量，以创建一个新的变量。

创建变量后，我们可以在 模块的参数、模板、控制结构 中使用这些变量。在 [GitHub 中有一个 Ansible 示例的目录](https://github.com/ansible/ansible-examples)，可以看到很多 Ansible 使用变量的例子

下面的示例就是在命令行中使用 debug 模块，查看了一下 inventory_hostname 这个默认变量的值

```bash
~]$ ansible -i ../inventory/ all -m debug -a 'msg={{inventory_hostname}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "hw-cloud-xngy-jump-server-linux-2"
}
```

# 变量的优先级

变量可以是自带的，就是由人们自行定义的，可以在多个地方定义变量，(e.g.在某些文件里定义变量、通过命令行传递变量等等。由于 ansible 所要处理的的文件有很多，不同类型的文件下定义的变量的优先级也不同)

下面的优先级列表由低到高，最下面的变量优先级最高

- command line values (eg “-u user”)
- **role defaults** # 定义在 `${ROLE}/defaults/main.yaml` 中的默认变量
- **inventory file or script group vars** # [Inventory 文件](/docs/9.运维/Ansible/Inventory%20文件.md#组变量)中的组变量，i.e. `[XXX:vars]`
- **inventory group_vars/all** # Inventory 文件所在目录下的 `group_vars/all` 文件。也可以是  `group_vars/all.yaml` 文件
- **playbook group_vars/all** # Playbook 根目录下的 `group_vars/all` 文件。也可以是  `group_vars/all.yaml` 文件
- **inventory group_vars/** # Inventory 文件所在目录下的 `group_vars/` 目录
- **playbook group_vars/** # Playbook 根目录下的 `group_vars/` 目录
- **inventory file or script host vars** # [Inventory 文件](/docs/9.运维/Ansible/Inventory%20文件.md#主机变量)中的主机变量。
- **inventory host_vars/** # Inventory 文件所在目录下的 `host_vars/` 目录
- **playbook host_vars/** # Playbook 根目录下的 `host_vars/` 目录
- **host facts / cached set_facts** #
- play vars #
- play vars_prompt #
- play vars_files #
- **role vars** # 定义在 `${ROLE}/vars/main.yml` 中的变量。针对每个 [Playbook Role(角色)](docs/9.运维/Ansible/Playbook/Playbook%20Role(角色).md) 的变量
- block vars (only for tasks in block) #
- task vars (only for the task) #
- include_vars #
- set_facts / registered vars #
- role (and include_role) params #
- include params #
- **extra vars** # 通过 ansible-playbook 命令行工具的 `-e, --extra-vars` 参数指定的变量

Note：可以说 ansible playbook 中写的所有内容都是变量。都是可以引用的，只不过引用的方式不同。

# 变量的定义与引用

变量名应为字母、数字、下划线。并且始终应该以字母开头。可以在 Inventory、Playbooks、命令行 中定义变量。Ansible 会加载它找到的每个可能的变量，然后根据[变量优先级规则](#变量的优先级)选择要应用的变量

可以通过 -e 选项直接定义一个变量，比如 `ansible -e "test_var=hello_world"`，这里定义了 test_var 变量，变量的值为 hello_world。

Ansible 使用 Jinja2 语法引用变量。Jinjia2 使用 `{{ VarName }}` 来引用变量，比如

```bash
~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "test_var=hello_world" -m debug -a 'msg={{test_var}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "hello_world"
}
```

上面的例子中，我们定义了一个名为 test_var 的变量，变量的值为 hello_world，并使用 debug 模块，引用 test_var 变量。

这只是最简单的变量的使用方式，命令行中不适合设置复杂格式的变量，更为复杂的类型的变量，通常在 YAML 或 JSON 格式的文件中定义，并直接引用文件即可定义变量(比如使用 `--extra-vars "@./test_var.yaml"` 选项，即可通过 test_var.yaml 文件定义变量)

## 变量的类型

### List(列表)变量

```bash
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ cat test_var.yaml
region:
- northeast
- southeast
- midwest
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "@./test_var.yaml" -m debug -a 'msg={{region[1]}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "southeast"
}
```

### Dictionary(字典)变量

可以通过两种方式引用字典变量

- 使用方 `[]` 进行引用
  - foo\['field1']
- 使用 `.` 进行引用(不推荐使用该方式引用变量，可能会与 Python 语法产生冲突)
  - foo.field1

```bash
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ cat test_var.yaml
foo:
  field1: one
  field2: two
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "@./test_var.yaml" -m debug -a msg="{{foo['field1']}}"
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "one"
}
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "@./test_var.yaml" -m debug -a msg="{{foo.field1}}"
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "one"
}
```

Note：如果使用 `.` 引用变量可能会引起问题，因为会与 python 字典的属性和方法冲突。所以，尽量使用 `[]` 来引用变量

### Registering(注册)变量

Registering 类型的变量适用于 Playbooks 中，通过 `register` 关键字将任务中的返回值注册为指定的变量，然后可以在 Playbooks 的后续任务中，引用注册的变量
比如

```yaml
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ cat roles/variables/tasks/main.yaml
- name: test
  command: whoami
  register: info
- name: debug
  debug:
    msg: "{{info}}"
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ ansible-playbook -i ../inventory/ variables.yaml

PLAY [test] **********************************************************************************************************************************************************************************************************************************************************************************************************************************************************

TASK [variables : test] **********************************************************************************************************************************************************************************************************************************************************************************************************************************************
changed: [hw-cloud-xngy-jump-server-linux-2]

TASK [variables : debug] *********************************************************************************************************************************************************************************************************************************************************************************************************************************************
ok: [hw-cloud-xngy-jump-server-linux-2] => {
    "msg": {
        "ansible_facts": {
            "discovered_interpreter_python": "/usr/bin/python3"
        },
        "changed": true,
        "cmd": [
            "whoami"
        ],
        "delta": "0:00:00.002390",
        "end": "2021-10-11 22:57:18.455061",
        "failed": false,
        "rc": 0,
        "start": "2021-10-11 22:57:18.452671",
        "stderr": "",
        "stderr_lines": [],
        "stdout": "root",
        "stdout_lines": [
            "root"
        ]
    }
}

PLAY RECAP ***********************************************************************************************************************************************************************************************************************************************************************************************************************************************************
hw-cloud-xngy-jump-server-linux-2 : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

```

在 test 任务中，我们将 command 模块执行的任务返回值，注册到 info 变量中；然后再 debug 任务中，引用了 info 变量。

有关在后续任务的条件中使用注册变量的更多示例，请参阅[条件](https://docs.ansible.com/ansible/latest/user_guide/playbooks_conditionals.html#playbooks-conditionals)。注册变量可以是简单变量、列表变量、字典变量或复杂的嵌套数据结构。每个模块的文档包括 RETURN 描述该模块返回值的部分。要查看特定任务的值，请使用-v.
注册的变量存储在内存中。您不能缓存已注册的变量以供将来使用。注册的变量仅在当前 playbook 运行的其余部分在主机上有效。

注册变量是主机级变量。当您使用循环在任务中注册变量时，注册的变量包含循环中每个项目的值。循环期间放置在变量中的数据结构将包含一个 results 属性，即来自模块的所有响应的列表。有关其工作原理的更深入示例，请参阅有关将寄存器与循环一起使用的[循环](https://docs.ansible.com/ansible/latest/user_guide/playbooks_loops.html#playbooks-loops)部分。

> 注意：如果任务失败或被跳过，Ansible 仍会注册一个处于失败或跳过状态的变量，除非根据标签跳过该任务。有关添加和使用标签的信息，请参阅[标签](https://docs.ansible.com/ansible/latest/user_guide/playbooks_tags.html#tags)。

### Nested(嵌套)变量

```bash
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ cat test_var.yaml
foo:
- field1:
    name: one
- field2:
    name: two
[desistdaydream@hw-cloud-xngy-jump-server-linux-2 ~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "@./test_var.yaml" -m debug -a msg="{{foo[0].field1.name}}"
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "one"
}

```

## 变量的定义方式

### 在命令行中定义变量

在使用 `ansible` 或 `ansible-playbook` 命令时，可以通过 --extra-vars 或 -e 选项，以在命令行中定义变量

可以通过多种方式在命令行定义变量

- **KEY=VALUE**

```bash
~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "test_var=hello_world" -m debug -a 'msg={{test_var}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "hello_world"
}
```

- **JSON 字符串**

```bash
~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars '{"test_var":"hello world"}' -m debug -a 'msg={{test_var}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "hello world"
}
```

- **来自 JSON 或 YAML 文件**

```bash
~/projects/DesistDaydream/ansible/playbooks]$ cat test_var.yaml
test_var: 'hello world'
~/projects/DesistDaydream/ansible/playbooks]$ ansible -i ../inventory/ all --extra-vars "@./test_var.yaml" -m debug -a 'msg={{test_var}}'
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "hello world"
}
```

### 在 Inventory 中定义变量

详见 [Inventory 文件](/docs/9.运维/Ansible/Inventory%20文件.md)

### 在 Playbooks 中定义变量

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
[root@cobbler playbook]# cat roles/test/defaults/main.yml
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

## 变量的引用方式

特殊情况不用加 `{{ }}` 而可以直接引用，比如在某些控制结构(比如 when)的语句中。

# Registering Variables(注册变量)

通常在剧本中，将给定命令的结果存储在变量中并在以后访问它可能很有用。

注意：

- 注册变量只适用于执行本注册任务的 host。假如在 host_A 注册了 Var_A，那么 host_B 想要引用 Var_A，则必须使用 `hostvars` 变量。

应用示例：

```yaml
- hosts: all
    tasks:
    - name: list contents of directory
    command: ls /root/
    register: contents # 将该任务执行后的ansible报告的信息保存在名为contents变量中
    - debug:
        msg: "{{contents}}" # 输出contents变量
    - debug：
        msg: "{{contents.stdout}}" # 输出contents下的stdout变量的值，值为anaconda-ks.cfg\nScripts
```

比如下面，就是是 contents 变量的值。这其中包括要执行的命令、命令执行的日期、执行结果，等等 ansible 执行该 playbook 后的信息。

```json
    TASK [debug] *************************************************************
    ok: [10.10.100.200] => {
        "msg": {
            "changed": true,
            "cmd": [
                "ls",
                "/root/"
            ],
            "delta": "0:00:00.004220",
            "end": "2019-11-11 15:02:17.326659",
            "failed": false,
            "rc": 0,
            "start": "2019-11-11 15:02:17.322439",
            "stderr": "",
            "stderr_lines": [],
            "stdout": "anaconda-ks.cfg\nScripts",
            "stdout_lines": [
                "anaconda-ks.cfg",
                "Scripts"
            ]
        }
    }
```

还可以将 register 与循环配合使用，通过命令获取的多个值注册到变量中，然后使用循环逐一读取变量的值

```yaml
    - name: retrieve the list of home directories
      command: ls /home
      register: home_dirs
    - name: add home dirs to the backup spooler
      file:
        path: /mnt/bkspool/{{ item }}
        src: /home/{{ item }}
        state: link
      loop: "{{ home_dirs.stdout_lines }}" # loop也可以使用这样的方式来获取每一行的值: "{{ home_dirs.stdout.split() }}"
```

这个例子就是查看/mnt/bkspool/目录下的内容，然后将其中所有文件注意拷贝到/home/目录下

# Special Variables(特殊的变量)

> 参考：
>
> - [官方文档，特殊变量](https://docs.ansible.com/ansible/latest/reference_appendices/special_variables.html)

无论是否定义任何变量，都可以使用 Ansible 提供的特殊变量访问有关主机的信息，一共有如下几种变量类型：

- **magic variables(魔法变量)**
- **facts variables(事实变量)**
- **connection variables(连接变量)**

## Magic Variables

> 官方文档：<https://docs.ansible.com/ansible/latest/reference_appendices/special_variables.html#magic>

魔术变量不能随意覆盖并且也没法覆盖，这是一种 Ansbie 提供的"内部变量" ，可以反映 Ansible 所管理主机的最简单的基本状态，比如该主机的主机名、在 inventory 文件中的定义都会转换成这里面变量的值、等等。

可以通过目标主机获取到 ansible 管理的所有主机的信息。最常用的魔术变量有以下几个

- **hostvars** # 每个目标主机下面都包含类似下图的信息。其中是每个组所包含的 hosts
  - 注意：通过 hostvars 变量，我们还可以获取到其他主机在执行任务是注册的变量，比如在 kubernetes 集群的 master-1 上生成了加入集群的指令，并注册为变量 join_cmd，正常是无法在其他主机直接使用的。这时候就要用到 hostvars 变量了。
  - ![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nsvz9y/1616125069735-9fbbff13-76a7-455a-9a5b-291800f65cc1.jpeg)
- **ansible_play_hosts** # 一个列表，是当前 play 中活动的主机列表，受序号限制，无法访问的主机不会被当做“活动”主机。
  - 该变量可以用于 for 循环，对列表中的主机进行遍历，逐一操作。
  - 等同于 ansible_play_batch
- **ansible_play_name** # 当前执行 paly 的名称。i.e.playbook 中 hosts 这个键的值，也就是当前的主机组名称
- **groups** # 默认值为 inbentory 下所有组及其组内的 host
- **group_names** # 默认值为当前主机所属组的列表。
- **inventory_hostname** # 默认值为 inventory 文件中配置的主机名称。i.e. ansible 的 hosts 文件的第一列内容
- **inventory_dir** # 默认值为 ansible 保存 hosts 文件的目录的绝对路径。默认路径为/etc/ansible/
- **play_hosts** # 默认值为当前 play 范围中可用的一组主机名
- **role_path** # 默认值为当前 role 的目录的绝对路径

应用实例：

**groups\["{{ansible\_play\_name}}"]** # 获取当前 play 下的主机列表

## Fact Variables

Ansible 在执行任务之前，会收集被控制节点的系统及应用程序的 facts(可以理解为：事实信息)。并将这些信息存储到 `${ansible_facts}` 变量中以供后续使用。

详见 《[Fact Variables](/docs/9.运维/Ansible/Ansible%20Variables/Fact%20Variables.md)》

## Connection Variables

# 应用示例

## 获取组中的主机数量

```yaml
      vars:
        HOST_COUNT: "{{ groups['组名'] | length }}"
```

获取 test 组中主机的总数量

```bash
~]$ ansible -i ../inventory/ all -m debug -a "msg={{ groups['test'] | length }}"
hw-cloud-xngy-jump-server-linux-2 | SUCCESS => {
    "msg": "1"
}
```
