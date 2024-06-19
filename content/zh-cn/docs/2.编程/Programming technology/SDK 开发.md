---
title: SDK 开发
---

# 概述

> 参考：

# 经典的 SDK 设计方式

- <https://github.com/huaweicloud/huaweicloud-sdk-go-v3>
- <https://github.com/wujiyu115/yuqueg>

目录结构示例

```bash
pkg/my_sdk/
├── README.md
├── core
│   ├── v1
│   │   └── client.go
│   └── v2
│       ├── client.go
│       └── log.go
├── index.go
├── services
│   ├── v1
│   │   ├── book.go
│   │   └── models
│   │       ├── model_request.go
│   │       └── model_response.go
│   └── v2
│       ├── doc.go
│       ├── group.go
│       ├── models
│       │   ├── model_request.go
│       │   └── model_response.go
│       ├── repo.go
│       ├── user.go
│       └── utils.go
└── yuque.log
```

上面的示例每个服务还可以放在单独的目录中，models 中也可以将文件分开，为每个服务创建一个单独的 model 文件。

## 核心客户端

里面包含向目标建立 HTTP 的逻辑

```go
type Client struct {
	token string
}

type RequestOption struct {
	Method string
	Data   map[string]string
}

func (c Client) Request(api string, options *RequestOption) ([]byte, error) {}
func (c Client) RequestObj(api string, container interface{}, options *RequestOption) (interface{}, error) {}
```

这里面的示例，是把将要响应的数据当做 container，作为参数传入，然后通过 HTTP 请求获取到的返回值填到 container，最后返回 container。

在华为云的 SDK 中，是另一种用法

```go
func (c *ElbClient) ListIpGroups(request *model.ListIpGroupsRequest) (*model.ListIpGroupsResponse, error) {
	requestDef := GenReqDefForListIpGroups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListIpGroupsResponse), nil
	}
}
```

使用类型断言，将返回的 interface 类型的数据固定下来后再返回

## 各个服务客户端

```go
type UserService struct {
	client *Client
}

func (c UserService) Get(login string) (*UserInfo, error) {
	_, err := c.client.RequestObj(url, &user, EmptyRO)
}
```

### 响应体与请求体

响应体与请求体的通常定义为一个 Struct，这些 Struct 通常放在服务目录下的 models 目录中
