---
title: run 运行容器
---

# 概述

> 参考：
>
> - [官方文档，参考-命令行参考-docker-Docker run 参考](https://docs.docker.com/engine/reference/run/)
> - [官方文档，参考-命令行参考-docker-docker run](https://docs.docker.com/engine/reference/commandline/run/)

# Syntax(语法)

**docker run \[OPTIONS] ImageName \[COMMAND] \[ARG...]**

## OPTIONS

- --add-host list Add a custom host-to-IP mapping (host:ip)
- -a, --attach list Attach to STDIN, STDOUT or STDERR
- --blkio-weight uint16 Block IO (relative weight), between 10 and 1000, or 0 to disable (default 0)
- --blkio-weight-device list Block IO weight (relative device weight) (default \[])
- --cgroup-parent string Optional parent cgroup for the container
- --cidfile string Write the container ID to the file
- **-d, --detach** # 让容器运行在后台并打印出容器的 ID
- --detach-keys string Override the key sequence for detaching a container
- --device list Add a host device to the container
- --device-cgroup-rule list Add a rule to the cgroup allowed devices list
- --device-read-bps list Limit read rate (bytes per second) from a device (default \[])
- --device-read-iops list Limit read rate (IO per second) from a device (default \[])
- --device-write-bps list Limit write rate (bytes per second) to a device (default \[])
- --device-write-iops list Limit write rate (IO per second) to a device (default \[])
- --disable-content-trust Skip image verification (default true)
- --dns list Set custom DNS servers
- --dns-option list Set DNS options
- --dns-search list Set custom DNS search domains
- --entrypoint string Overwrite the default ENTRYPOINT of the image
- **-e, --env \<LIST>**# 设定容器内的环境变量。LIST 格式为 `VAR=VALUE`，若要指定多个变量，则使用多次 --env 选项。
- --env-file list Read in a file of environment variables
- **--expose \<LIST>** # 等效于 Dockerfile 中的 EXPOSE 指令，仅暴露容器端口，不在宿主机暴露。
- --group-add list Add additional groups to join
- --health-cmd string Command to run to check health
- --health-interval duration Time between running the check (ms|s|m|h) (default 0s)
- --health-retries int Consecutive failures needed to report unhealthy
- --health-start-period duration Start period for the container to initialize before starting health-retries countdown (ms|s|m|h)(default 0s)
- --health-timeout duration Maximum time to allow one check to run (ms|s|m|h) (default 0s)
- -h, --hostname \<STRING> # 指定容器内的 hostname
- --init Run an init inside the container that forwards signals and reaps processes
- **-i, --interactive** # 即使没有 attach 到容器，也保持 STDIN(标准输入)开启。通常与 -t 一起使用
- --ip string IPv4 address (e.g., 172.30.100.104)
- --ip6 string IPv6 address (e.g., 2001:db8::33)
- --ipc string IPC mode to use
- --isolation string Container isolation technology
- --kernel-memory bytes Kernel memory limit
- -l, --label list Set meta data on a container
- --label-file list Read in a line delimited file of labels
- --link list Add link to another container
- --link-local-ip list Container IPv4/IPv6 link-local addresses
- --log-driver string Logging driver for the container
- --log-opt list Log driver options
- --mac-address string Container MAC address (e.g., 92:d0:c6:0a:29:33)
- --mount mount Attach a filesystem mount to the container
- **--name \<STRING>** # 为容器分配一个名称。默认为随机字符串
- **--network \<STRING>** # 连接一个容器到一个容器网络(default "default")，可以是 docker network ls 列出的网络，也可以是其余 container 的网络。STRING 包括下面几种
  - none # 容器使用自己的网络（类似--net=bridge），但是不进行配置
  - bridge # 通过 veth 接口将容器连接到默认的 Docker 桥(默认为 docker0 的网桥).
  - host # 直接使用宿主机的网络而不是独立的 network namespace
  - ContainerName # 连接到指定 container 的网络中
  - NetworkName # 连接到 docker network ls 所列出的其中一个 docker 网络上
- --network-alias list Add ne twork-scoped alias for the container
- --no-healthcheck Disable any container-specified HEALTHCHECK
- --oom-kill-disable Disable OOM Killer
- --oom-score-adj int Tune host's OOM preferences (-1000 to 1000)
- --pid string PID namespace to use
- --pids-limit int Tune container pids limit (set -1 for unlimited)
- --privileged Give extended privileges to this container
- **-p, --publish \[HostIP:]\[HostPort:]\<ContainerPort>**# 指明 Container 要映射到 Host 上的 IP 和端口。若只指明 HostIP 和 ContainerPort 则中间俩个冒号不可省。若不指定 HostIP，则第一个冒号可不写。要暴露多个端口则多次使用 -p 即可。
- **-P, --publish-all** # 将 Image 定义的 EXPOSE 要暴露的端口暴露给 host，随机分配 host 上的端口与之建立映射关系。一般从 10000 端口开始
- **--read-only** # 将容器的根文件系统挂载为只读模式
- **--rm** # 当容器退出时，删除它。包括创建的 volume 等一并删除
- --runtime string Runtime to use for this container
- --security-opt list Security Options
- --shm-size bytes Size of /dev/shm
- --sig-proxy Proxy received signals to the process (default true)
- --stop-signal string Signal to stop a container (default "SIGTERM")
- --stop-timeout int Timeout (in seconds) to stop a container
- --storage-opt list Storage driver options for the container
- --sysctl map Sysctl options (default map\[])
- --tmpfs list Mount a tmpfs directory
- **-t, --tty** # 为此命令分配一个 pseudo-TTY(伪终端)，可以支持终端登录，通常与-i 一起使用。
- **-u, --user \<STRING>** # 为容器进程指定运行的用户名/UID
  - STRING 格式：`<NAME|UID>[:<GROUP|GID>])`
- --userns string User namespace to use
- --uts string UTS namespace to use
- **-v, --volume \[SRC:]DST** # 为容器创建一个 Volume 并挂载到其中的目录上。若指定的 host 上的路径不存在，则自动创建这个目录；若不指定 SRC 则 docker 会自动创建一个。默认在 /var/lib/docker/volumes/ 目录下创建 volume 所用的目录
  - Note：使用 /HOST/PATH 与 VolumeName 的区别详见：《[Docker 存储](/docs/10.云原生/2.2.实现容器的工具/Docker/Docker%20存储.md)》
- --volume-driver string Optional volume driver for the container
- **--volumes-from \<ContainerName>** # 运行的新容器从 ContainerName 这个容器复制存储卷来使用
- **-w, --workdir \<STRING>** # 指定容器内的工作目录，让指定的目录执行当前命令

### 资源配置相关选项

- --cpu-period int Limit CPU CFS (Completely Fair Scheduler) period
- --cpu-quota int Limit CPU CFS (Completely Fair Scheduler) quota
- --cpu-rt-period int Limit CPU real-time period in microseconds
- --cpu-rt-runtime int Limit CPU real-time runtime in microseconds
- -c, --cpu-shares int CPU shares (relative weight)
- **--cpus \<INT>** # 容器可使用的最大 CPU 资源
- --cpuset-cpus string CPUs in which to allow execution (0-3, 0,1)
- --cpuset-mems string MEMs in which to allow execution (0-3, 0,1)
- **-m, --memory \<BYTES>** # 内存限制。容器能使用的最大内存
- --mem ory-reservation bytes Memory soft limit
- --memory-swap bytes Swap limit equal to memory plus swap: '-1' to enable unlimited swap
- --memory-swappiness int Tune container memory swappiness (0 to 100) (default -1)
- **--restart \<string>** # 容器的重启策略。`默认值：0`
- **--ulimit \<UlimitDesc>** # 为容器配置 Ulimit。`默认值：[]`
  - 比如：
    - --ulimit nofile=1000 # 限制容器最多能打开 1 万 个文件描述符
    - --ulimit nproc=10 # 限制容器最多能打开 10 个进程

### Linux Capabilities 相关选项

- --cap-add list # Add Linux capabilities
- --cap-drop list # Drop Linux capabilities

# 最佳实践

- docker run -d -p 80:80 httpd
  - 其过程可以简单的描述为
    - 从 Docker Hub 下载 httpd 镜像。镜像中已经安装好了 Apache HTTP Server。
    - 以后台启动 httpd 容器，并将容器的 80 端口映射到 host 的 80 端口。

以后台运行镜像 nginx:latest

- docker run -p 80:80 -v /data:/data -d nginx:latest

在运行 centos 容器的时候，执行 tail 命令。该命令是为了让容器启动后不自动关闭

- docker run -d centos tail -f /dev/null
