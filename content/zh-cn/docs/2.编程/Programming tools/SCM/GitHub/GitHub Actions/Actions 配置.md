---
title: Actions 配置
linkTitle: Actions 配置
weight: 20
---

# 概述

> 参考：
>
> -

在每个仓库的 Setting 页面可以为仓库的 Actions 进行一些配置

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/202404222254729.png)

# Workflow 权限

有的 Action 要想正常运行，要保证 Workflow 使用的默认 GITHUB_TOKEN 具有对 Realease 的<font color="#ff0000">读写权限</font>（比如 softprops/action-gh-release 需要将构建结果上传至 Release 中，这个属于写入操作）。可以在 `https://github.com/${USER}/${REOP}/settings/actions` 页面修改仓库的 Action 中关于 Workflow 的配置

Setting - Actions - General - Workflow permissions

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_acion/202404222221183.png)

参考: https://github.com/softprops/action-gh-release/issues/232#issuecomment-1375588379
