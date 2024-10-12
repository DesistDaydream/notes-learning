---
title: Jenkins 与 skaffold 示例
---

# 概述

准备操作

配置 webhook

配置 Jenkins Pipeline 自动 clone 代码并获取 Jenkinsfile 文件

代码示例

Jenkinsfile

skaffold.yaml

main.go

Dockerfile

k8s-pod.yaml

运行结果：

# 准备操作

# 配置 webhook

从 jenkins 项目中，获取 webhook 的 URL 和 TOKEN，填到 gitlab 指定项目的 webhooks 中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cyuibc/1616077495082-a11319fe-589d-4c0f-b45d-bb0b3bc52777.jpeg)

Jenkins 项目的 TOKEN 需要在构建触发器栏目中，点击 高级，然后点击 Generate 来生成。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cyuibc/1616077495099-5da918fd-e308-4ccd-9873-19a87c8ae444.jpeg)

在 GitLab 中填写 URL 和 TOKEN ，并点击 Add Webhook 即可，添加完成后，可以在最下方点击 Test 来测试连通性。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cyuibc/1616077495099-b091204a-a450-4b67-b215-cf3e5ec8a350.jpeg)

配置 Jenkins Pipeline 自动 clone 代码并获取 Jenkinsfile 文件

使用如下配置，开始构建前，让 Jenkins 自动获取代码仓库中的 Jenkinsfile 文件，并根据该文件执行 pipeline

在流水线类型的任务中，进行如下配置

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cyuibc/1616077495088-a4c61fb6-b61f-4697-a5ca-19877f84dad8.jpeg)

# 代码示例

## Jenkinsfile

    pipeline {
      agent any
      stages {
        stage('build') {
          steps {
            sh 'export TAG=1.0 && skaffold run'
          }
        }
      }
    }

skaffold.yaml

    apiVersion: skaffold/v2alpha1
    kind: Config
    build:
      artifacts:
      - image: 172.38.40.180/test/pipeline-test
      tagPolicy:
        envTemplate:
          template: "{{.TAG}}"
      local:
        push: true
    deploy:
      kubectl:
        manifests:
          - k8s-*

main.go

```go
package main
import (
  "fmt"
  "time"
)
func main() {
  for {
    fmt.Println("Hello Jenkins,Skafflod!")
    time.Sleep(time.Second * 5)
  }
}
```

Dockerfile

```Dockerfile
FROM golang:1.14.7
WORKDIR /src
COPY main.go .
RUN go build main.go
FROM ubuntu
COPY --from=0 /src/main .
CMD ["./main"]
```

k8s-pod.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: getting-started
spec:
  containers:
  - name: getting-started
    image: 172.38.40.180/test/pipeline-test:1.0
```

运行结果：

```bash
Started by GitLab push by DesistDaydream
# 从配置中指定的 URL 中获取 Jenkinsfile
Obtained Jenkinsfile from git http://10.20.5.5:10080//lich_wb/pipeline-skaffold-test.git
Running in Durability level: MAX_SURVIVABILITY
[Pipeline] Start of Pipeline
[Pipeline] node
Running on Jenkins in /opt/jenkins/jenkins-data/workspace/oc_dev/pipeline-skaffold-test
[Pipeline] {
[Pipeline] stage
# 检查 SCM，并 clone 代码到本地。
[Pipeline] { (Declarative: Checkout SCM)
[Pipeline] checkout
using credential git_global
  > /usr/bin/git rev-parse --is-inside-work-tree # timeout=10
Fetching changes from the remote Git repository
  > /usr/bin/git config remote.origin.url http://10.20.5.5:10080//lich_wb/pipeline-skaffold-test.git # timeout=10
Fetching upstream changes from http://10.20.5.5:10080//lich_wb/pipeline-skaffold-test.git
  > /usr/bin/git --version # timeout=10
using GIT_ASKPASS to set credentials
  > /usr/bin/git fetch --tags --progress http://10.20.5.5:10080//lich_wb/pipeline-skaffold-test.git +refs/heads/*:refs/remotes/origin/* # timeout=10
skipping resolution of commit remotes/origin/master, since it originates from another repository
  > /usr/bin/git rev-parse refs/remotes/origin/master^{commit} # timeout=10
  > /usr/bin/git rev-parse refs/remotes/origin/origin/master^{commit} # timeout=10
Checking out Revision 4309c942163c122bfc706076a51e5de8118d1e37 (refs/remotes/origin/master)
  > /usr/bin/git config core.sparsecheckout # timeout=10
  > /usr/bin/git checkout -f 4309c942163c122bfc706076a51e5de8118d1e37 # timeout=10
Commit message: "change"
  > /usr/bin/git rev-list --no-walk 504d09cea618173edcc56ee4e9ea0edb573b0ccf # timeout=10
[Pipeline] }
[Pipeline] // stage
[Pipeline] withEnv
[Pipeline] {
[Pipeline] stage
[Pipeline] { (build)
# 开始执行 skaffold，构建、推送、并部署到 kubernetes 集群
[Pipeline] sh
+ export TAG=1.0
+ TAG=1.0
+ skaffold run
Generating tags...
  - 172.38.40.180/test/pipeline-test -> 172.38.40.180/test/pipeline-test:1.0
Checking cache...
  - 172.38.40.180/test/pipeline-test: Found. Pushing
The push refers to repository [172.38.40.180/test/pipeline-test]
c22c4d1a8b07: Preparing
095624243293: Preparing
a37e74863e72: Preparing
8eeb4a14bcb4: Preparing
ce3011290956: Preparing
a37e74863e72: Layer already exists
8eeb4a14bcb4: Layer already exists
c22c4d1a8b07: Layer already exists
095624243293: Layer already exists
ce3011290956: Layer already exists
1.0: digest: sha256:7f7a76e3a0a0fafb43283a9269e6c370a883202d3a37640f43730487e9ad472b size: 1363
Tags used in deployment:
  - 172.38.40.180/test/pipeline-test -> 172.38.40.180/test/pipeline-test:1.0@sha256:7f7a76e3a0a0fafb43283a9269e6c370a883202d3a37640f43730487e9ad472b
Starting deploy...
  - pod/getting-started configured
Waiting for deployments to stabilize...
Deployments stabilized in 31.474009ms
You can also run [skaffold run --tail] to get the logs
[Pipeline] }
[Pipeline] // stage
[Pipeline] }
[Pipeline] // withEnv
[Pipeline] }
[Pipeline] // node
[Pipeline] End of Pipeline
Finished: SUCCESS
```
