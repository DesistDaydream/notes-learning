---
title: 华为 CDN
---

# CDN 整体部件名词解释

1. 上层平台术语
   1. CP：（ Content Provider） 内容提供者
   2. BOSS：(Business and Operation Support System)业务运营支撑系统
   3. CRS：（Content Resource System ）内容资源系统
   4. OMS 运营中心系统
      1. NMS：（Network Management Subsytem） 网络管理子系统
   5. CDNC：（ Content Delivery Network Center ）内容分发网络中心
   6. BMS：（ Business Management Subsystem） 业务管理子系统，
      1. 本节点业务管理、规则库管理(是否允许 HCS 缓存哪些域名的哪些内容等)、资源管理、系统管理等一般维护操作。
      2. 配置下发给本节点子系统。包括 TIS,HCS 等
2. 调度系统(Traffic Control)(系统只是一个概念一个解决方案，现场环境安装的全是组件，每个组件是组成一套系统的基础)
3. TCS：（Global Traffic Control System ）流量(全局)控制系统，该全局系统一般是总公司的系统，也包含下面说到的相关组件。当用户请求的时候，通过 TCS 的组件之间互相调度后再调度到 LTC 的组件。
4. LTC：（Local Traffic Control System ）本地(区域)流量控制系统。
5. DPI(特殊设备)：在 Cache 环境中，用户请求非签约模式的大文件内容时候，负责把请求镜像到 TIS-Mirror 上进行处理
6. DPI 非华为设备，一般是运营商与华为商量好，溯源中心的 DPI 部署在本省流量到外网的出口的位置上，当本地 HCS 需要去源站回源出去到外网的时候，会被该 DPI 劫持，然后把回源请求重定向到溯源中心的 HCS 上
7. TIS：Traffic Insight Subsystem 流量洞察子系统，负责流量镜像方式引流，包括 DNS 重定向、HTTP 重定向、P2P 重定向等方式。TIS 通过流量镜像方式获得终端用户上行请求的拷贝，对请求进行深度识别和处理后，通过重定向方式将终端用户引导到缓存服务器接受服务。
8. TIS-GSLB，用户请求内容时，通过 Local DNS 进行域名以及 cName DNS(DNS 别名)的解析，该组件通过 Local DNS 的地理位置，把离用户最近的 LBS 或者 LTC 的 IP 发送给客户(OTT,B2B 场景中)
9. GSLB：（Global Service Load Balance）全局服务负载均衡
10. LBS：（ Load Balancing Subsystem） 负载均衡子系统，接受用户的 HTTP 请求后，通过负载均衡调度策略，为用户选择一台 HCS(HMS)
11. SLB（Software Load Balancer），即软件负载均衡器，设置在一组功能相同或相似的集群服务器前端，对到达集群服务器组的流量进行合理分发，并在其中某一台服务器故障时，能将访问请求转移到其它可以正常工作的服务器的软件或网络设备
12. TIS-Mirror，用户请求内容时，收到 DPI 镜像过来的用户请求，来直接给用户分配 HCS(Cache 场景中)
13. TIS-DNS，用户请求内容时，Local DNS 本地 DNS 转发(Foward 方式)请求到该处，由 TIS_DNS 根据 BMS 规则进行域名匹配，把匹配到的 LBS(Cache 场景中)
14. TIS-HTTP：(Traffic Insight Subsystem HTTP)向 UE 返回 302(重定向)报文，报文中包含距离用户最近的边缘节点地市 HCS(HMS)前端 LBS 的虚 IP 地址。
15. TIS（Traffic Insight Suite）流量洞察套件，提供可观、可控、可预测的流量管理，用户行为预测、分析服务 ，是媒体使能平台（CDN）的媒体数据分析套件，负责用户访问日志的采集、处理、存储，对日志进行转换、分析热度模型和访问模型、历史访问及空间利用率统计，为 CDN 流量可视化和 CDN 现网资源运维、调配、优化提供分析支撑。
16. TA：（Traffic Analytics ）负责根据 CDN 的访问分析出 CDN 缓存系统的命中率、当前用户和内容的访问模型情况、当前最热的内容，从而为 CDN 运维，设备资源的配置、优化、业务规划部署、网路规划提供最佳的理论支撑。
17. UM：(Usage Mediator)用量采集器，负责从媒体服务器上接收用户访问的话单处理，并为外部报表系统提供数据源。也包括从 HMS、RRS 收集用户访问日志、调度日志并进行统一的日志规整，供流量洞察套件、VQE 套件实时分析。
18. ELB（Elastic Load Balancing）弹性负载均衡，用于实时监控系统的健康状况、下发配置变更以及安全防护。

