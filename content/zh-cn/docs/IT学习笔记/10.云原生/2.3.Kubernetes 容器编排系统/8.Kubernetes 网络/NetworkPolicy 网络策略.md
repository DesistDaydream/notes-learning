---
title: NetworkPolicy 网络策略
---

如果想要正常使用 NetworkPolicy 对象，则集群的网路插件必须要支持该功能

使用 NetworkPolicy 对象来定义网络策略(注意：该定义的网络策略只能适用于某个 namesapce 下的 Pod，可以在 metadata 中定义生效的 namespace)

一个网络策略中包含两中类型的规则，规则用于控制数据流量(Traffic)，

1. ingress #入口(入站)规则，列出所选择的 Pod，把入口规则应用其上

   1. port，定义允许还是拒绝哪些端口，ingress 规则中为，从外面来的可以从自己哪个端口进来

   - from #ingress 定义 from(目标过来的规则，即从哪来的可以入)

     - ipBlock #定义从哪来的 IP 段可以进来

       - except #定义完 ipBlock 后，定义 IP 段内的某些事不能进来

     - namespaceSelector #定义来自哪些 namesapce 的可以进来

     - podSelector #定义来自哪些 Pod 的可以进来

2. egress #出口(出站)规则，列出所选择的 Pod，把出口规则应用其上

   1. port，定义允许还是拒绝哪些端口，egress 规则中，到哪的端口可以出去

   - to #egress 定义 to(到目标的规则，即到哪的可以出)

     - ipBlock #定义到哪的 IP 段可以出去

       - except #定义完 ipBlock 后，定义 IP 段内的某些事不能出去

     - namespaceSelector #定义到哪些 namesapce 的可以出去

     - podSelector #定义到哪些 Pod 的可以出去

3. 规则说明

   1. 若不定义该规则，则默认全部允许

   2. 若定义了规则，那么规则默认拒绝，定义了哪些(ip、namespace、pod)那么定义的这些就允许

   3. 若定义的规则且内容为{}，说明全部允许

   4. 对 Pod1 使用入站规则设定 Pod2 不能进，那么 Pod2 无法 ping 通 Pod1，但是 Pod1 可以 ping 通 pod2

   5. 注意：出入规则为单向限制，仅对请求报文限制，响应报文不做限制。比如，对我使用入站规则禁止你 ping 我，但是我依然可以 ping 你，但是到某 IP 还是可以的，除非禁用了到那些 IP 的 Pod，那么就不能出也不能进了。该策略跟那种限制了出和入其中一个就都不行的不一样

网络策略中还包括规则的生效范围、(即把规则应用于 Pod 上和定义两个规则是否生效）。

1. podSelector #通过标签选择器，选择出把该策略的规则应用于哪些 Pod 上(默认为{}空，则对所有 Pod 生效)

2. policyTypes #列出网络策略涉及的规则类型。列出哪个规则，哪个规则生效，生效后默认拒绝所有，需要再定义规则以便允许，如果不定义规则那么什么都访问不了；如果两个规则都没列出，那么后面定义了哪个规则，则哪个规则生效。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zghpo8/1616117994284-0bc4626b-55c8-4f72-8e67-3e20bbe787c3.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zghpo8/1616117994290-5d61a239-4ce3-4fcd-94d8-a0f102cdb416.jpeg)
