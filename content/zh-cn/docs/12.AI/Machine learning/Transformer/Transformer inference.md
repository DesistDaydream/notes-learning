---
title: Transformer inference
linkTitle: Transformer inference
created: 2026-04-08T14:50
weight: 12
---

# 概述

> 参考：
>
> -

# NLP 推理的关联文件与配置

> [!Quote] 参考 [Transformer](/docs/12.AI/Machine%20learning/Transformer/Transformer.md#模型的关联文件与配置) 中的 模型的关联文件与配置，不同架构的模型配置文件有细微差别

# NLP 推理流程

> 参考：
>
> - [Transformers 文档，与模型聊天 - 聊天基础知识](https://huggingface.co/docs/transformers/en/conversations)

![500](/Excalidraw/AI/transformer-inference-flow.excalidraw.md)

> [!Note] 我们在 [Hugging Face](/docs/12.AI/Hugging%20Face.md) 的模型仓库中，通常可以看到除了代码之外的其他所有文件。至于代码，其实就是 transformers 依赖库。
>
> 还有人用 Model Packeage, Model Artifact, Model Repository, etc. 称呼这一整套完整的内容。

> [!Note] 下面的演示基于 Qwen3-0.6B 模型

聊天模型接受消息列表（聊天记录）作为输入。每条消息都是一个字典，包含 role 和 content 键。要发起聊天，只需添加一条 role 为 user 消息即可。

> 还可以选择添加一条 system 消息，为模型提供行为指令。e.g. `{"role": "system", "content": "你是一位乐于助人的助手。"},`

```python
chat = [
    {"role": "user", "content": "你好！"}
]
```

经过 Transfromer 处理（e.g. 分词、etc.）后，传给 Model。

Model 模型推理完成后，我们会看到模型的回复。像这样：

```json
[
    {
        "generated_text": [
            {"role": "user", "content": "你好！"},
            {"role": "assistant", "content": "你好！有什么可以帮助你的吗？"}
        ]
    }
]
```

如果想继续对话，需要将模型的回复更新到 chat 中。可以通过两种方式更新：

- 一是将文本**追加**到 `chat` 中（使用 `assistant` 角色）
- 二是将 chat 的值直接替换成 `response[0]["generated_text"]`， 其中包含完整的聊天记录，包括传入的消息以及回复的消息。

之后，在 chat 中继续追加新的 `user` 角色及其 content，继续对话。

通过重复此过程，我们可以随心所欲地持续聊天，**直到**<font color="#ff0000">模型超出上下文窗口或内存不足</font>为止。

## 加载分词器与模型

加载分词器

```python
tokenizer: PreTrainedTokenizerBase = AutoTokenizer.from_pretrained(model_name)
```

加载模型

```python
model = AutoModelForCausalLM.from_pretrained(
    model_name,
    torch_dtype="auto",
    device_map="auto",
)
```

## 分词器处理输入

准备要输入给模型的文本

```python
prompt = "hi, i'm DesistDaydream"
messages = [
    {"role": "user", "content": prompt},
]
```

### 应用聊天模板

聊天语言模型自带 **chat templates(聊天模板)**，这些模板定义了模型期望的聊天格式。我们可以使用 Transformer 的 `apply_chat_template()` 方法访问这些模板。利用分词器的 tokenizer_config.json 配置，将用户输入渲染成模型可以理解的格式，添加 **ControlTokens(控制 Token)**，e.g. `<|user|>` 、 `<|assistant|>`、`<|end_of_message|>` 、etc. ，这些 ControlTokens 使模型能够识别聊天结构。聊天格式多种多样，即使模型是从同一个基础模型微调而来，它们也可能使用不同的 ControlTokens 或 格式！

```python
text = tokenizer.apply_chat_template(
    messages,                    # 要使用模板处理的消息列表
    tokenize=False,              # 是否将文本转换为 Token 序列。默认为 Ture。关闭后只输出渲染后的文本。
    # return_tensors="pt",       # 当 `tokenize` 参数为 `True` 时，可以传入该参数，将返回特定框架的张量。否则只返回 TokenIDs
                                 # 'pt': 返回 PyTorch 的 `torch.Tensor` 对象。'np'：返回 NumPy 的 `np.ndarray` 对象。
    add_generation_prompt=True,  # 是否在文本末尾添加生成提示。默认为 False 。 与 enable_thinking 配套使用，才能关闭思考模式。
    enable_thinking=False,       # 在思考和非思考模式之间切换。默认为 True 。
)
```

当我们打印 text 时，可以看到类似下面得内容：

```python
>>> print(text)
<|im_start|>user
hi, i'm DesistDaydream<|im_end|>
<|im_start|>assistant
<think>

</think>
```

> [!Attention] 这里面有点拧巴，其实应该只有 im_start 和 im_end 及其中间的部分，后面是因为关闭了思考模式出现的。这是因为 Qwen3-0.6B 训练时规定的就是这样的。

> [!tip]
> 这个示例使用得是 Qwen3-0.6B 模型。可以看到，该模型的 ControlTokens 是 `<|im_start|>`, `<|im_end|>`, etc. 。这就是不同模型的聊天模板不同导致的。每个模型配套的资源都是有关联的，我们不能拿 A 模型的分词器去给 B 模型使用，模型内部是无法识别这些 Token 的。

### 编码

接下来，需要对这些已经渲染好的字符串进行 **tokenize(分词)**，得到 **Token sequence(Token 序列)**（默认使用 merges.txt 文件？）

```python
token_sequence = tokenizer.tokenize(text)
```

token_sequence 结果如下：

```python
['<|im_start|>', 'user', 'Ċ', 'hi', ',', 'Ġi', "'m", 'ĠDes', 'ist', 'Day', 'dream', '<|im_end|>', 'Ċ', '<|im_start|>', 'assistant', 'Ċ', '<think>', 'ĊĊ', '</think>', 'ĊĊ']
```

> [!important]
> 这个数组里，每个元素都是一个 **Token**，合起来就是 Token sequence

使用 **Vocabulary(词表)** 将 Token sequence 转为 **Token IDs()**（默认使用 vocab.json 文件？）

```python
token_ids = tokenizer.convert_tokens_to_ids(token_sequence)
```

token_ids 结果如下：

```python
[151644, 872, 198, 6023, 11, 600, 2776, 3874, 380, 10159, 56191, 151645, 198, 151644, 77091, 198, 151667, 271, 151668, 271]
```

使用 Token IDs 构造 **Tensor(张量)**。

```python
input_ids = torch.tensor([token_ids]).to(model.device)
```

input_ids 结果如下：

```python
tensor([[151644,    872,    198,   6023,     11,    600,   2776,   3874,    380,
          10159,  56191, 151645,    198, 151644,  77091,    198, 151667,    271,
         151668,    271]], device='cuda:0')
```

根据不同的模型处理输入（因果模型、掩码模型），这里使用的示例是 Qwen3-0.6B，需要添加掩码

```python
attention_mask = torch.ones_like(input_ids)
inputs_ids_with_mask = {"input_ids": input_ids, "attention_mask": attention_mask}
```

inputs_ids_with_mask 结果如下：

```python
{'input_ids': tensor([[151644,    872,    198,   6023,     11,    600,   2776,   3874,    380,
          10159,  56191, 151645,    198, 151644,  77091,    198, 151667,    271,
         151668,    271]]), 'attention_mask': tensor([[1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1]])}
```

> [!important] 传入模型的内容本质就是 input_ids
> 只不过 Qwen3-0.6B 这个模型需要 attention_mask 参数。但是还有很多模型不需要。所以 mask 只是 input_ids 的附属物。本质还是 input_ids

## 模型推理

```python
output_ids = model.generate(**inputs_ids_with_mask, max_new_tokens=32768)
```

推理结果是一个 **Tensor(张量)**，由于推理是黑盒，每次推理结果可能并不相同，以其中一次结果为例：

```python
tensor([[151644,    872,    198,   6023,     11,    600,   2776,   3874,    380,
          10159,  56191, 151645,    198, 151644,  77091,    198, 151667,    271,
         151668,    271,   9707,     11,    358,   2776,   3874,    380,  10159,
          56191,      0,   3555,    594,    389,    697,   3971,     30,  11162,
            236,    101, 144232, 151645]], device='cuda:0')
```

之后我们需要将反向走一遍分词器处理输入的逻辑，把 Tensor 一步步解码成人类可读的信息

> [!TODO] 推理原理是黑盒，略。如果有数学功底的话，可以在这里引入对 模型内部推理逻辑 的笔记。

## 分词器处理输出

### 解码

解构 Tensor 到 Token IDs。

>为了性能，通常会在这里就去掉用户输入，使用 `[len(input_ids[0]) :]` 只保留模型生成的部分返回给用户。

```python
output_token_ids = output_ids[0][len(input_ids[0]) :].tolist()
```

结果如下：

```python
[9707, 11, 358, 2776, 3874, 380, 10159, 56191, 0, 3555, 594, 389, 697, 3971, 30, 11162, 236, 101, 144232, 151645]
```

去掉 Control Token

```python
output_token_ids_filtered = []
for tid in output_token_ids:
    if tid not in tokenizer.all_special_ids:
        output_token_ids_filtered.append(tid)
```

转换 Token IDs 为 Token sequence。

```python
output_token_sequence = tokenizer.convert_ids_to_tokens(output_token_ids)
```

结果如下：

```python
['Hello', ',', 'ĠI', "'m", 'ĠDes', 'ist', 'Day', 'dream', '!', 'ĠWhat', "'s", 'Ġon', 'Ġyour', 'Ġmind', '?', 'ĠðŁ', 'İ', '¨', 'âľ¨']
```

合词，将所有 Token 合并成字符串

```python
output_text = tokenizer.convert_tokens_to_string(output_token_sequence)  # type: ignore
```

结果如下：

```text
Hello, I'm DesistDaydream! What's on your mind? 🎨✨
```

---

下面的内容与解码流程无关

模型的推理输出的内容是这样的，上面 解构 Tensor 的时候我们把 think 以及前面的内容全都去掉了

```python
<|im_start|>user
hi, i'm DesistDaydream<|im_end|>
<|im_start|>assistant
<think>

</think>

Hello, I'm DesistDaydream! What's on your mind? 🎨✨<|im_end|>
```

> [!quote] 其实 `convert_ids_to_tokens()`, `decode()` 这些方法里都有有 skip_special_tokens 参数，可以处理 ControlToken，有太多太多的包装方法可以直接处理。只不过这里为了演示流程，所以都没用这些。

### 应用聊天模板

将模型的回复，填充到 `{"role": "assistant", "content": "AssistantReply"}` 的 AssistantReply 处。

随后将这条消息 append 到 message 中，以供下一次对话一起传模型继续聊天。

```python
messages.append({"role": "assistant", "content": output_text})
```

最终，messages 将会变为：

```json
[
  {
    "role": "user",
    "content": "hi, i'm DesistDaydream"
  },
  {
    "role": "assistant",
    "content": "Hello, I'm DesistDaydream! What's on your mind? 🎨✨"
  }
]
```

## 最终

至此，与模型的一次对话完成，后续再想对话，只需再次对 `messages.append({"role": "user", "content": prompt}})`，把历史上下文都带着再次执行一遍推理流程，即可通过重复此过程，我们可以随心所欲地持续聊天，**直到**<font color="#ff0000">模型超出上下文窗口或内存不足</font>为止。

实际上，相对最简单的方式是：

```python
# 输入
messages = [
    {"role": "user", "content": "hi, i'm DesistDaydream"},
]

# 编码
inputs_ids_with_mask = tokenizer.apply_chat_template(
    messages,  # 要使用模板处理的消息列表
    tokenize=True,  # 是否将文本转换为 Token 序列。默认为 Ture。关闭后只输出渲染后的文本。
    return_tensors="pt",  # 返回 Tensor。只有当 tokenize 为 True 时才有效。
    add_generation_prompt=True,  # 是否在文本末尾添加生成提示。默认为 False 。
    enable_thinking=False,  # 在思考和非思考模式之间切换。默认为 True 。
).to(model.device)

# 推理
output_ids = model.generate(**inputs_ids_with_mask, max_new_tokens=32768)  # type: ignore

# 解码
assistant_reply = tokenizer.decode(
    output_ids[0][len(inputs_ids_with_mask["input_ids"][0]) :],
    skip_special_tokens=True,
)

# 输出
messages.append({"role": "assistant", "content": assistant_reply})
```

transformers 库甚至还有更直接的 pipeline 方法。随着发展，肯定包装得会越来越完善，调用起来越来越简单。

# 数据结构

> 参考：
>
> - [Transformers 文档，与模型聊天 - 聊天模板](https://huggingface.co/docs/transformers/en/chat_templating)

发送给模型的消息通常包含两部分

- **conversation** # 对话，i.e. 人们日常交流以及代码里常用 **messages**。
- **tools** # 工具。LLM 会根据工具信息，以及 conversation 中的内容判断是否需要使用工具。

TODO: 待整理。

从上面的推理流程可以看出来

所有语言模型（无论是否经过聊天训练）都遵循一个 **TokensSequence(Token 序列)**。训练因果语言模型时，通常先在庞大的文本语料库上进行 “**pre-training(预训练)**”，从而创建一个“基础”模型。这些基础模型随后通常会针对聊天进行 “**fine-tuned(微调)**”，这意味着使用格式化为消息序列的数据进行训练。

## Conversation

conversation 可用的常见 role 包括：

- user # 来自用户的消息
- assistant # 模型推理的结果
- system # 用于指示模型如何行动的指令（通常位于聊天开始处）
- tool # 工具执行后的响应结果

常见的 messages。每一次对话都是一个元素。下面是一个只有用户输入 ”你好“ 的 messages

```json
[
    {"role": "user", "content": "你好"},
]
```

对于多模态模型来说，输入的 content 应该包含类型，比如：

```json
messages = [
    {"role": "user", "content": [{"type": "text", "text": "你好"}]},
]
```

> [!Attention] 非多模态的模型基本无法识别这种消息结构。

如果有 图片、音视频、etc. 、etc. ，还可以使用其他 type，并指定 文件系统路径 或 URL，由分词器读取后，转换为张量

## Tools

> 参考：
>
> - [Transformers 文档，与模型聊天 - 使用工具](https://huggingface.co/docs/transformers/en/chat_extras#passing-tools)

> [!Attention] 模型本身无法直接调用工具
>
> 模型会响应结构化的信息，由 `tool_call` 包裹，就像这样: `<tool_call>格式化的工具调用参数</tool_call>`。
>
> 我们需要编写程序处理这些信息，由程序去执行具体的工具。

一个工具的典型结构示例如下：

```json
tools := [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的天气",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {
                        "type": "string",
                        "description": "城市名称，例如：北京",
                    }
                },
                "required": ["city"],
            },
        },
    }
]
```

通常，模型响应的信息中，工具调用应该放在 `"role":"assistant"` 的 messages 中，使用 `tool_calls` 键

```json
{
  "role": "assistant",
  "content": "获取天气",
  "tool_calls": [
    {
      "id": "call_00_tO9WlaA7bYyLG89jjiVV4rwe",
      "type": "function",
      "function": {
        "name": "get_weather",
        "arguments": "{\"name\": \"get_weather\", \"arguments\": {\"city\": \"北京\"}}"
      }
    }
  ]
}
```

工具的执行结果应该放在 `"role":"tool"`。之后，程序将这些 messages 交给模型进行分析。继续执行后续任务

# 最佳实践

## 性能和内存使用情况

https://huggingface.co/docs/transformers/en/conversations#performance-and-memory-usage

Transformer 默认以全精度 float32 加载模型，对于一个 80 亿的模型来说，这需要大约 32GB 的内存！使用 torch_dtype="auto" 参数可以减少内存占用，该参数通常会对使用 bfloat16 训练的模型使用 bfloat16 精度。
