---
title: SNMP(传统监控标准)
---

# 概述

> 参考：
> - [Wiki,SNMP](https://en.wikipedia.org/wiki/Simple_Network_Management_Protocol)

**Simple Network Management Protocol(简单网络管理协议，简称 SNMP)**。想实现该协议，通常需要由两部分完成(监控端和被监控端)，是在两端的两个进程之间进行通信，该进程都需要占用一个 socket

- 监控端：通常称为 NMS 端，管理端
- 被监控端：通常称为 Agent 端，NMS 要去收集被监控端的数据的时候，可能收集到的是一些很敏感的数据(CPU 使用率，带宽占用率等，这些一般是不公开的)，所以需要在被监控节点上部署一个专门的程序，这个程序能在本地取得一定的管理权限，然后接受监控端发送的数据收集指令，代为在被监控节点本地完成数据收集，所以被称为 Agent 端，代理端

SNMP 的工作模式，使用 udp 协议发送报文

- 监控端主动发送请求到被监控端的 agent 去收集数据
- 被监控节点主动向监控端报告自己所采集的数据
- 当监控端发现被监控端发生异常时，可以发送一些控制指令，将被监控端修改一些参数

# 实现 **SNMP 的组件**

- **Management Information Base(管理信息库，简称 MIB)** # 用来定义所有监控端的 objects，其中包括 objects 的名称、OID、数据类型、描述(干什么用的)。MIB 也可以看作是 SNMP 的服务端与代理端的沟通桥梁，只有具有统一的格式，才能确定数据。
  - **Object(对象)：**这个对象可以是一个具体需要采集到的数据，比如 内存、CPU、磁盘、网络接口等等，也可以是一种抽象的集合，比如地区、硬件、系统、硬件、网络等等。上面说的所有事物，每一个都是一个 Object。所以，Object 可以包含另一个 Object，这也是人称常常将 MIB 称为**树状**的原因
    - **Object Identifier(对象标识符，简称 OID)** # 每一个 Object 都有一个 OID
    - 数据存取格式：即每个 object 除了 OID 用作标示以外，还有数据内容需要遵循一定个格式规范
- **Structure of Managerment Intormation(管理信息结构,简称 SMI)** # 是 ASN.1 的子集
- **SNMP 本身**：一般通过 Net-SNMP 中的工具实现。

所谓的 **MIB**，其实主要是通过文件记录的内容。与其说是用文件记录，不如说 MIB 就是使用 ASN.1(标准的接口描述语言) 编写的代码。ASN.1 语言同样有类似 import、 function 这类的东西。只不过，记录 MIB 文件的语言，又与 ASN.1 有一些细微的区别，我们暂且称为 **MIB 语言 **吧~

可以这么说，**MIB 就是一门描述 OID 的编程语言。**

> **Abstract Syntax Notation One** (**ASN.1**) 是一种标准的 interface description language(接口描述语言)，用于定义以跨平台方式序列化和反序列化的数据结构。 它广泛用于电信和计算机网络，尤其是在密码学中。

# SNMP 安全

<https://www.dnsstuff.com/snmp-community-string#what-is-an-snmp-community-string>

## Community String(团体字符串)

SNMP 的 **Community String(团体字符串)** 是一种访问存储在路由器或其他设备中的统计信息的方法。有时简称为社区字符串或 SNMP 字符串，它包含与 GET 请求一起提供的用户凭据（ID 或密码）。

# OID 的格式

> 参考：
> - [OID 参考数据库](http://oidref.com/)

**在解释 MIB 语法前，首先要明确以下 OID 的格式，因为 MIB 就是用来记录 OID 的，只有明确了 OID 的格式，才能更好的理解 MIB 的语法。**

OID 有两种格式：

- 字符串：`MIB::OBJECT.INDEX`
  - MIB # 表示当前 Object 所属的 MIB 库
  - OBJECT # 表示当前 Object 的名称
  - INDEX # 表示当前 Object 的索引
- 数字：`.x.x.x.x.x.`
  - 每一个 Object 都对应一个数字，而 Object 总是属于某一个 MIB，所以，可以将 字符串 转换为 数字。

其中数字格式的 OID 是一种树状的结构，国际统一标准，效果如图所示，从 root 开始,每一层使用一个数字来表示

> 一般常见的都是从 `.1.3.6.1` 开始

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqmpt9/1616067085191-648071a5-d977-424f-b97d-e40625eb536e.jpeg)
这两种方式可以使用 snmptranslate 命令进行转换，转换成字符串后，人类可以通过英文了解到大概意思，比如下图，表示的是该设备内存的大小。这是 snmpwalk 命令获取内存这个 OID 的当前的值(下图中=后面的内容是该 Object 的数据类型以及值，数据类型与值以冒号分隔)

SNMPv2-MIB::sysDescr.0 这个 Object 用来输出系统描述信息。显示一些基本的系统信息

    [root@exporter ~]# snmpwalk -v 2c -c public 172.19.42.243 | more
    SNMPv2-MIB::sysDescr.0 = STRING: Linux HDM210235A3KHH209000234 3.14.17-ami #1 Thu Sep 10 10:55:48 CST 2020 armv6l
    [root@exporter ~]# snmptranslate -On iso
    .1
    [root@exporter ~]# snmptranslate -On SNMPv2-SMI::dod
    .1.3.6
    [root@exporter ~]# snmptranslate -On SNMPv2-MIB::sysDescr.0
    .1.3.6.1.2.1.1.1.0

**Object 中的 INDEX 说明：**
INDEX 是该 obejct 中项目的索引，默认值为 0。一个 object 里面可能包含多个项目，比如一个块磁盘里面可能包含多个分区，效果如下图，当一个 object 里面有多个项目时，INDEX 则不再是 0，而是一个随机数。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqmpt9/1616067085208-f6f5abc2-3713-4baf-8852-f1db9a9c7905.jpeg)
想要获取某个 object 的值，可以直接使用数字，也可以直接使用字符串
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqmpt9/1616067085202-a29526e8-85c9-4e45-91ea-d6c8baf46310.jpeg)

