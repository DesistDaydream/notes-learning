---
title: UUID 与 ULID
---

# 概述

> 参考：
> - [Wiki，UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)
> - [GitHub,ULID](https://github.com/ulid/spec)

## UUID

**universally unique identifier(通用唯一标识符，简称 UUID)** 是一个[128 位的](https://en.wikipedia.org/wiki/128-bit) [标签](https://en.wikipedia.org/wiki/Nominal_number)用于在计算机系统中的信息。还使用术语**全局唯一标识符**( **GUID** )，通常在[Microsoft](https://en.wikipedia.org/wiki/Microsoft)创建的软件中使用。[\[1\]](https://en.wikipedia.org/wiki/Universally_unique_identifier#cite_note-RFC_4122-1)
根据标准方法生成时，UUID 出于实用目的是唯一的。与大多数其他编号方案不同，它们的唯一性不依赖于中央注册机构或生成它们的各方之间的协调。虽然 UUID 被复制的[概率](https://en.wikipedia.org/wiki/Probability)不是零，但它足够接近零，可以忽略不计。[\[2\]](https://en.wikipedia.org/wiki/Universally_unique_identifier#cite_note-2)
因此，任何人都可以创建一个 UUID 并使用它来识别某些东西，几乎可以肯定的是，该标识符不会与已经或将要创建的标识符重复以识别其他东西。因此，独立方标有 UUID 的信息稍后可以合并到单个数据库中或在同一频道上传输，重复的可能性可以忽略不计。
UUID 的采用很普遍，许多计算平台为生成它们和解析它们的文本表示提供支持。

## ULID

**Universally Unique Lexicographically Sortable Identifier(通用唯一的字典可排序标识符，简称 ULID)**
