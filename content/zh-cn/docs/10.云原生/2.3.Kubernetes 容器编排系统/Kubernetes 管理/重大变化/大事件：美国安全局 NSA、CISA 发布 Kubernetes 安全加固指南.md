---
title: 大事件：美国安全局 NSA、CISA 发布 Kubernetes 安全加固指南
---

<https://mp.weixin.qq.com/s/sfYpGScDNgf2tsIYZOYAWw>

PDF原文：<https://www.nsa.gov/News-Features/Feature-Stories/Article-View/Article/2716980/nsa-cisa-release-kubernetes-hardening-guidance/>
中文翻译：<https://jimmysong.io/kubernetes-hardening-guidance/>

Kubernetes 最初由谷歌公司的工程师开发，随后由云原生计算基金会开源，它是当前最流行的容器协作软件。Kubernetes 主要用于基于云的基础设施内部，便于系统管理员使用软件容器部署新的 IT 资源。

Kubernetes 的攻击目标通常有以下三个原因：**数据窃取、计算能力窃取或拒绝服务**。传统上，数据盗窃是主要动机。

然而，由于 Kubernetes 和 Docker 模型和传统的单片软件平台之间存在巨大不同，因此很多系统管理员在安全配置 Kubernetes 方面问题颇多。多年来，多款密币挖掘僵尸网络都在攻击这类配置错误问题。

威胁行动者扫描互联网上被暴露的未认证的 Kubernetes 管理功能或者扫描在大型 Kubernetes 集群（如 Argo Workflow 或 Kubeflow）上运行的应用程序，获得对 Kubernetes 后端的访问权限，之后利用这种权限在受害者云基础设施上部署密币挖掘应用。这些攻击早在 2017 年初就已发生，而现在已发展为多个团伙为了利用同一个配置错误的集群而大打出手。

美国国家安全局 (NSA) 和网络安全与基础设施安全局 (CISA) 近日发布了一份 59 页的网络安全技术报告『Kubernetes 安全加固指南』。该报告详细介绍了对 Kubernetes 环境的威胁，并提供了配置指南以最大限度地降低风险。

CISA 和 NSA 发布的这份指南旨在为系统管理员提供关于未来 Kubernetes 配置的安全基线，以避免遭受此类攻击。另外除了基本的配置指南外，这份报告还详述了企业和政府机构可采取的基本缓解措施，阻止或限制 Kubernetes 安全事件的严重性，包括：

- 扫描容器和 Pods，查找漏洞或配置错误问题。
- 尽可能以最小权限运行容器和 Pods。
- 使用网络分割来控制攻陷事件造成的损害。
- 使用防火墙来限制不必要的网络连接和加密以保护机密性。
- 使用强认证和授权限制用户和管理员访问权限以及限制攻击面。
- 使用日志审计，以便管理员能够监控活动并收到关于潜在恶意活动的警报。
- 定期审计所有的 Kubernetes 设置并使用漏洞扫描确保安全风险得到控制，补丁已应用。

NSA 和 CISA 的这份指南侧重于安全挑战，并建议系统管理员尽可能强化他们的环境。NSA 发布此指南是支持美国国防部、国防工业基地和国家安全系统的使命的一部分。
