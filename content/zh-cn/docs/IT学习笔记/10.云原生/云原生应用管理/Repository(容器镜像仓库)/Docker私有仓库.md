---
title: Docker私有仓库
---

# 私有仓库

有时候使用 Docker Hub 这样的公共仓库可能不方便，用户可以创建一个本地仓库供私人使用。

本节介绍如何使用本地仓库。

docker-registry 是官方提供的工具，可以用于构建私有的镜像仓库。本文内容基于 docker-registry v2.x 版本。

# 安装运行 docker-registry

## 容器运行

你可以通过获取官方 registry 镜像来运行。

$ docker run -d -p 5000:5000 --restart=always --name registry registry

这将使用官方的 registry 镜像来启动私有仓库。默认情况下，仓库会被创建在容器的 /var/lib/registry 目录下。你可以通过 -v 参数来将镜像文件存放在本地的指定路径。例如下面的例子将上传的镜像放到本地的 /opt/data/registry 目录。

$ docker run -d \ -p 5000:5000 \ -v /opt/data/registry:/var/lib/registry \ registry

# 在私有仓库上传、搜索、下载镜像

创建好私有仓库之后，就可以使用 docker tag 来标记一个镜像，然后推送它到仓库。例如私有仓库地址为 127.0.0.1:5000。

先在本机查看已有的镜像。

$ docker image lsREPOSITORY TAG IMAGE ID CREATED VIRTUAL SIZEubuntu latest ba5877dc9bec 6 weeks ago 192.7 MB

使用 docker tag 将 ubuntu:latest 这个镜像标记为 127.0.0.1:5000/ubuntu:latest。

格式为 docker tag IMAGE\[:TAG] \[REGISTRY_HOST\[:REGISTRY_PORT]/]REPOSITORY\[:TAG]。

$ docker tag ubuntu:latest 127.0.0.1:5000/ubuntu:latest$ docker image lsREPOSITORY TAG IMAGE ID CREATED VIRTUAL SIZEubuntu latest ba5877dc9bec 6 weeks ago 192.7 MB127.0.0.1:5000/ubuntu:latest latest ba5877dc9bec 6 weeks ago 192.7 MB

使用 docker push 上传标记的镜像。

$ docker push 127.0.0.1:5000/ubuntu:latestThe push refers to repository \[127.0.0.1:5000/ubuntu]373a30c24545: Pusheda9148f5200b0: Pushedcdd3de0940ab: Pushedfc56279bbb33: Pushedb38367233d37: Pushed2aebd096e0e2: Pushedlatest: digest: sha256:fe4277621f10b5026266932ddf760f5a756d2facd505a94d2da12f4f52f71f5a size: 1568

用 curl 查看仓库中的镜像。

$ curl 127.0.0.1:5000/v2/\_catalog{"repositories":\["ubuntu"]}

这里可以看到 {"repositories":\["ubuntu"]}，表明镜像已经被成功上传了。

先删除已有镜像，再尝试从私有仓库中下载这个镜像。

$ docker image rm 127.0.0.1:5000/ubuntu:latest$ docker pull 127.0.0.1:5000/ubuntu:latestPulling repository 127.0.0.1:5000/ubuntu:latestba5877dc9bec: Download complete511136ea3c5a: Download complete9bad880da3d2: Download complete25f11f5fb0cb: Download completeebc34468f71d: Download complete2318d26665ef: Download complete$ docker image lsREPOSITORY TAG IMAGE ID CREATED VIRTUAL SIZE127.0.0.1:5000/ubuntu:latest latest ba5877dc9bec 6 weeks ago 192.7 MB

# 注意事项

如果你不想使用 127.0.0.1:5000 作为仓库地址，比如想让本网段的其他主机也能把镜像推送到私有仓库。你就得把例如 192.168.199.100:5000 这样的内网地址作为私有仓库地址，这时你会发现无法成功推送镜像。

这是因为 Docker 默认不允许非 HTTPS 方式推送镜像。我们可以通过 Docker 的配置选项来取消这个限制，或者查看下一节配置能够通过 HTTPS 访问的私有仓库。

Ubuntu 14.04, Debian 7 Wheezy

对于使用 upstart 的系统而言，编辑 /etc/default/docker 文件，在其中的 DOCKER_OPTS 中增加如下内容：

DOCKER_OPTS="--registry-mirror=https://registry.docker-cn.com --insecure-registries=192.168.199.100:5000"

重新启动服务。

$ sudo service docker restart

Ubuntu 16.04+, Debian 8+, centos 7

对于使用 systemd 的系统，请在 /etc/docker/daemon.json 中写入如下内容（如果文件不存在请新建该文件）

{ "registry-mirror": \[ "https://registry.docker-cn.com" ], "insecure-registries": \[ "192.168.199.100:5000" ]}

注意：该文件必须符合 json 规范，否则 Docker 将不能启动。

其他

对于 Docker for Windows 、 Docker for Mac 在设置中编辑 daemon.json 增加和上边一样的字符串即可。

Harbor

从 github 的 releases 下载安装包，然后解压，执行其中 install.sh 脚本
