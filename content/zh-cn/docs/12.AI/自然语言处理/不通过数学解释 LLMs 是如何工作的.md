---
title: 不通过数学解释 LLMs 是如何工作的
linkTitle: 不通过数学解释 LLMs 是如何工作的
weight: 20
---

# 概述

> 参考：
>
> - https://blog.miguelgrinberg.com/post/how-llms-work-explained-without-math
>   - [中文翻译：公众号 - 云原生实验室，大模型到底有没有智能？一篇文章给你讲明明白白](https://mp.weixin.qq.com/s/z93ZEVnjdpKM__-_0HvjRQ)

生成式人工智能 ([GenAI](https://en.wikipedia.org/wiki/Generative_artificial_intelligence)) 和大语言模型 (LLM[^LLM])，这两个词汇想必已在大家的耳边萦绕多时。它们如惊涛骇浪般席卷了整个科技界，登上了各大新闻头条。ChatGPT，这个神奇的对话助手，也许已成为你形影不离的良师益友。

[^LLM]:LLM: https://en.wikipedia.org/wiki/Large_language_model

然而，在这场方兴未艾的 GenAI 革命背后，有一个谜题久久萦绕在人们心头：**这些模型的智能究竟从何而来**？本文将为您揭开谜底，解析生成式文本模型的奥秘。我们将抛开晦涩艰深的数学，用通俗易懂的语言，带您走进这个神奇的算法世界。让我们撕下 “魔法” 的面纱，看清其中的计算机科学本质。

## LLM 的真面目

首先，我们要破除一个常见的误区。许多人误以为，这些模型是真的能够与人对话，回答人们的各种问题。然而，它们真正的能力远没有想象的那么复杂——**它们所做的，不过是根据输入的文本，预测下一个词语** (更准确地说，是下一个 token)。

Token，这个看似简单的概念，却是揭开 LLM 神秘面纱的钥匙。让我们由此出发，步步深入，一探究竟。

Token，这些文本的积木、语言的原子，正是 LLM 理解世界的基石。对我们而言，token 不过是单词、标点、空格的化身，但在 LLM 的眼中，它们是精简而高效的信息编码。有时，一个 token 可能代表一串字符，长短不一；有时，它可能是孤零零的一个标点符号。

**LLM 的词汇表，就是这些 token 的集合**，啥都有，样样全。这其中的奥秘，要追溯到 [BPE 算法](https://en.wikipedia.org/wiki/Byte_pair_encoding)。BPE 算法是如何炼制出这些 tokens 的？这个问题，值得我们细细探究。但在此之前，只需记住：[GPT-2 模型](https://github.com/openai/gpt-2)，这个自然语言处理界的明星，它的词汇表中有 50,257 个 token。

在 LLM 的世界里，每个 token 都有一个独一无二的数字身份证。而 Tokenizer，就是文本和 token 之间的 “翻译官”，将人类的语言转化为 LLM 能理解的编码，也将 LLM 的思维解码为人类的文字。如果你熟悉 Python，不妨亲自与 token 打个照面。只需安装 OpenAI 的 `tiktoken` 包：

```bash
pip install tiktoken
```

`

然后在 Python 中尝试以下操作：

```python
>>> import tiktoken
>>> encoding = tiktoken.encoding_for_model("gpt-2")

>>> encoding.encode("The quick brown fox jumps over the lazy dog.")
[464, 2068, 7586, 21831, 18045, 625, 262, 16931, 3290, 13]

>>> encoding.decode([464, 2068, 7586, 21831, 18045, 625, 262, 16931, 3290, 13])
'The quick brown fox jumps over the lazy dog.'

>>> encoding.decode([464])
'The'
>>> encoding.decode([2068])
' quick'
>>> encoding.decode([13])
'.'
```

在这个实验中，我们可以看到，对于 GPT-2 而言，token 464 表示单词 “The”，token 2068 表示 “quick” (含前导空格)，token 13 则表示句点。

由于 token 是算法生成的，有时会出现一些奇怪的现象。比如，同一个单词 “the” 的三个变体在 GPT-2 中被编码成了不同的 token：

```python
>>> encoding.encode('The')
[464]
>>> encoding.encode('the')
[1169]
>>> encoding.encode(' the')
[262]
```

BPE 算法并不总是将完整的单词映射为 token。事实上，不太常用的单词可能无法成为单独的 token，需要用多个 token 组合编码，比如这个 “Payment”，就要化身为 “Pay” 和 “ment” 的组合：

```python
>>> encoding.encode("Payment")
[19197, 434]

>>> encoding.decode([19197])
'Pay'
>>> encoding.decode([434])
'ment'
```

`

### 预测下一个 Token

语言模型就像一个 “水晶球”，给它一串文字，它就能预言下一个最可能出现的词语。这是它的看家本领。但模型并非真的手眼通天，它的预言能力其实基于扎实的概率计算。让我们一起掀开这层神秘的面纱，看看背后的真相。

如果你懂一点 Python，我们可以用几行代码来窥探语言模型的预言过程：

```python
predictions = get_token_predictions(['The', ' quick', ' brown', ' fox'])
```

`

这个 `get_token_predictions` 函数就是我们的 “水晶球”。它接受一个 token 列表作为输入，这些 token 来自用户提供的 prompt。在这个例子中，我们假设每个单词都是一个独立的 token。当然，在实际使用中，每个 token 都有一个对应的数字 ID，但为了简单起见，我们这里直接用单词的文本形式。

函数的返回结果是一个庞大的数据结构，里面记录了词汇表中每个 token 出现在输入文本之后的概率。以 GPT-2 模型为例，它的词汇表包含 50,257 个 token，因此返回值就是一个 50,257 维的概率分布。

现在再来重新审视一下这个例子。如果我们的语言模型训练有素，面对 “[The quick brown fox](https://en.wikipedia.org/wiki/The_quick_brown_fox_jumps_over_the_lazy_dog)” 这样一个烂大街的句子片段，它很可能会预测下一个词是 “jumps”，而不是 “potato” 之类风马牛不相及的词。在这个概率分布中，“jumps” 的概率值会非常高，而 “potato” 的概率值则接近于零。

> [!Tip]
>
> **The quick brown fox jumps over the lazy dog** (相应中文可简译为 “快狐跨懒狗”，完整翻译则是 “敏捷的棕色狐狸跨过懒狗”) 是一个著名的英语全字母句，常用于测试字体显示效果和键盘是否故障。此句也常以 “quick brown fox” 做为指代简称。

当然，语言模型的预测能力并非与生俱来，而是通过日积月累的训练得来的。在漫长的训练过程中，模型如饥似渴地汲取海量文本的营养，逐渐茁壮成长。训练结束时，它已经具备了应对各种文本输入的能力，可以利用积累的知识和经验，计算出任意 token 序列的下一个 token 概率。

现在是不是觉得语言模型的预测过程没那么神奇了？它与其说是魔法，不如说是一个基于概率的计算过程。这个过程虽然复杂，但并非不可理解。我们只要掌握了基本原理，就可以揭开它的神秘面纱，走近它，了解它。

### 长文本生成的奥秘

由于语言模型每次只能预测下一个 token 会是什么，因此生成完整句子的唯一方法就是在循环中多次运行该模型。每一轮迭代都会从模型返回的概率分布中选择一个新的 token，生成新的内容。然后将这个新 token 附加到下一轮中输入给模型的文本序列末尾，如此循环往复，直到生成足够长度的文本。

我们来看一个更完整的 Python 伪代码，展示具体的实现逻辑：

```python
def generate_text(prompt, num_tokens, hyperparameters):
    tokens = tokenize(prompt)
    for i in range(num_tokens):
        predictions = get_token_predictions(tokens)
        next_token = select_next_token(predictions, hyperparameters)
        tokens.append(next_token)
    return ''.join(tokens)
```

其中，`generate_text()` 函数接受一个用户输入的提示词 (prompt) 文本作为参数，这可以是一个问题或其他任意文本。

`tokenize()` 辅助函数使用类似 `tiktoken` 的分词库将提示文本转换成一系列等效的 `token(token)` 序列。在 for 循环内部，get\_token\_predictions() 函数调用语言模型来获取下一个 token 的概率分布，这一步与前面的示例类似。

`select_next_token()` 函数根据上一步得到的下个 token 概率分布，选择最合适的 token 来延续当前的文本序列。最简单的做法是选择概率最高的 token，在机器学习中被称为贪婪选择 (greedy selection)。更好的做法是用符合模型给出概率分布的随机数生成器来选词，这样可以让生成的文本更丰富多样。如果用同样的输入多次运行模型，这种方法还可以让每次产生的回应都略有不同。

为了让 token 选择过程更加灵活可控，可以用一些超参数 (hyperparameter) 来调整语言模型返回的概率分布，这些超参数作为参数传递给文本生成函数。通过调整超参数，你可以控制 token 选择的 “贪婪程度”。如果你用过大语言模型，你可能比较熟悉名为 `temperature` 的超参数。提高 temperature 的值可以让 token 的概率分布变得更加平缓，增加选中概率较低 token 的机会，从而让生成的文本显得更有创意和变化。此外，常用的还有 `top_p` 和 `top_k` 两个超参数，它们限定从概率最高的前 k 个或概率超过阈值 p 的 token 中进行选择，以平衡多样性和连贯性。

选定了一个新 token 后，循环进入下一轮迭代，将新 token 添加到原有文本序列的末尾，作为新一轮的输入，再接着生成下一个 token。`num_tokens` 参数控制循环的迭代轮数，决定要生成的文本长度。但需要注意的是，由于语言模型是逐词预测，没有句子或段落的概念，生成的文本常常会在句子中途意外结束。为了避免这种情况，我们可以把 `num_tokens` 参数视为生成长度的上限而非确切值，当遇到句号、问号等标点符号时提前结束生成过程，以保证文本在语义和语法上的完整性。

如果你已经读到这里且充分理解了以上内容，那么恭喜你！现在你对大语言模型的基本工作原理有了一个高层次的认识。如果你想进一步了解更多细节，我在下一节会深入探讨一些更加技术性的话题，但会尽量避免过多涉及晦涩难懂的数学原理。

## 模型训练

遗憾的是，不借助数学语言来讨论模型训练实际上是很困难的。这里先展示一种非常简单的训练方法。

既然 LLM 的任务是预测某些词后面跟随的词，那么一个简单的模型训练方式就是从训练数据集中提取所有连续的词对，并用它们来构建一张概率表。

让我们用一个小型词表和数据集来演示这个过程。假设模型的词表包含以下 5 个词：

```python
['I', 'you', 'like', 'apples', 'bananas']
```

为了保持示例简洁，我不打算将空格和标点符号视为独立的词。

我们使用由三个句子组成的训练数据集：

- I like apples
- I like bananas
- you like bananas

我们可以构建一个 5x5 的表格，在每个单元格中记录 “该单元格所在行的词” 后面跟随 “该单元格所在列的词” 的次数。下面是根据数据集中三个句子得到的表格：

| \- | I | you | like | apples | bananas |
| --- | --- | --- | --- | --- | --- |
| **I** |  |  | 2 |  |  |
| **you** |  |  | 1 |  |  |
| **like** |  |  |  | 1 | 2 |
| **apples** |  |  |  |  |  |
| **bananas** |  |  |  |  |  |
这个表格应该不难理解。数据集中包含两个 “I like” 实例，一个 “you like” 实例，一个 “like apples” 实例和两个 “like bananas” 实例。

现在我们知道了每对词在训练集中出现的次数，就可以计算每个词后面跟随其他词的概率了。为此，我们将表中每一行的数字转换为概率值。例如，表格中间行的 “like” 后面有一次跟随 “apples”，两次跟随 “bananas”。这意味着在 33.3%的情况下 “like” 后面是 “apples”，剩下 66.7%的情况下是 “bananas”。

下面是计算出所有概率后的完整表格。空单元格代表 0%的概率。

| \- | I | you | like | apples | bananas |
| --- | --- | --- | --- | --- | --- |
| **I** |  |  | 100% |  |  |
| **you** |  |  | 100% |  |  |
| **like** |  |  |  | 33.3% | 66.7% |
| **apples** | 25% | 25% | 25% |  | 25% |
| **bananas** | 25% | 25% | 25% | 25% |  |

“I”、“you” 和 “like” 这几行的概率很容易计算，但 “apples” 和 “bananas” 带来了问题。由于数据集中没有这两个词后面接其他词的例子，它们存在训练数据的空白。为了确保模型即使面对未见过的词也能做出预测，我决定将 “apples” 和 “bananas” 的后续词概率平均分配给其他四个可能的词。这种做法虽然可能产生不自然的结果，但至少能防止模型在遇到这两个词时陷入死循环。

训练数据存在 “空洞” 的问题对语言模型的影响不容忽视。**在真实的大语言模型中，由于训练语料极其庞大，这些空洞通常表现为局部覆盖率偏低，而不是整体性的缺失，因而不太容易被发现。语言模型在这些训练不足的领域或话题上会产生片面、错误或前后不一致的预测结果，但通常会以一种难以感知的形式表现出来。这就是语言模型有时会产生 “[Hallucination(幻觉)](https://en.wikipedia.org/wiki/Hallucination_(artificial_intelligence))” 的原因之一，所谓幻觉，就是指生成的文本表面上读起来通顺流畅，但实际包含了事实错误或前后矛盾的内容。**

借助上面给出的概率表，你现在可以自己想象一下 `get_token_predictions()` 函数会如何实现。用 Python 伪代码表示大致如下：

```python
def get_token_predictions(input_tokens):
    last_token = input_tokens[-1]
    return probabilities_table[last_token]
```

是不是比想象的要简单？该函数接受一个单词序列作为输入，也就是用户提示。它取这个序列的最后一个单词，然后返回概率表中与之对应的那一行。

举个例子，如果用 `['you', 'like']` 来调用这个函数，它会返回 “like” 所在的行，其中 “apples” 有 33.3%的概率接在后面组成句子，而 “bananas” 占剩下的 66.7%。有了这些概率信息，之前展示的 `select_next_token()` 函数在三分之一的情况下应该选择 “apples”。

当 “apples” 被选为 “you like” 的续词时，“you like apples” 这个句子就形成了。**这是一个在训练数据中不存在的全新句子，但它却非常合理**。希望这个例子能帮你认识到，**语言模型其实只是在重复使用和拼凑它在训练过程中学到的各种模式碎片，就能组合出一些看似原创的想法或概念**。

### 上下文窗口

上一节内容我使用 [Markov chain(马尔可夫链)](https://en.wikipedia.org/wiki/Markov_chain) 的方法训练了一个小语言模型。这种方法存在一个问题：它的上下文窗口只有一个标记，也就是说，**模型在预测下一个词时，只考虑了输入序列的最后一个词，而忽略了之前的所有内容**。这导致生成的文本缺乏连贯性和一致性，常常前后矛盾，逻辑跳跃。

为了提高模型的预测质量，一种直观的思路是**扩大上下文窗口的大小，比如增加到 2 个标记**。但这样做会导致概率表的规模急剧膨胀。以我之前使用的 5 个标记的简单词表为例，将上下文窗口增加到 2 个标记，就需要在原有的 5 行概率表基础上，额外增加 25 行来覆盖所有可能的双词组合。如果进一步扩大到 3 个标记，额外的行数将达到 125 行。可以预见，**随着上下文窗口的增大，概率表的规模将呈指数级爆炸式增长**。

更重要的是，即使将上下文窗口扩大到 2 个或 3 个标记，其改进效果仍然非常有限。要使语言模型生成的文本真正做到前后连贯、逻辑通顺，实际上需要一个远大于此的上下文窗口。**只有足够大的上下文，新生成的词才能与之前较远处提及的概念、思想产生联系，从而赋予文本连续的语义和逻辑**。

举个实际的例子，OpenAI 开源的 GPT-2 模型采用了 1024 个标记的上下文窗口。如果仍然沿用马尔可夫链的思路来实现这一尺度的上下文，**以 5 个标记的词表为例，仅覆盖 1024 个词长度的所有可能序列，就需要高达 $5^{1024}$ 行的概率表**。这是一个天文数字，我在 Python 中计算了这个值的具体大小，读者可以向右滚动来查看完整的数字：

```python
>>> pow(5, 1024)
55626846462680034577255817933310101605480399511558295763833185422180110870347954896357078975312775514101683493275895275128810854038836502721400309634442970528269449838300058261990253686064590901798039126173562593355209381270166265416453973718012279499214790991212515897719252957621869994522193843748736289511290126272884996414561770466127838448395124802899527144151299810833802858809753719892490239782222290074816037776586657834841586939662825734294051183140794537141608771803070715941051121170285190347786926570042246331102750604036185540464179153763503857127117918822547579033069472418242684328083352174724579376695971173152319349449321466491373527284227385153411689217559966957882267024615430273115634918212890625
```

这段 Python 代码示例生成了一个庞大的表格，但即便如此，它也只是整个表格的一小部分。因为除了当前的 1024 个 token 长度的序列，我们还需要生成更短的序列，譬如 1023 个、1022 个 token 的序列，一直到只包含 1 个 token 的序列。这样做是为了确保在输入数据 token 数量不足的情况下，模型也能妥善处理较短的序列。马尔可夫链虽然是一个有趣的文本生成方法，但在可扩展性方面确实存在很大的问题。

如今，1024 个 token 的上下文窗口已经不那么出色了。GPT-3 将其扩大到了 2048 个 token，GPT-3.5 进一步增加到 4096 个。GPT-4 一开始支持 8192 个 token 的上下文，后来增加到 32000 个，再后来甚至达到了 128000 个 token！目前，开始出现支持 100 万以上 token 的超大上下文窗口模型，使得模型在进行 token 预测时，能够拥有更好的一致性和更强的记忆能力。

总而言之，尽管马尔可夫链为我们提供了一种正确的思路来思考文本生成问题，但其固有的可扩展性不足，使其难以成为一个可行的、能够满足实际需求的解决方案。面对海量文本数据，我们需要寻求更加高效和可扩展的文本生成方法。

### 从马尔可夫链到神经网络

显然，我们必须摒弃使用概率表的想法。对于一个合理大小的上下文窗口，所需的表格大小将远超内存限制。**我们可以用一个函数来代替这个表格，该函数能够通过算法生成近似的下一个词出现概率，而无需将其存储在一个巨大的表格中**。这正是神经网络擅长的领域。

神经网络是一种特殊的函数，它接收一些输入，经过计算后给出输出。对于语言模型而言，输入是代表提示信息的词，输出是下一个可能出现的词及其概率列表。神经网络之所以特殊，是因为除了函数逻辑之外，它们对输入进行计算的方式还受到许多外部定义参数的控制。

最初，神经网络的参数是未知的，因此其输出毫无意义。神经网络的训练过程就是要找到那些能让函数在训练数据集上表现最佳的参数，并假设如果函数在训练数据上表现良好，它在其他数据上的表现也会相当不错。

在训练过程中，参数会使用一种叫做 [Backpropagation(反向传播)](https://en.wikipedia.org/wiki/Backpropagation) 的算法进行迭代调整，每次调整的幅度都很小。这个算法涉及大量数学计算，我们在这里就不展开了。**每次参数调整后，神经网络的预测都会变得更准一些**。参数更新后，网络会用训练数据集重新评估，结果为下一轮调整提供参考。这个过程会反复进行，直到函数能够在训练数据上很好地预测下一个词。

为了让你对神经网络的规模有个概念，**GPT-2 模型有大约 15 亿个参数，GPT-3 增加到了 1750 亿，而 GPT-4 据说有 1.76 万亿个参数**。在当前硬件条件下，训练如此规模的神经网络通常需要几周或几个月的时间。

有趣的是，由于参数数量巨大，并且都是通过漫长的迭代过程自动计算出来的，我们很难理解模型的工作原理。**训练完成的大语言模型就像一个难以解释的黑箱，因为模型的大部分 “思考” 过程都隐藏在海量参数之中。即使是训练它的人，也很难说清其内部的运作机制。**

### 层、Transformer 与 Attention 机制

你可能好奇神经网络函数内部进行了哪些神秘的计算。在精心调校的参数帮助下，它可以接收一系列输入标记，并以某种方式输出下一个标记出现的合理概率。

神经网络被配置为执行一系列操作，**每个操作称为一个 “层”**。第一层接收输入并对其进行转换。转换后的输入进入下一层，再次被转换。这一过程持续进行，直到数据到达最后一层并完成最终转换，生成输出或预测结果。

机器学习专家设计出不同类型的层，对输入数据进行数学转换。他们还探索了组织和分组层的方法，以实现期望的结果。有些层是通用的，而另一些则专门处理特定类型的输入数据，如图像，或者在大语言模型中的标记化文本。

**目前在大语言模型的文本生成任务中最流行的神经网络架构被称为 “[Transformer](https://en.wikipedia.org/wiki/Transformer_(deep_learning_architecture)”。使用这种架构的模型被称为 GPT，也就是 [Generative Pre-Trained Transformers(生成式预训练 Transformer)](https://en.wikipedia.org/wiki/Generative_pre-trained_transformer。**

Transformer 模型的独特之处在于其执行的 “[Attention](https://en.wikipedia.org/wiki/Attention_(machine_learning)” 层计算。这种计算允许模型**在上下文窗口内的标记之间找出关系和模式，并将其反映在下一个标记出现的概率中**。Attention 机制最初被用于语言翻译领域，作为一种找出输入序列中对理解句子意义最重要的标记的方法。这种机制赋予了现代语言模型在基本层面上 “理解” 句子的能力，它可以关注 (或集中 “注意力” 于) 重要词汇或标记，从而更好地把握句子的整体意义。正是这一机制，使 Transformer 模型在各种自然语言处理任务中取得了巨大成功。

## 大语言模型到底有没有智能？

通过上面的分析，你心中可能已经有了一个初步的判断：大语言模型在生成文本时是否表现出了某种形式的智能？

我个人并不认为大语言模型具备推理或提出原创想法的能力，但这并不意味着它们一无是处。得益于对上下文窗口中 token 进行的精妙计算，大语言模型能够捕捉用户输入中的模式，并将其与训练过程中学习到的相似模式匹配。**它们生成的文本大部分来自训练数据的片段，但将词语 (实际上是 token) 组合在一起的方式非常复杂，在许多情况下产生了感觉原创且有用的结果**。

不过，考虑到大语言模型容易产生幻觉，我不会信任任何将其输出直接提供给最终用户而不经过人工验证的工作流程。

未来几个月或几年内出现的更大规模语言模型是否能实现类似真正智能的能力？鉴于 GPT 架构的诸多局限性，我觉得这不太可能发生，但谁又说的准呢，也许将来出现一些创新手段，我们就能实现这一目标。
