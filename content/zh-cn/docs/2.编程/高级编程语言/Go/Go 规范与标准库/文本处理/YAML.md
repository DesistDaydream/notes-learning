---
title: YAML
linkTitle: YAML
date: 2024-10-19T19:17
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，go-yaml/yaml](https://github.com/go-yaml/yaml)

[YAML](docs/2.编程/无法分类的语言/YAML.md) 数据解析

YAML解析库 沿用了 JSON解析库 的相关说法。

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Nginx nginx  配置
type Nginx struct {
	Port    int    `yaml:"Port"`
	LogPath string `yaml:"LogPath"`
	Path    string `yaml:"Path"`
}

//Config   系统配置配置
type Config struct {
	Name      string `yaml:"SiteName"`
	Addr      string `yaml:"SiteAddr"`
	HTTPS     bool   `yaml:"Https"`
	SiteNginx Nginx  `yaml:"Nginx"`
}

func main() {
	var setting Config
	config, errRead := ioutil.ReadFile("./info.yaml")
	if errRead != nil {
		fmt.Print(errRead)
	}
	errUnmarshal := yaml.Unmarshal(config, &setting)
	if errUnmarshal != nil {
		log.Fatalf("error: %v", errUnmarshal)
	}

	fmt.Println(setting)
	fmt.Println(setting.Name)
	fmt.Println(setting.Addr)
	fmt.Println(setting.HTTPS)
	fmt.Println(setting.SiteNginx.Port)
	fmt.Println(setting.SiteNginx.LogPath)
	fmt.Println(setting.SiteNginx.Path)
}
```
