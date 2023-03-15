---
title: printf 与 echo 文本输出命令
---

#

printf, fprintf, sprintf, snprintf, vprintf, vfprintf, vsprintf, vsnprintf - 格式化输出转换

# printf 命令：格式化并且打印数据 format and print data

主要用于对 ARGUMENTs 进行格式化输出，ARGUMENTs 可以是字符串、数值等等，甚至可以通过变量引用；FORMAT 主要是对 ARGUMENTs 里的各种数据进行格式化输出，e.g.每个 ARGUMENT 是什么类型的(字符、整数、2 进制、16 进制等等)，各个 ARGUMENT 中间使用什么分隔符、是否换行等等。

## 语法格式：printf FORMAT \[ARGUMENT...]

主要用于按照 FORMAT 定义的格式来输出 ARGUMENT...给出的内容

### FORMAT 包括：格式替代符，自定义内容，格式控制符，这 3 个在使用的时候没有先后顺序

格式替代符 # 用于控制输出的每个 Argument 的类型。一个“格式替代符”对应后面一个 Argument，如果想要输出的类型与 Argument 给定的类型不符，则进行类型转换后输出 e.g.Argument 给了一个整数 100，而格式替代符使用的是%X,则会输出 64；若 Argument 不够 FORMAT 的个数，则以空白补充。一般情况格式替代符使用双引号引起来

1. %b #相对应的参数被视为含有要被处理的转义序列之字符串。

2. %c #ASCII 字符。显示相对应参数的第一个字符

3. %d, %i #十进制整数

4. %e, %E, %f #浮点格式

5. %g #%e 或%f 转换，看哪一个较短，则删除结尾的零

6. %G #%E 或%f 转换，看哪一个较短，则删除结尾的零

7. %o #不带正负号的八进制值

8. %s #字符串

9. %u #不带正负号的十进制值

10. %x #不带正负号的十六进制值，使用 a 至 f 表示 10 至 15

11. %X #不带正负号的十六进制值，使用 A 至 F 表示 10 至 15

12. %% #字面意义的%

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371965-2e43b313-e44b-4fbe-a44b-254e0e0b37fb.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371962-7d50c11b-356a-4098-acab-15272005b08a.jpeg)

自定义内容 # 在输出内容的格式中，可以自己添加任务字符串

格式控制符 # 用于控制输出内容整体的格式。i.e.每个参数有之间有多少空白符，在哪里换行等等

1. \a #警告字符，通常为 ASCII 的 BEL 字符

2. \b #后退

3. \c #抑制（不显示）输出结果中任何结尾的换行字符（只在%b 格式指示符控制下的参数字符串中有效），而且，任何留在参数里的字符、任何接下来的参数以及任何留在格式字符串中的字符，都被忽略

4. \f #换页（formfeed)

5. \n #换行

6. \r #回车（Carriage return）

7. \t #水平制表符 i.e.tab

8. \v #垂直制表符

9. \ # 一个字面上的反斜杠字符

10. \ddd #表示 1 到 3 位数八进制值的字符。仅在格式字符串中有效

11. \0ddd #表示 1 到 3 位的八进制值字符

### ARGUMENTs 包括各种想要输出的具体内容，可以是字符串、整数等、甚至可以引用变量

每个 Argument 使用空格进行分割，一个 Argument 中的内容传递给 FORMAT 中的“格式替代符”

EXAMPLE

```bash
[root@node3 ~]# printf "%-10s %-8s %-4s\n" 姓名 性别 体重kg; printf "%-10s %-8s %-4.2f\n" 郭靖 男 66.1234;printf "%-10s %-8s %-4.2f\n" 杨过 男 48.6543;printf "%-10s %-8s %-4.2f\n" 郭芙 女 47.9876
姓名     性别   体重kg
郭靖     男      66.12
杨过     男      48.65
郭芙     女      47.99
```

%s %c %d %f 都是格式替代符

%-10s 指一个宽度为 10 个字符（-表示左对齐，没有则表示右对齐），任何字符都会被显示在 10 个字符宽的字符内，如果不足则自动以空格填充，超过也会将内容全部显示出来。

%-4.2f 指格式化为小数，其中.2 指保留 2 位小数。

\n 表示“换行符”，i.e.输出完这一段就换行，否则每个 printf 输出的内容都在一行了

hi 是自定义内容，可以随便写

# Echo # 显示文本

Bash 脚本是非常流行的最简单的脚本语言。 与任何编程或脚本语言一样，您会在终端上遇到打印文本。 这可能发生在许多场景中，例如当您想要输出文件的内容或检查变量的值时。 程序员还可以通过在控制台上打印变量值来调试应用程序。 因此，在我们深入研究另一个教程的 bash 脚本之前，让我们看一下在终端中输出文本的不同方式。

