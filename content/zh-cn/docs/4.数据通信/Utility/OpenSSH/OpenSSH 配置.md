---
title: OpenSSH 配置
linkTitle: OpenSSH 配置
weight: 2
---

# 概述

> 参考：
>
> - [官方文档，手册 - sshd_config](https://man.openbsd.org/sshd_config)
> - [官方文档，手册 - ssh_config](https://man.openbsd.org/ssh_config)

# sshd_config 文件

**Port**(NUM) # 设定 sshd 服务监听的端口号

**ListenAddress**(IP) # 设定 sshd 服务监听的 IP 地址(全 0 为所有 IP)

**PermitRootLogin**(yes|no) # 设定是否允许 root 用户通过 ssh 直接登录

**AllowUsers <User1 User2 User3.......>** # 设定允许通过 ssh 登录的用户 User1,2,3 等等

**AllowGroups <Group1 Group2.........>** # 设定允许通过 ssh 登录的组 Group1,2,3,等等等

[KbdInteractiveAuthentication](https://man.openbsd.org/sshd_config#KbdInteractiveAuthentication) #

> 注意：该关键字是已经被启用的 ChallengeResponseAuthentication 关键字的替代品

**PermitTunnel**(STRING) # 指定是否允许 tun 设备转发。可用的值有 yes、point-to-point、ethernet、no。`默认值: no`

- https://man.openbsd.org/sshd_config#PermitTunnel

**GatewayPorts** # 指定是否允许远程主机连接到为客户端转发的端口。

**ForceCommand** #

- https://man.openbsd.org/sshd_config#ForceCommand

**Match** # 创建一个独立的配置块。满足 Match 指定的匹配条件时，Match 所在行下的所有配置都将覆盖

- https://man.openbsd.org/sshd_config#Match

## 连接行为

**ChrootDirectory**(STRING) # 用户登录成功，利用 [Chroot](/docs/1.操作系统/Kernel/Process/Chroot.md) 能力将用户看到的 `/` 目录改为指定的目录。`默认值: none`

> [!attention]
> ChrootDirectory 关键字指定的目录及其所有上级目录必须是 root 用户拥有，并且不能被组用户或其他用户写入。i.e. 这些目录的属主和属组必须是 root，权限最多设置为 755。
>
> 若不满足上述条件，登录将会失败，并产生 `client_loop: send disconnect: Broken pipe` 报错。

**UseDNS**(BOOLEAN) # 指定登陆时是否进行 DNS 解析

## 认证逻辑



## SFTP 子系统相关

**Subsystem** # 配置外部子系统（e.g. [SFTP Subsystem](/docs/4.数据通信/Utility/OpenSSH/SFTP%20Subsystem.md)）。可用参数详见对应子系统中的参数。

# ssh_config 文件

> [!Notes]
> 该配置文件中的配置信息，可以通过 ssh 命令的 -o 选项来覆盖配置文件中关键字的值

**Tunnel** # 启用隧道功能后, 客户端创建的 tun 设备的类型. 默认为 point-to-point. 该配置的功能与 sshd_config 中的 PermitTunnel 配置一样, 用来指定隧道功能下的 tun/tap 网络设备的类型.

**TunnelDevice LOCAL_TUN\[:REMOTE_TUN]** # 与 ssh 的 -w 选项功能一致.

**StrictHostKeyChecking**

If this flag is set to ’’accept-new’’ then ssh will automatically add new host keys to the user known hosts files, but will not permit connections to hosts with changed host keys. If this flag is set to ’’no’’ or ’’off’’, ssh will automatically add new host keys to the user known hosts files and allow connections to hosts with changed hostkeys to proceed, subject to some restrictions. If this flag is set to **ask** (the default), new host keys will be added to the user known host files only after the user has confirmed that is what they really want to do, and ssh will refuse to connect to hosts whose host key has changed. The host keys of known hosts will be verified automatically in all cases.

- **yes** # 则 ssh(1) 将永远不会自动将主机密钥添加到 \_~/.ssh/known_hosts \_文件中，并且拒绝连接主机密钥已更改的主机。尽管当 \_/etc/ssh/ssh_known_hosts \_文件维护不当或经常与新主机建立连接时可能很烦人，但这可以最大程度地防止中间人（MITM）攻击。此选项强制用户手动添加所有新主机。
- **accept-new** # 那么 ssh 会自动将新的主机密钥添加到用户已知的主机文件中，但不允许使用更改的主机密钥连接到主机。
- **no** 或 **off** # 则 ssh 会自动将新的主机密钥添加到用户已知的主机文件中，并允许在更改主机密钥的情况下继续进行主机连接。
- **ask(默认值)** # 则只有在用户确认他们确实要这样做之后，新的主机密钥才会添加到用户已知的主机文件中，并且 ssh 将拒绝连接到其主机密钥的主机已经改变。在所有情况下，都会自动验证已知主机的主机密钥。

效果如下：

```bash
# 严格检查，默认值 ask
~]# ssh 172.19.42.248
The authenticity of host '172.19.42.248 (172.19.42.248)' cant be established.
ECDSA key fingerprint is SHA256:dugyXVC21RvaDTtRp/cBTsqr0MPtjhBJmtjmzZTXljo.
Are you sure you want to continue connecting (yes/no/[fingerprint])? ^C
# 不严格检查，改为 no
~]# ssh -o 'StrictHostKeyChecking=no' 172.19.42.248
Warning: Permanently added '172.19.42.248' (ECDSA) to the list of known hosts.
root@172.19.42.248's password:
```

## 认证逻辑

**HostKeyAlgorithms**(STRING) # 设定主机密钥签名算法。多个算法以 `,` 分割，登陆时会按照顺序注意使用指定的算进行连接

- 可以使用 `+` 将指定的算法添加到默认的算法中；使用 `-` 将指定的算法从默认的算法中删除。
- 可以使用 `ssh -Q HostKeyAlgorithms` 查看都有哪些算法可用

> [!Tip]
> 有的 SSH 服务端设定了某些算法，若客户端默认算法没指定，则连接时将会报错: `Unable to negotiate with 1.1.1.1 port 22: no matching host key type found. Their offer: XXXX,YYYY`。此时就可以利用 `-o HostKeyAlgorithms=+XXXX,YYYY` 正常连接