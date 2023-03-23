---
title: rabbitmqctl 命令行工具详解
---

# rabbitmqctl 管理 RabbitMQ 节点的工具

官方文档：<https://www.rabbitmq.com/rabbitmqctl.8.html>

rabbitmqctl 是用于管理 RabbitMQ 服务器节点的命令行工具。它通过连接到专用 CLI 工具通信端口上的目标 RabbitMQ 节点并使用共享密钥（称为 cookie 文件）进行身份验证来执行所有操作。该命令主要的功能包括以下几点：

1. 停止节点运行

2. 获取节点状态、有效配置、健康检查

3. Virtual Hosts 管理

4. 用户和权限管理

5. Policy 管理

6. 查看 queues、connections、channels、exchanges 和 consumers 列表信息

7. 集群会员身份管理

**rabbitmqctl \[OPTIONS] COMMAND \[COMMAND_OPTIONS]**

OPTIONS

1. -q.--quiet # 安静模式输出，输出的信息减少

COMMAND 包括如下几大类：

1. Nodes # 节点管理

2. Cluster Managerment # 集群管理

3. Replication

4. Users Management # 用户管理

5. Access Control # 访问控制

6. Monitoring, observability and health checks

7. Parameters

8. Policies

9. Virtual hosts

10. Configuration and Environment

11. Definitions

12. Feature flags

13. Operations

14. Queues

15. Deprecated

COMMAND_OPTIONS 中有几个通用的

1. -p <VHost> # 用于当前适用的 COMMAND 所作用的 Virtual Host(虚拟主机)，VHost 为虚拟主机名称。默认为 / 这个 Virtual Host

# Nodes 类子命令

await_startup # Waits for the RabbitMQ application to start on the target node

reset # Instructs a RabbitMQ node to leave the cluster and return to its virgin state

rotate_logs # Instructs the RabbitMQ node to perform internal log rotation

shutdown # Stops RabbitMQ and its runtime (Erlang VM). Monitors progress for local nodes. Does not require a PID file path.

start_app # Starts the RabbitMQ application but leaves the runtime (Erlang VM) running

stop # Stops RabbitMQ and its runtime (Erlang VM). Requires a local node pid file path to monitor progress.

stop_app # Stops the RabbitMQ application, leaving the runtime (Erlang VM) running

wait # Waits for RabbitMQ node startup by monitoring a local PID file. See also 'rabbitmqctl await_online_nodes'

# Cluster Management 集群管理类子命令

await_online_nodes Waits for nodes to join the cluster

change_cluster_node_typeChanges the type of the cluster node

## cluster_status # 显示按节点类型分组的集群中的所有节点信息，以及当前正在运行的节点

force_boot Forces node to start even if it cannot contact or rejoin any of its previously known peers

force_reset Forcefully returns a RabbitMQ node to its virgin state

forget_cluster_node Removes a node from the cluster

join_cluster Instructs the node to become a member of the cluster that the specified node is in

rename_cluster_node Renames cluster nodes in the local database

update_cluster_nodes Instructs a cluster member node to sync the list of known cluster members from

# Replication

# Users Management 用户管理类子命令

## add_user # 在内部数据库中创建一个新用户

**rabbitmqctl add_user <UserNAME> <PASSWORD>**

EXAMPLE

1. rabbitmqctl add_user admin admin # 添加一个名为 admin 的用户，并设置其密码为 admin

authenticate_user # Attempts to authenticate a user. Exits with a non-zero code if authentication fails.

change_password # Changes the user password

clear_password # Clears (resets) password and disables password login for a user

delete_user # Removes a user from the internal database. Has no effect on users provided by external backends such as LDAP

## list_users # 列出所有用户和其标签

该命令效果如下：共两列，第一列为用户名，第二列为标签

    rabbitmq@hello-world-rabbitmq-server-0:/$ rabbitmqctl list_users
    Listing users ...
    user tags
    8E3s22eVBbIy3EINPFo0f8hBQ0FClORp [administrator]
    admin [administrator]
    test [monitoring]

set_user_tags # 设置指定用户的标签

**rabbitmqctl set_user_tags <UserNAME> <TAG>**

TAG 可以是任意值，但是有几个值在 rabbitmq 中具有特殊含义

1. administrator # 可登录管理控制台（启用 management plugin 的情况下），查看所有的信息，并且可以对用户、策略（policy）进行操作；

2. monitoring # 可登录管理控制台（启用 management plugin 的情况下），同时可以查看 rabbitmq 节点的相关信息（进程数、内存使用情况，磁盘使用情况等）；

3. policymaker # 可以登录管理控制台（启用 management plugin 的情况下），同时可以对策略（policy）进行操作；

4. management # 仅可登录管理控制台（启用 management plugin 的情况下），无法看到节点信息，也无法对策略进行管理；

5. 其他任意值 # 无法登录管理控制台，通常就是普通的生产者和消费者，这种 TAG 仅仅作为标识符。

EXAMPLE

1. rabbitmqctl set_user_tags admin administrator # 设置 admin 这个用户的标签为 administrator，让 admin 用户具有管理员权限。

2. rabbitmqctl set_user_tags admin # 移除 admin 这个用户的所有标签

