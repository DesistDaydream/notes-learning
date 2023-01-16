---
title: fs(文件系统相关参数)
---

> 参考：
> - [官方文档,Linux 内核用户和管理员指南-/proc/sys 文档-/proc/sys/fs 文档](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/fs.html)

## file-max 与 file-nr

### fs.file-max = 52706963

max-file 表示系统级别的能够打开的文件描述符的数量。是对整个系统的限制，并不是针对用户的。

> ulimit -n 控制进程级别能够打开的文件句柄的数量。提供对 shell 及其启动的进程的可用文件句柄的控制。这是进程级别的。

当系统尝试分配比 file-max 指定的值更多的文件描述符时，通常我们会看到如下报错：`VFS: file-max limit <number> reached`

### fs.file-nr = INT

file-nr 中的三个值分别表示：

- 已分配的文件描述符
- 已分配但未使用的文件描述符
- 最大文件描述符

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ cat /proc/sys/fs/file-max
9223372036854775807
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ cat /proc/sys/fs/file-nr
2400	0	9223372036854775807
```

通常情况下 已分配但未使用的文件描述符 的值总是为 0，这并不是错误的，只是意味着 `已分配的文件描述符=正在使用的文件描述符`

由于某些历史原因，内核虽然可以动态分配文件描述符，但是却无法再次释放它们~~~

## 其他

### fs.may_detach_mounts = 1

未知

### fs.nr_open = INT64

单个进程可分配的最大文件描述符数量。`默认值：1024 * 1024`，即 1048576。
