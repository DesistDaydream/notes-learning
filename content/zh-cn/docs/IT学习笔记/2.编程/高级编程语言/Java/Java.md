---
title: "Java"
linkTitle: "Java"
weight: 20
---

# 概述
> 参考：
> - [官网](https://www.java.com/)
> - [廖雪峰-Java教程，Java简介](https://www.liaoxuefeng.com/wiki/1252599548343744/1255876875896416#0)

Java 分三个版本：
- Java Standard Edition(标准版，简称 JSE)
- Java Enterprise Edition(企业版，简称 JEE)
- Java Micro Edition(微型版，简称 JME)
```ascii
┌───────────────────────────┐
│Java EE                    │
│    ┌────────────────────┐ │
│    │Java SE             │ │
│    │    ┌─────────────┐ │ │
│    │    │   Java ME   │ │ │
│    │    └─────────────┘ │ │
│    └────────────────────┘ │
└───────────────────────────┘
```


## Java 名词
JDK # Java Development Kit（Java 开发工具包）

JRE # Java Runtime Environment（Java 运行时环境）

JVM # Java Virtual Machin（Java 虚拟机）

JSR # Java Specification Request（Java 规范）

JCP # Java Community Process（Java 社区处理）

JRE 中包含运行 **Java 字节码** 的 JVM 和 库。但是，我们先要使用 JDK 将 Java 源码编译成 Java 字节码。因此 JDK 除了包含 JRE，还提供了编译器、调试器等开发工具。常说的安装 Java，其实就是指安装 JDK。

```ascii
  ┌─    ┌──────────────────────────────────┐
  │     │     Compiler, debugger, etc.     │
  │     └──────────────────────────────────┘
 JDK ┌─ ┌──────────────────────────────────┐
  │  │  │                                  │
  │ JRE │      JVM + Runtime Library       │
  │  │  │                                  │
  └─ └─ └──────────────────────────────────┘
        ┌───────┐┌───────┐┌───────┐┌───────┐
        │Windows││ Linux ││ macOS ││others │
        └───────┘└───────┘└───────┘└───────┘
```


## 学习资料
[菜鸟教程，Java](https://www.runoob.com/java/java-tutorial.html)

[廖雪峰，Java 教程](https://www.liaoxuefeng.com/wiki/1252599548343744)



## Andrioid
[Android与Java的关系](https://zhuanlan.zhihu.com/p/340609888)

# Hello World


# Java 语言关键字


# Java 语言规范
