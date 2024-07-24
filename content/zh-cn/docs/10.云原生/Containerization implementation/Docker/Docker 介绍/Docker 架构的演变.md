---
title: Docker 架构的演变
---

原文链接：<https://blog.csdn.net/csdnnews/article/details/90746002>

Docker 是如何工作的？这是一个简单的问题，但答案却是出乎意料的复杂。你可能很多次地听说过“守护进程(daemon)”和“运行时(runtime)”这两个术语，但可能从未真正理解它们的含义以及它们是如何配合在一起的。如果你像我一样，涉过源头去发现真相，那么，在你沉溺于代码之海时，你并不孤单。让我们面对现实吧，假想 Docker 的源代码是一顿意式大餐，而你正在狼吞虎咽你的美味意面。

就像一把叉子可以把面条送到你的口中，这篇文章会将 Docker 的技术的方方面面组织在一起并导入你饥饿的大脑。

为了更好地理解现在，我们首先需要回顾过去。2013 年，dotCloud 公司的 Solomon Hykes 在那年的 Python 大会上发表了 Linux 容器的未来的演讲（<https://www.youtube.com/watch?v=wW9CAH9nSLs>，需科学上网），第一次将 Docker 带入了公众的视线。让我们将他的 git 代码库回溯到 2013 年 1 月，这个 Docker 开发更轻松的时间。

**Docker 2013 是如何工作的？**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf39dp/1616122015461-9720f60e-baf0-4aa5-b9e5-985d6be8b373.png)

Docker 由两个主要组件组成，一个供用户使用的命令行应用程序和一个管理容器的守护进程。这个守护进程依赖两个子组件来执行它的任务：在宿主主机文件系统上用来存储镜像和容器数据的存储组件；以及用于抽象原始内核调用来构建 Linux 容器的 LXC 接口。

**命令行应用程序**

Docker 命令行应用程序是管理你的 Docker 运行副本已知的所有镜像和容器的人工界面。它相对简单，因为所有的管理都是由守护进程完成的。应用程序开始于一个 main 函数：

```go
funcmain() {
var err error
    ...
// Example:              "/var/run/docker.sock", "run"
    conn, err := rcli.CallTCP(os.Getenv("DOCKER"), os.Args[1:]...)
    ...
    receive_stdout := future.Go(func()error {
        _, err := io.Copy(os.Stdout, conn)
return err
    })
    ...
}
```

它会立即建立一个 TCP 连接，指向存储在 DOCKER 这个环境变量中的地址，这是 Docker 守护进程的地址。用户提供的参数发送给守护进程，然后命令行应用程序等待打印成功答复的结果。

**dockerd**

Docker 守护进程的代码存放在同一个代码库中，这个进程称为 dockerd。它运行在后台来处理用户请求和容器清理工作。启动后，dockerd 将监听在 8080 端口上传入的 HTTP 连接和在 4242 端口上的 TCP 连接。

```go
func main() {    ...    // d is the server, it will process requests made by the user    d, err := New()    ...    go func() {        if err := rcli.ListenAndServeHTTP(":8080", d); err != nil {            log.Fatal(err)        }    }()    if err := rcli.ListenAndServeTCP(":4242", d); err != nil {        log.Fatal(err)    }}
    ...
    // d is the server, it will process requests made by the user
    d, err := New()
    ...
    go func() {
        if err := rcli.ListenAndServeHTTP(":8080", d); err != nil {
            log.Fatal(err)
        }
    }()
    if err := rcli.ListenAndServeTCP(":4242", d); err != nil {
        log.Fatal(err)
    }
}
```

一旦命令接收到后，dockerd 将使用反射机制查找并调用要运行的函数。

**docker run 命令**

