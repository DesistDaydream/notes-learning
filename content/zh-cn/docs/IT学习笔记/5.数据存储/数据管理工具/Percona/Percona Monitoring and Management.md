---
title: Percona Monitoring and Management
weight: 2
---

# 概述

> 参考：
> - [官网介绍](https://www.percona.com/software/database-tools/percona-monitoring-and-management)
> - [官方文档](https://www.percona.com/doc/percona-monitoring-and-management/2.x/index.html)
> - [原文链接](https://www.cnblogs.com/okchy/p/13605701.html)

分析慢查询的:<https://www.percona.com/blog/2020/10/07/how-to-find-query-slowdowns-using-percona-monitoring-and-management/>

基于 pmm2 去排查故障的官方文档:<https://www.percona.com/blog/2020/07/15/mysql-query-performance-troubleshooting-resource-based-approach/>

**Percona Monitoring and Management(简称 PMM)** 是一个用于管理和监控 MySQL、PostgreSQL、MongoDB 和 ProxySQL 性能的开源平台。它是由 Percona 与管理数据库服务、支助和咨询领域的专家合作开发的。

PMM 是一种免费的开源解决方案，您可以在自己的环境中运行它，以获得最大的安全性和可靠性。它为 MySQL、PostgreSQL 和 MongoDB 服务器提供了全面的基于时间的分析，以确保您的数据尽可能高效地工作。

PMM 平台基于支持可伸缩性的客户机-服务器模型。它包括以下模块:

PMM 客户机安装在您想要监视的每个数据库主机上。它收集服务器指标、一般系统指标和查询分析数据，以获得完整的性能概述。

PMM 服务器是 PMM 的中心部分，它聚合收集到的数据，并在 web 界面中以表格、仪表板和图形的形式显示这些数据。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512427-90c73526-d534-4d64-a402-c2c42373abb2.png)

模块被打包以便于安装和使用。假设用户不需要了解组成每个模块的具体工具是什么，以及它们如何交互。然而，如果您想充分利用 PMM 的潜力，内部结构是重要的。

PMM 是一种工具的集合，它被设计成可以无缝地协同工作。有些是由 Percona 开发的，有些是第三方开源工具。

PMM Server

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512441-02d3d028-1588-42fc-a8a4-4d1736674427.png)

PMM 服务器在作为中央监视主机的机器上运行。它通过以下方式作为设备分发:

\*可用于运行容器的 Docker 映像

\*可以在 VirtualBox 或其他管理程序中运行的 OVA(打开虚拟设备)

\*您可以通过 Amazon Web 服务运行的 AMI (Amazon Machine Image)

PMM 服务器包括以下工具:

\*查询分析(QAN)允许您在一段时间内分析 MySQL 查询性能。除客户端 QAN 代理外，还包括:

QAN API 是存储和访问运行在 PMM 客户机上的 QAN 代理收集的查询数据的后端。

QAN Web App 是一个可视化收集查询分析数据的 Web 应用程序。

\*Metrics Monitor 提供了对 MySQL 或 MongoDB 服务器实例至关重要的指标的历史视图。它包括以下内容:

Prometheus 是一个第三方时间序列数据库，它连接到运行在 PMM 客户机上的出口商，并汇总出口商收集的指标。

ClickHouse 是一个第三方的面向列的数据库，它促进了查询分析功能。有关更多信息，请参见 ClickHouse 文档。

Grafana 是一个第三方的仪表盘和图形生成器，用于将普罗米修斯在直观的 web 界面中聚合的数据可视化。

Percona 仪表板是由 Percona 为 Grafana 开发的一套仪表板。

所有工具都可以从 PMM 服务器的 web 界面(登录页面)访问。

PMM Client

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512448-678fe0f9-6e67-49e3-bd57-f740c4f96042.png)

每个 PMM 客户机收集关于一般系统和数据库性能的各种数据，并将这些数据发送到相应的 PMM 服务器。

PMM 客户端包括以下内容:

PMM -admin 是用于管理 PMM 客户机的命令行工具，例如，用于添加和删除您想要监视的数据库实例。

PMM -agent 是一个客户端组件，一个最小的命令行接口，它是负责提供客户端功能的中心入口点:它进行客户端身份验证，获取存储在 PMM 服务器上的客户端配置，管理导出程序和其他代理。

node_exporters 是一个收集一般系统指标的 Prometheus 端口。

mysqld_exporters 是一个收集 MySQL 服务器指标的 Prometheus 端口。

mongodb_exporters 是一个收集 MongoDB 服务器指标的 Prometheus 端口。

postgres\_端口是一个收集 PostgreSQL 性能指标的 Prometheus 端口。

proxysql_exporters 是一个收集 ProxySQL 性能指标的 Prometheus 端口。

