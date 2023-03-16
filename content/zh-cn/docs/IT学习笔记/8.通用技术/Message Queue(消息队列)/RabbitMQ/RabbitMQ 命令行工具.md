---
title: RabbitMQ 命令行工具
---

# 概述

RabbitMQ 主要有四个命令行工具 Commend Line Tools，不同的命令行适用不用的场景

- rabbitmqctl：负责服务管理和进行操作
- rabbitmq-diagnostics：负责系统诊断和健康检查
- rabbitmq-plugins：负责插件管理
- rabbitmqadmin：用来操作 HTTP API

# rabbitmqctl 管理 RabbitMQ 节点

详见此处

# rabbitmq-diagnostics 负责系统诊断和健康检查

# rabbitmq-plugins 负责插件管理

rabbitmq-plugins \[--node ] \[--timeout ] \[--longnames] \[--quiet] \[]

Available commands:

Help:

autocomplete Provides command name autocomplete variants

help Displays usage information for a command

version Displays CLI tools version

Monitoring, observability and health checks:

directories Displays plugin directory and enabled plugin file paths

is_enabled Health check that exits with a non-zero code if provided plugins are not enabled on target node

Plugin Management:

1. disable # 禁用一个或多个插件
2. enable # 启用一个或多个插件
3. list # 列出插件及其状态
4. set # 启用一个或多个插件，并禁用其他的插件。

# rabbitmqadmin 用来操作 HTTP API

官方文档：<https://www.rabbitmq.com/management-cli.html>

参考：<https://www.cnblogs.com/xishuai/p/rabbitmq-cli-rabbitmqadmin.html>

rabbitmqadmin 是一个 python 脚本。从 rabbitmq 的 http api 接口可以直接获取该脚本

1. curl -L 172.38.40.212:45672/cli/rabbitmqadmin -o /usr/local/bin/rabbitmqadmin
2. chmod 755 /usr/local/bin/rabbitmqadmin

生成命令行配置文件

rabbitmqadmin 默认使用 ~/.rabbitmqadmin.conf 作为配置文件。配置文件用来指定 rabbitmqadmin 连接 rabbitmq http api 的信息，比如 ip、port、登录用户名和密码、所使用的 vhost 等等。注意：配置文件中的内容可以通过命令行参数代替。

1. cat > /root/.rabbitmqadmin.conf << EOF
2. \[default]
3. hostname = 172.19.42.214
4. port = 45672
5. username = admin
6. password = admin
7. declare_vhost = /
8. vhost = /
9. EOF

生成命令补全

1. echo "source <(rabbitmqadmin --bash-completion)" >> /root/.bashrc

rabbitmqadmin \[OPTIONS] COMMAND

OPTIONS

1. --config=CONFIG, -c CONFIG # 指定 rabbitmqadmin 的配置文件，默认为: ~/.rabbitmqadmin.conf
2. --node=NODE, -N NODE # 指定配置文件中节点(即配置中间中 \[] 内的值就是节点名) \[default:'default' only if configuration file is specified]
3. --host=HOST, -H HOST # connect to host HOST \[default: localhost]
4. --port=PORT, -P PORT # connect to port PORT \[default: 15672]
5. --path-prefix=PATH_PREFIX # use specific URI path prefix for the RabbitMQ HTTP API. /api and operation path will be appended to it.(default: blank string) \[default: ]
6. --vhost=VHOST, -V VHOST # 指定要连接的 VHost \[default: all vhosts for list,'/' for declare]
7. --username=USERNAME, -u USERNAME # 指定连接时所使用的用户名，默认为 guest
8. --password=PASSWORD, -p PASSWORD # 指定连接时所使用的密码，默认为 guest
9. --base-uri=URI, -U URI # connect using a base HTTP API URI. /api and operation path will be appended to it. Path will be ignored.--vhost has to be provided separately.
10. --quiet, -q # suppress status messages \[default: True]
11. --ssl, -s # connect with ssl \[default: False]
12. --ssl-key-file=SSL_KEY_FILE # PEM format key file for SSL
13. --ssl-cert-file=SSL_CERT_FILE # PEM format certificate file for SSL
14. --ssl-ca-cert-file=SSL_CA_CERT_FILE # PEM format CA certificate file for SSL
15. --ssl-disable-hostname-verification # Disables peer hostname verification
16. --ssl-insecure, -k # Disables all SSL validations like curl's '-k' argument
17. --request-timeout=REQUEST_TIMEOUT, -t REQUEST_TIMEOUT # HTTP request timeout in seconds \[default: 120]
18. --format=FORMAT, -f FORMAT # format for listing commands - one of \[raw_json,pretty_json, tsv, long, table, kvp, bash] \[default:table]
19. --sort=SORT, -S SORT # sort key for listing queries
20. --sort-reverse, -R # reverse the sort order
21. --depth=DEPTH, -d DEPTH # maximum depth to recurse for listing tables \[default:1]
22. --bash-completion # Print bash completion script \[default: False]
23. --version # Display version and exit

