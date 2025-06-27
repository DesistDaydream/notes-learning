---
title: Kubelet 启动流程
---

# 概述

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/frggyc/1645579695318-7cd1eb33-cb95-4223-b6fd-299d52112b20.png)

# 启动

首先从 kubelet 的 `main()` 函数开始，调用 `app.NewKubeletCommand()` 方法以获取配置文件中的参数、校验参数、为参数设置默认值。主要逻辑为：

- 解析命令行参数；
- 为 kubelet 初始化 feature gates 参数；
- 加载 kubelet 配置文件；
- 校验配置文件中的参数；
- 检查 kubelet 是否启用动态配置功能；
- 初始化 kubeletDeps，kubeletDeps 包含 kubelet 运行所必须的配置，是为了实现 dependency injection，其目的是为了把 kubelet 依赖的组件对象作为参数传进来，这样可以控制 kubelet 的行为；
- 调用 `Run()` 函数；

## main() - 入口

源码：`[cmd/kubelet/kubelet.go](https://github.com/kubernetes/kubernetes/blob/master/cmd/kubelet/kubelet.go)`

```go
func main() {
 command := app.NewKubeletCommand()

 // kubelet 使用一个配置文件，并对标志和该配置文件进行自己的特殊解析。
    // 完成后，它会初始化日志记录。因此，它不像其他更简单的命令那样使用 cli.Run()
 code := run(command)
 os.Exit(code)
}

func run(command *cobra.Command) int {
 defer logs.FlushLogs()
 rand.Seed(time.Now().UnixNano())

 command.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)
 if err := command.Execute(); err != nil {
  return 1
 }
 return 0
}
```

## NewKubeletCommand() - Cobra 库的基本逻辑

源码：`[cmd/kubelet/app/server.go](https://github.com/kubernetes/kubernetes/blob/master/cmd/kubelet/app/server.go)`

```go
func NewKubeletCommand() *cobra.Command {
 cleanFlagSet := pflag.NewFlagSet(componentKubelet, pflag.ContinueOnError)
 cleanFlagSet.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

    // 1、kubelet 配置分两部分:
    // KubeletFlag: 指那些不允许在 kubelet 运行时进行修改的配置集，或者不能在集群中各个 Nodes 之间共享的配置集。
    // KubeletConfiguration: 指可以在集群中各个Nodes之间共享的配置集，可以进行动态配置。
    kubeletFlags := options.NewKubeletFlags()
    kubeletConfig, err := options.NewKubeletConfiguration()

    cmd := &cobra.Command{
        Use: componentKubelet,
        DisableFlagParsing: true,
        ......
        RunE: func(cmd *cobra.Command, args []string) error {
            // 2、解析命令行标志，这里禁用了 cobra 的标志解析
   if err := cleanFlagSet.Parse(args); err != nil {
    return fmt.Errorf("failed to parse kubelet flag: %w", err)
   }
            ...... // 一系列验证行为


            // 3、初始化 feature gates 配置
   if err := utilfeature.DefaultMutableFeatureGate.SetFromMap(kubeletConfig.FeatureGates); err != nil {
    return fmt.Errorf("failed to set feature gates from initial flags-based config: %w", err)
   }

   // 验证初始的 kubelet 标志
   if err := options.ValidateKubeletFlags(kubeletFlags); err != nil {
    return fmt.Errorf("failed to validate kubelet flags: %w", err)
   }

            if kubeletFlags.ContainerRuntime == "remote" && cleanFlagSet.Changed("pod-infra-container-image") {
                klog.Warning("Warning: For remote container runtime, --pod-infra-container-image is ignored in kubelet, which should be set in that      remote runtime instead")
            }

            // 4、加载 kubelet 配置文件
            if configFile := kubeletFlags.KubeletConfigFile; len(configFile) > 0 {
                kubeletConfig, err = loadConfigFile(configFile)
                ......
            }
            // 5、校验配置文件中的参数
   if err := kubeletconfigvalidation.ValidateKubeletConfiguration(kubeletConfig); err != nil {
    return fmt.Errorf("failed to validate kubelet configuration, error: %w, path: %s", err, kubeletConfig)
   }

            // 6、检查 kubelet 是否启用动态配置功能
   if utilfeature.DefaultFeatureGate.Enabled(features.DynamicKubeletConfig) {
    return fmt.Errorf("cannot set feature gate %v to %v, feature is locked to %v", features.DynamicKubeletConfig, true, false)
   }

            // 配置和标志解析完成后，开始初始化日志
            logs.InitLogs()

            // 从 kubeletFlags 和 kubeletConfig 构造一个 kubeletServer
            kubeletServer := &options.KubeletServer{
                KubeletFlags:         *kubeletFlags,
                KubeletConfiguration: *kubeletConfig,
            }

            // 7、初始化 kubeletDeps
            kubeletDeps, err := UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)

            // 检查 kubelet 是否以 root 权限启动
            if err := checkPermissions(); err != nil {
   }

   // 8、调用 Run 方法，即开始运行 kubelet
   return Run(ctx, kubeletServer, kubeletDeps, utilfeature.DefaultFeatureGate)
        },
    }
    kubeletFlags.AddFlags(cleanFlagSet)
    options.AddKubeletConfigFlags(cleanFlagSet, kubeletConfig)
    options.AddGlobalFlags(cleanFlagSet)
    ......

    return cmd
}
```