为了使从 PMM 客户机到 PMM 服务器的数据传输更加安全，所有端口都能够使用 SSL/TLS 加密的连接，并且它们与 PMM 服务器的通信受到 HTTP 基本身份验证的保护。

参考：

端口：以下端口必须在 pmm server 和 client 之间开放;

pmm server 需要开放 80 或 443 端口用于 pmm client 访问 pmm web。

pmm client 端必须开放以下默认端口采集数据，可以通过 pmm-admin addc 命令进行修改。

42000 For PMM to collect genenal system metrics.

42001 This port is used by a service which collects query performance data and makes it available to QAN.

42002 For PMM to collect MySQL server metrics.

42003 For PMM to collect MongoDB server metrics.

42004 For PMM to collect ProxySQL server metrics.

# 部署 PMM

参考：官方文档

安装步骤

docker 部署 pmm 与 mysql 监控

安装 docker

yum install -y yum-utils device-mapper-persistent-data lvm2

yum-config-manager --add-repo <http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo>

yum makecache fast

yum -y install docker-ce

systemctl start docker

docker run hello-world

PMM Server：192.168.24.90

PMM Client：192.168.24.92

1.Docker 安装 PMM Server

这里使用的 Docker 标签适用于最新版本的 PMM 2(2.9.1)，但是您可以指定任何可用的标签来使用相应版本的 PMM 服务器。

度量收集消耗磁盘空间。PMM 需要为每个被监视的数据库节点提供大约 1GB 的存储空间，数据保留时间设置为一周。(默认情况下，数据保留时间为 30 天。)要减小 Prometheus 数据库的大小，可以考虑禁用表统计信息。

尽管一个受监控的数据库节点的最小内存量为 2 GB，但内存使用不会随着节点数量的增加而增加。例如，16GB 足够用于 20 个节点。

\#版本可自选

docker create -v /opt/prometheus/data -v /opt/consul-data -v /var/lib/mysql -v /var/lib/grafana --name pmm-data percona/pmm-server:2 /bin/true

2.启动

\#必须开启防火墙

docker run -d -p 80:80 -p 443:443 --volumes-from pmm-data --name pmm-server --restart always percona/pmm-server:2

端口默认是 80 ，如果 80 端口被占用，可改为其它端口号 比如 81

3.查看 docker 运行状态

docker ps -a

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512459-b4a07422-0748-4be5-bbdd-db952ae5cad9.png)

4.浏览器访问，地址一般是 http://ip 地址:端口，也可以直接输 ip 地址； 注意，一般端口默认为 80，默认用户名：admin，默认密码：admin

例：<http://192.168.24.90:80>

5.安装 pmm-client 客户端。

yum install <https://repo.percona.com/yum/percona-release-latest.noarch.rpm> -y

yum install pmm2-client -y

6.连接 PMM Server。

pmm-admin config --server-insecure-tls --server-url=\https://admin:admin@\<IP Address>:443

例：pmm-admin config --server-insecure-tls --server-url=https://admin:admin@192.168.24.90:443

注：PMM2 不需要像 PMM1 输入指定命令添加 Linux 主机监控

当你使用 pmm-admin config 配置了要监控的节点时，PMM2 从那时自动开始收集 Linux 指标。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512457-0d30d8f9-e92b-4ea3-b25a-2b7da2d86621.png)

7.登陆浏览器访问主机监控数据

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512459-6f938df2-c730-4504-bf71-e101ced88649.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512488-8a3a80a1-199f-451b-a1ef-b8915f8d5e15.png)

\#pmm-admin 管理命令

annotate \[\<flags>] \<text>

在 Grafana 图表中添加注释

config \[\<flags>] \[\<node-address>] \[\<node-type>] \[\<node-name>]

配置本地 pmm-agent

list \[\<flags>]

显示在此节点上运行的服务和代理

status

显示关于本地 pmm 代理的信息

summary \[\<flags>]

获取系统数据以进行诊断

add external --listen-port=LISTEN-PORT \[<flags>]

将外部监视添加

add mongodb \[<flags>] \[<name>] \[<address>]

监控 MongoDB

add mysql \[<flags>] \[<name>] \[<address>]

监控 MySQL

add postgresql \[<flags>] \[<name>] \[<address>]

监控 PostgreSQL

add proxysql \[<flags>] \[<name>] \[<address>]

监控 ProxySQL

register \[<flags>] \[<node-address>] \[<node-type>] \[<node-name>]

注册当前节点到 PMM 服务器

remove \[<flags>] <service-type> \[<service-name>]

从监控中删除服务

7.添加 mysql 监控。

MySQL 服务器添加指定权限用户

create user pmm@'%' identified by 'pmmpassword';

grant select,process,super,replication client on _._ to 'pmm'@'%';

grant update,delete,drop on performance_schema.\* to 'pmm'@'%';