# MIB 语法

> 参考：[CSDN](https://blog.csdn.net/shanzhizi/article/details/15340305?spm=a2c4e.10696291.0.0.149419a4qn7Vvh)、[CSDN2](https://blog.csdn.net/fuhanghang/article/details/104656348)

注释的写法：在每行开头写两个 - 的行

MIB 编写完成后，呈现的是一种树状的结构。MIB 实际上就是很多 Objects 的集合。Objects 的结构就像树状一样。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqmpt9/1616067085173-e9ba88c6-a0c7-417e-a19e-d3eee254cc58.jpeg)

## MIB HelloWorld

MIB 文件总是以 `XXXX DEFINITIONS ::= BEGIN` 开始，最后一行以 `END` 结束，在 BEGIN 和 END 之间，就是用来定义 Object 的代码块。其中 XXXX 为库的名称。比如：

    # 这是一个名为 SNMPv2-SMI 的库。这给地方相当于 func main() { 而最后的 END 就是 }
    SNMPv2-SMI DEFINITIONS ::= BEGIN
    # 从这里开始，就是具体的代码逻辑
    # 定义了一个名为 org 的对象，OID 号为 3，属于 iso 这个对象。而 iso 这个对象又属于 iso
    org            OBJECT IDENTIFIER ::= { iso 3 }  --  "iso" = 1
    # 定义了一个名为 dod 的对象，OID 号为 6，属于 org 这个对象。
    dod            OBJECT IDENTIFIER ::= { org 6 }
    # 定义了一个名为 internet 的对象，OID 号为 1，属于 dod 这个对象
    internet       OBJECT IDENTIFIER ::= { dod 1 }

    directory      OBJECT IDENTIFIER ::= { internet 1 }

    mgmt           OBJECT IDENTIFIER ::= { internet 2 }
    mib-2          OBJECT IDENTIFIER ::= { mgmt 1 }
    ......略
    # 结束
    END

在代码块中，通常定义两种东西：

- **OBJECT # 对象。**就是上文描述的一个可以采集到的监控数据对象
- **MACRO # 宏。**一种类似 function 的代码块，可以实现一个完成的功能，并让其他 MIB 调用。

