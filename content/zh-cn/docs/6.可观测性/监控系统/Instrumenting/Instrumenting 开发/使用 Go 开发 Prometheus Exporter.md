---
title: 使用 Go 开发 Prometheus Exporter
---

Exporter 是 Prometheus 监控的核心，如果你遇到一些应用不存在相应的 Exporter，那么我们可以自己去编写 Exporter。下面我们简单介绍如何使用 Golang 来快速编写一个 Exporter。

**1. 安装 GO 和依赖包**

按照 <https://golang.org/doc/install> 上的步骤进行安装配置 GO 环境，创建一个名为 my_first_exporter 的文件夹。

    $ go mod init my_first_exporter
    $ go get github.com/prometheus/client_golang
    $ go get github.com/joho/godotenv
    --> creates go.mod file
    --> Installs dependency into the go.mod file

1
2
3
4
5
Go

**2. 创建入口点和导入依赖包**

    package main
    import (
     "github.com/joho/godotenv"
     "github.com/prometheus/client_golang/prometheus"
     "github.com/prometheus/client_golang/prometheus/promhttp"
    )

1
2
3
4
5
6
7
Go

**3. 创建 main() 函数**

    func main()

1
Go

**4. 添加 prometheus metrics 端点，并在某个服务端口上监听**

    func main() {
       http.Handle("/metrics", promhttp.Handler())
       log.Fatal(http.ListenAndServe(":9101", nil))
    }

1
2
3
4
Go

**5. 使用 curl 请求外部服务接口**

比如我们这里监控的应用程序是 MirthConnect，所以我需要进行两个 API 接口调用：

- 获取 channel 统计数据

- 获取 channel id 和名称映射


    curl -k --location --request GET '<https://apihost/api/channels/statistics>' \\
    --user admin:admin
    curl -k --location --request GET '<https://apihost/api/channels/idsAndNames>' \\
    --user admin:admin

1
2
3
4
5
6
7
Go

**6. 将 curl 调用转换为 go http 调用，并解析结果**

如果你是 Go 新手，这应该是最困难的一步。对于我这里的例子，端点返回的是 XML 格式的数据，所以我必须用 `"encoding/xml"` 包来反序列化 XML。转换成功后意味着我的 GO 程序可以执行和 curl 命令一样的 API 调用。

**7. 声明 Prometheus metrics**

在 Prometheus 中，每个 metric 指标都由以下几个部分组成：`metric name/metric label values/metric help text/metric type/measurement` ，例如：

    Example:
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code=”200"} 1.829159e+06
    promhttp_metric_handler_requests_total{code=”500"} 0
    promhttp_metric_handler_requests_total{code=”503"} 0

1
2
3
4
5
6
7
8
Go

对于应用 scrapers，我们将定义 Prometheus metrics 描述信息，其中包括 metric 名称、metric label 标签以及 metric 帮助信息。

    messagesReceived = prometheus.NewDesc(
     prometheus.BuildFQName(namespace, "", "messages_received_total"),
     "How many messages have been received (per channel).",
     []string{"channel"}, nil,
    )

1
2
3
4
5
Go

**8. 定义一个结构体实现 Prometheus 的 Collector 接口**

Prometheus 的 client 库提供了实现自定义 Exportor 的接口，Collector 接口定义了两个方法 Describe 和 Collect，实现这两个方法就可以暴露自定义的数据：

- Describe(chan<- \*Desc)

- Collect(chan<- Metric)

如下所示：

    type Exporter struct {
     mirthEndpoint, mirthUsername, mirthPassword string
    }
    func NewExporter(mirthEndpoint string, mirthUsername string, mirthPassword string) *Exporter {
     return &Exporter{
      mirthEndpoint: mirthEndpoint,
      mirthUsername: mirthUsername,
      mirthPassword: mirthPassword,
     }
    }
    func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
    }
    func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
    }

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
Go

**9. 在 Describe 函数中，把第 7 步的 metric 描述信息发送给它**

    func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
     ch <- up
     ch <- messagesReceived
     ch <- messagesFiltered
     ch <- messagesQueued
     ch <- messagesSent
     ch <- messagesErrored
    }

1
2
3
4
5
6
7
8
9
10
Plain Text