其中一个运行的函数是 CmdRun，它对应这个 docker run（注：这个命令创建一个新的容器并运行一个命令）命令。

    func(srv *Server)CmdRun(stdin io.ReadCloser, stdout io.Writer, args ...string)error {
        flags := rcli.Subcmd(stdout, "run", "[OPTIONS] IMAGE COMMAND [ARG...]", "Run a command in a new container")
        ...
    // Choose a default image if needed
    if name == "" {
            name = "base"
        }
    // Choose a default command if needed
    iflen(cmd) == 0 {
            *fl_stdin = true
            *fl_tty = true
            *fl_attach = true
            cmd = []string{"/bin/bash", "-i"}
        }
        ...
    // Find the image
        img := srv.images.Find(name)
        ...
    // Create Container
        container := srv.CreateContainer(img, *fl_tty, *fl_stdin, *fl_comment, cmd[0], cmd[1:]...)
        ...
    // Start Container
        container.Start()
        ...
    }

用户通常会提供一个镜像和一个命令，以便 dockerd 运行。当它们被省略时，将默认使用镜像 base 和命令/bin/bash。

**查找镜像**

然后，我们通过将名称（或 ID）映射到文件系统上的一个位置来找到指定的镜像（假设前面已经执行过 docker pull （注：这个命令从镜像仓库中拉取或者更新指定镜像）命令，镜像已经生成）。

    type Index struct {
        Path    string// "/var/lib/docker/images/index.json"
        ByName  map[string]*History
        ById    map[string]*Image
    }

    func(index *Index)Find(idOrName string) *Image {
        ...
    // Lookup by ID
    if image, exists := index.ById[idOrName]; exists {
    return image
        }
    // Lookup by name
    if history, exists := index.ByName[idOrName]; exists && history.Len() > 0 {
    return (*history)[0]
        }
    returnnil
    }

在这个版本的 Docker 中，所有镜像都存储在/var/lib/docker/images 文件夹中。想进一步了解 Docker 镜像中有什么，请参阅我以前的博客文章（<https://cameronlonsdale.com/2018/11/26/whats-in-a-docker-image/>）。

**创建容器**

接着我们开始创建容器。dockerd 创建了一个结构（struct）来保存与这个容器相关的所有元数据，并将它存储在一个列表中以便于访问。

    container := &Container{    // Examples
        Id:         id,         // "09906fa3"
        Root:       root,       // /var/lib/docker/containers/09906fa3/"
        Created:    time.Now(),
        Path:       command,    // "/bin/bash"
        Args:       args,       // ["-i"]
        Config:     config,

    // "/var/lib/docker/containers/09906fa3/rootfs"
    // "/var/lib/docker/containers/09906fa3/rw"
        Filesystem: newFilesystem(path.Join(root, "rootfs"), path.Join(root, "rw"), layers),
        State:      newState(),

    // "/var/lib/docker/containers/09906fa3/config.lxc"
        lxcConfigPath: path.Join(root, "config.lxc"),
    stdout:        newWriteBroadcaster(),
    stderr:        newWriteBroadcaster(),
        stdoutLog:     new(bytes.Buffer),
        stderrLog:     new(bytes.Buffer),
    }
    ...
    // Create directories
    os.Mkdir(root, 0700);
    container.Filesystem.createMountPoints();
    ...
    // Generate LXC Config file
    container.generateLXCConfig();

当创建这个 struct 时，dockerd 会在下列路径为容器创建下面这个唯一目录：/var/lib/docker/containers/<ID>。在这个目录下中有两个子目录：/rootfs 只读根文件系统层（来自联合挂载（union mounted)的镜像），/rw 提供一个单独的读写层，供容器来创建临时文件。

最后，使用我们新创建的容器数据填充一个模板，用以生成 LXC 配置文件。更多关于 LXC 的信息，请参见下一节。

**运行容器**

