---
title: Go 常见问题
---

# goroutines 与 os.Chdir

> 参考：
>
> - [GitHub 项目,golang/go-issue-27658](https://github.com/golang/go/issues/27658)


简而言之，错误报告是，如果两种不同的 Goroutines 同时调用 OS.Chdir，则这将生效是不可预测的。

那是真实的。 OS.Chdir 是一个流程范围的属性，而不是每个 goroutine 或 per-thread 属性。即使我们能够弄清楚改变的方式 - 没有什么意思 - 没有什么想到的 - 我们现在无法改变它，因为它会破坏一个大峡谷中调用 OS.Chdir 的现有 Go 程序，并期望它会影响另一个大花序。
Closing as unfortunate.
关闭不幸。
