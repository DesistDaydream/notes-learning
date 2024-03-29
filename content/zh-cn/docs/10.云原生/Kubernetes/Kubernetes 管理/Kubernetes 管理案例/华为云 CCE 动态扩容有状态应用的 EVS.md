---
title: 华为云 CCE 动态扩容有状态应用的 EVS
---

想要为华为云 CCE 集群中的 Statefulset 的 PVC 扩容

只修改 `statefulset.spec.volumeClaimTemplates.spec.resources.requests` 字段下的内容是无法真正扩容成功的，仅仅在下面的页面中容量显示会变化，但是真实硬盘并没有变化。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vcyh4v/1648526233603-107d40e6-e618-4e00-a2da-9c1952ce28f4.png)

同时还需要手动对 PVC 进行扩容，只有对 PVC 进行了扩容操作，华为云的 CSI 才会检测到变化并执行扩容操作

- 直接修改 PVC 中的 `pvc.spec.resources.requests` 字段中存储的容量
- 在华为云 web 控制台修改存储容量，如下图
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vcyh4v/1648526376264-fae48b42-637d-4300-8539-0483d17e5a06.png)
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vcyh4v/1648526400481-12a024c0-e4d9-46eb-ac8f-dc2580ca823b.png)

PVC 扩容完成后，真实的硬盘容量将会正常扩容
