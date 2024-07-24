---
title: Containerd 部署
linkTitle: Containerd 部署
date: 2024-07-24T20:09
weight: 2
---

# 概述

> 参考：
>
> - [GitHub 文档，containerd/containerd/docs/getting-started.md](https://github.com/containerd/containerd/blob/main/docs/getting-started.md)

我们可以在官方 README 中的 [Runtime Requirements](https://github.com/containerd/containerd#runtime-requirements) 处找到当前 Containerd 版本所依赖的各种组件所需的版本，比如 runc 的版本等。

- 依赖的 runc 版本通常记录在 [containerd/script/setup/runc-version](https://github.com/containerd/containerd/blob/main/script/setup/runc-version) 文件中

# 安装 Containerd

是否需要 libseccomp2 依赖？待验证

## 使用包管理器安装

### CentOS

```bash
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum install -y containerd.io
```

### Ubuntu

```bash
sudo apt-get -y install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get -y update
sudo apt-get -y install containerd.io
```

### 配置 unit 文件

略

## 使用二进制文件安装

通常，我们使用二进制安装 Containerd 时，除了 Containerd 的本体，还需要安装 runc 与 CNI。

注意：在 1.6.0 版本的更新说明中，对 Releases 中的包进行一些调整，将在未来的 2.0 版本之后弃用一些东西

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gg4dmt/1654268589875-dff18ee5-d643-489f-9c13-b8bc1dd4e99d.png)

这里说的主要是对那些整合包的弃用，让 containerd 的 realease 更纯粹，那些带着 cni 或 cri 的整合包，都没有了。并且，根据 1.6 版本的官方文档的安装说明，CRI 功能已经整合在 containerd 中，所以更无须下载整合包了。

### 安装 Containerd

在 [release](https://github.com/containerd/containerd/releases) 页面下载二进制程序压缩包，解压并将二进制程序放到 $PATH 中

```bash
export ARCH="amd64"
export CONTAINER_VERSION="1.6.16"
export OS="linux"
wget https://github.com/containerd/containerd/releases/download/v${CONTAINER_VERSION}/containerd-${CONTAINER_VERSION}-${OS}-${ARCH}.tar.gz
tar Cxzvf /usr/local containerd-${CONTAINER_VERSION}-${OS}-${ARCH}.tar.gz
```

### 安装 runc

从 [GitHub 项目，opencontainers/runc 的 Releases](https://github.com/opencontainers/runc/releases) 处下载 `runc.<ARCH>` 二进制文件，拷贝到 `/usr/local/sbin/runc` 处

```bash
export ARCH=amd64
cp runc.${ARCH} /usr/local/sbin/runc && chmod 755 /usr/local/sbin/runc
```

### 安装 CNI 插件

从 [GitHub 项目，containernetworking/plugins 的 Releases](https://github.com/containernetworking/plugins/releases) 处下载 `cni-plugins-<OS>-<ARCH>-<VERSION>.tgz` 文件，解压到 `/opt/cni/bin/` 目录下

```bash
export OS=linux
export ARCH=amd64
export VERSION=v1.1.1
mkdir -p /opt/cni/bin
tar Cxzvf /opt/cni/bin cni-plugins-${OS}-${ARCH}-${VERSION}.tgz
```

### 配置 Unit 文件

从 [GitHub 项目文件，containerd/containerd/containerd.service](https://github.com/containerd/containerd/blob/main/containerd.service) 中下载用于 Systemd 的 Unit 文件。(对于 cri-containerd-.... 类型的 release 压缩文件中包含 Unit 文件)

这是一个 1.4.4 版本的 continerd.service 文件样例：

```bash
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target

[Service]
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd

Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=1048576
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
```

# 配置并启动 Containerd

## 添加 containerd 配置文件

通过命令生成配置文件

```bash
containerd config default > /etc/containerd/config.toml
```

## 修改内核参数

```bash
cat > /etc/sysctl.d/containerd.conf << EOF
net.ipv4.ip_forward = 1
EOF
sysctl -p /etc/sysctl.d/*
```

## 启动 containerd

```bash
systemctl daemon-reload
systemctl enable containerd --now
```

# rootless 模式

强烈推荐开启 cgroup v2，否则最好不要使用 rootless 模式，开启参考：https://rootlesscontaine.rs/getting-started/common/cgroup2/

```bash
[INFO] Checking RootlessKit functionality
[INFO] Checking cgroup v2
[WARNING] Enabling cgroup v2 is highly recommended, see https://rootlesscontaine.rs/getting-started/common/cgroup2/
[INFO] Checking overlayfs
[INFO] Requirements are satisfied
[INFO] Creating "/home/desistdaydream/.config/systemd/user/containerd.service"
[INFO] Starting systemd unit "containerd.service"
+ systemctl --user start containerd.service
+ sleep 3
+ systemctl --user --no-pager --full status containerd.service
● containerd.service - containerd (Rootless)
     Loaded: loaded (/home/desistdaydream/.config/systemd/user/containerd.service; disabled; vendor preset: enabled)
     Active: active (running) since Mon 2021-09-13 21:48:44 CST; 3s ago
   Main PID: 2625 (rootlesskit)
     CGroup: /user.slice/user-1000.slice/user@1000.service/containerd.service
             ├─2625 rootlesskit --state-dir=/run/user/1000/containerd-rootless --net=slirp4netns --mtu=65520 --slirp4netns-sandbox=auto --slirp4netns-seccomp=auto --disable-host-loopback --port-driver=builtin --copy-up=/etc --copy-up=/run --copy-up=/var/lib --propagation=rslave /usr/local/bin/containerd-rootless.sh
             ├─2635 /proc/self/exe --state-dir=/run/user/1000/containerd-rootless --net=slirp4netns --mtu=65520 --slirp4netns-sandbox=auto --slirp4netns-seccomp=auto --disable-host-loopback --port-driver=builtin --copy-up=/etc --copy-up=/run --copy-up=/var/lib --propagation=rslave /usr/local/bin/containerd-rootless.sh
             ├─2651 slirp4netns --mtu 65520 -r 3 --disable-host-loopback --enable-sandbox --enable-seccomp 2635 tap0
             └─2659 containerd

9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.601967589+08:00" level=info msg="loading plugin \"io.containerd.grpc.v1.cri\"..." type=io.containerd.grpc.v1
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602049316+08:00" level=info msg="Start cri plugin with config {PluginConfig:{ContainerdConfig:{Snapshotter:overlayfs DefaultRuntimeName:runc DefaultRuntime:{Type: Engine: PodAnnotations:[] ContainerAnnotations:[] Root: Options:map[] PrivilegedWithoutHostDevices:false BaseRuntimeSpec:} UntrustedWorkloadRuntime:{Type: Engine: PodAnnotations:[] ContainerAnnotations:[] Root: Options:map[] PrivilegedWithoutHostDevices:false BaseRuntimeSpec:} Runtimes:map[runc:{Type:io.containerd.runc.v2 Engine: PodAnnotations:[] ContainerAnnotations:[] Root: Options:map[BinaryName: CriuImagePath: CriuPath: CriuWorkPath: IoGid:0 IoUid:0 NoNewKeyring:false NoPivotRoot:false Root: ShimCgroup: SystemdCgroup:false] PrivilegedWithoutHostDevices:false BaseRuntimeSpec:}] NoPivot:false DisableSnapshotAnnotations:true DiscardUnpackedLayers:false} CniConfig:{NetworkPluginBinDir:/opt/cni/bin NetworkPluginConfDir:/etc/cni/net.d NetworkPluginMaxConfNum:1 NetworkPluginConfTemplate:} Registry:{ConfigPath: Mirrors:map[] Configs:map[] Auths:map[] Headers:map[]} ImageDecryption:{KeyModel:node} DisableTCPService:true StreamServerAddress:127.0.0.1 StreamServerPort:0 StreamIdleTimeout:4h0m0s EnableSelinux:false SelinuxCategoryRange:1024 SandboxImage:k8s.gcr.io/pause:3.5 StatsCollectPeriod:10 SystemdCgroup:false EnableTLSStreaming:false X509KeyPairStreaming:{TLSCertFile: TLSKeyFile:} MaxContainerLogLineSize:16384 DisableCgroup:false DisableApparmor:false RestrictOOMScoreAdj:false MaxConcurrentDownloads:3 DisableProcMount:false UnsetSeccompProfile: TolerateMissingHugetlbController:true DisableHugetlbController:true IgnoreImageDefinedVolumes:false NetNSMountsUnderStateDir:false} ContainerdRootDir:/var/lib/containerd ContainerdEndpoint:/run/containerd/containerd.sock RootDir:/var/lib/containerd/io.containerd.grpc.v1.cri StateDir:/run/containerd/io.containerd.grpc.v1.cri}"
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602101237+08:00" level=info msg="Connect containerd service"
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602157949+08:00" level=info msg="Get image filesystem path \"/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs\""
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602220361+08:00" level=warning msg="Running containerd in a user namespace typically requires disable_cgroup, disable_apparmor, restrict_oom_score_adj set to be true"
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602782645+08:00" level=warning msg="failed to load plugin io.containerd.grpc.v1.cri" error="failed to create CRI service: failed to create cni conf monitor: failed to create cni conf dir=/etc/cni/net.d for watch: mkdir /etc/cni/net.d: permission denied"
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.602807316+08:00" level=info msg="loading plugin \"io.containerd.grpc.v1.introspection\"..." type=io.containerd.grpc.v1
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.603004589+08:00" level=info msg=serving... address=/run/containerd/containerd.sock.ttrpc
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.603068510+08:00" level=info msg=serving... address=/run/containerd/containerd.sock
9月 13 21:48:44 hw-cloud-xngy-jump-server-linux-2 containerd-rootless.sh[2659]: time="2021-09-13T21:48:44.603088751+08:00" level=info msg="containerd successfully booted in 0.056085s"
+ systemctl --user enable containerd.service
Created symlink /home/desistdaydream/.config/systemd/user/default.target.wants/containerd.service → /home/desistdaydream/.config/systemd/user/containerd.service.
[INFO] Installed "containerd.service" successfully.
[INFO] To control "containerd.service", run: `systemctl --user (start|stop|restart) containerd.service`
[INFO] To run "containerd.service" on system startup automatically, run: `sudo loginctl enable-linger desistdaydream`
[INFO] ------------------------------------------------------------------------------------------
[INFO] Use `nerdctl` to connect to the rootless containerd.
[INFO] You do NOT need to specify $CONTAINERD_ADDRESS explicitly.

```

没法使用 host network 模式，找不到解决办法，好 TM 麻烦
