---
title: "Ascend Plugin"
linkTitle: "Ascend Plugin"
created: "2026-05-08T20:08"
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，vllm-project/vllm-ascend](https://github.com/vllm-project/vllm-ascend)
> - [官方文档](https://docs.vllm.ai/projects/ascend/en/latest/index.html)

vllm-project/vllm-ascend 项目是 vLLM 的昇腾插件，让 vllm 可以在 [NPU](/docs/0.计算机/NPU.md) 设备上使用

支持的模型: https://docs.vllm.ai/projects/ascend/en/latest/user_guide/support_matrix/supported_models.html

# 部署

## 容器化部署

可用的镜像位置

- https://quay.io/repository/ascend/vllm-ascend
- https://quay.io/repository/ascend/cann

如果在国内，可以使用 `daocloud` 或其他镜像站点来加速下载：

```bash
TAG=v0.18.0
docker pull m.daocloud.io/quay.io/ascend/vllm-ascend:$TAG
# 或者
docker pull quay.nju.edu.cn/ascend/vllm-ascend:$TAG
```

---

```bash
export IMAGE=quay.io/ascend/vllm-ascend:v0.18.0
# 使用容器运行
docker run --rm \
--name vllm-ascend \
--shm-size=10g \
--device /dev/davinci0 \
--device /dev/davinci1 \
--device /dev/davinci2 \
--device /dev/davinci3 \
--device /dev/davinci4 \
--device /dev/davinci5 \
--device /dev/davinci6 \
--device /dev/davinci7 \
--device /dev/davinci_manager \
--device /dev/devmm_svm \
--device /dev/hisi_hdc \
-v /usr/local/dcmi:/usr/local/dcmi \
-v /usr/local/bin/npu-smi:/usr/local/bin/npu-smi \
-v /usr/local/Ascend/driver/lib64/:/usr/local/Ascend/driver/lib64/ \
-v /usr/local/Ascend/driver/version.info:/usr/local/Ascend/driver/version.info \
-v /etc/ascend_install.info:/etc/ascend_install.info \
-v /root/.cache:/root/.cache \
-p 8000:8000 \
-it $IMAGE bash
```

# 使用 vLLM

> 参考：
>
> - [官方文档，入门 - 使用](https://docs.vllm.ai/projects/ascend/zh-cn/v0.18.0/quick_start.html#usage)

> [!Attention] 昇腾插件比较特殊，由于某些底层硬件的缺陷，在使用 `vllm serve` 时需要添加额外的参数才能保证正常启动。
>
> 在这个章节，只记录一些特殊情况，还有一些情况记录在最佳实践中。一般来说，vLLM 使用模型还是比较简单的，只需要 `vllm serve ${Model}` 即可。

# 最佳实践

## 昇腾 Atlas 300I Duo 上的使用案例

> 由于使用的是 Atlas 300I Duo（芯片: 310p NPU），需要参考[这里](https://docs.vllm.ai/projects/ascend/zh-cn/latest/tutorials/hardwares/310p.html)，修改一些参数的默认值，才可以加载模型。

前提：使用 `HF_ENDPOINT=https://hf-mirror.com hf download Qwen/Qwen3-0.6B --local-dir /root/models/qwen3-0.6B` 将模型下载到

一、启动 vLLM

```bash
# Atlas 300 推理系列支持较晚，使用 rc 版本的。
export IMAGE=quay.io/ascend/vllm-ascend:v0.18.0rc1-310p-openeuler
docker run --rm \
--name vllm-ascend \
--network=host \
--shm-size=10g \
--device /dev/davinci0 \
--device /dev/davinci1 \
--device /dev/davinci2 \
--device /dev/davinci3 \
--device /dev/davinci_manager \
--device /dev/devmm_svm \
--device /dev/hisi_hdc \
-v /usr/local/dcmi:/usr/local/dcmi \
-v /usr/local/bin/npu-smi:/usr/local/bin/npu-smi \
-v /usr/local/Ascend/driver/lib64/:/usr/local/Ascend/driver/lib64/ \
-v /usr/local/Ascend/driver/version.info:/usr/local/Ascend/driver/version.info \
-v /etc/ascend_install.info:/etc/ascend_install.info \
-v /root/.cache:/root/.cache \
-v /root/models:/root/models \
-it $IMAGE bash
```

二、加载模型，提供推理服务

```bash
vllm serve /root/models/qwen3-32B \
    --enforce-eager \
    --dtype float16 \
    --served-model-name qwen3-32b \
    --enable-auto-tool-choice \
    --tool-call-parser hermes \
    --tensor-parallel-size 4
```

> Attention: `--enforce-eager` 和 `--dtype` 是在 Atlas 300I Duo 加速卡上的适配参数。否则无法加载模型。

三、验证一下推理服务是否可以响应

```bash
# 列出模型
curl -s http://localhost:8000/v1/models | python3 -m json.tool
# 测试文本补全
curl -s http://localhost:8000/v1/completions \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "你好",
    "max_completion_tokens": 64,
    "top_p": 0.95,
    "top_k": 50,
    "temperature": 0.6
  }' | jq .
```

四、进入对话模式开始对话

```bash
docker exec -it vllm-ascend vllm chat
```

由于实现了 WebAPI 的推理服务，其它程序可以通过 `http://localhost:8000/v1` 使用 OpenAI 兼容的接口获取推理结果。

五、验证成功后，可以直接使用下面的命令一键拉起命令，这样也能通过 docker logs 检查日志

```bash
export IMAGE=quay.io/ascend/vllm-ascend:v0.18.0rc1-310p-openeuler
docker run -d \
--name vllm-ascend \
--network=host \
--shm-size=10g \
--device /dev/davinci0 \
--device /dev/davinci1 \
--device /dev/davinci2 \
--device /dev/davinci3 \
--device /dev/davinci_manager \
--device /dev/devmm_svm \
--device /dev/hisi_hdc \
-v /usr/local/dcmi:/usr/local/dcmi \
-v /usr/local/bin/npu-smi:/usr/local/bin/npu-smi \
-v /usr/local/Ascend/driver/lib64/:/usr/local/Ascend/driver/lib64/ \
-v /usr/local/Ascend/driver/version.info:/usr/local/Ascend/driver/version.info \
-v /etc/ascend_install.info:/etc/ascend_install.info \
-v /root/.cache:/root/.cache \
-v /root/models:/root/models \
$IMAGE vllm serve /root/models/qwen3-32B \
    --served-model-name qwen3-32b \
    --enable-auto-tool-choice \
    --tool-call-parser hermes \
    --tensor-parallel-size 4 \
    --enforce-eager \
    --dtype float16
```

# 重大变化

[issue 7394](https://github.com/vllm-project/vllm-ascend/issues/7394) # 让 Qwen3.5 系列模型可以在 Atlas 300I Duo 上部署。创建于 2025-03-17
