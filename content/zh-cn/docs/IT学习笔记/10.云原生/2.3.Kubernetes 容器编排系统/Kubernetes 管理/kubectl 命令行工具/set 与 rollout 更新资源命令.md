---
title: set 与 rollout 更新资源命令
---

# set # 在对象上设定特定的特性

**kubectl set COMMAND \[OPTIONS]**
COMMAND

- **env** # Update environment variables on a pod template
- **image** # 更新一个 Pod 模板的镜像
- **resources** # Update resource requests/limits on objects with pod templates
- **selector** # Set the selector on a resource
- **serviceaccount** # Update ServiceAccount of a resource
- **subject** # Update User, Group or ServiceAccount in a RoleBinding/ClusterRoleBinding

## kubectl set image # 更新资源

更新资源的现有容器映像。可能的资源包括(不区分大小写)pod (po), replicationcontroller (rc), deployment (deploy), daemonset (ds), replicaset (rs)

### Syntax(语法)

**kubectl set image (-f FILENAME | TYPE NAME) CONTAINER_NAME_1=CONTAINER_IMAGE_1 ... CONTAINER_NAME_N=CONTAINER_IMAGE_N**

OPTIONS

EXAMPLE

- 更新(myapp-deploy 这个控制器下的 pod 模板中的名为 myapp 这个容器中的镜像)到 ikubernetes/myapp:v3 这个版本,并且暂停 myapp-deploy 这个 deployment 的滚动更新，该命令的作用是只更新一个 pod，暂停后续更新，以便查看新版本在生产环境中的运行情况，正常之后，再开始暂停的更新即可
  - kubectl set image deployment myapp-deploy myapp=ikubernetes/myapp:v3 && kubectl rollout pause deployment myapp-deploy

# rollout # 管理资源的滚动更新

**kubectl rollout COMMAND \[OPTIONS]**
COMMAND

- **history** #滚动更新视图
- **pause** #标记提供的资源以暂停这个资源
- **restart**# 重启一个资源
- **resume** #重新开始被暂停的资源
- **status**#显示滚动更新的状态
- **undo** #撤销以前的滚动更新

## history # 查看滚动更新的历史情况

kubectl rollout history (TYPE NAME | TYPE/NAME) \[flags] \[options]

REVISION # 指明更新的版本序号

CHANGE-CAUSE # 指明该次更新执行的具体命令，只有在 apply 的时候使用--record 参数，该项才有内容

EXAMPLE

- **kubectl rollout history deployment myapp-deploy**

## pause # 暂停滚动更新的资源

kubectl rollout pause RESOURCE \[options]

EXAMPLE

- **kubectl rollout pause deployment myapp-deploy**

## restart # 重启一个资源

滚动重启指定的资源

### Syntax(语法)

**kubectl rollout restart RESOURCE \[OPTIONS]**

**EXAMPLE**

- 滚动重启 monitoring 名称空间下，名为 node-exporter 的 daemonset 类型资源
  - **kubectl rollout restart -n monitoring daemonset node-exporter**

## resume # 重新开始被暂停的资源

kubectl rollout resume RESOURCE \[options]

EXAMPLE

- kubectl rollout resume deployment myapp-deploy

## status # 查看更新状态

语法结构
kubectl rollout status (TYPE NAME | TYPE/NAME) \[flags] \[options]

EXAMPLE

- kubectl rollout status deployment myapp # 查看 myapp 这个 deployment 的更新状态

## undo # 回滚

语法结构
kubectl rollout undo (TYPE NAME | TYPE/NAME) \[flags] \[OPTIONS]

OPTIONS

- --to-revision=NUM # 指定要回滚到哪个修订版，默认为 0，最后的修正版。查看修订版的 NUM 可以使用 kubectl rollout history 命令，该命令显示出的 REVISION 下面的数字就是 NUM，注意：如果我从 V1 更新到 V2 再更新到 V3，然后回滚到 V1，那么此时我再回滚到上一版的时候，指的是回滚到 V3 版

EXAMPLE

- kubectl rollout undo deployment myapp # 回滚 myapp 这个 deployment 的更新状态
