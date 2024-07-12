---
title: GitLab
linkTitle: GitLab
date: 2024-05-15T20:03
weight: 1
---


# 概述

> 参考：
>
> - [官方文档](https://docs.gitlab.com/)
> - https://www.qikqiak.com/post/gitlab-install-on-k8s/


# GitLab 部署

官方文档：<https://docs.gitlab.com/ee/install/README.html>

## 通过官方的 linux 软件包安装

<https://about.gitlab.com/install/> 根据该页面选择想要运行 GitLab 的 Linux 发行版，可以通过 yum、apt 等方式直接安装 GitLab 及其所需的依赖。

## 使用 Docker 运行 GitLab

官方文档：<https://docs.gitlab.com/omnibus/docker/>

```bash
docker run --detach \
   --hostname 10.10.100.151 \
   --publish 443:443 --publish 80:80 --publish 9022:22 \
   --name gitlab \
   --restart always \
   --volume /root/gitlab/config:/etc/gitlab \
   --volume /root/gitlab/logs:/var/log/gitlab \
   --volume /root/gitlab/data:/var/opt/gitlab \
   gitlab/gitlab-ce:latest
```

部署完成后，使用 root 和 第一次打开 web 页面时设置的密码，即可登录管理员账户。

# 最佳实践

## 从 GitHub 导入仓库到 GitLab 并定时同步

https://www.jianshu.com/p/0959d021c281

> Notes: GitLab 中有**镜像仓库**的能力，在 `项目设置 - 仓库 - 镜像仓库` 中。但是<font color="#ff0000">只有商业版才能使用</font>从别的 git 仓库同步到 GitLab 的功能。

这里以从 GitHub 导入到 GitLab 为例

一、导入仓库

GitLab 选择 `新建项目/仓库 - 导入项目` 选择 GitHub 项目（需要 [GitHub 的 Personal access tokens](https://github.com/settings/tokens) 以通过 GitHub 认证）

二、关联仓库

假如已经成功导入了 DesistDaydream/net_tool 项目。需要通过 [git CLI](/docs/2.编程/Programming%20tools/SCM/Git/git%20CLI.md) 让两个仓库互相管理起来

> 在设备上使用 [OpenSSH Utilities](/docs/4.数据通信/Utility/OpenSSH/OpenSSH%20Utilities.md) 创建 id_rsa.pub，并拷贝到 [GitLab](https://gitlab.com/-/profile/keys) 和 [GitHub](https://github.com/settings/keys) 的 SSH 密钥中

先 clone 刚才已导入的仓库，然后查看一下当前的 remote 信息

```bash
~]# git clone git@gitlab.com:DesistDaydream/net_tool.git
~]# git remote -v
origin  git@gitlab.com:DesistDaydream/net_tool.git (fetch)
origin  git@gitlab.com:DesistDaydream/net_tool.git (push)
```

添加 GitHub 远程仓库，此时本地项目中，包含了两个 remote 信息

```bash
~]# git remote add github git@github.com:DesistDaydream/net_tool.git
~]# git remote -v
github  git@github.com:DesistDaydream/net_tool.git (fetch)
github  git@github.com:DesistDaydream/net_tool.git (push)
origin  git@gitlab.com:DesistDaydream/net_tool.git (fetch)
origin  git@gitlab.com:DesistDaydream/net_tool.git (push)
```

三、同步代码

从 GitHub 仓库中同步代码到本地，然后再从本地同步代码到 GitLab

```bash
~]# git pull github main
~]# git push origin main  --force
```

# GitLab CI

> 参考:
>
> - [官方文档，主题 - 使用 CI/CD 构建你的应用](https://docs.gitlab.com/ee/topics/build_your_application.html)

## .gitlab-ci.yml

> 参考:
>
> - [官方文档，CI/CD YAML 语法参考](https://docs.gitlab.com/ee/ci/yaml/)

顶层字段

- **stages**(\[]STRING)
- **${JOB_NAME}**(OBJECT) # 

```yaml
stages:
  - build
  - release

build-job:
  tags:
    - instance-runner-tc
  stage: build
  script:
    - bash observability/monitoring/packaging_server.sh
  artifacts:
    paths:
      - server.tar.gz

release-job:
  tags:
    - instance-runner-tc
  stage: release
  image: registry.gitlab.com/gitlab-org/release-cli:latest
  # 解决 release-cli 无法信任自建 CA 问题。https://gitlab.com/gitlab-org/release-cli/-/issues/47
  before_script:
    - sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
    - apk --no-cache add openssl ca-certificates
    - mkdir -p /usr/local/share/ca-certificates/extra
    - openssl s_client -connect ${CI_SERVER_HOST}:${CI_SERVER_PORT} -servername ${CI_SERVER_HOST} -showcerts </dev/null 2>/dev/null | sed -e '/-----BEGIN/,/-----END/!d' | tee "/usr/local/share/ca-certificates/${CI_SERVER_HOST}.crt" >/dev/null
    - update-ca-certificates
  script:
    - echo "检查 CI_COMMIT_TAG ${CI_COMMIT_TAG}"
    - ls -l
    - test -f server.tar.gz && echo "server.tar.gz exists" || echo "server.tar.gz does not exist"
  release:
    name: "Release $CI_COMMIT_TAG"
    # tag_name: $CI_COMMIT_TAG
    tag_name: "v1.0.0"
    description: "Release created by GitLab CI"
    assets:
      links:
        - name: "server.tar.gz"
          url: "${CI_PROJECT_URL}/-/jobs/artifacts/${CI_COMMIT_REF_NAME}/raw/server.tar.gz?job=build-job"
          link_type: "package"
  needs:
    - job: build-job
      artifacts: true

```