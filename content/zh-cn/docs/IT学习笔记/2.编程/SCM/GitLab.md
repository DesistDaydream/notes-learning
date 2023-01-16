---
title: GitLab
---

# GitLab 介绍

<https://www.qikqiak.com/post/gitlab-install-on-k8s/>

官方文档：<https://docs.gitlab.com/>

# GitLab 部署

官方文档：<https://docs.gitlab.com/ee/install/README.html>

## 通过官方的 linux 软件包安装

<https://about.gitlab.com/install/> 根据该页面选择想要运行 GitLab 的 Linux 发行版，可以通过 yum、apt 等方式直接安装 GitLab 及其所需的依赖。

## 使用 Docker 运行 GitLab

官方文档：<https://docs.gitlab.com/omnibus/docker/>

    docker run --detach \
       --hostname 10.10.100.151 \
       --publish 443:443 --publish 80:80 --publish 9022:22 \
       --name gitlab \
       --restart always \
       --volume /root/gitlab/config:/etc/gitlab \
       --volume /root/gitlab/logs:/var/log/gitlab \
       --volume /root/gitlab/data:/var/opt/gitlab \
       gitlab/gitlab-ce:latest

部署完成后，使用 root 和 第一次打开 web 页面时设置的密码，即可登录管理员账户。
