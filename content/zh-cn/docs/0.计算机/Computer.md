---
title: Computer
linkTitle: Computer
date: 2024-03-15T21:31
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Computer](https://en.wikipedia.org/wiki/Computer)


# 计算机工作的原理

> 参考：
> 
> - [公众号-码农的荒岛求生，你管这破玩意叫 CPU ？](https://mp.weixin.qq.com/s/Yntk83Z5cuZ2OhfnZzvWUA)
> - [B 站-幼麟实验室，迷你计算机(小白入门)：计算机工作的原理](https://www.bilibili.com/video/BV12d4y1272j)

每次回家开灯时你有没有想过，用你按的简单开关实际上能打造出复杂的 CPU 来，只不过需要的数量会比较多，也就**几十亿**个吧。

# 伟大的发明-晶体管

过去 200 年人类最重要的发明是什么？蒸汽机？电灯？火箭？这些可能都不是，最重要的也许是这个小东西：

![](https://mmbiz.qpic.cn/mmbiz_jpg/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3cXJpq4H8DRCBwzFmdKfLXeibXdEK5HbB7ukbqzLTZibQOcIBgqVZkxYA/640)

这个小东西就叫**晶体管**，你可能会问，晶体管有什么用呢？

实际上晶体管的功能简单到不能再简单，给一端通上电，那么电流可以从另外两端通过，否则不能通过，其本质就是一个开关。

就是这个小东西的发明让三个人获得了诺贝尔物理学奖，可见其举足轻重的地位。

**无论程序员编写的程序多么复杂，软件承载的功能最终都是通过这个小东西简单的开闭完成的**，除了神奇二字，我想不出其它词来。


## AND、OR、NOT
现在有了晶体管，也就是开关，在此基础之上就可以搭积木了，你随手搭建出来这样三种组合：

- 两个开关只有同时打开电流才会通过，灯才会亮
- 两个开关中只要有一个打开电流就能通过，灯就会亮
- 当开关关闭时电流通过灯会亮，打开开关灯反而电流不能通过灯会灭

天赋异禀的你搭建的上述组合分别就是：与门，AND Gate、或门，OR gate、非门，NOT gate，用符号表示就是这样：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3ZbdZfNMWa7OibAXw4a7J9XVTtOyPSVicdgI4icAK9GUzMa0eBIhkcxh0Q/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=284&id=ln8oX&margin=%5Bobject%20Object%5D&originHeight=284&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)


**道生一、一生二、二生三、三生万物**

最神奇的是，你随手搭建的三种电路竟然有一种很 amazing 的特性，那就是：任何一个逻辑函数最终都可以通过 AND、OR 以及 NOT 表达出来，这就是所谓的**逻辑完备性**，就是这么神奇。

也就是说**给定足够的 AND、OR 以及 NOT 门，就可以实现任何一个逻辑函数，除此之外我们不需要任何其它类型的逻辑门电路，**这时我们认为 AND、OR、NOT 门就是逻辑完备的。

这一结论的得出吹响了计算机革命的号角，这个结论告诉我们计算机最终可以通过简单的 AND、OR、NOT 门构造出来，这些简单的逻辑门电路就好比基因。
![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3efwiapNDVRibXjXaLvQ3lG7HkYreVhjiay72hsyAbpQaPEfpNAeDSdPmQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=435&id=SHQ0p&margin=%5Bobject%20Object%5D&originHeight=748&originWidth=739&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=430)
老子有云：**道生一、一生二、二生三、三生万物，实乃异曲同工之妙**。

虽然，我们可以用 AND、OR、NOT 来实现所有的逻辑运算，但我们真的需要把所有的逻辑运算都用 AND、OR、NOT 门实现出来吗？显然不是，而且这也不太可行。

## 逻辑门

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/computer/logic-gate-202304281056.png)

# 计算能力是怎么来的

现在能生成万物的基础元素与或非门出现了，接下来我们着手设计 CPU 最重要的能力：计算，以加法为例。

由于 CPU 只认知 0 和 1，也就是二进制，那么二进制的加法有哪些组合呢：

- 0 + 0，结果为 0，进位为 0
- 0 + 1，结果为 1，进位为 0
- 1 + 0，结果为 1，进位为 0
- 1 + 1，结果为 0，进位为 1，二进制嘛！

注意进位一列，只有当两路输入的值都是 1 时，进位才是 1 ，看一下你设计的三种组合电路，这就是与门啊，有没有！

再看下结果一列，当两路输入的值不同时结果为 1，输入结果相同时结果为 0，这就是异或啊，有没有！我们说过与或非门是逻辑完备可以生万物，异或逻辑当然不在话下，用一个与门和一个异或门就可以实现二进制加法：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3JxIJrOIlU0Vk38nwdcvogOutwQDWYPhNVRueicC5Men9natwfwxLkTQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=677&id=qyvVK&margin=%5Bobject%20Object%5D&originHeight=677&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)

上述电路就是一个简单的加法器，就问你神奇不神奇，加法可以用与或非门实现，其它的也一样能实现，逻辑完备嘛。

除了加法，我们也可以根据需要将不同的算数运算设计出来，负责计算的电路有一个统称，这就是所谓的 arithmetic/logic unit，ALU，CPU 中专门负责运算的模块，本质上和上面的简单电路没什么区别，就是更加复杂而已。

现在，通过与或非门的组合我们获得了计算能力，计算能力就是这么来的。

但，只有计算能力是不够的，电路需要能**记得住**信息。


# 神奇的记忆能力

到目前为止，你设计的组合电路比如加法器天生是没有办法存储信息的，它们只是简单的根据输入得出输出，但输入输出总的有个地方能够保存起来，这就是需要电路能保存信息。

电路怎么能保存信息呢？你不知道该怎么设计，这个问题解决不了你寝食难安，吃饭时在思考、走路时在思考，蹲坑时在思考，直到有一天你在梦中遇一位英国物理学家，他给了你这样一个简单但极其神奇的电路：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3Ix9HYprZmoBm2XYlZcDvXLB0TEXXeS7Ct3Gv1NRgtwufMmpy1CqcjA/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=701&id=tATWF&margin=%5Bobject%20Object%5D&originHeight=701&originWidth=1000&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1000)