### CRI 与 CNI 的默认配置

<https://github.com/kubernetes/kubernetes/blob/release-1.22/cmd/kubelet/app/options/container_runtime.go>

```bash
const (
 // When these values are updated, also update test/utils/image/manifest.go
 defaultPodSandboxImageName    = "k8s.gcr.io/pause"
 defaultPodSandboxImageVersion = "3.5"
)

var (
 defaultPodSandboxImage = defaultPodSandboxImageName +
  ":" + defaultPodSandboxImageVersion
)

// NewContainerRuntimeOptions will create a new ContainerRuntimeOptions with
// default values.
func NewContainerRuntimeOptions() *config.ContainerRuntimeOptions {
 dockerEndpoint := ""
 if runtime.GOOS != "windows" {
  dockerEndpoint = "unix:///var/run/docker.sock"
 }

 return &config.ContainerRuntimeOptions{
  ContainerRuntime:          kubetypes.DockerContainerRuntime,
  DockerEndpoint:            dockerEndpoint,
  DockershimRootDirectory:   "/var/lib/dockershim",
  PodSandboxImage:           defaultPodSandboxImage,
  ImagePullProgressDeadline: metav1.Duration{Duration: 1 * time.Minute},

  CNIBinDir:   "/opt/cni/bin",
  CNIConfDir:  "/etc/cni/net.d",
  CNICacheDir: "/var/lib/cni/cache",
 }
}
```

## Run() - 启动 kubelet

`Run()` 函数仅仅调用 `run()` 函数以执行启动逻辑。

- `Run()` 使用给定的依赖(即 \*kubelet.Dependencies 参数)运行指定的 KubeletServer。这不应该退出。 kubeDeps 参数可能是 nil - 如果是这样，它是从 KubeletServer 上的设置初始化的。否则，假定调用者已设置 Dependencies 对象，并且不会生成默认对象。

源码：`cmd/kubelet/app/server.go`

```go
func Run(ctx context.Context, s *options.KubeletServer, kubeDeps *kubelet.Dependencies, featureGate featuregate.FeatureGate) error {
 if err := initForOS(s.KubeletFlags.WindowsService, s.KubeletFlags.WindowsPriorityClass); err != nil {
  return fmt.Errorf("failed OS init: %w", err)
 }
 if err := run(ctx, s, kubeDeps, featureGate); err != nil {
  return fmt.Errorf("failed to run Kubelet: %w", err)
 }
 return nil
}
```

## run() - 运行 kubelet 前配置及检查

`run()` 函数中主要是为 kubelet 的启动做一些基本的配置及检查工作，主要逻辑为：

