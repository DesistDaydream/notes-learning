---
title: C
linkTitle: C
weight: 1
---

# 概述

> 参考：
>
> - [ISO C 工作组官网](http://www.open-std.org/jtc1/sc22/wg14/)
> - [官网？](https://www.c-language.org/)
> - [标准文件？](https://www.iso.org/standard/82075.html)
> - [Wiki, C Programming Language](https://en.wikipedia.org/wiki/C_(programming_language))

TODO: 头文件的由来

- https://www.zhihu.com/question/280665935

# 学习资料

[菜鸟教程，C](https://www.runoob.com/cprogramming/c-tutorial.html)（快速上手尝试，简单直接）

[网道，C](https://wangdoc.com/clang/)

https://www.learn-c.org/

# Hello World

代码：`hello_world.c`

```c
#include <stdio.h>

int main(void) {
  printf("Hello World\n");
  return 0;
}
```

编译（将会生成 a.out 文件）

```bash
gcc hello_world.c
```

运行

```shell
$ ./a.out
Hello World
```

# C 语言的版本

不同的标准，甚至对 `main()` 函数的写法都有定义，这么多年下来，下面这些写法在某种情况下都是正确的 ^[[知乎，\[C语言\]你真的了解C语言吗之main函数(一)](https://zhuanlan.zhihu.com/p/55329391)]

```c
main()
main(void)
main(int argc,char *argv[])
main(int argc,char **argv)
void main()
void main(void)
void main(int argc,char *argv[])
void main(int argc,char **argv)
int main()
int main(void)
int main(int argc,char *argv[])
int main(int argc,char **argv)
```

**一、K & RC**

    1978 年

**二、ANSI C**

    1989 年

**三、C95**

    1995 年

**四、C99**

    1999 年

**五、C11**

    2011 年

**六、C17**

    2017 年

**七、C2x**

    202x 年

# C 范儿

## 环境变量

## 项目结构

- bin/ # 编译后的可执行文件
- build/ # 编译产生的临时文件（如 .o 文件）
- src/ # 源代码
- include/ # 头文件

## 命令规范

## 代码格式

## 编码风格

## 依赖管理

## 构建方式

