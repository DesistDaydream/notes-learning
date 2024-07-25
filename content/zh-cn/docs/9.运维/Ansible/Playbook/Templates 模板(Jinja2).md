---
title: Templates 模板(Jinja2)
weight: 4
---

# 概述

> 参考：
>
> - [官方文档,用户指南-传统目录-使用剧本-模板(Jinja2)](https://docs.ansible.com/ansible/latest/user_guide/playbooks_templating.html)
> - 朱双印博客,jinja2 模板
>   - https://www.zsythink.net/archives/2999
>   - https://www.zsythink.net/archives/3021
>   - https://www.zsythink.net/archives/3037
>   - https://www.zsythink.net/archives/3051
> - [骏马金龙，9. 如虎添翼的力量：解锁强大的 Jinja2 模板](https://www.junmajinlong.com/ansible/9_power_of_jinja2/)

Jinja2 的内容较多，但对于学习 Ansible 来说，只需要学习其中和 template 相关的一部分 (其它的都和开发有关或 Ansible 中用不上) 以及 Ansible 对 Jinja2 的扩展功能即可。

详见 Python 编程语言部分的 《[Jinja](/docs/2.编程/高级编程语言/Python/Jinja.md)》章节

尽管在编写 Playbook 时可以不用在意是否要用 Jinja2，但 Ansible 的运行离不开 Jinja2，当 Ansible 开始执行 playbook 或任务时，总是会先使用 Jinja2 去解析所有指令的值，然后再执行任务。另一方面，在编写任务的过程中也会经常用到 Jinja2 来实现一些需求。所以，Jinja2 可以重要到成为 Ansible 的命脉。

严格地说，playbook 中所有地方都使用了 Jinja2，包括几乎所有指令的值、template 模板文件、copy 模块的 content 指令的值、lookup 的 template 插件、等等。它们会先经过 Jinja2 渲染，然后再执行相关任务。

例如，下面的 playbook 中分别使用了三种 Jinja2 特殊符号。

```yaml
---
- hosts: localhost
  gather_facts: no
  tasks:
    - debug:
        msg: "hello world, {{inventory_hostname}}"
    - debug:
        msg: "hello world{# comment #}"
    - debug:
        msg: "{% if True %}hello world{% endif %}"
```

> 注：jinja2 原生的布尔值应当是小写的 true 和 false，但也支持首字母大写形式的 True 和 False。

执行结果：

```text
TASK [debug] ************************
ok: [localhost] => {
    "msg": "hello world, localhost"
}

TASK [debug] ************************
ok: [localhost] => {
    "msg": "hello world"
}

TASK [debug] ************************
ok: [localhost] => {
    "msg": "hello world"
}
```

再比如模板文件 a.conf.j2 中使用这三种特殊语法：

```yaml
{# Comment this line #}
variable value: {{inventory_hostname}}
{% if True %}
in if tag code: {{inventory_hostname}}
{% endif %}
```

对应的模板渲染任务：

```yaml
- template:
    src: a.conf.j2
    dest: /tmp/a.conf
```

执行后，将在 / tmp/a.conf 中生成如下内容：

```yaml
variable value: localhost
in if tag code: localhost
```

有些指令比较特殊，它们已经使用隐式的 {{}} 进行了预包围，例如 debug 模块的 var 参数、条件判断`when`指令，所以这时就不要手动使用 {{}} 再包围指令的值。例如：

```yaml
- debug:
    var: inventory_hostname
```

但有时候也确实是需要在 var 或 when 中的一部分使用 {{}} 来包围表示这是一个变量或是一个表达式，而非字符串的。例如：

```yaml
- debug:
    var: hostvars['{{php}}']
  vars:
    - php: 192.168.200.143
```

# Ansible 扩展的测试函数

模板引擎是多功能的，可以用在很多方面，所以 Jinja2 自身置的大多数功能都是通用功能。使用 Jinja2 的工具可能会对 Jinja2 进行功能扩展，比如 Flask 扩展了一些功能，Ansible 也对 Jinja2 扩展了一些功能。

Ansible 扩展的测试函数官方手册：<https://docs.ansible.com/ansible/latest/user_guide/playbooks_tests.html>。

### 测试字符串

Ansible 提供了三个正则测试函数：

- match()
- search()
- regex()

它们都返回布尔值，匹配成功时返回 true。

其中，match() 要求从字符串的首字符开始匹配成功。

例如：

```text
"hello123world" is match("\d+")    -> False
"hello123world" is match(".*\d+")  -> True
"hello123world" is search("\d+")   -> True
"hello123world" is regex("\d+")    -> True
```

### 版本号大小比较

Ansible 作为配置服务、程序的配置管理工具，经常需要比较版本号的大小是否符合要求。Ansible 提供了一个`version`测试函数可以用来测试版本号是否大于、小于、等于、不等于给定的版本号。

语法：

```text
version('VERSION',CMP)
```

其中 CMP 可以是如下几种：

```text
<, lt, <=, le, >, gt, >=, ge, ==, =, eq, !=, <>, ne
```

例如：

```text
{{ ansible_facts["distribution_version"] is version("7.5","<=") }}
```

判断操作系统版本号是否小于等于 7.5。

### 子集、父集测试

- `A is subset(B)`测试 A 是否是 B 的子集
- `A is superset(B)`测试 A 是否是 B 的父集

例如：

```yaml
- debug:
    msg: '{{[1,2,3] is subset([1,2,3,4])}}'
```

### 成员测试

Jinja2 自己有一个`in`操作符可以做成员测试，Ansible 另外还实现了一个 contains 测试函数，主要目的是为了结合 select、reject、selectattr 和 rejectattr 筛选器。

官方给了一个示例：

```yaml
vars:
  lacp_groups:
    - master: lacp0
      network: 10.65.100.0/24
      gateway: 10.65.100.1
      dns4:
        - 10.65.100.10
        - 10.65.100.11
      interfaces:
        - em1
        - em2

    - master: lacp1
      network: 10.65.120.0/24
      gateway: 10.65.120.1
      dns4:
        - 10.65.100.10
        - 10.65.100.11
      interfaces:
          - em3
          - em4

tasks:
  - debug:
      msg: "{{ (lacp_groups|selectattr('interfaces', 'contains', 'em1')|first).master }}"
```

此外，Ansible 还实现了`all`和`any`测试函数，`all()`测试表示当序列中所有元素都返回 true 时，all() 返回 true，`any()`测试表示当序列中只要有元素返回 true，any() 就返回 true。

仍然是官方给的示例：

```yaml
  mylist:
      - 1
      - "{{ 3 == 3 }}"
      - True
  myotherlist:
      - False
      - True
tasks:
  - debug:
      msg: "all are true!"
    when: mylist is all

  - debug:
      msg: "at least one is true"
    when: myotherlist is any
```

### 测试文件

Ansible 提供了测试文件的相关函数：

- is exists：是否存在
- is directory：是否是目录
- is file：是否是普通文件
- is link：是否是软链接
- is abs：是否是绝对路径
- is same_file(F)：是否和 F 是硬链接关系
- is mount：是否是挂载点

```yaml
- debug:
    msg: "path is a directory"
  when: mypath is directory

- debug:
    msg: "path is {{ (mypath is abs)|ternary('absolute','relative')}}"

- debug:
    msg: "path is the same file as path2"
  when: mypath is same_file(path2)

- debug:
    msg: "path is a mount"
  when: mypath is mount
```

### 测试任务的执行状态

每个任务的执行结果都有 4 种状态：成功、失败、changed、跳过。

Ansible 提供了相关的测试函数：

- succeeded、success
- failed、failure
- changed、change
- skipped、skip

```yaml
- shell: /usr/bin/foo
  register: result
  ignore_errors: True

- debug:
    msg: "it failed"
  when: result is failed

- debug:
    msg: "it changed"
  when: result is changed

- debug:
    msg: "it succeeded in Ansible >= 2.1"
  when: result is succeeded

- debug:
    msg: "it succeeded"
  when: result is success

- debug:
    msg: "it was skipped"
  when: result is skipped
```

# Ansible 扩展的 Filter

Ansible 扩展了非常多的 Filter，非常非常多，本来我只想介绍一部分。但是想到有些人不愿看英文，我还是将它们全都写出来，各位权当看中文手册好了。实际上它们也都非常容易，绝大多数筛选器用法几乎都不用动脑，一看便懂。

### 类型转换类筛选器

例如：

```text
{{"123"|int}}
{{"123"|float}}
{{123|string}}
{{range(1,6)|list}}
{{123|bool}}
```

注意，没有 dict 筛选器转换成字典类型。

### 获取当前时间点

Ansible 提供的 now() 可以获取当前时间点。

例如：

```yaml
- debug:
    msg: "{{now()}}"
```

得到结果：

```text
    ok: [localhost] => {
      "msg": "2020-01-25 00:27:17.563627"
    }
```

可以指定输出的格式化字符串，支持的格式化字符串参考 python 官方手册：<https://docs.python.org/3/library/datetime.html#strftime-strptime-behavior>。

例如：

```yaml
- debug:
    msg: '{{now().strftime("%Y-%m-%d %H:%M:%S.%f")}}'
```

### YAML、JSON 格式化

Ansible 提供了几个和 YAML、JSON 格式化相关的 Filter：

    to_yaml
    to_json
    to_nice_yaml
    to_nice_json

它们都可使用 indent 参数指定缩进的层次。

`to_yaml`和`to_json`适用于调试，`to_nice_yaml`和`to_nice_json`适用于用户查看。

例如：

    - debug:
        msg: '{{f1|to_nice_json(indent=2)}}'
      vars:
        f1:
          father: "Bob"
          mother: "Alice"
          Children:
            - Judy
            - Tedy

### 参数忽略

Ansible 提供了一个特殊变量 omit，可以用来忽略模块的参数效果。

官方手册给了一个非常有代表性的示例，如下：

```yaml
- name: touch files with an optional mode
  file:
    dest: "{{ item.path }}"
    state: touch
    mode: "{{ item.mode | default(omit) }}"
  loop:
    - path: /tmp/foo
    - path: /tmp/bar
    - path: /tmp/baz
      mode: "0444"
```

当所迭代的元素中不存在 mode 项，则使用默认值，默认值设置为特殊变量 omit，使得 file 模块的 mode 参数被忽略，相当于未书写该参数。只有给定了 mode 项时，mode 参数才生效。

### 列表元素连接

`join`可以将列表各个元素根据指定的连接符连接起来：

    {{ [1,2,3] | join("-") }}

### 列表压平

前面的文章曾介绍过 flatten 筛选器，它可以将嵌套列表压平。

例如：

```yaml
- debug:
    msg: "{{ [3, [4, 2] ] | flatten }}"
- debug:
    msg: "{{ [3, [4, [2]] ] | flatten(levels=1) }}"
```

### 并集、交集、差集

Ansible 提供了集合理论类的求值操作：

- unique：去重
- union：并集，即两个集合中所有元素
- intersect：交集，即两个集合中都存在的元素
- difference：差集，即返回只在第一个集合中，不在第二个集合中的元素
- symmetric_difference：对称差集，即返回两个集合中不重复的元素

  - name: return [1,2,3]
      debug:
        msg: "{{ [1,2,3,2,1] | unique }}"
  - name: return [1,2,3,4]
      debug:
        msg: "{{ [1,2,3] | union([2,3,4]) }}"
  - name: return [2,3]
      debug:
        msg: "{{ [1,2,3] | intersect([2,3,4]) }}"
  - name: return [1]
      debug:
        msg: "{{ [1,2,3] | difference([2,3,4]) }}"
  - name: return [1,4]
      debug:
        msg: "{{ [1,2,3] | symmetric_difference([2,3,4]) }}"

### dict 和 list 转换

- `dict2items()`：将字典转换为列表
- `items2dict()`：将列表转换为字典

对于`dict2items`，例如：

    - debug:
        msg: "{{ p | dict2items }}"
      vars:
        p:
          name: junmajinlong
          age: 28

得到：

    [
      {
        "key": "name",
        "value": "junmajinlong"
      },
      {
        "key": "age",
        "value": 28
      }
    ]

对于`items2dict`，例如：

    - debug:
        msg: "{{ p | items2dict }}"
      vars:
        p:
          - key: name
            value: junmajinlong
          - key: age
            value: 28

得到：

    {
      "age": 28,
      "name": "junmajinlong"
    }

默认情况下，`dict2items`和`items2dict`都使用`key`和`value`来转换，但它们都允许使用`key_name`和`value_name`自定义转换的名称。

例如：

```yaml
- debug:
    msg: "{{  files | dict2items(key_name='file', value_name='path') }}"
  vars:
    files:
      users: /etc/passwd
      groups: /etc/group
```

得到：

    [
      {
        "file": "users",
        "path": "/etc/passwd"
      },
      {
        "file": "groups",
        "path": "/etc/group"
      }
    ]

### zip 和 zip_longest

`zip`和`zip_longest`可以将多个列表的元素一一对应并组合起来。它们的区别是，zip 以短序列为主，`zip_longest`以最长序列为主，缺失的部分使用 null 填充。

例如：

```yaml
- name: return [[1,"a"], [2,"b"]]
  debug:
    msg: "{{ [1,2] | zip(['a','b']) | list }}"

- name: return [[1,"a"], [2,"b"]]
  debug:
    msg: "{{ [1,2] | zip(['a','b','c']) | list }}"

- name: return [[1,"a"], [2,"b"], [null, "c"]]
  debug:
    msg: "{{ [1,2] | zip_longest(['a','b','c']) | list }}"

- name: return [[1,"a","aa"], [2,"b","bb"]]
  debug:
    msg: "{{ [1,2] | zip(['a','b'], ['aa','bb']) | list }}"
```

在 Python 中经常会将 zip 的运算结果使用 dict() 构造成字典，Jinja2 中也可以：

    - name: !unsafe 'return {"a": 1, "b": 2}'
      debug:
        msg: "{{ dict(['a','b'] | zip([1,2])) }}"

### 子元素 subelements

subelements 筛选器在前一章节详细解释过，这里不再介绍。

### random 生成随机数

Jinja2 自身内置了一个 random 筛选器，Ansible 也有一个 random 筛选器，比 Jinja2 内置的定制性要更强一点。

    "{{ ['a','b','c'] | random }}"
    # => 'c'

    # 生成[0,60)的随机数
    "{{ 60 | random }} * * * * root /script/from/cron"
    # => '21 * * * * root /script/from/cron'

    # 生成[0,100)的随机数，步进为10
    {{ 101 | random(step=10) }}
    # => 70

    # 生成[1,100)的随机数，步进为10
    {{ 101 | random(1, 10) }}
    # => 31
    {{ 101 | random(start=1, step=10) }}
    # => 51

    # 指定随机数种子。
    # 下面指定为主机名，所以同主机生成的随机数相同，但不同主机的随机数不同
    "{{ 60 | random(seed=inventory_hostname) }} * * * * root /script/from/cron"

### shuffle 打乱顺序

 ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lmdoi2/1638717920341-716b62fb-5d66-43db-a26d-9a14afaa42fa.png)

### json_query

可查询 Json 格式的数据，`json_query`在 Ansible 中非常实用，是必学 Filter 之一。

Ansible 的`json_query`基于 jmespath，所以需要先安装 jmespath：

```shell
pip3 install jmespath
```

jmespath 的查询语法相关示例可参见其官方手册：

- 入门手册： <http://jmespath.org/tutorial.html>
- 示例：<http://jmespath.org/examples.html>

下面我列出 Ansible 中给出的示例。

例如，对于下面的数据结构：

    {
      "domain_definition": {
        "domain": {
          "cluster": [
            {"name": "cluster1"},
            {"name": "cluster2"}
          ],
          "server": [
            {
              "name": "server11",
              "cluster": "cluster1",
              "port": "8080"
            },
            {
              "name": "server12",
              "cluster": "cluster1",
              "port": "8090"
            },
            {
              "name": "server21",
              "cluster": "cluster2",
              "port": "9080"
            },
            {
              "name": "server22",
              "cluster": "cluster2",
              "port": "9090"
            }
          ],
          "library": [
            {
              "name": "lib1",
              "target": "cluster1"
            },
            {
              "name": "lib2",
              "target": "cluster2"
            }
          ]
        }
      }
    }

使用

    {{ domain_definition | json_query('domain.cluster[*].name') }}

可以获取到名称 cluster1 和 cluster2。

使用

    {{ domain_definition | json_query('domain.server[*].name') }}

可以获取到 server11、server12、server21 和 server22。

使用

    - debug:
        var: item
      loop: "{{ domain_definition | json_query(server_name_cluster1_query) }}"
      vars:
        server_name_cluster1_query: "domain.server[?cluster=='cluster1'].port"

可以迭代 8080 和 8090 两个端口。

注意上面使用了问号`?`，这表示后面的是一个表达式。

使用

    "{{domain_definition | json_query('domain.server[?cluster==`cluster2`].{name1: name, port1: port}')}}"

可得到：

    [
      {
        "name1": "server21",
        "port1": "9080"
      },
      {
        "name1": "server22",
        "port1": "9090"
      }
    ]

注意上面使用了反引号``而不是单双引号，因为单双引号都被使用过了，再使用就不方便，可读性也差。

### ip 地址筛选

Ansible 提供了非常丰富的功能来完成 IP 地址的筛选，用法非常多，绝大多数关于 IP、网络地址类的计算、查询需求都能解决。

使用它需要先安装 python 的 netaddr 包：

```shell
pip3 install netaddr
```

完整用法参考官方手册：<https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters_ipaddr.html#playbooks-filters-ipaddr>。

下面是我整理的一部分用法。

检测是否是合理的 IP 地址：

    {{ myvar | ipaddr }}

检测是否是合理的 ipv4 地址、ipv6 地址：

    {{ myvar | ipv4 }}
    {{ myvar | ipv6 }}

从列表中筛选出合理的 IP 地址：

    test_list = ['192.24.2.1', 'host.fqdn', '::1', '192.168.32.0/24', 'fe80::100/10', True, '']

    # {{ test_list | ipaddr }}
    ['192.24.2.1', '::1', '192.168.32.0/24', 'fe80::100/10']

    # {{ test_list | ipv4 }}
    ['192.24.2.1', '192.168.32.0/24']

    # {{ test_list | ipv6 }}
    ['::1', 'fe80::100/10']

获取 IP 地址部分：

    {{ '192.0.2.1/24' | ipaddr('address') }}
    {{ ipvar | ipv4('address') }}
    {{ ipvar | ipv6('address') }}

检测或找出公网 IP 和私网 IP：

    # {{ test_list | ipaddr('public') }}
    ['192.24.2.1']

    # {{ test_list | ipaddr('private') }}
    ['192.168.32.0/24', 'fe80::100/10']

### 正则表达式筛选器

Ansible 提供了几个正则类的 Filter，主要有：

- `regex_search`：普通正则匹配
- `regex_findall`：全局匹配
- `regex_replace`：正则替换

例如：

    {{ 'foobar' | regex_search('(foo)') }}

    # 匹配失败时返回空
    {{ 'ansible' | regex_search('(foobar)') }}

    # 多行模式、忽略大小写的匹配
    {{ 'foo\nBAR' | regex_search("^bar", multiline=True, ignorecase=True) }}

    # 全局匹配
    # 每次匹配到的内容将存放在一个列表中
    {{ 'DNS servers 8.8.8.8 and 8.8.4.4' | regex_findall('\\b(?:[0-9]{1,3}\\.){3}[0-9]{1,3}\\b') }}

    # 正则替换
    # "ansible" 替换为 "able"
    {{ 'ansible' | regex_replace('^a.*i(.*)$', 'a\\1') }}

    # "foobar" 替换为 "bar"
    {{ 'foobar' | regex_replace('^f.*o(.*)$', '\\1') }}

    # 使用命名捕获，"localhost:80" 替换为 "localhost, 80"
    {{ 'localhost:80' | regex_replace('^(?P<host>.+):(?P<port>\\d+)$', '\\g<host>, \\g<port>') }}

    # "localhost:80" 替换为 "localhost"
    {{ 'localhost:80' | regex_replace(':80') }}

### URL 处理筛选器

`urlsplit`筛选器可以从一个 URL 中提取 fragment、hostname、netloc、password、path、port、query、scheme、以及 username。如果不传递任何参数，则直接返回一个包含了所有字段的字典。

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('hostname') }}
    # => 'www.acme.com'

    {{ "http://user:password@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('netloc') }}
    # => 'user:passwd@www.acme.com:9000'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('username') }}
    # => 'user'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('password') }}
    # => 'passwd'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('path') }}
    # => '/dir/index.html'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('port') }}
    # => '9000'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('scheme') }}
    # => 'http'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('query') }}
    # => 'query=term'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit('fragment') }}
    # => 'fragment'

    {{ "http://user:passwd@www.acme.com:9000/dir/index.html?query=term#fragment" | urlsplit }}
    # =>
    #   {
    #       "fragment": "fragment",
    #       "hostname": "www.acme.com",
    #       "netloc": "user:password@www.acme.com:9000",
    #       "password": "passwd",
    #       "path": "/dir/index.html",
    #       "port": 9000,
    #       "query": "query=term",
    #       "scheme": "http",
    #       "username": "user"
    #   }

### 编写注释的筛选器

在模板渲染中，有可能需要在目标文件中生成一些注释信息。Ansible 提供了`comment`筛选器来完成该任务。

    {{ "Plain style (default)" | comment }}

会得到：

    #
    # Plain style (default)
    #

可以自定义注释语法：

    {{ "My Special Case" | comment(decoration="! ") }}

得到：

    !
    ! My Special Case
    !

extract 结合 map 使用时，可以根据索引 (列表索引或字典索引) 从列表或字典中提取对应元素的值。

    {{ [0,2] | map('extract', ['x','y','z']) | list }}
    {{ ['x','y'] | map('extract', {'x': 42, 'y': 31}) | list }}

得到：

    ['x', 'z']
    [42, 31]

### dict 合并

combine 筛选器可以将多个 dict 中同名 key 进行合并 (以覆盖的方式合并)。

    {{ {'a':1, 'b':2} | combine({'b':3}) }}

得到：

    {'a':1, 'b':3}

使用`recursive=True`参数，可以递归到嵌套 dict 中进行覆盖合并：

    {{ {'a':{'foo':1, 'bar':2}, 'b':2} | combine({'a':{'bar':3, 'baz':4}}, recursive=True) }}

将得到：

    {'a':{'foo':1, 'bar':3, 'baz':4}, 'b':2}

可以合并多个 dict，如下：

    {{ a | combine(b, c, d) }}

d 中同名 key 会覆盖 c，c 会覆盖 b，b 会覆盖 a。

### hash 值计算

计算字符串的 sha1：

    {{ 'test1' | hash('sha1') }}

计算字符串的 md5：

    {{ 'test1' | hash('md5') }}

计算字符串的 checksum，默认即`hash('sha1')`值：

    {{ 'test2' | checksum }}

计算 sha512 密码 (随机 salt):

    {{ 'passwd' | password_hash('sha512') }}

计算 sha256 密码 (指定 salt):

    {{ passwd' | password_hash('sha256', 'mysalt') }}

同一主机生成的密码相同：

    {{ 'passwd' | password_hash('sha512', 65534 | random(seed=inventory_hostname) | string) }}

根据字符串生成 UUID 值：

    {{ hostname | to_uuid }}

### base64 编解码筛选器

    {{ encoded_str | b64decode }}
    {{ decoded_str | b64encode }}

例如，Ansible 有一个`slurp`模块，它的作用类似于 fetch 模块，它可以从目标节点中读取指定文件的内容，然后以 base64 方式编码返回，所以要获取其原始内容，需要 base64 解码。

例如：

```yaml
- slurp:
    src: "/var/run/sshd.pid"
  register: sshd_pid
- debug:
    msg: "base64_pid: {{sshd_pid.content}}"
- debug:
    msg: "sshd_pid: {{sshd_pid.content|b64decode}}"
```

结果：

    TASK [slurp] ***************
    ok: [localhost]

    TASK [debug] ******************
    ok: [localhost] => {
        "msg": "base64_pid: MTE4OAo="
    }

    TASK [debug] *****************
    ok: [localhost] => {
        "msg": "base64_pid: 1188\n"
    }

### 文件名处理

- `basename`：获取字符串中的文件名部分
- `dirname`：获取字符串中目录名部分
- `expanduser`：扩展家目录，即将`~`替换为家目录
- `realpath`：获取软链接的原始路径
- `splitext`：扩展名分离

对于 splitext 筛选器，例如：

    {{"nginx.conf"|splitext}}
    #=> ("nginx",".conf")

    {{'/etc/my.cnf'|splitext}}
    #=> ("/etc/my",".cnf")

### 日期时间类处理

相对来说，Ansible 中处理日期时间的机会是比较少的。但是 Ansible 也提供了比较方便的处理日期时间的方式。

使用`strftime`将当前时间或给定时间 (只能以 epoch 数值指定) 按照给定的格式输出：

    # 将当前时间点以year-month-day hour:min:sec格式输出
    {{ '%Y-%m-%d %H:%M:%S' | strftime }}

    # 将指定的时间按照指定格式输出
    {{ '%Y-%m-%d' | strftime(0) }}          # => 1970-01-01
    {{ '%Y-%m-%d' | strftime(1441357287) }} # => 2015-09-04

使用`to_datetime`可以将日期时间字符串转换为 Python 日期时间对象，既然得到了对象，就可以进行时间比较、时间运算等操作。

    # 计算时间差(单位秒)
    # 默认解析的日期时间字符串格式为%Y-%m-%d %H:%M:%S，但可以自定义格式
    {{ (("2016-08-14 20:00:12" | to_datetime) - ("2015-12-25" | to_datetime("%Y-%m-%d"))).total_seconds() }}
    #=>20203212.0

    # 计算相差多少天。只考虑天数，不考虑时分秒等
    {{ (("2016-08-14 20:00:12" | to_datetime) - ("2015-12-25" | to_datetime('%Y-%m-%d'))).days  }}

### human_to_bytes 和 human_readable

`human_readable`将数值转换为人类可读的数据量大小单位：

```yaml
{{1|human_readable}}

{{1|human_readable(isbits=True)}}

{{10240|human_readable}}

{{102400000|human_readable}}

{{102400000|human_readable(unit="G")}}

{{102400000|human_readable(isbits=True, unit="G")}}
```

`human_to_bytes`将人类可读的单位转换为数值：

    {{'0'|human_to_bytes}}        #= 0
    {{'0.1'|human_to_bytes}}      #= 0
    {{'0.9'|human_to_bytes}}      #= 1
    {{'1'|human_to_bytes}}        #= 1
    {{'10.00 KB'|human_to_bytes}} #= 10240
    {{   '11 MB'|human_to_bytes}} #= 11534336
    {{  '1.1 GB'|human_to_bytes}} #= 1181116006
    {{'10.00 Kb'|human_to_bytes(isbits=True)}} #= 10240

### 其它筛选器

`quote`为字符串加引号，比如方便编写 shell 模块的命令：

```yaml
- shell: echo {{ string_value | quote }}
```

`ternary`根据 true、false 来决定返回哪个值：

```yaml

{{ (gender == "male") | ternary('Mr','Ms') }}


{{ enabled | ternary('no shutdown', 'shutdown', omit) }}
```

`product`生成笛卡尔积：

    {{['foo', 'bar'] | product(['com'])|list}}
    #=>[["foo","com"], ["bar","com"]]

# 使用示例

```yaml
- name: 将foo.j2文件输出到指定主机的/etc/file.con
  template:
    src: /mytemplates/foo.j2 # 指定源文件，是一个用jinja2语言写的文件
    dest: /etc/file.conf # 指定要生成的目的文件
    mode: 0744 #必须添加一个前导零，以便Ansible的YAML解析器知道它是一个八进制数（例如0644或01777）或将其引号（例如'644'或'1777'），以便Ansible接收字符串并可以从字符串进行自己的转换成数字。
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

## 完全自定义的 nginx 虚拟主机配置

在生产中，一个开发不太完善的系统可能时不时就要去 nginx 虚拟主机中添加一个 location 配置段落，如果有多个 nginx 节点要配置，无疑这是一件让人悲伤的事情。

值得庆幸，Ansible 通过 Jinja2 模板可以很容易地解决这个看上去复杂的事情。

首先提供相关的变量文件`vhost_vars.yml`，内容如下：

```yaml
servers:
  - server_name: www.abc.com
    listen: 80
    locations:
      - match_method: ""
        uri: "/"
        root: "/usr/share/nginx/html/abc/"
        index: "index.html index.htm"
        gzip_types:
          - css
          - js
          - plain

      - match_method: "="
        uri: "/blogs"
        root: "/usr/share/nginx/html/abc/blogs/"
        index: "index.html index.htm"

      - match_method: "~"
        uri: "\\.php$"
        fastcgi_pass: "127.0.0.1:9000"
        fastcgi_index: "index.php"
        fastcgi_param: "SCRIPT_FILENAME /usr/share/www/php$fastcgi_script_name"
        include: "fastcgi_params"

  - server_name: www.def.com
    listen: 8080
    locations:
      - match_method: ""
        uri: "/"
        root: "/usr/share/nginx/html/def/"
        index: "index.html index.htm"

      - match_method: "~"
        uri: "/imgs/.*\\.(png|jpg|jpeg|gif)$"
        root: "/usr/share/nginx/html/def/imgs"
        index: "index.html index.htm"
```

从上面提供的变量文件来看，应该能看出来它的目的是为了能够自动生成一个或多个 server 段，而且允许随意增删改每个 server 段中的 location 及其它指令。这样一来，编写 nginx 虚拟主机配置的任务就变成了编写这个变量文件。

需注意，每个 location 段有两个变量名`match_method`和`uri`，作用是生成 nginx location 配置项的前一部分，即`location METHOD URI {}`。除这两个变量名外，剩余的变量名都会直接当作 nginx 配置指令渲染到配置文件中，所以它们都需和 nginx 指令名相同，比如 index 变量名渲染后会得到 nginx 的 index 指令。

剩下的就是写一个 Jinja2 模板文件，模板中 Jinja2 语句块标签部分我没有使用缩进，这样比较容易控制格式。文件内容如下：

    {# 负责渲染每个指令 #}
    {% macro config(key,value) %}
    {% if (value is sequence) and (value is not string) and (value is not mapping) %}
    {# 如果指令是列表 #}
    {% for item in value -%}
    {# 如生成的结果是：gzip_types css js plain; #}
    {{ key ~ ' ' ~ item if loop.first else item}}{{' ' if not loop.last else ';'}}
    {%- endfor %}
    {% else %}
    {# 如果指令不是列表 #}
    {{key}} {{value}};
    {% endif %}
    {% endmacro %}

    {# 负责渲染location指令 #}
    {% macro location(d) %}
    location {{d.match_method}} {{d.uri}} {
    {% for item in d|dict2items if item.key != "match_method" and item.key != "uri" %}
        {{ config(item.key, item.value) }}
    {%- endfor %}
      }
    {% endmacro %}

    {% for server in servers %}
    server {
    {% for item in server|dict2items %}
    {# 非location指令部分 #}
    {% if item.key != "locations" %}
      {{ config(item.key,item.value) }}
    {%- else %}
    {# 各个location指令部分 #}
    {% for l in item.value|default([],true) %}
      {{ location(l) }}
    {% endfor %}
    {% endif %}
    {%- endfor %}
    }
    {% endfor %}

然后使用 template 模块去渲染即可：

```yaml
- hosts: localhost
  gather_facts: no
  vars_files:
    - vhost_vars.yml
  tasks:
    - template:
        src: "vhost.conf.j2"
        dest: /tmp/vhost.conf
```

渲染得到的结果：

    server {
      server_name www.abc.com;
      listen 80;
      location  / {
        root /usr/share/nginx/html/abc/;
        index index.html index.htm;
        gzip_types css js plain;  }

      location = /blogs {
        root /usr/share/nginx/html/abc/blogs/;
        index index.html index.htm;
      }

      location ~ \.php$ {
        fastcgi_pass 127.0.0.1:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME /usr/share/www/php$fastcgi_script_name;
        include fastcgi_params;
      }

    }
    server {
      server_name www.def.com;
      listen 8080;
      location  / {
        root /usr/share/nginx/html/def/;
        index index.html index.htm;
      }

      location ~ /imgs/.*\.(png|jpg|jpeg|gif)$ {
        root /usr/share/nginx/html/def/imgs;
        index index.html index.htm;
      }

    }
