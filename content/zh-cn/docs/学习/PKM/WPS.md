---
title: WPS
linkTitle: WPS
date: 2024-04-02T22:08
weight: 20
---

# 概述

> 参考：
> -

垃圾软件

安装时修改路径后提示没有权限，然后点继续就直接安装到 C 盘了，点取消才能安装到其他路径。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wps/wps_fuck_1.png)


# 关联文件与配置

- **C:/ProgramData/kingsoft/** # 安装程序运行时下载的文件保存路径
- **C:/rogram Files (x86)/Kingsoft/** # 不知道干啥用的
- C:/Users/DesistDaydream/AppData/Local/kingsoft/ # 不知道干啥用的
- C:/Users/DesistDaydream/AppData/Roaming/kingsoft/ # 不知道干啥用的
- **%我的文档%/KingsoftData/** # 
  - 暂时不知道win中“我的文档”的变量是什么。
- **%我的文档%/WPS Cloud Files/**
  - `./${ACCOUNT_ID}/cachedata/${RANDOM_NUM}/` # 金山文档通过 wps 方式打开后，会将文件缓存到该目录。每个文件占一个目录。
