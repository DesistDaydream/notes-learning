---
title: Files 类模块
---

# 概述

> 参考：
> - [官方文档 2.9，用户指南-使用模块-模块索引-文件模块](https://docs.ansible.com/ansible/2.9/modules/list_of_files_modules.html)

Files 类别模块用来处理文件、文本

- [acl – Set and retrieve file ACL information](https://docs.ansible.com/ansible/2.9/modules/acl_module.html#acl-module)
- [archive – Creates a compressed archive of one or more files or trees](https://docs.ansible.com/ansible/2.9/modules/archive_module.html#archive-module)
- [assemble – Assemble configuration files from fragments](https://docs.ansible.com/ansible/2.9/modules/assemble_module.html#assemble-module)
- [blockinfile # 添加、更新、删除指定的多行文本。Insert/update/remove a text block surrounded by marker lines](https://docs.ansible.com/ansible/2.9/modules/blockinfile_module.html#blockinfile-module)
- [copy # 用于将文件从本地或远程设备上复制到远程设备上的某个位置。Copy files to remote locations](https://docs.ansible.com/ansible/2.9/modules/copy_module.html#copy-module)
- [fetch – Fetch files from remote nodes](https://docs.ansible.com/ansible/2.9/modules/fetch_module.html#fetch-module)
- [file # 管理文件和文件熟悉，用于创建文件、目录等。Manage files and file properties](https://docs.ansible.com/ansible/2.9/modules/file_module.html#file-module)
- [find – Return a list of files based on specific criteria](https://docs.ansible.com/ansible/2.9/modules/find_module.html#find-module)
- [ini_file – Tweak settings in INI files](https://docs.ansible.com/ansible/2.9/modules/ini_file_module.html#ini-file-module)
- [iso_extract – Extract files from an ISO image](https://docs.ansible.com/ansible/2.9/modules/iso_extract_module.html#iso-extract-module)
- [lineinfile # 与 sed 命令类似，修改指定文件中匹配到的行或添加行。Manage lines in text files](https://docs.ansible.com/ansible/2.9/modules/lineinfile_module.html#lineinfile-module)
- [patch – Apply patch files using the GNU patch tool](https://docs.ansible.com/ansible/2.9/modules/patch_module.html#patch-module)
- [read_csv – Read a CSV file](https://docs.ansible.com/ansible/2.9/modules/read_csv_module.html#read-csv-module)
- [replace – Replace all instances of a particular string in a file using a back-referenced regular expression](https://docs.ansible.com/ansible/2.9/modules/replace_module.html#replace-module)
- [stat # 获取文件或文件系统状态 Retrieve file or file system status](https://docs.ansible.com/ansible/2.9/modules/stat_module.html#stat-module)
- [synchronize – A wrapper around rsync to make common tasks in your playbooks quick and easy](https://docs.ansible.com/ansible/2.9/modules/synchronize_module.html#synchronize-module)
- [tempfile – Creates temporary files and directories](https://docs.ansible.com/ansible/2.9/modules/tempfile_module.html#tempfile-module)
- [template # 根据文件模板，在远程主机上生成新文件。Template a file out to a remote server](https://docs.ansible.com/ansible/2.9/modules/template_module.html#template-module)
- [unarchive # 解压缩一个归档文件。就是 tar 命。Unpacks an archive after (optionally) copying it from the local machine](https://docs.ansible.com/ansible/2.9/modules/unarchive_module.html#unarchive-module)
- [xattr – Manage user defined extended attributes](https://docs.ansible.com/ansible/2.9/modules/xattr_module.html#xattr-module)
- [xml – Manage bits and pieces of XML files or strings](https://docs.ansible.com/ansible/2.9/modules/xml_module.html#xml-module)

# [blockinfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/blockinfile_module.html) # 添加、更新、删除指定的多行文本

blockinfile 模块使用示例：

```yaml
- name:写入多行文本
    blockinfile:
    path: /etc/hosts #指定要添加文本的文件
    block: | # 注意要使用 | 符号，否则将没有换行。
      10.0.13.77 iptv-k8s-master-1.tjiptv.net
      10.0.13.82 iptv-k8s-master-2.tjiptv.net
```

添加结果

```shell
[root@lichenhao ~]# cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
# BEGIN ANSIBLE MANAGED BLOCK
10.0.13.77 iptv-k8s-master-1.tjiptv.net
10.0.13.82 iptv-k8s-master-2.tjiptv.net
# END ANSIBLE MANAGED BLOCK
```

Note：

- blockinfile 模块会将 block 字段下面的所有内容当做一个文本块来看，将这一块内容全部添加到指定文件中 block 位置
  - block 位置是由 # BEGIN ANSIBLE MANAGED BLOCK 与 # END ANSIBLE MANAGED BLOCK 两行夹在中间的所有位置。
- 如果再次执行该任务，则会将 block 下指定的文本块覆盖到目标文件 ANSIBLE 所表示的那几行，而不会添加到文件末尾

# [copy](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/copy_module.html) # 用于将文件拷贝到被管理设备上的某个位置

## 参数

- 必选参数
  - **src: <PATH>** # 待拷贝的源文件路径。默认从 Ansible 控制节点搜索路径，搜索逻辑可以被 remote_src 参数修改
- 可选参数
  - **remote_src: <BOOLEAN>** # 若开启 remote_src 参数，则 src 参数将会从被管理节点搜索待拷贝的源文件。`默认值:false`

## 应用示例

```bash
~]$ ansible all -m copy -a "src=/etc/hosts dest=/tmp/hosts"
~]$ ansible all -m file -a "dest=/tmp/hosts mode=755 owner=lichenhao group=lichenhao state=directory"
```

```yaml
- name: 拷贝文件
  ansible.builtin.copy:
    src: /etc/hosts
    dest: /tmp/hosts
- name: 创建目录
  copy:
    dest: /tmp/hosts
    mode: 0755
    owner: lichenhao
    group: lichenhao
    state: directory
```

# [file](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/file_module.html#file-module) # 用于创建文件、目录等

file 模块示例：

```yaml
- name: 创建名为k8s的目录
  file:
    path: /etc/ssl/k8s # 指定要创建的路径
    owner: k8s #可省。默认所属用户为root
    group: k8s #可省。默认所属组为k8s
    state: directory #指定要创建的类型为目录
```

Note: state 还可以使用 link 用来创建软链接

# 文本替换

lineinfile、replace、blockinfile 这三个模块可以实现文本编辑功能，类似于 sed。

- replace 模块可以对多行中的文本执行操作。
- lineinfile 模块对一行全部内容执行操作。
- blockinfile 模块可以在文件中插入、更新、删除一行

对于其他情况的文本编辑，则需要使用 copy、template 等模块直接生成文件，而不是编辑文件。

虽然都是行匹配，但是 replace 只会替换匹配到行的字符串，而 lineinfile 则会替换匹配到整行。

## [lineinfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/lineinfile_module.html) # 行替换

- 此模块确保文件中包含特定行，或使用向后引用的正则表达式替换现有行。
- 当我们只想更改文件中的单行时，这主要有用。如果匹配到多行，那么只会对最后一行进行操作

### 参数

必填参数：

- **path** # 要修改的文件

其他参数：

- **regexp** # 正则表达式，用以匹配需要修改的行

注意

- 如果 regexp 和 insertbefore 同时存在，则仅在找不到与 regexp 匹配的情况下才使用 insertbefore。不可与 backref 一起使用。
- insertafter EOF 与 insertbefore BOF 分别表示在文件末尾以及文件开头插入指定的行

### 应用示例

取消 UseDNS no 行前的 # 号

```yaml
- name: 修改指定行的内容
  lineinfile:
    dest: /etc/ssh/sshd_config # 指定要修改的文件
    regexp: "#(UseDNS\\s*no)" # 使用正则表达式进行内容匹配
    line: '\1' # 引用正则表达式中()匹配的内容
    backrefs: yes # 指定是否可以进行引用，如果不指定，则匹配到的行会变成\1而不是通过\1引用的内容。
```

在指定匹配到的最后一行之前或者之后添加行

```yaml
- name: 在指定行之前或者之后添加行
  lineinfile:
    dest: /usr/lib/systemd/system/docker.service
    regexp: "^ExecStartPost=" # (可省略)如果未匹配到regexp中的内容，则执行insertabefore。如果匹配到了，则不再匹配insertbefore，直接将匹配到的行替换成line指定的行。
    insertbefore: "^ExecReload=" # 在该字段指定的行之前添加line指定的内容。如果想再指定的行之后添加，则使用insertafter关键字
    line: "ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT" # 要添加的行
```

配置 Docker 启动参数

```yaml
- name: 配置docker启动参数
  block:
    - lineinfile:
        dest: /usr/lib/systemd/system/docker.service
        regexp: "^ExecStartPost="
        insertbefore: "^ExecReload="
        line: "ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT"
    - lineinfile:
        dest: /usr/lib/systemd/system/docker.service
        regexp: "(ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock)"
        line: '\1 {{docker.options}}'
        backrefs: yes
```

## [replace](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/replace_module.html) # 字符替换

### 参数

必填参数：

- **path** # 要修改的文件。
- **regexp** # 正则表达式，匹配到需要修改的字符串
- **replace** # 替换 regexp 参数匹配到的字符串。若不指定该参数，则会删除 regexp 参数匹配到的字符串。

其他参数：

-

### 应用示例

```yaml
- replace:
    path: "/etc/apt/apt.conf.d/10periodic"
    regexp: "1"
    replace: "0"
```

# [stat](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/stat_module.html#stat-module) # 获取文件或文件系统状态

stat 模块类似于 Linux 中的 `stat` 命令。常用来在改变文件之间，获取某些文件状态，并注册为变量，以便为后续任务进行判断。

## 应用示例

# [template](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/template_module.html) # 根据文件模板，在远程主机上生成新文件

template 模块使用 Jinja2 模板语言处理文件，并将渲染后的文件传送到目标主机上。

## 包含变量

- ansible_managed # (configurable via the defaults section of ansible.cfg) contains a string which can be used to describe the template name, host, modification time of the template file and the owner uid.
- template_host # contains the node name of the template’s machine.
- template_uid # is the numeric user id of the owner.
- template_path # is the path of the template.
- template_fullpath # is the absolute path of the template.
- template_destpath # is the path of the template on the remote system (added in 2.8).
- template_run_date # is the date that the template was rendered.

## 参数

必选参数

- **dest** # 渲染后的模板文件被发送到目标主机的位置
- **src** # Ansible 控制节点上 Jinja2 格式的模板文件路径。

Jinja 行为控制

- **lstrip_blocks **# 是否移除前导空白符和制表符。`默认值：no`
- **trim_blocks** # 是否移除换行符。`默认值：yes`

其他参数

- **backup: <BOOLEAN>** # 是否创建一个包含时间戳信息的备份文件。`默认值：no`

## 返回值

```json
    "msg": {
        "changed": true,
        "checksum": "d176c556a237d7d62f8e1a95ffcec7c06c1e6851",
        "dest": "/tmp/template_variables.conf",
        "diff": [],
        "failed": false,
        "gid": 0,
        "group": "root",
        "md5sum": "2f549414973ec5b547e40d0f49357ce5",
        "mode": "0644",
        "owner": "root",
        "size": 329,
        "src": "/home/lichenhao/.ansible/tmp/ansible-tmp-1634401131.4389606-251148042480663/source",
        "state": "file",
        "uid": 0
    }
```

## 应用示例

```yaml
ansible -m template -a 'src=/mytemplates/foo.j2 dest=/tmp/foo.conf lstrip_blocks=yes'
```

模板文件示例：

```shell
{
{% if docker.registryMirrors is defined %} #如果docker.registryMirrors变量存在，则执行最后一行之前的语句
  "registry-mirrors": [{% for mirror in docker.registryMirrors %} #输出 "registry-mirrors": 后执行for循环，将docker.registryMirrors变量的多个值逐一传递给mirror变量，直到docker.registryMirros变量里的值全部引用完成

    "{{ mirror}}"{% if not loop.last %},{% endif %} #输出 mirror 变量的值。如果循环没有结束，则输出一个逗号
  {%- endfor %} #结束for循环

  ],
{% endif %} #结束if结构
}
```

输出结果示例：

```json
{
  "registry-mirrors": [
    "https://ac1rmo5p.mirror.aliyuncs.com",
    "https://123.123.123"
  ]
}
```

更多 template 模块原理及应用，详见 [Playbook 章节中的 Templates](✏IT 学习笔记/🛠️9.运维/Ansible/Playbook/Templates%20 模板(Jinja2).md 模板(Jinja2).md)

# [unarchive](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/unarchive_module.html) # 解压缩一个归档文件。就是 tar 命令。

## 参数

## 应用示例