- 为 kubelet 设置默认的 FeatureGates，kubelet 所有的 FeatureGates 可以通过命令参数查看，k8s 中处于 Alpha 状态的 FeatureGates 在组件启动时默认关闭，处于 Beta 和 GA 状态的默认开启
- 校验 kubelet 的参数
- 尝试获取 kubelet 的 lock file，需要在 kubelet 启动时指定 --exit-on-lock-contention 和 --lock-file，该功能处于 Alpha 版本默认为关闭状态
- 将当前的配置文件注册到 http server /configz URL 中
- 检查 kubelet 启动模式是否为 standalone 模式，此模式下不会和 apiserver 交互，主要用于 kubelet 的调试
- 初始化 kubeDeps，kubeDeps 中包含 kubelet 的一些依赖，主要有 KubeClient、EventClient、HeartbeatClient、Auth、cadvisor、ContainerManager
- 为进程设置 oom 分数，默认为 -999，分数范围为 \[-1000, 1000]，越小越不容易被 kill 掉
- 调用 PreInitRuntimeService() 函数，以初始化运行时
- 调用 RunKubelet() 函数，以继续执行运行 kubelet 的后续代码
- 检查 kubelet 是否启动了动态配置功能
- 启动 Healthz http server
- 如果使用 systemd 启动，通知 systemd kubelet 已经启动

源码：`cmd/kubelet/app/server.go`

```go
func run(ctx context.Context, s *options.KubeletServer, kubeDeps *kubelet.Dependencies, featureGate featuregate.FeatureGate) (err error) {
 // 1、根据初始 KubeletServer 设置默认的特性
 err = utilfeature.DefaultMutableFeatureGate.SetFromMap(s.KubeletConfiguration.FeatureGates)

    // 2、验证初始的 KubeletServer(我们先设置了特性，因为这个验证依赖于特性)
 err := options.ValidateKubeletServer(s)

 // 如果使用 cgroups v1 启用 MemoryQoS，则发出警告
 if utilfeature.DefaultFeatureGate.Enabled(features.MemoryQoS) && !isCgroup2UnifiedMode() {
 }
 // 3、获取 kubelet 的锁文件
 if s.ExitOnLockContention && s.LockFilePath == "" {
 }

 // 4、将当前的配置文件注册到 /configz 端点上，也就是说，通过 HTTP 访问 kubelet
 err = initConfigz(&s.KubeletConfiguration)

 // 5、即将获取 client，判断是否为 standalone 模式
 standaloneMode := true
 if len(s.KubeConfig) > 0 {
  standaloneMode = false
 }

    // 6、初始化 kubeDeps
 kubeDeps, err = UnsecuredDependencies(s, featureGate)
 hostName, err := nodeutil.GetHostname(s.HostnameOverride)
 nodeName, err := getNodeName(kubeDeps.Cloud, hostName)

 // 7、如果是 standalone 模式将所有 client 设置为 nil
 switch {
 case standaloneMode:
  kubeDeps.KubeClient = nil
  kubeDeps.EventClient = nil
  kubeDeps.HeartbeatClient = nil
  // 8、为 kubeDeps 初始化 KubeClient、EventClient、HeartbeatClient 模块
 case kubeDeps.KubeClient == nil, kubeDeps.EventClient == nil, kubeDeps.HeartbeatClient == nil:
        clientConfig, closeAllConns, err := buildKubeletClientConfig(ctx, s, nodeName)
        kubeDeps.KubeClient, err = clientset.NewForConfig(clientConfig)
        ......
 }
 // 9、初始化 auth 模块
 if kubeDeps.Auth == nil {
  auth, runAuthenticatorCAReload, err := BuildAuth(nodeName, kubeDeps.KubeClient, s.KubeletConfiguration)
  kubeDeps.Auth = auth
  runAuthenticatorCAReload(ctx.Done())
 }

    // 10、设置 CGroup Root
 var cgroupRoots []string
 nodeAllocatableRoot := cm.NodeAllocatableRoot(s.CgroupRoot, s.CgroupsPerQOS, s.CgroupDriver)
 cgroupRoots = append(cgroupRoots, nodeAllocatableRoot)
 kubeletCgroup, err := cm.GetKubeletContainer(s.KubeletCgroups)

 // 11、初始化 CAdvisor
 if kubeDeps.CAdvisorInterface == nil {
  imageFsInfoProvider := cadvisor.NewImageFsInfoProvider(s.RemoteRuntimeEndpoint)
  kubeDeps.CAdvisorInterface, err = cadvisor.New(imageFsInfoProvider, s.RootDirectory, cgroupRoots, cadvisor.UsingLegacyCadvisorStats(s.RemoteRuntimeEndpoint))
 }

 // Setup event recorder if required.
 makeEventRecorder(kubeDeps, nodeName)

    // 12、初始化 ContainerManager(即.容器管理器)
 if kubeDeps.ContainerManager == nil {
        ......
 }

    // 14、为 kubelet 进程设置 oom 分数
 // TODO(vmarmol): 通过container config执行此操作。
 oomAdjuster := kubeDeps.OOMAdjuster

    // 在 RunKubelet()(即.运行 kubelet) 初始化 runtime 服务
 err = kubelet.PreInitRuntimeService(&s.KubeletConfiguration, kubeDeps, s.RemoteRuntimeEndpoint, s.RemoteImageEndpoint)

    // 15、设置和运行 kubelet
 err := RunKubelet(s, kubeDeps, s.RunOnce)

    // 16、启动 Healthz http server
 if s.HealthzPort > 0 {
  mux := http.NewServeMux()
  healthz.InstallHandler(mux)
  go wait.Until(func() {
   err := http.ListenAndServe(net.JoinHostPort(s.HealthzBindAddress, strconv.Itoa(int(s.HealthzPort))), mux)
  }, 5*time.Second, wait.NeverStop)
 }

 if s.RunOnce {
  return nil
 }
 // 17、如果使用了 systemd，则向 systemd 发送启动信号
 go daemon.SdNotify(false, "READY=1")

 select {
 case <-done:
  break
 case <-ctx.Done():
  break
 }

 return nil
}
```

