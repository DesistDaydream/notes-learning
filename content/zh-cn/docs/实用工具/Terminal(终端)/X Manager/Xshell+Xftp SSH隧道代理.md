---
title: Xshell+Xftp SSH隧道代理
---

# 概述

出于安全考虑，公司的一组应用服务器仅允许特定 P 远程 SSH 访问，带来安全防护的同时也增加了进行 SSH 登录和 SFTP 上传维护的繁琐，在授权的 IP 服务器上搭建 VPN 作为跳板是一种解决方案，本文阐述的，是另一种更加简单的安全访问方式，主要是基于日常维护所使用的 Xshell 和 Xftp 工具来配置（这两个工具实在是太方便了）。

为了方便阐述，先上一张网络结构示意图，如下:
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329680-2d7e8805-4094-4508-91f4-3b7e6696b35a.jpeg)
如上图所示，服务器组 C 仅允许授权的服务器 B 访问，计算机 A 无法直接访问维护服务器组 C，本文就是要通过 SSH 隧道代理配置，实现计算机 A 可以直接通过安全隧道代理访问服务器组 C。

1、打开 Xshell，新建管理维护机 A 到授权 IP 服务器 B 的正常连接。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329609-f1e0bc09-7dad-4644-ba12-8f8f2af2730e.jpeg)

配置此连接的登录帐号和密码，如下图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329658-3027a6d6-78fc-4553-b2ba-0e01e8718c73.jpeg)
点击左侧的“隧道”，在右侧点击“添加”，添加 SSH 隧道。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329664-0bbc8564-8c91-471d-a8e0-219afdc0a189.jpeg)

在“转移规则”配置窗口，选择 SOCKS4/5，侦听端口默认为 1080，如果此端口不可用，则需要修改。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329653-72d53f54-8e60-4525-a86c-924e034e9f84.jpeg)

确定后点击连接到代理服务器，连接后在下方的转移规则里如果出现下图所示的信息，则代理服务器配置成功完成。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329770-1cfa9c7f-71c9-47ba-b030-4c133eb5f885.jpeg)

2、在 Xshell 新建到服务器组 C 的连接

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329699-d21ef1c6-790c-46ee-beae-fbc6e74ef486.jpeg)

点击左侧的“代理”，在右侧点击代理服务器后的“浏览”按钮

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329719-5e602fbb-98e1-4b5d-9b04-dd896a1ce785.jpeg)

在弹出的代理服务器设置窗口，输入名称（自己定义），类型选择“SOCKS5”，主机填写"localhost"，端口 1080（如果第一步配置时修改了端口，此处也需要对应一致）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329690-70893f3e-f8a2-4f68-aaaa-3f40e943ada9.jpeg)

点击确定后列表代理里会出现刚刚配置的代理

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329719-bb507607-6d78-484c-8618-17998d058b38.jpeg)

确定后返回刚才的连接窗口，在代理服务器下拉列表中选择刚才建立的代理

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329783-e8ec6628-01a4-43cd-9f17-66bf4630ff39.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329725-d854eb26-addd-4291-a25b-b47963440891.jpeg)

确定后，即可如同直连一样使用计算机 A 远程连接维护服务器组 C 了。

在此处配置成功之后，Xftp 无需配置，只需要在新建到服务器组 C 的连接属性里选择刚才建立的代理服务器即可。

# xshell 通过堡垒机登录服务器

参考：原文链接

具体方案：

1. 首先通过 ssh 的方式登录堡垒机（其实到这一步也是可以的，登录上堡垒机后会提示你要链接的哪个服务器，然后选择登录即可），并通过隧道对服务器建立监听

2. 通过 ssh 链接上面创建隧道的监听

## 配置登录堡垒机

`新建会话`，填写堡垒机的`地址`与`端口号以及堡垒机的账号密码`，如图所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329768-b5f521e5-5e0f-43c1-97f1-9e7ab4c514ee.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329835-560e52dc-0eb6-4e14-9160-a41fcecf9e66.png)

## 配置隧道

**通过本地（拔出）方式：**

接下来继续配置连接内网服务器的隧道，点击隧道再点击添加按钮进入隧道添加页面，源主机为本机 localhost，侦听端口可以在有效范围内随便填写，这里为了区分连接内网哪台服务器，所以用内网服务器 ip 最后一位加 22 即 522 作为侦听端口。目标主机就是我们要通过跳板机访问的内网主机，端口是 22。同样的操作再配置一个连接 192.168.100.6 的隧道，端口不能与 522 冲突，按刚才的规则可以用 622 端口。具体配置如下图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329710-35d6f97a-21e4-4d3e-b1c1-67e064248b7c.png)

**或者通过 Dynamic：**

`类型`选择`Dynamic(SOCKS4/5)`，`侦听端口`任意填写，保存连接并登录，如图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329732-f80aa9d2-51a2-4f56-8a78-e8e897161c03.png)

## 通过 ssh 链接隧道登录服务器

**通过本地（拔出）方式：**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329869-c57675df-1fa1-4835-93c9-379604b1c544.png)

**通过 Dynamic：**

**重现新建**正常会话窗口，依次填写目标机器`地址`、`端口号`、`用户名`与`密码`，在`代理`选项中，`名称`任意填写，`类型`选择`SOCKS5`，`主机`填写`localhost`，`端口`与**第三步**保持一致，返回后选择此服务器并确定，如图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329849-c16fe5e0-5282-429f-8ed9-b3a20431cf43.png)

# xshell 登录脚本配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329720-23e177e6-1ff6-4d57-94c6-8ab63575650f.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329825-d3f82356-a350-4a21-bb36-9d7d2135e256.png)

通过这个脚本可以实现自动连接服务器然后登录操作。

# 快速命令集的使用

- **功能位置：**

1. 工具 》 快速命令集

2. 文件 》 属性 》 高级 》 快速命令集 》 浏览

- **设置快速命令集：**

选择快速命令集 》 新建 》添加 》 然后设置标签（这个很重要哟，下面就可以看到了，一定要是自己懂的，切记不要随便起一些乱七不糟自己都会忘记意思的） 》 选择操作，这里可以选择执行脚本和发送文本，我选择了发送文本，执行脚本需要自己上传脚本。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329790-abe6998f-ad84-4829-9c8e-92849a82196b.png)

- **设置显示快速命令集**

1. 查看 》 勾选快速命令集

2. 在这个位置会显示配置的快速命令集

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eaepgt/1622790329702-fd86989d-737a-471e-9ad8-f81710a0bdb7.png)