# Access Control 访问控制子命令

## clear_permissions # 撤销指定用户关于 vhost 的权限

**rabbitmqctl clear_permissions \[-p VHOST] UserNAME**

EXAMPLE

1. rabbitmqctl clear_permissions -p test admin # 撤销 admin 用户关于 / 这个 vhost 的所有全新

clear_topic_permissions # Clears user topic permissions for a vhost or exchange

list_permissions # Lists user permissions in a virtual host

list_topic_permissions # Lists topic permissions in a virtual host

list_user_permissions # Lists permissions of a user across all virtual hosts

list_user_topic_permissions # Lists user topic permissions

## list_vhosts # 列出所有 virtual hosts

## set_permissions # 设置指定用户关于 vhost 的权限

**rabbitmqctl set_permissions \[-p VHost] UserName CONF WRITE READ**

一个用户对于 vhost 来说，有 CONF(配置)、WRITE(写)、READ(读) 这三个权限。可以使用正则表达式。

EXAMPLE

1. rabbitmqctl set_permissions -p / admin '._' '._' '.\*' # 为 admin 用户授予关于 / 这个 vhost 下所有资源的所有权限。

set_topic_permissions # Sets user topic permissions for an exchange

# Monitoring, observability and health checks 监控，可观察性以及健康检查子命令

list_bindings # Lists all bindings on a vhost

list_channels # Lists all channels in the node

list_ciphers # Lists cipher suites supported by encoding commands

list_connections # Lists AMQP 0.9.1 connections for the node

list_consumers # Lists all consumers for a vhost

## list_exchanges # 列出交换器的详细信息

**rabbitmqctl list_exchanges \[-p VHost] \[ExchangeInfoItem]**

ExchangeInfoItem # 该参数用于指示要在结果中包括交换器的哪些信息。多个信息使用逗号或空格分隔。默认显示 exchanges 的 name 与 type 信息。

可用的信息有如下几个：

1. name、type、durable、auto_delete、internal、arguments、policy

EXAMPLE

list_hashes # Lists hash functions supported by encoding commands

## list_queues # 列出队列及其属性

**rabbitmqctl list_queues \[-p VHost] \[QueueInfoItem]**

QueueInfoItem 用于指示要在结果中包括哪些队列信息项。结果中的列顺序将与参数的顺序匹配。不使用该参数时，默认显示 name 与 message 信息。

QueueInfoItem

1. name # 队列的名称

2. pid # 队列的 Erlang 进程标识符(其中包含队列在哪个节点)

3. messages # 队列深度，即 ready 和 unacknowledged 两种状态的消息总和。

EXAMPLE

1. rabbitmqctl list_queues # 列出 / vhost 下的队列名及其消息数量，效果如下

    rabbitmq@hello-world-rabbitmq-server-0:/$ rabbitmqctl list_queues
    Timeout: 60.0 seconds ...
    Listing queues for vhost / ...
    name messages
    test 2

list_unresponsive_queues # Tests queues to respond within timeout. Lists those which did not respond

ping # Checks that the node OS process is up, registered with EPMD and CLI tools can authenticate with it

report # Generate a server status report containing a concatenation of all server status information for support purposes

schema_info # Lists schema database tables and their properties

## status # 显示节点的状态

# Runtime Parameters and Policies 运行时参数和策略相关子命令

clear_global_parameter Clears a global runtime parameter

clear_parameter Clears a runtime parameter.

list_global_parameters Lists global runtime parameters

list_parameters Lists runtime parameters for a virtual host

set_global_parameter Sets a runtime parameter.

set_parameter Sets a runtime parameter.

clear_operator_policy Clears an operator policy

clear_policy Clears (removes) a policy

list_operator_policies Lists operator policy overrides for a virtual host

list_policies Lists all policies in a virtual host

set_operator_policy Sets an operator policy that overrides a subset of arguments in user policies

set_policy Sets or updates a policy

# Virtual hosts 虚拟主机相关子命令

## add_vhost # 创建一个 vhost

**rabbitmqctl add_vhost VHOST**

EXAMPLE

1. rabbitmqctl add_vhost test # 创建一个名为 test 的 vhost

clear_vhost_limits # Clears virtual host limits

## delete_vhost # 删除一个 vhost

> 注意：删除 vhost 将删除其所有交换，队列，绑定，用户权限，参数和策略。

**rabbitmqctl delete_vhost VHOST**

EXAMPLE

1. rabbitmqctl delete_vhost test # 删除一个名为 test 的 vhost

list_vhost_limits # Displays configured virtual host limits

restart_vhost # Restarts a failed vhost data stores and queues

set_vhost_limits # Sets virtual host limits

trace_off #

trace_on #

Configuration and Environment

Definitions

Feature flags

Operations

# Queues Operations 队列的相关操作子命令

## delete_queue # 删除一个队列

**rabbitmqctl delete_queue QueueName \[-p VHost] \[OPTIONS]**

OPTIONS

1. --if-empty,-e # 如果队列为空，则删除。(没有准备好传递的消息)

2. --if-unused,-u # 仅当队列没有消费者时才删除。

## purge_queue # 清洗一个队列(删除队列中所有的消息)

# Deprecated
