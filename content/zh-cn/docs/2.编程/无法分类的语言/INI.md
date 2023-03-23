---
title: INI
---

# 概述

> 参考：
> - [Wiki,INI](https://en.wikipedia.org/wiki/INI_file)

## INI 格式规范

- `;` 分号表示注释

# INI 原语

## Key/Value pair(键/值对)

INI 格式的文件主要结构是 **Key/Value pair(键/值对)** 格式。Key 与 Value 以 `=` 符号分割。有的地方也称为 **Properties(属性)**。

## Sections(部分)

**Selections(部分)** 是 `键值对` 的集合，也称为 Hash Tables(哈希表) 或 Dictionaries(字典)，以 `[]` 符号表示。从 Table 的 `[]` 符号开始到下一个 `[]` 符号为止，所有键值对都属于该 Sections。

人们日常生活中描述的 第一部分、第二部分、我这部分 等等，这就是 部分的意思，表示一个整体的其中一部分。一个 INI 有很多部分，比如可以说：有 main 部分、logging 部分 等等。

Sections 也有章节的意思，但是不如 Chapter 这个词用来表示章节更合适。