**10. 将接口调用逻辑从第 6 步移到 Collect 函数中**

直接将采集的数据发送到 `prometheus.Metric` 通道中。

    func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
     channelIdNameMap, err := e.LoadChannelIdNameMap()
     if err != nil {
      ch <- prometheus.MustNewConstMetric(
       up, prometheus.GaugeValue, 0,
      )
      log.Println(err)
      return
     }
     ch <- prometheus.MustNewConstMetric(
      up, prometheus.GaugeValue, 1,
     )
     e.HitMirthRestApisAndUpdateMetrics(channelIdNameMap, ch)
    }

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
Go

当执行 api 调用时，确保使用`prometheus.MustNewConstMetric(prometheus.Desc, metric type, measurement)`发送测量值，如果你需要传入额外的标签，可以像下面这样在参数列表的后面加入：

    channelError, _ := strconv.ParseFloat(channelStatsList.Channels[i].Error, 64)
    ch <- prometheus.MustNewConstMetric(
     messagesErrored, prometheus.GaugeValue, channelError, channelName,
    )

1
2
3
4
Go

**11. 在 main 函数中声明 exporter**

    exporter := NewExporter(mirthEndpoint, mirthUsername, mirthPassword)
    prometheus.MustRegister(exporter)

1
2
Go

到这里其实这个 Exporter 就可以使用了，每次访问 metrics 路由的时候，它会执行 api 调用，并以 Prometheus Text 文本格式返回数据。下面的步骤主要是方便部署了。

**12. 将硬编码的 api 路径放到 flag 中**

前面我们硬编码了好多参数，比如应用程序的网址、metrics 路由地址以及 exporter 端口，我们可以通过从命令行参数中来解析这些值使程序更加灵活。

    var (
    listenAddress = flag.String("web.listen-address", ":9141",
     "Address to listen on for telemetry")
    metricsPath = flag.String("web.telemetry-path", "/metrics",
     "Path under which to expose metrics")
    )
    func main() {
       flag.Parse()
       ...
       http.Handle(*metricsPath, promhttp.Handler())
       log.Fatal(http.ListenAndServe(*listenAddress, nil))
    }

1
2
3
4
5
6
7
8
9
10
11
12
Go

**13. 将凭证放入环境变量**

如果应用端点改变了或者登录凭证改变了怎么办？我们可以从环境变量中来加载这些数据，在这个例子中，我们使用 godotenv 这个包来帮助将变量值存储在本地的一个目录中：

    import (
      "os"
    )
    func main() {
     err := godotenv.Load()
     if err != nil {
      log.Println("Error loading .env file, assume env variables are set.")
     }
     mirthEndpoint := os.Getenv("MIRTH_ENDPOINT")
     mirthUsername := os.Getenv("MIRTH_USERNAME")
     mirthPassword := os.Getenv("MIRTH_PASSWORD")
    }

1
2
3
4
5
6
7
8
9
10
11
12
Go

整个 Exporter 完整的代码如下所示：