在 MIB 中除了可以定义 Object，还可以定义一些类似于编程语言的 function 的东西，在 MIB 语言里称为 **MACRO(宏指令)**

> 这里提个题外话，本人玩魔兽，经常就会编写一些 **宏**，魔兽中的宏，其实跟变成语言中的 function 是类似的概念，都是通过一组代码，实现一个行为。

## MIB 语言关键字

### DEFINITIONS # 一般只出现在文件开头，用来定义一个该 MIB 的名称。与文件最后一行的 END 关键字组合，构成了一个完成 MIB 代码

### OBJECT IDENTIFIER # 对象标识符，就是用来定义一个对象

    org            OBJECT IDENTIFIER ::= { iso 3 } --  "iso" = 1

这是定义 Object 的标准格式

- 关键字的左边 # 是对象的名称
- 关键字的右边 # 使用 `::={ }` 这种格式，其中 `{ }` 中的内容是 `{该对象所属对象的OID名 该对象的OID号}`

也就是说 org 这个对象如果用数字格式表示的话，就是 `.1.3`，如果用字符串表示的话需要一个前提，就是指定 MIB，因为 org 这个对象可以定义在不同的 MIB 中，只不过大家懒得重复造轮子，大部分情况都会直接导入 SNMPv2-SMI 这个 MIB，以便可以直接使用其中已经定义好的的各种 Object。但是 iso 是不用指定 MIB 的，因为 iso 是基础，是根下面的一级，就算属于，也只会属于 根 MIB

比如咱用 snmptranslate 转换以下

    [root@exporter ~]# snmptranslate -On iso
    .1
    [root@exporter ~]# snmptranslate -On org
    org: Unknown Object Identifier (Sub-id not found: (top) -> org)
    [root@exporter ~]# snmptranslate -On SNMPv2-SMI::org
    .1.3

除了这种最基本的定义 Object 的格式，**还可以通过自定义的 MACRO 来定义 Object**，这种方式更为复杂，但是可以更详细得描述一个 Object 的数据格式，以及代表什么监控项。详见 OBJECT-TYPE 宏。这种定义 Object 的方式，常用来定义一个具体的待采集的具体的监控项，比如 内存大小、内存使用率、磁盘大小 等等

### MACRO # 宏，用来定义一个宏指令

    MODULE-IDENTITY MACRO ::=
    BEGIN
        TYPE NOTATION ::=
                      "LAST-UPDATED" value(Update ExtUTCTime)
                      "ORGANIZATION" Text
                      "CONTACT-INFO" Text
                      "DESCRIPTION" Text
                      RevisionPart
        VALUE NOTATION ::=
                      value(VALUE OBJECT IDENTIFIER)
        RevisionPart ::=
                      Revisions
                    | empty
        Revisions ::=
                      Revision
                    | Revisions Revision
        Revision ::=
                      "REVISION" value(Update ExtUTCTime)
                      "DESCRIPTION" Text
        -- a character string as defined in section 3.1.1
        Text ::= value(IA5String)
    END

MACRO 也是以 BEGIN 开头，END 结尾，上面的例子定义个一个名为 MODULE-IDENITY 的宏指令，可以让其他 MIB 导入后直接使用。

### IMPORTS # 从其他 MIB 中导入 Object(对象) 或 MARCRO(宏指令)

    IMPORTS
        MODULE-IDENTITY, OBJECT-TYPE, NOTIFICATION-TYPE,
        TimeTicks, Counter32, snmpModules, mib-2
            FROM SNMPv2-SMI
        DisplayString, TestAndIncr, TimeStamp
            FROM SNMPv2-TC
        MODULE-COMPLIANCE, OBJECT-GROUP, NOTIFICATION-GROUP
            FROM SNMPv2-CONF;

SNMPv2-MIB 这个库导入了多个对象和宏，来自于三个库。宏一般都是全大写的字符串：

- 从 SNMPv2-SMI 库导入了 MODULE-IDENTITY, OBJECT-TYPE, NOTIFICATION-TYPE, TimeTicks, Counter32, snmpModules, mib-2
- 从 FROM SNMPv2-TC 库导入了 DisplayString, TestAndIncr, TimeStamp
- 从 SNMPv2-CONF 库 导入了 MODULE-COMPLIANCE, OBJECT-GROUP, NOTIFICATION-GROUP

