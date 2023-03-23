---
title: Sublime
---

# 概述

在使用 Git 进行提交操作中，常见需要把 CRLF 转换成 LF 的警告。

个人目前使用代码编辑器是 Sublime Text 3，可以在设置中避免这个问题。设置如下：

Perference->Setting-User 中加入配置 "default_line_ending": "unix"

这个参数有三个可用选项，system, windows, unix

注意添加逗号，效果如下：

    {
        	"ignored_packages":
        	[
        		"Vintage"
        	],
        	"default_line_ending": "unix",
    }
