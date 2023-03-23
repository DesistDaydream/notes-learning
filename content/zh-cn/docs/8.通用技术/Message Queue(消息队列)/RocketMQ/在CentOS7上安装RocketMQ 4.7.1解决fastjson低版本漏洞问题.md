---
title: 在CentOS7上安装RocketMQ 4.7.1解决fastjson低版本漏洞问题
---

### 文章目录- 在 CentOS7 上安装 RocketMQ 4.7.1 解决 fastjson 低版本漏洞问题

- 前言

- 安装过程

- 下载和解压 RocketMQ

- 调低 RocketMQ 的 JVM 大小

- bin/runserver.sh

- bin/runserver.sh

- bin/tools.sh

- 启动 Name Server

- 启动 Broker

- 查看 RocketMQ 进程

- 测试 RocketMQ

- 测试发送消息和接收消息

- 关闭 RocketMQ

- 关闭 Broker

- 关闭 Name Server

- 修改 Name Server 的端口

- 安装 RocketMQ 控制台

- Troubleshooting

- 参考文档

- 扩展阅读

# 在 CentOS7 上安装 RocketMQ 4.7.1 解决 fastjson 低版本漏洞问题## 前言

阿里的 fastjson 的低版本（<=1.2.68）被爆出有安全漏洞，而 RocketMQ 4.7.0 使用了 fastjson 1.2.62，因此需要将 RocketMQ 升级到 RocketMQ 4.7.1(fastjson 1.2.69)。本文描述了在 CentOS7 上安装 RocketMQ 4.7.1 的过程，仅作为开发测试环境使用：

1. 单机部署，Name Server 和 Broker 都装在一台服务器上；

2. 调低了 RocketMQ 默认的 JVM 大小；

3. 没有设置开机自启动和守护进程。

## 安装过程

   服务器上已经安装了 OpenJDK 8，并设置了 JAVA_HOME。

### 下载和解压 RocketMQ

   在 RocketMQ 官网上找到下载 RocketMQ 4.7.1 的链接，下载和解压 RocketMQ： # 下载
   wget http://ftp.cuhk.edu.hk/pub/packages/apache.org/rocketmq/4.7.1/rocketmq-all-4.7.1-bin-release.zip # 解压
   unzip rocketmq-all-4.7.1-bin-release.zip # 安装到/usr/local/rocketmq
   mv rocketmq-all-4.7.1-bin-release /usr/local
   ln -s /usr/local/rocketmq-all-4.7.1-bin-release /usr/local/rocketmq
   1
   2
   3
   4
   5
   6
   7
   8
   91
   2
   3
   4
   5
   6
   7
   8
   9
   10
   11
   12
   13
   14
   15
   16
   17
   18
   19
   Plain Text

### 调低 RocketMQ 的 JVM 大小

   RocketMQ 的默认 JVM 太大，不适合在开发测试环境中使用，需要调低 JVM 大小。在 RocketMQ 的安装目录（本例为`/usr/local/rocketmq`)，查找 sh 脚本中的 JVM 参数设置：
   find . -name '\*.sh' | xargs egrep 'Xms'1
   Plain Text 需要修改以下 sh 脚本的 JVM 参数：

- bin/runserver.sh

- bin/runbroker.sh

- bin/tools.sh
  > 修改前记得先备份相应脚本，具体 JVM 大小根据实际情况设定。
  > bin/runserver.sh 修改前：
        JAVA_OPT="${JAVA_OPT} -server -Xms4g -Xmx4g -Xmn2g -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=320m"1
  Plain Text 修改后：
  JAVA_OPT="${JAVA_OPT} -server -Xms256m -Xmx256m -Xmn128m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=320m"1
Plain Textbin/runserver.sh修改前：
    JAVA_OPT="${JAVA_OPT} -server -Xms4g -Xmx4g -Xmn2g -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=320m"1
  Plain Text 修改后：
  JAVA_OPT="${JAVA_OPT} -server -Xms256m -Xmx256m -Xmn128m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=320m"1
Plain Textbin/tools.sh修改前：
    JAVA_OPT="${JAVA_OPT} -server -Xms1g -Xmx1g -Xmn256m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=128m"1
  Plain Text 修改后：
  JAVA_OPT="${JAVA_OPT} -server -Xms256m -Xmx256m -Xmn128m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=128m"1
  Plain Text

