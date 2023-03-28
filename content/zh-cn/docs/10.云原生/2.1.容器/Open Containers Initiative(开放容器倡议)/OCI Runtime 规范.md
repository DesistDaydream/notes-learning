---
title: OCI Runtime 规范
---

# 概述

> 参考：
>
> - [GitHub 项目,opencontainers/runtime-spec/spec.md](https://github.com/opencontainers/runtime-spec/blob/main/spec.md)
> - [GitHub 项目,opencontainers/runtime-tools](https://github.com/opencontainers/runtime-tools)
> - [思否大佬](https://segmentfault.com/a/1190000009583199)

OCI Runtime 规范用来指定一个 Container 的配置、执行环境和生命周期。

容器的配置被指定为 config.json ，并详细说明了可以创建容器的字段。指定执行环境是为了确保容器内运行的应用程序在运行时之间具有一致的环境，以及为容器的生命周期定义的常见操作。

由于容器运行起来，需要一个运行环境，比如是运行在 linux 上、还是 windows 上~~所以，OCI Runtime 标准，会根据不同的平台，制定不同的规范。现阶段有 4 中平台规范。这点是根 OCI Image 规范不太一样的地方。

- linux：[runtime.md](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md)，[config.md](https://github.com/opencontainers/runtime-spec/blob/master/config.md)，[config-linux.md](https://github.com/opencontainers/runtime-spec/blob/master/config-linux.md)和[runtime-linux.md](https://github.com/opencontainers/runtime-spec/blob/master/runtime-linux.md)。
- solaris：[runtime.md](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md)，[config.md](https://github.com/opencontainers/runtime-spec/blob/master/config.md)和[config-solaris.md](https://github.com/opencontainers/runtime-spec/blob/master/config-solaris.md)。
- windows：[runtime.md](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md)，[config.md](https://github.com/opencontainers/runtime-spec/blob/master/config.md)和[config-windows.md](https://github.com/opencontainers/runtime-spec/blob/master/config-windows.md)。
- vm：[runtime.md](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md)，[config.md](https://github.com/opencontainers/runtime-spec/blob/master/config.md)和[config-vm.md](https://github.com/opencontainers/runtime-spec/blob/master/config-vm.md)。

由于我们日常使用 linux，所以下面就只研究 linux 平台的 OCI Runtime 规范

runtime 规范有如下几个，所有人必须遵守该规范来使用 runtime ：

- [Filesystem Bundle](https://github.com/opencontainers/runtime-spec/blob/master/bundle.md) # 文件系统捆绑。bundle 是以某种方式组织的一组文件，包含了容器所需要的所有信息，有了这个 bundle 后，符合 runtime 标准的程序(e.g.runc)就可以根据 bundle 启动容器了(哪怕没有 docker，也可以启动一个容器)。
- [Runtime and Lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md) # 使用一个 runtime 创建的容器实体必须能够对同一容器使用本规范中定义的操作。
  - [Linux-specific Runtime and Lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime-linux.md) # 这是关于 linux 平台的 Runtime 与 Lifecycle
- [Configuration](https://github.com/opencontainers/runtime-spec/blob/master/config.md) # Configuration 包含对容器执行标准操作(比如 create、start、stop 等)所必须的元数据。这包括要运行的过程、要注入的环境变量、要使用的沙盒功能等等。不同平台(linux、window 等)，有不同的规范。
  - [Linux-specific configuration](https://github.com/opencontainers/runtime-spec/blob/master/config-linux.md) # 这是关于 linux 平台的 Configuration

# Filesystem Bundle

官方详解：<https://github.com/opencontainers/runtime-spec/blob/master/bundle.md>

bundle 中包含了运行容器所需要的所有信息，有了这个 bundle 后，符合 runtime 标准的程序（比如 runc）就可以根据 bundle 启动容器了。

一个容器的 Bundle 必须包含以下内容：

- **config.json** # 此文件包含配置数据。该文件必须留在 bundle 目录的根目录中，并且必须名为 config.json。config.json 内的写法参考下文 [Configuration 章节](#AGQoL)。
- **rootfs** # 容器的根文件系统目录，不过 bundle 对 rootfs 没有要求，名字也可以随便改，只要在 config.json 文件中的 `.root.path` 字段配置好就可以。

## 应用示例

rootfs OCI 镜像规范中，blobs 目录下的镜像层文件。config.json 可以通过 OCI 官方提供的 [runtime-tools](https://github.com/opencontainers/runtime-tools) 工具生成，现在我们操作一下:

> 也可以使用 `runc spec` 命令生成 config.json 文件
> 这里接着 [OCI Image 规范中的实验](/docs/10.云原生/2.1.容器/Open%20Containers%20Initiative(开放容器倡议)/OCI%20Image%20 规范.md Image 规范.md)中的 [Layers 文件](/docs/10.云原生/2.1.容器/Open%20Containers%20Initiative(开放容器倡议)/OCI%20Image%20 规范.md Image 规范.md)章节，使用 lchdzh/k8s-debug 镜像。

```bash
~]# cd /root/test_dir/k8s-debug/layers
~]# oci-runtime-tool generate --output config.json
```

> config.json 文件中的 root.path 字段用来指定 rootfs 路径；process.args 字段用来指定运行容器时所使用的命令

```json
{
    ......
 "process": {
  "args": [
   "sh"
  ],
 "root": {
  "path": "rootfs"
 },
    ......
}
```

此时我们将其中一个镜像层目录改名，改为 rootfs，省的改 config.json 文件了~~~:D

验证一下，可以发现通过了

```bash
$ oci-runtime-tool validate
Bundle validation succeeded.
```

这时，使用 OCI 中的 runc 工具，即可运行出来一个容器

```bash
~/test_dir/k8s-debug/layers]# runc run hello
cat /etc/os-release
NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.12.1
PRETTY_NAME="Alpine Linux v3.12"
HOME_URL="https://alpinelinux.org/"
BUG_REPORT_URL="https://bugs.alpinelinux.org/"

ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
```

可以看到，直接运行的时候，虽然没有 shell 提示符，但是各种名称空间已经是单独的了，这就是通过 alpine:latest 镜像启动的容器。并且没有任何网络设备接入。

### rootfs 中只有一个文件的话

任何目录都可以作为容器的根文件系统并启动，比如我们直接使用 `COPY --from=grpcurl  /go/bin/grpcurl /usr/bin/grpcurl` 镜像层作为容器的根文件系统

```bash
$ mv rootfs 2
$ mv 1 config.json
# 这里注意，由于这个镜像层只有一个 grpcurl 命令
# 所以无法执行 Bundle 中 config.conf 中 .process.args 字段中指定的 sh 命令
# 所以 .process.args 字段的值需要修改为 /usr/bin/grpcurl
$ cat config.json  |grep grpc
   "/usr/bin/grpcurl"
```

此时运行一下看看，可以看到可以正常运行一个容器，只不过这个容器是瞬时的而已。

```bash
~/test_dir/k8s-debug/layers]# runc run hello
Too few arguments.
Try '/usr/bin/grpcurl -help' for more details.
```

## docker 中的 Filesystem Bundle

下面是在 docker 环境下 bundle 位置的示例：

```shell
[root@lichenhao containerd]# pwd
/run/containerd
[root@lichenhao containerd]# tree
.
├── containerd.sock
├── io.containerd.runtime.v1.linux
│   └── moby
│       └── 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
│           ├── config.json
│           ├── init.pid
│           ├── log.json
│           └── rootfs
└── io.containerd.runtime.v2.task
```

而 rootfs 实际在另外一个目录，通过 config.json 中可以看到

```shell
[root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# pwd
/run/containerd/io.containerd.runtime.v1.linux/moby/28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10
[root@lichenhao 28f5bed704dc80bed6dbaa8af514d2191d8d4ab0339bb3a663e66609ccd34c10]# cat config.json | jq . | more
.....
  "root": {
    "path": "/var/lib/docker/overlay2/d976eddf7575a3464486d92539229146f3df66080a3265195791ebb0d24b24dd/merged"
  },
.....
```

# Runtime and Lifecycle

官方详解：<https://github.com/opencontainers/runtime-spec/blob/master/runtime.md>

Runtime 与 Lifecycle 规范部分，定义了跟容器 runtime 相关的三部分内容

1. Container 的状态
2. Container 的生命周期
3. 对使用一个 runtime 创建出来的 container 实体，可以实现的操作。这些操作包括增删改查等等。

## Container State(容器的状态) 规范

- ociVersion (string, REQUIRED) is version of the Open Container Initiative Runtime Specification with which the state complies.
- id (string, REQUIRED) is the container's ID. This MUST be unique across all containers on this host. There is no requirement that it be unique across hosts.
- status (string, REQUIRED) is the runtime state of the container. The value MAY be one of:
  - creating: the container is being created (step 2 in the [lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#lifecycle))
  - created: the runtime has finished the [create operation](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#create) (after step 2 in the [lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#lifecycle)), and the container process has neither exited nor executed the user-specified program
  - running: the container process has executed the user-specified program but has not exited (after step 5 in the [lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#lifecycle))
  - stopped: the container process has exited (step 7 in the [lifecycle](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#lifecycle))

Additional values MAY be defined by the runtime, however, they MUST be used to represent new runtime states not defined above.

- pid (int, REQUIRED when status is created or running on Linux, OPTIONAL on other platforms) is the ID of the container process. For hooks executed in the runtime namespace, it is the pid as seen by the runtime. For hooks executed in the container namespace, it is the pid as seen by the container.
- bundle (string, REQUIRED) is the absolute path to the container's bundle directory. This is provided so that consumers can find the container's configuration and root filesystem on the host.
- annotations (map, OPTIONAL) contains the list of annotations associated with the container. If no annotations were provided then this property MAY either be absent or an empty map.

除了上述属性外，还可以包含其他属性。有关获取容器状态的信息的方式，请参考下文 Query State 段落。

## Lifecycle(生命周期) 规范

生命周期描述了从创建容器到容器不再存在发生的事件的时间轴。

- OCI compliant runtime's [create](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#create) command is invoked with a reference to the location of the bundle and a unique identifier.
- The container's runtime environment MUST be created according to the configuration in [config.json](https://github.com/opencontainers/runtime-spec/blob/master/config.md). If the runtime is unable to create the environment specified in the [config.json](https://github.com/opencontainers/runtime-spec/blob/master/config.md), it MUST [generate an error](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#errors). While the resources requested in the [config.json](https://github.com/opencontainers/runtime-spec/blob/master/config.md) MUST be created, the user-specified program (from [process](https://github.com/opencontainers/runtime-spec/blob/master/config.md#process)) MUST NOT be run at this time. Any updates to [config.json](https://github.com/opencontainers/runtime-spec/blob/master/config.md) after this step MUST NOT affect the container.
- The [prestart hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#prestart) MUST be invoked by the runtime. If any prestart hook fails, the runtime MUST [generate an error](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#errors), stop the container, and continue the lifecycle at step 12.
- The [createRuntime hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#createRuntime-hooks) MUST be invoked by the runtime. If any createRuntime hook fails, the runtime MUST [generate an error](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#errors), stop the container, and continue the lifecycle at step 12.
- The [createContainer hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#createContainer-hooks) MUST be invoked by the runtime. If any createContainer hook fails, the runtime MUST [generate an error](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#errors), stop the container, and continue the lifecycle at step 12.
- Runtime's [start](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#start) command is invoked with the unique identifier of the container.
- The [startContainer hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#startContainer-hooks) MUST be invoked by the runtime. If any startContainer hook fails, the runtime MUST [generate an error](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#errors), stop the container, and continue the lifecycle at step 12.
- The runtime MUST run the user-specified program, as specified by [process](https://github.com/opencontainers/runtime-spec/blob/master/config.md#process).
- The [poststart hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#poststart) MUST be invoked by the runtime. If any poststart hook fails, the runtime MUST [log a warning](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#warnings), but the remaining hooks and lifecycle continue as if the hook had succeeded.
- The container process exits. This MAY happen due to erroring out, exiting, crashing or the runtime's [kill](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#kill) operation being invoked.
- Runtime's [delete](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#delete) command is invoked with the unique identifier of the container.
- The container MUST be destroyed by undoing the steps performed during create phase (step 2).
- The [poststop hooks](https://github.com/opencontainers/runtime-spec/blob/master/config.md#poststop) MUST be invoked by the runtime. If any poststop hook fails, the runtime MUST [log a warning](https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#warnings), but the remaining hooks and lifecycle continue as if the hook had succeeded.

这里以 runc 为例，说明容器的生命周期

- 执行命令 runc create 创建容器，参数中指定 bundle 的位置以及容器的 ID，容器的状态变为 creating
- runc 根据 bundle 中的 config.json，准备好容器运行时需要的环境和资源，但不运行 process 中指定的进程，这步执行完成之后，表示容器创建成功，修改 config.json 将不再对创建的容器产生影响，这时容器的状态变成 created。
- 执行命令 runc start 启动容器
- runc 执行 config.json 中配置的 prestart 钩子
- runc 执行 config.json 中 process 指定的程序，这时容器状态变成了 running
- runc 执行 poststart 钩子。
- 容器由于某些原因退出，比如容器中的第一个进程主动退出，挂掉或者被 kill 掉等。这时容器状态变成了 stoped
- 执行命令 runc delete 删除容器，这时 runc 就会删除掉上面第 2 步所做的所有工作。
- runc 执行 poststop 钩子

**Errors 与 Warnings 规范**

## Standard Operations(标准操作) 规范，用来规范容器必须支持的操作

除非另有说明，否则 runtime 必须支持以下操作(注意：这些操作未指定任何命令行 API，并且参数是常规操作的输入。)

- **Query State(查询状态)** # 返回容器的状态，包含上面介绍的那些内容.
- **Create(创建容器**) # 创建容器，这一步执行完成后，容器创建完成，修改 bundle 中的 config.json 将不再对已创建的容器产生影响
- **Start(启动容器)** # 启动容器，执行 config.json 中 process 部分指定的进程
- **Kill(停止容器)** # 通过给容器发送信号来停止容器，信号的内容由 kill 命令的参数指定
- **Delete(删除容器)** # 删除容器，如果容器正在运行中，则删除失败。删除操作会删除掉 create 操作时创建的所有内容。

# Configuration

官方详解：<https://github.com/opencontainers/runtime-spec/blob/master/config-linux.md>

Configuration 规范定义了《Filesystem Bundle 规范》章节中，config.json 文件中应该包含哪些内容

Configuration 中包含对容器实施 standard operations(标准操作内容见上文) 所必须的元数据。这包括要运行的过程、要注入的环境变量、要使用的沙盒功能等等。

config.json 文件的样例，可以参考官方：<https://github.com/opencontainers/runtime-spec/blob/master/config.md#configuration-schema-example>
