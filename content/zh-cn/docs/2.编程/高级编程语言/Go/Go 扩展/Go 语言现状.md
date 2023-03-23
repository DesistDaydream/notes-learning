---
title: Go 语言现状
---

#

<https://blog.jetbrains.com/zh-hans/go/2021/02/19/the-state-of-go/>

JetBrains 官方博客发表了一篇题为"The state of Go"的文章，他们通过深入研究有关 Go 的信息，发现了不少鲜为人知的事实，并提供了翔实的数据作为支撑。此外，JetBrains 还邀请到了知名的 Go 语言专家 Florin Pățan 针对各项数据发表了自己的见解。

**一、Go 开发者**

**数量 & 所处地区**

全球大约有 110 万名职业 Go 开发者（特指在工作中专门将 Go 作为主力编程语言的群体），如果把主要使用其他编程语言但同时兼职使用 Go 的专业开发者计算在内，这个数字可能接近 270 万。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972517-506bcf4e-f5f3-43d8-88f1-9b85ede79a96.png)

从 Go 开发者在全球地区的分布来看，生活在亚洲的职业 Go 开发者最多，大约有 57 万。

Go 语言专家 Florin 表示这在他的预期之内。他认为亚洲的 Go 开发者数量之所以高居榜首是因为那里有大量的开发者来自腾讯、阿里巴巴和华为等大型公司，这些公司一般都会雇佣许多开发者。

如果再细分下去，以国家为维度查看使用 Go 作为主力编程语言的开发者分布情况，中国所占的比例最高，全球有 16% 的 Go 开发者来自中国。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972543-961f7b1c-370b-4d43-9756-eb3b723dfd01.png)

Florin 对此同样不感到意外，不过他表示本以为俄罗斯会排第二，美国会进入前五，然而事实却是日本的 Go 开发者数量排名第二，美国排到了第七。谈及中国位居榜首的原因，Florin 认为最重要的是中国拥有数量相当庞大的开发者，根据他自己所认识的公司来看，诸如 PingCAP、腾讯、和华为都拥有大量开发者帮助构建与微服务相结合的内部工具、基础设施和后端服务。

Florin 提到了俄罗斯的 Go 社区非常活跃，Go 在那里也非常流行，不过他对 Go 开发者在日本和乌克兰的分布情况感到意外，因为他本以为德国和印度会更高，Florin 表示自己四五年前在柏林的时候，所认识的每家初创公司都使用了 Go。

**二、使用 Go 开发的软件类型**

根据 JetBrains 2020 年开发者生态调查的结果，Web 服务是使用 Go 进行开发的最受欢迎的领域，所占份额为 36％。其次分别是实用程序、IT 基础设施、工具库和系统软软件开发。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972487-8b847288-8c5a-44d3-a38f-44e1d6730fd1.png)

Florin 认为，对于 Web 服务，首要任务是创建速度足够快的 API 服务器。他们不一定需要框架，因此开发者可以使用 Go 快速启动并运行。他希望未来这张图不会发生大变化，希望看到 Web 服务获得更多的分享，因为使用 Go 入门很简单。

**三、使用 Go 的热门行业**

根据 JetBrains 2020 年开发者生态调查的结果，Go 开发者主要从事 IT 服务行业，其次是金融和金融科技，云计算/平台、大数据、移动开发和其他行业。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972505-b1a6ddeb-0312-48d6-b6b8-327c2c848bc1.png)

Florin 表示没想到 Go 在移动开发行业也能占有一席之地，毕竟它的移动开发历史比较匮乏。人们可能会使用 Go 来为移动应用程序提供 Web 服务或后端，但是仅此而已。

**四、Go 工具**

**Go Web 框架 Top5**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972515-b8149686-408b-441c-8b46-fdfe94948085.png)

## **包管理器**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972512-11a950c2-8c33-4a90-a613-3d3614390609.png)

## **Go routers**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972506-44b7a88a-aa39-4f53-bb7c-bdbadf13b393.png)

## **测试框架**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972578-6ccdd127-f895-4736-a472-2b528ffe800b.png)

**五、讨论度最高的 Go 工具和其他语言**

讨论的高频词：JSON、goroutine、PostgreSQL、MySQL、Dockers……

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972517-ce60297d-694e-440a-bc19-00e27764b78e.png)

**六、行业见解**

根据 JetBrains 2020 年开发者生态调查的结果，Go 是十大主要编程语言之一，被专业开发者采用的所占比例为 7％。Florin 认为，许多人并不倾向于以 Go 作为其第一门学习的编程语言，通常是从其他语言（例如 PHP 和 Ruby）迁移过来，据他所知主要是从 C++ 和 C# 迁移到 Go。

Florin 表示，Go 相对 PHP 的优势在于类型安全，因为 Go 是静态类型语言，而 PHP 是动态语言。这意味着编译器会帮助开发者完成大部分工作，以确保他们编写的代码能够正确编译和运行，并在运行时不会出现问题。Go 与 C++ 相比的优势是简单。在 Go 中，一切都非常简单。此外在不进行任何特殊优化的情况下，使用 Go 还会获得性能方面的提升，这对公司来说是重要的生产力优势。

Florin 还提到了 Go 采用率持续增长的另一个原因，由于许多流行的 IT 基础设施都是用 Go 编写，例如 Kubernetes、Docker 和 Vault，因此尽管许多公司的主力技术栈可能是 Java 或者其他语言，但他们也会配置针对 Go 的团队，尤其是在维护和修补此类基础设施项目方面。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/epdpuz/1616161972547-76dff786-98f7-44f4-9d56-626b26397496.png)
