---
title: MIT  开源许可证
---

这一回终于把 MIT 协议讲明白了

以下文章来源于微月人话 ，作者卫 sir

微月人话

简单而深入

以前看到过李笑来讲的发生在他身上的故事，说他当年 ( 2001 年 ) 住在双榆树，经常去双安商场的地下超市买东西，有一次买了个什么东西觉得不好，要退，超市服务员说按规定，该类商品售出一概不能退，李笑来大怒，说你把书面规定拿出来，有的话我就不退了，如果没有，那我就一定要退，最后叫来了超市经理，经理一看这来者不善啊，也吵不过李笑来，就给退了。

讲这个故事想说明什么呢，其实我们都明白，20 多年前的中国超市，很多管理规定都是口头上的，怎么会写成白纸黑字呢。

从超市服务员的角度看，李笑来这行为就是捣乱，是胡搅蛮缠；李笑来则肯定觉得是在维护自己正当的消费者权益；最受震动的应该是超市管理层，如果是我，我会立刻要求起草一个关于退换货的规定，我可真不想再遇到下一个这样的人。

这就是契约精神，说好的，都写下来，写下来的，我都认。

现在说 MIT 协议。

MIT 协议就是大名鼎鼎的开源软件许可协议 The MIT License，又称 MIT 许可证。

有人在两年前专门做过分析 1，MIT 是 Github 上使用率最高的许可证，第二名到第五名是 Apache 2.0、GPL 2.0、GPL 3.0 和 bsd-3-clause。

注：本文中，“MIT 协议”和“MIT 许可证”等同。

什么是开源许可证？

开源许可证是这样的，我把源码放网上了，如果还不错，就有很多人问我了，说你那个代码能不能让我用用？你那个代码我能不能放在我的产品里啊？你那个代码我用了，怎么那么多 Bug 啊？你那个代码我想当作教学案例使用，请问是不是可以啊？还有，你那个代码我用了，感觉不错，而且我还改了很多地方，我也把它放网上了，而且我还改了个名，你没有意见吧？你有意见我也准备改名了，因为现在这个软件中，我写的代码，比你写的多多了！

（这都是比较有版权意识的，怕不问你就用以后惹上官司。）

我可懒得回答这么多问题，我把这些可能问到的问题，都写成一段话，放在我的代码里，意思就是说：

我允许你们 XXX，我许可你们 XXXX，你们可以 XXXX，但是，你们必须 XXXX，如果你们 XXXX 了，你们就必须 XXXX，对了，对于 XXXX 这些情况，我可不负责。

你要同意，就用，不同意就别用。如果你用了，但违反了许可证的要求，我可能会告你啊！

这就是许可证。

你可以自己写一个许可证，但是如果你很懒的话 ( 一般人都很懒 ) ，你可以用别人写的比较好的许可证。

写的比较好的开源许可证有很多种，比如 GPL、BSD、MIT、Apache 等等，MIT 只是其中的一个。

你可以挑一个合你胃口的，这些许可证模版都是免费的，毕竟也没人指望这个卖钱。

至于它们的区别，可以看看下面这张图接受一下科普。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654522-76d7d7b6-c4cd-4bbc-87b4-835c7eb24fcb.jpeg)

“新蜂商城”事件

前段时间有一个叫做“新蜂商城”（简称“新蜂”）的开源项目有点新闻，它使用了 MIT 许可证，其作者被人告知说，哎，你的开源软件被人在网上卖哎，你不管管吗？

作者一看，还真是，有 up 主在 b 站上卖，有人在闲鱼上卖，虽然卖的也不贵，有卖 300 的，有卖几十的，但一眼望过去，很显然卖的就是自己的“新蜂”嘛！

具体可以戳这里 -> ……这是什么骚操作

然后就有点热闹，有人说这是侵权，要声援作者，控诉有人如此无耻；有人说这不算侵权，MIT 是很宽松的协议，基本上你什么都可以干，只要保留作者版权即可。

作者倒是没有想去怎么处理，作者只是觉得“我抽出下班时间，断断续续写了半年的项目，不是心血也算是我的小作品吧，开源出去就是给你这么玩的？佩服。”

大概作者还是比较年轻，不懂人世间的险恶吧！

其实这应该是预料中的事。

不应该有任何心理波澜。

更何况，MIT 许可证，允许别人卖你的源码！

从契约精神讲，说好允许别人做，就让别人做吧。

这里面有个新手不太能想明白的问题，为啥有许可证允许别人卖自己的开源软件？

这要谈到开源的精神了，还没有一点概念的读者可以先看看这篇：开源的 7 大理念

早期玩开源的人，开源自己的代码，大多不是为了卖软件，大多都有着开放、共享、自由、打破垄断等等比较理想化的情结，唯一可图的大概也就是个“名”，所以几乎所有许可证都要求保留作者名字。

