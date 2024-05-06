---
title: Security software
linkTitle: Security software
date: 2024-05-06T13:08
weight: 1
---


# 概述

> 参考：
>
> -

[Snort](/docs/7.信息安全/Security%20software/Snort.md)

[Yara](#yara)

# Yara

> 参考：
>
> - [GitHub 项目，virustotal/yara](https://github.com/virustotal/yara)
> - [官网](https://virustotal.github.io/yara/)
> - [Wiki，YARA](https://en.wikipedia.org/wiki/YARA)

YARA 是一款旨在（但不限于）帮助恶意软件研究人员识别和分类恶意软件样本的工具。使用 YARA，您可以根据 **文本** 或 **二进制** 模式创建恶意软件系列（或您想要描述的任何内容）的描述。每个描述（也称为规则）由一组字符串和一个决定其逻辑的布尔表达式组成。让我们看一个例子：

```json
rule silent_banker : banker
{
    meta:
        description = "This is just an example"
        threat_level = 3
        in_the_wild = true

    strings:
        $a = {6A 40 68 00 30 00 00 6A 14 8D 91}
        $b = {8D 4D B0 2B C1 83 C0 27 99 6A 4E 59 F7 F9}
        $c = "UVODFRYSIHLNWPEJXQZAKCBGMT"

    condition:
        $a or $b or $c
}
```

# 安全系统提供商

数美科技 https://www.ishumei.com/ # 在线业务风控解决方案提供商
