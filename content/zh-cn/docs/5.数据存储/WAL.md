---
title: WAL
linkTitle: WAL
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Write-ahead logging](https://en.wikipedia.org/wiki/Write-ahead_logging)

在计算机科学中，**Write-ahead logging(预写日志记录，简称 WAL)** 是一系列用于在[数据库](/docs/5.数据存储/数据库/数据库.md)系统中提供原子性和持久性（ACID 属性中的两个）的技术。

在使用 WAL 的系统中，所有的修改都先被写入到日志中，然后再被应用到系统状态中。通常包含 redo 和 undo 两部分信息。为什么需要使用 WAL，然后包含 redo 和 undo 信息呢？举个例子，如果一个系统直接将变更应用到系统状态中，那么在机器断电重启之后系统需要知道操作是成功了，还是只有部分成功或者是失败了（为了恢复状态）。如果使用了 WAL，那么在重启之后系统可以通过比较日志和系统状态来决定是继续完成操作还是撤销操作。

`redo log` 称为重做日志，每当有操作时，在数据变更之前将操作写入 redo log，这样当发生断电之类的情况时系统可以在重启后继续操作。`undo log` 称为撤销日志，当一些变更执行到一半无法完成时，可以根据撤销日志恢复到变更之间的状态。

现代文件系统通常至少对文件系统元数据使用 WAL 的变体；这就是所谓的 [Journaling File System](/docs/1.操作系统/Kernel/Filesystem/磁盘文件系统/Journaling%20File%20System.md)。