- ELB 管理节点采用集群模式。
- 接收用户配置文件，并发送给 LVS 处理。
- 接收 LVS 上报的监控数据，管理负载均衡系统的健康状态。

1. 内容，资源缓存系统术语
2. HCS：（ Hybrid Cache Subsystem ）混合缓存子系统，缓存各种内容内容，并为用户提供内容，供用户下载，观看等，多台 HCS 内部可以互相查找内容，当用户在 HCS1 上没找到所需内容的时候，HCS1 会通知用户去 HCS2 上获取内容
3. CSS-BT：(Cache Subsystem of BT)BT 缓存子系统，负责缓存 BT 热门资源，在 UE 请求资源时返回给 UE。
4. CSS-XL：(Cache Subsystem of XunLei)XL 缓存子系统，负责缓存迅雷热门资源，在 UE 请求资源时返回给 UE。
5. DMFS：（DigitalMedia File Server） 数字媒体文件服务器对应 MDN
6. 网管系统相关术语(网管指的是网络中所有设备的管理者，可以管理各种设备的设备，比如 I2000,cacti 等)
7. DV：（DigitalView ）数字化视图，可接收接入设备上报的告警并展示，还可以对设备进行运维管理(包括部署等操作，需要安装扩展插件 DF 等)
8. DF:（Digita Foundry）数字化部署，把 DF 安装后接入 DV，可以通过 DV 的部署界面在 web 界面上部署各种应用以及软件
9. UOA：（Uniform Operation & Maintenance Agent）统一操作与维护(简称 OM)代理，统一网管代理是设备接入网管的统一接口，该代理模块负责将网元设备接入网管系统。网管与 UOA 之间的通信基于 SNMP（Simple Network Management Protocol）协议。
10. BMU：(Board Management Unit)单板管理单元，
11. 报表，数据收集等相关术语
12. Report 报表系统：里面包含 BDI，ISA 组件
13. BDI：（Big Data Integration）大数据集成
14. ISA：（Interactive Self Analytics）交互式自我分析
15. Zookeeper 是基于 Java 的开源软件，是 Hadoop 和 Hbase 的重要组件。主要作用是集群状态管理。MK 平台也是利用该组件来实现集群状态管理。
16. 媒体服务相关术语
17. CMI：Content Management Interface 内容管理接口，负责接收外部内容管理系统分发内容等请求，把请求转交给 MM
18. MM：Media Manager 媒体管理，即内容调度服务，处理 CMI 下发的任务，通过 MC 的健康检查机制，来决定调度哪台 HMS 去内容源获取内容
19. MC：Management Component 管理组件
20. 向网管提供设备管理接口，网管通过 MC 对 CDN 进行配置管理。
21. 监控并记录 CDN 设备健康状态，获取告警和日志信息。通过设备健康状态来决定由哪台 HCS(HMS)设备来获取资源
22. UM：Usage Management 使用管理，即鉴权，鉴别用户有没有权利来使用该内容。并且在订户体验点播、录播、网络时移以及直播（单播）业务时，记录订户的播放记录，生成播放记录文件，供报表子系统读取。
23. RRS：Request Routing System 路由请求系统，主要用于当用户请求媒体内容时，负责调度其中一台 HMS 给用户提供服务
24. HMS：Huawei Media Server 华为媒体服务器，一般情况下有很多很多台，存储媒体内容，给用户提供媒体内容
25. LRM：Live TV & Recording Management 直播与录像管理。直录播管理服务，负责直播频道和录播任务的管理、录制索引的管理。提升热点直播频道录制节目的分发速度。
26. GPM：Global Popularity Management 全局热度管理。负责 POP 点内分片热度数据的收集、统计、分析、排序、搬迁、空间均衡和热点均衡。
27. 其他
28. FS：（FusionSphere） 融合领域，华为的云操作系统
29. DNS： Domain Name Server 域名服务器
30. Local DNS
31. CP 授权 DNS,授权 DNS 通过 Cname(别名)的方式返回查询的域名的别名

