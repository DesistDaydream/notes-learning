---
title: Kubernetes并发控制和资源变更
---

# 概述

> 参考：
> - [公众号，云原生实验室，Kubernetes 是如何控制并发和资源变更的](https://mp.weixin.qq.com/s/pLmKnu-PY6hdcO7qmdvb5A)
>   - [原文，简书，Kubernetes 并发控制和资源变更](https://www.jianshu.com/p/ac830694a2cf)

## 并发控制

> 并发控制指的是当多个用户同时更新运行时，用于保护数据库完整性的各种技术。并发机制不正确可能导致脏读、幻读和不可重复读等此类问题。并发控制的目的是保证一个用户的工作不会对另一个用户的工作产生不合理的影响。

### 悲观锁

悲观锁在操作数据时比较悲观，认为别人会同时修改数据。因此操作数据时直接把数据锁住，直到操作完成后才会释放锁；上锁期间其他人不能修改数据。

悲观锁**主要用于数据争用激烈的环境**，以及发生并发冲突时使用锁保护数据的成本要低于回滚事务的成本的环境中。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/c5cda339-e90c-4fc0-9636-e84a6a8ad8fe/13618762-8d6f0d002e92278f.png)

**优点**

- 是“先取锁再访问”的保守策略，为数据处理的安全提供了保证。

**缺点**

- 在效率方面，处理加锁的机制会让数据库产生额外的开销，还有增加产生[死锁](https://links.jianshu.com/go?to=https%3A%2F%2Fzh.wikipedia.org%2Fwiki%2F%25E6%25AD%25BB%25E9%2594%2581)的机会；
- 在只读型事务处理中由于不会产生冲突，也没必要使用锁，这样做只能增加系统[负载](https://links.jianshu.com/go?to=https%3A%2F%2Fzh.wikipedia.org%2Fwiki%2F%25E8%25B4%259F%25E8%25BD%25BD)；
- 会降低了并行性，一个事务如果锁定了某行数据，其他事务就必须等待该事务处理完才可以处理那行数据。

### 乐观锁

乐观锁在操作数据时非常乐观，认为别人不会同时修改数据。因此乐观锁不会上锁，只是在执行更新的时候判断一下在此期间别人是否修改了数据：如果别人修改了数据则放弃操作，否则执行操作。

乐观并发控制多数**用于数据争用不大、冲突较少的环境中**，这种环境中，偶尔回滚事务的成本会低于读取数据时锁定数据的成本，因此可以获得比其他并发控制方法更高的吞吐量。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/c5cda339-e90c-4fc0-9636-e84a6a8ad8fe/13618762-64ac5f6be8446a30.png)

**优点**

- 不会产生任何锁和死锁
- 有更高的吞吐量

**缺点**

- [ABA 问题](https://links.jianshu.com/go?to=https%3A%2F%2Fzh.wikipedia.org%2Fwiki%2F%25E6%25AF%2594%25E8%25BE%2583%25E5%25B9%25B6%25E4%25BA%25A4%25E6%258D%25A2%23ABA%25E9%2597%25AE%25E9%25A2%2598)是乐观锁一个常见的问题
- 循环时间长开销大

**乐观锁一般会使用版本号机制或 CAS 算法实现：**

#### **版本号机制**

一般是在数据表中加上一个数据版本号 version 字段，表示数据被修改的次数，当数据被修改时，version 值会加一。当线程 A 要更新数据值时，在读取数据的同时也会读取 version 值，在提交更新时，若刚才读取到的 version 值为当前数据库中的 version 值相等时才更新，否则重试更新操作，直到更新成功。

#### **CAS 算法**

即**compare and swap（比较与交换）**，是一种有名的**无锁算法**。无锁编程，即不使用锁的情况下实现多线程之间的变量同步，也就是在没有线程被阻塞的情况下实现变量的同步，所以也叫非阻塞同步（Non-blocking Synchronization）。**CAS 算法**涉及到三个操作数

- 需要读写的内存值 V
- 进行比较的值 A
- 拟写入的新值 B

当且仅当 V 的值等于 A 时，CAS 通过原子方式用新值 B 来更新 V 的值，否则不会执行任何操作（比较和替换是一个原子操作）。一般情况下是一个**自旋操作**，即**不断的重试**。

## Kubernetes 并发控制

在 Kubernetes 集群中，外部用户及内部组件频繁的数据更新操作，导致系统的数据并发读写量非常大。假设采用悲观并行的控制方法，将严重损耗集群性能，因此 Kubernetes 采用乐观并行的控制方法。

### Resource Version

Kubernetes 通过定义资源版本字段实现了乐观并发控制，资源版本 (ResourceVersion)字段包含在 Kubernetes 对象的元数据 (Metadata)中。这个字符串格式的字段标识了对象的内部版本号。

通过 API Server 获取到的所有对象中，都有一个”resourceVersion”的字段。如：

    apiVersion: v1
    kind: Pod
    metadata:
      resourceVersion: "879232"
      selfLink: /api/v1/namespaces/default/pods/nginx-1zr5x
      uid: 9910eaf7-f0f3-11e7-a0b3-0800274a4ec3

该 Pod 的 resourceVersion 为 879232，更新该 Pod 时，Kubernetes 会比较该 resourceVersion 和 ETCD 中对象的 resourceVersion，在一致的情况下都会更新，一旦发生更新，该对象的 resourceVersion 值也会改变。

#### Resource Version 生成机制

下面的代码是 Kubernetes 从 ETCD 中获取对象的过程，我们可以从其中发现 Resource Version 的来源。

    func (s *store) Get(ctx context.Context, key string, resourceVersion string, out runtime.Object, ignoreNotFound bool) error {

      key = path.Join(s.pathPrefix, key)
        startTime := time.Now()

        getResp, err := s.client.KV.Get(ctx, key, s.getOps...)

        metrics.RecordEtcdRequestLatency("get", getTypeName(out), startTime)
        if err != nil {
            return err
        }

        if len(getResp.Kvs) == 0 {
            if ignoreNotFound {
                return runtime.SetZeroValue(out)
            }
            return storage.NewKeyNotFoundError(key, 0)
        }
        kv := getResp.Kvs[0]

        data, _, err := s.transformer.TransformFromStorage(kv.Value, authenticatedDataString(key))
        if err != nil {
            return storage.NewInternalError(err.Error())
        }

        return decode(s.codec, s.versioner, data, out, kv.ModRevision)
    }

从代码中我们可以看到，Resource Version 使用的是 ETCD 的 ModRevision。

ResourceVersion 字段在 Kubernetes 中除了用在上述并发控制机制外，还用在 Kubernetes 的 list-watch 机制中。Client 端的 list-watch 分为两个步骤，先 list 取回所有对象，再以增量的方式 watch 后续对象。Client 端在 list 取回所有对象后，将会把最新对象的 ResourceVersion 作为下一步 watch 操作的起点参数，也即 Kube-Apiserver 以收到的 ResourceVersion 为起始点返回后续数据，保证了 list-watch 中数据的连续性与完整性。

#### ETCD Version

ETCD 共四种\_version\_

- Revision
- ModRevision
- Version
- CreateRevision

关于他们的区别可以看下这个[issue：what is different about Revision, ModRevision and Version?](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fetcd-io%2Fetcd%2Fissues%2F6518)

> the Revision is the current revision of etcd. It is incremented every time the v3 backed is modified (e.g., Put, Delete, Txn). ModRevision is the etcd revision of the last update to a key. Version is the number of times the key has been modified since it was created. Get(..., WithRev(rev)) will perform a Get as if the etcd store is still at revision rev.

| 字段           | 作用范围 | 说明                               |
| -------------- | -------- | ---------------------------------- |
| Version        | Key      | 单个 Key 的修改次数，单调递增      |
| Revision       | 全局     | Key 在集群中的全局版本号，全局唯一 |
| ModRevison     | Key      | Key 最后一次修改时的 Revision      |
| CreateRevision | 全局     | Key 创建时的 Revision              |

《[Kubernetes 对象版本控制 ResourceVersion 和 Generation 原理分析](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.dianduidian.com%2Fpost%2Fkubernetes-resourceversion%25E5%258E%259F%25E7%2590%2586%25E5%2588%2586%25E6%259E%2590%2F)》中详细讲解了 Etcd Version 的变化过程。

### Generation

Generation 表示对象元配置信息（包括 spec 和 annotations）变更的次数。

    apiVersion: apps/v1
    kind: Deployment
    metadata:
      annotations:
        deployment.kubernetes.io/revision: "1"
      creationTimestamp: "2022-03-29T06:40:30Z"
      generation: 2

以 Deployment 为例，当某个 Deployment 对象被创建时，其 Generation 被设置为 1：

    func (deploymentStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
        deployment := obj.(*apps.Deployment)
        deployment.Status = apps.DeploymentStatus{}
        deployment.Generation = 1

        pod.DropDisabledTemplateFields(&deployment.Spec.Template, nil)
    }

每次当该 Deployment 对象的 spec 或 annotations 发生变化时，其 Generation + 1：

    func (deploymentStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
        newDeployment := obj.(*apps.Deployment)
        oldDeployment := old.(*apps.Deployment)
        newDeployment.Status = oldDeployment.Status

        pod.DropDisabledTemplateFields(&newDeployment.Spec.Template, &oldDeployment.Spec.Template)

        // Spec updates bump the generation so that we can distinguish between
        // scaling events and template changes, annotation updates bump the generation
        // because annotations are copied from deployments to their replica sets.
        // 当 spec 或 annotations 发生变化时，其 Generation + 1
        if !apiequality.Semantic.DeepEqual(newDeployment.Spec, oldDeployment.Spec) ||
            !apiequality.Semantic.DeepEqual(newDeployment.Annotations, oldDeployment.Annotations) {
            newDeployment.Generation = oldDeployment.Generation + 1
        }
    }

## 资源变更

### Create

Kubernetes 对象的创建流程如下：

1. 判断对象的 resourceVersion 是否合法，如果 resourceVersion != 0，则抛出错误
2. 对待处理对象做一些预处理：把 resourceVersion 和 selfLink 置为空
3. 对待处理对象进行编码，转换成二进制，进而转换成可被 ETCD 接受的格式
4. 判断 key 是否已存在，如果不存在，则存入 ETCD，否则返回错误信息
5. 记录执行耗时
6. 返回存储好的数据，并将 ETCD 中更新后的 Reversion 设置为 resourceVersion

<!---->

    func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {

        if version, err := s.versioner.ObjectResourceVersion(obj); err == nil && version != 0 {
            return errors.New("resourceVersion should not be set on objects to be created")
        }

        if err := s.versioner.PrepareObjectForStorage(obj); err != nil {
            return fmt.Errorf("PrepareObjectForStorage failed: %v", err)
        }

        data, err := runtime.Encode(s.codec, obj)
        if err != nil {
            return err
        }

        key = path.Join(s.pathPrefix, key)

        opts, err := s.ttlOpts(ctx, int64(ttl))
        if err != nil {
            return err
        }

        newData, err := s.transformer.TransformToStorage(data, authenticatedDataString(key))
        if err != nil {
            return storage.NewInternalError(err.Error())
        }

        startTime := time.Now()
        txnResp, err := s.client.KV.Txn(ctx).If(
            notFound(key),
        ).Then(
            clientv3.OpPut(key, string(newData), opts...),
        ).Commit()

        metrics.RecordEtcdRequestLatency("create", getTypeName(obj), startTime)
        if err != nil {
            return err
        }
        if !txnResp.Succeeded {
            return storage.NewKeyExistsError(key, 0)
        }

        if out != nil {

            putResp := txnResp.Responses[0].GetResponsePut()
            return decode(s.codec, s.versioner, data, out, putResp.Header.Revision)
        }
        return nil
    }

### Update

Kubernetes 实现了 Update 和 Patch 两个对象更新的方法，两者提供不同的更新操作方式，但冲突判断机制是相同的。

对于 Update，客户端更新请求中包含的是整个 obj 对象，服务器端将对比该请求中的 obj 对象和服务器端最新 obj 对象的 ResourceVersion 值。如果相等，则表明未发生冲突，将成功更新整个对象。反之若不相等则返回 409 冲突错误， Kube-Apiserver 中冲突判断的代码片段如下。
Kubernetes 对象的更新流程如下：

1. 获取当前更新请求中 obj 对象的 ResourceVersion 值，及服务器端最新 obj 对象 (existing) 的 ResourceVersion 值
2. 如果当前更新请求中 bj 对象的 ResourceVersion 值等于 0，即客户端未设置该值，则判断是否要硬改写 (AllowUnconditionalUpdate)，如配置为硬改写策略，将直接更新 obj 对象
3. 如果当前更新请求中 obj 对象的 ResourceVersion 值不等于 0，则判断两个 ResourceVersion 值是否一致，不一致返回冲突错误 (OptimisticLockErrorMsg)

\[图片上传失败...(image-67a94f-1648537253951)]

上图展示了多个用户同时 update 某一个资源对象时会发生的事情。而如果如果发生了 Conflict 冲突，对于 User A 而言应该做的就是做一次重试，再次获取到最新版本的对象，修改后重新提交 update，因此：

1. 用户修改 YAML 后提交 update 失败，是因为 YAML 文件中没有包含 resourceVersion 字段。对于 update 请求而言，应该取出当前 K8s 中的对象做修改后提交；
2. 如果两个用户同时对一个资源对象做 update，不管操作的是对象中同一个字段还是不同字段，都存在并发控制的机制确保两个用户的 update 请求不会发生覆盖。

Update 流程相关代码实现如下：

    func (s *store) GuaranteedUpdate(
        ctx context.Context, key string, out runtime.Object, ignoreNotFound bool,
        preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, suggestion ...runtime.Object) error {
        trace := utiltrace.New("GuaranteedUpdate etcd3", utiltrace.Field{"type", getTypeName(out)})
        defer trace.LogIfLong(500 * time.Millisecond)

        v, err := conversion.EnforcePtr(out)
        if err != nil {
            return fmt.Errorf("unable to convert output object to pointer: %v", err)
        }
        key = path.Join(s.pathPrefix, key)

        getCurrentState := func() (*objState, error) {
            startTime := time.Now()
            getResp, err := s.client.KV.Get(ctx, key, s.getOps...)
            metrics.RecordEtcdRequestLatency("get", getTypeName(out), startTime)
            if err != nil {
                return nil, err
            }
            return s.getState(getResp, key, v, ignoreNotFound)
        }

        var origState *objState
        var mustCheckData bool
        if len(suggestion) == 1 && suggestion[0] != nil {
            origState, err = s.getStateFromObject(suggestion[0])
            if err != nil {
                return err
            }
            mustCheckData = true
        } else {
            origState, err = getCurrentState()
            if err != nil {
                return err
            }
        }
        trace.Step("initial value restored")

        transformContext := authenticatedDataString(key)
        for {
            if err := preconditions.Check(key, origState.obj); err != nil {

                if !mustCheckData {
                    return err
                }



                origState, err = getCurrentState()
                if err != nil {
                    return err
                }
                mustCheckData = false

                continue
            }

            ret, ttl, err := s.updateState(origState, tryUpdate)
            if err != nil {

                if !mustCheckData {
                    return err
                }



                origState, err = getCurrentState()
                if err != nil {
                    return err
                }
                mustCheckData = false

                continue
            }

            data, err := runtime.Encode(s.codec, ret)
            if err != nil {
                return err
            }
            if !origState.stale && bytes.Equal(data, origState.data) {



                if mustCheckData {
                    origState, err = getCurrentState()
                    if err != nil {
                        return err
                    }
                    mustCheckData = false
                    if !bytes.Equal(data, origState.data) {

                        continue
                    }
                }

                if !origState.stale {
                    return decode(s.codec, s.versioner, origState.data, out, origState.rev)
                }
            }

            newData, err := s.transformer.TransformToStorage(data, transformContext)
            if err != nil {
                return storage.NewInternalError(err.Error())
            }

            opts, err := s.ttlOpts(ctx, int64(ttl))
            if err != nil {
                return err
            }
            trace.Step("Transaction prepared")

            startTime := time.Now()
            txnResp, err := s.client.KV.Txn(ctx).If(
                clientv3.Compare(clientv3.ModRevision(key), "=", origState.rev),
            ).Then(
                clientv3.OpPut(key, string(newData), opts...),
            ).Else(
                clientv3.OpGet(key),
            ).Commit()
            metrics.RecordEtcdRequestLatency("update", getTypeName(out), startTime)
            if err != nil {
                return err
            }
            trace.Step("Transaction committed")
            if !txnResp.Succeeded {
                getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
                klog.V(4).Infof("GuaranteedUpdate of %s failed because of a conflict, going to retry", key)
                origState, err = s.getState(getResp, key, v, ignoreNotFound)
                if err != nil {
                    return err
                }
                trace.Step("Retry value restored")
                mustCheckData = false
                continue
            }
            putResp := txnResp.Responses[0].GetResponsePut()

            return decode(s.codec, s.versioner, data, out, putResp.Header.Revision)
        }
    }

### Patch

相比 Update 请求包含整个 obj 对象，Patch 请求实现了更细粒度的对象更新操作，其请求中只包含需要更新的字段。例如要更新 pod 中 container 的镜像，可使用如下命令：

    kubectl patch pod my-pod -p '{"spec":{"containers":[{"name":"my-container","image":"new-image"}]}}'

服务器端只收到以上的 patch 信息，然后通过如下代码将该 patch 更新到 Etcd 中。

Kubernetes 对象的 Patch 更新流程如下：

1. 首先判断 patch 的类型，根据类型选择相应的 mechanism
2. 利用 DefaultUpdatedObjectInfo 方法将 applyPatch (应用 Patch 的方法)添加到 admission chain 的头部
3. 最终还是调用上述 Update 方法执行更新操作

<!---->

     func (p *patcher) patchResource(ctx context.Context, scope *RequestScope) (runtime.Object, bool, error) {
        p.namespace = request.NamespaceValue(ctx)
        switch p.patchType {
        case types.JSONPatchType, types.MergePatchType:
            p.mechanism = &jsonPatcher{
                patcher:      p,
                fieldManager: scope.FieldManager,
            }
        case types.StrategicMergePatchType:
            schemaReferenceObj, err := p.unsafeConvertor.ConvertToVersion(p.restPatcher.New(), p.kind.GroupVersion())
            if err != nil {
                return nil, false, err
            }
            p.mechanism = &smpPatcher{
                patcher:            p,
                schemaReferenceObj: schemaReferenceObj,
                fieldManager:       scope.FieldManager,
            }

        case types.ApplyPatchType:
            p.mechanism = &applyPatcher{
                fieldManager: scope.FieldManager,
                patch:        p.patchBytes,
                options:      p.options,
                creater:      p.creater,
                kind:         p.kind,
            }
            p.forceAllowCreate = true
        default:
            return nil, false, fmt.Errorf("%v: unimplemented patch type", p.patchType)
        }

        wasCreated := false
        p.updatedObjectInfo = rest.DefaultUpdatedObjectInfo(nil, p.applyPatch, p.applyAdmission)
        result, err := finishRequest(p.timeout, func() (runtime.Object, error) {

            options := patchToUpdateOptions(p.options)
            updateObject, created, updateErr := p.restPatcher.Update(ctx, p.name, p.updatedObjectInfo, p.createValidation, p.updateValidation, p.forceAllowCreate, options)
            wasCreated = created
            return updateObject, updateErr
        })
        return result, wasCreated, err
    }

相比 Update，Patch 的主要优势在于客户端**不必提供全量的 obj 对象信息**。客户端只需以 patch 的方式提交要修改的字段信息，服务器端会将该 patch 数据应用到最新获取的 obj 中。省略了 Client 端获取、修改再提交全量 obj 的步骤，降低了数据被修改的风险，更大大减小了冲突概率。 由于 Patch 方法在传输效率及冲突概率上都占有绝对优势，目前 Kubernetes 中几乎所有更新操作都采用了 Patch 方法，我们在编写代码时也应该注意使用 Patch 方法。

不过，patch 的复杂点在于，目前 K8s 提供了 4 种 patch 策略：json patch、merge patch、strategic merge patch、apply patch(server-side apply)。通过 kubectl patch -h 命令我们也可以看到这个策略选项（默认采用 strategic)

#### json patch

在 RFC6902 协议的定义中，JSON Patch 是执行在资源对象上的一系列操作，如下所示：

    {
        "op": "add",
        "path": "/spec/containers/0/image",
        "value": "busybox:latest"
    }

- op: 表示对资源对象的操作，主要有以下六种操作。
  1. add
  2. replace
  3. remove
  4. move
  5. copy
  6. test
- path: 表示被作资源对象的路径. 例如/spec/containers/0/image 表示要操作的对象是“spec.containers\[0].image”
- value: 表示预修改的值。

**新增容器**：

    kubectl patch deployment/foo --type='json' -p \
      '[{"op":"add","path":"/spec/template/spec/containers/1","value":{"name":"nginx","image":"nginx:alpine"}}]'

**修改已有的容器镜像**:

    kubectl patch deployment/foo --type='json' -p \
      '[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"app-image:v2"}]'

根据 http patch 原子性的定义，当某个 op(操作)不成功,则整个 patch 都不成功。

#### merge patch

merge patch 必须包含一个对资源对象的部分描述，json 对象。该 json 对象被提交到服务端，并和服务端的当前对象进行合并，从而创建新的对象。完整的替换列表，也就是说，新的列表定义会替换原有的定义。

例如(设置 label)：

    kubectl patch deployment/foo --type='merge' -p '{"metadata":{"labels":{"test-key":"foo"}}}'

使用 merge patch 也有如下限制：

- 如果 value 的值为 null,表示要删除对应的键，因此我们无法将 value 的值设置为 null, 如下，表示删除键 f

<!---->

    {
     "a":"z",
     "c": {
      "f": null
     }
    }

- merge patch 无法单独更新一个列表(数组)中的某个元素，因此不管我们是要在 containers 里新增容器、还是修改已有容器的 image、env 等字段，都要用整个 containers 列表(数组)来提交 patch：

<!---->

    kubectl patch deployment/foo --type='merge' -p \
      '{"spec":{"template":{"spec":{"containers":[{"name":"app","image":"app-image:v2"},{"name":"nginx","image":"nginx:alpline"}]}}}}'

#### strategic merge patch

参考《[kubernetes 中 update 与 patch 的区别](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.kubeclub.cn%2Fkubernetes%2F122.html)》

#### apply patch

参考《[kubernetes 中 update 与 patch 的区别](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.kubeclub.cn%2Fkubernetes%2F122.html)》

### Delete

Kubernetes 对象的删除流程如下：

1. 判断目标对象类型是否正确：是否为指针类型，是否不为 nil
2. 删除之前，先从 ETCD 中获取对应的数据，并判断该删除操作是否满足前置条件
3. 通过比对 ModVersion 判断这段时间内目标对象是否被其他进程/线程修改，如果未被修改，则执行删除操作；否则执行 Get 操作，删除失败，打印错误信息，并重新尝试删除
4. 删除成功，返回被删除的数据

```go
func (s *store) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions, validateDeletion storage.ValidateObjectFunc) error {

    v, err := conversion.EnforcePtr(out)
    if err != nil {
        return fmt.Errorf("unable to convert output object to pointer: %v", err)
    }

    key = path.Join(s.pathPrefix, key)
    return s.conditionalDelete(ctx, key, out, v, preconditions, validateDeletion)
}

func (s *store) conditionalDelete(ctx context.Context, key string, out runtime.Object, v reflect.Value, preconditions *storage.Preconditions, validateDeletion storage.ValidateObjectFunc) error {
    startTime := time.Now()

    getResp, err := s.client.KV.Get(ctx, key)
    metrics.RecordEtcdRequestLatency("get", getTypeName(out), startTime)
    if err != nil {
        return err
    }
    for {

        origState, err := s.getState(getResp, key, v, false)
        if err != nil {
            return err
        }
        if preconditions != nil {
            if err := preconditions.Check(key, origState.obj); err != nil {
                return err
            }
        }

        if err := validateDeletion(ctx, origState.obj); err != nil {
            return err
        }
        startTime := time.Now()



        txnResp, err := s.client.KV.Txn(ctx).If(
            clientv3.Compare(clientv3.ModRevision(key), "=", origState.rev),
        ).Then(
            clientv3.OpDelete(key),
        ).Else(
            clientv3.OpGet(key),
        ).Commit()
        metrics.RecordEtcdRequestLatency("delete", getTypeName(out), startTime)
        if err != nil {
            return err
        }
        if !txnResp.Succeeded {

            getResp = (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
            klog.V(4).Infof("deletion of %s failed because of a conflict, going to retry", key)
            continue
        }

        return decode(s.codec, s.versioner, origState.data, out, origState.rev)
    }
}
```

## 参考文档

[Kubernetes 并发控制与数据一致性的实现原理](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.cnblogs.com%2Fhuaweiyuncce%2Fp%2F10001148.html)

[kubernetes 中 update 与 patch 的区别](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.kubeclub.cn%2Fkubernetes%2F122.html)

[Kubernetes 对象版本控制 ResourceVersion 和 Generation 原理分析](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.dianduidian.com%2Fpost%2Fkubernetes-resourceversion%25E5%258E%259F%25E7%2590%2586%25E5%2588%2586%25E6%259E%2590%2F)

###
