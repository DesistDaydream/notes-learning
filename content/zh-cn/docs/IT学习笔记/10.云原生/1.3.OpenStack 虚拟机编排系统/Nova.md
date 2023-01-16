---
title: Nova
---

# Nova 的子组件以及工作流程

Nova 的子组件

1. API

   1. nova-api 接收和响应客户的 API 调用。 除了提供 OpenStack 自己的 API，nova-api 还支持 Amazon EC2 API。 也就是说，如果客户以前使用 Amazon EC2，并且用 EC2 的 API 开发了些工具来管理虚机，那么如果现在要换成 OpenStack，这些工具可以无缝迁移到 OpenStack，因为 nova-api 兼容 EC2 API，无需做任何修改。

2. Compute Core

   1. nova-scheduler #虚机调度服务，负责决定在哪个计算节点上运行虚机

   2. nova-compute #管理虚机的核心服务，通过调用 Hypervisor API 实现虚机生命周期管理

   3. Hypervisor #计算节点上跑的虚拟化管理程序，虚机管理最底层的程序。 不同虚拟化技术提供自己的 Hypervisor。 常用的 Hypervisor 有 KVM，Xen， VMWare 等

   4. nova-conductor #nova-compute 经常需要更新数据库，比如更新虚机的状态。 出于安全性和伸缩性的考虑，nova-compute 并不会直接访问数据库，而是将这个任务委托给 nova-conductor，这个我们后面详细讨论。

3. Console Interface

   1. nova-console 用户可以通过多种方式访问虚机的控制台：nova-novncproxy，基于 Web 浏览器的 VNC 访问 nova-spicehtml5proxy，基于 HTML5 浏览器的 SPICE 访问 nova-xvpnvncproxy，基于 Java 客户端的 VNC 访问

   2. nova-consoleauth 负责对访问虚机控制台请求提供 Token 认证

   3. nova-cert 提供 x509 证书支持

4. Database Nova 会有一些数据需要存放到数据库中，一般使用 MySQL。数据库安装在控制节点上。 Nova 使用命名为 “nova” 的数据库。

5. Message Queue 在前面我们了解到 Nova 包含众多的子服务，这些子服务之间需要相互协调和通信。为解耦各个子服务，Nova 通过 Message Queue 作为子服务的信息中转站。 所以在架构图上我们看到了子服务之间没有直接的连线，是通过 Message Queue 联系的。OpenStack 默认是用 RabbitMQ 作为 Message Queue。 MQ 是 OpenStack 的核心基础组件。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cu87t4/1616123186194-d43ec1f0-3506-4784-8c70-81fc866de03d.jpeg)

1. 客户（可以是 OpenStack 最终用户，也可以是其他程序）向 API（nova-api）发送请求：“帮我创建一个虚机”

2. API 对请求做一些必要处理后，向 Messaging（RabbitMQ）发送了一条消息：“让 Scheduler 创建一个虚机”

3. Scheduler（nova-scheduler）从 Messaging 获取到 API 发给它的消息，然后执行调度算法，从若干计算节点中选出节点 A

4. Scheduler 向 Messaging 发送了一条消息：“在计算节点 A 上创建这个虚机”

5. 计算节点 A 的 Compute（nova-compute）从 Messaging 中获取到 Scheduler 发给它的消息，然后在本节点的 Hypervisor 上启动虚机。

6. 在虚机创建的过程中，Compute 如果需要查询或更新数据库信息，会通过 Messaging 向 Conductor（nova-conductor）发送消息，Conductor 负责数据库访问。

7. 以上是创建虚机最核心的步骤。

程序之间的调用通常分两种：同步调用和异步调用。

同步调用

API 直接调用 Scheduler 的接口是同步调用。 其特点是 API 发出请求后需要一直等待，直到 Scheduler 完成对 Compute 的调度，将结果返回给 API 后 API 才能够继续做后面的工作。

异步调用

API 通过 Messaging 间接调用 Scheduler 就是异步调用。 其特点是 API 发出请求后不需要等待，直接返回，继续做后面的工作。 Scheduler 从 Messaging 接收到请求后执行调度操作，完成后将结果也通过 Messaging 发送给 API。在 OpenStack 这类分布式系统中，通常采用异步调用的方式，其好处是：

1. 解耦各子服务。 子服务不需要知道其他服务在哪里运行，只需要发送消息给 Messaging 就能完成调用。

2. 提高性能 异步调用使得调用者无需等待结果返回。这样可以继续执行更多的工作，提高系统总的吞吐量。

3. 提高伸缩性 子服务可以根据需要进行扩展，启动更多的实例处理更多的请求，在提高可用性的同时也提高了整个系统的伸缩性。而且这种变化不会影响到其他子服务，也就是说变化对别人是透明的。

image 说明：创建 instance 时会使用一个从 glance 下载下来的 image(暂时称作 image base）作为 backing file，当 instance 创建完成后生成了一个新的 image(暂时称作 using image)，启动时候后对 image 进行的操作都是在 using image 中进行的，base image 不受影响，当使用相同 image 创建多个 instance 的时候，就不会下载新的 image，而是直接使用已经下载好的在\_base 文件夹下的 base image 来直接创建。

1. \~/data/nova/instances/ 这是保存 using image 的路径，在该文件夹下的 image ID 是文件夹形式，可以进入该文件夹查看该 using image 的日志磁盘信息等附加信息，这些信息都是可变的。

2. \~/data/nova/instances/\_base/ 这是从 glance 下载下来的 base image 的保存路径，这个基本镜像是不变的。

注：VMware 中的快照就相当于这里的 base image(也可叫做模板)，基于这个快照使用的系统就是 using image。我们也可以把 using image 创建成一个新的 new base image，然后新的 image 发送给 glance，然后再从 glance 下载 new base image 作为 base image 2 来使用。