### 启动 Name Server # 后台启动

  nohup sh bin/mqnamesrv >/dev/null 2>&1 &1
  2
  Plain Text
  > 可以将启动 Name Server 命令保存为脚本，以方便下次启动。### 启动 Broker
  > 启动 Broker 时需要指定要连接的 Name Server：
        # 后台启动
        nohup sh bin/mqbroker -n localhost:9876 >/dev/null 2>&1 &1
  2
  Plain Text
  > 可以将启动 Broker 命令保存为脚本，以方便下次启动。### 查看 RocketMQ 进程 ps -ef | grep -v grep | grep rocketmq
        11
  2
  3
  Plain Text

### 测试发送消息和接收消息

  使用 RocketMQ 自带的消息生产者和消费者示例来测试发送消息和接收消息：
  export NAMESRV_ADDR=localhost:9876
  sh bin/tools.sh org.apache.rocketmq.example.quickstart.Producer
  sh bin/tools.sh org.apache.rocketmq.example.quickstart.Consumer
  1
  2
  31
  2
  3
  4
  5
  6
  7
  Plain Text

### 关闭 Broker sh bin/mqshutdown broker

  11
  2
  3
  Plain Text
  > 可以将关闭 Broker 命令保存为脚本，以方便下次关闭。### 关闭 Name Server sh bin/mqshutdown namesrv
        11
  2
  3
  Plain Text
  > 关闭 Name Server 前需要先关闭 Broker；可以将关闭 Name Server 命令保存为脚本，以方便下次关闭。## 修改 Name Server 的端口
  > RocketMQ Name Server 的默认端口为 9876，可以通过以下方法修改 Name Server 的端口：

1. 新增一个 Name Server 配置文件 namesrv.conf，保存内容为：listenPort=100761

2. 1

3. 启动 Name Server 时指定配置文件：nohup sh bin/mqnamesrv -c namesrv.conf >/dev/null 2>&1 &1

4. 1

5. 查看 RocketMQ 进程：ps -ef | grep rocketmq1

6. 1

7. 查看 RocketMQ Name Server 的端口号：netstat -tnlp | grep \<nameserver_pid>1

8. 1

9. 修改后 Broker 需要指定新的 Name Server 地址（端口）。

## 安装 RocketMQ 控制台 docker run -d --name console \

    -e "JAVA_OPTS=-Drocketmq.namesrv.addr=172.38.40.247:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false" \
    -p 8080:8080 \
    -t styletang/rocketmq-console-ng1
   2
   3
   4
   Plain Text 克隆 rocketmq-externals 项目，并编译 rocketmq-console。命令示例：
   git clone https://github.com/apache/rocketmq-externals.git
   cd rocketmq-externals/rocketmq-console
   mvn clean package -Dmaven.test.skip=true1
   2
   3
   Plain Text 将`target/rocketmq-console-ng*.jar` 放到和 RocketMQ 安装目录（本例为`/usr/local/rocketmq`)下。在 RocketMQ 安装目录下新建一个启动 RocketMQ 控制台的脚本来启动 RocketMQ 控制台：
   nohup java -jar rocketmq-console-ng\*.jar --server.port=8080 --rocketmq.config.namesrvAddr=localhost:9876 > /dev/null 2>&1 &1
   Plain Text 默认 RocketMQ 控制台不需要密码登录，请参考 RocketMQ 使用文档 进行配置。参见：

- <https://github.com/apache/rocketmq-externals>

- <https://github.com/apache/rocketmq-externals/tree/master/rocketmq-console>

## Troubleshooting

  问题 1: 启动 Name Server 和 Broker，或测试时报错`Please set the JAVA_HOME variable in your environment, We need java(x64)!` 但是系统已经安装了 OpenJDK8，并且已经设置了 JAVA_HOME。解决方法: 运行`which java`来查看 java 的路径，比如为`/usr/bin/java`。修改 bin/runserver.sh 和 bin/runbroker.sh 和 bin/tools.sh，注释掉校验 JAVA_HOME 语句，并明确指定 JAVA 路径： #[ ! -e "$JAVA_HOME/bin/java" ] && JAVA_HOME=$HOME/jdk/java
    #[ ! -e "$JAVA_HOME/bin/java" ] && JAVA_HOME=/usr/java #[ ! -e "$JAVA_HOME/bin/java" ] && error_exit "Please set the JAVA_HOME variable in your environment, We need java(x64)!"
  # export JAVA_HOME
  export JAVA="/usr/bin/java"1
  2
  3
  4
  5
  6
  Plain Text

## 参考文档- RocketMQ 安装部署教程详解

- RocketMQ 安装详细说明

## 扩展阅读- 十分钟入门 RocketMQ

- RocketMQ 核心设计理念

- RocketMQ 中文文档(非官方）
