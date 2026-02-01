---
title: Linux系统安装问题
linkTitle: Linux系统安装问题
weight: 100
---

# 概述

> 参考：
>
> - 

使用 U 盘安装 CentOS7 卡在”Starting dracut initqueue hook…”

使用 U 盘安装系统的过程中遇到卡在 Starting dracut initqueue hook 这里的情况，过一会会报 timeout 的错误。这是因为安装程序没有找到安装文件的位置。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532248-128a82ff-30cc-4b4d-9db2-efb3509b83fc.png)

安装程序是按照卷标寻找分区的，可以在开机过程选择 Install CentOS 7 后按 tab 编辑开机选项，uefi 启动模式按 e 编辑。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532228-b481918d-40fb-4a5d-9dd8-481d5d55422e.png)

**问题原因：在安装 centos 时，设备无法找到 U 盘来启动安装程序**

在 ios 安装程序，找到/isolinux/isolinux.cfg(uefi 模式的配置文件路径为：/EFI/BOOT/grub.cfg)这个文件，效果如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532249-a2840410-d1be-4f90-b9de-29703f1040b9.png)

此处 inst.stage2=hd:LABEL=CentOS\x207\x20x86_64。这就是造成超时的原因，inst.stage2 这里应该是指向一个具体的路径，如果是 DVD，它的标签就是“CentOS 7 x86_64”，而 U 盘则可能是你自己定义的标签。 这就造成了 DVD 能正常安装，U 盘就不行了。

而为什么 U 盘的标签不是默认的 CentOS 7 x86_64 呢，是因为标签(LABEL)长度超出了 windows 的卷标长度限制(主要是因为这个 U 盘是在 windos 下制作的。。。)，并且 Windows 限制卷标只能使用大写字母，就算输入的是小写，实际上也是大写。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532266-b8d3cc0c-b9aa-421d-a3e0-03c1bb252e65.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532253-69173fc7-c704-459b-8057-7ef2d24837b4.png)

**解决方式：**

### 解决方式 1

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532236-5667c5d2-d80e-428a-a9ae-28a8c02a9f96.png)

1. 直接修改 /isolinux/isolinux.cfg 文件中的内容，将 CentOS\x207\x20x86_64 修改为 CENTOS7
2. 然后在 windows 中修改 U 盘的标签名为 CENTOS7

这样就可以保证使用 U 盘安装时，服务器可以根据标签来读取到 U 盘并进行后续安装了

### 解决方式 2

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nebxk0/1616168532287-1f22f5eb-1b38-4bde-8048-9a5a8065373a.png)

1. 在 windows 中修改 U 盘的标签名为 CENTOS7
2. 在安装系统界面，如右图，按 e 或者 tab 进入编辑模式
3. 修改 LABEL=后面的内容为 CENTOS7

以上两种解决方式都是让安装程序来根据 U 盘标签来识别 U 盘

**注意 UEFI 启动模式则是修改/EFI/BOOT/grub.cfg 这个文件**

---

### 解决方式 3

首先等待 timeout 报错完成，这是会进入 dracut 的简易终端，cd /dev 执行 ls 查看有哪些分区，我的分区有 sda 和 sdb、sdb4，所以 sdb4 就是镜像所在的分区，编辑 inst.stage2=hd:LABEL=CentOS\x207\x2086_64 为 inst.stage2=hd:/dev/sdb4 按 enter 开机，efi 安装按 ctrl+x 启动。之后就能正常开机了。

## 因此，具体操作有 2 个

**1. 直接修改 /isolinux/isolinux.cfg 里 hd:LABEL= 为 U 盘的具体标签，或者将 U 盘的标签修改为 “CentOS 7 x86_64” ，按照正常流程安装即可。如果害怕空格影响，就把 isolinux.cfg 的 label 去掉\x20，同时 U 盘标签也去掉空格。**

**2. 在选择安装 CentOS 时，选择 Install CentOS 7 ,然后修改 按 e 键，进入修改状态，将 hd:LABEL= 修改为 U 盘的标签，或者修改为当前 U 盘在安装机的具体路径，一般为 /dev/hda1 等，栗子：“hd:LABEL=/dev/hda1”，然后按 Ctrl+x 开始执行安装。**

当然，我还是推荐修改标签的方式来解决。