为了更好地传播自己的代码，开源作者通常允许别人卖自己的源码。初期会考虑这样的情况：如果有人把 Linux 源码刻录成光盘发放，是不是应该收点成本费呢，再或者服务费？

那如果有人以此牟取暴利呢，岂不是很亏？通常不会。因为既然能在网上免费下载源码，明白人就不会再去花大价钱去买。

而且，如果作者发现真有人能使用自己的源码牟取暴利，完全可以不授权让他从中牟利，改许可证就可以了。这种事也不是没有发生过，而且还不是个例。参见开源公司被云厂商“寄生”，咋整？

比如开源云原生 SQL 数据库 CockroachDB 宣布修改开源协议，从原本的 Apache-2.0 协议修改为 BSL ( Business Source License ) ，该协议要求用户唯一不能做的是在没有取得授权的情况下以商业形式用 CockroachDB 提供数据库即服务 ( DBaaS ) 。

BSL 由 MySQL 的开发者迈克尔·蒙蒂·维德纽斯 ( Michael "Monty" Widenius ) 在 2013 年设计。它有三个主要特点，一是非商业性使用没有限制，商业性使用有限制；二是许可证中可以附加使用者自己的要求；三是有一个 change date，自此时间开始，源码将会由 BSL 转变为作者指定的其他许可证，如 GPL 等开源许可证。

依我看，闲鱼上把“新蜂”卖个几十块钱，还搭上售后服务，也算正常吧。

MIT 到底说了什么 ( 学英语！)

先看原文：

Copyright ( C )

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files ( the "Software" ) , to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

英语不好的直接晕倒！

英语就算还可以的，想弄明白，也得看好一阵！

剖析其句子结构，其实就是下面这样的，注意放在括号里面的，都不是句子的主干成分：

第一段：

