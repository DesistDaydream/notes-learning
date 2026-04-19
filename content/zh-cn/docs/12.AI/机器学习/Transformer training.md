---
title: Transformer training
linkTitle: Transformer training
created: 2026-04-13T15:17
weight: 101
---

# 概述

> 参考：
>
> -

# NLP 训练流程

![](Excalidraw/AI/transformer-training-flow.excalidraw.md)

# NLP 训练的关联文件与配置

> [!Quote] 参考 [Transformer](/docs/12.AI/机器学习/Transformer.md#模型的关联文件与配置) 中的 模型的关联文件与配置，不同架构的模型配置文件有细微差别

# TODO1

## `model.forward(**inputs)` 在做什么

`model` 是 `Qwen3ForCausalLM`，`forward()` 是它的前向传播方法。

`**inputs` 展开后大概包含：

- `input_ids` — token id 序列
- `attention_mask` — 哪些位置是真实 token，哪些是 padding
- `labels` — 监督信号，告诉模型"正确答案是什么"

模型拿到这些输入后，做的事情是：

```
input_ids → Embedding层 → 多层Transformer → logits（每个位置上词表的得分）
```

得到 logits 后，因为 `labels` 存在，模型内部会自动计算 **Cross-Entropy Loss**：

```
loss = CrossEntropy(logits, labels)
```

具体逻辑是把 `logits[0..n-1]` 和 `labels[1..n]` 对齐（即"用前一个 token 预测下一个 token"），计算平均负对数似然。

---

## `res["loss"]` 是什么

`res` 是一个 `CausalLMOutputWithPast` 对象，支持按 key 访问，包含：

|key|内容|
|---|---|
|`loss`|标量，当前 batch 的平均交叉熵损失|
|`logits`|shape `[batch, seq_len, vocab_size]` 的得分矩阵|
|`past_key_values`|KV cache，推理时用|

`res["loss"]` 就是取出那个**标量 loss**，后续：

```python
loss.backward()   # 从这个标量出发，反向传播，计算所有参数的梯度
optimizer.step()  # 按梯度更新参数，让 loss 下降
```

---

## 一句话总结

`model.forward(**inputs)` = 模型做一次完整的"预测 + 算分"；`res["loss"]` = 模型预测得有多错，这个值越小说明模型越接近训练目标。

# TODO2

## labels 是什么

### 语言模型的训练任务

语言模型做的事情本质上只有一件：**预测下一个 token**。

给定一段文字，模型要学会"看到前面，猜后面"：

```
输入序列：  我  爱  吃
预测目标：  爱  吃  饭
```

---

### labels 就是"正确答案序列"

假设训练数据是 `"我爱吃饭"`，tokenize 后得到：

```
input_ids: [我, 爱, 吃, 饭]
labels:    [我, 爱, 吃, 饭]   ← 通常和 input_ids 完全一样
```

但模型内部计算 loss 时，会自动做一个**错位对齐**：

```
模型看到：  我  爱  吃  饭
要预测的：  爱  吃  饭  ???
```

即：看到"我"，预测"爱"；看到"我爱"，预测"吃"……以此类推。

---

### 为什么叫"监督信号"

这是**有监督学习**的概念：

- 模型预测出一个概率分布（logits）
- labels 告诉模型正确答案是什么
- 对比两者，算出 loss
- 用 loss 反向传播，调整参数

没有 labels，模型就不知道自己预测得对不对，也就无法训练。

---

### 微调场景下的特殊处理

在指令微调（SFT）时，通常不希望模型去"学"用户的问题部分，只学模型回答部分。

所以 labels 里，**问题对应的位置会被设成 -100**：

```
input_ids: [你好,  请问,  天空,  为什么,  是,  蓝色,  的,  ？,  因为,  散射]
labels:    [-100, -100, -100,  -100,  -100, -100, -100, -100,  因为,  散射]
                  ↑ 这些位置计算 loss 时会被忽略          ↑ 只学这里
```

`-100` 是 PyTorch CrossEntropyLoss 的默认 `ignore_index`，遇到它直接跳过不算 loss。