flush privileges;

\#授权密码如报错：Your password does not satisfy the current policy requirements

set global validate_password_policy=LOW;

MySQL8.0 版本设置：set global validate_password.policy=LOW;

查询分析获得从 MySQL 中获取指标数据有两种可能的来源：慢查询日志和 Performance Schema

添加 Performance Schema 数据字典监控

pmm-admin add mysql --query-source=perfschema --username=pmm --password=pmmpassword ps-mysql

添加慢日志监控

pmm-admin add mysql --query-source=slowlog --username=pmm --password=pmmpassword sl-mysql

查看运行的服务

pmm-admin list

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512489-86f410ae-d99d-49a2-a505-2a9826ff1172.png)

9.pmm 服务器页面查看

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512483-5837728e-5d7e-4102-b343-9935ee58ac83.png)

点击 Query Analytics 进入 SQL 语句分析

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/hxu3v9/1616133512502-b408704d-70af-4efe-97b1-db073b551516.png)

10.MySQL 最佳配置

慢日志设置

如果你使用 Percona 分支版本，正确的慢查询日志配置将以最小的开销提供最多的信息。否则，如果支持请使用 PerformanceSchema。

按定义，慢查询日志应该只记录慢查询。这些查询的执行时间超过了特定的阈值。这个阈值由参数 long_query_time 指定。

在高负载的应用中，频繁快速的查询比罕见的慢速查询对性能的影响要大的多。为全面分析你的查询流量，设置 long_query_time 为 0，这样所有的查询语句都将被记录。

然而，记录所有的查询将消耗 I/O 带宽，并导致慢查询日志增长很快。为了限制记录到慢查询日志中的查询数量，使用 Percona 分支版本中的查询采样功能。

查询采样可能导致一些罕见的慢查询无法被记录到慢查询日志中。为了避免这种情况，使用 slow_query_long_always_write_time 参数指定哪类查询应该被忽略采样。也就是说，长时间的慢查询应该始终记录到慢查询日志中。

Performance Schema 设置

PMM 查询数据的默认源是慢查询日志。这在 MySQL5.1 及以后的版本中可用。从 MySQL5.6 版本（包括 Percona 分支版本 5.6 及以后版本），你可以选择从 Performance Schema 中解析查询数据，而不是慢查询日志。从 MySQL5.6.6 开始，默认启用 Performance Schema。

Performance Schema 不像慢查询日志那样有丰富的数据，但是它具有所有关键的数据，并且通常解析很快。如果您没有使用 Percona 分支版本（支持对慢查询日志采样），Performance Schema 是更好的选择。

启用 Performance Schema ，要将 performance_schema 参数设置为 ON:

SHOW VARIABLES LIKE 'performance_schema';

+--------------------+-------+

| Variable_name | Value |

+--------------------+-------+

| performance_schema | ON |

+--------------------+-------+

如果这个参数没有设置为 ON，在 my.cnf 配置文件中添加以下内容并重启 MySQL 服务。

\[mysql]

performance_schema=ON

如果您使用了自定义的 Performance Schema 配置，确认 statement_digest 消费者已经启用：

select \* from setup_consumers;

+----------------------------------+---------+

| NAME | ENABLED |

+----------------------------------+---------+

| events_stages_current | NO |

| events_stages_history | NO |

| events_stages_history_long | NO |

| events_statements_current | YES |

| events_statements_history | YES |

| events_statements_history_long | NO |

| events_transactions_current | NO |

| events_transactions_history | NO |

| events_transactions_history_long | NO |

| events_waits_current | NO |

| events_waits_history | NO |

| events_waits_history_long | NO |

| global_instrumentation | YES |

| thread_instrumentation | YES |

| statements_digest | YES |

+----------------------------------+---------+

15 rows in set (0.00 sec)

重要

Performance Schema 生产者在 MySQL5.6.6 及之后的版本中默认启用。它在 MySQL5.6 之前的版本中完全不可用。如果某些生产者没有被启用，您在 MySQLPerformanceSchemaDashboard 的 dashboard 中看不到相应的图。启用所有的生产者，在启动 MySQL 服务时设置 --performance_schema_instrument 选项为 '%=on'。

mysqld --performance-schema-instrument='%=on'

这个选项会带来额外的负载，请小心使用。

如果实例已经在运行，配置 QAN agent 从 Performance Schema 中收集数据：

1.打开 PMM Query Analytics dashboard。

2.点击 Settings 按钮。

3.打开 Settings 部分。

4.从收集下拉列表中选择 PerformanceSchema。

5.点击 Apply 保存更改。

如果您使用 pmm-admin 工具添加一个新的监控实例，使用 --query-sourceperfschema 选项：

使用 root 用户或者 sudo 命令执行以下命令