Permission is ( hereby ) granted ( , free of charge, ) to any person ( obtaining a copy of ( this software and associated documentation files ( the "Software” ) ) ) , to deal in the Software without restriction, ( including without limitation the rights ( to use, copy, modify, merge, publish, distribute, sublicense, ( and/or ) sell copies of the Software, and to permit persons ( to whom the Software is furnished ) to do so ) ) , subject to the following conditions:

其大致意思就是：只要满足以下条件，许可被赋予任何 ( 获得本软件拷贝及相关文档的 ) 人 ，这些人可以免费地、没有限制地处理本软件 ( 包括随意地使用、拷贝、合并、发布、分发、再许可、卖拷贝，以及再授权其他人 ( 那些被装配了本软件的 ) 做上面说的这些事 ) 。

注意，and to permit persons 这里的 to permit，是和 to use，to copy 并列的，后面那个 to do so，则是指代前面的的一系列权利，to whom the Software is furnished 是修饰 persons 的，指的是被装配、被配发本软件的人。

再有一点比较有趣的是，这一长段被许可的选项中，后面出现了个 and/or。

and/or 一般用于连接两个选项，比如 A and/or B，意思就是说 A or B or both。这段文字中，虽然 and/or 只连接最后两个选项，但其用意似乎却是连接所有。比如：A, B, C, D and/or E，想表达的是：A and/or B and/or C and/or D and/or E，意思就是 A、B、C、D、E 这些选项可以任意组合选用。这个写法貌似严谨，实则多余。因为没有它完全不影响对文本意思的理解，有了它反而还增加了疑惑。（在文学上或是法律用语上，对 and/or 用法的批评都比较多 2。）

第二段：

The above copyright notice and this permission notice shall be included ( in all ( copies or substantial portions ) of the Software ) .

上面这句翻译过来就是：以上版权声明和许可声明都必须包含 ( 在所有的本软件的完整拷贝或者实质性成分中 ) 。所谓实质性成分，可以这样理解，比如你对这个软件做了修改，但只改了 5%，那么版权声明和许可证声明都必须包含，如果你改了 95%，那就未必需要了，具体多少需要，要看是不是实质上 ( substantially ) 仍然是人家的。

第三段：

THE SOFTWARE IS PROVIDED "AS IS”, ( WITHOUT WARRANTY OF ANY KIND ( , EXPRESS OR IMPLIED, ) ( INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. ) ） IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE ( FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, ( WHETHER ( IN AN ACTION OF CONTRACT, TORT ) OR OTHERWISE ( , ARISING FROM, OUT OF OR IN CONNECTION WITH ( THE SOFTWARE OR ( ( THE USE OR OTHER DEALINGS ) IN THE SOFTWARE. ) ) ) ) )

这段话，主要是说本软件是“AS IS”的，“AS IS”的意思就是“就这样的”，售出 ( 或免费提供 ) 后一概不负任何责任，“别再找我，就这样了”。有任何问题你就接受或者自己想办法处理吧，反正我这里不再管了，就是这个意思。有时候我在“闲鱼”上卖二手东西的时候，也真的想加一个标识“AS IS”，但是考虑到很多人不知道这个词，也就算了。

上面这段全大写的英文很长，但还好，不算很难，基本上就是：本软件是 AS IS 的 ( 不提供任何保证， ( 不管是显式的还是隐式的，包括但不限于适销性保证、适用性保证、非侵权性保证 ) ) ，在任何情况下， ( 对于任何的权益追索、损失赔偿或者任何追责 ) ，作者或者版权所有人都不会负责。( 无论这些追责产生自合同、侵权，还是直接或间接来自于本软件以及与本软件使用或经营有关的情形 )

适销性 ( MERCHANTABILITY ) 就是说商品一旦售出 ( 虽然可能是免费的 ) ，如果确有问题，就可以退换货，一般来说适销性是缺省的保证 ( 即便没有明示 ) ，法律会支持消费者对不合格产品的退换货权利 3。但对于 AS IS 这类商品而言，就是说你觉得不好使也别找我。你觉得根本没法用 ( FITNESS FOR A PARTICULAR PURPOSE ) 也别找我，而且我也不保证我这软件是不是侵权了。

MIT 协议用了几乎一半的篇幅来说这个，足以见其重要性，这是西方人很强的商品经济意识和法律意识导致的。

这是一种撇清关系的态度，就是说，出了什么事都别找我，更别去法院告我。我给你们贡献源码，可不是想给自己找麻烦的，我又不挣你们的钱。

试想一下，如果在一个医疗设备的软件中，使用了“本软件”，结果由于 bug 导致医疗事故，结果追查下来，还要找我麻烦，那我怎么受得了！

把这个 license 里所有的主干挑出来，其实就是：

Permission is granted to any person to deal in the Software without restriction，subject to the following conditions:

The above copyright notice and this permission notice shall be included.

THE SOFTWARE IS PROVIDED "AS IS", IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE.

全文翻译过来，大约是这样，下面译版的版权归我哦：

MIT 开源许可协议 ( 中文版翻译：卫 sir，遵循：CC-BY 协议 )

版权 ( C ) <这里填年份> <这里填版权者姓名>

特此向任何得到本软件拷贝及相关文档 ( 以下统称“本软件” ) 的人授权：被授权人有权使用、复制、修改、合并、发布、发行、再许可、售卖本软件拷贝、并有权向被供应人授予同等的权利，但必须满足以下条件：

在本软件的所有副本或实质性使用中，都必须包含以上版权声明和本授权声明。

本软件是“按原样“提供的，不附带任何明示或暗示的保证，包括没有任何有关适销性、适用性、非侵权性保证以及其他保证。在任何情况下，作者或版权持有人，对任何权益追索、损害赔偿以及其他追责，都不负任何责任。无论这些追责产生自合同、侵权，还是直接或间接来自于本软件以及与本软件使用或经营有关的情形。

翻译成大白话缩略版，大约是下面这样的：

版权是我 XXX 的，源自 XXXX 这一年！

我授权任何人，可以干任何事，包括卖拷贝！

但是，你必须：

保留我的版权和许可！

这软件就这样的！爱用不用，出了事别找我！

我怎么用人家代码才算合规？

使用一个用了 MIT 协议的源代码，你只需要保留人家的版权和许可证信息即可。

也就是说要包含 LICENSE 文件，这个文件有完整的 MIT 许可证，其中也会有作者的版权信息。

人家源码里面的版权和许可信息头，你也得保留。

由于 MIT 协议并不要求使用者公开源码，如果你发行的仅仅是可执行软件（不带源码），那就要在软件的某个界面上说明。

比如 Google Chrome 浏览器使用很多开源软件，在其界面中都予以了明示。

在 Chrome 浏览器的“关于”中，有一句话：“Google Chrome 的诞生离不开 Chromium 开源项目以及其他开源软件。”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654514-af86080c-b38d-410b-87da-382264e29ef2.jpeg)

点击“开源软件”，会打开一个页面，列出了一长串的开源软件、其 LICENSE 和主页（或代码托管地）。

下面是部分截图：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654516-35908e7c-5b5c-4806-8bc9-0563ea70d0e9.jpeg)

这就很规范了。

我在里面随便找个一个使用 MIT 协议的软件：brotli

进入其代码仓库后，可以看到其源码头部是这么写的。

/\* Copyright 2013 Google Inc. All Rights Reserved.

Distributed under MIT license.

See file LICENSE for detail or copy at <https://opensource.org/licenses/MIT>

\*/

截图如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654546-2a9951cc-1854-40fc-ab21-8a040710fcfa.jpeg)

