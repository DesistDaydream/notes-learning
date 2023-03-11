---
title: "整合RocketMq提示RemotingTooMuchRequestException: sendDefaultImpl call timeout"
---

在云服务器上安装 RocketMq 后，项目整合测试

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xuvc3b/1616130360509-382458c2-a13e-443d-b718-bdf91b083bb3.png)

启动好 nameServer 和 Broker 之后, 启动生产者会报这样的错误

    Exception in thread "main" org.apache.rocketmq.remoting.exception.RemotingTooMuchRequestException: sendDefaultImpl call timeout at org.apache.rocketmq.client.impl.producer.DefaultMQProducerImpl.sendDefaultImpl(DefaultMQProducerImpl.java:588) at org.apache.rocketmq.client.impl.producer.DefaultMQProducerImpl.send(DefaultMQProducerImpl.java:1223) at org.apache.rocketmq.client.impl.producer.DefaultMQProducerImpl.send(DefaultMQProducerImpl.java:1173) at org.apache.rocketmq.client.producer.DefaultMQProducer.send(DefaultMQProducer.java:214) at com.baojian.mob.base.producer.SyncProducer.main(SyncProducer.java:41)15:22:31.455 [NettyClientSelector_1] INFO RocketmqRemoting - closeChannel: close the connection to remote address[] result: true15:22:32.049 [NettyClientSelector_1] INFO RocketmqRemoting - closeChannel: close the connection to remote address[] result: true

1
Plain Text

**原因： BrokerIP 展示的是云服务器的本地 IP，不是公网 IP;**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xuvc3b/1616130360476-d1af823e-9914-4962-8925-f5cd19599ee0.png)

**解决方法：**

在 conf/broker.conf 中 加入 两行配置

namesrvAddr = 你的公网 IP:9876

brokerIP1=你的公网 IP

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/xuvc3b/1616130360465-f2ca9153-4a43-4646-b792-9f8dc7c5d986.png)

**重新启动 broker**

启动 broker 的指令要修改下, 要将这个配置文件指定加载

    nohup sh mqbroker -n localhost:9876 -c ../conf/broker.conf autoCreateTopicEnable=true &

1
Plain Text
