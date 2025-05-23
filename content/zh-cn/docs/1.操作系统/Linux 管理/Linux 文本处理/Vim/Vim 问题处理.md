---
title: Vim 问题处理
linkTitle: Vim 问题处理
weight: 20
---

# 概述

> 参考：
>
> -

# Vim 黄色阴影处理

经常会出现字符上会出现黄色的阴影部分，虽然不会影响到自己的编辑，但是看着就是不习惯，解决这个问题的方式：

- 在 vim 编辑器中执行 :nohlsearch 即可
- 如果希望以后都不在出现这个情况，可以在 /etc/vim/vimrc 文件中添加 set nohlsearch，使之永久生效

# Vim中复制粘贴缩进错乱问题的解决方案

不知道大家是否会有这种困扰，例如在 Android Studio 有一段缩进优美的代码实现，例如：

```bash
public void sayHello() {
    String msg = "Hello Vim Paste Mode";
    System.out.println(msg);
}
```

当你把这段缩进优美的代码直接 ctrl+c，ctrl+v 到 Vim 的时候，就会出现如下恶心的情况：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dssefy/1616166238841-3e8c3b7d-589a-4ded-9c81-77f47e2ada04.png)

可以看到，这种直接粘贴的方式会导致代码丢失和缩进错乱等情况。

## **解决方案**

vim 进入 paste 模式，命令如下：

```
:set paste
```

进入 paste 模式之后，再按 i 进入插入模式，进行复制、粘贴就很正常了。

命令模式下，输入下面的命令以解除 paste 模式。

```
:set nopaste
```

paste 模式主要帮我们做了如下事情：

- textwidth 设置为 0
- wrapmargin 设置为 0
- set noai
- set nosi
- softtabstop 设置为 0
- revins 重置
- ruler 重置
- showmatch 重置
- formatoptions 使用空值
