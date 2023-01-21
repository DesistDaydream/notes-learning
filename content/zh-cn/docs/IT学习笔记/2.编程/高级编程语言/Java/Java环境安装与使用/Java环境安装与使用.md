---
title: "Java 环境安装与使用"
linkTitle: "Java 环境安装与使用"
weight: 20
---

# 概述
> 参考：
> - [廖雪峰-Java 教程，安装 JDK](https://www.liaoxuefeng.com/wiki/1252599548343744/1280507291631649)


# 安装 Java

安装好的 JavaSE 包含很多可执行程序
-   java：这个可执行程序其实就是JVM，运行Java程序，就是启动JVM，然后让JVM执行指定的编译后的代码；
-   javac：这是Java的编译器，它用于把Java源码文件（以`.java`后缀结尾）编译为Java字节码文件（以`.class`后缀结尾）；
-   jar：用于把一组`.class`文件打包成一个`.jar`文件，便于发布；
-   javadoc：用于从Java源码中自动提取注释并生成文档；
-   jdb：Java调试器，用于开发阶段的运行调试。


## Windows 安装

从[这里](https://www.oracle.com/java/technologies/downloads/)下载 JavaSE

解压到指定目录(我通常是在 `D:/Tools/Java/jdk-${VERSION}`)，将该目录添加到 JAVA_HOME 环境变量

将 `%JAVA_HOME%/bin` 添加到 PATH 变量中。