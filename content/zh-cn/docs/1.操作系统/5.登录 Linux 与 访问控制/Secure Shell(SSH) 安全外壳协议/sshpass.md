---
title: sshpass
---

# 概述

> 参考：
> - [源码](https://sourceforge.net/projects/sshpass/)
> - [manual](https://man.cx/sshpass)

Sshpass 是一种使用 SSH 所谓的“交互式键盘密码身份验证”非交互式执行密码身份验证的工具。大多数用户应该改用 SSH 更安全的公钥认证。

# Syntax(语法)

**OPTIONS**

- **-p <STRING>** # 指定 ssh 时使用的密码

# 问题示例

## 使用 sshpass -p 密码 ssh root@ip 地址 没有任何反应，解决办法找到了

添加-o StrictHostKeyChecking=no 选项【表示远程连接时不提示是否输入 yes/no】

sshpass -p molihuacha.1 ssh -o StrictHostKeyChecking=no root@188.131.150.204

[
](https://blog.csdn.net/weixin_41831919/article/details/109660760)
