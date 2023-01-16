---
title: Containerd 问题汇总
---

# 与老版本不兼容问题

使用 nerdctl 通过 containerd 运行容器时报错：

```bash
FATA[0000] failed to create shim: OCI runtime create failed: unable to retrieve OCI runtime error (open /run/containerd/io.containerd.runtime.v2.task/default/210729ebc4386d8e89132a3dea24fa0d67643587af119247837a0f1009d82fa7/log.json: no such file or directory): runc did not terminate successfully: exit status 127: unknown
```

本质是 runc 问题

```bash
~]# runc -v
runc: symbol lookup error: runc: undefined symbol: seccomp_api_get

~]# ldd /usr/bin/runc
	linux-vdso.so.1 (0x00007fffbfbee000)
	libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fd802a37000)
	libseccomp.so.2 => /lib/x86_64-linux-gnu/libseccomp.so.2 (0x00007fd802a15000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd802823000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fd8036f9000)

```

<https://github.com/containerd/containerd/issues/6209>

主要是 libseccomp 这个包的版本问题，`seccomp_api_get`是从 libseccomp 的 2.4.0 版本开始支持的，但是 RedHad 只能更新到 2.3.X。

Ubuntu 可以通过 `apt install libseccomp-dev` 安装解决，但是 CentOS 还不知道怎么解决

runc 1.0.3 - 1.1.1 之间的几个版本报这个错，1.0.2 和 1.1.1 没问题，但是 containerd 1.5.9 默认带的就是 runc 1.0.3

# 各种命令与配置文件

ctr、nerdctl 等命令好像不会操作应用了 /etc/containerd/config.toml 配置文件的 Containerd，但是 crictl 会

```bash
~]# ctr image pull reg.superstor.com/lchdzh/k8s-debug:v1
INFO[0000] trying next host                              error="failed to do request: Head \"https://reg.superstor.com/v2/lchdzh/k8s-debug/manifests/v1\": x509: certificate is not valid for any names, but wanted to match reg.superstor.com" host=reg.superstor.com
ctr: failed to resolve reference "reg.superstor.com/lchdzh/k8s-debug:v1": failed to do request: Head "https://reg.superstor.com/v2/lchdzh/k8s-debug/manifests/v1": x509: certificate is not valid for any names, but wanted to match reg.superstor.com

~]# nerdctl pull reg.superstor.com/lchdzh/k8s-debug:v1
INFO[0000] trying next host                              error="failed to do request: Head \"https://reg.superstor.com/v2/lchdzh/k8s-debug/manifests/v1\": x509: certificate is not valid for any names, but wanted to match reg.superstor.com" host=reg.superstor.com
FATA[0000] failed to resolve reference "reg.superstor.com/lchdzh/k8s-debug:v1": failed to do request: Head "https://reg.superstor.com/v2/lchdzh/k8s-debug/manifests/v1": x509: certificate is not valid for any names, but wanted to match reg.superstor.com

~]# crictl pull reg.superstor.com/lchdzh/k8s-debug:v1
Image is up to date for sha256:c690d4fd64d6622c3721a1db686c2e4cfb559dd1d9f9ff825584a8f56ec02c7f

```

已经配置了私有镜像注册中心，但是 ctr 和 nerdctl 却没有效果。
