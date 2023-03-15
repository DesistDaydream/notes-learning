---
title: sublime 相关问题
---

#

# 报错 1：There are no packages available for installation

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/frqwsc/1616161737847-06e6b281-5be4-4877-ad55-a3b5112335a4.jpeg)

原因：

打开如下配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/frqwsc/1616161737912-401ac1c0-f667-4c60-9493-991ddf5a395f.jpeg)

配置中红框位置为 sublime 安装包所需要的仓库内容，该连接被墙，无法打开

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/frqwsc/1616161737893-78f8885e-c90f-4d58-8e7c-86e3c2ce4d6b.jpeg)

# 解决办法

通过翻墙提前把该连接下的 channel_v3.json 文件下载下来保存在本地

然后在 首选项—Package Settings—Package Control—Settings-User 的配置中添加如下内容

    {
     "bootstrapped": true,
     "channels":
     [
      "E:\\Tools\\channel_v3.json" # 指定channel_v3.json文件的路径，让sublime去读取本地的仓库文件
     ],
     "in_process_packages":
     [
     ],
     "installed_packages":
     [
      "ChineseLocalizations",
      "ConvertToUTF8",
      "Package Control"
     ]
    }
