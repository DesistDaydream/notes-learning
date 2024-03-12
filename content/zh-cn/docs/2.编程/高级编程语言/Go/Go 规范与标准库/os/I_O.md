---
title: I/O
linkTitle: I/O
date: 2023-07-04T10:36
weight: 2
---

# 概述

> 参考：

Go 写文件的几种姿势，你喜欢哪一种？<https://mp.weixin.qq.com/s/56g5k17Zt4iytbWkYcouig>

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
