---
title: "ChatBot"
linkTitle: "ChatBot"
created: "2026-03-13T09:35"
weight: 100
---

# 概述

> 参考：
>
> - [Wiki, Chatbot](https://en.wikipedia.org/wiki/Chatbot)

ChatBot(聊天机器人) 是一种软件应用程序。可以在 [即时通信](/docs/Utils/即时通信/即时通信.md) 中让用户通过 ChatBot 实现一定的 自动化控制、情感陪伴、etc.

https://github.com/zhayujie/chatgpt-on-wechat # <font color="#ff0000">Python 实现的</font>。基于大模型搭建的聊天机器人，同时支持 微信公众号、企业微信应用、飞书、钉钉 等接入（这些接入目标抽象为 Channel），可选择GPT3.5/GPT-4o/GPT-o1/ DeepSeek/Claude/文心一言/讯飞星火/通义千问/ Gemini/GLM-4/Claude/Kimi/LinkAI（这些是 AI 模型）。能处理文本、语音和图片，访问操作系统和互联网，支持基于自有知识库进行定制企业智能客服。

- https://github.com/hanfangyuan4396/dify-on-wechat # <font color="#ff0000">Python 实现的</font>。基于 chatgpt-on-wechat，相当于其下游分支。为 channel 和 model 添加了对接目标。channel 的 wechat bot 可以对接 **Gewechat**；model 可以对接 **Dify**。
- 加入了 Ollama 支持，个人不更新项目
  - https://github.com/kaina404/chatgpt-on-wechat/tree/feature/ollama_support
  - https://github.com/Joycc/chatgpt-on-wechat/tree/master

微信的 Chatbot 可以使用多种方式实现：

- 模拟 ipad
- 模拟 web # 截至 2025-03-22 被查的厉害，而且容易失败，失败几次就会被官方警告
- 模拟 Windows 桌面端

[GitHub 项目，yincongcyincong/MuseBot](https://github.com/yincongcyincong/MuseBot) # 一个基于 **Golang** 构建的 **智能机器人**，集成了 **LLM API**，实现 AI 驱动的自然对话与智能回复。 它支持 **OpenAI**、**DeepSeek**、**Gemini**、**Doubao**、**Qwen** 等多种大模型，  并可无缝接入 **Telegram**、**Discord**、**Slack**、**Lark（飞书）**、**钉钉**、**企业微信**、**QQ**、**微信** 等聊天平台，为用户带来更加流畅、多平台联通的 AI 对话体验。

## Wechaty

> 参考：
>
> - [GitHub 项目，wechaty/wechaty](https://github.com/wechaty/wechaty)

[ECMAScript](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript.md) 实现的，一个简化构建聊天机器人的对话式 [RPA](/docs/12.AI/Automation/RPA.md) SDK，只需 6 行 JavaScript、Python、Go、Java 代码，即可创建微信机器人，并具有跨平台支持，包括 Linux，Windows，MacOS 和 Docker。

## Gewechat

https://github.com/Devo919/Gewechat

[Java](/docs/2.编程/高级编程语言/Java/Java.md) 实现的 微信机器人框架，个人微信二次开发，最简单易用的免费二开框架，微信 ipad 登录（非HOOK破解桌面端）

> https://github.com/hanfangyuan4396/dify-on-wechat/blob/master/docs/gewechat/README.md 这里重构了镜像 registry.cn-chengdu.aliyuncs.com/tu1h/wechotd:alpine ，可以不依赖 cgroup 和 --privilege
