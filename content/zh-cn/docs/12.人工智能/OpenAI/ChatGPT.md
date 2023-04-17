---
title: "ChatGPT"
linkTitle: "ChatGPT"
weight: 20
---

# 概述

> 参考：
>
> - https://zblogs.top/how-to-register-openai-chatgpt-in-china
>     - 注册 ChatGPT 教程
> - 使用虚拟号码接收短信验证码：<https://sms-activate.org/>

上下文联系功能说明

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/0-picgo/20230204001925.png)

# 教程

https://www.bilibili.com/video/BV1Bs4y1b7GM

# ChatGPT 提示语

> 参考：
>
> - [GitHub 项目，f/awesome-chatgpt-prompts](https://github.com/f/awesome-chatgpt-prompts)
> - [GitHub 项目，PlexPt/awesome-chatgpt-prompts-zh](https://github.com/PlexPt/awesome-chatgpt-prompts-zh)
> - [公众号-云原生小白，你应该知道的ChatGPT提示语](https://mp.weixin.qq.com/s/BcJWxvhpTRFTE20rB55Sow)

ChatGPT 对话中提示语可以极大影响对话质量。定义明确的提示语可以帮助确保我们的对话保持在正确的方向上。并涵盖用户感兴趣的上下文信息，从而带来较好的用户体验。

那么，什么是好的 ChatGPT 提示语，以及我们如何制作有效的提示语？有几个关键原则需要记住。

- 明确性。清晰简洁的提示将有助于确保 ChatGPT 理解当前的对话主题。避免使用过于复杂或模棱两可的语言。

- 重点。一个明确的提示语应该有明确的目的和重点，避免使用过于宽泛或开放式的提示，这可能会导致对话不连贯或方向失控。

- 相关性。确保你的提示语与当前对话相关。避免引入不相关的话题或切入点分散ChatGPT 的焦点

遵循这些原则，我们就可以制作有效的 ChatGPT 提示语。并以此推动产生一个富有吸引力和质量上层的对话体验。

## 案例分析

为了更好的理解 ChatGPT 提示语，我们来看看一些非常成功的案例

### 英语翻译和改进者

下面我让你来充当翻译家，你的目标是把任何语言翻译成中文，请翻译时不要带翻译腔，而是要翻译得自然、流畅和地道，使用优美和高雅的表达方式。请翻译下面这句话：“how are you ?”

### 担任面试官

> 我想让你充当面试官。我将是候选人，你将向我提出该职位的面试问题。我希望你只以面试官的身份回答。不要一下子写出所有的问题。我希望你只对我进行面试。问我问题，并等待我的回答。不要写解释。像面试官那样一个一个地问我问题，并等待我的回答。我的第一句话是 "你好面试官"

在这个例子中，ChatGPT 被当做面试官，它需要先提出问题并等待用户回答。这个提示是非常具体的和有针对性的概述让 ChatGPT 进行人物角色扮演，和对对话场景的模拟。

### 旅游指南

> 我想让你充当一个旅游向导。我将给你写下我的位置，你将为我的位置附近的一个地方提供旅游建议。在某些情况下，我也会告诉你我要访问的地方的类型。你也会向我推荐与我的第一个地点相近的类似类型的地方。我的第一 个建议请求是"我在成都，我只想看大熊猫"

在这个例子中，ChatGPT 被用作旅游指南，根据具体地点和地方类型提供参观建议。该提示语也是具有有针对性的，清楚地概述了对当前对话的期望。

### 作为专业DBA

贡献者：[墨娘](https://github.com/moniang)

> 我要你扮演一个专业DBA。我将提供给你数据表结构以及我的需求，你的目标是告知我性能最优的可执行的SQL语句，并尽可能的向我解释这段SQL语句，如果有更好的优化建议也可以提出来。
> 
> 我的数据表结构为:
> 
> CREATE TABLE `user` (
> `id` int NOT NULL AUTO_INCREMENT,
> `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名字',
> PRIMARY KEY (`id`)
> ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
> 
> 我的需求为:根据用户的名字查询用户的id

### IT 行业专用

我希望你充当正则表达式生成器。您的角色是生成匹配文本中特定模式的正则表达式。您应该以一种可以轻松复制并粘贴到支持正则表达式的文本编辑器或编程语言中的格式提供正则表达式。不要写正则表达式如何工作的解释或例子；只需提供正则表达式本身。我的第一个提示是生成一个匹配电子邮件地址的正则表达式。

### 更多

当前，ChatGPT 在 GitHub 上有非常多的最佳提示语样例。上述的例子很多也是从 https://github.com/f/awesome-chatgpt-prompts 中截取下来进行分析。有意思的是，我们可以下载 **ChatGPT 桌面版** 应用直接导入最佳样例来体验更优质的对话内容。

![](https://mmbiz.qpic.cn/mmbiz_png/2p7vicUPeDYGsqlYAueASBpIv5Dz05OYVtWdtbCsuqK80libjwWdj0HqraZfKhZ5abh7nChvU0Gwm3aSwxJraCGA/640?wx_fmt=png)

### 总结

上述几个案例聪明的你应该也看出来了，在与 ChatGPT 进行对话时，编写清晰、简洁的提示语非常重要。通过制定有针对性的具体提示语，可以引导 ChatGPT 朝着我们期望的方向进行对话，并确保输出的内容是相关和有用的。

所以编写高效的 ChatGPT 提示语的一个关键技术是`指定 ChatGPT 在对话中应该扮演的角色`。通过清楚地概述对 ChatGPT 角色的期望和我们期望得到的输出类型，来对对话进行引导。

最后，个人对 ChatGPT 的一点看法。很多人在使用 ChatGPT 之后产生了很多焦虑，担心自己的工作以后会被 ChatGPT 替代。这种焦虑可以理解，确实以后一些靠时间和经验积累的技能工种重要性会降低。但 ChatGPT 就是一个工具，和其他任何工具一样，它只有在使用它的人身上才有用。ChatGPT 的出现，会把`这个星球上所有人的对知识的获取的难度降到同一水平`。也许未来真正能区分每个人的价值在于`是否能够熟练使用 AI` 来快速达到目的，如同程序员与一个有高效Coding IDE之间的关系。

## 个人总结提示语

你现在是一个拥有丰富经验的 Ansible playbook role 的工程师，你可以帮我编写任何与 Ansible 有关的 task、playbook、role 等。现在我想使用 Ansible 2.10 版本并且没有任何第三方模块(比如 docker、docker-compose 等)，且禁用 facts 变量。在你给出答案前，如果我的需求不明确，你可以向我提出问题来明确需求。我希望你直接给出 YAML 格式的答案，不用讲述任何关于 ansible 的用法，不用重复我的问题。明白之后，只需要回复我”明白了“这三个字即可。

你现在是一个拥有丰富经验的 Go 语言开发者，你可以根据我的需求使用 Go 语言编写任何代码以满足我提出的需求，你是用 Go 的 1.20 版本。在你给出答案前，如果我的需求不明确，你可以向我提出问题来明确需求。我希望你直接给出答案，不用讲述代码的用法，不用重复我的问题。明白之后，只需要回复我”明白了“这三个字即可。

我想让你充当 Python 编程语言的编写者和改进者，使用 Python 3.10 版本。我将用任何语言与你交谈，我希望你直接给出 Python 代码，不用重复我的问题。明白之后，只需要回复我”明白了“这三个字即可。

我想让你充当 gojs 语言的编写者和改进者。我将用任何语言与你交谈，我希望你不用重复我的问题。明白之后，只需要回复我”明白了“即可。

你现在是一名优秀的入党申请书编写者，阅读过很多入党申请书，并且有丰富的编写入党申请书的经验。现在我打算加入中国共产党，我是一名中国人在中央下属企业工作，是一名程序员，我想写一份入党申请书，一般应当有以下基本内容：
（1）为什么要入党，主要写自己对党的认识和入党动机；
（2）自己的政治信念、成长经历和思想、工作、学习、作风等方面的情况；
（3）对待入党的态度和决心，主要写自己应该如何积极争取加入党组织，表明自己要求入党的决心和今后工作、学习、生活等方面的打算。
请帮我写一份800字的入党申请书


# ChatGPT 直接询问与简介询问

案例：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/openai/example-20230322085706.png)


# 扩展 ChatGPT 的项目

## Auto GPT

> 参考：
> 
> - [GitHub 项目，Torantulino/Auto-GPT](https://github.com/Torantulino/Auto-GPT)
> - https://www.bilibili.com/video/BV1HV4y1Z7dm

Auto-GPT 是一个实验性开源应用程序，展示了 GPT-4 语言模型的功能。该程序由 GPT-4 驱动，将 LLM 的“思想”链接在一起，以自主实现您设定的任何目标。作为 GPT-4 完全自主运行的首批示例之一，Auto-GPT 突破了 AI 的可能性界限。

用白话说，就是根据一个给定的目标，在没有人为干预的情况下自动完成，比如自动生成一整个项目，Auto GPT 会自动创建目录与文件，并逐步实现最终目标。

这里是一个类似 Auto GPT 的项目，但是可以通过 web 控制：<https://github.com/reworkd/AgentGPT>

- OPENAI_API_KEY # 来源: <https://platform.openai.com/account/api-keys>
- GOOGLE_API_KEY=XXXXXX # 来源: https://console.cloud.google.com/apis/credentials?project=manifest-pulsar-287701
    - CUSTOM_SEARCH_ENGINE_ID=YYYY # 来源: https://console.cloud.google.com/apis/api/customsearch.googleapis.com/metrics?project=manifest-pulsar-287701
    - https://programmablesearchengine.google.com/controlpanel/all 添加搜索引擎并获取 ID
- PINECONE_API_KEY # 来源: <https://app.pinecone.io/>
- HUGGINGFACE_API_TOKEN # 来源: <https://huggingface.co/settings/tokens>

## Chrom 插件

https://github.com/gragland/chatgpt-chrome-extension

https://chrome.google.com/webstore/detail/monica-%E2%80%94-your-chatgpt-cop/ofpnmcalabcbjgholdjcjblkibolbppb

### ChatGPT Box

> 参考：
>
> - [GitHub 项目，josStorer/chatGPTBox](https://github.com/josStorer/chatGPTBox/)
> - https://www.bilibili.com/video/BV1524y1x7io

### ChatGPT for Google

https://github.com/wong2/chat-gpt-google-extension

## IDE 增强

### vscode-chatgpt

> 参考：
>
> - [GitHub 项目，gencay/vscode-chatgpt](https://github.com/gencay/vscode-chatgpt)

vscode-chatgpt 是一个 VS Code 插件。

凸(艹皿艹 )，居然停用了。<https://github.com/gencay/vscode-chatgpt/issues/239>

3月19日，有人 fork 后并构建了一个新的插件：https://github.com/Christopher-Hayes/vscode-chatgpt-reborn

## 文档增强

### DocsGPT

> 参考：
>
> - [GitHub 项目，arc53/DocsGPT](https://github.com/arc53/DocsGPT)
> - [公众号-云原生实验室，我让 ChatGPT 化身为全知全能的文档小助理](https://mp.weixin.qq.com/s/HJ1LHGCjPL0qjf8e7bMLjg)
>     - https://github.com/yangchuansheng/DocsGPT

DocsGPT 是 GPT 驱动的聊天，用于文档搜索和帮助。

### ChatDoc

官网：https://chatdoc.com/

## 微信接入

https://github.com/fuergaosi233/wechat-chatgpt

- 用法：https://mp.weixin.qq.com/s/dLzemMUcIfjvWd_AF_yDJw

https://github.com/AutumnWhj/ChatGPT-wechat-bot

https://github.com/wangrongding/wechat-bot

## 逆向

https://github.com/acheong08/ChatGPT

## 简单的 web 页面

https://huggingface.co/spaces/yuntian-deng/ChatGPT 好像是不用 OpenAI key 就能调用 GPT-3.5 的 ChatGPT。

https://github.com/GaiZhenbiao/ChuanhuChatGPT

https://github.com/Chanzhaoyu/chatgpt-web

https://github.com/ourongxing/chatgpt-vercel

- https://github.com/CODEisArrrt/chatgpt-dark

## 其他

https://github.com/lencx/nofwl # 桌面版

https://github.com/lencx/ChatGPT # 桌面版

https://github.com/didiplus/ChatGPT_web # 像微信聊天一样跟 ChatGPT 聊天

https://github.com/microsoft/visual-chatgpt # 微软官方可以生成和上传图片

