---
title: Rook+Ceph 现存问题
---

#

k8s 版本：1.18.8

# rook-ceph 的问题

## 使用 pvc 后出现

Failed to list \*v1beta1.VolumeSnapshotContent: the server could not find the requested resource (get volumesnapshotcontents.snapshot.storage.k8s.io)

解决方式：暂无

## 清理 rook-ceph 后出现

Operation for "/var/lib/kubelet/plugins/rbd.csi.ceph.com/csi.sock" failed."

问题跟踪：<https://github.com/rook/rook/issues/4359>

## Failed to list \*v1.PartialObjectMetadata: the server could not find the requested resource

问题跟踪：<https://github.com/kubernetes/kubernetes/issues/79610>

CRD 删除后，controller-manager 依然会持续跟踪，并报告错误

解决方式：

暂无