### 其他

## 总结

经过多年的发展，现在已经发展出了很多非常常见的 MIB 和 MACRO。很多 snmp 相关的工具，默认都已经导入了这些基本的 MIB。

所以当各大公司，想要根据自己产品，再定义新的 MIB 时，绝大多数情况，都会直接引用这些常见的 MIB，基于此，再定义自己的 MIB

### 常见的 MIB

**SNMPv2-SMI**

非常基本的 MIB，从 根 开始定义每一个 Object，并定义了以下常见的 MACRO。很多 MIB 中，都会导入这个 MIB。

就像这个 库的名字一样，它代表了 v2 版本的 SNMP 应该具有的一般 OBJECT 和 MACRO。

    SNMPv2-SMI DEFINITIONS ::= BEGIN
    org            OBJECT IDENTIFIER ::= { iso 3 }  --  "iso" = 1
    dod            OBJECT IDENTIFIER ::= { org 6 }
    internet       OBJECT IDENTIFIER ::= { dod 1 }
    directory      OBJECT IDENTIFIER ::= { internet 1 }
    ......略
    MODULE-IDENTITY MACRO ::=
    BEGIN
    ......略
    END
    ......略
    END

**SNMPv2-MIB**

**SNMPv2-TC**

**IF-MIB**

### 常见的 **MACRO(宏指令)**

经过了这么多年的发展，现在有很多常用的 MACRO 作为默认自带的。比如其中 MODULE-IDENTITY 这个宏，就是对 Object 的一个抽象，其实也可以当作 Object 的一种。

有一个名为** SNMPv2-SMI** 的库，可以当作一个基本库，每次大家定义新库时，总会导入这个库，**这些常见的 MACRO 很多都是在 SNMPv2-SMI 这个库中定义的**

**MODULE-IDENTITY**

很多 Object 的集合

**OBJECT IDENTIFIER**

**TEXTUAL-CONVENTION**

定义了对标准数据类型的进行扩展的语法。很多 MIB 定义中都会先定义一些**基于标准类型的扩展类型**，如：

    KBytes ::= TEXTUAL-CONVENTION（文本约定）
        STATUS current
        DESCRIPTION
            "Storage size, expressed in units of 1024 bytes."
        SYNTAX Integer32 (0..2147483647)

**OBJECT-TYPE**

用来定义一个 Object，

    sysDescr OBJECT-TYPE
        SYNTAX      DisplayString (SIZE (0..255))
        MAX-ACCESS  read-only
        STATUS      current
        DESCRIPTION
                "A textual description of the entity.  This value should
                include the full name and version identification of
                the system's hardware type, software operating-system,
                and networking software."
        ::= { system 1 }

需要注意，这里面定义 Object 已经不是使用 **OBJECT IDENTIFIER **这个关键字了，而是通过 MACRO 来定义的，这个宏就是 OBJECT-TYPE，

**OBJECT-GROUP**

**MODULE-COMPLIANCE**