```go
package main
import (
 "crypto/tls"
 "encoding/xml"
 "flag"
 "io/ioutil"
 "log"
 "net/http"
 "os"
 "strconv"
 "github.com/joho/godotenv"
 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
)
/*
<map>
  <entry>
    <string>101af57f-f26c-40d3-86a3-309e74b93512</string>
    <string>Send-Email-Notification</string>
  </entry>
</map>
*/
type ChannelIdNameMap struct {
 XMLName xml.Name       `xml:"map"`
 Entries []ChannelEntry `xml:"entry"`
}
type ChannelEntry struct {
 XMLName xml.Name `xml:"entry"`
 Values  []string `xml:"string"`
}
/*
<list>
  <channelStatistics>
    <serverId>c5e6a736-0e88-46a7-bf32-5b4908c4d859</serverId>
    <channelId>101af57f-f26c-40d3-86a3-309e74b93512</channelId>
    <received>0</received>
    <sent>0</sent>
    <error>0</error>
    <filtered>0</filtered>
    <queued>0</queued>
  </channelStatistics>
</list>
*/
type ChannelStatsList struct {
 XMLName  xml.Name       `xml:"list"`
 Channels []ChannelStats `xml:"channelStatistics"`
}
type ChannelStats struct {
 XMLName   xml.Name `xml:"channelStatistics"`
 ServerId  string   `xml:"serverId"`
 ChannelId string   `xml:"channelId"`
 Received  string   `xml:"received"`
 Sent      string   `xml:"sent"`
 Error     string   `xml:"error"`
 Filtered  string   `xml:"filtered"`
 Queued    string   `xml:"queued"`
}
const namespace = "mirth"
const channelIdNameApi = "/api/channels/idsAndNames"
const channelStatsApi = "/api/channels/statistics"
var (
 tr = &http.Transport{
  TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
 }
 client = &http.Client{Transport: tr}
 listenAddress = flag.String("web.listen-address", ":9141",
  "Address to listen on for telemetry")
 metricsPath = flag.String("web.telemetry-path", "/metrics",
  "Path under which to expose metrics")
 // Metrics
 up = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "up"),
  "Was the last Mirth query successful.",
  nil, nil,
 )
 messagesReceived = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "messages_received_total"),
  "How many messages have been received (per channel).",
  []string{"channel"}, nil,
 )
 messagesFiltered = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "messages_filtered_total"),
  "How many messages have been filtered (per channel).",
  []string{"channel"}, nil,
 )
 messagesQueued = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "messages_queued"),
  "How many messages are currently queued (per channel).",
  []string{"channel"}, nil,
 )
 messagesSent = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "messages_sent_total"),
  "How many messages have been sent (per channel).",
  []string{"channel"}, nil,
 )
 messagesErrored = prometheus.NewDesc(
  prometheus.BuildFQName(namespace, "", "messages_errored_total"),
  "How many messages have errored (per channel).",
  []string{"channel"}, nil,
 )
)
type Exporter struct {
 mirthEndpoint, mirthUsername, mirthPassword string
}
func NewExporter(mirthEndpoint string, mirthUsername string, mirthPassword string) *Exporter {
 return &Exporter{
  mirthEndpoint: mirthEndpoint,
  mirthUsername: mirthUsername,
  mirthPassword: mirthPassword,
 }
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
 ch <- up
 ch <- messagesReceived
 ch <- messagesFiltered
 ch <- messagesQueued
 ch <- messagesSent
 ch <- messagesErrored
}
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
 channelIdNameMap, err := e.LoadChannelIdNameMap()
 if err != nil {
  ch <- prometheus.MustNewConstMetric(
   up, prometheus.GaugeValue, 0,
  )
  log.Println(err)
  return
 }
 ch <- prometheus.MustNewConstMetric(
  up, prometheus.GaugeValue, 1,
 )
 e.HitMirthRestApisAndUpdateMetrics(channelIdNameMap, ch)
}
func (e *Exporter) LoadChannelIdNameMap() (map[string]string, error) {
 // Create the map of channel id to names
 channelIdNameMap := make(map[string]string)
 req, err := http.NewRequest("GET", e.mirthEndpoint+channelIdNameApi, nil)
 if err != nil {
  return nil, err
 }
 // This one line implements the authentication required for the task.
 req.SetBasicAuth(e.mirthUsername, e.mirthPassword)
 // Make request and show output.
 resp, err := client.Do(req)
 if err != nil {
  return nil, err
 }
 body, err := ioutil.ReadAll(resp.Body)
 resp.Body.Close()
 if err != nil {
  return nil, err
 }
 // fmt.Println(string(body))
 // we initialize our array
 var channelIdNameMapXML ChannelIdNameMap
 // we unmarshal our byteArray which contains our
 // xmlFiles content into 'users' which we defined above
 err = xml.Unmarshal(body, &channelIdNameMapXML)
 if err != nil {
  return nil, err
 }
 for i := 0; i < len(channelIdNameMapXML.Entries); i++ {
  channelIdNameMap[channelIdNameMapXML.Entries[i].Values[0]] = channelIdNameMapXML.Entries[i].Values[1]
 }
 return channelIdNameMap, nil
}
func (e *Exporter) HitMirthRestApisAndUpdateMetrics(channelIdNameMap map[string]string, ch chan<- prometheus.Metric) {
 // Load channel stats
 req, err := http.NewRequest("GET", e.mirthEndpoint+channelStatsApi, nil)
 if err != nil {
  log.Fatal(err)
 }
 // This one line implements the authentication required for the task.
 req.SetBasicAuth(e.mirthUsername, e.mirthPassword)
 // Make request and show output.
 resp, err := client.Do(req)
 if err != nil {
  log.Fatal(err)
 }
 body, err := ioutil.ReadAll(resp.Body)
 resp.Body.Close()
 if err != nil {
  log.Fatal(err)
 }
 // fmt.Println(string(body))
 // we initialize our array
 var channelStatsList ChannelStatsList
 // we unmarshal our byteArray which contains our
 // xmlFiles content into 'users' which we defined above
 err = xml.Unmarshal(body, &channelStatsList)
 if err != nil {
  log.Fatal(err)
 }
 for i := 0; i < len(channelStatsList.Channels); i++ {
  channelName := channelIdNameMap[channelStatsList.Channels[i].ChannelId]
  channelReceived, _ := strconv.ParseFloat(channelStatsList.Channels[i].Received, 64)
  ch <- prometheus.MustNewConstMetric(
   messagesReceived, prometheus.GaugeValue, channelReceived, channelName,
  )
  channelSent, _ := strconv.ParseFloat(channelStatsList.Channels[i].Sent, 64)
  ch <- prometheus.MustNewConstMetric(
   messagesSent, prometheus.GaugeValue, channelSent, channelName,
  )
  channelError, _ := strconv.ParseFloat(channelStatsList.Channels[i].Error, 64)
  ch <- prometheus.MustNewConstMetric(
   messagesErrored, prometheus.GaugeValue, channelError, channelName,
  )
  channelFiltered, _ := strconv.ParseFloat(channelStatsList.Channels[i].Filtered, 64)
  ch <- prometheus.MustNewConstMetric(
   messagesFiltered, prometheus.GaugeValue, channelFiltered, channelName,
  )
  channelQueued, _ := strconv.ParseFloat(channelStatsList.Channels[i].Queued, 64)
  ch <- prometheus.MustNewConstMetric(
   messagesQueued, prometheus.GaugeValue, channelQueued, channelName,
  )
 }
 log.Println("Endpoint scraped")
}
func main() {
 err := godotenv.Load()
 if err != nil {
  log.Println("Error loading .env file, assume env variables are set.")
 }
 flag.Parse()
 mirthEndpoint := os.Getenv("MIRTH_ENDPOINT")
 mirthUsername := os.Getenv("MIRTH_USERNAME")
 mirthPassword := os.Getenv("MIRTH_PASSWORD")
 exporter := NewExporter(mirthEndpoint, mirthUsername, mirthPassword)
 prometheus.MustRegister(exporter)
 http.Handle(*metricsPath, promhttp.Handler())
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(`<html>
             <head><title>Mirth Channel Exporter</title></head>
             <body>
             <h1>Mirth Channel Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
 })
 log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
