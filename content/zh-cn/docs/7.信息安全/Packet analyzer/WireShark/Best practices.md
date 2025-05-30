---
title: Best practices
linkTitle: Best practices
weight: 21
---

# 概述

> 参考：
>
> - 

# 从 pcap 文件中还原出原始文件

一个文件本身是一堆 Bytes 的组合，而 pcap 中每个 tcp 之类的包的 payload 包含的就是文件传输时文件的所有 Bytes，只不过被拆分到多个包中

首先使用 `_ws.col.info contains "TCP segment of a reassembled PDU"` 过滤出来所有被分片的包（这些通常都是包含文件的包）

选中其中一个包，按 `ctrl + alt + shift + t` 快捷键，追踪该包所述的 TCP 流

一个文件本质上就是红框中的内容（只不过这里把 Bytes 以 ASCII 表示，所以人类不可读）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/best_practices/restore_files_tcp_stream_head.png)

在窗口下方选择**该会话其中一个方向**，并选择**原始数据**后，点击另存为，就保存下来文件了。假如命名为 response_file.hex

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/best_practices/restore_files_2.png)

此时我们需要做的就是想办法把该文件开头/结尾不需要的部分去掉 (●ˇ∀ˇ●)

使用 [Visual Studio Code](docs/2.编程/Programming%20environment/IDE/Visual%20Studio%20Code/Visual%20Studio%20Code.md) 安装一个十六进制编辑器的插件，并用打开 response_file.hex 会看到如下内容。

![800](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/best_practices/restore_files_response_file_hex_1.png)

我们需要把 HTTP 的 Header 删掉，根据 WireShark 的 TCP 流看，HTTP Header 后面跟着两个换行（从 [ASCII 表](docs/8.通用技术/编码与解码/字符的编码与解码/ASCII%20表.md) 查时 0D 0A），所以也就是截止到 50 4B 之前的所有内容都要删掉

然后回到 WireShark 中看看 TCP 流的尾部，是一个换行及后面的内容，也要删掉

![1000](https://notes-learning.oss-cn-beijing.aliyuncs.com/wireshark/best_practices/restore_files_4.png)

同样的，在 VSCode 中删掉，最后保存剩余的内容为新文件，此时，该文件就是原始文件了。MD5 也是一样的。

