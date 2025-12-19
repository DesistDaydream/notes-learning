---
title: "Algorithm"
linkTitle: "Algorithm"
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Algorithm](https://en.wikipedia.org/wiki/Algorithm)

**Algorithm(算法)**

# Complexity

> 参考：
>
> - [Wiki, Computational complexity](https://en.wikipedia.org/wiki/Computational_complexity)
> - [Wiki, Time complexity](https://en.wikipedia.org/wiki/Time_complexity)
> - [B 站，常见的大O表示法有哪些？时间复杂度是什么？](https://www.bilibili.com/video/BV1DY4y1H7DG)

**Complexity(复杂度)** 指运行某算法所需的资源量。通常不会精确计算，而是用一种表示法来表示一个数值的范围。

特别关注的是计算时间和存储空间。问题的复杂度则是指解决该问题的最佳算法所具有的复杂度。

常用 **O()** 表示。如下图时间复杂度所示，越靠左上角的复杂度越糟糕。

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/program/algorithm/complexity_o.png)

时间复杂度通常指一个函数运行完成需要执行某些行代码的次数。比如：

- 恒定时间的复杂度是 $O(1)$
- for 循环复杂度是 $O(n)$，i.e. 循环 n 次。查找
- 嵌套 for 循环复杂度通常是 $O(n^2)$
- 二分查找的复杂度是 $O(\log_{}n)$
- etc.

> [!Note]
>
> 复杂度的计算并不是一个精确的计算，而是一种在宏观上通过表示法来表示某种无法确定传入参数的算法所需要消耗的资源
>
> 比如，上面的嵌套循环 2 次方只是方便描述，是一种抽象的表示法，并不是真的只循环 2 次并且每次都一样，只是用来表示一种宏观上需要消耗的时间。
>
> 如果 n 也确定了，2 也确定了，比如 1000 个元素遍历 2 遍，那这种情况的时间复杂度应该用 $O(1)$ 表示，因为传入参数是固定的，没有复杂度。