这是两个 NAND 门的组合，不要紧张，NAND 也是有你设计的与或非门组合而成的，所谓 NAND 门就是与非门，先与然后取非，比如给定输入 1 和 0，那么与运算后为 0，非运算后为 1，这就是与非门，这些不重要。

比较独特的是该电路的组合方式，**一个 NAND 门的输出是两一个 NAND 门的输入**，该电路的组合方式会自带一种很有趣的特性，**只要给 S 和 R 段输入 1，那么这个电路只会有两种状态**:

- 要么 a 端为 1，此时 B=0、A=1、b=0；
- 要么 a 端为 0，此时 B=1、A=0、b=1;

不会再有其他可能了，**我们把 a 端的值作为电路的输出**。

此后，你把 S 端置为 0 的话 (R 保持为 1)，那么电路的输出也就是 a 端永远为 1，这时就可以说我们把 1 存到电路中了；而如果你把 R 段置为 0 的话 (S 保持为 1)，那么此时电路的输出也就是 a 端永远为 0，此时我们可以说把 0 存到电路中了。

就问你神奇不神奇，**电路竟然具备存储信息的能力了**。

现在为保存信息你需要同时设置 S 端和 R 端，但你的输入是有一个 (存储一个 bit 位嘛)，为此你对电路进行了简单的改造：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3oAg248jVapQ9xO4Go54Bvxu4kVajVk0BIokDOIibya6sb0bY0OEF0lQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=378&id=S5OL4&margin=%5Bobject%20Object%5D&originHeight=378&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)

这样，当 D 为 0 时，整个电路保存的就是 0，否则就是 1。


# 寄存器与内存的诞生

现在你的电路能存储一个比特位了，想存储多个比特位还不简单，复制粘贴就可以了：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3p3UERMqaXH1Pu4hwnicb6HwQVlSRmkW64WqGF26wiaKHc1JCoTibecNGQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=509&id=JwKnu&margin=%5Bobject%20Object%5D&originHeight=509&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)

我们管这个组合电路就叫**寄存器**，你没有看错，我们常说的寄存器就是这个东西。

你不满足，还要继续搭建更加复杂的电路以存储更多信息，同时提供寻址功能，就这样**内存**也诞生了。

寄存器及内存都离不开上一节那个简单电路，只要通电，这个电路中就保存信息，但是断电后很显然保存的信息就丢掉了，**现在你应该明白为什么内存在断电后就不能保存数据了吧**。


# 硬件还是软件？

现在我们的电路可以计算数据、也可以存储信息，但现在还有一个问题，那就是尽管我们可以用 AND、OR、NOT 表达出所有的逻辑函数，但我们真的有必要把所有的逻辑运算都用与或非门实现出来吗？这显然是不现实的。

