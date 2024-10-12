---
title: Jenkins 页面详解
---

#

流水线任务

General

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ab7yqi/1616077861969-3249ea3f-1fe0-4313-aa5f-912913c7e656.jpeg)

GitHub 项目 #

This build requires lockable resources #

Throttle builds #

丢弃旧的构建 #

参数化构建过程 # 指定各种参数以便在脚本或者 pipeline 中使用。

1. 字符参数 # 指定变量，名称为变量名，默认值为变量值。可以直接在脚本中引用。

关闭构建

在必要的时候

构建触发器

用来指定触发本次任务的方式，比如提交代码时触发，或者定时触发等等。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ab7yqi/1616077861983-48fd85c4-6f88-4884-b8b9-57888efcd5b3.jpeg)

流水线

用来指定 Jenkinsfile 的内容。

定义 # 指定 Jenkinsfile 的获取方式。

1. Pipeline script # 直接在网页编写 Jenkinsfile 内容。

2. Pipeline script from SCM # 从项目根目录中查找 Jenkinsfile 文件

   1. SCM # 指定获取 Jenkinsfile 的 SCM 及其 URL、认证、分支。

   2. 脚本路径 # 指定 Jenkinsfile 所在位置的绝对路径

自由风格的软件项目

General

GitHub 项目 #

This build requires lockable resources #

Throttle builds #

丢弃旧的构建 #

参数化构建过程 # 指定各种参数以便在脚本或者 pipeline 中使用。

1. 字符参数 # 指定变量，名称为变量名，默认值为变量值。可以直接在脚本中引用。

关闭构建

在必要的时候
