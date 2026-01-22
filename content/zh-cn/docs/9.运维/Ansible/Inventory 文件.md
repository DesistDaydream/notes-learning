---
title: Inventory 文件
linkTitle: Inventory 文件
weight: 3
---

# 概述

> 参考：
>
> - [官方文档，用户指南 - 如何建立你的 Inventory](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html)

Ansible 可同时操作属于一个组的多台主机，组和主机之间的关系通过 Inventory 文件配置。默认的文件路径为 /etc/ansible/hosts，也可以在 `ansible`、`ansible-playbook` 命令中使用 -i 选项指定其他的 Inventory 文件。

除默认文件外，还可以同时使用多个 inventory 文件，也可以从动态源，或云上拉取 inventory 配置信息。详见[动态 Inventory](https://docs.ansible.com/ansible/latest/inventory_guide/intro_dynamic_inventory.html#intro-dynamic-inventory)。

## Inventory文件格式

最常见的格式是 [INI](/docs/2.编程/无法分类的语言/INI.md) 和 [YAML](/docs/2.编程/无法分类的语言/YAML.md) 格式，下面这是一个 INI 格式的 Inventory 示例

```ini
# 例1:定义一个单独的主机。未分组的机器。Note:需要在“例2”中中括号定义组之前指定
green.example.com
192.168.100.1
# 例2:定义一个主机组。组名为webservers的主机集合
[webservers]
alpha.example.org
192.168.1.100
# 再定义一个主机组。组名为dbservers的主机集合
[dbservers]
192.168.2.100
# 定义一个主机的另一种方式。使用正则表达式来指定多个主机
www[001:006].example.com #www001.example.com一直到www006.example.com一共6台主机
db-[a:f]-node.example.com
# 可以在主机ip或主机名后面添加参数，以此来控制Ansible与远程主机的交互方式。
# 详细的参数信息见官网：https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#connecting-to-hosts-behavioral-inventory-parameters
# 指定该主机要执行操作的主机ip、用户名和密码
www.desistdaydream.com ansible_host=10.10.100.200 ansible_user="root" ansible_password="my@password"
# 例3:组的引用,可以把一个或多个组作为另一个组的子成员。all_host 组包含 webservers 与 dbservers 两个组中所有的 hosts
[all_host:children]
webservers
dbservers
```

Note：该文件的第一列为一般推荐使用主机名来表示，如果需要指定该主机的ip地址，则使用 ansbile_ssh_host 或者 ansible_host 参数来指定ip。

因为，Ansible 默认变量 inventory_hostname 的值为 inventory 文件中第一列的内容，还有另一个变量是 inventory_hostname_short ，这个变量的值是主机名的短格式

所以，在 Ansbile 里，自动就会将第一列认定为主机名，如果使用 ip 作为第一列的表示形式，那么与 Ansible 理念不符(至于为什么还可以用 ip 表示，可能是为了大家方便，所以第一列才可以使用 ip 的吧~~~)

对应的 YAML 格式 Inventory 文件示例：

```yaml
# 定义一个名为 all_host 的组，通过 children 字段把其他组引入
all_host:
  # 例:组的引用,可以把一个或多个组作为另一个组的子成员。all_host 组包含 webservers 与 dbservers 两个组中所有的 hosts
  children:
    # 例:定义一个主机组。组名为 dbservers 的主机集合
    dbservers:
      hosts:
        192.168.2.100:
        # 定义一个主机的另一种方式。使用正则表达式来指定多个主机
        db-[a:f]-node.example.com:
        www[001:006].example.com:
        # 可以在主机ip或主机名后面添加参数，以此来控制Ansible与远程主机的交互方式。
        # 详细的参数信息见官网：https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#connecting-to-hosts-behavioral-inventory-parameters
        # 指定该主机要执行操作的主机ip、用户名和密码
        www.desistdaydream.com:
          ansible_host: 10.10.100.200
          ansible_password: my@password
          ansible_user: root
    # 例:定义一个主机组。组名为 webservers 的主机集合
    webservers:
      hosts:
        192.168.1.100:
        alpha.example.org:
# 例:定义一个单独的主机。未分组的机器。
ungrouped:
  hosts:
    192.168.100.1:
    green.example.com:

```

# 主机与组

/etc/ansible/hosts 文件支持 [INI](/docs/2.编程/无法分类的语言/INI.md)、[YAML](/docs/2.编程/高级编程语言/Go/Go%20规范与标准库/文本处理/YAML.md) 语法书写：

```ini
mail.example.com

[webservers]
foo.example.com
bar.example.com

[dbservers]
one.example.com
two.example.com
three.example.com
```

这里的 `[]` 中的字符是组名，用于对系统进行分类,便于对不同系统进行个别的管理。

一个系统可以属于不同的组,比如一台服务器可以同时属于 webserver组 和 dbserver组. 这时属于两个组的变量都可以为这台主机所用,至于变量的优先级关系将于以后的章节中讨论.

如果有主机的SSH端口不是标准的22端口, 可在主机名之后加上端口号,用冒号分隔. SSH 配置文件中列出的端口号不会在 paramiko 连接中使用,会在 openssh 连接中使用.

端口号不是默认设置时,可明确的表示为:

```
badwolf.example.com:5309
```

假设你有一些静态IP地址,希望设置一些别名,但不是在系统的 host 文件中设置,又或者你是通过隧道在连接,那么可以设置如下:

```
jumper ansible_ssh_port=5555 ansible_ssh_host=192.168.1.50
```

在这个例子中,通过 “jumper” 别名,会连接 192.168.1.50:5555.记住,这是通过 inventory 文件的特性功能设置的变量. 一般而言,这不是设置变量(描述你的系统策略的变量)的最好方式.后面会说到这个问题.

一组相似的 hostname , 可简写如下:

```
[webservers]
www[01:50].example.com
```

数字的简写模式中,01:50 也可写为 1:50,意义相同.你还可以定义字母范围的简写模式:

```
[databases]
db-[a:f].example.com
```

对于每一个 host,你还可以选择连接类型和连接用户名:

```
[targets]
localhost ansible_connection=local
other1.example.com ansible_connection=ssh ansible_ssh_user=mpdehaan
other2.example.com ansible_connection=ssh ansible_ssh_user=mdehaan
```

所有以上讨论的对于 inventory 文件的设置是一种速记法,后面我们会讨论如何将这些设置保存为 ‘host_vars’ 目录中的独立的文件.

## 默认组

Inventory 文件中有两个默认的组，名称为：`all` 和 `ungrouped`(这两个名称是隐藏的)。all 组包含所有主机，ungrouped 组包含除了 all 有之外没有属组的主机。

每个主机至少属于 2 个组

- all 和 ungrouped
- all 和 某组

## 子组

> [!Warning]
> Ansible 组名全局唯一，没有树的逻辑
>
> Ansible 的所有组都是平级的，<font color="#ff0000">只能让一个组引用另一个组，而不能让一个组属于另一个组</font>
>
> 如果 a 组中有名为 dev 的组，b 组中有名为 dev 的组。那么当处理 a 组时，会连同 b 组的 dev 一起处理。
>
> https://stackoverflow.com/questions/71032136/why-in-ansible-i-cant-have-same-key-in-many-children-groups
>
> https://stackoverflow.com/questions/64867207/ansible-intersection-of-two-host-groups-using-yaml-inventory
>
> PS: 既然没有从属逻辑，为啥要用 children 这种关键字和描述。。。( ╯□╰ ) 与其叫子组，不如叫引用组

要让某个组引用另一个组，使用 `children` 关键字

```yaml
east:
  hosts:
    foo.example.com:
    one.example.com:
    two.example.com:
west:
  hosts:
    bar.example.com:
    three.example.com:
prod:
  children:
    east:
test:
  children:
    west:
```

这个 Inventory 的效果与下面这个相同

```yaml
east:
  hosts:
    foo.example.com:
    one.example.com:
    two.example.com:
west:
  hosts:
    bar.example.com:
    three.example.com:
prod:
  hosts:
    foo.example.com:
    one.example.com:
    two.example.com:
test:
  hosts:
    bar.example.com:
    three.example.com:
```

# 主机变量与组变量

## 主机变量

前面已经提到过，分配变量给主机很容易做到，这些变量定义后可在 playbooks 中使用

```
[atlanta]
host1 http_port=80 maxRequestsPerChild=808
host2 http_port=303 maxRequestsPerChild=909
```

## 组变量

也可以定义属于整个组的变量

```
[atlanta]
host1
host2
[atlanta:vars]
ntp_server=ntp.atlanta.example.com
```

把一个组作为另一个组的子成员

可以把一个组作为另一个组的子成员,以及分配变量给整个组使用. 这些变量可以给 /usr/bin/ansible-playbook 使用,但不能给 /usr/bin/ansible 使用:

```
[atlanta]
host1
host2
[raleigh]
host2
host3
[southeast:children]
atlanta
raleigh
[southeast:vars]
some_server=foo.southeast.example.com
halon_system_timeout=30
self_destruct_countdown=60
escape_pods=2
[usa:children]
southeast
northeast
southwest
northwest
```

如果我们需要存储一个列表或 hash 值，或者更喜欢把 host 和 group 的变量分开配置，请看下一节的说明.

# 组织 host_vars(主机变量) 和 group_vars(组变量)

在 Inventory 主文件中保存所有的变量并不是最佳的方式。我们通常在**独立的文件**中定义这些变量，这些独立文件与 inventory 文件保持关联. 不同于 inventory 文件(INI 格式)，这些独立文件的格式为 YAML。

假设有一个主机名为 ‘foosball’，主机同时属于两个组

- raleigh
- webservers

那么以下配置文件中的变量可以为 ‘foosball’ 主机所用。依次为 ‘raleigh’ 的组变量，’webservers’ 的组变量，’foosball’ 的主机变量：

```
${PrefixDirPath}/group_vars/raleigh
/etc/ansible/group_vars/webservers
/etc/ansible/host_vars/foosball
```

> [!Attention]
>
> `group_vars/` 目录下文件名必须是**组名**才可以将变量的值应用相同组名的组中的主机
>
> - 如上所示：group_vars/raleigh 中的组变量适用于 raleigh 组。
> - 文件名也可以使用 all 和 ungrouped 用于为所有主机或所有未分组的主机定义变量。
>
> `host_vars/` 目录下文件名必须时**主机名**（Inventory 定义的主机名，不是系统中的主机名）

举例来说, 假设你有一些主机，属于不同的数据中心，并依次进行划分。每一个数据中心使用一些不同的服务器。比如 ntp 服务器，database 服务器等等。那么 ‘raleigh’ 这个组的组变量定义在文件 ‘/etc/ansible/group_vars/raleigh’ 之中，可能类似这样：

```
---
ntp_server: acme.example.org
database_server: storage.example.org
```

这些定义变量的文件不是一定要存在，因为这是可选的特性。

还有更进一步的运用,你可以为一个主机，或一个组，创建一个目录，目录名就是主机名或组名。目录中的可以创建多个文件，文件中的变量都会被读取为主机或组的变量。如下 ‘raleigh’ 组对应于 /etc/ansible/group_vars/raleigh/ 目录，其下有两个文件 db_settings 和 cluster_settings，其中分别设置不同的变量：

```
/etc/ansible/group_vars/raleigh/db_settings
/etc/ansible/group_vars/raleigh/cluster_settings
```

‘raleigh’ 组下的所有主机，都可以使用 ‘raleigh’ 组的变量。当变量变得太多时，分文件定义变量更方便我们进行管理和组织。还有一个方式也可参考，详见 Ansible Vault 关于组变量的部分。注意，分文件定义变量的方式只适用于 Ansible 1.4 及以上版本。

我们可以将 `group_vars/` 和 `host_vars/` 目录添加到 playbook 目录下。如果两个目录下都存在,那么 playbook 目录下的配置会覆盖 inventory 目录的配置。

把我们的 Inventory 文件 和 变量 放入 git repo 中，以便跟踪他们的更新，这是一种非常推荐的方式。

# 主机匹配模式

> 参考：
>
> - [官方文档，用户指南 - 传统目录 - 模式：针对主机和组](https://docs.ansible.com/ansible/latest/user_guide/intro_patterns.html)

主机列表的正则匹配

ansible支持主机列表的正则匹配

- 全量: `all/*`
- 逻辑或: `:`
- 逻辑非: `!`
- 逻辑与: `＆`
- 切片： `[]`
- 正则匹配： 以 `~` 开头

- `ansible all -m ping` # 所有默认inventory文件中的机器
- `ansible "*" -m ping` # 同上
- `ansible 121.28.13.X -m  ping` # 所有122.28.13.X机器
- `ansible  web1:web2  -m  ping` # 所有属于组web1或属于web2的机器
- `ansible  web1:!web2  -m  ping` # 属于组web1，但不属于web2的机器
- `ansible  web1&web2  -m  ping` # 属于组web1又属于web2的机器
- `ansible webserver[0]  -m  ping` # 属于组webserver的第1台机器
- `ansible webserver[0:5]  -m  ping` # 属于组webserver的第1到4台机器
- `ansible "~(beta|web).example.(com|org)"  -m ping`

# Inventory 参数详解

> 参考：
>
> - [官方文档，用户指南 - 传统目录 - 如何构建你的 Inventory - 连接到主机:Inventory 参数](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#connecting-to-hosts-behavioral-inventory-parameters)

> Tips: 配置 Inventory 中主机的参数本质就是配置使用哪个 Ansible Plugin 或者 Ansible  Module 来连接受管理节点

如同前面提到的，通过设置下面的参数，可以控制 ansible 与远程主机的交互方式：

**ansible_connection** # 指定 ansible 与远程主机的 connector(连接器)，默认为 ssh 的 smart 类型。'smart' 方式会根据是否支持 ControlPersist，来判断 'ssh' 方式是否可行

- smart、ssh、paramiko # 这三种类型都是 ssh 连接器下的类型。默认为 smart
- local
- docker

通用连接参数

- **ansible_host** # 将要连接的远程主机名.可以设为ip
- **ansible_port** # 将要连接的远程主机端口号.默认端口为22
- **ansible_user** # 将要连接的远程主机的用户名
- **ansible_password** # 将要连接的远程主机的密码(这种方式并不安全,我们强烈建议使用 --ask-pass 或 -k 或 SSH 密钥)

只适用于SSH连接所用参数

- **ansible_sudo_pass** # sudo 密码(这种方式并不安全,我们强烈建议使用 --ask-sudo-pass)
- **ansible_sudo_exe** # sudo 命令路径(适用于1.8及以上版本)
- **ansible_ssh_private_key_file** # ssh 使用的私钥文件.适用于有多个密钥,而你不想使用 SSH 代理的情况.
- **ansible_shell_type** # 目标系统的shell类型.默认情况下,命令的执行使用 'sh' 语法,可设置为 'csh' 或 'fish'.
- **ansible_ssh_extra_args** # 参数值可以作为 [ssh](/docs/4.数据通信/Utility/OpenSSH/ssh.md) 的命令行参数
  - e.g. `ansible_ssh_extra_args: "-o HostKeyAlgorithms=+ssh-dss"` 相当于 `ssh -o HostKeyAlgorithms=+ssh-dss`

> 更多 SSH 相关参数可以参考 [Connection Plugins](/docs/9.运维/Ansible/Ansible%20Plugins/Connection%20Plugins.md) 的 SSH 章节

权限提升参数

- **ansible_becom=yes|no** # 是否允许提升权限执行操作。
- **ansible_become_user=\<STRING>** # 权限提升执行操作时所使用的用户。`默认值：root`
- **ansible_become_password=\<STRING>** # 权限提升执行操作时所使用的用户的密码。(这种方式并不安全,我们强烈建议使用 --ask-become-pass 或 -K)

远程主机环境参数

- **ansible_python_interprete** # 目标主机的 python 路径.适用于的情况: 系统中有多个 Python, 或者命令路径不是"/usr/bin/python",比如  \*BSD, 或者 /usr/bin/python
  - 不是 2.X 版本的 Python.我们不使用 "/usr/bin/env" 机制,因为这要求远程用户的路径设置正确,且要求 "python" 可执行程序名不可为 python以外的名字(实际有可能名为python26). 与 ansible_python_interpreter 的工作方式相同,可设定如 ruby 或 perl 的路径....

一个主机文件的例子:

```text
some_host         ansible_ssh_port=2222     ansible_ssh_user=manager
aws_host          ansible_ssh_private_key_file=/home/example/.ssh/aws.pem
freebsd_host      ansible_python_interpreter=/usr/local/bin/python
ruby_module_host  ansible_ruby_interpreter=/usr/bin/ruby.1.9.3
```

# 加密 Inventory 中的密码

假如现在有如下主机清单

```ini
[test]
hw-cloud-xngy-jump-server-linux-2 ansible_host=192.168.0.249 ansible_port=10022
[test:vars]
ansible_user=desistdaydream
ansible_password={{ test_password }}
ansible_become=yes
ansible_become_password={{ test_become_password }}
```

其中的密码，是通过变量引用的，而这些变量所在的文件是可以加密的，加密后，即可保证操作便捷的同时保证安全性

我们可以在 Inventory/ 目录下创建一个 group_vars 目录，并在 group_vars 目录创建一个与 主机组名称同名的文件，效果如下：

```bash
../inventory/
├── group_vars
│   └── test
└── test
```

在 group_vars/test 中，内容应该如下(其中的变量替换成自己的密码)：

```bash
test_password: ${PASSWORD}
test_become_password: ${PASSWORD}
```

然后使用 `ansible-vault encrypt group_vars/test` 命令加密该文件，加密后的文件内容如下：

```bash
$ANSIBLE_VAULT;1.1;AES256
33393764306539363338646334396661323930396235303837663131366562303237666337643864
6132363363383364316561326263633564323134336466660a623732613966653036326433313666
36356665663962613630393063303361353839313839636332313332666264363265646331333965
3164393363643639650a306233376438386333343961313735666161396365663235343430666437
64306364646266363563333437323364356332393639323436396136343438383662653133323634
33613239316237353839313632383530303638393966363133383834363662353135306563323635
386564623764653966303265653136353165
```

这时候，我们执行 Playbooks 时，如果不指定解密所需的密码，将会提示如下报错

```bash
~]$ ansible-playbook -i ../inventory/ variables.yaml

PLAY [test] *******************************************************************************************************************************************************************
ERROR! Attempting to decrypt but no vault secrets found
```

只需要添加 `--ask-vault-pass` 参数并输入密码，Ansible 即可在运行中解密文件，并获取其中的变量值

```bash
~]$ ansible-playbook -i ../inventory/ variables.yaml --ask-vault-pass
Vault password:

PLAY [test] *******************************************************************************************************************************************************************

TASK [variables : test] *******************************************************************************************************************************************************
changed: [hw-cloud-xngy-jump-server-linux-2]

TASK [variables : debug] ******************************************************************************************************************************************************
ok: [hw-cloud-xngy-jump-server-linux-2] => {
    "msg": {
        "ansible_facts": {
            "discovered_interpreter_python": "/usr/bin/python3"
        },
        "changed": true,
        "cmd": [
            "whoami"
        ],
        "delta": "0:00:00.002581",
        "end": "2021-10-10 22:10:39.338994",
        "failed": false,
        "rc": 0,
        "start": "2021-10-10 22:10:39.336413",
        "stderr": "",
        "stderr_lines": [],
        "stdout": "root",
        "stdout_lines": [
            "root"
        ]
    }
}
```