## MIB 文件简单示例：

    # SNMPv2-MIB 是该 MIB 库的名字，也是在调用该 MIB 时所用的名字。
    SNMPv2-MIB DEFINITIONS ::= BEGIN

    # 类似于编程中导入某某包，这里就是导入一些 MIB 信息。比如Object 类型、
    IMPORTS
        MODULE-IDENTITY, OBJECT-TYPE, NOTIFICATION-TYPE,
        TimeTicks, Counter32, snmpModules, mib-2
            FROM SNMPv2-SMI
    ......略

    # 该 MIB 库的一些基本信息，介绍、维护者的邮箱名字等等
    # snmpMIB 是最重要的标识之一，标识该模块的标识。一般情况，一个文件中只有一个 MODULE-IDENTITY
    # 且该文件中后面所有的 Object 都属于 snmpMIB 这个 MODULE-IDENTITY
    # 这些模块可以在 IMPORTS 中被 MIB 导入。
    snmpMIB MODULE-IDENTITY
        DESCRIPTION
                "The MIB module for SNMP entities.

                 Copyright (C) The Internet Society (2002). This
                 version of this MIB module is part of RFC 3418;
                 see the RFC itself for full legal notices.
                "
        REVISION      "200210160000Z"

    # snmpMIBObjects 是一个 Objcet 的标识符
    # 属于 snmpMIB 这个模块
    snmpMIBObjects OBJECT IDENTIFIER ::= { snmpMIB 1 }


    # system 是一个 Object 的标识符。
    # 属于 mib-2 这个模块
    system   OBJECT IDENTIFIER ::= { mib-2 1 }

    # sysDescr 这个 Object 的信息。
    # ::= { system 1 } 表示 sysDescr 属于 system 这个 Object
    sysDescr OBJECT-TYPE
        SYNTAX      DisplayString (SIZE (0..255))
        MAX-ACCESS  read-only
        STATUS      current
        DESCRIPTION
                "A textual description of the entity.  This value should
                include the full name and version identification of
                the system's hardware type, software operating-system,
                and networking software."
        ::= { system 1 }

    # sysObjectID 这个 Object 的信息。
    # ::= { system 2 } 表示 sysObjectID 属于 system 这个 Object
    sysObjectID OBJECT-TYPE
        SYNTAX      OBJECT IDENTIFIER
        MAX-ACCESS  read-only
        STATUS      current
        DESCRIPTION
                "The vendor's authoritative identification of the
                network management subsystem contained in the entity.
                This value is allocated within the SMI enterprises
                subtree (1.3.6.1.4.1) and provides an easy and
                unambiguous means for determining `what kind of box' is
                being managed.  For example, if vendor `Flintstones,
                Inc.' was assigned the subtree 1.3.6.1.4.1.424242,
                it could assign the identifier 1.3.6.1.4.1.424242.1.1
                to its `Fred Router'."
        ::= { system 2 }

    END

如果仔细观察，可以看到一个树状的关系网：

- 在 `IMPORTS` 字段中导入了 mib-2
- 在 `system   OBJECT IDENTIFIER ::= { mib-2 1 }` 中可以看导 system 这个 Object 属于 mib-2
- 在 `sysDescr OBJECT-TYPE ...略... ::= { system 1 }` 中可以看到 sysDescr 这个 Object 属于 system

# 一些官方的 MIB

- Cisco: ftp://ftp.cisco.com/pub/mibs/v2/v2.tar.gz
- APC: <https://download.schneider-electric.com/files?p_File_Name=powernet432.mib>
- Servertech: ftp://ftp.servertech.com/Pub/SNMP/sentry3/Sentry3.mib
- Palo Alto PanOS 7.0 enterprise MIBs: <https://www.paloaltonetworks.com/content/dam/pan/en_US/assets/zip/technical-documentation/snmp-mib-modules/PAN-MIB-MODULES-7.0.zip>
- Arista Networks: <https://www.arista.com/assets/data/docs/MIBS/ARISTA-ENTITY-SENSOR-MIB.txt> <https://www.arista.com/assets/data/docs/MIBS/ARISTA-SW-IP-FORWARDING-MIB.txt> <https://www.arista.com/assets/data/docs/MIBS/ARISTA-SMI-MIB.txt>
- Synology: <https://global.download.synology.com/download/Document/Software/DeveloperGuide/Firmware/DSM/All/enu/Synology_MIB_File.zip>
- MikroTik: <http://download2.mikrotik.com/Mikrotik.mib>
- UCD-SNMP-MIB (Net-SNMP): <http://www.net-snmp.org/docs/mibs/UCD-SNMP-MIB.txt>
- Ubiquiti Networks: <http://dl.ubnt-ut.com/snmp/UBNT-MIB> <http://dl.ubnt-ut.com/snmp/UBNT-UniFi-MIB> <https://dl.ubnt.com/firmwares/airos-ubnt-mib/ubnt-mib.zip>

<https://github.com/librenms/librenms/tree/master/mibs> 也是不错的 MIB 来源，这里收录了很多的 MIB 文件

推荐使用 <http://oidref.com> 浏览 MIBs.
