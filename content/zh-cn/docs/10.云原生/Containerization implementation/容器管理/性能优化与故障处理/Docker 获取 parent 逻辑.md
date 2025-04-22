---
title: Docker 获取 parent 逻辑
---

每次 kubelet 获取镜像列表时，docker 都会获取一遍镜像的 parent

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fnmwl8/1648798871972-b0aab87e-1e9d-47c9-8053-1e976e8a8f70.png)

具体逻辑在这里 [image/store.go](https://github.com/moby/moby/blob/20.10/image/store.go#L202)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fnmwl8/1648798901545-0898c7d4-448a-47a1-a53a-077aa24b5539.png)