```

**14. 编写一个 Makefile 文件，方便在不同平台上快速构建**

Makefile 可以让你在开发过程中省去很多多余的操作，比如我们要构建多个平台的构建程序，可以创建如下所示的 Makefile 文件。

    linux:
       GOOS=linux GOARCH=amd64 go build
    mac:
       GOOS=darwin GOARCH=amd64 go build

只要调用 `make mac` 或 `make linux` 命令就可以看到不同的可执行文件出现。

**15. 编写一个 service 文件，将这个 go 应用作为守护进程运行**

我们可以为这个 Exporter 编写一个 service 文件或者 Dockerfile 文件来管理该应用，比如这里我们直接在 Centos 7 上用 systemd 来管理该应用。这可以编写一个如下所示的 service 文件：

```systemverilog
[Unit]
Description=mirth channel exporter
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
WorkingDirectory=/mirth/mirthconnect
EnvironmentFile=/etc/sysconfig/mirth_channel_exporter
ExecStart=/mirth/mirthconnect/mirth_channel_exporter
[Install]
WantedBy=multi-user.target
```

到这里就完成了使用 Golang 编写一个简单的 Prometheus Exporter。

> “原文链接：<https://medium.com/teamzerolabs/15-steps-to-write-an-application-prometheus-exporter-in-go-9746b4520e26>”
