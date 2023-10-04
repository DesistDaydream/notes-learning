---
title: Rancher 配置
linkTitle: Rancher 配置
date: 2023-11-23T11:49
weight: 20
---

# 概述

> 参考：
> 
> -



Rancher Server URL 的修改

当 Rancher Server URL 变更后(比如从 40443 变到 60443)，则还需要连带修改以下几部分

- Rancher Web 页面最上面的标签，进入`系统设置`，修改`server-url`。

- k8s 集群中，修改 cattle-system 名称空间中，名为`cattle-credentials-XXX`的 secret 资源中的 .data.url 字段的值，这个值是用 base64 编码的。

- echo -n "https://X.X.X.X:60443" | bas64 ，通过该命令获取编码后的 url，然后填入 .data.url 字段中

- k8s 集群中，修改 cattle-cluster-agent-XX 和 cattle-node-agent-XX 这些 pod 的 env 参数，将其中的 CATTLE_SERVER 的值改为想要的 URL。

- cattle-node-agent 在 2.5.0 版本之后没有了，就不用改了。

导入集群的 yaml 文件位置

打开 `https://RancherIP/v3/cluster/集群ID/clusterregistrationtokens` 页面

在 data 字段下，可以看到获取 yaml 文件的 URL，可能会有多组，一般以时间最新的那组为准。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ggn0dn/1616114779749-bd6fd7cc-32cb-41b8-9122-2047f125c4a7.png)
