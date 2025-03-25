---
title: C
linkTitle: C
weight: 1
---

# 概述

> 参考：
>
> - [ISO C 工作组官网](http://www.open-std.org/jtc1/sc22/wg14/)
> - <https://www.learn-c.org/>
> - [Wiki, C Programming Language](<https://en.wikipedia.org/wiki/C_(programming_language)>)
> - [网道，C](https://wangdoc.com/clang/)

# Hello World

代码：`hello_world.c`

```c
#include <stdio.h>

int main(void) {
  printf("Hello World\n");
  return 0;
}
```

编译

```bash
gcc hello_world.c
```

运行

```shell
$ ./a.out
Hello World
```