## PreInitRuntimeService() - 初始化 runtime

源码：`[pkg/kubelet/kubelet.go](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/kubelet.go)`
PreInitRuntimeService()

## RunKubelet() - 运行 kubelet

`RunKubelet()` 中主要是两个行为

- 调用 `createAndInitKubelet()` 执行 kubelet 组件的初始化
- 然后调用 `startKubelet()` 启动 kubelet 中的组件。

源码：`cmd/kubelet/app/server.go`

```go
func RunKubelet(kubeServer *options.KubeletServer, kubeDeps *kubelet.Dependencies, runOnce bool) error {
 hostname, err := nodeutil.GetHostname(kubeServer.HostnameOverride)
 if err != nil {
  return err
 }
 // Query the cloud provider for our node name, default to hostname if kubeDeps.Cloud == nil
 nodeName, err := getNodeName(kubeDeps.Cloud, hostname)
 if err != nil {
  return err
 }
 hostnameOverridden := len(kubeServer.HostnameOverride) > 0
 // Setup event recorder if required.
 makeEventRecorder(kubeDeps, nodeName)

 var nodeIPs []net.IP
 if kubeServer.NodeIP != "" {
  for _, ip := range strings.Split(kubeServer.NodeIP, ",") {
   parsedNodeIP := netutils.ParseIPSloppy(strings.TrimSpace(ip))
   if parsedNodeIP == nil {
    klog.InfoS("Could not parse --node-ip ignoring", "IP", ip)
   } else {
    nodeIPs = append(nodeIPs, parsedNodeIP)
   }
  }
 }

 if len(nodeIPs) > 2 || (len(nodeIPs) == 2 && netutils.IsIPv6(nodeIPs[0]) == netutils.IsIPv6(nodeIPs[1])) {
  return fmt.Errorf("bad --node-ip %q; must contain either a single IP or a dual-stack pair of IPs", kubeServer.NodeIP)
 } else if len(nodeIPs) == 2 && kubeServer.CloudProvider != "" {
  return fmt.Errorf("dual-stack --node-ip %q not supported when using a cloud provider", kubeServer.NodeIP)
 } else if len(nodeIPs) == 2 && (nodeIPs[0].IsUnspecified() || nodeIPs[1].IsUnspecified()) {
  return fmt.Errorf("dual-stack --node-ip %q cannot include '0.0.0.0' or '::'", kubeServer.NodeIP)
 }
 // 1、默认启动特权模式
 capabilities.Initialize(capabilities.Capabilities{
  AllowPrivileged: true,
 })

 credentialprovider.SetPreferredDockercfgPath(kubeServer.RootDirectory)
 klog.V(2).InfoS("Using root directory", "path", kubeServer.RootDirectory)

 if kubeDeps.OSInterface == nil {
  kubeDeps.OSInterface = kubecontainer.RealOS{}
 }

 if kubeServer.KubeletConfiguration.SeccompDefault && !utilfeature.DefaultFeatureGate.Enabled(features.SeccompDefault) {
  return fmt.Errorf("the SeccompDefault feature gate must be enabled in order to use the SeccompDefault configuration")
 }

    // 2、调用 createAndInitKubelet() 函数，以初始化 kubelet
 k, err := createAndInitKubelet(&kubeServer.KubeletConfiguration,
  kubeDeps,
  &kubeServer.ContainerRuntimeOptions,
  hostname,
  hostnameOverridden,
  nodeName,
  nodeIPs,
  kubeServer.ProviderID,
  kubeServer.CloudProvider,
  kubeServer.CertDirectory,
  kubeServer.RootDirectory,
  kubeServer.ImageCredentialProviderConfigFile,
  kubeServer.ImageCredentialProviderBinDir,
  kubeServer.RegisterNode,
  kubeServer.RegisterWithTaints,
  kubeServer.AllowedUnsafeSysctls,
  kubeServer.ExperimentalMounterPath,
  kubeServer.KernelMemcgNotification,
  kubeServer.ExperimentalCheckNodeCapabilitiesBeforeMount,
  kubeServer.ExperimentalNodeAllocatableIgnoreEvictionThreshold,
  kubeServer.MinimumGCAge,
  kubeServer.MaxPerPodContainerCount,
  kubeServer.MaxContainerCount,
  kubeServer.MasterServiceNamespace,
  kubeServer.RegisterSchedulable,
  kubeServer.KeepTerminatedPodVolumes,
  kubeServer.NodeLabels,
  kubeServer.NodeStatusMaxImages,
  kubeServer.KubeletFlags.SeccompDefault || kubeServer.KubeletConfiguration.SeccompDefault,
 )
 if err != nil {
  return fmt.Errorf("failed to create kubelet: %w", err)
 }

 // NewMainKubelet should have set up a pod source config if one didn't exist
 // when the builder was run. This is just a precaution.
 if kubeDeps.PodConfig == nil {
  return fmt.Errorf("failed to create kubelet, pod source config was nil")
 }
 podCfg := kubeDeps.PodConfig

 if err := rlimit.SetNumFiles(uint64(kubeServer.MaxOpenFiles)); err != nil {
  klog.ErrorS(err, "Failed to set rlimit on max file handles")
 }

 // process pods and exit.
 if runOnce {
  if _, err := k.RunOnce(podCfg.Updates()); err != nil {
   return fmt.Errorf("runonce failed: %w", err)
  }
  klog.InfoS("Started kubelet as runonce")
 } else {
        // 3、调用 startKubelet() 函数继续执行后续 kubelet 启动逻辑
  startKubelet(k, podCfg, &kubeServer.KubeletConfiguration, kubeDeps, kubeServer.EnableServer)
  klog.InfoS("Started kubelet")
 }
 return nil
}
```