## 显示信息 COMMAND

list connections \[...] # 显示连接信息

list channels \[...]

list consumers \[...]

list exchanges \[...]

list queues \[...]

list bindings \[...]

list users \[...] # 显示所有用户

list vhosts \[...] # 显示所有 vhosts

list permissions \[...]

list nodes \[...] # 显示所有节点

list parameters \[...]

list policies \[...]

list operator_policies \[...]

list vhost_limits \[...]

show overview \[...]

## 创建资源 COMMAND

### exchange 声明交换器

rabbitmqadmin declare exchange name=... type=... \[auto_delete=... durable=... internal=... arguments=...]

EXAMPLE

1. rabbitmqadmin --vhost=test declare exchange name=test.topic type=topic #

### queue 声明队列

rabbitmqadmin declare queue name=... \[auto_delete=... durable=... arguments=... node=... queue_type=...]

OPTIONS

1. durable= # 该队列是否持久化。默认为 true
2. queue_type= # 该队列类型。默认为 classic。TYPE 可用值为 classic、quorum
3. Note：quorum 类型下，durable 只能为 true。

EXAMPLE

1. rabbitmqadmin --vhost=test declare queue name=test-1 queue_type=quorum #

### binding 绑定一个交换器和队列

rabbitmqadmin declare binding source=... destination=... \[destination_type=... routing_key=... arguments=...]

OPTIONS

EXAMPLE

1. rabbitmqadmin --vhost=test declare binding source=test.topic destination=test routing_key=my.# #

### vhost 声明虚拟主机

rabbitmqadmin declare vhost name=... \[tracing=...]

OPTIONS

EXAMPLE

1. rabbitmqadmin declare vhost name=test #

declare user name=... password=... OR password_hash=... tags=... \[hashing_algorithm=...]

declare permission vhost=... user=... configure=... write=... read=...

declare parameter component=... name=... value=...

declare policy name=... pattern=... definition=... \[priority=... apply-to=...]

declare operator_policy name=... pattern=... definition=... \[priority=... apply-to=...]

declare vhost_limit vhost=... name=... value=...

## 删除、清理资源 COMMAND

delete exchange name=...

delete queue name=...

delete binding source=... destination_type=... destination=... \[properties_key=...]

delete vhost name=...

delete user name=...

delete permission vhost=... user=...

delete parameter component=... name=...

delete policy name=...

delete operator_policy name=...

delete vhost_limit vhost=... name=...

close connection name=...

purge queue name=...

## Broker 定义 COMMAND

export

import

## 发布与消费消息 COMMAND

### publish 发布消息

rabbitmqadmin publish routing_key=... \[payload=... properties=... exchange=... payload_encoding=...]

EXAMPLE

1. rabbitmqadmin publish routing_key=test payload="hello world" #
2. rabbitmqadmin publish routing_key=my.test exchange=my.topic payload="hello world by my.test" #
3. rabbitmqadmin publish routing_key=my.test.test exchange=my.topic payload="hello world by my.test.test" #

### get 消费消息

rabbitmqadmin get queue=... \[count=... ackmode=... payload_file=... encoding=...]

OPTIONS

1. count= # 指定要获取的最大消息数。默认为 1，即只获取队列中最先进入的消息。
2. ackmode= # 指定是否从队列中删除消息。默认为 ack_requeue_true
3. ack_requeue_true # 重新排队，不删除消息
4. ack_requeue_fale # 删除消息

EXAMPLE

1. rabbitmqadmin get queue=test # 获取 test 队列中的消息
2. rabbitmqadmin get queue=test ackmode=ack_requeue_false # 获取 test 队列中的消息，并将获取到的消息删除

- If payload is not specified on publish, standard input is used

- If payload_file is not specified on get, the payload will be shown on

  standard output along with the message metadata

- If payload_file is specified on get, count must not be set
