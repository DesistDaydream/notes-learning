---
title: "Markdown"
linkTitle: "Markdown"
weight: 20
---

# 概述
>
> 参考：
>
> - [官网](https://commonmark.org/)
> - [官方规范](https://spec.commonmark.org/)
> - [Wiki, Markdown](https://en.wikipedia.org/wiki/Markdown)
> - [GitHub 项目，DavidAnson/markdownlint](https://github.com/DavidAnson/markdownlint/tree/main/doc) # 一种 MarkDown 的格式规范

Markdown 是一种轻量级的标记语言

# Markdown 中的 LaTex 数学公式

> 参考：
>
> - [Markdwon + LaTex 表达数学式子](https://blog.ypingcn.com/notes/Markdown/LaTex-math/)
> - [简书，Markdown 数学公式语法](https://www.jianshu.com/p/e74eb43960a1)

[基础的 Markdown 语法](https://blog.ypingcn.com/notes/Markdown/basic/) 中无法满足数学公式的表达需求，此时可以借助 Latex 语法完成。在 Markdown 中由 `$` 符号包围的部分编写的 Latex 语法，可以解析成数学公式。一共有两种格式

- **行内格式** # 使用 `$`。比如 `$2^{10}=1024$` 的解析效果为：$2^{10}=1024$
- **独行格式** # 使用 `$$`。比如 `$$LaTex code$$` 的解析效果如下

$$LaTex code$$


## 速查

Markdown 中 Latex 基本符号速查表

| 显示字符       | 输入字符             | 显示字符      | 输入字符              | 显示字符 | 输入字符                |
| ---------- | ---------------- | --------- | ----------------- | ---- | ------------------- |
| `#`        | `\#`             | `$`       | `\$`              | `%`  | `\%`                |
| `&`        | `\&`             | `~`       | `\~`              | `_`  | `\_`                |
| `^`        | `\^`             | `\\`      | `\\`              | `{`  | `\{`                |
| `}`        | `\}`             |           |                   |      |                     |
| ≤          | `\le`            | ≥         | `\ge`             | ≡    | `\equiv`            |
| ≠          | `\ne`            |           |                   |      |                     |
| 文本底线对齐的省略号 | `\ldots`         | 文本中对齐的省略号 | `\cdots`          |      |                     |
| 圆括号        | `()`             | 方括号       | `[]`              | 竖线   | ``                  |
| 花括号        | `\{\}`           | 双竖线       | `\`               |      |                     |
| 长圆括号       | `\left( \right)` | 长方括号      | `\left[ \\right]` | 长花括号 | `\left\\{ \right\}` |
| 换行         | `\\`             | 空格        | `\space`          |      |                     |
| ←          | `\leftarrow`     | →         | `\rightarrow`     | 文字   | `\mbox{ }`          |

## 字符相关

### **字符**

插入 `# $ % & ~ _ ^ \ { }` 这些符号类字符需多加 `\` 符号（类似于编程语言中的转义字符），其他可以直接插入。

`\\` 对应换行符 ，`\space` 对应空格。

小于等于、大于等于、恒等于、不等于分别为 `\le \ge \equiv \ne` （ l 意为 less ，e 意为 equal ，g 意为 greater ）

### **省略号**

`\ldots \cdots` 分别表示与文本底线对齐和与文本中对齐的省略号。（l 意为 line ，c 意为 center ）

### **括号**

圆括号、方括号和竖线直接输入，花括号前需添加 `\` ，双竖线对应`\|` 。

长圆括号、长方括号、长花括号对应`\left( \right)` `\left[ \right]` `\left\{ \right\}`

### **箭头**

左右箭头对应 `\leftarrow \rightarrow`

### **插入文字**

`\mbox{ }` 用于插入文字（显示效果不是斜体字） 。

## 数学式子

### **标准函数**

欲输入 `sin` 时，应用 `\sin(x)`

### **分数**

二分之一对于 `\frac{1}{2}`

### **根号**

根号二对应 `\sqrt{2}` ，开 n 次方为`\sqrt[n]{expression}`

### **导数 偏导数**

对 x 导数对应 `\mathrm{d}x`

对 x 的偏导数对应`\partial x`

### **积分**

f(x) 对 x 从 a 到 b 的积分 `\int_a^b f(x) \mathrm{d}x` ， 多重积分则多次输入 `\int` ，两个符号之间添加 `\!\!\!` 调整正确的间隔。

### **极限**

x 到正无穷的极限 `\lim_{x\to+\infty}`

### **求和**

1 到 n 的和对应 `\sum_{1}^{n}`

### **向量**

向量 ab 对应 `\vec{ab}`

### **排列组合**

从 n 中选 m 的组合数和排列数为 `\mathrm{C}_n^m \mathrm{A}_n^m`

[具体内容参见文末参考资料](#1)。

## 上下标与希腊字母

`^` 表示上标 `_`表示下标，同时出现上下标时，先上标后下标与先下标后上标的效果相同。

用 `\` 加相应的拼写即可，第一个字母大写则显示大写字母，小写则显示小写字母。

| 输入          | 展示 | 输入     | 展示 |
| ------------- | ---- | -------- | ---- |
| \alpha        | α    | \beta    | β    |
| \gamma        | γ    | \Gamma   | Γ    |
| \theta        | θ    | \Theta   | Θ    |
| \delta        | δ    | \Delta   | Δ    |
| \triangledown | ▽    | \epsilon | ϵ    |
| \zeta         | ζ    | \eta     | η    |
| \kappa        | κ    | \lambda  | λ    |
| \mu           | μ    | \nu      | ν    |
| \xi           | ξ    | \pi      | π    |
| \sigma        | σ    | \tau     | τ    |
| \upsilon      | υ    | \phi     | ϕ    |
| \omega        | ω    |          |      |

## 多行的数学公式

例子

```latex
\begin{eqnarray}
\cos 2\theta & = & \cos^2 \theta - \sin^2 \theta \\
& = & 2 \cos^2 \theta - 1.
\end{eqnarray}
```

解析结果

$$
\begin{eqnarray}
\cos 2\theta & = & \cos^2 \theta - \sin^2 \theta \\
& = & 2 \cos^2 \theta - 1.
\end{eqnarray}
$$

& 是对齐点，具体例子中表示多行式子在等号之间对齐。

```latex
f(n) =
\begin{cases}
n+1,  & \mbox{if }n \mbox{ is even} \\
n-1, & \mbox{if }n \mbox{ is odd}
\end{cases}
```

条件定义式。奇数加一，偶数减一。解析结果

$$f(n) =
\begin{cases}
n+1,  & \mbox{if }n \mbox{ is even} \\
n-1, & \mbox{if }n \mbox{ is odd}
\end{cases}$$

## 矩阵

表示 3 x 3 的矩阵，c 表示居中对齐，l 是左对齐，r 是右对齐

```latex
\begin{array}{ccc}
a & b & c \\
d & e & f \\
g & h & i
\end{array}
```

解析结果：

$$
\begin{array}{ccc}
a & b & c \\
d & e & f \\
g & h & i
\end{array}
$$

参考资料 ：

\#1 [帮助:数学公式 - 维基百科，自由的百科全书](https://zh.wikipedia.org/wiki/Help:%E6%95%B0%E5%AD%A6%E5%85%AC%E5%BC%8F)

\#2 [LaTeX/数学公式 - 维基教科书，自由的教学读本](https://zh.wikibooks.org/wiki/LaTeX/%E6%95%B0%E5%AD%A6%E5%85%AC%E5%BC%8F)

# EXAMPLE

上标、下标与组合
--------

1.  上标符号，符号：`^`，如：$x^4$
2.  下标符号，符号：`_`，如：$x_1$
3.  组合符号，符号：`{}`，如：${16}_{8}O{2+}_{2}$

汉字、字体与格式
--------

- 汉字形式，符号：`\mbox{}`，如：$V_{\mbox{初始}}$
- 字体控制，符号：`\displaystyle`，如：$\displaystyle \frac{x+y}{y+z}$
- 下划线符号，符号：`\underline`，如：$\underline{x+y}$
- 标签，符号`\tag{数字}`，如：$\tag{11}$
- 上大括号，符号：`\overbrace{算式}`，如：$\overbrace{a+b+c+d}^{2.0}$
- 下大括号，符号：`\underbrace{算式}`，如：$a+\underbrace{b+c}_{1.0}+d$
- 上位符号，符号：`\stacrel{上位符号}{基位符号}`，如：$\vec{x}\stackrel{\mathrm{def}}{=}{x_1,\dots,x_n}$

占位符
---

- 两个quad空格，符号：`\qquad`，如：$x \qquad y$
- quad空格，符号：`\quad`，如：$x \quad y$
- 大空格，符号`\`，如：$x \ y$
- 中空格，符号`\:`，如：$x : y$
- 小空格，符号`\,`，如：$x , y$
- 没有空格，符号``，如：$xy$
- 紧贴，符号`\!`，如：$x ! y$

定界符与组合
------

- 括号，符号：`（）\big(\big) \Big(\Big) \bigg(\bigg) \Bigg(\Bigg)`，如：$（）\big(\big) \Big(\Big) \bigg(\bigg) \Bigg(\Bigg)$
- 中括号，符号：`[]`，如：$[x+y]$
- 大括号，符号：`\{ \}`，如：${x+y}$
- 自适应括号，符号：`\left \right`，如：$\left(x\right)$，$\left(x{yz}\right)$
- 组合公式，符号：`{上位公式 \choose 下位公式}`，如：${n+1 \choose k}={n \choose k}+{n \choose k-1}$
- 组合公式，符号：`{上位公式 \atop 下位公式}`，如：$\sum_{k_0,k_1,\ldots>0 \atop k_0+k_1+\cdots=n}A_{k_0}A_{k_1}\cdots$

四则运算
----

- 加法运算，符号：`+`，如：$x+y=z$
- 减法运算，符号：`-`，如：$x-y=z$
- 加减运算，符号：`\pm`，如：$x \pm y=z$
- 减甲运算，符号：`\mp`，如：$x \mp y=z$
- 乘法运算，符号：`\times`，如：$x \times y=z$
- 点乘运算，符号：`\cdot`，如：$x \cdot y=z$
- 星乘运算，符号：`\ast`，如：$x \ast y=z$
- 除法运算，符号：`\div`，如：$x \div y=z$
- 斜法运算，符号：`/`，如：$x/y=z$
- 分式表示，符号：`\frac{分子}{分母}`，如：$\frac{x+y}{y+z}$
- 分式表示，符号：`{分子} \voer {分母}`，如：${x+y} \over {y+z}$
- 绝对值表示，符号：`||`，如：$|x+y|$

高级运算
----

1.  平均数运算，符号：`\overline{算式}`，如：$\overline{xyz}$
2.  开二次方运算，符号：`\sqrt`，如：$\\sqrt x$
3.  开方运算，符号：`\sqrt[开方数]{被开方数}`，如：$\sqrt[3]{x+y}$
4.  对数运算，符号：`\log`，如：$\log(x)$
5.  极限运算，符号：`\lim`，如：$\lim^{x \to \infty}_{y \to 0}{\frac{x}{y}}$
6.  极限运算，符号：`\displaystyle \lim`，如：$\displaystyle \lim^{x \to \infty}_{y \to 0}{\frac{x}{y}}$
7.  求和运算，符号：`\sum`，如：$\sum^{x \to \infty}_{y \to 0}{\frac{x}{y}}$
8.  求和运算，符号：`\displaystyle \sum`，如：$\displaystyle \sum^{x \to \infty}_{y \to 0}{\frac{x}{y}}$
9.  积分运算，符号：`\int`，如：$\int^{\infty}_{0}{xdx}$
10.  积分运算，符号：`\displaystyle \int`，如：$\displaystyle \int^{\infty}_{0}{xdx}$
11.  微分运算，符号：`\partial`，如：$\frac{\partial x}{\partial y}$
12.  矩阵表示，符号：`\begin{matrix} \end{matrix}`，如：$\left[ \begin{matrix} 1 &2 &\cdots &4\5 &6 &\cdots &8\\vdots &\vdots &\ddots &\vdots\13 &14 &\cdots &16\end{matrix} \right]$

逻辑运算
----

- 等于运算，符号：`=`，如：$x+y=z$
- 大于运算，符号：`>`，如：$x+y>z$
- 小于运算，符号：`<`，如：$x+y<z$
- 大于等于运算，符号：`\geq`，如：$x+y \geq z$
- 小于等于运算，符号：`\leq`，如：$x+y \leq z$
- 不等于运算，符号：`\neq`，如：$x+y \neq z$
- 不大于等于运算，符号：`\ngeq`，如：$x+y \ngeq z$
- 不大于等于运算，符号：`\not\geq`，如：$x+y \not\geq z$
- 不小于等于运算，符号：`\nleq`，如：$x+y \nleq z$
- 不小于等于运算，符号：`\not\leq`，如：$x+y \not\leq z$
- 约等于运算，符号：`\approx`，如：$x+y \approx z$
- 恒定等于运算，符号：`\equiv`，如：$x+y \equiv z$

集合运算
----

- 属于运算，符号：`\in`，如：$x \in y$
- 不属于运算，符号：`\notin`，如：$x \notin y$
- 不属于运算，符号：`\not\in`，如：$x \not\in y$
- 子集运算，符号：`\subset`，如：$x \subset y$
- 子集运算，符号：`\supset`，如：$x \supset y$
- 真子集运算，符号：`\subseteq`，如：$x \subseteq y$
- 非真子集运算，符号：`\subsetneq`，如：$x \subsetneq y$
- 真子集运算，符号：`\supseteq`，如：$x \supseteq y$
- 非真子集运算，符号：`\supsetneq`，如：$x \supsetneq y$
- 非子集运算，符号：`\not\subset`，如：$x \not\subset y$
- 非子集运算，符号：`\not\supset`，如：$x \not\supset y$
- 并集运算，符号：`\cup`，如：$x \cup y$
- 交集运算，符号：`\cap`，如：$x \cap y$
- 差集运算，符号：`\setminus`，如：$x \setminus y$
- 同或运算，符号：`\bigodot`，如：$x \bigodot y$
- 同与运算，符号：`\bigotimes`，如：$x \bigotimes y$
- 实数集合，符号：`\mathbb{R}`，如：`\mathbb{R}`
- 自然数集合，符号：`\mathbb{Z}`，如：`\mathbb{Z}`
- 空集，符号：`\emptyset`，如：$\emptyset$

数学符号
----

- 无穷，符号：`\infty`，如：$\infty$
- 虚数，符号：`\imath`，如：$\imath$
- 虚数，符号：`\jmath`，如：$\jmath$
- 数学符号，符号`\hat{a}`，如：$\hat{a}$
- 数学符号，符号`\check{a}`，如：$\check{a}$
- 数学符号，符号`\breve{a}`，如：$\breve{a}$
- 数学符号，符号`\tilde{a}`，如：$\tilde{a}$
- 数学符号，符号`\bar{a}`，如：$\bar{a}$
- 矢量符号，符号`\vec{a}`，如：$\vec{a}$
- 数学符号，符号`\acute{a}`，如：$\acute{a}$
- 数学符号，符号`\grave{a}`，如：$\grave{a}$
- 数学符号，符号`\mathring{a}`，如：$\mathring{a}$
- 一阶导数符号，符号`\dot{a}`，如：$\dot{a}$
- 二阶导数符号，符号`\ddot{a}`，如：$\ddot{a}$
- 上箭头，符号：`\uparrow`，如：$\uparrow$
- 上箭头，符号：`\Uparrow`，如：$\Uparrow$
- 下箭头，符号：`\downarrow`，如：$\downarrow$
- 下箭头，符号：`\Downarrow`，如：$\Downarrow$
- 左箭头，符号：`\leftarrow`，如：$\leftarrow$
- 左箭头，符号：`\Leftarrow`，如：$\Leftarrow$
- 右箭头，符号：`\rightarrow`，如：$\rightarrow$
- 右箭头，符号：`\Rightarrow`，如：$\Rightarrow$
- 底端对齐的省略号，符号：`\ldots`，如：$1,2,\ldots,n$
- 中线对齐的省略号，符号：`\cdots`，如：$x_1^2 + x_2^2 + \cdots + x_n^2$
- 竖直对齐的省略号，符号：`\vdots`，如：$\vdots$
- 斜对齐的省略号，符号：`\ddots`，如：$\ddots$

希腊字母
----

| 字母 | 实现 | 字母 | 实现 |
| --- | --- | --- | --- |
| A | `A` | α | `\alhpa` |
| B | `B` | β | `\beta` |
| Γ | `\Gamma` | γ | `\gamma` |
| Δ | `\Delta` | δ | `\delta` |
| E | `E` | ϵ | `\epsilon` |
| Z | `Z` | ζ | `\zeta` |
| H | `H` | η | `\eta` |
| Θ | `\Theta` | θ | `\theta` |
| I | `I` | ι | `\iota` |
| K | `K` | κ | `\kappa` |
| Λ | `\Lambda` | λ | `\lambda` |
| M | `M` | μ | `\mu` |
| N | `N` | ν | `\nu` |
| Ξ | `\Xi` | ξ | `\xi` |
| O | `O` | ο | `\omicron` |
| Π | `\Pi` | π | `\pi` |
| P | `P` | ρ | `\rho` |
| Σ | `\Sigma` | σ | `\sigma` |
| T | `T` | τ | `\tau` |
| Υ | `\Upsilon` | υ | `\upsilon` |
| Φ | `\Phi` | ϕ | `\phi` |
| X | `X` | χ | `\chi` |
| Ψ | `\Psi` | ψ | `\psi` |
| Ω | `\v` | ω | `\omega` |

