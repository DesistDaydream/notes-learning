---
title: gitlab-ci.yml
linkTitle: gitlab-ci.yml
date: 2024-07-18T10:52
weight: 2
---

# 概述

> 参考:
>
> - [官方文档，CI/CD YAML 语法参考](https://docs.gitlab.com/ee/ci/yaml/)

GitLab 默认使用 .gitlab-ci.yml 作为 Pipeline 的配置文件。

顶层字段

- **default**(OBJECT) # 有些字段可以作为全局定义，以便让其效果作用在该 Pipeline 的所有 Jobs 中。
- **workflow**([workflow](#workflow)) # 控制运行 Pipeline 的时机
- **include**(\[][include](#include)) # 使用 include 在 CI/CD 配置中包含外部 YAML 文件。将文件 gitlab-ci.yml 成多个文件以提高可读性，或减少在多个位置重复相同的配置。
- **stages**(\[]STRING) # 定义要执行哪些阶段，也就是说要执行哪些 Job。`${JOB_NAME}.stage` 字段的值，就是被 stages 识别 JOB 的唯一标识符。
- **variables**(map\[STRING]STRING) # 定义全局变量，可以被所有 Job 使用。
- **${JOB_NAME}**([Jobs](#jobs)) #

# workflow

# include

# Jobs

https://docs.gitlab.com/ee/ci/yaml/#job-keywords

定义 Pipeline 中每个 Job 如何运行的

---

与 Job 是否运行以及如何运行的前置条件相关字段

**image**(STRING) # 运行 Job 要使用的容器镜像。TODO: 若不指定默认使用的时什么？好像是创建 runner 时可以指定默认使用的 image？

**stage**(STRING) # Job 的阶段。顶层字段 stages 所使用的就是该字段定义的值。`默认值: test`

**tags**(\[]STRING) # 运行 Job 要使用的 runner。只有具有相同 Tag 的 runner 才会被调度执行该 Job

**needs**(\[]OBJECT) # 指定运行 Job 需要满足的条件。

---

Job 要执行的具体内容的相关字段

**befor_script**(\[]STRING) # 运行 script 字段指定的脚本之前要运行的脚本

**script**(\[]STRING) # 执行 Job 时要运行的脚本

**release**([release](#release)) # 与 Release 相关的配置。可以创建 Release

---

其他

**artifacts**([artifacts](#artifacts)) # 定义如何处理在 Job 运行中产生的 [Artifacts(工件)](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab%20CI/Artifacts.md)

## artifacts

与 [Artifacts(工件)](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab%20CI/Artifacts.md) 相关的配置，用以指定在 Job 中要将哪些文件保存为 Artifact，以及如何管理这些 Artifacts。

**paths**(\[]STRING) # 指定要作为 Artifact 的文件或目录的路径。该字段是生成 Artifacts 的最基本配置。指定路径下的文件将会被打包成归档文件，归档文件名称使用 name 字段指定的值。

**expire_in**(STRING) # Artifacts 的过期时间。过期的 Artifacts 将会被删除。可用的格式有: `never、42、42 seconds、3 mins 4 sec、2 hrs 20 min、2h20min、6 mos 1 day、47 yrs 6 mos and 4d、3 weeks and 2 days`。`默认值: never`。若不指定单位只使用数字，则默认单位: 秒。

**name**(STRING) # 指定 [Archive File(归档文件)](/docs/1.操作系统/Filesystem/Archive%20File(归档文件).md) 类型的 Artifacts 的名字。`默认值: artifacts`，默认生成 artifacts.zip 文件。

**reports**([reports](#reports)) # 处理 [Report](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab%20CI/Artifacts.md#Report%20Artifacts) 类型的 Artifacts。用以指定一些额外的内容作为 Artifacts。

- **dotenv**(STRING) # 收集指定的环境变量文件作为 Artifact，该文件中定义的变量可以作为下游 JOB 的环境变量。文件格式为 `VAR_NAME=VAR_VALUE`，其中 VAR_NAME 可以在下游 JOB 的环境变量。

## release

创建 Release 时，需要使用 GitLab 研发的 release-cli 程序。可以使用 registry.gitlab.com/gitlab-org/release-cli:latest 容器，内部包含 release-cli 程序；也可以使用任意镜像手动安装 release-cli。

# Variables(变量)

> 参考:
>
> - [官方文档，CI - 变量](https://docs.gitlab.com/ee/ci/variables/)

在[这里](https://docs.gitlab.com/ee/ci/variables/predefined_variables.html)可以看到 GitLab 内部预定义的所有变量列表，都是以 CI 开都的，比如 CI_SERVER_HOST、CI_COMMIT_TAG、etc. 。

# Release(发布的版本)

最佳实践

- [将 Artifacts 上传到 Package registry，并在 Release Asset 中使用 Package 的 URL](https://docs.gitlab.com/ee/user/project/releases/release_fields.html#use-a-generic-package-for-attaching-binaries)

# 简单示例

```yaml
stages:
  - build
  - release

build-job:
  tags:
    - instance-runner-tc
  stage: build
  script:
    - bash observability/monitoring/packaging_server.sh
  artifacts:
    paths:
      - server.tar.gz

release-job:
  tags:
    - instance-runner-tc
  stage: release
  image: registry.gitlab.com/gitlab-org/release-cli:latest
  # 解决 release-cli 无法信任自建 CA 问题。https://gitlab.com/gitlab-org/release-cli/-/issues/47
  before_script:
    - sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
    - apk --no-cache add openssl ca-certificates
    - mkdir -p /usr/local/share/ca-certificates/extra
    - openssl s_client -connect ${CI_SERVER_HOST}:${CI_SERVER_PORT} -servername ${CI_SERVER_HOST} -showcerts </dev/null 2>/dev/null | sed -e '/-----BEGIN/,/-----END/!d' | tee "/usr/local/share/ca-certificates/${CI_SERVER_HOST}.crt" >/dev/null
    - update-ca-certificates
  script:
    - echo "检查 CI_COMMIT_TAG ${CI_COMMIT_TAG}"
    - ls -l
    - test -f server.tar.gz && echo "server.tar.gz exists" || echo "server.tar.gz does not exist"
  release:
    name: "Release $CI_COMMIT_TAG"
    # tag_name: $CI_COMMIT_TAG
    tag_name: "v1.0.0"
    description: "Release created by GitLab CI"
    assets:
      links:
        - name: "server.tar.gz"
          url: "${CI_PROJECT_URL}/-/jobs/artifacts/${CI_COMMIT_REF_NAME}/raw/server.tar.gz?job=build-job"
          link_type: "package"
  needs:
    - job: build-job
      artifacts: true
```
