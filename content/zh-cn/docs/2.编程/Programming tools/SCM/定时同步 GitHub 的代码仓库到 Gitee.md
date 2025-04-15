---
title: 定时同步 GitHub 的代码仓库到 Gitee
linkTitle: 定时同步 GitHub 的代码仓库到 Gitee
weight: 20
---


# 概述

> 参考：
>
> - 

# 利用 GitHub Action 同步

该功能已经有很多实现了，这篇文章以 <https://github.com/Yikun/hub-mirror-action> 项目为例。这个项目的基本逻辑是这样的：

- 通过 GitHub Actions 构建一个 Docker 容器，在 Docker 容器中，引入 Gitee 的私钥，这样可以在容器中使用 git 命令向 Gitee push 代码而不用输入密码了
- 容器启动后，在容器内 pull github 上的代码，并 push 到 gitee 上。

首先先来一个最基本的 Action 的 workflow 文件示例

```yaml
name: Gitee repos mirror periodic job
on:
  # 取消 push 的注释后，向本仓库推送代码即可开始 Gitee 同步
  # push:
  schedule:
    # 每天北京时间9点跑
    - cron: "0 1 * * *"
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Cache phpdragon src repos
        # 使用 github 官方提供的 action 来引用发行版的主要版本
        uses: actions/cache@v1
        with:
          path: /home/runner/work/phpdragon/phpdragon-cache
          key: ${{ runner.os }}-phpdragon-repos-cache
      - name: Mirror the Github organization repos to Gitee.
        # 这里我将对方项目原封不动 copy 到自己的 github 上了，所以这个步骤就直接使用自己的 action
        uses: DesistDaydream/hub-mirror-action@main
        with:
          # 必选，需要同步的Github用户（源）
          src: github/DesistDaydream
          # 必选，需要同步到的Gitee的用户（目的）
          dst: gitee/DesistDaydream
          # 必选，Gitee公钥对应的私钥
          dst_key: ${{ secrets.GITEE_PRIVATE_KEY }}
          # 必选，Gitee对应的用于创建仓库的token
          dst_token: ${{ secrets.GITEE_TOKEN }}
          # 黑、白名单，静态名单机制，可以用于更新某些指定库
          # static_list: repo_name
          black_list: "eHualu,kubernetesAPI,v2ray"
          # white_list: 'repo_name,repo_name2'
          # force_update: true
```

必选参数

- `src` 需要被同步的源端账户名，如 github/DesistDaydream，表示 Github 的 DesistDaydream 账户。
- `dst` 需要同步到的目的端账户名，如 gitee/DesistDaydream，表示 Gitee 的 DesistDaydream 账户。
- `dst_key` 与 Gitee 公钥对应的私钥，给 GitHub Actions 激活后创建的容器中 git 命令认证所用。
- `dst_token` 用于自动创建不存在的仓库。
- 注意： dst_key 与 dst_token 的值是通过 GitHub 的 Secrets 功能 引用的，类似于 k8s 的 secret 功能。

可选参数

- `account_type` 默认为 user，源和目的的账户类型，可以设置为 org（组织）或者 user（用户），目前仅支持**同类型账户**（即组织到组织，或用户到用户）的同步。
- `clone_style` 默认为 https，可以设置为 ssh 或者 https。
- `cache_path` 默认为'', 将代码缓存在指定目录，用于与 actions/cache 配合以加速镜像过程。
- `black_list` 默认为'', 配置后，黑名单中的 repos 将不会被同步，如“repo1,repo2,repo3”。
- `white_list` 默认为'', 配置后，仅同步白名单中的 repos，如“repo1,repo2,repo3”。
- `static_list` 默认为'', 配置后，仅同步静态列表，不会再动态获取需同步列表（黑白名单机制依旧生效），如“repo1,repo2,repo3”。
- `force_update` 默认为 false, 配置后，启用 git push -f 强制同步，**注意：开启后，会强制覆盖目的端仓库**。
- `debug` 默认为 false, 配置后，启用 debug 开关，会显示所有执行命令。

# 获取并配置 GitHub 连接 Gitee 所需的认证信息

认证信息比较敏感，详见 TOKEN 文章

## 获取 Gitee TOKEN

使用该连接：<https://gitee.com/profile/personal_access_tokens/new>，添加`令牌描述`后，点击`提交`以生成 TOKEN

## 获取密钥对

随便找一个有 ssh-keygen 命令的主机，用于生成一对密钥。使用 ssh-keygen 命令生成密钥对，ssh-keygen 命令用法详见此处

```
ssh-keygen -t rsa -C 我的邮箱
```

## 配置 Gitee 公钥

在 [Gitee 的配置页面中](https://gitee.com/profile/sshkeys)添加公钥信息。以便 GitHub 使用 私钥连接时，可以通过认证。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wn0scx/1640568955462-d3dbe873-2a73-4539-a477-0cfa71fb8a43.png)

## 配置 私钥 和 TOKEN

在 GitHub 以加密的方式传入到容器中。如果不加密，直接写到代码仓库中，那其他人就都看到了。。。该操作主要是针对 代码仓库而言的，因为 私钥和 TOKEN 的信息，是需要在 Action 中引用的，而 Action 本身就是一摞代码~

在[项目仓库的 Setting 中的 Secrets](https://github.com/DesistDaydream/hub-mirror-action/settings/secrets) 中[添加](https://github.com/DesistDaydream/hub-mirror-action/settings/secrets/new) 私钥 与 TOKEN 的变量。

### 添加 TOKEN

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/wn0scx/1640569010998-1d5f41bd-359d-4b4c-ae4d-d4352ba41ab7.png)

### 添加私钥

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wn0scx/1616903594321-e357ab96-5486-42f9-ba85-9fdf869e9fbb.png)

# 使用 Gitee 的 仓库镜像管理

> 参考：
>
> - [公众号，更优雅的GitHub/Gitee仓库镜像同步：一次提交，同时更新两个平台](https://mp.weixin.qq.com/s/BkK979TRAURg8MCxXwx6Ow)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/scm/20250415092345553.png)

需要如下 GitHub 的权限

- **Repository permissions**
    - **Webhooks** # 读/写。用于 Gitee 自动为 GitHub 仓库创建 [Webhook](https://github.com/DesistDaydream/notes-learning/settings/hooks)