## createAndInitKubelet() - 初始化 kubelet

`createAndInitKubelet()` 中主要调用了 一个函数，两个方法 来完成 kubelet 的初始化：

- `NewMainKubelet()` # 实例化 kubelet 对象，并对 kubelet 依赖的所有模块进行初始化；
- `kubelet.BirthCry()` # 向 apiserver 发送一条 kubelet 启动了的 event；
- `kubelet.StartGarbageCollection` # 启动垃圾回收服务，回收 container 和 images；

代码：`cmd/kubelet/app/server.go`

### NewMainKubelet() - 实例化 kubelet

`NewMainKubelet()` 是初始化 kubelet 的一个函数，主要逻辑为：

- 初始化 PodConfig 即监听 pod 元数据的来源(file，http，apiserver)，将不同 source 的 pod configuration 合并到一个结构中；
- 初始化 containerGCPolicy、imageGCPolicy、evictionConfig 配置；
- 启动 serviceInformer 和 nodeInformer；
- 初始化 containerRefManager、oomWatcher；
- 初始化 kubelet 对象；
- 初始化 secretManager、configMapManager；
- 初始化 livenessManager、podManager、statusManager、resourceAnalyzer；
- 调用 kuberuntime.NewKubeGenericRuntimeManager 初始化 containerRuntime；
- 初始化 pleg；
- 初始化 containerGC、containerDeletor、imageManager、containerLogManager；
- 初始化 serverCertificateManager、probeManager、tokenManager、volumePluginMgr、pluginManager、volumeManager；
- 初始化 workQueue、podWorkers、evictionManager；
- 最后注册相关模块的 handler；

