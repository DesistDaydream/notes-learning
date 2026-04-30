---
title: "llama.cpp"
linkTitle: "llama.cpp"
created: "2026-04-30T13:15"
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，ggml-org/llama.cpp](https://github.com/ggml-org/llama.cpp)

LLaMA C++，其中 + 是 plus，简化为 2 个 p，所以是 llama.cpp。

llama.cpp 的主要目标是在本地和云端，以最小的设置和最先进的性能实现 LLM 推理。（尤其是消费级 CPU, GPU 和苹果 Mac）

大多数深度学习框架（e.g PyTorch）依赖 Python 和庞大的库，而 llama.cpp 使用纯 C/C++ 编写。

- 无依赖性： 不需要安装复杂的 Python 环境或庞大的 CUDA 驱动。
- 极致优化： 针对不同的 CPU 指令集（如 ARM 的 NEON、x86 的 AVX2）进行了底层优化。

llama.cpp 使用的模型格式为 **GPT-Generated Unified Format(GPT 生成的统一格式，简称 GGUF)**。这个 GGUF 文件中包含了所有必要的 元数据、分词器信息、模型权重。让 llama.cpp 使用一个单一的文件，即可进行推理。

> [!Tip] [Transformer](/docs/12.AI/Machine%20learning/Transformer/Transformer.md) 的模型通常包含 权重、配置、分词器配置、etc. 多个文件组成一个完整的 模型。

# 安装

从 [Release](https://github.com/ggml-org/llama.cpp/releases) 下载并使用即可，二进制文件直接用。

Windows 注意：有两个文件 [Windows x64 (CUDA 13)](https://github.com/ggml-org/llama.cpp/releases/download/b8840/llama-b8840-bin-win-cuda-13.1-x64.zip) 和 [CUDA 13.1 DLLs](https://github.com/ggml-org/llama.cpp/releases/download/b8840/cudart-llama-bin-win-cuda-13.1-x64.zip)，把 DLLs 文件解压到 llama.cpp 所在目录相同的目录下即可。

有时候我们的设备在运行从 Release 中下载的二进制文件时可能会出现下面这种报错，这是操作系统的 glibc / libstdc++ 版本太低的错误。为了防止升级 glibc 导致系统被破坏，我们可以选择使用源码 [编译安装](#编译安装)

```bash
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by ./llama-cli)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.30' not found (required by ./llama-cli)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by ./llama-cli)
./llama-cli: /usr/lib64/libm.so.6: version `GLIBC_2.38' not found (required by libllama-common.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.30' not found (required by libllama-common.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libllama-common.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `CXXABI_1.3.13' not found (required by libllama-common.so.0)
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by libllama-common.so.0)
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by libmtmd.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libmtmd.so.0)
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by libllama.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libllama.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `CXXABI_1.3.13' not found (required by libllama.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libggml-cpu.so.0)
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by libggml-cann.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libggml-cann.so.0)
./llama-cli: /usr/lib64/libstdc++.so.6: version `GLIBCXX_3.4.29' not found (required by libggml-base.so.0)
./llama-cli: /usr/lib64/libc.so.6: version `GLIBC_2.38' not found (required by libggml-base.so.0)
```

## 编译安装

```bash
git clone https://github.com/ggerganov/llama.cpp
cd llama.cpp
cmake -B build
cmake --build build --config Release -j$(nproc)
```

构建结果保存在当前目录的 `./build/bin/` 目录下

---

启用 [CANN](/docs/12.AI/Computing%20platform/CANN.md) 计算平台构建

https://github.com/ggml-org/llama.cpp/blob/master/docs/build.md#cann

```bash
cmake -B build -DGGML_CANN=on -DCMAKE_BUILD_TYPE=release
cmake --build build --config release -j$(nproc)
```

# 关联文件与配置

# 最佳实践

开始交互式对话

```powershell
.\llama-cli.exe -m D:\appdata\models\desistdayream.gguf
```

### 将模型包转换为 GGUF 文件

使用 llama.cpp 项目中的 convert_hf_to_gguf.py 文件将训练好的模型权重及相关文件，转换为单一的 GGUF 格式的文件

```powershell
python .\convert_hf_to_gguf.py `
  D:\appdata\models\desistdaydream\ `
  --outfile D:\appdata\models\desistdayream.gguf
```

```bash
python3 convert_hf_to_gguf.py \
  /root/models/qwen3.6 \
  --outfile /root/models/qwen3.6.gguf
```