其他

OPIN：（open Network integtation platform）开放的网络集成平台，电软的下一代平台， 在 ENIP2.0 的基础上继承发展，并扩展 ONIP3.0 的范围到云服务、IT 集成和互联网领域，为产品和解决方案提供开放的、易集成的、弹性可伸缩的、敏捷的核心平台。ONIP 是一系列平台的套件集合，以 SOA 架构为基础，服务和组件可以做为独立的交付件，按需应用到各解决方案中； 同时 ONIP 的领域平台集成了公共技术能力和领域部件，为解决方案提供集成的领域平台。

CIE 运维组件：（Carrier-grade Infrastructure application Environment）电信基础设施应用环境。包含 OMU,DMU。主要提供物理设备接入，虚拟设备接入，安装部署(OS,基础软件以及业务软件)等

1. OMU：Unit 操作维护单元，实现系统操作维护管理功能，提供 B/S 管理界面，实现对系统中设备管理、部署等功能的操作维护能力，包括配置、监控、告警、日志等 OM 能力。OMU 支持分布式管理多个 DMU
2. 提供 Web 界面，接收用户的软件安装请求。
3. 与 DMU 通信，将安装任务下发给 DMU。
4. DMU：Unit 域管理单元，实现系统全局设备管理功能。对所在区域内各类设备的集中管理，包括计算设备、存储设备、网络设备、虚拟资源的管理，及提供对业务板、虚拟机、存储设备的部署能力。同事 DMU 作为系统全局服务的协调单元，实现系统中各类自动化服务和 RAS 决策协调功能
5. 与 OMU 通信，接收下发的任务。
6. 通过 PXE、TFTP、NFS 等协议与单板建立安装通道，完成软件安装。

HACS：(High Availability Cluster Server)高可用集群服务器

GDR：（Geographical Disaster Recovery，地理容灾恢复）数据库双机

1. SCP：service control point 业务控制点
2. SAU：signaling access unit 信令接入单元
3. STP：signaling transfer point 信令转接点
4. MSC：mobile switching center 移动交换中心

UVP：(Universal Virtualization Platform)统一虚拟化平台，虚拟化类型，实现虚拟化的组件

USM：（Universal Server Manager）通用服务器管理，是一款采用 B/S（Browser/Server）架构的设备管理软件。USM 定位于提高华为 ATAE 设备的可管理性。USM 通过局域网管理 ATAE 设备，具有零客户端安装、网络任意点接入管理，安全、单一、集中的管理特点。使设备管理工作变得轻松、便捷、低成本。

UDB：（URL Database）用于记录 URL 热度，根据访问情况统计内容访问热度

# BMS 相关问题

升级规则库

测试规则库是否可以正常的方法

BMS 相关问题

/home/icache/rule/baseline 目录下可以查看规则库中设定的所有缓存的域名地址，就算 web 界面设置禁用了在这里也能看到，可以从这里总结规则库条目(即规则库里缓存的域名地址的数量)

测试单个资源在源站与本地 cache 之间的区别，通过放行用户访问请求达到测试目的

1. 业务管理—黑名单管理—创建—URL 关键字黑名单，选择节点，URL 中输入源站资源的 URL
2. 通过业务管理—配置报文下发监控 查看规则是否生效

升级规则库

详情见相关产品线产品文档，比如《(For Engineer) Internet Cache V200R002C34 产品文档 06》

注意事项：

1. 注意升级规则库之前，要往下翻页，只把第一页选中不行
2. 在准备下次升级之前，一定要查看升级结果，等待升级结果中全部已经完成再进下一次升级

产品分类

视频，网页，下载共三类

测试规则库是否可以正常的方法

1. tailf /home/icache/hcs/hcs/log/accesslog # 使用该命令查看日志内容是否会产生 TCP_HIT 200 GET 这样的状态信息日志，如果全是 MISS 则不正常
