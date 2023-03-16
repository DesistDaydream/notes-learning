---
title: Server 与 Runner 部署
---

# GitHub Server

官方文档：<https://docs.drone.io/server/provider/github/>

1. 在 GitHub 上创建 OAuth Application，通过此连接创建。也可以点击 头像—Settings—Developer settings—OAuth Apps 来进入创建页面

   1. 在创建页面需要指明该 OAuth Application 的 name、URL、callback URL 这三个信息。下面是三个信息的例子：

      1. Application name：Grone

      2. Homepage URL：<http://10.10.100.150> # 这个 URL 指的是将要运行 Drone 的 GitHub Server 端设备 IP。也可以使用域名来表示。

      3. Authorization callback URL：<http://10.10.100.150/login>

   2. 创建完成后，会生成两个信息 Client ID 与 Client Secret。下面是这两个信息的例子

      1. Client ID：df6dd0a1d49a4a26a151

      2. Client Secret：26507ea48eab7612f0bd8865f9d6091608baa7f3

   3. OAuth Application 的相关信息在安装 Drone 的 GitHub Server 端时会用到

2. 创建一个共享密钥，该密钥用来在 Server 与 Runner 两端进行通信时进行验证。使用如下命令即可

   1. openssl rand -hex 16 #

   2. 生成字符串为：3de4792f96c2aa9e483b36a54b23d70a

3. 使用 docker 镜像运行 Drone 的 GitHub Server 端

    docker run \
      --volume=/var/lib/drone:/data \
      --env=DRONE_GITHUB_CLIENT_ID=df6dd0a1d49a4a26a151 \
      --env=DRONE_GITHUB_CLIENT_SECRET=26507ea48eab7612f0bd8865f9d6091608baa7f3 \
      --env=DRONE_RPC_SECRET=3de4792f96c2aa9e483b36a54b23d70a \
      --env=DRONE_SERVER_HOST=10.10.100.150 \
      --env=DRONE_SERVER_PROTO=http \
      --publish=80:80 \
      --publish=443:443 \
      --restart=always \
      --detach=true \
      --name=drone \
      drone/drone:1

1. Server 端启动后， 通过 10.10.100.150 来访问，会弹出让 GitHub 授权的页面

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqgger/1616077734543-9c1fcf23-8238-46a8-b2b6-75e5be31e563.jpeg)

1. 授权完成，即可登录到 Drone Server 端的 web 界面。可以看到，已经同步了我 GitHub 上的所有仓库，我可以使用这些仓库的代码来执行流水线操作。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqgger/1616077734552-f99408e2-e44f-439e-bc0a-7286a3c644fd.jpeg)

1. 选择一个仓库，进入，点击 Activate Repository ，这样这个仓库就会激活，并在 GitHub 上该仓库的设置中，添加 webhooks 信息，当 GitHub 收到代码提交时，就会自动向 Drone Server 发送信息，执行自动化 pipeline。

2. 但是！值得注意的是，由于 Drone Server 无法手动触发 pipeline ，所以 Server 会等待 GitHub 上代码变更之后的 WebHook 触发。但是查看下图可以看到，由于 Server 是安装在内网的，Github 服务器是无法访问的，所以就算部署完 Server ，也无法通过 GitHub 代码提交来触发 pipeline。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/qqgger/1616077734546-284ab667-96a1-4dcf-a36d-73f507314160.jpeg)

# GitLab Server

GitLab Server 的部署方式与 GitHub 基本一致

官方文档：<https://docs.drone.io/server/provider/gitlab/>

1. 创建 OAuth Application，通过此连接创建。也可以点击 头像—Settings(设置)—Applications(应用) 来进入创建页面

   1. 在创建页面需要指明 Name(名称) 与 Redirect URL

      1. Name：Grone

      2. Redirect URL：<http://10.10.100.150/login> # 这个 URL 指的是将要运行 Drone 的 GitLab Server 端设备 IP。也可以使用域名来表示。

   2. 创建完成后，会生成两个信息 Application ID 与 Secret。下面是这两个信息的例子

      1. Application ID(应用程序 ID)：cc1443ac09e700fb93676bd22b42e4e79f47db9db7b0c48938e4358fcbb01031

      2. Secret(密码)：edaebba6a1e0a5ef80b5f89b2bfd3857fe9c088e14299a27cb16147b71098a2b

   3. Application 的相关信息在安装 Drone 的 GitLab Server 端时会用到

2. 创建一个共享密钥，该密钥用来在 Server 与 Runner 两端进行通信时进行验证。使用如下命令即可

   1. openssl rand -hex 16 #

   2. 生成字符串为：3de4792f96c2aa9e483b36a54b23d70a

3. 使用 docker 镜像运行 Drone 的 GitLab Server 端

    docker run \
      --volume=/var/lib/drone:/data \
      --env=DRONE_GITLAB_SERVER=http://10.10.100.151/ \
      --env=DRONE_GITLAB_CLIENT_ID=cc1443ac09e700fb93676bd22b42e4e79f47db9db7b0c48938e4358fcbb01031 \
      --env=DRONE_GITLAB_CLIENT_SECRET=edaebba6a1e0a5ef80b5f89b2bfd3857fe9c088e14299a27cb16147b71098a2b \
      --env=DRONE_RPC_SECRET=3de4792f96c2aa9e483b36a54b23d70a \
      --env=DRONE_SERVER_HOST=10.10.100.150 \
      --env=DRONE_SERVER_PROTO=http \
      --publish=80:80 \
      --publish=443:443 \
      --restart=always \
      --detach=true \
      --name=drone \
      drone/drone:1

1. Server 端启动后， 通过 10.10.100.150 来访问，会弹出让 GitLab 授权的页面

2. 授权完成即可开始使用 Drone。

Note：

1. 如果 GitLab 在配置 webhook 时报错 url is blocked Requests to the local network are not allowed，那么修改 GitLab 的设置即可解决。修改路径如下：

2. 登录管理员账号，在 Admin area 中(页面左上角的)，左侧 Settings -> Network -> Outbound requests，勾选 Allow requests to the local network from hooks and services

# Docker Runner

官方文档：<https://docs.drone.io/runner/docker/overview/>

DRONE_RPC_HOST # 指定 Drone Server 的地址和端口

DRONE_RPC_SECRET # 指定 Runner 与 Server 交互的认证信息。需要与 Server 端的该参数的值保持一致。

    docker run -d \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -e DRONE_RPC_PROTO=http \
      -e DRONE_RPC_HOST=10.10.100.150 \
      -e DRONE_RPC_SECRET=3de4792f96c2aa9e483b36a54b23d70a \
      -e DRONE_RUNNER_CAPACITY=2 \
      -e DRONE_RUNNER_NAME=${HOSTNAME} \
      -p 3000:3000 \
      --restart always \
      --name runner \
      drone/drone-runner-docker:1

# Kuberntes Runner

官方文档：