这就好比厨师，你没有听说只专做一道菜的厨师然后酒店要把各个菜系厨师雇全才能做出一桌菜来吧！

中国菜系博大精深，千差万别，但制作每道菜品的方式大同小异，其中包括刀工、颠勺技术等，这些是基本功，制作每道菜品都要经过这些步骤，变化的也无非就是食材、火候、调料等，这些放到菜谱中即可，这样给厨师一个菜谱他就能制作出任意的菜来，在这里厨师就好比硬件，菜谱就好比软件。

![](https://mmbiz.qpic.cn/mmbiz_jpg/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3jo7pp7gFSicg2HtaFOz4yrK7nq2YB5sW70aNw3B8lmprnSRlMibfMNqA/640?wx_fmt=jpeg#crop=0&crop=0&crop=1&crop=1&height=719&id=ny8Pi&margin=%5Bobject%20Object%5D&originHeight=719&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)

同样的道理，**我们没有必要为所有的计算逻辑实现出对应的硬件**，硬件只需要提供最基本的功能，最终所有的计算逻辑都通过这些最基本的功能用软件表达出来就好，这就是所谓的软件一词的来源，**硬件不可变，但软件可变**，**不变的是硬件但提供不同的软件就能让硬件实现全新的功能**，**无比天才的思想**，人类真的是太聪明了。

同样一台计算机硬件，安装上 word 你就能编辑文档，安装上微信你就能读到码农的荒岛求生公众号、安装上游戏你就能玩王者农药，硬件还是那套硬件，提供不同的软件就是实现不同的功能，**每次打开电脑使用各种 App 时没有在内心高呼一声牛逼你都对不起计算机这么伟大的发明创造**，这就是为什么计算机被称为**通用**计算设备的原因，这一思想是计算机科学祖师爷图灵提出的。

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3Ky3ZJWe3H6oQS38iaVIgAB5dtCXlcuL4vR3ssQPmBqAtBVYJxNWFNPQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=504&id=sdZvM&margin=%5Bobject%20Object%5D&originHeight=504&originWidth=1024&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1024)

扯远了，接下来我们看下硬件是怎么提供所谓的基本功能的。


# 硬件的基本功

让我们来思考一个问题，CPU 怎么能知道自己要去对两个数进行加法计算，以及哪两个数进行加法计算呢？

很显然，你得告诉 CPU，该怎么告诉呢？还记得上一节中给厨师的菜谱吗？没错，CPU 也需要一张菜谱告诉自己该接下来该干啥，在这里菜谱就是机器指令，指令通过我们上述实现的组合电路来执行。

接下来我们面临另一个问题，那就是这样的指令应该会很多吧，废话，还是以加法指令为例，你可以让 CPU 计算 1+1，也可以计算 1+2 等等，实际上单单加法指令就可以有无数种组合，显然 CPU 不可能去实现所有的指令。

实际上 CPU 只需要提供**加法操作**，你提供**操作数**就可以了，CPU 说：“我可以打人”，你告诉 CPU 该打谁、CPU 说：“我可以唱歌”，你告诉 CPU 唱什么，CPU 说我可以做饭，你告诉 CPU 该做什么饭，CPU 说：“我可以炒股”，你告诉 CPU 快滚一边去吧韭菜。

因此我们可以看到 CPU 只提供**机制**或者说功能 (打人、唱歌、炒菜，加法、减法、跳转)，我们提供**策略**(打谁、歌名、菜名，操作数，跳转地址)。

CPU 表达机制就通过指令集来实现的。


## 指令集

指令集告诉我们 CPU 可以执行什么指令，每种指令需要提供什么样的操作数。不同类型的 CPU 会有不同的指令集。

指令集中的指令其实都非常简单，画风大体上是这样的：

- 从内存中读一个数，地址是 abc
- 对两个数加和
- 检查一个数是不是大于 6
- 把这数存储到内存，地址是 abc
- 等等

看上去很像碎碎念有没有，这就是机器指令，我们用高级语言编写的程序，比如对一个数组进行排序，**最终都会等价转换为上面的碎碎念指令，然后 CPU 一条一条的去执行，很神奇有没有**。

接下来我们看一条可能的机器指令：

