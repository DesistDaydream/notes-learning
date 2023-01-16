---
title: Zone 与 Domain
---

# Delegation(授权)

Zone 是通过授权实现的，而授权，主要授予的就是 Domain 的管理权。

管理域的组织可以将域进一步划分成子域。每个子域都可以被授权给其他组织来管理，这意味着这些被授权的组织得负责维护子域中所有的数据。他们可以自由地 改变域中的数据，甚至将子域进一步划分成更多的子域，然后再授权给其他组织管理。父域仅仅保留指向子域数据来源的指针，这样父域便可将查询者引到该处。例如，stanford.edu 域被授权给在斯坦福大学的人，由他们负责管理校园网络。

比如，我想买一个域名 desistdaydream.com，我就需要去找管理 com 域的组织购买，我付给他们钱后，他们就给我授权，让我可以管理 desistdaydream.com 域。此时，该域的管理权则由我全权负责，而我还可以对 desistdaydream.com 域再次划分，比如 a.desistdaydream.com、b.desistdaydream.com、c.desistdaydream.com 等等，然后将这些域授权出去。

> 同时，我还可以只授权一个域，比如我只授权某人管理 a.b.desistdaydream.com 域，那么某人则无法将 a.b.desistdaydream.com 域再进行子域划分。

**所以，授权就是指将管理子域的责任交给另一个组织的行为。**

**zone **和 **domain **的区别

1. zone 是同 Delegation(授权) 联系在一起的，为了管理上的方便，我们把域的某部分授权出去让别人代为管理，这部分就是一个 zone 。为什么说是为了管理上的方便呢？因为这样一个很大的域就可以实现分散管理，而不是集中由一两台服务器来管理。而 zone 的划分就是通过 “授权机制”来实现的。这也是 设计 DNS 系统得初衷。
2. 并不能说 domain 就比 zone 大，反过来也一样。例如 edu 域可以包含多个 zone ：berkeley.edu 、purdue.edu 。但 edu 也可以看成是根域 "." 下的一个被授权出去的 zone ，它含有 berkeley.edu 、purdue.edu 等几个域。
3. 域是以域名进行分界的，而 zone 是以授权范围来定界的。一个 zone 有可能跨域多个域名。例如 berkeley 域是所有以 berkeley.edu 结尾的域名空间；而 edu zone 可以包括 berkey 和 purdue 这两个域，都统一归 edu 这个 zone 管理。
4. 一个域和一个 zone 可能具有相同的域名，但包含的节点却不同。例如使用了授权的域
5. Name Server 在加载数据时是以 zone 为单位，而不是以 Domain 为单位。

总结：
domain 这是从逻辑上进行划分，体现域名的树性结构，根域、com 域、edu 域等；

上图中 虚线 内就是一个 Zone

Note：ZONE 文件是 DNS 上保存域名配置的文件，对 BIND 来说一个域名对应一个 ZONE 文件，现以 abc.com 的 ZONE 文件为例展开(该 ZONE 存在于权威 DNS 上)。
