---
title: MCP
linkTitle: MCP
weight: 51
---

# 概述

> 参考：
>
> - [GitHub 组织，modelcontextprotocol](https://github.com/modelcontextprotocol)
> - [官网](https://modelcontextprotocol.io/introduction)

**Model Context Protocol(模型上下文协议，简称 MCP)** 是一种开放协议，能够实现 LLM 应用程序与外部数据源和工具之间的无缝集成。无论是在构建 AI 驱动的 IDE、增强聊天界面，还是创建自定义 AI 工作流，MCP 都提供了一种标准化的方式，将 LLM 连接到所需的上下文。

# 规范

> 参考：
>
> - [规范](https://modelcontextprotocol.io/specification/2025-03-26)

## 架构

https://modelcontextprotocol.io/specification/2025-03-26/architecture

![800](Excalidraw/AI/mcp-arch.excalidraw.md)

MCP 的工程化实现本质也是一个类似 C/S 的架构，但是在其中多了 Model(模型) 的参与。一套完整的 MCP 系统通常包含如下几类组件：

> Tips: 有些模型进行了 MCP 的训练，可以正常识别 MCP 的上下文内容，并且也可以返回 MCP 标准的上下文

- **Host(主持人)** # 充当内容的容器和协调器。创建和管理多个 MCP Client 实例。协调 AI/LLM 集成和采样。管理跨客户端的上下文聚合。
    - 类似 [Web](/docs/Web/Web.md) 中的 User Agent 概念。
- **MCP Client** # 由 Host 创建并维护。与 MCP Server 建立 1:1 的连接关系。
- **MCP Server** # 通过 MCP 原语公开 Resources, Tools, Prompts，提供后端程序的功能，用于与 MCP Client 通信。
    - 类似各种可以提供 API 的 Web 应用程序。只不过使用了 MCP 协议包装了 API 或其他各种功能。
    - MCP Server 可以是多种形态，并不只局限于是监听了 TCP 的进程。甚至可以是一个本地文件系统上独立的 py 文件或二进制文件，只在使用时，由 MCP Client 直接使用以从标准输入和输出中进行交互。
- **Model** # AI 模型。为问题提供 分析、总结 能力。

一套完整的 MCP 系统的运行过程通常如下所述：

用户向 Host 询问问题，Host 拉起多个 MCP Client，每个 MCP Client 注册一个 MCP Server，并获取其可以提供的能力

1. Host 带着 所有 MCP Server 的信息及其能力、问题内容，交给 Model
2. Model 分析问题后，决定是否要调用 MCP Server 以及如何调用，返回 MCP 上下文内容
3. Host 带着 MCP 上下文通过 MCP Client 向 MCP Server 发起请求获取数据
4. MCP 返回包含响应数据的 MCP 上下文
5. Host 拿着包含响应数据的 MCP 上下文交给 Model 分析各种结果
6. Model 以人类可读的方式返回结果
7. 若一个问题过于复杂，需要拆解问题后，重复上述 1 - 6 步

这里 https://github.com/CassInfra/KubeDoor/blob/1.3.0/src/kubedoor-mcp/kubedoor-mcp.py 有一个非常简单直观的 MCP Server 示例，其中利用 MCP 库的 FastMCP 实例化了一个 mcp。后面通过 @mcp.tool() 多次定义 Tools，每个 Tool 都由一个 API 提供能能力。Tools 中的 Args、Returns 本质就是结构化的 Prompts。当 Client 与 Server 建立连接时，Server 会将所有 Tools 响应给 Client

> [!Tip]
> 规范定义的架构看似复杂，实际上，由于加入了 AI，比传统的规范有了更大的灵活度。本质上， AI 需要的就是结构化的 Prompts，MCP Server 的各种能力定义本质也是 Prompts，只要这些定义清晰，那么 AI 也是可以处理的，规范的内容更多的是为了让 Host、Client、Server 之间的交互可以更加开放，不会受到某个公司限制，只要满足规范约定，都可以交互。
> 
> 比如下面的简单示例，并没有完全遵守规范的所有交互数据的格式，仅仅让 Tools 的定义清晰明了，模型同样可以理解并返回正确的内容。这个示例是将 MCP Server 与 MCP Client 的行为都放在同一个行为里

```python
import json

from apis.api import LlmAPI
from datetime import datetime

# MCP Server 声明了两个工具
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "Get weather of an location, the user shoud supply a location first",
            "parameters": {
                "type": "object",
                "properties": {
                    "location": {
                        "type": "string",
                        "description": "The city and state, e.g. San Francisco, CA",
                    },
                    "date": {
                        "type": "string",
                        "description": "date, e.g. 2025-04-21",
                    }
                },
                "required": ["location", "date"]
            },
        }
    },
    {
        "type": "function",
        "function": {
            "name": "get_date",
            "description": "Get the current time"
        }
    },
]

# 这里是两个工具要执行的具体逻辑
def get_weather(location, date):
    return "{}摄氏度".format(24.4)
def get_date():
    return datetime.strftime(datetime.now(), "%Y-%m-%d")

# 定义 MCP Client 如何调用模型
web_search = {
    "enable": True,
    "enable_citation": True,
    "enable_trace": True
}
llm_api = LlmAPI(api_key="",
                 model_name="deepseek-chat",
                 is_print=True,
                 stream=True,
                 base_url="https://api.deepseek.com",
                 web_search=web_search
                 )

# 这里模拟了 MCP Client 与 MCP Server 初始化时，获取到的 Tools 列表。并没有完全按照规范定义
# Notes: 这里面省略了规范中定义的 Server 与 Client 之间交互过程
tools_map = {"get_weather": get_weather, "get_date": get_date}

# 这部分模拟了 MCP Client 与 AI Model 交互的逻辑
res = llm_api("获取今天北京的天气情况", tools)

print("================== 检查 AI 返回的将要调用的工具的上下文 ==================")
print(needCallToolsContext)

# 这部分模拟了 MCP Client 的逻辑，通过遍历直接调用 MCP Server 的 Tools
# 这是一个循环调用，每次调用都带着上一次 MCP Server 的响应信息让 AI Model 处理
for i in range(2):
    # 从 AI Model 返回的工具调用上下文中，解析出来需要向 MCP Server 发送的的数据
    # ！！！注意：这里可以看出来，Deepseek 的训练结果包含了 MCP 相关内容，在没有任何额外提示的前提下，也可以响应标准 MCP 格式的信息！！！
    sid1 = res.choices[0].message.tool_calls[0].id
    function_name = res.choices[0].message.tool_calls[0].function.name
    params = json.loads(res.choices[0].message.tool_calls[0].function.arguments)

    # ********模拟 MCP Client 向 MCP Server 发起请求获取响应信息的逻辑********
    tool_response = tools_map[function_name](**params)

    # 模拟 MCP Client 使用 MCP Server 的响应信息让 AI Model 总结响应信息
    res = llm_api(tool_response, tools, sid1)

    print(res)
```

# MCP Server

> 参考：
>
> - [规范 - 服务端功能](https://modelcontextprotocol.io/specification/2025-03-26/server)

MCP Server 的能力被抽象为如下几类（当 Client 注册 Server 时，需要将自身所具有的能力提供给 Client）：

- **Prompts(提示)** # 指导语言模型交互的预定义模板或指令。
    - TODO: 好像是让 Client 可以填充的提示词模板？
- **Resources(资源)** # 提供额外上下文的结构化数据或内容。e.g. 文件内容、git 历史、etc. ，甚至一个数据库之类的东西可以当作一个 Resource
- **Tools(工具)** # 执行操作或检索信息的可执行函数。
    - 每个 Tool 都是一个可以实现的具体功能。可以简单理解为一个 API、一个具体的功能、etc. 。在设计 MCP Server 时，Tools 通常是占比最大的部分。

https://github.com/mark3labs/mcp-go

https://github.com/ThinkInAIXYZ/go-mcp

- https://mp.weixin.qq.com/s/LFIUVdTznkr7tWZ4_TnXGA

## Tools

https://modelcontextprotocol.io/specification/2025-03-26/server/tools

数据类型

https://modelcontextprotocol.io/specification/2025-03-26/server/tools#data-types

# MCP Client

> 参考：
>
> - [规范 - 客户端功能](https://modelcontextprotocol.io/specification/2025-03-26/client/roots)

> 通常来说，各种 AI 工具都内置了 MCP Client，以便对接各自自定义的 MCP Server

# 历史
