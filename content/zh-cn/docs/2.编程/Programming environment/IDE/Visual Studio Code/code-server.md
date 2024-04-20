---
title: code-server
---

# 概述

> 参考：
>
> - [GitHub 项目，coder/code-server](https://github.com/coder/code-server)
> - [官网](https://coder.com/)

code-server 可以让 VS Code 运行在任何机器上，并在浏览器中访问。

code-server 是一个免费的基于浏览器的 IDE，而 Coder 是该企业提供的收费版本。

# 部署

```shell
mkdir -p /opt/code-server/user-data-dir

nerdctl run --name code-server --network host -d \
-v "/opt/code-server/user-data-dir:/root" \
-u "$(id -u):$(id -g)" \
-e "DOCKER_USER=$USER" \
codercom/code-server:latest
```

## 个性部署

```shell
nerdctl run --name code-server --network host -d \
-v "/opt/code-server/user-data-dir:/root" \
-v "/usr/local/go:/usr/local/go" \
-v "/root/go:/root/go" \
-v "/root/projects:/root/projects" \
-u "$(id -u):$(id -g)" \
-e "DOCKER_USER=$USER" \
-e "GOPROXY=https://goproxy.cn,https://goproxy.io,https://mirrors.aliyun.com/goproxy/,direct" \
codercom/code-server:latest \
--bind-addr 0.0.0.0:58080
```

# code-server 关联文件

**~/** # code-server 运行时生成的持久化数据都在当前用户的家目录下。刚部署完只有如下几个文件

```shell
.
├── .config
│   └── code-server
│       └── config.yaml
└── .local
    └── share
        └── code-server
            ├── coder-logs
            │   ├── code-server-stderr.log
            │   └── code-server-stdout.log
            └── heartbeat
```

**~/.config/code-server/config.yaml** # 登录密码、监听端口 等基本信息
**~/.vscode-remote/data/** #
**~/.local/share/code-server/** # 默认的用户数据(即.持久化数据)路径。可以通过 `--user-data-dir` 命令行参数指定。

> 默认的目录中有一些文件是不属于用户数据。比如 coder-logs、hearteat 等。

- **./coder.json** #
- .**/coder-logs/** #
- .**/extensions/** # 默认扩展保存路径。可以通过 `--extensions-dir` 命令行参数指定。
- **./logs/** #
- .**/machineid** #
- .**/User/** #
  - **./settings.json** # 默认配置
