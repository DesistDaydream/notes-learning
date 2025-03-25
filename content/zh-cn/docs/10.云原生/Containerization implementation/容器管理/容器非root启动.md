---
title: 容器非root启动
linkTitle: 容器非root启动
weight: 20
---


# 概述

> 参考：
>
> - [官方文档，Rootless 模式](https://docs.docker.com/engine/security/rootless/)

# 容器非 root 启动改造的经验

> 参考：
>
> - [zhangguanzhang，容器非 root 启动改造的经验](https://zhangguanzhang.github.io/2023/11/03/non-root-containers/)

最近业务容器的非 root 启动改造实战案例经验，后续有新的也更新进来

改造
---------------------

### 前提须知

这里列举些基础知识

#### 使用 root 不安全的举例

虽然 linux 有 user namespace 隔离技术，但是 docker 不支持类似 podman 那样的给每个容器设置范围性的 uidmap 映射（当然 k8s 现在也不支持），并且容器默认配置下的权限虽然去掉了一些。但是容器内还是能对挂载进去的进行修改的，比如帖子 [rm -rf * 前一定一定要看清当前目录](https://www.v2ex.com/t/976554) 老哥的操作：

```plaintext
docker run --rm -v /mnt/sda1:/mnt/sda1 -it alpine
cp /mnt/sda1/somefile.tar.gz .
tar xzvf somefile.tar.gz
cd somefile-v1.0
ls
# 看了看内容觉得不是自己想要的，回上一级目录准备删掉：
cd ..
rm -rf *
```

嗯，alpine 默认的 workdir 是 `/` ，所以删除 `rm -rf /*`。当然还有其他不安全的，所以在业务角度上，我们需要给容器内进程设置在非 root 下最小的运行权限。

#### 设置-USER-还是使用-docker-entrypoint-sh-入口 "设置 USER 还是使用 docker-entrypoint.sh 入口

Dockerfile 里设置 `USER` 或者 run 的时候设置 `-u user:group` 只能针对于一些简单的进程，例如大部分 exporter 和一些只是用 http API 的进程，这几天我测试后也提交了一些 pr：

- [danielqsj/kafka_exporter](https://github.com/danielqsj/kafka_exporter/pull/410)
- [ClickHouse/clickhouse_exporter](https://github.com/ClickHouse/clickhouse_exporter/pull/83)
- [kubernetes addonresizer](https://github.com/kubernetes/autoscaler/pull/6242/files)

对于很多挂载目录持久化数据的，例如各种中间件，例如 mysql，redis ，单纯设置 USER 的话，需要在容器启动之前设置目录的权限。other 权限为 7 的话，很不安全，所以只能是 owner、group 权限，但是容器内的用户名和宿主机用户名是不一致的，只能设置 uid、gid。使用这些需要数据持久化的容器，会存在：

- 直接 -v 挂载或者 docker volume
- k8s 上使用 hostPath
- 固定 pv
- sc 下使用 pvc
- 别人的 k8s 集群或者实例上去部署

如果你提前修改目录权限，上面最后俩场景根本无法自动化，而且说不定某天新版本官方镜像里 Dockerfile 里换基础镜像的同时忘记在添加用户时候设置 uid 和 gid ，uid 和 gid 就变了，只能是加启动脚本里处理。

对此，[mysql docker 镜像的官方启动脚本](https://github.com/docker-library/mysql/blob/master/5.7/docker-entrypoint.sh) 给了很好的参考，Dockerfile 制作镜像就创建了指定 uid、gid 的 mysql 用户，然后启动容器的时候都是 `ENTRYPOINT CMD` （k8s 里对应 command、args） 的形式启动：

```plaintext
docker-entrypoint.sh mysqld
```

或者可以通过 cmdline 设置 mysql 启动端口

```plaintext
docker run xxx mysql:5.7 --port 4306
```

mysql 脚本里包含对于权限以外的信息比较多，不方便举例，这里使用 redis 举例：

```bash
#!/bin/sh

set -e

if [ "${1#-}" != "$1" ] || [ "${1%.conf}" != "$1" ]; then

 set -- redis-server "$@"
fi



if [ "$1" = 'redis-server' -a "$(id -u)" = '0' ]; then

 find . \! -user redis -exec chown redis '{}' +

 exec gosu redis "$0" "$@"
fi




um="$(umask)"
if [ "$um" = '0022' ]; then
 umask 0077
fi

exec "$@"
```

例如下面执行流程：

```bash
$ docker run -d -name redis7 -v $PWD/redis-ctr-data:/data --net host redis:7 --port 7777
$ docker top redis7
UID                 PID                 PPID                C                   STIME               TTY                 TIME                CMD
systemd+            1041135             1041116             1                   15:47               ?                   00:00:00            redis-server *:7777
$ docker exec redis7 id redis
uid=999(redis) gid=999(redis) groups=999(redis)
$ grep 999 /etc/passwd
systemd-coredump:x:999:999:systemd Core Dumper:/:/usr/sbin/nologin
```

docker top 显示的用户，是按照宿主机上 uid 显示的，[gosu](https://github.com/tianon/gosu) 是 golang 实现 [su-exec](https://github.com/ncopa/su-exec)，切换指定用户执行命令，exec 是执行后面的命令，替换当前的 shell 进程，这样在 docker stop 给容器内 pid 为 1 的进程发送信号，业务进程能收到信号进行优雅退出，而没 exec 的话，pid 为 1 的进程是 shell 脚本，它不会转发信号的。

`ENTRYPOINT` 使用脚本当作入口的形式，最后业务切用户执行，即使使用 docker exec 还是使用镜像默认的 USER root，排查问题也方便。 也推荐使用镜像之前，先看官方的启动脚本，例如 mongodb 官方镜像是支持类似 redis 这种非 root 启动的，但是我们 k8s 里是：

```plaintext
...
    - name: {{ NODE_NAME }}
      image: xxx/mongo:xxx
      command:
        - mongod
        - "--port"
```

这样覆盖了 entrypoint，没有使用官方启动脚本执行，就是 root 用户，改为下面的不覆盖就行：

```plaintext
- name: {{ NODE_NAME }}
  image: xxx/mongo:xxx
  args: # <--- 这里
    - mongod
    - "--port"
```

要注意一个点，su-exec 在 alpine 里可以包管理安装，非 alpine 的基础镜像使用 gosu 可以参考 redis 官方镜像，以及 su-exec 不是静态编译的，可能某些系统上有问题，自行测试下看看

### [](#案例实战 "案例实战")案例实战[](#案例实战)

这列梳理一些我做的案例。先说一些知识点：

- 产生 pid 和 sock 文件的，可以放 /tmp 下
- 业务进程非 root 对 `/dev/stdxxx` 没权限的，可以脚本里 `chmod a+w /dev/std*`
- 如果自己业务镜像产生的数据会被其他容器挂载操作数据，你的业务进程最好创建用户的时候使用固定同样的 `uid:gid` ，例如我们的 mysql-backup 备份 mysql 数据用到的用户 `uid:gid` 保持和 mysql 官方镜像一致，这样不需要修改 mysql 数据目录权限和 owner
- 不要 `chmod -R 777` 目录

#### [](#机器码处理 "机器码处理")机器码处理[](#机器码处理)

获取机器码一般是使用 `dmidecode -s system-uuid` ，但是容器内你以 root 执行会报错：

```plaintext
$ docker run --rm -ti debian:11
$ apt update && apt-get install -y dmidecode
$ dmidecode -s system-uuid
/dev/mem: No such file or directory
```

所以之前我们都是读取 `/sys/devices/virtual/dmi/id/product_uuid`，但是非 root 后无法读取，因为该文件权限为 `0400`:

```plaintext
$ ls -l /sys/devices/virtual/dmi/id/product_uuid
-r-------- 1 root root 4096 Nov  3 08:48 /sys/devices/virtual/dmi/id/product_uuid
```

且该文件是[内核设置的权限](https://github.com/torvalds/linux/blob/master/drivers/firmware/dmi-id.c#L61)，无法被更改。

后面尝试发现一些信息:

```plaintext
$ strace dmidecode -s system-uuid
...
openat(AT_FDCWD, "/sys/firmware/dmi/tables/smbios_entry_point", O_RDONLY)
...
openat(AT_FDCWD, "/sys/firmware/dmi/tables/DMI", O_RDONLY)
```

发现读取了这俩文件，搜索资料发现是 dmi table，例如 root 下可以这样获取机器码：

```plaintext
dmidecode -t 1  < /sys/firmware/dmi/tables/DMI
dmidecode -t 1 -u < /sys/firmware/dmi/tables/DMI
```

该文件内容按照 DMI 规范字节结构解析可以得到不少信息。然后找到了一个 go 库，在 linux 上尝试成功：

```golang
package main

import (
 "fmt"
 "log"

 "github.com/digitalocean/go-smbios/smbios"
)

func main() {

 rc, ep, err := smbios.Stream()
 if err != nil {
  log.Fatalf("failed to open stream: %v", err)
 }

 defer rc.Close()


 d := smbios.NewDecoder(rc)
 ss, err := d.Decode()
 if err != nil {
  log.Fatalf("failed to decode structures: %v", err)
 }

 major, minor, _ := ep.Version()

 for _, s := range ss {
  if s.Header.Type == 1 {
   d := s.Formatted

   if major > 0x02 || (major == 0x02 && minor >= 0x06) {
    fmt.Printf("UUID: %02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X\n",
     d[7], d[6], d[5], d[4],
     d[9], d[8], d[11], d[10], d[12], d[13],
     d[14], d[15], d[16], d[17], d[18], d[19],
    )
   } else {
    fmt.Printf("UUID: %02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X\n",
     d[4], d[5], d[6], d[7],
     d[8], d[9], d[10], d[11], d[12], d[13],
     d[14], d[15], d[16], d[17], d[18], d[19],
    )
   }
  }
 }
}
```

机器上测试：

```bash
$ dmidecode -s system-uuid | tr a-z A-Z
66C0F667-71A0-xxxx-xxxx-4AC0A21F5428
$ go build -o /tmp/uuid-go test.go
$ chmod a+r /sys/firmware/dmi/tables/DMI
$ su - guanzhang
guanzhang@guan:~$ /tmp/uuid-go
UUID: 66C0F667-71A0-xxxx-xxxx-4AC0A21F5428
```

然后把宿主机的 `/sys/firmware/dmi/tables` 挂载到 `/rootfs/sys/firmware/dmi/tables` 里，在 gosu 之前 `chmod a+r /rootfs/sys/firmware/dmi/tables/DMI`，业务使用上面的库 hack 后，从指定路径的 DMI 信息即可获取到机器码。

#### [](#etcd "etcd")etcd[](#etcd)

没啥说的，加了 gosu 后再加启动脚本：

```bash
#!/bin/bash

set -e

if [ "${1:0:1}" = '-' ]; then
 set -- etcd "$@"
fi


if [ "$1" = 'etcd' ] || [ "$1" = '/usr/local/bin/etcd' ];then
    if [ "$(id -u)" = '0' -a -n "$RUN_USER" ]; then
     find /var/lib/etcd \! -user ${RUN_USER} -exec chown ${RUN_USER} '{}' +
     exec gosu ${RUN_USER} "$@"
    fi
fi

exec "$@"
```

为了不影响其他分支，这里我用了 env 作为开关，[wurstmeister/kafka-docker](https://github.com/wurstmeister/kafka-docker) 也是一样：

```bash
#!/bin/bash

set -e

if [ "${1:0:1}" = '-' ]; then
 set -- start-kafka.sh "$@"
fi


if [ "$1" = 'start-kafka.sh' ] || [ "$1" = '/usr/bin/start-kafka.sh' ];then
    if [ "$(id -u)" = '0' -a -n "$RUN_USER" ]; then
  find $(readlink -f ${KAFKA_HOME}) \! -user ${RUN_USER} -exec chown ${RUN_USER} '{}' +
  find /kafka \! -user ${RUN_USER} -exec chown ${RUN_USER} '{}' +
     exec gosu ${RUN_USER} "$@"
    fi
fi

exec "$@"
```

其他的很多都是类似这样，不再举例，自行制作

#### [](#coredns "coredns")coredns[](#coredns)

coredns 1.11.0 才开始非 root 启动，我们业务使用的是 1.10.1 的，不升级避免客户现场出现问题，所以重做镜像最稳妥：

```dockerfile
ARG DEBIAN_IMAGE=debian:stable-slim
ARG BASE=gcr.io/distroless/static-debian12:nonroot
FROM coredns/coredns:1.10.1 as bin

FROM  ${DEBIAN_IMAGE} AS build
SHELL [ "/bin/sh", "-ec" ]

RUN export DEBCONF_NONINTERACTIVE_SEEN=true \
           DEBIAN_FRONTEND=noninteractive \
           DEBIAN_PRIORITY=critical \
           TERM=linux ; \
    apt-get -qq update ; \
    apt-get -yyqq upgrade ; \
    apt-get -yyqq install ca-certificates libcap2-bin; \
    apt-get clean
COPY --from=bin /coredns /coredns
RUN setcap cap_net_bind_service=+ep /coredns

FROM  ${BASE}
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /coredns /coredns
USER nonroot:nonroot
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
```

非 root 用户是无法监听 1024 以下端口的，coredns 监听 53 端口是因为使用了 `setcap cap_net_bind_service=+ep /coredns`，但是这个属性属于扩展属性，docker 构建多层 COPY 会不支持而丢失，必须使用 buildkit 构建，否则 cap 信息丢失，部署上去无法监听 53 端口：

```plaintext
DOCKER_BUILDKIT=1 docker build --platform=amd64  . -t coredns/coredns:1.10.1  --load
```

#### [](#consul "consul")consul[](#consul)

consul 镜像也支持，但是 chown 的时候没带 -R 选项。

```plaintext
if [ "$(stat -c %u "$CONSUL_DATA_DIR")" != "${CONSUL_UID}" ]; then
  chown ${CONSUL_UID}:${CONSUL_GID} "$CONSUL_DATA_DIR"
fi
```

这里会存在一个问题，如果之前是覆盖了 entrypoint 使用 root 启动的，再切正确姿势下，因为 data 目录下子目录没被 chown，consul 在 data 下子目录写入 node-id 会报错没权限，所以我是这样 hack 重做镜像的：

```dockerfile
ARG VER=1.8.3
FROM consul:${VER}
RUN sed -ri -e 's/(chown)(\s+consul:)/\1 -R\2/' \
        -e '1s@/usr/bin/dumb-init\s+@@' \
    /usr/local/bin/docker-entrypoint.sh
```

去掉 `dumb-init` 是因为客户要求容器内所有进程都是非 root，不去掉 pid 为 1 的就是 root 用户 dumb-init sh 进程

#### [](#docker-sock-文件 "docker.sock 文件")docker.sock 文件[](#docker-sock-文件)

有些进程是需要挂载 `/var/run` 为了使用宿主机的 `/var/run/docker.sock` 和宿主机 docker 通信的，这里我们使用 cadvisor 举例：

```dockerfile
ARG VER=v0.37.5
FROM gcr.m.daocloud.io/cadvisor/cadvisor:${VER}
ARG GO_SU=1.17
RUN set -eux; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    apk update; \
    apk add --no-cache \
      curl \
    curl -o /usr/local/bin/gosu -sSL https://github.com/tianon/gosu/releases/download/${GO_SU}/gosu-amd64; \
    chmod a+x /usr/local/bin/gosu; \
    gosu --version; \
    rm -rf /var/cache/apk/* /tmp/*
COPY docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["cadvisor", "-logtostderr"]
```

```bash
#!/bin/sh
set -e

if [ "${1:0:1}" = '-' ]; then
 set -- cadvisor "$@"
fi

if [ "$1" = 'cadvisor' ] || [ "$1" = '/usr/bin/cadvisor' ];then
    if [ "$(id -u)" = '0' -a -n "$RUN_USER" ]; then
        if [ -S /var/run/docker.sock ];then
            group_id=`stat -c "%g" /var/run/docker.sock`
            if ! getent group | cut -d: -f3 | grep -wq $group_id; then
                if ! addgroup -g ${group_id} docker;then

                    group_failed=true
                fi
            fi
            if [ -z "$group_failed" ];then
                group_name=$(stat -c "%G" /var/run/docker.sock)
                if ! id -nG ${RUN_USER} | grep -w ${group_name};then

                    adduser ${RUN_USER} ${group_name}
                fi
            else

                setfacl -m u:${RUN_USER}:rw /var/run/docker.sock
            fi

        fi

        exec gosu $RUN_USER $@
    fi
fi

exec $@
```

- cadvisor 挂载了宿主机的 rootfs ，改为纯非 root 不行，但是 cadvisor 镜像内有个 `operator` 用户的 gid 是 0，利用启动脚本和 docker 权限来改造成非 root 启动。
- docker.sock 权限是 `0660`，利用 shell 把 operator 用户加到 docker 组里即可（必须取 gid）。这里要注意的是，不同版本 alpine 和其他 rootfs 的 adduser/addgroup 参数不一样，自行注意 shell 兼容

设置 “RUN_USER” 为 `operator` ，然后设置宿主机的 docker 的 data-root 下面权限（可以使用 systemd 的`ExecStartPost=`）：

```plaintext
/var/lib/docker/image：750   ok
/var/lib/docker/image/overlay2：750 ok
/var/lib/docker/image/overlay2/layerdb：750 ok
```

cadvisor 参数为：

```yaml
...
      args:
      - -docker_only=true
      - -housekeeping_interval=20s
      - -disable_metrics=accelerator,cpu_topology,tcp,udp,percpu,sched,process,hugetlb,referenced_memory,resctrl
```

promtail 官方在 [Promtail should run as non-root in docker](https://github.com/grafana/loki/issues/1813#issuecomment-694969497) 推荐在 k8s 里设置安全上下文，有少部分老哥指定的 gid 为 0 可以运行，但是还是很多人有权限问题，所以应该也要像上面这样，用 gid 为 0 和配合脚本设置权限。

#### [](#cron "cron")cron[](#cron)

非 root 无法使用 cron 启动，使用 [go-crond](https://github.com/webdevops/go-crond)

```plaintext
exec gosu  user1 go-crond   --default-user=user1  --include=/etc/cron.d --allow-unprivileged
```

# 参考

- [k8s 社区关于支持 user namespace 提议](https://github.com/kubernetes/enhancements/issues/127)
- [dmi 信息规范](https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.3.0.pdf)
- [dmidecode 源码](https://github.com/mirror/dmidecode/blob/master/dmidecode.c#L448)