`NewMainKubelet()` 中对 kubelet 依赖的所有模块进行了初始化，每个模块对应的功能在上篇文章“kubelet 架构浅析”有介绍，至于每个模块初始化的流程以及功能会在后面的文章中进行详细分析。

源码：`pkg/kubelet/kubelet.go`

## startKubelet() - 开始运行 kubelet，并监听端口

在 `startKubelet()` 中通过调用 `k.Run()` 来启动 kubelet 中的所有模块以及主流程，然后启动 kubelet 所需要的 http server，在 v1.16 中，kubelet 默认仅启动健康检查端口 10248 和 kubelet server 的端口 10250。

源码：`cmd/kubelet/app/server.go`

```go
func startKubelet(k kubelet.Bootstrap, podCfg *config.PodConfig, kubeCfg *kubeletconfiginternal.KubeletConfiguration, kubeDeps *kubelet.Dependencies, enableServer bool) {
 // kubelet 真正启动，在这里面会有开始同步循环并监听 Pod 的逻辑调用，即启动 syncLoop
 go k.Run(podCfg.Updates())

 // 同时启动 kubelet 的 HTTP 服务，并开始监听端口
 if enableServer {
  go k.ListenAndServe(kubeCfg, kubeDeps.TLSOptions, kubeDeps.Auth)
 }
 if kubeCfg.ReadOnlyPort > 0 {
  go k.ListenAndServeReadOnly(netutils.ParseIPSloppy(kubeCfg.Address), uint(kubeCfg.ReadOnlyPort))
 }
 if utilfeature.DefaultFeatureGate.Enabled(features.KubeletPodResources) {
  go k.ListenAndServePodResources()
 }
}
```

# 运行

## kubelet.Run()

