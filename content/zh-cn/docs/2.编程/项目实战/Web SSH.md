---
title: "Web SSH"
linkTitle: "Web SSH"
date: "2023-07-20T17:57"
weight: 20
---

# 概述

> 参考：
> 
> - [公众号-云原生运维圈，Gin+Xterm.js实现WebSSH远程Kubernetes Pod](https://mp.weixin.qq.com/s/GKK6JHohFNRVPke7zn1apQ)

**一**

**Xterm.js简介**

![](https://mmbiz.qpic.cn/sz_mmbiz_png/vNHkZdwXUGPicU8ALTCgMicSnVTvm1Jjf0kYaWIeHjdfQB1qRES4icz0WibNgKq64IB2nxERlpAU28OQXdSpjU6ictw/640?wx_fmt=png)

xterm.js （https://xtermjs.org/）是一个开源的 JavaScript 库，它模拟了一个终端接口，可以在网页中嵌入一个完全功能的终端。这个库非常灵活，并且具有很多定制选项和插件系统。

下面是一些使用 xterm.js 的基本步骤：

*   首先，需要在项目中安装 xterm.js。你可以直接从 npm 安装：  

npm install xterm  

*   然后在 HTML 中创建一个容器来承载终端

`<div id="terminal"></div>  
`
*   在你的 JavaScript 文件中，导入 Terminal 类并创建一个新的实例

`import { Terminal } from 'xterm';

const term = new Terminal();

*   把这个终端附加到 HTML 元素上

`term.open(document.getElementById('terminal'));  
`

*   现在你就可以向终端写入数据了

`term.write('Hello, World!');  

*   如果你想读取用户在终端中的输入，可以监听 onData 事件

`term.onData(data => {  
  console.log(data);  
});  

以上只是最基础的使用方法。xterm.js 提供了许多其他功能，如主题定制、附加插件（例如 FitAddon 可以自动调整终端大小，WebLinksAddon 可以捕获 URL 并将其变为可点击链接）、设置光标样式、更改字体大小等等。你可以访问 xterm.js 的 GitHub （https://github.com/xtermjs/xterm.js）仓库 或者 文档 来获取更详细的信息。

**二**

**使用Gin、client-go的SPDYExecutor来执行远程命令**

``package main

import (  
 "context"  
 "encoding/json"  
 "github.com/gin-gonic/gin"  
 "github.com/gorilla/websocket"  
 corev1 "k8s.io/api/core/v1"  
 "k8s.io/client-go/kubernetes"  
 "k8s.io/client-go/kubernetes/scheme"  
 "k8s.io/client-go/tools/clientcmd"  
 "k8s.io/client-go/tools/remotecommand"  
 "log"  
 "net/http"  
)

// websocket 升级器配置  
var upgrader = websocket.Upgrader{  
 CheckOrigin: func(r *http.Request) bool {  
  return true  
 },  
}

// WSClient 结构体，封装了 WebSocket 连接和 resize 通道，用于在 WebSocket 和 remotecommand 之间进行数据交换。  
type WSClient struct {  
 // WebSocket 连接对象  
 ws *websocket.Conn  
 // TerminalSize 类型的通道，用于传输窗口大小调整事件  
 resize chan remotecommand.TerminalSize  
}

// MSG 结构体，用于解析从 WebSocket 接收到的消息。  
type MSG struct {  
 // 消息类型字段  
 MsgType string `json:"msg_type"`  
 Rows uint16 `json:"rows"`  
 Cols uint16 `json:"cols"`  
 // 输入消息的数据字段  
 Data string `json:"data"`  
}

// WSClient 的 Read 方法，实现了 io.Reader 接口，从 websocket 中读取数据。  
func (c *WSClient) Read(p []byte) (n int, err error) {  
 // 从 WebSocket 中读取消息  
 _, message, err := c.ws.ReadMessage()  
 if err != nil {  
  return 0, err  
 }  
 var msg MSG  
 if err := json.Unmarshal(message, &msg); err != nil {  
  return 0, err  
 }

 // 根据消息类型进行不同的处理  
 switch msg.MsgType {  
 // 如果是窗口调整消息  
 case "resize":  
  winSize := remotecommand.TerminalSize{  
   Width: msg.Cols,  
   Height: msg.Rows,  
  }  
  // 将 TerminalSize 对象发送到 resize 通道  
  c.resize <- winSize  
  return 0, nil  
 // 如果是输入消息  
 case "input":  
  copy(p, msg.Data)  
  return len(msg.Data), err  
 }  
 return 0, nil  
}

// WSClient 的 Write 方法，实现了 io.Writer 接口，将数据写入 websocket。  
func (c *WSClient) Write(p []byte) (n int, err error) {  
 // 将数据作为文本消息写入 WebSocket  
 err = c.ws.WriteMessage(websocket.TextMessage, p)  
 return len(p), err  
}

// Next WSClient 的 Next 方法，用于从 resize 通道获取下一个 TerminalSize 事件。  
func (c *WSClient) Next() *remotecommand.TerminalSize {  
 // 从 resize 通道读取 TerminalSize 对象  
 size := <-c.resize  
 return &size  
}

// podSSH 函数，这是主要的 SSH 功能逻辑，使用 kubernetes client-go 的 SPDY executor 来执行远程命令。  
func podSSH(wsClient *WSClient, q query) {  
 // 使用 kubeconfig 文件初始化 kubernetes 客户端配置  
 // 请注意，你应该替换 ./config 为你的 kubeconfig 文件路径  
 restClientConfig, err := clientcmd.BuildConfigFromFlags("", "./config")  
 if err != nil {  
  log.Fatalf("Failed to build config: %v", err)  
 }

 // 根据配置创建 kubernetes 客户端  
 clientSet, err := kubernetes.NewForConfig(restClientConfig)  
 if err != nil {  
  log.Fatalf("Failed to create client: %v", err)  
 }  
 // 构造一个用于执行远程命令的请求  
 request := clientSet.CoreV1().RESTClient().Post().  
  Resource("pods").  
  Namespace(q.Namespace).  
  Name(q.PodName).  
  SubResource("exec").  
  VersionedParams(&corev1.PodExecOptions{  
   Container: q.ContainerName,  
   Command: []string{  
    q.Command,  
   },  
   Stdout: true,  
   Stdin:  true,  
   Stderr: true,  
   TTY:    true,  
  }, scheme.ParameterCodec)  
 // 创建 SPDY executor，用于后续的 Stream 操作  
 exec, err := remotecommand.NewSPDYExecutor(restClientConfig, "POST", request.URL())  
 if err != nil {  
  log.Fatalf("Failed to initialize executor: %v", err)  
 }

 // 开始进行 Stream 操作，即通过 websocket 执行命令  
 err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{  
  Stderr:            wsClient,  
  Stdout:            wsClient,  
  Stdin:             wsClient,  
  Tty:               true,  
  TerminalSizeQueue: wsClient,  
 })  
 if err != nil {  
  log.Fatalf("Failed to start stream: %v", err)  
 }  
}

// query 结构体，用于解析和验证查询参数  
type query struct {  
 Namespace     string `form:"namespace" binding:"required"`  
 PodName       string `form:"pod_name" binding:"required"`  
 ContainerName string `form:"container_name" binding:"required"`  
 Command       string `form:"command" binding:"required"`  
}

func main() {  
 router := gin.Default()  
 router.LoadHTMLGlob("templates/*")  
 router.GET("/", func(ctx *gin.Context) {  
  ctx.HTML(http.StatusOK, "ssh.html", nil)  
 })

 // 设置 /ssh 路由  
 router.GET("/ssh", func(ctx *gin.Context) {  
  var r query  
  if err := ctx.ShouldBindQuery(&r); err != nil {  
   ctx.JSON(http.StatusBadRequest, gin.H{  
    "err": err.Error(),  
   })  
   return  
  }  
  // 将 HTTP 连接升级为 websocket 连接  
  ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)  
  if err != nil {  
   log.Printf("Failed to upgrade connection: %v", err)  
   return  
  }  
  // 使用 podSSH 函数处理 websocket 连接  
  podSSH(&WSClient{  
   ws:     ws,  
   resize: make(chan remotecommand.TerminalSize),  
  }, r)  
 })  
 router.Run(":9191")  
}

``

*   后端全部代码  
    

https://gitee.com/KubeSec/pod-webssh/tree/master/pod-ssh

**三**

**使用vue-admin-template和Xterm.js实现Web终端**

https://github.com/PanJiaChen/vue-admin-template

https://github.com/xtermjs/xterm.js

*   下载vue-admin-template项目
    

https://github.com/PanJiaChen/vue-admin-template.git

*   安装xterm.js及插件
    

`npm install  
npm install xterm  
npm install --save xterm-addon-web-links  
npm install --save xterm-addon-fit  
npm install -S xterm-style`

*   打开vue-admin-template项目，在src/views目录下新建目录pod-ssh，在pod-ssh目录下新建index.vue代码如下  

```
<template>  
  <div class="app-container">  
    <!-- 使用 Element UI 的表单组件创建一个带有标签和输入框的表单 -->  
    <el-form ref="form" :model="form" :inline="true" label-width="120px">  
      <el-form-item label="namespace"> <!-- namespace 输入框 -->  
        <el-input v-model="form.namespace" />  
      </el-form-item>  
      <el-form-item label="pod name"> <!-- pod 名称输入框 -->  
        <el-input v-model="form.pod_name" />  
      </el-form-item>  
      <el-form-item label="container name"> <!-- 容器名称输入框 -->  
        <el-input v-model="form.container_name" />  
      </el-form-item>  
      <el-form-item label="Command"> <!-- 命令选择框 -->  
        <el-select v-model="form.command" placeholder="bash">  
          <el-option label="bash" value="bash" />  
          <el-option label="sh" value="sh" />  
        </el-select>  
      </el-form-item>  
      <el-form-item> <!-- 提交按钮 -->  
        <el-button type="primary" @click="onSubmit">Create</el-button>  
      </el-form-item>  
      <div id="terminal" /> <!-- 终端视图容器 -->  
    </el-form>  
  </div>  
</template>  
  
<script>  
import { Terminal } from 'xterm' // 导入 xterm 包，用于创建和操作终端对象  
import { common as xtermTheme } from 'xterm-style' // 导入 xterm 样式主题  
import 'xterm/css/xterm.css' // 导入 xterm CSS 样式  
import { FitAddon } from 'xterm-addon-fit' // 导入 xterm fit 插件，用于调整终端大小  
import { WebLinksAddon } from 'xterm-addon-web-links' // 导入 xterm web-links 插件，可以捕获 URL 并将其转换为可点击链接  
import 'xterm/lib/xterm.js' // 导入 xterm 库  
  
export default {  
  data() {  
    return {  
      form: {  
        namespace: 'default', // 默认命名空间为 "default"  
        command: 'bash', // 默认 shell 命令为 "bash"  
        pod_name: 'nginx', // 默认 Pod 名称为 "nginx"  
        container_name: 'nginx' // 默认容器名称为 "nginx"  
      },  
    }  
  },  
  methods: {  
    onSubmit() {  
      // 创建一个新的 Terminal 对象  
      const xterm = new Terminal({  
        theme: xtermTheme,  
        rendererType: 'canvas',  
        convertEol: true,  
        cursorBlink: true  
      })  
  
      // 创建并加载 FitAddon 和 WebLinksAddon  
      const fitAddon = new FitAddon()  
      xterm.loadAddon(fitAddon)  
      xterm.loadAddon(new WebLinksAddon())  
  
      // 打开这个终端，并附加到 HTML 元素上  
      xterm.open(document.getElementById('terminal'))  
  
      // 调整终端的大小以适应其父元素  
      fitAddon.fit()  
  
      // 创建一个新的 WebSocket 连接，并通过 URL 参数传递 pod, namespace, container 和 command 信息  
      const ws = new WebSocket('ws://127.0.0.1:9191/ssh?namespace=' + this.form.namespace + '&pod_name=' + this.form.pod_name + '&container_name=' + this.form.container_name + '&command=' + this.form.command)  
  
      // 当 WebSocket 连接打开时，发送一个 resize 消息给服务器，告诉它终端的尺寸  
      ws.onopen = function() {  
        ws.send(JSON.stringify({  
          msg_type: 'resize',  
          rows: xterm.rows,  
          cols: xterm.cols  
        }))  
      }  
  
      // 当从服务器收到消息时，写入终端显示  
      ws.onmessage = function(evt) {  
        xterm.write(evt.data)  
      }  
  
      // 当发生错误时，也写入终端显示  
      ws.onerror = function(evt) {  
        xterm.write(evt.data)  
      }  
  
      // 当窗口尺寸变化时，重新调整终端的尺寸，并发送一个新的 resize 消息给服务器  
      window.addEventListener('resize', function() {  
        fitAddon.fit()  
        ws.send(JSON.stringify({  
          msg_type: 'resize',  
          rows: xterm.rows,  
          cols: xterm.cols  
        }))  
      })  
  
      // 当在终端中键入字符时，发送一个 input 消息给服务器  
      xterm.onData((b) => {  
        ws.send(JSON.stringify({  
          msg_type: 'input',  
          data: b  
        }))  
      })  
    }  
  }  
}  
</script>  
  
<style scoped>  
.line{  
  text-align: center;  
}  
</style>
```

在src/router/index.js文件中增加路由  

```
{  
    path: '/pod-ssh',  
    component: Layout,  
    children: [  
      {  
        path: 'pod-ssh',  
        name: 'SSH',  
        component: () => import('@/views/pod-ssh/index'),  
        meta: { title: 'SSH', icon: 'form' }  
      }  
    ]  
  },
```

*   启动项目

`npm install  
npm install  
`

*   前端全部代码  

https://gitee.com/KubeSec/pod-webssh/tree/master/pod-webssh

**四**

测试

*   在kubernetes中创建测试的Pod  
    

`apiVersion: v1  
kind: Pod  
metadata:  
  name: nginx  
spec:  
  containers:  
  - name: nginx  
    image: nginx:1.14.2  
    ports:  
    - containerPort: 80`

访问http://localhost:9528/#/pod-ssh/pod-ssh

![](https://mmbiz.qpic.cn/sz_mmbiz_png/vNHkZdwXUGPicU8ALTCgMicSnVTvm1Jjf0kYaWIeHjdfQB1qRES4icz0WibNgKq64IB2nxERlpAU28OQXdSpjU6ictw/640?wx_fmt=png)
