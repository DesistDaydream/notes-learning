---
title: Y.Windows 管理
weight: 10
---

# 概述

>


# 查看崩溃信息

<https://its401.com/article/CRJ297486/120602345>
特别生气！！！！某一天突然发现拖拽文件拖拽到其他文件夹就会导致资源管理器卡死，然后还以为是自己拖错了，然后越来越频繁。
然后疯狂百度 1.打开控制面板
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qnpbng/1654348939190-57e80915-99a7-4521-992f-683029eff444.png) 2.再进入安全和维护
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qnpbng/1654348939189-d155d75e-ee86-4dce-8ddf-66eb077b7138.png) 3.点击维护查看可靠性历史记录
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qnpbng/1654348939189-f0bf328c-11eb-449c-91d7-a67f2d7c6e84.png) 4.点击关键信息随便个事件进去可以看见因为啥文件导致卡死的。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qnpbng/1654348939160-131bfda2-8f0e-466b-b0d9-941a4542d9c0.png)
我是因为 wps 的某个文件，把 wps 卸了就好了。 5.如果还没解决的话，可以试试利用 Dism 修复系统
管理员打开 cmd 命令行。
直接输入这两条就好了
DISM /Online /Cleanup-image /ScanHealth //这一条指令用来扫描全部系统文件，并扫描计算机中映像文件与官方系统不一致的情况。 DISM /Online /Cleanup-image /RestoreHealth //计算机必须联网，这种命令的好处在于可以在修复时，系统未损坏部分可以继续运行