我们的容器终于创建出来了！但它还没有运行，现在让我们启动它。

    func(container *Container)Start()error {
    // Mount file system if not mounted
        container.Filesystem.EnsureMounted();

        params := []string{
    "-n", container.Id,
    "-f", container.lxcConfigPath,
    "--",
            container.Path,
        }
        params = append(params, container.Args...)

    // /usr/bin/lxc-start -n 09906fa3 -f /var/lib/docker/containers/09906fa3/config.lxc -- /bin/bash -i
        container.cmd = exec.Command("/usr/bin/lxc-start", params...)
        ...
    }
    第一步是确保容器的文件系统已挂载。

第一步是确保容器的文件系统已挂载。

    func(fs *Filesystem)Mount()error {
        ...
        rwBranch := fmt.Sprintf("%v=rw", fs.RWPath)
        roBranches := ""
    for _, layer := range fs.Layers {
            roBranches += fmt.Sprintf("%v=ro:", layer)
        }
        branches := fmt.Sprintf("br:%v:%v", rwBranch, roBranches)

    // Mount the branches onto "/var/lib/docker/containers/09906fa3/rootfs"
        syscall.Mount("none", fs.RootFS, "aufs", 0, branches);
        ...
    }

使用 AUFS 联合挂载文件系统（union mount file system），镜像的各个层以 read-only 方式挂载在彼此的顶部，以向容器呈现一个一致的视图。read-write 路径被挂载在最顶层，为容器提供临时存储。

然后，为了启动容器，dockerd 使用刚才生成的 LXC 模板来运行另一个程序 lxc-start。

**LXC**

LXC（Linux Containers）是一个抽象层，它为用户空间应用程序提供了一个简单的 API 来创建和管理容器。事实是，容器不是真实的东西，在 Linux 内核中没有称作容器的对象。容器是一组内核对象的集合，它们协同工作以提供进程隔离。因此，简单的 lxc-start 命令实际上被翻译成下面的设置和应用：

- 内核名字空间 (ipc, uts, mount, pid, network 和 user)

- Apparmor 和 SELinux profiles

- Seccomp policies

- Chroots (使用 pivot_root)

- Kernel capabilities

- CGroups (control groups)

**容器清理**

    func(container *Container)monitor() {
    // Wait for the program to exit
        container.cmd.Wait()
        exitCode := container.cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

    // Cleanup
        container.stdout.Close()
        container.stderr.Close()
        container.Filesystem.Umount();

    // Report status back
        container.State.setStopped(exitCode)
        container.save()
    }

最后，dockerd 将监控容器直到完成，并在容器完成后清理不必要的数据。

**总结**

概括起来，使用 Docker 2013 来管理一个容器的运行涉及以下步骤：

    dockerd is sent the run command (using the command-line application or otherwise)
        ↳ dockerd finds the specified image on the file system
        ↳ A container struct is created and stored for future use
        ↳ Directories on the filesystemare setup foruseby the container
        ↳ LXC is instructed tostart the container
        ↳ dockerd monitors the containeruntil completion

**Docker 发生了什么变化？**

Docker 自从被引入以来已经过去 6 年了，容器化模式在此期间已经得到了迅猛发展。无论是大小企业都采用了 Docker，特别在它和容器编排系统 Kubernetes 的结合之后。

通过开源的力量，起初的 3 位贡献者已经发展到 1800 多人，每个人都为项目带来了新的想法。为了促进可扩展性，Open Container Initiative（OCI）标准于 2015 年发布，它定义了容器的镜像和运行时规范的开放标准。镜像规范概述了容器镜像的结构，运行时规范则描述了在其平台上运行容器时的实现应该遵循的接口和行为标准。因此，社区开发了广泛的容器管理项目，涵盖了从原生容器到被虚拟机隔离的容器。在微软的支持下，该行业现在也拥有了符合 OCI 标准的原生 Windows 容器。

所有这些变化都反映在 moby 代码库中。基于这种历史背景，我们可以开始分析 Docker 2019 及其组件了。

**Docker 2019 是如何工作的？**
经过 6 年和 36207 次的代码提交，moby 代码库已经发展成为一个大型的合作项目，它正在影响和依赖许多组件。![1.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf39dp/1616122191830-b3b387ff-df23-4ff2-90ef-4e833bdfddee.jpeg)

