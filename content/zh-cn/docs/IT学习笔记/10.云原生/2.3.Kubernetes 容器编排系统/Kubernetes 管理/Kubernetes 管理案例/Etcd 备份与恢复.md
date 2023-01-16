---
title: Etcd 备份与恢复
---

系统环境：

- Etcd 版本：3.4.3
- Kubernetes 版本：1.18.8
- Kubernetes 安装方式：Kubeadm

# 备份 Etcd 数据

本人采用的是 Kubeadm 安装的 Kubernetes 集群，采用镜像方式部署的 Etcd，所以操作 Etcd 需要使用 Etcd 镜像提供的 Etcdctl 工具。如果是非镜像方式部署 Etcd，可以直接使用 Etcdctl 命令备份数据。

```bash
# 备份现有 Etcd 数据和manifests
mkdir -p /root/backup/kubernetes/
cp -r /var/lib/etcd/member /root/backup/kubernetes/member-$(date +%F)
cp -r /etc/kubernetes/manifests /root/backup/kubernetes/manifests-$(date +%F)

# 通过运行 Etcd 镜像，并且使用镜像内部的 etcdctl 工具连接 etcd 集群，执行数据快照备份：
docker run --rm --name etcdctl \
  -v /root/backup/kubernetes:/backup \
  -v /etc/kubernetes/pki/etcd:/etc/kubernetes/pki/etcd:ro \
  --env ETCDCTL_API=3 \
  registry.aliyuncs.com/k8sxio/etcd:3.4.13-0 \
  /bin/sh -c "etcdctl --endpoints=https://172.38.40.212:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --key=/etc/kubernetes/pki/etcd/healthcheck-client.key \
  --cert=/etc/kubernetes/pki/etcd/healthcheck-client.crt \
  snapshot save /backup/etcd.db-$(date +%F)"
```

# 恢复 Etcd 数据

注意：恢复数据前先停止 Kubernetes 相关组件！！防止恢复数据过程还会持续写入数据导致问题。然后进入 Etcd 镜像使用 etcdctl 工具执行恢复操作。

```bash
# 停止所有 master 节点上 k8s 系统组件
# 移除且备份 /etc/kubernetes/manifests 目录
$ mv /etc/kubernetes/manifests /root/backup
# 查看 kube-apiserver、etcd 镜像是否停止
$ docker ps|grep etcd && docker ps|grep kube-apiserver
# 清理已经损坏的 etcd 数据
rm -rf /var/lib/etcd/member
```

恢复 Etcd 数据
运行 Etcd 镜像，然后执行数据恢复，默认会恢复到 /default.etcd/member/ 目录下，这里使用 mv 命令在移动到挂载目录 /var/lib/etcd/ 下。

```bash
# 下面命令的 TIME 变量改成想要恢复的数据哪个时间的
# 在其中一个节点上恢复数据
docker run --rm              \
-v /root/backup:/backup        \
-v /var/lib/etcd:/var/lib/etcd \
--env ETCDCTL_API=3            \
registry.aliyuncs.com/k8sxio/etcd:3.4.3-0        \
/bin/sh -c "etcdctl snapshot restore /backup/etcd.db-${TIME}; mv /default.etcd/member/ /var/lib/etcd/"
# 在其他 master 节点上，将数据拷贝过来
scp -r 172.38.40.212:/var/lib/etcd/member /var/lib/etcd/
```

恢复 Kube-Apiserver 与 Etcd 镜像

```bash
# 将 /etc/kubernetes/manifests 目录恢复，使 Kubernetes 重启 Kube-Apiserver 与 Etcd 镜像：
cp -r /root/backup/manifests/ /etc/kubernetes/
```

恢复完成！！！

# 使用 k8s 的 CronJob 定时备份 etcd

注意：正式使用时修改 .spec.schedule 的值来改变备份周期(示例为每月 1 号执行备份)

基于 kubernetes-v1.23.2 编写

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: etcd-backup
  namespace: etcd-backup
spec:
  schedule: "0 0 1 * *"
  jobTemplate:
    spec:
      template:
        metadata:
          labels:
            name: etcd-backup
        spec:
          affinity:
            nodeAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                  - matchExpressions:
                      - key: node-role.kubernetes.io/master
                        operator: Exists
          containers:
            - command:
                - sh
                - -c
                - "etcdctl --endpoints=https://localhost:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt
                  --key=/etc/kubernetes/pki/etcd/healthcheck-client.key --cert=/etc/kubernetes/pki/etcd/healthcheck-client.crt
                  snapshot save /backup/etcd.db-$(printf '%(%Y-%m-%d-%H:%M:%S)T') "
              env:
                - name: ETCDCTL_API
                  value: "3"
              image: registry.aliyuncs.com/dd_k8s/etcd:3.5.1-0
              name: etcd
              volumeMounts:
                - mountPath: /backup
                  name: etcd-backup-db
                - mountPath: /etc/localtime
                  name: host-time
                  readOnly: true
                - mountPath: /etc/kubernetes/pki/etcd
                  name: etcd-certs
                  readOnly: true
          hostNetwork: true
          restartPolicy: Never
          tolerations:
            - effect: NoSchedule
              key: node-role.kubernetes.io/master
              operator: Exists
          volumes:
            - name: etcd-backup-db
              persistentVolumeClaim:
                claimName: etcd-backup
            - hostPath:
                path: /etc/localtime
              name: host-time
            - hostPath:
                path: /etc/kubernetes/pki/etcd
              name: etcd-certs
```
