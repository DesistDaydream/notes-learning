---
title: Git常见问题
---

# warning LF will be replaced by CRLF

## 发现问题

windows 平台下使用 git add，git deploy 文件时经常出现“warning: LF will be replaced by CRLF” 的提示。

网上很多解决办法提到：

设置 core.autocrlf=false，windows 也用 LF 换行。

除了记事本，其他编辑器都可以正常编辑。

而没有给出具体原因和分析，现在加以补充。

## 分析问题

格式化与多余的空白字符，特别是在跨平台情况下，有时候是一个令人发指的问题。由于编辑器的不同或者文件行尾的换行符在 Windows 下被替换了，一些细微的空格变化会不经意地混入提交，造成麻烦。虽然这是小问题，但它会极大地扰乱跨平台协作。

其实，这是因为在文本处理中，CR（CarriageReturn），LF（LineFeed），CR/LF 是不同操作系统上使用的换行符，具体如下：

换行符‘\n’和回车符‘\r’

- 回车符就是回到一行的开头，用符号 r 表示，十进制 ASCII 代码是 13，十六进制代码为 0x0D，回车（return）；

- 换行符就是另起一行，用 n 符号表示，ASCII 代码是 10，十六制为 0x0A， 换行（newline）。

所以我们平时编写文件的回车符应该确切来说叫做回车换行符。

## 应用情况

- Dos 和 Windows 平台： 使用回车（CR）和换行（LF）两个字符来结束一行，回车+换行(CR+LF)，即“\r\n”；

- Mac 和 Linux 平台：只使用换行（LF）一个字符来结束一行，即“\n”；

- 最早 Mac 每行结尾是回车 CR 即'\r'，后 mac os x 也投奔了 unix。

许多 Windows 上的编辑器会悄悄把行尾的换行（LF）字符转换成回车（CR）和换行（LF），或在用户按下 Enter 键时，插入回车（CR）和换行（LF）两个字符。

## 影响

- 一个直接后果是，Unix/Mac 系统下的文件在 Windows 里打开的话，所有文字会变成一行；

- 而 Windows 里的文件在 Unix/Mac 下打开的话，在每行的结尾可能会多出一个^M 符号。

- Linux 保存的文件在 windows 上用记事本看的话会出现黑点。

这些问题都可以通过一定方式进行转换统一，例如，在 linux 下，命令 unix2dos 是把 linux 文件格式转换成 windows 文件格式，命令 dos2unix 是把 windows 格式转换成 linux 文件格式。

## 解决问题

Git 默认行为会在提交时自动地把 CRLF(回车换行) 转换成 LF(换行)，而在检出代码时把 LF(换行) 转换成 CRLF(回车换行) .

情况一：

你可以用 git config --global core.autocrlf true 来打开此项功能。 如果是在 Windows 系统上，把它设置成 true，这样在检出代码时，换行会被转换成回车和换行：

```bash
# 提交时转换为LF，检出时转换为CRLF
$ git config --global core.autocrlf true
```

情况二：

如果使用以换行（LF）作为行结束符的 Linux 或 Mac，你不需要 Git 在检出文件时进行自动的转换。然而当一个以回车（CR）和换行（LF）作为行结束符的文件不小心被引入时，你肯定想让 Git 修正。 所以，你可以把 core.autocrlf 设置成 input 来告诉 Git 在提交时把回车和换行转换成换行，检出时不转换：（这样在 Windows 上的检出文件中会保留回车和换行，而在 Mac 和 Linux 上，以及版本库中会保留换行。）

```bash
# 提交时转换为LF，检出时不转换
$ git config --global core.autocrlf input
```

情况三：

如果你是 Windows 程序员，且正在开发仅运行在 Windows 上的项目，可以设置 false 取消此功能，把回车保留在版本库中：

```bash
# 提交检出均不转换
$ git config --global core.autocrlf false
```

你也可以在文件提交时进行 safecrlf 检查

```bash
# 拒绝提交包含混合换行符的文件
git config --global core.safecrlf true
# 允许提交包含混合换行符的文件
git config --global core.safecrlf false
# 提交包含混合换行符的文件时给出警告
git config --global core.safecrlf warn
```

## 通俗解释

- git 的 Windows 客户端基本都会默认设置 core.autocrlf=true，设置 core.autocrlf=true, 只要保持工作区都是纯 CRLF 文件，编辑器用 CRLF 换行，就不会出现警告了；
- Linux 最好不要设置 core.autocrlf，因为这个配置算是为 Windows 平台定制；
- Windows 上设置 core.autocrlf=false，仓库里也没有配置 .gitattributes，很容易引入 CRLF 或者混合换行符（Mixed Line Endings，一个文件里既有 LF 又有 CRLF）到版本库，这样就可能产生各种奇怪的问题。