从一个非常简单的角度来看，Moby 2019 有两个新的主要组件，一个是在容器的生命周期中管理容器的 containerd 组件，另一个是符合 OCI 标准的运行时（例如 runc）组件，它是用于创建容器的最低用户级抽象（替换 LXC）。

**命令行应用程序**

命令行应用程序的控制流程大部分没有改变。今天，JSON 格式的 HTTP（S）报文是与 dockerd 通信的标准。

为了实现可扩展性，API 和 docker 二进制文件是分开的。程序代码位于 docker/cli 目录，它依赖 moby/moby/client 包作为接口与 dockerd 对话。

**dockerd**

dockerd 负责监听用户请求，并根据预先定义的路由对这些请求进行处理。

    // Container Routes
    func(r *containerRouter)initRoutes() {
        r.routes = []router.Route{
    // HEAD
            router.NewHeadRoute("/containers/{name:.*}/archive", r.headContainersArchive),
    // GET
            router.NewGetRoute("/containers/json", r.getContainersJSON),
            router.NewGetRoute("/containers/{name:.*}/export", r.getContainersExport),
            router.NewGetRoute("/containers/{name:.*}/changes", r.getContainersChanges),
            ...
        }
    } 

这个引擎仍然负责各种各样的任务，例如与镜像注册项的交互，并在文件系统上设置目录供容器使用。默认驱动程序会把一个镜像联合挂载到/var/lib/docker/overlay2/下的目录。

dockerd 不再负责管理运行容器的生命周期。随着项目的发展，决定将容器监管分开到一个名为 containerd 的单独项目。这样 docker 守护进程可以继续创新，而不必担心破坏运行时的实现。

尽管 docker/engine 是从 moby/moby 分支出来的，考虑到可能的代码差异，到目前为止，它们共享相同的代码提交树。

**docker run 命令**

    funcrunContainer(dockerCli command.Cli, opts *runOptions, copts *containerOptions, containerConfig *containerConfig)error {
        ...
    // Create the container
        createResponse = createContainer(ctx, dockerCli, containerConfig, &opts.createOptions)
        ...
    // Start the container
        client.ContainerStart(ctx, createResponse.ID, types.ContainerStartOptions{});
        ...
    }

docker run 命令首先请求守护进程创建一个容器。此请求被路由到 postContainersCreate 命令。

**创建容器**

稍后调用几个函数，我们来创建一个容器。

    func (daemon *Daemon) create(opts createOpts) (retC *container.Container, retErr error) {
        ...
    // Find the Image
        img = daemon.imageService.GetImage(params.Config.Image)
        ...
    // Create container object
        container = daemon.newContainer(opts.params.Name, os, opts.params.Config, opts.params.HostConfig, imgID, opts.managed);
        ...
    // Set RWLayer for container after mount labels have been set
        rwLayer = daemon.imageService.CreateLayer(container, setupInitLayer(daemon.idMapping))
        container.RWLayer = rwLayer
        ...
    // Create root directory 
        idtools.MkdirAndChown(container.Root, 0700, rootIDs);
        ...
    // Windows or Linux specific setup
        daemon.createContainerOSSpecificSettings(container, opts.params.Config, opts.params.HostConfig);
        ...
    // Store in a map for future lookup
        daemon.Register(container);
        ...
    }

首先，我们创建一个对象来存储容器的元数据。

然后像以前一样，我们创建一个根目录，其中包含镜像数据和读写（read-write）层，供容器使用。现在的区别在于联合挂载文件系统（union mount file system）已经成长到能够支持 btrfs, OverlayFS 和其它文件系统。为了方便这一点，驱动系统抽象了实现。

最后，容器对象被添加到守护进程的容器列表(map)中，以供将来使用。

**启动容器**

