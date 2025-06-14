---
title: VS Code 中的 WSL
---

# 概述

> 参考：
>
> - [官方文档，远程-WSL](https://code.visualstudio.com/docs/remote/wsl)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/pz2gno/1636522141650-89ca683f-8e68-4305-879f-ce3e78f565fd.png)

使用指定发行版的 WSL

> 使用第一个的话，go 环境始终无法获取环境变量

## WSL 代理

VS Code 中的 WSL 将会继承 VS Code 的代理配置。

若启动 VS Code 之前，已启动系统代理(比如 Clash for windows)，则 VS Code 中的 WSL 也会使用代理，而此时 VS Code 的默认代理是 127.0.0.1:7890，此时 WSL 继承了这个配置，那么访问 WSL 中的这个地址是有问题的（此时 Ubuntu 中可并没有运行 Clash，也自然不会监听本地 7890 端口）

所以，若想让插件使用本地连接，启动 VS Code 之前不要开启系统代理。

# VS Code 部分扩展联网失败的问题

https://kawnnor.com/wsl-vscode-proxy

在 VS Code 中，通过 Remote - WSL 打开 Ubuntu 子系统中的项目

- IntelliCode 扩展无法下载模型，WakaTime 扩展无法上传统计数据
- ChatGPT 无法访问

## 猜测原因

Clash for Windows + TUN Mode 对 Ubuntu 子系统中的 vscode-server 没有起到代理的作用，但是 IntelliCode 和 WakaTime 在没有代理的情况下应该也可以正常工作，所以怀疑 vscode-server 被配置了错误的代理。

开启 Clash for Windows 之后，Windows 系统代理一般被设置为 `127.0.0.1:7890`，VS Code 会继承这个代理设置，Ubuntu 子系统中的 vscode-server 应该也继承了这个代理设置。但是 Ubuntu 子系统在 `127.0.0.1:7890` 上并没有代理服务，导致 vscode-server 联网失败。

## 解决方法

> 说白了，就是让 wsl 启动时自动配置各种代理，以便让 vscode 可以继承 wsl 中的代理配置，这样 vscode 中的插件也可以使用这个代理配置了

打开 VS Code 设置，并搜索 "proxy" 可以看到：

![vscode_proxy_settings.png](https://cdn.hashnode.com/res/hashnode/image/upload/v1659638904404/U2-J89Pys.png?auto=compress,format&format=webp)

VS Code 的代理设置，会从 "http_proxy" 和 "https_proxy" 环境变量中继承，所以在 Ubuntu 子系统中正确设置 "http_proxy" 和 "https_proxy" 环境变量应该可以解决问题。

在 Ubuntu 子系统的 `~/.bashrc` 文件末尾加上下面代码：

```
# 获取 Host IP
WINDOWS_IP=$(ip route | grep default | awk '{print $3}')
PROXY_HTTP="http://${WINDOWS_IP}:7890"

# 设置环境变量
export http_proxy="${PROXY_HTTP}"
export https_proxy="${PROXY_HTTP}"
```

重新在 VS Code 中打开 Ubuntu 子系统的远程项目，发现 IntelliCode 和 WakaTime 仍然无法正常工作，经排除发现 Windows 防火墙阻止了 Ubuntu 子系统访问 Host IP。

在 Windows 防火墙中增加一条规则：

```
New-NetFirewallRule -DisplayName "WSL" -Direction Inbound -InterfaceAlias "vEthernet (WSL)" -Action Allow
```

另外，Clash for Windows 需要启用 Allow LAN。

至此，问题解决。
