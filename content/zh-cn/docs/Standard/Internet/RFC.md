# 概述

> 参考：
> 
> - [Wiki，RFC](https://en.wikipedia.org/wiki/Request_for_Comments)
> - [RFC 存储库搜索](https://www.rfc-editor.org/search/rfc_search.php)
> - [RFC Editor 存储库搜索](https://www.rfc-editor.org/)

**Request for Comments(征求意见，简称 RFC)**是互联网协会（ISOC）及其相关机构的出版物，最突出的互联网工程工作组（IETF），互联网的主要技术开发和标准制定机构。

# RFC 文档阅读方法

参考：[知乎](https://zhuanlan.zhihu.com/p/44635072)

- Obsoleted by: NUM # 是当前文档的下一版。可以描述为：当前文档被 NUM 淘汰。也就是说：对于 NUM 文档来说，当前文档已过时。
- Updated by: NUM # 是当前文档的早期版本。可以描述为：当前文档被 NUM 更新。也就是说：对于 NUM 文档来说，当前文档是已更新的。
- Obsoletes: NUM # 对于当前版本，NUM 是过时的
- Updates: NUM # 对于当前版本，NUM 是最新的

## 如何阅读 RFC

来源： [How to read RFC?](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/standards/)

无论好坏，请求注释文档（RFC）包含了我们在 Internet 上遇到的许多协议。这些 RFC 文档被开发人员视为圣经，他们会试着去发现隐藏的含义，即使无法理解也无关紧要。虽然这通常会导致挫败感 - 但更重要的是 - RFC 其中的操作性和安全性思考。
然而，根据我对 HTTP 和[其他一些事情的](https://link.zhihu.com/?target=https%3A//datatracker.ietf.org/person/Mark%2520Nottingham)经验和收获，我们通过对 RFC 如何构建和发布的一些了解，可以更容易理解正在查看的 RFC 内容。

- [从哪儿开始阅读？](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23where-to-start)
- [它是什么类型的 RFC？](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23what-kind-of-rfc-is-it)
- [RFC 是最新的吗？](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23is-it-current)
- [理解 RFC 背景](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23understanding-context)
- RFC 语法要求[SHOULD](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23should)

[RFC 阅读技巧实例](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23reading-examples)

- [ABNF](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23on-abnf)
- [考虑安全因素](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23security-considerations)
- [更多](https://link.zhihu.com/?target=https%3A//www.mnot.net/blog/2018/07/31/read_rfc%23finding-out-more)

## 1.从哪儿开始阅读？

查找 RFC 的规范位置是 [RFC Editor 网站](https://www.rfc-editor.org/)。但是，正如我们将在下面看到的那样，RFC 编辑器缺少一些关键信息，因此大多数人都使用 [tools.ietf.org](https://link.zhihu.com/?target=https%3A//tools.ietf.org/)。
即使找到合适的 RFC 也很困难，因为有这么多的 RFC（目前，将近 9,000 份）。显然，您可以使用通用 Web 搜索引擎找到它们，并且 RFC 编辑器在其站点上具有出色的搜索功能。
另一个方式是[EveryRFC](https://link.zhihu.com/?target=https%3A//everyrfc.org/)，我将它放在一起，以便按标题和关键字搜索 RFC，并按标签进行探索。
毫无疑问，纯文本 RFC 难以阅读，这种情况可以通过一些方式改变; RFC 编辑器正在包含一种[新的 RFC 格式](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/rse/format-faq/)，具有更令人满意的演示和自定义选项。与此同时，如果您想要更多可用的 RFC，您可以将第三方存储库用于选定的存储库; 例如，[greenbytes](https://link.zhihu.com/?target=https%3A//greenbytes.de/tech/webdav/)保留与 WebDAV 相关的 RFC 列表，[HTTP 工作组](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/)维护与 HTTP 相关的 RFC。

## 2.它是什么类型的 RFC？

所有 RFC 在顶部都有一个横幅，看起来像这样：

    Internet Engineering Task Force (IETF)                  R. Fielding, Ed.
    Request for Comments: 7230                                         Adobe
    Obsoletes: 2145, 2616                                    J. Reschke, Ed.
    Updates: 2817, 2818                                           greenbytes
    Category: Standards Track                                      June 2014
    ISSN: 2070-1721

在左上角，这个说“互联网工程任务组（IETF）”。这表明这是 IETF 的产品; 虽然它并不广为人知，但还有其他方法可以发布不需要 IETF 共识的 RFC; 例如，[独立媒体](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/about/independent/)。
实际上，许多“媒体”在 RFC 编辑器文档上发布了 RFC 文档。但要注意，**只有 IETF 流表明整个 IETF 已经审核并已就协议规范达成了共识**。
在较旧的文档（在 RFC5705 之前）在那里 IETF 被称为“网络工作组”，因此你需要多挖掘一下，看看它们是否代表了 IETF 的共识; 查看[RFC 编辑器站点](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/)“Status of this Memo ”部分可以了解标识情况。
在确认“请求文档”编号之前。**它被称为“Internet-Draft 互联网草案”，而不是 RFC** ; 说明这还只是一个提案，任何人都 [可以写一个](https://link.zhihu.com/?target=https%3A//datatracker.ietf.org/submit/)草案。仅仅因为某些东西是互联网草案但并不意味着它将被 IETF 采用。
RFC 的**类别**是“标准跟踪”，“信息”，“实验”或“最佳实践”之一。这些之间的区别有时是模糊的，但如果它是由 IETF（互联网工程组）产生的，那么它就有了合理的审查。但请注意，即使 IETF 已达成共识，**信息**和**实验**也不是标准。
最后，RFC 文档的**作者**会在标题的右侧著名。与学术界不同，这不是一份全面的文件清单; 通常，这是在“致谢”部分的底部附近完成的。在 RFC 中，这实际上是“谁编写了文档。”通常，您会看到附加的“Ed。”，这表明它们充当编辑人员，这是因为 RFC 或者草案是预先存在的（就像 RFC 修订时一样））。

## 3.RFC 是最新的吗？

“RFC 是一系列档案文件; 他们不能改变，即使是一个字符也不能改变”（参见[RFC7158 和 RFC7159 之间](https://author-tools.ietf.org/iddiff?url1=rfc7158&url2=rfc7159)的[差异](https://author-tools.ietf.org/iddiff?url1=rfc7158&url2=rfc7159)这样做是极端的;）。
因此，知道您正在查看正确的文档非常重要。标题包含几个元数据，可以帮助明白这是关于什么内容的 RFC：

- **Obsoletes**列出了完全被取代的 RFC 文档（取代文档会比旧文档大）; 即，你应该使用这个文件，而不是那个。请注意，旧版本的协议不一定会在较新版本的协议出现时废弃; 例如，HTTP / 2 不会废弃 HTTP / 1.1，因为实现旧协议仍然是合法的（也是必要的）。但是，RFC7230 确实废弃了 RFC2616，RFC2616 是该协议的参考。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283937874-fed16b60-0dda-4fbd-9d00-f9d3d8a43c20.jpeg)

- **update**列出了本文档进行实质性更改的 RFC; 换句话说，如果你正在阅读这边 RFC 文档，同时也应该阅读这些 Update RFC 文章。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283937793-3b3e3793-6810-4627-b420-190386ac8d0c.jpeg)

遗憾的是，ASCII 格式文本 RFC（例如，在 RFC 编辑器站点上的 RFC）不会告诉你当前正在查看的文档哪些文档更新或废弃了。这就是大多数人在[http://tools.ietf.org](https://link.zhihu.com/?target=http%3A//tools.ietf.org)上使用 RFC 存储库的原因，它将这些信息放在这样的[Banner](https://link.zhihu.com/?target=https%3A//tools.ietf.org/html/rfc2616)：

    [Docs] [txt|pdf] [draft-ietf-http...] [Tracker] [Diff1] [Diff2] [Errata]
    Obsoleted by: 7230, 7231, 7232, 7233, 7234, 7235          DRAFT STANDARD
    Updated by: 2817, 5785, 6266, 6585                          Errata Exist

工具页面上的每个数字都是一个链接，因此您可以轻松找到当前文档。
即使是最新的 RFC 也经常出现问题。在工具横幅中，您还会在右侧看到“Errata Exist”上面的勘误表链接警告。
**勘误表 Errata**是对文档的更正和澄清，更正和澄清是不值得发布新的 RFC。但有时它们会对 RFC 的实现方式产生重大影响（例如，如果规范中的错误是一个严重的错误），那么它们值得通过。
例如，这是[RFC7230](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/errata_search.php%3Frfc%3D7230)的[勘误表 Errata](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/errata_search.php%3Frfc%3D7230)。当阅读勘误表时，请记住他们的状态; 许多人的修订被拒绝，因为有人误解了规范。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283937883-ea04793f-1070-497d-85db-9870a38e9d99.jpeg)

## 4.理解 RFC 背景

对于开发人员来说，查看 RFC 中的语句，实现他们看到的内容，但实际会与作者的意图相反，这种情况比想象中更为常见。
这是因为在选择性地阅读规范时，以一种不能被误解的方式编写规范是极其困难的（任何仿佛圣经一样的文档都会出现类似情况）。
因此，不仅需要阅读直接相关的文本，还必须（至少）阅读它引用的任何相关内容，无论是相同的规范还是不同的规范。如果您无法阅读整个 RFC，那么通过阅读任何可能相关的部分也会有很大帮助。
例如，HTTP 消息头被[定义](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7230.html%23http.message)为由 CRLF（回车\r 换行\n）分隔，但是如果你在[这里](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7230.html%23message.robustness)跳过，你会看到“收件人可以将单个 LF（换行\n）识别为行终止符并忽略任何前面的 CR（回车\n）。”
同样重要的是要记住，许多协议都设立了[IANA 注册](https://link.zhihu.com/?target=https%3A//www.iana.org/protocols)管理[机构](https://link.zhihu.com/?target=https%3A//www.iana.org/protocols)来管理其扩展点; 这些注册管理机构不是规范，是事实的来源。例如，HTTP 方法的规范列表在[此注册表中](https://link.zhihu.com/?target=https%3A//www.iana.org/assignments/http-methods/http-methods.xhtml)，而不是任何 HTTP 规范。

## 5.RFC 语法要求

几乎所有的 RFC 都有类似于顶部的样板：

    The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
    "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
    "OPTIONAL" in this document are to be interpreted as described in
    BCP 14 [RFC2119] [RFC8174] when, and only when, they appear in all
    capitals, as shown here.

这些[RFC2119](https://link.zhihu.com/?target=https%3A//tools.ietf.org/html/rfc2119)关键字有助于定义互操作性，但它们有时也[会使](https://link.zhihu.com/?target=https%3A//tools.ietf.org/html/rfc2119)开发人员感到困惑。看到规范说出如下内容是很常见的：

    The Foo message MUST NOT contain a Bar header.

该要求被置于协议假象“Foo message”上。如果您要发送一个，很明显它不需要包含 Bar header; 如果你包含一个，它将不是一个符合要求的消息。
但是，收件人的行为不太清楚; 如果你看到带有“Bar header“的“Foo message“，你会怎么做？
一些开发人员会拒绝包含“Bar header“的“Foo message“，即使规范没有说明这样做。其他人仍将处理消息，但剥离 Bar header，或忽略 Bar header - 即使规范明确指出需要处理所有 header。
所有这些事情都可能 -或者无意中 - 导致一些互操作性问题。正确的做法是遵循标题的正常处理，**除非有相反的特定要求**。
这是因为通常会写规范以便明确指定行为; 换句话说，允许未明确禁止的所有内容（只要没说的都是允许，但是一旦严格规定就需要严格遵循）。因此，过多地阅读规范可能会无意中造成伤害，因为你会被引入其他人必须解决的新思维。
在理想的世界中，规范将根据处理消息的人的行为来定义，如下所示：

    Senders of the Foo message MUST NOT include a Bar header. Recipients
    of a Foo message that includes a Bar header MUST ignore the Bar header,
    but MUST NOT remove it.

如果没有这个，最好在规范的其他地方寻找有关错误处理的更加一般的建议（例如，HTTP 的[一致性和错误处理](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7230.html%23conformance)部分）。
另外，请记住要求的目标 ：大多数规范都有一套高度发展的术语，用于区分协议中的不同角色。
例如，HTTP 具有[代理](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7230.html%23intermediaries)，[代理](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7230.html%23intermediaries)是一种中介，它实现客户端和服务器（但不是用户代理或源服务器）; 他们需要关注针对所有这些角色的要求。
同样，HTTP 根据具体情况区分“生成”消息并仅在某些要求中“转发”消息。注意这种特定的术语可以为你节省很多猜测。

## SHOULD

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283938058-68f4c249-bd67-431b-b432-50f3646251dc.jpeg)

是的，应该得到自己的部分。尽管努力消除它，但这个多余的术语困扰着许多 RFC。RFC2119 将其描述为：

    SHOULD  This word, or the adjective "RECOMMENDED", mean that there
            may exist valid reasons in particular circumstances to ignore a
            particular item, but the full implications must be understood and
            carefully weighed before choosing a different course.

在实践中，作者经常使用 SHOULD 和 SHOULD NOT，这表示“我们希望你这样做，但我们知道我们不能总是要求这么做。”
例如，在[HTTP 方法](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7231.html%23method.overview)的[概述中](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7231.html%23method.overview)，我们看到：

    When a request method is received that is unrecognized or not
    implemented by an origin server, the origin server SHOULD respond
    with the 501 (Not Implemented) status code. When a request method
    is received that is known by an origin server but not allowed for
    the target resource, the origin server SHOULD respond with the 405
    (Method Not Allowed) status code.

这个 SHOULD 不是必须的，因为服务器可能会合理地决定采取另一个行动; 如果请求被认为是来自攻击者的客户端，则可能丢弃连接，或者要求 HTTP 身份验证，则可能会在到达 405 之前使用 401（未经过身份验证）强制执行该请求。
SHOULD 并不意味着服务器可以自由地忽略一个要求，这看起来对 RFC 不太尊重。
有时候，我们会[看到](https://link.zhihu.com/?target=https%3A//httpwg.org/specs/rfc7231.html%23multipart.types)一个遵循这种形式的 SHOULD 部分：

    A sender that generates a message containing a payload body SHOULD
    generate a Content-Type header field in that message unless the
    intended media type of the enclosed representation is unknown to
    the sender.

注意“除非 unless” - 它指定了应该允许的“特殊情况”。可以说这可以指定为 SHOULD，因为 Unless 部分仍然适用，但这种规范风格有点普遍。

## 6.RFC 阅读实例

另一个非常常见的缺陷是浏览示例的规范，并实现它们的功能。
不幸的是，示例通常得到作者最少的关注，因为它们需要随协议的每次更改而更新。
因此，它们通常是规范中最不可靠的部分。是的，作者应该在发布前绝对仔细检查这些例子，但错误确实会出现纰漏。
此外，即使是一个完美的例子也可能无法说明是关于你正在寻找的协议的相关; 为简洁起见，它们经常被截断，或者在解码发生后显示。
尽管需要更多时间，但最好还是阅读实际文本; 示例 Examples 不是 RFC 规范。

## 7.ABNF

[增强型 BNF](https://link.zhihu.com/?target=https%3A//tools.ietf.org/html/rfc5234)通常用于定义伪协议。例如：

    FooHeader = 1#foo
    foo       = 1*9DIGIT [ ";" "bar" ]

一旦你习惯了它，ABNF 提供了一个易于理解的协议元素应该是什么样子的草图。
但是，ABNF 是“有理想和抱负的” - 它确定了一个理想的消息形式，而你生成的那些消息确实需要与之匹配。它没有指定如何处理未能匹配的已接收消息。事实上，很多规范关于在处理所有协议上很难说清楚和 ABNF 的关系。
如果你试图严格执行他们的 ABNF，大多数协议将会严重失败，但有时它很重要。在上面的例子中，分号周围不允许有空格，但你可以打赌有些人会把它放在那里，有些实现过程中会接受它。
因此，请确保阅读 ABNF 周围的文本以了解其他要求或上下文联系，并意识到如果没有直接要求，您可能必须将解析调整为比 ABNF 暗示的更容易接受的输入。
一些规范开始承认 ABNF 的期望性质并指定包含错误处理的显式解析算法。如果指定 ABNF，应严格遵循这些，以确保互操作性。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283937918-0dd17044-0b0b-4764-ac24-263f9f14447e.jpeg)

## 8.考虑安全因素

自[RFC3552](https://link.zhihu.com/?target=https%3A//tools.ietf.org/html/rfc3552)以来，RFC 样式包含了“安全注意事项”部分。
因此，如果没有关于安全性的实质性部分，很少发布 RFC; 审核流程不允许草案只是说“此协议没有考虑安全因素”。
因此，无论您是在实施还是部署协议，都必须阅读并确保您了解“安全注意事项”部分; 如果你不这样做，很可能会有一些东西会让你不知所措（没有考虑的安全因素）。
遵循其参考（如果有的话）也是一个好主意。如果没有，请尝试查找用于理解所讨论问题的一些术语。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283938062-7651f582-76f3-42a2-9c97-024799c28338.jpeg)

## 9.发现更多

如果 RFC 没有回答您的问题，或者您不确定其文本的意图，最好的办法是找到最相关的[工作组](https://link.zhihu.com/?target=https%3A//datatracker.ietf.org/wg/)并在他们的邮件列表中提出问题。如果没有涉及相关主题的活动工作组，请尝试相应[区域](https://link.zhihu.com/?target=https%3A//ietf.org/topics/areas/)的邮件列表。
提交勘误通常不是您应该采取的第一步 - 第一步一般是：与某人交谈。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zwxdg5/1618283938138-fce8f68d-637b-4752-9a0b-e6a544c76123.jpeg)

许多工作组现在正在使用 Github 来管理他们的规范; 如果您对有效规范有疑问，请继续提交问题。如果它已经是 RFC，通常最好使用邮件列表，除非您找到截然不同的指示。
我确信还有更多关于如何阅读 RFC 的文章，有些人会质疑我在这里写的内容，但这就是我对它们的看法。我希望它很有用。

## 附加

[https://www.rfc-editor.org/rfc-index-100a.html](https://link.zhihu.com/?target=https%3A//www.rfc-editor.org/rfc-index-100a.html)
[https://tools.ietf.org/rfc/index](https://link.zhihu.com/?target=https%3A//tools.ietf.org/rfc/index)
发布于 2018-09-15