现在容器已创建成功，但尚未运行。接下来，我们请求启动它。

    func(daemon *Daemon)containerStart(container *container.Container, checkpoint string, checkpointDir string, resetRestartManager bool)(err error) {
        ...
    // Create OCI spec for container
        spec = daemon.createSpec(container);
        ...
    // Call containerd to create the container according to spec
        daemon.containerd.Create(ctx, container.ID, spec, createOptions)
        ...
    // Call containerd to start running process inside of container
        pid = daemon.containerd.Start(context.Background(), container.ID, checkpointDir,
            container.StreamConfig.Stdin() != nil || container.Config.Tty,
            container.InitializeStdio);
        ...
    }

这就是 containerd 起作用的地方，首先我们请求按照 OCI 规格（spec）来创建一个容器。然后，开始在这个容器内运行一个进程。然后把随后的监管交给 containerd 来处理。

**containerd**

containerd 这个术语有令人困惑的地方。它被认为是一个运行时（runtime），但是又不实现 OCI 运行时规范，因此它是一个和 runc 不一样的的运行时。containerd 是一个守护进程，它使用了符合 OCI 规格的运行时来监管容器的生命周期。正如 Michael Crosby 所描述的，containerd 是容器的监管者。
![2.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf39dp/1616122241290-e69d939f-fadf-4db9-851d-8ac7148cbed0.jpeg)
(Source: Docker)

containerd 被设计成监控容器的通用基础层，专注于速度和简单性。

很简单，创建一个容器所需要的只是它的规格（spec）和一个对根文件系统所在的位置进行编码的包。

**创建容器**

    func(l *local)Create(ctx context.Context, req *api.CreateContainerRequest, _ ...grpc.CallOption)(*api.CreateContainerResponse, error) {
        l.withStoreUpdate(ctx, func(ctx context.Context, store containers.Store)error {
    // Update container data store with new container record
            container := containerFromProto(&req.Container)
            created := store.Create(ctx, container)
            resp.Container = containerToProto(&created)
    returnnil
        });
        ...
    }

dockerd（通过一个 GRPC 客户端）请求 containerd 创建一个容器。containerd 接收这个请求后，将容器规格(spec)存储在位于/var/lib/containerd/下的文件系统支持的数据库中。

**启动容器**

    // Start create and start a task for the specified containerd id
    func(c *client)Start(ctx context.Context, id, checkpointDir string, withStdin bool, attachStdio libcontainerdtypes.StdioCallback)(int, error) {
        ctr, err := c.getContainer(ctx, id)
        ...
        t, err = ctr.NewTask(ctx, ...)
        ...
        t.Start(ctx);
        ...
    returnint(t.Pid()), nil
    }

启动一个容器涉及创建和启动一个称为 Task 的新对象，该对象代表着容器内的一个进程。

**创建 Task**

1
Plain Text

    func (l *local) Create(ctxcontext.Context, r *api.CreateTaskRequest, _ ...grpc.CallOption) (*api.CreateTaskResponse, error) {
    container := l.getContainer(ctx, r.ContainerID)
        ...
        rtime l.getRuntime(container.Runtime.Name)
        ...
        c = rtime.Create(ctx, r.ContainerID, opts)
        ...
    }

Task 的创建由底层容器的运行时负责。containerd 复用了 OCI 运行时，因此我们需要查找哪个运行时被用来创建这个 Task。第一个和默认的运行时是 runc。这个运行时的创建命令最终是运行一个外部进程 runc，但它是通过间接调用 docker-shim 进程来完成的。

如果 containerd 崩溃了，运行中的容器的信息将会丢失。为了防止这种情况发生，containerd 为每个容器创建一个称为垫片（shim）的管理进程。shim 进程将调用一个 OCI 运行时来创建和启动一个容器，然后执行监视容器的职责，以捕获退出代码并管理标准 IO。

在嵌套代码中，shim 将在执行 create 命令时使用 go-runc bindings 库来启动/run/containerd/runc 命令。。有关 runc 的更多信息，请参见下一节。

