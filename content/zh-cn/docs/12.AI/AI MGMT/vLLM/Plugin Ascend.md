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
>
> e.g. `--enforce-eager` 和 `--dtype float16` 这俩参数，就要在 Atlas 300I DUO 设备上用。

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

> [!Attention] 32B 的模型只用一块卡显存不够
> `--tensor-parallel-size 1` 的话会报错: `RuntimeError: NPU out of memory. Tried to allocate 502.00 MiB (NPU 0; 43.24 GiB total capacity; 41.70 GiB already allocated; 41.70 GiB current active; 349.23 MiB free; 41.72 GiB reserved in total by PyTorch) If reserved memory is >> allocated memory try setting max_split_size_mb to avoid fragmentation.`

```bash
vllm serve \
  --enforce-eager --dtype float16 \
  --model /root/models/qwen3-32B --served-model-name qwen3-32b \
  --enable-auto-tool-choice --tool-call-parser hermes \
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
$IMAGE vllm serve \
  --enforce-eager --dtype float16 \
  --model /root/models/qwen3-32B --served-model-name qwen3-32b \
  --enable-auto-tool-choice --tool-call-parser hermes \
  --tensor-parallel-size 4
```

# 重大变化

[issue 7394](https://github.com/vllm-project/vllm-ascend/issues/7394) # 让 Qwen3.5 系列模型可以在 Atlas 300I Duo 上部署。创建于 2025-03-17

# 基准测试

## Atlas 300I DUO

### 离线吞吐量基准测试

**qwen3-0.6B**

```bash
vllm bench throughput \
  --enforce-eager --dtype float16 \
  --model Qwen/Qwen3-0.6B \
  --input-len 256 --output-len 128 \
  --num-prompts 200 \
  --tensor-parallel-size 4
```

结果：

```bash
# 第一次
Throughput: 1.84 requests/s, 2124.55 total tokens/s, 236.06 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
# 第二次
Throughput: 2.13 requests/s, 2458.73 total tokens/s, 273.19 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
# 第三次
Throughput: 2.16 requests/s, 2489.63 total tokens/s, 276.63 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
```

--tensor-parallel-size=1 结果：

```bash
# 第一次
Throughput: 1.64 requests/s, 1887.75 total tokens/s, 209.75 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
# 第二次
Throughput: 1.64 requests/s, 1883.84 total tokens/s, 209.32 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
```

---

**qwen3-32B**

```bash
vllm bench throughput \
  --enforce-eager --dtype float16 \
  --model /root/models/qwen3-32B \
  --input-len 256 --output-len 128 --num-prompts 200 \
  --tensor-parallel-size 4
```

结果：

```bash
# 第一次
Throughput: 0.55 requests/s, 638.73 total tokens/s, 70.97 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
# 第二次
Throughput: 0.56 requests/s, 640.16 total tokens/s, 71.13 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
# 第三次
Throughput: 0.55 requests/s, 631.02 total tokens/s, 70.11 output tokens/s
Total num prompt tokens:  204800
Total num output tokens:  25600
```

### 在线基准测试

**qwen3-0.6B**

```bash
# 启动推理服务
vllm serve \
  --enforce-eager --dtype float16 \
  --model Qwen/Qwen3-0.6B \
  --tensor-parallel-size 4
# 对聊天接口基准测试
export bench_method="--num-prompts 20 --max-concurrency 1"
export model_id=$(curl -s http://localhost:8000/v1/models | python3 -c "import sys, json; print(json.load(sys.stdin)['data'][0]['id'])")
vllm bench serve \
  --model ${model_id} \
  --backend openai-chat --endpoint /v1/chat/completions \
  --dataset-name random \
  --random-input-len 256 --random-output-len 128 \
  ${bench_method}
```

结果：

<table>
<tr> 
<th>测试方法: --num-prompts 20 --max-concurrency 1</th>
<th>测试方法: --num-prompts 200 --max-concurrency 8</th> 
</tr>
<tr>
<td><pre>
============ Serving Benchmark Result ============
Successful requests:                     20        
Failed requests:                         0         
Maximum request concurrency:             1         
Benchmark duration (s):                  155.71    
Total input tokens:                      5120      
Total generated tokens:                  2560      
Request throughput (req/s):              0.13      
Output token throughput (tok/s):         16.44     
Peak output token throughput (tok/s):    18.00     
Peak concurrent requests:                2.00      
Total token throughput (tok/s):          49.32     
---------------Time to First Token----------------
Mean TTFT (ms):                          123.72    
Median TTFT (ms):                        88.58     
P99 TTFT (ms):                           656.70    
-----Time per Output Token (excl. 1st token)------
Mean TPOT (ms):                          60.33     
Median TPOT (ms):                        58.51     
P99 TPOT (ms):                           85.91     
---------------Inter-token Latency----------------
Mean ITL (ms):                           59.85     
Median ITL (ms):                         58.04     
P99 ITL (ms):                            149.20    
==================================================
</pre></td>
<td><pre>
============ Serving Benchmark Result ============
Successful requests:                     200       
Failed requests:                         0         
Maximum request concurrency:             8         
Benchmark duration (s):                  182.42    
Total input tokens:                      51200     
Total generated tokens:                  25600     
Request throughput (req/s):              1.10      
Output token throughput (tok/s):         140.34    
Peak output token throughput (tok/s):    160.00    
Peak concurrent requests:                16.00     
Total token throughput (tok/s):          421.02    
---------------Time to First Token----------------
Mean TTFT (ms):                          705.70    
Median TTFT (ms):                        661.14    
P99 TTFT (ms):                           1847.60   
-----Time per Output Token (excl. 1st token)------
Mean TPOT (ms):                          51.89     
Median TPOT (ms):                        52.02     
P99 TPOT (ms):                           52.91     
---------------Inter-token Latency----------------
Mean ITL (ms):                           51.89     
Median ITL (ms):                         52.17     
P99 ITL (ms):                            53.87     
==================================================
</pre></td>
</tr>
</table>

---

**qwen3-32B**

```bash
# 启动推理服务
vllm serve \
  --enforce-eager --dtype float16 \
  --model /root/models/qwen3-32B \
  --served-model-name qwen3-32b \
  --tensor-parallel-size 4
# 对聊天接口基准测试
export bench_method="--num-prompts 20 --max-concurrency 1"
export model_id=$(curl -s http://localhost:8000/v1/models | python3 -c "import sys, json; print(json.load(sys.stdin)['data'][0]['id'])")
vllm bench serve \
  --model ${model_id} \
  --tokenizer /root/models/qwen3-32B \
  --backend openai-chat --endpoint /v1/chat/completions \
  --dataset-name random \
  --random-input-len 256 --random-output-len 128 \
  ${bench_method}
```

> [!Attention] 由于 bench serve 的缺陷，需要手动添加 --tokenize 参数
> 原因是 vllm bench serve 在初始化 tokenizer 时，把 --model qwen3-32b 当成了 HuggingFace 的 repo_id ，本地没查到需要去下载，但 qwen3-32b 不是 namespace/name 格式，所以报错。

结果：

<table>
<tr> 
<th>测试方法: --num-prompts 20 --max-concurrency 1</th>
<th>测试方法: --num-prompts 200 --max-concurrency 8</th> 
</tr>
<tr>
<td><pre>
============ Serving Benchmark Result ============
Successful requests:                     20        
Failed requests:                         0         
Maximum request concurrency:             1         
Benchmark duration (s):                  334.72    
Total input tokens:                      5120      
Total generated tokens:                  2560      
Request throughput (req/s):              0.06      
Output token throughput (tok/s):         7.65      
Peak output token throughput (tok/s):    9.00      
Peak concurrent requests:                2.00      
Total token throughput (tok/s):          22.94     
---------------Time to First Token----------------
Mean TTFT (ms):                          328.02    
Median TTFT (ms):                        293.29    
P99 TTFT (ms):                           852.37    
-----Time per Output Token (excl. 1st token)------
Mean TPOT (ms):                          129.19    
Median TPOT (ms):                        127.83    
P99 TPOT (ms):                           157.93    
---------------Inter-token Latency----------------
Mean ITL (ms):                           128.18    
Median ITL (ms):                         127.89    
P99 ITL (ms):                            140.57    
==================================================
</pre></td>
<td><pre>
============ Serving Benchmark Result ============
Successful requests:                     200       
Failed requests:                         0         
Maximum request concurrency:             8         
Benchmark duration (s):                  541.15    
Total input tokens:                      51200     
Total generated tokens:                  25600     
Request throughput (req/s):              0.37      
Output token throughput (tok/s):         47.31     
Peak output token throughput (tok/s):    58.00     
Peak concurrent requests:                16.00     
Total token throughput (tok/s):          141.92    
---------------Time to First Token----------------
Mean TTFT (ms):                          2070.00   
Median TTFT (ms):                        2212.61   
P99 TTFT (ms):                           2743.95   
-----Time per Output Token (excl. 1st token)------
Mean TPOT (ms):                          154.09    
Median TPOT (ms):                        149.52    
P99 TPOT (ms):                           203.62    
---------------Inter-token Latency----------------
Mean ITL (ms):                           152.89    
Median ITL (ms):                         144.33    
P99 ITL (ms):                            450.19    
==================================================
</pre></td>
</tr>
</table>
