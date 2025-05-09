---
title: "services"
linkTitle: "services"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考 - Compose 文件参考 - Services 顶级元素](https://docs.docker.com/reference/compose-file/services/)


## build

指定为构建镜像上下文路径：

例如 webapp 服务，指定为从上下文路径 ./dir/Dockerfile 所构建的镜像：

```yaml
version: "3"
services:
  webapp:
    build: ./dir
```

或者，作为具有在上下文指定的路径的对象，以及可选的 Dockerfile 和 args：

```yaml
version: "3"
services:
  webapp:
    build:
      context: ./dir
      dockerfile: Dockerfile-alternate
      args:
        buildno: 1
      labels:
        - "com.example.description=Accounting webapp"
        - "com.example.department=Finance"
        - "com.example.label-with-empty-value"
      target: prod
```

- context：上下文路径。
- dockerfile：指定构建镜像的 Dockerfile 文件名。
- args：添加构建参数，这是只能在构建过程中访问的环境变量。
- labels：设置构建镜像的标签。
- target：多层构建，可以指定构建哪一层。

## cap_add 与 cap_drop

添加或删除容器拥有的宿主机的内核功能。等价于 [docker run 命令中的的 --cap-add 标志](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20CLI/run.md#特权%20与%20Linux%20Capabilities)

```yaml
cap_add:
  - ALL # 开启全部权限
cap_drop:
  - SYS_PTRACE # 关闭 ptrace权限
```

## cgroup_parent

为容器指定父 cgroup 组，意味着将继承该组的资源限制

```yaml
cgroup_parent: m-executor-abcd
```

## command

覆盖容器启动的默认命令

```yaml
command: ["bundle", "exec", "thin", "-p", "3000"]
```

## container_name

指定自定义容器名称，而不是生成的默认名称。

```yaml
container_name: my-web-container
```

## depends_on

设置依赖关系。

- docker-compose up ：以依赖性顺序启动服务。在以下示例中，先启动 db 和 redis ，才会启动 web。
- docker-compose up SERVICE ：自动包含 SERVICE 的依赖项。在以下示例中，docker-compose up web 还将创建并启动 db 和 redis。
- docker-compose stop ：按依赖关系顺序停止服务。在以下示例中，web 在 db 和 redis 之前停止。

```yaml
version: "3"
services:
  web:
    build: .
    depends_on:
      - db
      - redis
  redis:
    image: redis
  db:
    image: postgres
```

注意：web 服务不会等待 redis db 完全启动 之后才启动。

## deploy

该字段有独立的规范，可以定义一些 services 的部署方式和生命周期。e.g. 如何使用 GPU、etc.

https://docs.docker.com/reference/compose-file/deploy/

## devices

指定设备映射列表。

```yaml
devices:
  - "/dev/ttyUSB0:/dev/ttyUSB0"
```

## dns

自定义 DNS 服务器，可以是单个值或列表的多个值。

```yaml
dns: 8.8.8.8
dns:
  - 8.8.8.8
  - 9.9.9.9
```

## dns_search

自定义 DNS 搜索域。可以是单个值或列表。

```yaml
dns_search: example.com
dns_search:
  - dc1.example.com
  - dc2.example.com
```

## entrypoint

覆盖容器默认的 entrypoint

```yaml
entrypoint: /code/entrypoint.sh
```

也可以是以下格式：

```yaml
entrypoint:
    - php
    - -d
    - zend_extension=/usr/local/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so
    - -d
    - memory_limit=-1
    - vendor/bin/phpunit
```

## env_file

从文件添加环境变量。可以是单个值或列表的多个值。

```yaml
env_file: .env
```

也可以是列表格式：

```yaml
env_file:
  - ./common.env
  - ./apps/web.env
  - /opt/secrets.env
```

## environment

添加环境变量。您可以使用数组或字典、任何布尔值，布尔值需要用引号引起来，以确保 YML 解析器不会将其转换为 True 或 False。

```yaml
environment:
  RACK_ENV: development
  SHOW: 'true'
```

## expose

暴露端口，但不映射到宿主机，只被连接的服务访问。

仅可以指定内部端口为参数：

```yaml
expose:
 - "3000"
 - "8000"
```

## extra_hosts

添加主机名映射。类似 docker client --add-host。

```yaml
extra_hosts:
 - "somehost:162.242.195.82"
 - "otherhost:50.31.209.229"
```

以上会在此服务的内部容器中 /etc/hosts 创建一个具有 ip 地址和主机名的映射关系：

```yaml
162.242.195.82  somehost
50.31.209.229   otherhost
```

## healthcheck

用于检测 docker 服务是否健康运行。

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost"] # 设置检测程序
  interval: 1m30s # 设置检测间隔
  timeout: 10s # 设置检测超时时间
  retries: 3 # 设置重试次数
  start_period: 40s # 启动后，多少秒开始启动检测程序
```

## image

指定容器运行的镜像。以下格式都可以：

```yaml
image: redis
image: ubuntu:14.04
image: tutum/influxdb
image: example-registry.com:4000/postgresql
image: a4bc65fd # 镜像id
```

## logging

服务的日志记录配置。
driver：指定服务容器的日志记录驱动程序，默认值为json-file。有以下三个选项

```yaml
driver: "json-file"
driver: "syslog"
driver: "none"
```

仅在 json-file 驱动程序下，可以使用以下参数，限制日志得数量和大小。

```yaml
logging:
  driver: json-file
  options:
    max-size: "200k" # 单个文件大小为200k
    max-file: "10" # 最多10个文件
```

当达到文件限制上限，会自动删除旧得文件。

syslog 驱动程序下，可以使用 syslog-address 指定日志接收地址。

```yaml
logging:
  driver: syslog
  options:
    syslog-address: "tcp://192.168.0.42:123"
```

## network_mode

设置容器的网络模式，可用的值有如下几种

- **host** # 使用宿主机网络。i.e. 让容器加入 1 号进程的网络名称空间
- **none** # 关闭所有容器网络。
- **service:${ServiceName}** # 让该容器加入其他容器的网络，让两个容器共享 Network namespace。
  - Notes: ServiceName 就是顶层字段 services 的下级字段的名称
  - 关于容器网络更详细的内容详见 [Docker Network](/docs/10.云原生/Containerization%20implementation/Docker/Docker%20Network.md)

> [!Warning]
> network_mode 与 [networks](#networks) 字段互斥，若使用了 networks 字段，则相当于之前老版本将 network_mode 设置为 bridge
>
> TODO: 使用 `service:${ServiceName}` 模式的场景下，若被加入的容器若重启后，加入者无法自动重新加入

e.g. 下面的例子中，other-service 会加入 some-service 的 Network namespace

```yaml
services:
  some-service:
    networks:
      - some-network
  other-service:
    network_mode: service:some-service
```

## networks

配置容器连接的网络，引用顶级 networks 下的条目 。networks 节点有两种语法

> [!Attention] 配置该字段后，相当于让该容器使用 Bridge 驱动的网络。在老版本相当于将 network_mode 设置为 bridge

```yaml
services:
  some-service:
    networks:
      - some-network
      - other-network
```

**aliases**(\[]STRING) # 网络别名。指定的 aliases 会记录到 Docker 内嵌的 DNS Server 中，可以通过别名访问容器。（注意，要想使用内嵌 DNS 需要连接到自定义网络而不是 Docker 的默认网络）

```yaml
services:
  some-service: # 自定义 service 名称
    networks:
      some-network: # 让 service 连接到 some-network 网络
        aliases:
         - alias1 # 在 some-network 网络中的别名
      other-network:
        aliases:
         - alias2
networks:
  some-network:
    # Use a custom driver
    driver: custom-driver-1
  other-network:
    # Use a custom driver which takes special options
    driver: custom-driver-2
```

> [!Note] network 节点下，使用 aliases 与不使用 aliases 的语法是不一样的

## ports

配置端口映射，有 3 种语法

- `HostPort:ContainerPort` # HOST 为主机上的端口，CONTAINER 为容器内的端口。
- `ContainerPort` # 仅指定容器内的端口，主机端口随机选择
- `HostIP:HostPort:ContainerPort` # 将容器内的端口映射到主机上指定 IP 的端口上。

## restart

- no：是默认的重启策略，在任何情况下都不会重启容器。
- always：容器总是重新启动。
- on-failure：在容器非正常退出时（退出状态非0），才会重启容器。
- unless-stopped：在容器退出时总是重启容器，但是不考虑在 Docker 守护进程启动时就已经停止了的容器

```yaml
restart: "no"
restart: always
restart: on-failure
restart: unless-stopped
```

Notes：Swarm 集群模式，请改用 restart_policy。

## secrets

存储敏感数据，例如密码：

```yaml
version: "3.1"
services:
mysql:
  image: mysql
  environment:
    MYSQL_ROOT_PASSWORD_FILE: /run/secrets/my_secret
  secrets:
    - my_secret
secrets:
  my_secret:
    file: ./my_secret.txt
```

## security_opt

修改容器默认的 schema 标签。

```yaml
security-opt：
  - label:user:USER   # 设置容器的用户标签
  - label:role:ROLE   # 设置容器的角色标签
  - label:type:TYPE   # 设置容器的安全策略标签
  - label:level:LEVEL  # 设置容器的安全等级标签
```

## stop_grace_period

指定在容器无法处理 SIGTERM (或者任何 stop_signal 的信号)，等待多久后发送 SIGKILL 信号关闭容器。

```yaml
stop_grace_period: 1s # 等待 1 秒
stop_grace_period: 1m30s # 等待 1 分 30 秒
```

默认的等待时间是 10 秒。

## stop_signal

设置停止容器的替代信号。默认情况下使用 SIGTERM 。

以下示例，使用 SIGUSR1 替代信号 SIGTERM 来停止容器。

```yaml
stop_signal: SIGUSR1
```

## sysctls

设置容器中的内核参数，可以使用数组或字典格式。

```yaml
sysctls:
  net.core.somaxconn: 1024
  net.ipv4.tcp_syncookies: 0
sysctls:
  - net.core.somaxconn=1024
  - net.ipv4.tcp_syncookies=0
```

## tmpfs

在容器内安装一个临时文件系统。可以是单个值或列表的多个值。

```yaml
tmpfs: /run
tmpfs:
  - /run
  - /tmp
```

## ulimits

覆盖容器默认的 ulimit。

```yaml
ulimits:
  nproc: 65535
  nofile:
    soft: 20000
    hard: 40000
```

## volumes

将主机的数据卷或着文件挂载到容器里。

```yaml
version: "3.7"
services:
  db:
    image: postgres:latest
    volumes:
    - "/localhost/postgres.sock:/var/run/postgres/postgres.sock"
    - "/localhost/data:/var/lib/postgresql/data"
```