为了在终端上输出文本，Echo 是您需要知道的最重要的命令。 正如名称本身所示，echo 在终端的标准输出上显示数字或字符串。 它还有许多选项，如下表所示。

根据 Linux 文档，以下是 echo 命令的语法。

echo \[OPTIONS] \[ARG....]

OPTIONS

1. -n # 不打印后面的换行符

2. -E # 禁用反斜杠转义字符的解释，默认自带选项。

3. -e # 启用反斜杠转义的解释

输出 args，以空格分隔，后跟换行符。返回状态始终为 0。 -E 选项禁用这些转义符的解释，即使在默认情况下已解释它们的系统上也是如此。 xpg_echo shell 选项可用于动态确定 echo 默认情况下是否扩展这些转义字符。回声并不意味着选项的结束。 echo 解释以下转义序列

如果给出-e 选项，则启用对以下反斜杠转义字符的解释：

1. \a 显示警告字符

2. \b 退格符。

3. \c 在输出中禁止另外跟在最终参数后面的换行字符。所有跟在 \c 序列后的字符都被忽略。

4. \e Escape

5. \E

6. \f 换页

7. \n 新行

8. \r 回车

9. \t 水平选项卡

10. \v 垂直制表符

11. \ backslash

发送文本到标准输出

要输出终端上的任何字符串、数字或文本，请键入以下命令并按 enter。

echo "Hello World Linux 公社www.linuxidc.com"

以下输出将显示在终端上

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371980-9a7c412a-4bdc-4673-a511-ed9283a06d8b.jpeg)

打印一个变量

让我们声明一个变量并在终端上打印它的值。假设 x 是我们在 160 处初始化的一个变量。

x=160

现在，我们将在终端上输出变量的值。

echo $x

终端将打印 160。同样，您也可以将字符串存储在变量中并将其输出到终端。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371967-852dae9b-770e-40c9-b132-2ca2788f3e60.jpeg)

试一试，让我们知道这对你来说是否容易。

删除单词之间的空格

这是我最喜欢的 echo 选项之一，因为它消除了句子中不同单词之间的所有空格，并将它们混在一起。在这个特性中，我们将使用表 1 中提到的两个选项。

echo -e "欢迎来到 \bLinux \b 公社 \bwww \b.linuxidc \b.com"

显示：欢迎来到 Linux 公社www.linuxidc.com

从上面的示例中可以看到，我们正在启用反斜杠转义的解释以及添加退格。输出如下所示。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166372007-65bfcc66-1d0b-4927-b69e-7f065e1cb924.jpeg)

以新行输出单词

在使用 bash 脚本时，echo 的这个选项非常方便。大多数情况下，你需要在完成后移动到下一行。因此，这是最好的选择。

echo -e "欢迎来到 \nLinux \n 公社 \nwww \n.linuxidc \n.com"

输出将在单独的一行中显示每个单词，如下面的屏幕截图所示。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371988-02f5d8d5-7ae2-4475-8cdd-659b85c73f30.jpeg)

输出带声音的文本

这是一个使用 bell 或 alert 输出文本的简单选项。为此，键入以下命令。

echo -e "hello \a 欢迎来到 Linux 公社www.linuxidc.com"

确保系统的音量足够大，以便在终端输出文本时能够听到微小的铃声。

删除后面新行

echo 的另一个选项是删除后面的换行符，以便在同一行输出所有内容。为此，我们使用“\c”选项，如下图所示。

echo -e "欢迎来到 Linux 公社www.linuxidc.com \c 你是谁你是谁"

显示以下输出

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166371985-a27b4bd3-aca2-4074-992f-721954070d00.jpeg)

将回车符添加到输出中

要在输出中添加特定的回车符，我们有“\r”选项。

echo -e "Linux 公社www.linuxidc.com \r 欢迎您的来到"

在终端上显示以下输出。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166372027-a7366654-9e20-413e-bfae-f13e15fc6616.jpeg)

在输出中使用选项卡

在终端上打印输出时，您也可以添加水平和垂直标签。 这些产品可以用于更清洁的产品。 要添加水平制表符，您必须添加“\t”，对于垂直制表符，请添加“\v”。 我们将为这些中的每一个做一个样本，然后组合一个。

echo -e "www.linuxidc.com \t 欢迎来到 Linux 公社"

这个命令的输出如下所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166372029-9dcd1465-f711-418d-9414-46e2511d9b37.jpeg)

echo -e "www.linuxidc.com \v 欢迎来到 Linux 公社"

这个命令的输出如下所示

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/nfqmgm/1616166372005-1aca898b-377c-4bd4-82d2-70fd41a3985a.jpeg)

这就是在终端上打印文本的所有选项。这是一个需要学习的重要特性，因为当您开始使用 bash 脚本时，它将进一步帮助您。确保你实现了每一个选项并努力练习。如果本教程帮助您解决了问题，请告诉我们。
