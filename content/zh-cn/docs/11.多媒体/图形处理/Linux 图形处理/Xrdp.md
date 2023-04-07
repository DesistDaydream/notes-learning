---
title: Xrdp
---

# 概述

> 参考：
> - [GitHub 项目，neutrinolabs/xrdp](https://github.com/neutrinolabs/xrdp)
> - <https://tecadmin.net/how-to-install-xrdp-on-ubuntu-20-04/>

**XRDP** 使用 **Microsoft 的 Remote Desktop Protocol(远程桌面协议，简称 RDP)** 提供了一个图形界面以登录到计算机。XRDP 接受来自各种 RDP 客户端的连接。比如：

- FreeRDP
- rdesktop
- KRDC
- NeutrinoRDP
- Windows MSTSC (Microsoft Terminal Services Client, aka mstsc.exe)
- Microsoft Remote Desktop (found on Microsoft Store, which is distinct from MSTSC)

XRDP 主要针对 Linux 操作系统，可以让 Windows，使用其自带的远程桌面，通过 RDP 协议，远程控制装有 XRDP 的 Linux 系统，以实现远程图形化控制。

# 为 Ubuntu 安装 Xrdp

远程桌面协议允许用户访问远程系统桌面。 XRDP 服务向您提供使用 Microsoft RDP（远程桌面协议）的远程计算机的图形登录。 XRDP 还支持双向剪贴板传输（文本，位图，文件），音频重定向和驱动器重定向（在远程计算机上安装本地客户端驱动器）。
XRDP 是 Ubuntu 系统的易于安装和可配置的服务。但您还可以使用 VNC Server 访问 Ubuntu 系统的远程桌面。在 Ubuntu 20.04 系统上找到安装 VNC 服务器的教程。
本教程可帮助您在 Ubuntu 20.04 Linux 系统上安装远程桌面（XRDP）。还提供在系统上安装桌面环境的说明。

## 安装桌面环境

默认情况下，Ubuntu Server 没有已安装的桌面环境。 XRDP 服务器旨在仅控制桌面系统。因此，您需要向系统添加桌面环境。
使用 `tasksel` 工具安装桌面环境

```bash
sudo apt update && sudo apt upgrade
apt install tasksel -y
```

安装 Tasksel 后，使用命令 `tasksel` 启动程序：

应该看到以下界面：
[![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ff73o6/1630033188975-4b9f059c-5f81-4263-8965-bb83636058bf.png)](https://tecadmin.net/wp-content/uploads/2020/11/tasksel-ubuntu-desktop.png)
使用箭头键向下滚动列表并查找 `Ubuntu minimal desktop`。接下来，按空格键选择它，然后按 Tab 键选择确定，然后按 Enter 安装 `Ubuntu minimal desktop` 桌面。

安装所有包后，`reboot` 重启操作系统

## 安装 xrdp 包

XRDP 软件包可在默认系统存储库下提供。您可以通过执行以下命令在 Ubuntu 系统上安装远程桌面。

```bash
sudo apt install xrdp -y
```

## 配置 XRDP

在安装过程中，XRDP 在您的系统中添加了名为“XRDP”的系统中的用户。 XRDP 会话使用证书密钥文件“/etc/sl/private/ssl-cert-snakeoil.key”，它与远程桌面播放一个重要的角色。
要正确工作，请使用以下命令将 XRDP 用户添加到“SSL-Cert”组。

```bash
sudo usermod -a -G ssl-cert xrdp
```

有时用户面临黑屏的问题出现在后台。因此，我是一个介绍在背景中解析黑屏问题的步骤。在文本编辑器中编辑 XRDP 文件 /etc/xrdp/startwm.sh：

在 测试 和执行 XSession 的命令之间添加这些命令:

```bash
Unset DBUS_SESSION_ADDRESS
Unset XDG_RUNTIME_DIR
```

如下所示：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ff73o6/1630060168909-03b02b6b-c076-479b-ab43-ccad7d2993ff.png)
通过运行下面给出的命令重新启动 XRDP 服务：

```bash
sudo systemctl restart xrdp
```

## 连接到远程桌面

xrdp 默认监听 3389 端口，防火墙放通该端口

XRDP 服务已成功安装并准备连接。在 Windows 客户端上，打开 RDP 客户端并输入 Ubuntu 系统的 IP 地址。

您可以通过在运行窗口或命令行中键入“mstsc”来启动 RDP 客户端。
[![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ff73o6/1630033188968-627c5a91-5c70-4d84-8b6d-5513705883ad.png)](https://tecadmin.net/wp-content/uploads/2021/06/connect-xrdp.png)
一旦连接成功，远程系统提示用于身份验证。输入远程 Ubuntu 系统的登录凭据以获取远程桌面访问。
[![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ff73o6/1630033190060-961f359c-adbe-4058-84fc-fa92a0ff780e.png)](https://tecadmin.net/wp-content/uploads/2021/06/xrdp-enter-login-credentials.png)

# Xrdp 配置

**/etc/xrdp/\*** #

- **./sesman.ini** #
- **./xrdp.ini** #

# Ubuntu Desktop 可能需要优化的地方

## sorry, that didn’t work please try again(抱歉，登录失败，请再试一次) 问题

注释掉下面文件文件中 `auth	required	pam_succeed_if.so user != root quiet_success` 这一行

```bash
sudo sed -i.bak 's/^auth.*quiet_success/#&/' /etc/pam.d/gdm-autologin
sudo sed -i.bak 's/^auth.*quiet_success/#&/' /etc/pam.d/gdm-password
```

## 关闭 Ubuntu 20.04 的桌面动画效果

```bash
gsettings set org.gnome.desktop.interface enable-animations false
```

注意：改命令仅对当前用户生效。

## xrdp_mm_process_login_response: login failed

问题原因：远程桌面没有正确关闭，虽然在 windows 系统关闭远程桌面连接，但是在里 linux 上的进程还在运行，导致连接数量达到上限，出现问题。

解决：通过设置 sesman.in 文件内的参数解决：

```bash
sudo sed -i.bak 's/KillDisconnected=false/KillDisconnected=true/' /etc/xrdp/sesman.ini
```

可以修改会话设置 ：将最大会话限制该大 MaxSessions=50; 将 KillDisconnected=1；则每次断开连接时，linux 都会关闭会话进程。 然后重启/etc/init.d/xrdp restart 就可解决问题

## 其他未知问题

sudo sed -i.bak "4 a # Improved Look n Feel Method\ncat <<EOF > ~/.xsessionrc\nexport GNOME_SHELL_SESSION_MODE=ubuntu\nexport XDG_CURRENT_DESKTOP=ubuntu:GNOME\nexport XDG_CONFIG_DIRS=/etc/xdg/xdg-ubuntu:/etc/xdg\nEOF\n" /etc/xrdp/startwm.sh

sudo sed -i 's/allowed_users=console/allowed_users=anybody/' /etc/X11/Xwrapper.config\[

]\(https://blog.csdn.net/u014447845/article/details/80291678)
在 /etc/profile.d/ 目录中添加 `source <(kubectl completion bash)` 这种命令会导致 xrdp 登录失败，现象是输入完用户名和密码之后就闪退。
\[

]\(https://blog.csdn.net/u014447845/article/details/80291678)
