---
title: I/O
linkTitle: I/O
date: 2023-07-04T10:36
weight: 20
---

# 概述

> 参考：
>
> - [Go 标准库，io](https://pkg.go.dev/io)
> - [Go 标准库，bufio](https://pkg.go.dev/bufio)
> - [公众号 - 云原生生态圈，Go 写文件的几种姿势，你喜欢哪一种？](https://mp.weixin.qq.com/s/56g5k17Zt4iytbWkYcouig)
> - [Introduction to bufio package in Golang](https://medium.com/golangspec/introduction-to-bufio-package-in-golang-ad7d1877f762)
>   - [Go 语言中文网，Go 语言 bufio 包介绍](https://studygolang.com/articles/11824)

> [!Notes]
> 想要理解 io 标准库的逻辑，必须要理解 [Method AND Interface](docs/2.编程/高级编程语言/Go/Go%20规范与标准库/Method%20AND%20Interface/Method%20AND%20Interface.md) 中的 Interface(接口) 的概念，这是理解 I/O 的前提，否则不要往下阅读！

[io.Reader](https://pkg.go.dev/io#Reader)、[io.Writer](https://pkg.go.dev/io#Writer) 是 io 包中的接口，用于处理 I/O 操作。

- 所有实现了 io.Reader 接口的类型都可以作为输入源，例如 文件、网络连接、etc.
- 所有实现了 io.Writer 接口的类型都可以作为输出目标，例如 文件、网络连接、etc.

**[bufio](https://pkg.go.dev/bufio)** 包用来帮助处理 **[buffered I/O(I/O 缓存)](https://www.quora.com/In-C-what-does-buffering-I-O-or-buffered-I-O-mean/answer/Robert-Love-1)**，通过 I/O 缓存我们可以减少对系统调用，提高性能。

## 读取用户的输入 

大多数的程序都是处理输入，产生输出；这也正是`计算`的定义。但是程序如何获取要处理的输入数据呢？有一些程序生成自己的数据，但是通常情况下，输入来自于程序外部，e.g.文件、网络连接、其他程序的输出、敲键盘的用户、命令行参数或其它类似的输入源。想要使用 Go 语言的输入输出功能，一般不外乎下面 3 步

1. 获取输入源的定位符，e.g.文件描述符、用户的标准输入、etc.
2. 通过输入源的定位符，把输入内容放到缓冲区并充缓冲器发送给变量
3. 打印缓冲区的变量即可实现输出输入源提供的数据

`fmt` 包的 `Scan` 和 `Sscan` 开头的函数。
  
`Scanln` 扫描来自标准输入的文本，将空格分隔的值一次存放到后续的参数内，直到碰到换行。`Scanf` 与 `Scanln` 类似，除了 `Scanf` 的第一个参数作用格式字符串，用来决定如何读取。以 `Sscan` 和以 `Sscan` 开头的函数则是从字符串读取，除此之外，与 `Scanf` 相同。

## 文件读写

在 Go 语言中，文件使用指向 os.File 类型的指针来表示的，也叫做文件描述符。所以想要通过代码读取一个文件的内容则需要先获取文件描述符。

## 拷贝文件

https://golang.org/pkg/io/#Copy

# 最佳实践

## 读取文件时显示进度

```go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func progress(ch <-chan int64) {
	for rate := range ch {
		fmt.Printf("\rrate:%3d%%", rate)
	}
}

//读取文件
func readFile() {
	rateCh := make(chan int64)
	defer close(rateCh)

	file, err := os.Open("./test_file/read_file_show_progess.jpg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileTotalSize := fileInfo.Size()

	go progress(rateCh)
	resultFileByte := make([]byte, 0)

	for {
		readBuf := make([]byte, 1024)
		n, err := file.Read(readBuf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		resultFileByte = append(resultFileByte, readBuf...)
		time.Sleep(1 * time.Millisecond)

		go func() {
			rateCh <- int64(float32(len(resultFileByte)) / float32(fileTotalSize) * 100)
		}()
	}

	ioutil.WriteFile("./test_file/read_file_show_progess_result.jpg", resultFileByte, 0600)
}

func main() {
	readFile()
}

```
