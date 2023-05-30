---
title: Access Control(访问控制)
---

# 概述

> 参考：
> - [Wiki,DAC](https://en.wikipedia.org/wiki/Discretionary_access_control)

Linux 使用 **Discretionary Access Control(自主访问控制，简称 DAC)** 概念控制所有文件的基本权限。

Linux 中每个文件都具有三个拥有者：

- **user** # 文件的属主，拥有文件的**一个 Linux Account(账户)**
- **group** # 文件的属组，拥有文件的**一组 Linux Account(账户)**
- **other** # 文件的其他，拥有该文件的**其他 Linux Account(账户)**

上述三个角色，可以被赋予三个基本权限

- **read** # 读，简写为 r
- **write** # 写，简写为 w
- **execute** # 执行，简写为 x

我们使用 `ls -l` 命令查看文件，可以从第 1 列看到文件的类型与权限

```bash
~]# ls -lh
total 20K
lrwxrwxrwx.   1 root root    7 May 24  2019 bin -> usr/bin
dr-xr-xr-x.   5 root root 4.0K May 24  2019 boot
drwxr-xr-x   20 root root 3.1K May 14 09:38 dev
drwxr-xr-x.  82 root root 8.0K Jun 21 19:42 etc
......
```

第 1 列共 11 个字符，中间 9 个字符用以表示文件的基本权限，最后一个字符是文件的 ACL 与 SELinux 属性。格式与说明如下：

| 示例   | 文件类型 | 属主权限 |     |      | 属组权限 |     |      | 其他权限 |     |      | ACL 与 SELinux |
| ------ | -------- | -------- | --- | ---- | -------- | --- | ---- | -------- | --- | ---- | -------------- |
|        |          | 读       | 写  | 执行 | 读       | 写  | 执行 | 读       | 写  | 执行 |                |
| 示例 1 | l        | r        | w   | x    | r        | w   | x    | r        | w   | x    | .              |
| 示例 2 | d        | r        | -   | x    | r        | -   | x    | r        | -   | x    | .              |
| 示例 3 | d        | r        | w   | x    | r        | -   | x    | r        | -   | x    |                |

中间 9 个字符分为 3 个组，分别对应文件的 3 种拥有者：

- 第一组为文件的属主权限
- 第二组为文件的属组权限
- 第三组为文件的其他权限

每组的 3 个字符都符合如下规则：

- 第一个字符表示是否有“读取”权限，为 `r` 或 `-`
- 第二个字符表示是否有“写入”权限，为 `w` 或 `-`
- 第三个字符表示是否有“执行”权限，为 `x` 或 `-`

若为 `-` 符号时，表示没有对应的权限

简单示例

- **-rw-------** # 表明了文件的拥有者对文件有 读、写 的权限，但是没有运行的权限。也很好理解，因为这是一个普通文件，默认没有可执行的属性。记住：如果有 w 权限（写的权限），那么表明也有删除此文件的权限。
- **----r-----** # 表明文件所在的群组内的用户只可以读此文件，但不能写也不能执行。
- **-------r--** # 表示其他用户只可以读此文件，但不能写也不能执行。

# 权限管理工具

> 参考：
> 
> - [Manual(手册)，chmod](https://man7.org/linux/man-pages/man1/chmod.1.html)
> - [Manual(手册)，chown](https://man7.org/linux/man-pages/man1/chown.1.html)

Linux 所说的权限，就是用户和组的权限。这是最基本的权限。后面的文章中还会介绍高级权限。

## chmod # 修改文件的访问权限命令

Linux/Unix 的文件权限分为三级 : 文件拥有者、文件所属组、其他。利用 chmod 可以控制文件如何被他人所调用。

### Syntax(语法)

**chmod \[OPTIONS] MODE\[,MODE]... FILE...**
**chmod \[OPTIONS] OCTAL-MODE FILE...**
**chmod \[OPTIONS] --reference=RFILE FILE...**

MODE 格式如下 : `[ugoa][[+-=][rwxX]…][,…]`，其中

- **\[ugoa]** # u 表示该档案的拥有者，g 表示与该档案的拥有者属于同一个群体(group)者，o 表示其他以外的人，a 表示这三者皆是。
- **\[+-=]** # `+` 表示增加权限、`-` 表示取消权限、`=` 表示唯一设定权限。
- **\[rwxX]** # 表示可读取，w 表示可写入，x 表示可执行。

**OPTIONS**

- -c, --changes # like verbose but report only when a change is made
- **-f, --silent, --quiet** # 静默模式
- -**v, --verbose** # 诊断模式，显示完整的执行过程
- --no-preserve-root do not treat '/' specially (the default)
- --preserve-root # fail to operate recursively on '/'
- --reference=RFILE use RFILE's mode instead of MODE values
- **-R, --recursive** # 递归方式设置，即修改指定目录及其所有子目录和其内文件的权限。

### EXAMPLE

- 修改当前目录及子目录内所有目录类型文件，并将这些目录的权限改为 755
  - **find ./\* -type d -exec chmod 755 {} ;**
- 修改当前目录及子目录内所有普通类型文件，并将这些文件的权限改为 644
  - **find ./\* -type f -exec chmod 644 {} ;**
- 文件 file.txt 的所有者增加读和运行的权限。
  - chmod u+rx file.txt
- 文件 file.txt 的群组其他用户增加读的权限。
  - chmod g+r file.txt
- 文件 file.txt 的其他用户移除读的权限。
  - chmod o-r file.txt
- 文件 file.txt 的群组其他用户增加读的权限，其他用户移除读的权限。
  - chmod g+r o-r file.txt
- 文件 file.txt 的群组其他用户和其他用户均移除读的权限。
  - chmod go-r file.txt
- 文件 file.txt 的所有用户增加运行的权限。
  - chmod +x file.txt
- 文件 file.txt 的所有者分配读，写和执行的权限；群组其他用户分配读的权限，不能写或执行；其他用户没有任何权限。
  - chmod u=rwx,g=r,o=- file.txt
- 递归执行赋权，设置 newname 文件夹权限
  - chmod -R 700 /home/newname

## chown # 改变文件的所有者命令

### Syntax(语法)

**chown \[OPTION] \[OWNER]\[:\[GROUP]] FILE...**
**chown \[OPTION] --reference=RFILE FILE...**

-c, --changes like verbose but report only when a change is made

-f, --silent, --quiet suppress most error messages

-v, --verbose output a diagnostic for every file processed

      --dereference      affect the referent of each symbolic link (this is

                         the default), rather than the symbolic link itself

-h, --no-dereference affect symbolic links instead of any referenced file

                         (useful only on systems that can change the

                         ownership of a symlink)

      --from=当前所有者:当前所属组

                         只当每个文件的所有者和组符合选项所指定时才更改所

有者和组。其中一个可以省略，这时已省略的属性就不

需要符合原有的属性。

      --no-preserve-root  do not treat '/' specially (the default)

      --preserve-root    fail to operate recursively on '/'

      --reference=RFILE  use RFILE's owner and group rather than

                         specifying OWNER:GROUP values

-R, --recursive operate on files and directories recursively

The following options modify how a hierarchy is traversed when the -R

option is also specified. If more than one is specified, only the final

one takes effect.

-H if a command line argument is a symbolic link

                         to a directory, traverse it

-L traverse every symbolic link to a directory

                         encountered

-P do not traverse any symbolic links (default)

### EXAMPLE

- 改变文件的用户（用 ls -l 可以快速查看原用户和组），后接新的所有者的用户名，再接文件名：
  - chown newname file.txt
- chown 命令也可以改变文件的群组，如下：
  - chown newname:friends file.txt # 将 file.txt 文件的用户改为 newname，所属用户组修改为 friends
  - chown root /u # 将 /u 的属主更改为"root"
  - chown root:staff /u # 和上面类似，但同时也将其属组更改为"staff"
  - chown -hR root /u # 将 /u 及其子目录下所有文件的属主更改为"root"

## chattr 与 lsattr

chattr 工具可以修改文件属性，添加了 `i` 属性的文件将无法被编辑，即使是 root 用户也不行。

chattr # 改变 Linux 文件系统上的文件属性

EXAMPLE

- 为 /etc/passwd 文件添加 i 属性
  - chattr +i /etc/passwd
- 为 /etc/passwd 文件去除 i 属性
  - chattr -i /etc/passwd

lsattr # 查看 Linux 文件系统上的文件属性