pmm-admin add mysql --username=pmm --password=pmmpassword --query-source='perfschema' ps-mysql 127.0.0.1:3306

更多信息，请执行 pmm-admin add mysql--help。

MySQL InnoDB 指标

为图形收集指标和统计信息会增加开销。您可以使用收集和绘制低开销的指标，在故障排除时启用高开销的指标。

InnoDB 指标提供了有关 InnoDB 操作的详细信息。尽管您可以选择捕获指定的计数器，但是即使始终启用它们，它们的开销也很低。启用所有的 InnoDB 指标，设置全局参数 innodb_monitor_enable 为 all:

SET GLOBAL innodb_monitor_enable

Percona 分支版本的特殊设置

默认情况下，并非所有 Metrics Monitor 的 dashboard 都可以用于所有 MySQL 分支和配置：Oracle 版，Percona 版或者 MariaDB。一些图形适用于 Percona 版本和专有的插件和额外的配置。

MySQL 用户统计信息（userstat）

用户统计信息是 Percona 分支版本和 MariaDB 分支版本的功能。它提供了用户活动、单个表和索引访问的信息。在某些情况下，收集用户统计信息可能会带来高昂的开销，所以请谨慎使用此功能。

启用收集用户统计信息，设置 userstat 参数为 1。

查询相应时间插件

查询响应时间分布是 Percona 分支版的可用功能。它提供了不同查询组的查询响应时间变化的信息，通常可以在导致严重问题之前发现性能问题。

启用收集查询响应时间：

1.安装 QUERY_RESPONSE_TIME 插件 mysql>INSTALL PLUGIN QUERY_RESPONSE_TIME_AUDIT SONAME'query_response_time.so';mysql>INSTALL PLUGIN QUERY_RESPONSE_TIME SONAME'query_response_time.so';mysql>INSTALL PLUGIN QUERY_RESPONSE_TIME_READ SONAME'query_response_time.so';mysql>INSTALL PLUGIN QUERY_RESPONSE_TIME_WRITE SONAME'query_response_time.so';

2.设置全局参数 query_response_time_stats 为 ON。 mysql>SET GLOBAL query_response_time_stats=ON;

相关信息：Percona 分支版官方文档

query_response_time_stats: <https://www.percona.com/doc/percona-server/5.7/diagnostics/responsetimedistribution.html#queryresponsetime_stats>

Response time 介绍: <https://www.percona.com/doc/percona-server/5.7/diagnostics/responsetimedistribution.html#installing-the-plugins>

logslowrate_limit

log_slow_rate_limit 参数定义了慢查询日志记录查询的比例。一个好的经验是每秒记录 100 个查询。例如如果您的 Percona 实例 QPS 为 10000，您应该设 log_slow_rate_limit 为 100,这样慢日志会记录每 100 个查询。

注意

当使用查询采样时，设置 log_slow_rate_type 为 query，以便它应用的是查询而不是会话。最好设置 log_slow_verbosity 为 full，以便在慢查询日志中记录每个记录的查询语句的最大的信息量。

logslowverbosity

- log_slow_verbosity 参数指定了慢查询日志中包含多少信息。最好设置 log_slow_verbosity 为 full，以便存储有关每个记录的查询语句的最大信息量。

slowqueryloguseglobal_control

默认情况下，慢查询日志只适用于新会话。如果希望调整慢查询日志设置并将这些设置应用于现有连接，请将 slow_query_log_use_global_control 设置为 all。

为 PMM 配置 MySQL8.0

MySQL8（在 8.0.4 版本中）改变了对客户端身份验证的方式。 default_authentication_pluging 参数设置为 caching_sha2_password。默认值的改变意味着 MySQL 的驱动需要支持 SHA-256 身份验证。另外，在使用 caching_sha2_password 时，必须对 MySQL8 的加密通道进行加密。

PMM 使用的 MySQL 驱动还不支持 SHA-256 身份验证。

为支持当前 MySQL 的版本，PMM 需要设置专有的 MySQL 用户。该 MySQL 用户应该使用 mysql_native_password 插件。虽然 MySQL 被配置支持 SSL 客户端，但是到 MySQL 服务器的连接没有加密。

有两种解决方法监控 MySQL8.0.4 及以上版本

1.更改你打算用于 PMM 的 MySQL 用户

2.改变 MySQL 的全局配置

更改 MySQL 用户

假设你已经创建了你打算用于 PMM 的 MySQL 用户，请使用以下方法更改：

然后，将此用户传递给 pmm-admin add 作为 --username 的参数值

这是首选的方法，因为这只会降低一个用户的安全性。

更改全局 MySQL 的配置

一种不太安全的方法是在添加监控前将 default_authentication_plugin 设置为 mysql_native_password。然后，重启 MySQL 服务，应用这个更改。
