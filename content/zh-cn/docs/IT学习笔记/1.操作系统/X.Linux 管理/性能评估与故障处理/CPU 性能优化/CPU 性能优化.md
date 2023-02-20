---
title: CPU 性能优化
---

# 系统性能调优之绑定 cpu

原文链接：https://mp.weixin.qq.com/s/jiQz01hg8DiA1zucqjMZkQ

支持超线程的numa架构
------------

物理硬件视角，

*   将多个CPU封装在一起，这个封装被称为插槽Socket；
    
*   Core是socket上独立的硬件单元；
    
*   通过intel的超线程HT技术进一步提升CPU的处理能力，OS看到的逻辑上的核Processor的数量。
    

每个硬件线程都可以按逻辑cpu寻址，因此这个处理器看上去有八块cpu。![](https://mmbiz.qpic.cn/mmbiz_png/ibD9iaaPDn99gtSUiaaFg901xOL7aeib35k3l1cnIaQCeOuYLaOO6jkpcjlibwMn7H99ia4pY6ic7veFZwdkzuUblC9SQ/640?wx_fmt=png)

对于操作系统的视角：

*   CPU(s)：8
    
*   NUMA node0 CPU(s)：0，4
    
*   NUMA node1  CPU(s)：1，5
    
*   NUMA node2 CPU(s)：2，6
    
*   NUMA node3 CPU(s)：3，7
    

![](https://mmbiz.qpic.cn/mmbiz_png/ibD9iaaPDn99gtSUiaaFg901xOL7aeib35k3Go5f7Q3n45NiaGYcUFCFaPicPnGEU4tiboXy3ljpYS7GrlwIAzrpOu55A/640?wx_fmt=png)

操作系统视角.png

L1缓分成两种，一种是指令缓存，一种是数据缓存。L2缓存和L3缓存不分指令和数据。L1和L2缓存在第一个CPU核中，L3则是所有CPU核心共享的内存。L1、L2、L3的越离CPU近就越小，速度也越快，越离CPU远，速度也越慢。再往后面就是内存，内存的后面就是硬盘。我们来看一些他们的速度：

*   L1 的存取速度：4 个CPU时钟周期
    
*   L2 的存取速度：11 个CPU时钟周期
    
*   L3 的存取速度：39 个CPU时钟周期
    
*   RAM内存的存取速度 ：107 个CPU时钟周期
    

如果 CPU 所要操作的数据在缓存中，则直接读取，这称为缓存命中。命中缓存会带来很大的性能提升，因此，我们的代码优化目标是提升 CPU 缓存的命中率。![](https://mmbiz.qpic.cn/mmbiz_png/ibD9iaaPDn99gtSUiaaFg901xOL7aeib35k3MDcGSMaDov6SEc811uibXAQB5FhfarrlSNgWrNzUicUzWfTjCbibaq2dw/640?wx_fmt=png)

在主流的服务器上，一个 CPU 处理器会有 10 到 20 多个物理核。同时，为了提升服务器的处理能力，服务器上通常还会有多个 CPU 处理器（也称为多 CPU Socket），每个处理器有自己的物理核（包括 L1、L2 缓存），L3 缓存，以及连接的内存，同时，不同处理器间通过总线连接。通过lscpu来看：

`root@ubuntu:~# lscpu  
Architecture:          x86_64  
CPU(s):                32  
Thread(s) per core:    1  
Core(s) per socket:    8  
Socket(s):             4  
L1d cache:             32K  
L1i cache:             32K  
L2 cache:              256K  
L3 cache:              20480K  
NUMA node0 CPU(s):     0-7  
NUMA node1 CPU(s):     8-15  
NUMA node2 CPU(s):     16-23  
NUMA node3 CPU(s):     24-31  
`

你可能注意到，三级缓存要比一、二级缓存大许多倍，这是因为当下的 CPU 都是多核心的，每个核心都有自己的一、二级缓存，但三级缓存却是一颗 CPU 上所有核心共享的。![](https://mmbiz.qpic.cn/mmbiz_png/ibD9iaaPDn99gtSUiaaFg901xOL7aeib35k3AZ3lZlxEz64mnaUTticqnOgVKRoC7NcmZyiaACsneOu3LREqIFjlJZgQ/640?wx_fmt=png)
但是，有个地方需要你注意一下：如果应用程序先在一个 Socket 上运行，并且把数据保存到了内存，然后被调度到另一个 Socket 上运行，此时，应用程序再进行内存访问时，就需要访问之前 Socket 上连接的内存，这种访问属于远端内存访问。和访问 Socket 直接连接的内存相比，远端内存访问会增加应用程序的延迟。

### 常用性能监测工具

Linux系统下，CPU与内存子系统性能调优的常用性能监测工具有top、perf、numactl这3个工具。1） top工具 top工具是最常用的Linux性能监测工具之一。通过top工具可以监视进程和系统整体性能。

*   top                                         查看系统整体的资源使用情况
    
*   top后输入1                            查看看每一个逻辑核cpu的资源使用情况
    
*   top -p $PID -H                      查看某个进程内所有检查的CPU资源使用情况
    
*   top后输入F，并选择P选项    查看线程执行过程中是否调度到其他cpu上执行，上下文切换过多时，需要注意。
    

2） perf工具 perf工具是非常强大的Linux性能分析工具，可以通过该工具获得进程内的调用情况、资源消耗情况并查找分析热点函数。以CentOS为例，使用如下命令安装perf工具：

*   perf top                                        查看占用 CPU 时钟最多的函数或者指令，因此可以用来查找热点函数。
    
*   perf -g record -- sleep 1 -p $PID  记录进程在1s内的系统调用。
    
*   perf -g latency --sort max             查看上一步记录的结果，以调度延迟排序。
    
*   perf report                                   查看记录
    

3） numactl工具 numactl工具可用于查看当前服务器的NUMA节点配置、状态，可通过该工具将进程绑定到指定CPU核上，由指定CPU核来运行对应进程。以CentOS为例，使用如下命令安装numactl工具：

*   numactl -H                      查看当前服务器的NUMA配置。
    
*   numastat                          查看当前的NUMA运行状态。
    

### 优化方法

1） NUMA优化，减少跨NUMA访问内存 不同NUMA内的CPU核访问同一个位置的内存，性能不同。内存访问延时从高到低为：跨CPU>跨NUMA，不跨CPU>NUMA内。因此在应用程序运行时要尽可能地避免跨NUMA访问内存，这可以通过设置线程的CPU亲和性来实现。常用的修改方式有如下：（1）将设备中断绑定到特定CPU核上。可以通过如下命令绑定：

`echo $cpuNumber > /proc/irq/$irq/smp_affinity_list  
 例子：echo 0-4 > /proc/irq/78/smp_affinity_list  
      echo 3,8 > /proc/irq/78/smp_affinity_list  
`

（2）通过numactl启动程序，如下面的启动命令表示启动程序./mongod，mongo就只能在CPU core 0到core7运行(-C控制)。

`numactl -C 0-7 ./mongod  
`

（3）可以使用 taskset 命令把一个程序绑定在一个核上运行。

`taskset -c 0 ./redis-server  
`

（4）在C/C++代码中通过sched\_setaffinity函数来设置线程亲和性。（5）很多开源软件已经支持在自带的配置文件中修改线程的亲和性，例如Nginx可以修改nginx.conf文件中worker\_cpu_affinity参数来设置Nginx线程亲和性。

2绑核注意事项
-------

在 CPU 的 NUMA 架构下，对 CPU 核的编号规则，并不是先把一个 CPU Socket 中的所有逻辑核编完，再对下一个 CPU Socket 中的逻辑核编码，而是先给每个 CPU Socket 中每个物理核的第一个逻辑核依次编号，再给每个 CPU Socket 中的物理核的第二个逻辑核依次编号。![](https://mmbiz.qpic.cn/mmbiz_png/ibD9iaaPDn99gtSUiaaFg901xOL7aeib35k3Go5f7Q3n45NiaGYcUFCFaPicPnGEU4tiboXy3ljpYS7GrlwIAzrpOu55A/640?wx_fmt=png)
注意的是在多个进程要进行亲和性绑核的，你一定要注意 NUMA 架构下 CPU 核的编号方法，这样才不会绑错核。

### 预告

下一节，我们将聊聊如何通过提L1与L2缓存命中率来提高应用程序性能。