注意这个用法，第一行是写明了版权，下面则是对 MIT 许可证的一个引用。这样显得比较简洁。

毕竟 MIT 许可证中唯一需要填写的就是版权。所以把第一句的信息放这里就可以了，MIT 许可证就不用动了。

注意版权是指这个软件的著作权是谁的，许可证是指版权所有者允许别人怎么用这个软件。

版权后面那句“All Rights Reserved”写不写都可以，这只是一个形式。真打到法院去，所有的权利都遵循著作权相关法律。

顺便解释一下版权中的年份概念：

软件在发布时都会声明版权，其中会包含年份，比如 Copyright 2012， Copyright 2008-2012 等。如果只包含一个年份，说明这是首次发布的年份。如果包含时间段，则第一个年份为第一次发布的年份，第二个是当前版本发布的年份。

比如微软的 Windows XP 版权声明：Copyright© 1985-2001 Microsoft Corp。说明 Windows 第一个版本发布在 1985 年，Windows XP 版本发布年份为 2001 年。

我是作者，我如何使用 MIT 协议

如果我的开源软件要使用 MIT 协议，我应该怎么做？

通常应该这么做：

- 在一级目录下，给出一个 LICENSE（或 COPYING）文件，里面就是这个许可证的全文。

- 在所有的源码头部，说明版权，说明许可。

注：有的项目使用了多个许可证，会建一个 licenses 目录放这些许可证。（比如 cockroachdb）

版权那一行，你还可以写上你软件的主页或者代码存放地，一般来说，使用者都不应该删除这行内容的。

举个例子，“木兰许可证”是这样指导人们使用的：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654557-38671027-cd5b-405c-8ac4-7daac1cc0b96.jpeg)

那么，“新蜂”是怎么做的？

我专门去“新蜂”的 github 仓库看了一眼，没错，作者在一级目录下放了 LICENSE 文件。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654543-b5460bf5-e188-422f-ac69-9e867e9216bb.jpeg)

打开 LICENSE 文件，可以看到：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654557-3d9aa8cf-b207-4ac5-a2fd-fd8a04c0fd9c.jpeg)

年份写成了 2019-2029，版权所有人写成了 newbee-mall。

这是需要改进的，年份不应该写到 2029 去，应该写到当前版本发布的年份。

版权所有人应该写作者“十三”（或“13”）而不是软件的名字。

看看源码里面是怎么写的：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tamq2d/1616161654547-b5472820-8f0a-431b-ac22-f5dec273610f.jpeg)

翻了一遍，没有发现在任何源码的头部有版权和许可信息，这个做法也不规范。

可见，作者十三对开源许可证并没有研究过。

一些问题解答

我能改许可证吗？ 当然可以。（MIT 允许你干任何事）

所以说，MIT 协议和其他协议的兼容性很强（其实是单向兼容），比如，完全可以把一个 MIT 协议的软件改为 GPL 的，但一个 GPL 的可改不回 MIT 的。

我能自己写一个许可证吗？ 当然可以。

我能不写任何许可吗？ 当然可以，不写许可，就是“保留所有权利”。你可以让他们打电话问你，写邮件问你，发微信问你，你再告诉他们可以干什么。

保留许可的意思基本上就是说，他们除了可以看你的源码 ( 因为你自己把它开源了 ) ，除了著作权法里面说的权利，基本上啥也不能干。

如果我保留所有权利，那他有运行的权利吗？ 如果他是为了学习、研究或者欣赏，是可以运行的。

他还可以评论您的软件。

因为我国著作权法赋予了他这样的权利：

现行的中华人民共和国著作权法 ( 2010 修正 ) 第二十二条中规定：

在下列情况下使用作品，可以不经著作权人许可，不向其支付报酬，但应当指明作者姓名、作品名称，并且不得侵犯著作权人依照本法享有的其他权利：

( 一 ) 为个人学习、研究或者欣赏，使用他人已经发表的作品；

( 二 ) 为介绍、评论某一作品或者说明某一问题，在作品中适当引用他人已经发表的作品；

……

如果有人没有按照我的许可做，怎么办？ 你可以告他。

不过，能不能打赢官司又是另一个话题了！

---

1. <https://www.kaggle.com/mrisdal/safely-analyzing-github-projects-popular-licenses> ↩

2. <https://en.wikipedia.org/wiki/And/or> ↩

3. <https://consumer.findlaw.com/consumer-transactions/what-is-the-warranty-of-merchantability.html> ↩

推荐阅读

树莓派销量突然猛增

80%的代码曾由一人提交，这项目何以从 ASF 毕业

红帽借“订阅”模式成开源一哥，首创者升任总裁

Git 15 周年：当年的分道扬镳，成就了今天的开源传奇

Windows 中现在有独立的 Linux 文件夹系统
