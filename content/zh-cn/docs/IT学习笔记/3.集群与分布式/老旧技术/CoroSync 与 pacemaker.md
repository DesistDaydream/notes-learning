---
title: CoroSync 与 pacemaker
---

# AIS：OpenAIS（Application Interface Standard）应用接口规范

OpenAIS 提供了一种集群模式，包含集群框架。集群成员管理、通信方式、集群检测，但没有集群资源管理功能

组件包括：AMF，CLM，CPKT，EVT 等接口

CoroSync （The Corosync Cluster Engine）集群管理引擎，是 OpenAIS 的一个子组件

corosync+pacemaker

crmsh:suse 研发的命令行工具

pcs:RedHat 研发的命令行工具

# CoroSync 与 Pacemaker 的使用方法

CoroSync 可以直接用过源安装，自动解决依赖关系，包名就是 corosync

Cluster 中的各个 Node 通过相同的组播地址来组成集群，这些 Node 上会拥有相同的配置文件和 corosync-keygen 密钥，如果一个 Node 没有配置文件，则 Pacemaker 也无法控制这个节点，因为这个节点更本就不在这个集群中

配置文件在/etc/corosync.conf，该文件定义 CoroSync 底层是如何让 Cluster 内的各个 Node 互相通信的

1. totem { # 图腾，定义 Cluster 的 Messaging Infrastructure 层，各 Node 如何进行通信的

2. version：NUM # totem 的版本号

3. cluster_name：NAME # 定义集群名称

4. secauth：on|off # 定义是否打开安全认证功能，如果打开了，那么就需要使用 corosync-keygen 命令生成密钥，其余节点需要都拥有该密钥文件

5. threads：NUM # 定义线程数，默认为 0

6. interface { # 定义多个节点之间通过哪个接口，基于哪个多播地址，监听什么端口，完成多播通信的

7. ringnumber：0 # 一般保持为 0 就行

8. bindnetaddr:IP # 指明把组播地址绑定在哪个网络地址上，注意：是网络地址

9. mcastaddr:IP # 指明使用的组播地址

10. mcastport: NUM # 指明组播监听的端口号，默认是 5405

11. logging { # 指定日志系统信息

12. to_logfile: yes|no # 定义是否把信息写进日志

13. logfile: /PATH # 定义日志文件路径，默认为/var/log/corosync/corosync.log

14. logger_subsys { # 定义 log 子系统

15. debug:on|off # 定义 debug 是开还是关

16. service { # 该方法定义 pacemaker 作为 CoroSync 的一个插件来运行

         ver:    0

          name:   pacemaker

         use_mgmtd:yes

}

corosync-keygen # 如果打开了 secauth，那么需要使用该命令生成一个 corosync 的密钥文件，该命令使用后需要从/dev/random 设备中读取 1024bits 字节数，可以用手输入，或者安装某些程序

/var/log/XXXXX/corosync.log 日志的使用方法

1. grep -e "Corosync Cluster Engin" -e "configuration file" corosync.log # 查看 CoroSync 引擎是否正常启动

2. grep "TOTEM" corosync.log # 查看初始化成员节点通知是否正常发出

3. grep "pcmk_startup" corosync.log # 查看 pacemaker 是否正常启动了

# CRMSH 说明

语法结构

1. crm \[OPTIONS] \[COMMAND] \[SubCommand .....]

2. \#直接输入 crm 命令可以进入 crm shell 界面，在 crm 终端中继续输入相应命令进行操作,或者直接在 crm 后面跟后续命令再接输入的命令中的子命令中的子命令，比如 crm configure show，show 是 configure 的子命令

3. COMMAN 与 SubCommand

- status # 显示当前 Cluster 的状态，如图所示

- configure # 对 cluster 进行配置

- show # 显示配置信息

- show xml # 以 xml 方式显示配置信息

- verify # 验证配置完但是还没 commit 的信息是否正确

- commit # 提交被改变的配置到 CIB，把在内存中的配置信息写到磁盘中

- primitive|group|clone # 定义一个{主|组|克隆}资源，包括该资源代理的类型以及该类型的配置等等信息(注意：如果在定义资源的时候不进行 monitor 监控定义，那么在该资源出现异常无法提供服务的时候，这个资源不会从一个 node 自动切换到另一个 node)

- primitive webip ocf:heartbeat:IPaddr params ip=192.168.100.70 nic=enp0s3 cidr_netmask=16 op monitor interval=10s timeout=20s # 定义一个叫 webip 的资源，资源 params(参数)为：ip 地址，网卡，掩码；monitor(监控)信息：间隔时间为 10 秒，超时时间为 20 秒

- 间隔时间:每隔多少秒检查一下该资源是否可用

- 超时时间：当检测资源是否可用的时间超过多少秒时，则视为该资源失效，需要转移

- primitive webserver lsb:apache2 op monitor interval=10s timeout=20s # 定义主资源，资源名为：webserver；资源 params(参数)为：apache2 的进程；monitor(监控)信息：间隔时间为 10 秒，超时时间为 20 秒

- primitive webstore ocf:heartbeat:Filesystem params device="192.168.0.80:/web/htdocs" directory="/var/www/html" fstype="nfs" op monitor interval=20s timeout=40s op start timeout=60s op stop timeout=60s # 定义一个叫 webstore 的资源，资源 params(参数)为：设备,目录,文件系统类型,monitor(监控)信息:间隔时间为 20 秒，超时时间为 40 秒，资源启动超时时间为 60s，资源停止超时时间为 60 秒

- group webservice webip webserver # 定义一个叫 webservice 的组资源，把 webip 和 webserver 两个资源放到同一个资源组

- colocation|location|order # 定义{排列约束|位置约束|顺序约束}

- colocation webserver_with_webip inf: webserver webip # 定义 webserver 和 webip 排列约束名为 webserver_with_webip，规则为资源在同一个 node 上，约束值 inf 为无穷大(注意：当需要约束的资源同属一个组资源的时候，不用定义该约束)

- order webip_before_webserver Mandatory: webip webserver # 定义 webserver 和 webip 顺序约束名为 webip_before_webserver,规则为 webip 第一个启动，webip 不启动，则 webserver 资源也无法启动

- location webip_on_node2 webip rule 50: # uname eq k8s-node2 # 定义 webip 资源的位置约束名为 webip_on_node2，规则为 webip 资源倾向于在 node2 上运行，倾向性 50 分，默认分数为 0

- property # 配置 Cluster 的全局属性，配置在 show 展示的信息中 property 中的内容，cib 启动程序选项

- property default-resource-stickiness=50 # 定义资源默认对当前 node 的粘性，就算另一个 node 的位置约束数值更高，已经在当前节点的资源也不会转移走

- property stonith-enabled=false # 定义是否启用 stonith 设备，默认为 true 启动，如果没有 stonith 设备的话，该项配置为 false，否则集群无法使用，两节点集群需要配置为 false

- property no-quorum-policy=ignore # 定义在当前节点没有法定票数的时候策略是忽略，两节点集群需要配置为忽略

- edit # 修改一个 CIB objects(对象)，以文本形式修改 show 命令展示出的配置信息

- delete # 删除一个 CIB objects(对象)

- delete webip_before_webserver # 删除一个名字为 webip_before_websrver 的已经被定义了的对象(对象就是各种定义的名字)

- node # 对各 node 进行配置

- crm node online # 在执行该命令的 node 上线

- crm node standby # 在执行该命令的 node 下线变成备用

- resource # 对资源进行配置，启动，停止，重启，迁移等

- \[un]migrate # 迁移一个 resource 到其余 node\[迁移一个其余 node 的 resource 至当前 node]

- ra # 对资源代理的查看和控制

- classes # 查看有哪些资源代理类别

- list # 列出 classes 下可支持的所有资源代理

# 其他说明

资源类型分类：

1. primitive：主资源，只能运行于集群内的某单个节点，也称作 native

2. group：组资源，包含一个或多个资源，这些资源可通过“组”这个资源统一进行调度

3. clone：克隆资源，可以在同一集群内的多个节点运行多份克隆，每个节点都有一份该资源

4. master/slave:主从资源，在同一集群内部与两个节点运行两份资源，其中一个主一个从

HA 的逻辑架构：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yal1l0/1616132755527-d5c42b90-751c-4b93-a85c-c02ac1c9b1d2.jpeg)

1. 基础设施消息层 Messaging Infrastructure：在 HA 集群的每个节点中部署同一个应用程序以传递各节点的心跳信息以及集群事物信息，从而完成将这些节点构建成一个集群模式，但是该层对集群工作状态的生成的以及对各种事物做出决策，需要由集群管理层实现；

2. 集群成员关系层 Membership：该层功能可以被合并到 Mesaging Infrastructure 层，这一层主要是用来做决策的，根据法定票数来决定集群的运行，起作用主要是：添加或删除一个节点的时候做决策。在这层用到的一个进程是 CCM（Cluster Consensus Menbership）通过下层的心跳消息来决定其集群的拓扑。

3. 资源分配层 Resource Allocation:该层提供集群服务。

   1. CRM 集群管理器（cluster resource manager）：相当于一个由 pacemaker 提供的可操作的命令行工具，通过调用消息层的 API，为上层提供 API 调用接口和管理接口，系统管理员可以通过 CRM 实现管理操作，配置资源的各种约束(即资源运行在哪个 Node 上的定义)等，CRM 类似于一个董事长，可进行决策。

      1. LRM 本地资源管理器（Local Resource Manager）：为 CRM 提供支撑，负责接收 CRM 的指令交给 RA 来实现。

   2. DC 指派的协调员（Designated Coordinator）：集群中只有一个 DC，在集群中各个节点之间是平等的，但是需要一个协调所有关系的角色，这就是 DC，DC 所在节点即为主节点。对于集群的所有的计算、控制都在这个节点上进行，同时修改 CIB 通常也是在这个节点上进行然后同步到其他节点上去。

      1. PE 策略引擎（Policy Engine）：通过对这个集群上的所有节点进行计算，来进行决策。

      2. TE 转移引擎（Transition Engine）：执行决策，是进行监督决策的执行，但本身并不执行，如果是对本机的决策，则直接传给本地的

   3. CIB 集群信息库（Cluster Information Base）：要保证集群中所有节点的服务配置文件完全相同，在服务器之间同步数据的最好方式是使用 xml 格式的数据（其中包括了节点的个数，类型，总权重，以及各个服务器的配置信息等），这个文件就叫 CIB。

      1. cibadmin 命令管理 CIB

4. 资源层 Resources:这一层才是真正提供集群服务资源的。

   1. 资源代理 RA（Resource Agent）：负责实现为下层 LRM 提供启动，停止，监控的资源管理功能，RA 就是很多脚本，大体分为几个类型

      1. heartbeat V1

      2. LSB 遵循 linux 规范的 startup 脚本

      3. OCF 开放式集群格式，类似于 LSB，但是比其更有通用性

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/yal1l0/1616132755522-d36e0eff-da9d-482a-b363-1b6fedee53e2.jpeg)

- 消息层 # HA 的底层通讯模块，负责节点间的通信、心跳信息传递，集群事务消息传递。

- CRM 层 # HA 的资源调度模块，负责资源的配置、管理和调度。

- LRM #负责在每一个节点上执行 CRM 所发出的指令，类似于每个节点安装一个 LRM 然后接收控制节点上 CRM 发出的指令，也可以每个节点都安装 CRM，只不过 LRM 只能对自己所在的节点进行操作，而 CRM 可以给所有节点的 LRM 发送指令

- agent # HA 的资源代理，通过接受 LRM 层的 CRM 指令完成对所在节点的指定类型资源的管理。一个代理对应一类资源(图中的 A,B...X 指的就是资源)

HA 各层的解决方案：

1. Messaging Layer：

   1. heartbeat：分 V1,V2，V3 版本

   2. corosync：华为 HACS 用该方案

   3. Cman

   4. keepalived(工作方式完全不同于以上三种)

2. 资源分配层 Resource Allocation:CRM：cluster resource manager 集群资源管理器(管理)

   1. Heartbeat V1-HAresources：(配置接口：不通过命令管理，而是通过配置文件管理，配置文件的名称为 haresources)

   2. Heartbeat V2-crm：在各节点运行一个 crmd 进程，配置接口：命令行客户端程序 crmsh(crm 命令)；GUI 客户端：hb_gui

   3. pacemaker：从 Heartbeat V3 独立出来，可以以插件或独立方式运行.配置接口：CLI 接口:crmsh，pcs；GUI：hawk，LCMC，pacemaker-mgmt

   4. rgmanager：resource group manager 资源组管理器。配置接口：CLI：clustat，cman_tool；GUI：Conga

   5. CRM 分类

      1. HA-aware：资源自身可直接调用 HA 集群底层的 HA 功能

      2. 非 HA-aware：必须要借助于 CRM 完成在 HA 集群上实现 HA 功能

3. LRM：Local resource manager 本地资源管理器

   1. 由 CRM 通过子程序提供

4. RA:Resource Agent 资源代理，就是一个脚本，LRM 负责激活该脚本，CRM 负责决定在哪个资源上激活该脚本

   1. heartbeat legacy：heartbeat 传统类型的 RA，通常位于/etc/ha.d/haresources.d 目录下

   2. LSB：Linux Standard Base，Linux 基本库，通常位于/etc/rc.d/init.d/目录下，这些脚本至少接收 4 个参数{start|stop|restart|status}

   3. OCF：(probider)Open Cluster Framework，开发集群框架

   4. STONITH：专用于实现调用 STONIT 功能设备的资源，通常为 clone 类型

5. 解决方案的组合方式

   1. corosync+pacemaker

   2. cman+rgmanager

   3. 等

资源的约束(constraints)：启动完一个 HA 集群之后，HA 最终是要支撑业务的运行，业务是由各个节点的资源来提供的。一个资源更倾向于运行在哪个节点上，资源和资源之间是互斥的还是互相吸引的，以及资源和资源之间谁先谁后，是有一定的约束的，通过约束来指挥资源以最优的方式运行于何种节点上

1. 资源的约束关系

   1. location：位置约束，定义资源运行于某节点的倾向性，用数值来表示，倾向性指该资源对于该节点之前的感情，类似于女人喜欢钱，但是有另一个更有钱的男人追求她，她更倾向于前者，那么后追的人再有钱也不会跟着去的

   2. colocation：排列约束，定义资源彼此间是否可以运行在同一节点的倾向性

   3. order:顺序约束，定义资源在同一个节点上启动的先后顺序，先启动没起来，则后启动的也起不来；后启动的如果

心跳信息传递机制：通过 Messaging Layer 层的组件，主节点主动通报类似于"我还活着"的信息。因为多节点集群，只有在主动通报的情况下才会效率最大，否则每个节点都去查看主节点信息就太浪费资源了。该消息通过组播网络进行发送

udp 单播 组播 广播

组播地址：用于标识一个 ip 组播域，d 类地址空间，范围是 224.0.0.0-239.255.255.255

1. 永久组播地址：224.0.0.0-224.0.0.255

2. 临时组播地址：224.0.1.0-238.255.255.255，自己用

3. 本地组播地址：239.0.0.0-239.255.255.255，仅在特定本地范围内有效

HA 集群可能发生的故障以及应对情况

故障一：partitioned cluster 分裂集群，集群脑裂 cluster split brain：当 HA 的各节点之间无法互相通信的时候，无法报告自己的状态信息，每个节点都会认为自己是可用的，别的节点都挂掉了，那么每个节点都会抢夺资源，这就是脑裂

一般情况下，HA 集群的节点个数推荐为奇数个，这样才能保证产生分裂时，总有那么一方拥有的票数最大，vote system 即是决定该功能的系统

1. vote system:投票系统，HA 中的各节点无法探测彼此的心跳信息时，必然无法协调工作，而此种状态即为 partitioned cluster；必须按照少数服从多数的原则，Current DC(当前 DC)的状态分两种

   1. with quorum #当前具有法定票数，具有这种状态的即为主节点

   2. without quorum 当前节点没有拥有法定票数，如果出现该状态，则会出现几种处理方式

      1. stopped #停止

      2. ignore #忽略

      3. freeze #冻结

      4. suicide #自杀

      5. 那么就会被隔离，隔离方式分为两类；

         1. 节点：STONITH：Shooting The Others Node In The Head，直接关掉该节点电源

         2. 资源：fencing 隔离

故障二：正常情况下，当用户请求到路由器的时候，会把请求转发给 HA 中的浮动 ip，该浮动 ip 所在主机的 mac 地址会发送给路由器，路由器更新 arp 表，绑定浮动 ip 与该主机 mac 的关系；当该主机不可用时，用户请求再到路由器时，路由器查询自己的 arp 表，依然把请求转发到那台已经不可用的主机上，这时 HA 就不可用了。

- 解决方法：当 HA 集群主节点不可用。浮动 ip 转移以后，浮动 ip 所在的新节点会通过一个 arp 脚本 send_arp，广播自己的 mac 与该浮动 ip 的绑定关系，让集群内所有的其余节点和前端路由设备更新自己的 arp 表

HA 的工作模型：

1. A/P：两节点集群，Active，Passive 两个主备模型，一主一备，只有主节点运行 service 的所有 resource

2. A/A：两节点集群：Active，Active 两个双活模型，互备，比如：设备 1 运行一个 web service，设备 2 运行 mail service，web sevice 在设备 1 是主，设备 2 是备；mail service 在设备 1 是备，设备 2 是主；每个服务都会有自己独立的浮动 IP(类似于 VRRP 的多个虚拟 IP，每个虚拟 IP 当一个网段的网关)。一般情况下

3. N/M：n 个节点 m 个服务(不是资源，多个资源组成服务)，通常 n>m

4. N/N：n 个节点 n 个服务

### HA 案例一：HA Web Services

所需资源：ip，httpd，filesystem，floating ip

约束关系：

1. 使用“组”资源，或通过排列约束让资源运行于同一节点

2. 顺序约束：有次序地启动资源

解决方案组合：

heartbeat v2+haresources 或者 crm

配置 HA 集群的前提：

1. 节点之间时间必须同步，使用 NTP 协议实现

2. 节点间需要通过主机名互相通信，必须可以解析主机名至 IP 地址

   1. 建议名称解析功能使用 hosts 文件来实现

   2. 通信中使用的名字与节点名字必须保持一致：“uname -n”命令或者“hostname”命令展示出的名字一样

3. 考虑仲裁设备是否会用到，两节点必会用

4. 建立各节点之间的 root 用户能够基于密钥认证

   1. \#ssh-keygen -t rsa -P ' ' -f FILE

   2. \#ssh-copy-id -i /root/.ssh/id_rsa.pub root@HOSTNAME

5. 注意：定义称为集群服务中的资源，一定不能开机自启动，需要由 crm 统一管理

配置文件说明

1. /etc/ha.d 目录下

   1. ha.cf:主配置文件，定义各节点上的 heartbeat HA 集群的基本属性，以及集群事务信息如何传递的。

   2. authkeys：集群内节点间彼此传递消息时使用加密算法及密

   3. haresources：为 heartbeat v1 提供资源管理器配置接口，v1 版本专用配置接口，定义拥有哪些资源，以及节点资源运行的倾向性的

### HA 案例二：HA Mysql

## HA 中常用的术语

1. SPOF：Single Point of Failure 单点故障，HA 是由于 SPOF 的隐患所产生的高可用冗余方案

2. MTBF：mean time between failure 平均无故障时间

3. MTTR：mean time to repair 平均修复时长

4. A=MTBF/(MTBF+MTTR)，则 0

5. director：调度器

6. failover：失效转移，故障转移

7. failback：失效转回，故障转回

8. message layer

9. DC：designated coordinator # 指派的协调员，在分配为 DC（designated coordinator）的机器上创建 ha.cf

10. CIB：Cluster Information Base

异地多活模式，多个机房多活动中心，数据基于机房之间同步

Real Server：让 director 对其做健康状态检测，并且根据检测的结果自动完成添加或者移除等管理

1. 基于协议层次检查(lvs 不具备健康检测机制)

   1. ip:icmp

   2. 传输层：检测端口的开放状态

   3. 应用层：请求获取关键性资源，根据是否获取到资源来进行判断

2. 检查频度

3. 状态判断