`kubelet.Run()` 方法是**启动 kubelet 的核心方法**，其中会启动 kubelet 的依赖模块以及主循环逻辑，这是实现了 [Bootstrap](/docs/10.云原生/Kubernetes/Kubernetes%20开发/源码解析/Kubelet/Kubelet.md#Bootstrap%20接口) 接口的 [Kubelet](/docs/10.云原生/Kubernetes/Kubernetes%20开发/源码解析/Kubelet/Kubelet.md#Kubelet%20结构体) 结构体的方法。该方法的主要逻辑为：

- 注册 logServer；
- 判断是否需要启动 cloud provider sync manager；
- 调用 kl.initializeModules 首先启动不依赖 container runtime 的一些模块；
- 启动 volume manager；
- 执行 kl.syncNodeStatus 定时同步 Node 状态；
- 调用 kl.fastStatusUpdateOnce 更新容器运行时启动时间以及执行首次状态同步；
- 判断是否启用 NodeLease 机制；
- 执行 kl.updateRuntimeUp 定时更新 Runtime 状态；
- 执行 kl.syncNetworkUtil 定时同步 iptables 规则；
- 执行 kl.podKiller 定时清理异常 pod，当 pod 没有被 podworker 正确处理的时候，启动一个 goroutine 负责 kill 掉 pod；
- 启动 statusManager；
- 启动 probeManager；
- 启动 runtimeClassManager；
- 启动 pleg；
- 调用 kl.syncLoop 监听 pod 变化；

在 `kubelet.Run()` 方法中主要调用了两个方法 kl.initializeModules 和 kl.fastStatusUpdateOnce 来完成启动前的一些初始化，在初始化完所有的模块后会启动主循环。

源码：`pkg/kubelet/kubelet.go`

```go
func (kl *Kubelet) Run(updates <-chan kubetypes.PodUpdate) {
    // 1、在 HTTP Server 的 /logs/ 端点上注册日志服务
 kl.logServer = http.StripPrefix("/logs/", http.FileServer(http.Dir("/var/log/")))

 // 2、判断是否需要启动云提供商的同步管理器
 if kl.cloudResourceSyncManager != nil {
  go kl.cloudResourceSyncManager.Run(wait.NeverStop)
 }

    // 3、调用 kubelet.initializeModules 首先启动不依赖容器 runtime 的一些模块
 if err := kl.initializeModules(); err != nil {
  kl.recorder.Eventf(kl.nodeRef, v1.EventTypeWarning, events.KubeletSetupFailed, err.Error())
 }

 // 4、启动 Volume 管理器
 go kl.volumeManager.Run(kl.sourcesReady, wait.NeverStop)

 if kl.kubeClient != nil {
  // Introduce some small jittering to ensure that over time the requests won't start
  // accumulating at approximately the same time from the set of nodes due to priority and
  // fairness effect.
        // 5、执行 kl.syncNodeStatus 定时同步 Node 状态
  go wait.JitterUntil(kl.syncNodeStatus, kl.nodeStatusUpdateFrequency, 0.04, true, wait.NeverStop)
  // 6、调用 kl.fastStatusUpdateOnce 更新容器运行时启动时间以及执行首次状态同步
        go kl.fastStatusUpdateOnce()

  // 7、弃用 NodeLease 机制
  go kl.nodeLeaseController.Run(wait.NeverStop)
 }
    // 8、执行 kl.updateRuntimeUp 定时更新 Runtime 状态
 go wait.Until(kl.updateRuntimeUp, 5*time.Second, wait.NeverStop)

 // 9、执行 kl.syncNetworkUtil 定时同步 iptables 规则
 if kl.makeIPTablesUtilChains {
  kl.initNetworkUtil()
 }

    // Start component sync loops.(启动组件的同步循环)
 kl.statusManager.Start()

 // 启动 runtime 类管理器
 if kl.runtimeClassManager != nil {
  kl.runtimeClassManager.Start(wait.NeverStop)
 }

 // 启动 Pod 生命周期事件生成器
 kl.pleg.Start()
    // 13、调用 kubelet.syncLoop 监听 Pod 变化
 kl.syncLoop(updates, kl)
}
```

### kubelet.initializeModules() - 启动不依赖容器 Runtime 的模块

`initializeModules()` 中启动的模块是不依赖于容器 Runtime 的，并且不依赖于尚未初始化的模块，其主要逻辑为：

- 调用 kl.setupDataDirs 创建 kubelet 所需要的文件目录；
- 创建 ContainerLogsDir /var/log/containers；
- 启动 imageManager，image gc 的功能已经在 RunKubelet 中启动了，此处主要是监控 image 的变化；
- 启动 certificateManager，负责证书更新；
- 启动 oomWatcher，监听 oom 并记录事件；
- 启动 resourceAnalyzer；

### kubelet.fastStatusUpdateOnce()

`fastStatusUpdateOnce()` 会不断尝试更新 pod CIDR，一旦更新成功会立即执行 updateRuntimeUp 和 syncNodeStatus 来进行运行时的更新和节点状态更新。此方法只在 kubelet 启动时执行一次，目的是为了通过更新 pod CIDR，减少节点达到 ready 状态的时延，尽可能快的进行 runtime update 和 node status update。

### kubelet.updateRuntimeUp() - 启动依赖容器 Runtime 的模块

`updateRuntimeUp()` 方法在容器运行时首次启动过程中初始化运行时依赖的模块，并在 kubelet 的 runtimeState 中更新容器运行时的启动时间。updateRuntimeUp() 方法主要逻辑：

- 获取容器运行时状态
- 检查 network 以及 runtime 是否处于 ready 状态
- 如果 network 以及 runtime 都处于 ready 状态，则调用 `kubelet.initializeRuntimeDependentModules()` 初始化依赖容器 Runtime 的模块：
  - 包括 cadvisor、containerManager、evictionManager、containerLogManager、pluginManage 等。

### kubelet.initializeRuntimeDependentModules() - 启动依赖容器 Runtime 的模块

该方法的主要逻辑为：

- 启动 cadvisor；
- 获取 CgroupStats；
- 启动 containerManager、evictionManager、containerLogManager；
- 将 CSI Driver 和 Device Manager 注册到 pluginManager，然后启动 pluginManager；

# kubelet.syncLoop() - 同步循环，Kubelet 开始运行

相见[《Kubelet 同步循环》](SyncLoop 模块 # Kubelet 同步循环.md)章节
