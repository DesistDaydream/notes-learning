---
title: Standardized Glossary(标准化术语)
---

# 部署环境

Proof of Concept(概念验证，简称 poc)
Development(开发，简称 dev)
Test(测试，简称 test)
(灰度，简称 pre)
Production(生产，简称 pro)

# 服务器命名规范

> 参考：
> 
> - <https://codeantenna.com/a/VDpjUR86Hx>
> - [RFC-1178，为你的计算机选择一个名字](https://datatracker.ietf.org/doc/html/rfc1178)

## 推荐小规模集群 hostname 命名规范

规则: UN/LOCODE 码-机房标记(可选)-随机字符-系统和版本(可选)-云服务商缩写(可选)-环境(可选)-域名(可选)

简洁示例: suz-ba91 lxa-4f97

完整示例: gzu-spe-a904-rhel7-ecs-ctyun.prd.21cn.com can-4th-b69d-win2012-bms-ctyun.tst.21cn.com

过程:

- 查询 UN/LOCODE 城市代码, https://service.unece.org/trade/locode/cn.htm ;
- 查询操作系统发行版本, 执行命令: hostnamectl ;
- 截取 uuid.online 生成的 ID 4 位字符, http://www.uuid.online/ ;
- 根据上述规则和数据, 组成 hostname 名称 ;
- 执行命令: hostnamectl set-hostname .

说明: bms 代表物理机，ecs 代表虚拟机；系统和版本参见附录 3；环境缩写参见附录 4.

## hostname 命名规则总结

1. 公有云服务器
   规则: 云服务商缩写-IATA 城市代码-系统和版本-随机字符-域名
   示例: aws-tko-ctos7-44rr4.colinleefish.com
2. 标准化别名结构(Standardized CNAME Structure)
   规则: OrenTirosh 记忆编码项目特定选择的 1633 个词之一(只有 4-7 个字母),
   示例：crimson melody verona banjo
   DNS A Records 和 CNAME Records 示例：
   melody.example.com. A 192.0.2.12
   melody.lan.example.com. A 10.0.2.12
   melody.oob.example.com. A 10.42.2.12
   web02.prd.nyc.example.com. CNAME melody.example.com.
   说明: 适用于 1500 个左右全局服务器命名.
3. IBM hostname 命名格式
   a.标准域名服务器（DNS）主机名字符串，例如，xmachine.manhattan.ibm.com
   b.缺省的简短 DNS 主机名字符串，例如，xmachine
   c.数字 IP 地址，例如，127.1.255.3
4. ansible 方案
   规则: 项目名-环境-模块-ip
   示例: hnds-online-app-242
5. IDC 方案
   规则: {IDC}-{业务 bu}-\[{项目名}\[{编号}]]-{应用名}{机器编号}.vivo.lan
   示例: jsyz01-op-cmdb-mysql001.aa.lan
6. YouTube 方案
   规则: {数据中心}{区域}{节点}-in-{楼层}.{域名}
   示例:lga34s13-in-f14.1e100.net nuq04s29-in-f14.1e100.net
7. ServerDensity 方案
   示例: hcluster3-web1.sjc.sl.serverdensity.net
8. aws 方案
   示例：ec2-34-194-228-249.compute-1.amazonaws.com
9. 小规模服务器群方案
   规则：以单词 / 动物 / 人物命名，适合
   示例：lyre.riseup.net
   devianza.investici.org
   confino.investici.org
   perdizione.investici.org
   cryptonomicon.mit.edu
   Random Name Generator 网站：
   https://www.behindthename.com/random/
10. google 方案
    规则：必须符合 RFC 1035 要求
    示例: test.example.com
    说明：主机名必须包含一系列与正则表达式 [a-z](https://codeantenna.com/a/%5B-a-z0-9%5D*%5Ba-z0-9%5D)? 匹配的标签，各个标签用点连接起来。每个标签的长度为 1-63 个字符，整个序列不得超过 253 个字符。

## 常见问题

a. 不以用途(如 db1/nginx1 等)来命名的原因
因使用云主机时要突出的内容并不是用途, 故标记了地区/供应商/系统版本等.
b. 云主机供应商缩写非权威
目前，没有权威机构编制了云主机供应商的代号
c. 为什么选用 UN/LOCODE 码，而不是 IATA 城市代码
使用全拼太长, 取首字母易混淆, 例如 sz 无法区分深圳和苏州；
UN/LOCODE 码比 IATA 码能覆盖更多特定的位置，而且具有定义良好的标准。
d. 使用 5 个随机字符而不是 ip 地址后 3 位，或者使用 001 编号.
随机字符可以解决标记冲突的问题, 既足够使用又不至于太长.
ip 地址后三位有冲突的风险，而标记数字在服务器过少(如只有 001 编号)时, 显得突兀。

## 附录

a. UN/LOCODE 码城市代码示例

- 广州 can
- 深圳 snz
- 拉萨 lxa
- 杭州 haz
- 苏州 suz
- 贵州 gzu

b. 云主机供应商缩写

- Amazon Web Services aws
- Microsoft Azure maz
- Linode lnd
- DigitalOcean don
- Vultr vlt
- Bandwagon bwg
- 阿里云 aliyun
- Ucloud ucd
- 腾讯云 qcd
- 天翼云 ctyun

c. 系统和版本缩写与示例

- Red Hat Enterprise Linux rhel rhel7
- CentOS ctos ctos7
- Fedora fdr fdr7
- Oracle Linux orl orl7
- Ubuntu ubt ubt1604
- FreeBSD fbd fbd10
- CoreOS crs crs1068
- Windows win win2012

d. 软件应用环境缩写

- 开发环境 development dev
- 集成环境 integration intgr
- 测试环境 testing tst
- QA 验证 QA qa
- 模拟环境 staging stg
- 生产环境 production prd

e. 主机功能编号

- app Application Server (non-web)
- sql Database Server
- ftp SFTPserver
- mta Mail Server
- dns Name Server
- cfg Configuration Management (puppet/ansible/etc.)
- mon Monitoring Server (nagios, sensu, etc.)
- prx Proxy/Load Balancer (software)
- ssh SSHJump/Bastion Host
- sto Storage Server
- vcs Version Control Software Server (Git/SVN/CVS/etc.)
- vmm Virtual Machine Manager
- web Web Server
- con Console/Terminal Server
- fwl Firewall
- lbl Load Balancer (physical)
- rtr L3 Router
- swt L2 Switch
- vpn VPN Gateway
- pdu Power Distribution Unit
- ups Uninterruptible Power Supply

## 参考资料
[1. 怎么制定一套合适的服务器命名方案](https://cloud.tencent.com/developer/article/1114482)
\[2. Airline and Location Code Search]https://www.iata.org/en/publications/directories/code-search/?airport.search=shenzhen
[3. 我如何标记自己的公有云服务器实例](https://www.jianshu.com/p/9cxmD4)
[4. 我如何标记自己的公有云服务器实例](http://v.colinlee.fish/posts/how-do-i-mark-my-public-cloud-instances.html)
[5. 什么是 staging server](https://blog.csdn.net/blade2001/article/details/7194895)
[6. 软件生命周期中要经历的几种环境](https://blog.csdn.net/yy19890521/article/details/82345963)
[7. 概要文件、节点、服务器、主机和单元的命名注意事项](https://www.ibm.com/support/knowledgecenter/zh/SSFTBX_8.5.6/com.ibm.wbpm.imuc.doc/topics/cins_naming.html)
[8. RFC 1178 - Choosing a name for your computer](http://www.faqs.org/rfcs/rfc1178.html)
[9. Ops：命名规范](https://www.cnblogs.com/William-Guozi/p/Ops_nameRules.html)
[10. 对服务器 rDNS/Hostname 命名的一次探索](https://nova.moe/explore-in-server-rdns-and-hostname/)
[11. 主机名命名规范](https://www.cnblogs.com/kaishirenshi/p/10249072.html)
[12. 创建使用自定义主机名的虚拟机实例](https://cloud.google.com/compute/docs/instances/custom-hostname-vm#limitations)


原