![](https://mmbiz.qpic.cn/mmbiz_png/8g3rwJPmya2J1UjGyUZ2x0Yxh82cWJu3ssjP2NAoNF1txpoYr8icPGz9M2Q3Pof0Wpq62rgNgMZ8lKYY21U6YlQ/640?wx_fmt=png#crop=0&crop=0&crop=1&crop=1&height=201&id=ZOwu9&margin=%5Bobject%20Object%5D&originHeight=201&originWidth=1080&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=&width=1080)

这条指令占据 16 比特，其中前四个比特告诉 CPU 这是加法指令，这意味着该 CPU 的指令集中可以包含 2^4 也就是 16 个机器指令，这四个比特位告诉 CPU 该做什么，剩下的 bit 告诉 CPU 该怎么做，也就是把寄存器 R6 和寄存器 R2 中的值相加然后写到寄存器 R6 中。

可以看到，机器指令是非常繁琐的，现代程序员都使用高级语言来编写程序，关于高级程序语言以及机器指令的话题请参见[《你管这破玩意叫编程语言？](http://mp.weixin.qq.com/s?__biz=MzU2NTYyOTQ4OQ==&mid=2247485439&idx=1&sn=5045e795fe3a881ec719ffd0ea41302a&chksm=fcb980a1cbce09b7cb79cac0964d082bda3f8b94701012ab5fbd911d630bd5fef6017feb6dd9&scene=21#wechat_redirect)》。


### 指挥家：让我们演奏一曲

现在我们的电路有了计算功能、存储功能，还可以通过指令告诉该电路执行什么操作，还有一个问题没有解决。

我们的电路有很多部分，用来计算的、用来存储的，以最简单的加法为例，假设我们要计算 1+1，这两个数分别来自寄存器 R1 和 R2，要知道寄存器中可以保存任意值，**我们怎么能确保加法器开始工作时 R1 和 R2 中在这一时刻保存的都是 1 而不是其它数**？

即，我们靠什么来协调或者说靠什么来同步电路各个部分让它们协同工作呢？就像一场成功的交响乐演离不开指挥家一样，我们的计算组合电路也需要这样一个指挥家。

负责指挥角色的就是时钟信号。

时钟信号就像指挥家手里拿的指挥棒，**指挥棒挥动一下整个乐队会整齐划一的有个相应动作**，同样的，时钟信号每一次电压改变，整个电路中的各个寄存器 (也就是整个电路的状态) 会更新一下，这样我们就能确保整个电路协同工作不会这里提到的问题。

现在你应该知道 CPU 的主频是什么意思了吧，主频是说一秒钟指挥棒挥动了多少次，显然主频越高 CPU 在一秒内完成的操作也就越多。


## 大功告成

现在我们有了可以完成各种计算的 ALU、可以存储信息的寄存器以及控制它们协同工作的时钟信号，这些统称 **Central Processing Unit**，简称就是 **CPU**。


# 总结

一个小小的开关竟然能构造出功能强大的 CPU ，这背后理论和制造工艺的突破是人类史上的里程碑时刻，说 CPU 是智慧的结晶简直再正确不过。

本文从一枚开关开始讲解了 CPU 构造的基本原理，希望这篇对大家理解 CPU 有所帮助。

##### _**参考资料:**_

1. [**程序员应如何理解 CPU**](http://mp.weixin.qq.com/s?__biz=MzU2NTYyOTQ4OQ==&mid=2247483850&idx=1&sn=b90a78604fa174f0e7314227a3002bdc&chksm=fcb98694cbce0f82024467c835c6e3b4984773b1a2f6c1625d573066c36b14420d996819bed7&scene=21#wechat_redirect)[**？**](http://mp.weixin.qq.com/s?__biz=MzU2NTYyOTQ4OQ==&mid=2247483736&idx=1&sn=4da1eec64e42567a0fdf4ae6d4e9344e&chksm=fcb98606cbce0f10090d950ec468b0a1e28087cd158a850bc7dc4c262fd2612a319851987220&scene=21#wechat_redirect)
2. [**CPU 空闲时在干嘛？**](http://mp.weixin.qq.com/s?__biz=MzU2NTYyOTQ4OQ==&mid=2247485379&idx=1&sn=77ccd2258f0280dfb536ad3d389cd43a&chksm=fcb9809dcbce098b6f89dc59e71cf7fb6af6e5152e40f2f84f7c33ba9ea62e5bc8390ffd0553&scene=21#wechat_redirect)
3. [**CPU 与进程、线程、操作系统**](http://mp.weixin.qq.com/s?__biz=MzU2NTYyOTQ4OQ==&mid=2247484768&idx=1&sn=049db350af9e5eea5cf3523ceb83f447&chksm=fcb9823ecbce0b28ca28d021e68c78138cde4a1b86ea7209c0c667d3d544d223d8b2aecbccec&scene=21#wechat_redirect)


