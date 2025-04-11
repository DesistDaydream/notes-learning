---
title: Filters Plugins
linkTitle: Filters Plugins
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，用户指南-目录-使用过滤器操作数据](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html)
> - [https://www.zsythink.net/archives/2862](https://www.zsythink.net/archives/2862)

在本博客中，ansible是一个系列文章，我们会尽量以通俗易懂的方式总结ansible的相关知识点。

ansible系列博文直达链接：ansible轻松入门系列

现在我有一个需求，我想要将获取到的变量的值中的所有字母都变成大写，如果想要在playbook中实现这个需求，我该怎么办呢？我可以借助一个叫做"过滤器"的东西，帮助我完成刚才的需求，"过滤器（filters）"可以帮助我们对数据进行处理，这样解释可能不够直观，不如这样，我们先来看一个过滤器的小例子，然后结合示例解释过滤器是个什么东西，示例如下：

```yaml
- hosts: test70
  remote_user: root
  gather_facts: no
  vars:
    testvar: 1a2b3c
  tasks:
  - debug:
      msg: "{{ testvar | upper }}"
```

如上例所示，testvar变量的值中包含三个小写字母，在使用debug模块输出这个变量的值时，我们使用了一个管道符，将testvar变量传递给了一个名为"upper"的东西，"upper"就是一个"过滤器"，执行上例playbook后你会发现，testvar中的所有小写字母都被变成了大写。

通过上述示例，你一定已经明白了，过滤器是一种能够帮助我们处理数据的工具，其实，ansible中的过滤器功能来自于jinja2模板引擎，我们可以借助jinja2的过滤器功能在ansible中对数据进行各种处理，而上例中的upper就是一种过滤器，这个过滤器的作用就是将小写字母变成大写，你一定已经发现了，当我们想要通过过滤器处理数据时，只需要将数据通过管道符传递给对应的过滤器即可，当然，过滤器不只有upper，还有很多其他的过滤器，这些过滤器有些是jinja2内置的，有些是ansible特有的，如果这些过滤器都不能满足你的需求，jinja2也支持自定义过滤器。

这篇文章我们就来总结一些常用的过滤器的用法，在总结时，不会区分它是jinja2内置的过滤器，还是ansible所独有的，我们总结的目的是在ansible中使用这些过滤器，如果你想要了解jinja2中有哪些内置过滤器，可以参考jinja2的官网链接，如下

[http://jinja.pocoo.org/docs/2.10/templates/#builtin-filters](http://jinja.pocoo.org/docs/2.10/templates/#builtin-filters)

# 字符串操作有关的过滤器

```yaml
- hosts: test70
  remote_user: root
  vars:
    testvar: "abc123ABC 666"
    testvar1: "  abc  "
    testvar2: '123456789'
    testvar3: "1a2b,@#$%^&"
  tasks:
  - debug:
      #将字符串转换成纯大写
      msg: "{{ testvar | upper }}"
  - debug:
      #将字符串转换成纯小写
      msg: "{{ testvar | lower }}"
  - debug:
      #将字符串变成首字母大写,之后所有字母纯小写
      msg: "{{ testvar | capitalize }}"
  - debug:
      #将字符串反转
      msg: "{{ testvar | reverse }}"
  - debug:
      #返回字符串的第一个字符
      msg: "{{ testvar | first }}"
  - debug:
      #返回字符串的最后一个字符
      msg: "{{ testvar | last }}"
  - debug:
      #将字符串开头和结尾的空格去除
      msg: "{{ testvar1 | trim }}"
  - debug:
      #将字符串放在中间，并且设置字符串的长度为30，字符串两边用空格补齐30位长
      msg: "{{ testvar1 | center(width=30) }}"
  - debug:
      #返回字符串长度,length与count等效,可以写为count
      msg: "{{ testvar2 | length }}"
  - debug:
      #将字符串转换成列表，每个字符作为一个元素
      msg: "{{ testvar3 | list }}"
  - debug:
      #将字符串转换成列表，每个字符作为一个元素，并且随机打乱顺序
      #shuffle的字面意思为洗牌
      msg: "{{ testvar3 | shuffle }}"
  - debug:
      #将字符串转换成列表，每个字符作为一个元素，并且随机打乱顺序
      # 在随机打乱顺序时，将ansible_date_time.epoch的值设置为随机种子
      #也可以使用其他值作为随机种子，ansible_date_time.epoch是facts信息
      msg: "{{ testvar3 | shuffle(seed=(ansible_date_time.epoch)) }}"
```

跟数字操作有关的过滤器，示例如下

```yaml
- hosts: test70
  remote_user: root
  vars:
    testvar4: -1
  tasks:
  - debug:
      #将对应的值转换成int类型
      # ansible中，字符串和整形不能直接计算，比如{{ 8+'8' }}会报错
      #所以，我们可以把一个值为数字的字符串转换成整形后再做计算
      msg: "{{ 8+('8'  int) }}"
  - debug:
      #将对应的值转换成int类型,如果无法转换,默认返回0
      # 使用int(default=6)或者int(6)时，如果无法转换则返回指定值6
      msg: "{{ 'a'  int(default=6) }}"
  - debug:
      #将对应的值转换成浮点型，如果无法转换，默认返回'0.0'
      msg: "{{ '8'  float }}"
  - debug:
      #当对应的值无法被转换成浮点型时，则返回指定值’8.8‘
      msg: "{{ 'a'  float(8.88) }}"
  - debug:
      #获取对应数值的绝对值
      msg: "{{ testvar4  abs }}"
  - debug:
      #四舍五入
      msg: "{{ 12.5  round }}"
  - debug:
      #取小数点后五位
      msg: "{{ 3.1415926  round(5) }}"
  - debug:
      #从0到100中随机返回一个随机数
      msg: "{{ 100  random }}"
  - debug:
      #从5到10中随机返回一个随机数
      msg: "{{ 10  random(start=5) }}"
  - debug:
      #从5到15中随机返回一个随机数,步长为3
      #步长为3的意思是返回的随机数只有可能是5、8、11、14中的一个
      msg: "{{ 15  random(start=5,step=3) }}"
  - debug:
      #从0到15中随机返回一个随机数,这个随机数是5的倍数
      msg: "{{ 15  random(step=5) }}"
  - debug:
      #从0到15中随机返回一个随机数，并将ansible_date_time.epoch的值设置为随机种子
      #也可以使用其他值作为随机种子，ansible_date_time.epoch是facts信息
      #seed参数从ansible2.3版本开始可用
      msg: "{{ 15  random(seed=(ansible_date_time.epoch)) }}" |
```

列表操作相关的过滤器，示例如下

```yaml
- hosts: test70
  remote_user: root
  vars:
    testvar7: [22,18,5,33,27,30]
    testvar8: [1,[7,2,[15,9]],3,5]
    testvar9: [1,'b',5]
    testvar10: [1,'A','b',['QQ','wechat'],'CdEf']
    testvar11: ['abc',1,3,'a',3,'1','abc']
    testvar12: ['abc',2,'a','b','a']
  tasks:
  - debug:
      #返回列表长度,length与count等效,可以写为count
      msg: "{{ testvar7  length }}"
  - debug:
      #返回列表中的第一个值
      msg: "{{ testvar7  first }}"
  - debug:
      #返回列表中的最后一个值
      msg: "{{ testvar7  last }}"
  - debug:
      #返回列表中最小的值
      msg: "{{ testvar7  min }}"
  - debug:
      #返回列表中最大的值
      msg: "{{ testvar7  max }}"
  - debug:
      #将列表升序排序输出
      msg: "{{ testvar7  sort }}"
  - debug:
      #将列表降序排序输出
      msg: "{{ testvar7  sort(reverse=true) }}"
  - debug:
      #返回纯数字非嵌套列表中所有数字的和
      msg: "{{ testvar7  sum }}"
  - debug:
      #如果列表中包含列表，那么使用flatten可以'拉平'嵌套的列表
      #2.5版本中可用,执行如下示例后查看效果
      msg: "{{ testvar8  flatten }}"
  - debug:
      #如果列表中嵌套了列表，那么将第1层的嵌套列表‘拉平’
      #2.5版本中可用,执行如下示例后查看效果
      msg: "{{ testvar8  flatten(levels=1) }}"
  - debug:
      #过滤器都是可以自由结合使用的，就好像linux命令中的管道符一样
      #如下，取出嵌套列表中的最大值
      msg: "{{ testvar8  max }}"
  - debug:
      #将列表中的元素合并成一个字符串
      msg: "{{ testvar9  join }}"
  - debug:
      #将列表中的元素合并成一个字符串,每个元素之间用指定的字符隔开
      msg: "{{ testvar9  join(' , ') }}"
  - debug:
      #从列表中随机返回一个元素
      #对列表使用random过滤器时，不能使用start和step参数
      msg: "{{ testvar9  random }}"
  - debug:
      #从列表中随机返回一个元素,并将ansible_date_time.epoch的值设置为随机种子
      #seed参数从ansible2.3版本开始可用
      msg: "{{ testvar9  random(seed=(ansible_date_time.epoch)) }}"
  - debug:
      #随机打乱顺序列表中元素的顺序
      #shuffle的字面意思为洗牌
      msg: "{{ testvar9  shuffle }}"
  - debug:
      #随机打乱顺序列表中元素的顺序
      # 在随机打乱顺序时，将ansible_date_time.epoch的值设置为随机种子
      #seed参数从ansible2.3版本开始可用
      msg: "{{ testvar9  shuffle(seed=(ansible_date_time.epoch)) }}"
  - debug:
      #将列表中的每个元素变成纯大写
      msg: "{{ testvar10  upper }}"
  - debug:
      #将列表中的每个元素变成纯小写
      msg: "{{ testvar10  lower }}"
  - debug:
      #去掉列表中重复的元素，重复的元素只留下一个
      msg: "{{ testvar11  unique }}"
  - debug:
      #将两个列表合并，重复的元素只留下一个
      #也就是求两个列表的并集
      msg: "{{ testvar11  union(testvar12) }}"
  - debug:
      #取出两个列表的交集，重复的元素只留下一个
      msg: "{{ testvar11  intersect(testvar12) }}"
  - debug:
      #取出存在于testvar11列表中,但是不存在于testvar12列表中的元素
      #去重后重复的元素只留下一个
      #换句话说就是:两个列表的交集在列表1中的补集
      msg: "{{ testvar11  difference(testvar12) }}"
  - debug:
      #取出两个列表中各自独有的元素,重复的元素只留下一个
      #即去除两个列表的交集，剩余的元素
      msg: "{{ testvar11  symmetric_difference(testvar12) }}" |
```

变量未定义时相关操作的过滤器，示例如下

```yaml
- hosts: test70
  remote_user: root
  gather_facts: no
  vars:
    testvar6: ''
  tasks:
  - debug:
      #如果变量没有定义，则临时返回一个指定的默认值
      #注：如果定义了变量，变量值为空字符串，则会输出空字符
      #default过滤器的别名是d
      msg: "{{ testvar5  default('zsythink') }}"
  - debug:
      #如果变量的值是一个空字符串或者变量没有定义，则临时返回一个指定的默认值
      msg: "{{ testvar6  default('zsythink',boolean=true) }}"
  - debug:
      #如果对应的变量未定义,则报出“Mandatory variable not defined.”错误，而不是报出默认错误
      msg: "{{ testvar5  mandatory }}" |
```

其实，说到上例中的default过滤器，还有一个很方便的用法，default过滤器不仅能在变量未定义时返回指定的值，还能够让模块的参数变得"可有可无"。

这样说不太容易理解，不如我们先来看一个工作场景，然后根据这个工作场景来描述所谓的"可有可无"，就容易理解多了，场景如下：

假设，我现在需要在目标主机上创建几个文件，这些文件大多数都不需要指定特定的权限，只有个别文件需要指定特定的权限，所以，在定义这些文件时，我将变量定义为了如下样子

```yaml
vars:
  paths:
    - path: /tmp/testfile
      mode: '0444'
    - path: /tmp/foo
    - path: /tmp/bar
```

如上所示，我一共定义了3个文件，只有第一个文件指定了权限，第二个文件和第三个文件没有指定任何权限，这样定义目的是，当这三个文件在目标主机中创建时，只有第一个文件按照指定的权限被创建，之后的两个文件都按照操作系统的默认权限进行创建，为了方便示例，我只定义了3个文件作为示例，但是在实际工作中，你获得列表中可能有几十个这样的文件需要被创建，这些文件中，有些文件需要特定的权限，有些不需要，所以，我们可能需要使用循环来处理这个问题，但是在使用循环时，我们会遇到另一个问题，问题就是，有的文件有mode属性，有的文件没有mode属性，那么，我们就需要对文件是否有mode属性进行判断，所以，你可能会编写一个类似如下结构的playbook

```yaml
- hosts: test70
  remote_user: root
  gather_facts: no
  vars:
    paths:
      - path: /tmp/test
        mode: '0444'
      - path: /tmp/foo
      - path: /tmp/bar
  tasks:
  - file: dest={{item.path}} state=touch mode={{item.mode}}
    with_items: "{{ paths }}"
    when: item.mode is defined
  - file: dest={{item.path}} state=touch
    with_items: "{{ paths }}"
    when: item.mode is undefined |
```

上例中，使用file模块在目标主机中创建文件，很好的解决我们的问题，但是上例中，我们一共循环了两遍，因为我们需要对文件是否有mode属性进行判断，然后根据判断结果调整file模块的参数设定，那么有没有更好的办法呢？当然有，这个办法就是我们刚才所说的"可有可无"，我们可以将上例playbook简化成如下模样：

```yaml
- hosts: test70
  remote_user: root
  gather_facts: no
  vars:
    paths:
      - path: /tmp/test
        mode: '0444'
      - path: /tmp/foo
      - path: /tmp/bar
  tasks:
  - file: dest={{item.path}} state=touch mode={{item.mode  default(omit)}}
    with_items: "{{ paths }}" |
```

上例中，我们并没有对文件是否有mode属性进行判断，而是直接调用了file模块的mode参数，将mode参数的值设定为了"{{item.mode | default(omit)}}"，这是什么意思呢？它的意思是，如果item有mode属性，就把file模块的mode参数的值设置为item的mode属性的值，如果item没有mode属性，file模块就直接省略mode参数，'omit'的字面意思就是"省略"，换成大白话说就是：[有就用，没有就不用，可以有，也可以没有]，所谓的"可有可无"就是这个意思，是不是很方便？我觉得聪明如你一定看懂了，快动手试试吧~

施主加油吧，这篇文章就总结到这里，希望能够对你有所帮助~掰掰~