如果 containerd 真的崩溃了，可以通过与 shim 通信并从/var/run/containerd/目录读取状态信息来恢复。

**启动 Task**

既然容器已经创建了，启动 Task 需要做的事情就是简单地指示 shim 程序通过调用 runc start 命令来启动进程。

**Runc**

runc 是一个命令行工具，用于根据 OCI 规格生成和运行容器。当它执行与 LXC 类似的工作时，它抽象出创建容器所需的 Linux 内核调用。
![3.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf39dp/1616122274035-4eb4c0d9-e3ba-47fe-a049-2e7cee2b89c7.jpeg)

（这只可爱的仓鼠属于 runc）

runc 只是 OCI 运行时规格（spec）的一个实现，更多的可用于在各种系统上创建容器的实现可以在这里找到（<https://github.com/opencontainers/runtime-spec/blob/master/implementations.md>）。

**创建容器(runc create)**

    var createCommand = cli.Command{
        Name:  "create",
        Description: `The create command creates an instance of a container for a bundle. The bundle
    is a directory with a specification file named "` + specConfig + `" and a root
    filesystem.`,
        ...
        Action: func(context *cli.Context)error {
            spec := setupSpec(context)
            ...
            status := startContainer(context, spec, CT_ACT_CREATE, nil)
    // exit with the container's exit status so any external supervisor is
    // notified of the exit with the correct exit status.
            os.Exit(status)
    returnnil
    },

当 runc 创建一个容器时，它会在容器内设置名字空间、cgroups 甚至 init 进程。当创建结束时，进程暂停，等待信号开始运行。

**启动容器(runc start)**

    var startCommand = cli.Command{
        Name:  "start",
        Usage: "executes the user defined process in a created container",
        ...
        Action: func(context *cli.Context)error {
            container := getContainer(context)
            status := container.Status()
            ...
    switch status {
    case libcontainer.Created:
    return container.Exec()
    case libcontainer.Stopped:
    return errors.New("cannot start a container that has stopped")
    case libcontainer.Running:
    return errors.New("cannot start an already running container")
            }
        },
    }

最后，runc 向暂停的进程发送一个信号以开始容器的启动。

**直观总结**

概括起来，使用 Docker 2019 来管理一个容器的运行涉及以下步骤：

- 一个创建容器的 POST 请求发送给 dockerd;

- dockerd 查找被请求的镜像;

- 一个容器对象被创建，并存储起来以供将来使用

- 设置文件系统的目录结构供容器使用;

- 一个启动容器的 POST 请求发送给 dockerd;

- 为容器创建 OCI 规格（spec）;

- containerd 被调用来创建容器;

- containerd 将容器规格储存在数据库中;

- containerd 被调用来启动容器;

- containerd 为容器创建一个 Task

- Task 使用 shim 来调用 runc create 指令

- containerd 启动这个 Task

- Task 使用 shim 来调用 runc start 指令

- shim/containerd 继续监控容器直到任务完成

借助 containerd 的架构图，可以让我们直观地理解整个过程。
![p.jpg](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf39dp/1616122118090-96da77da-0b52-41c5-8942-9787f054a66a.jpeg)

**结论**

表面上看，Docker 及其配套项目显得杂乱无章，但实际上其底层的结构清晰并且已经实现了模块化。也就是说，发现所有上面的这些信息不是一件轻松的事。它散落在代码、博客帖子、会议讨论、文档和会议笔记中。拥有清晰的“自文档化”的代码是一个伟大的目标，但当涉及到大型系统时，我认为这还不够。有时，你只需要用简单的语言写下系统的外观，以及每个组件的职责。

非常感谢这些项目的所有贡献者，特别是那些编写解释这些系统的文档的人。

希望这篇文章这有助于帮助你理解 Docker 是如何运行容器的。我知道将来我会多次使用这篇文章作为参考。
