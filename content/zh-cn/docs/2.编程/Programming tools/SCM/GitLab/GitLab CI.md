---
title: GitLab CI
linkTitle: GitLab CI
date: 2024-07-16T14:15
weight: 20
---

# 概述

> 参考:
>
> - [官方文档，主题 - 使用 CI/CD 构建你的应用](https://docs.gitlab.com/ee/topics/build_your_application.html)

在 WebUI 左侧导航栏 **Build** 标签中可以看到下面几个标签

- **Pipelines** # 查看根据 .gitlab-ci.yml 生成的流水线
- **Jobs** # 查看 Pipelines 中的所有 Jobs
- **Pipeline editor** # 在线编辑 .gitlab-ci.yml 文件
- **Pipeline schedules**
- **Artifacts** # 查看 [Artifacts](#Artifacts(工件))

# .gitlab-ci.yml

> 参考:
>
> - [官方文档，CI/CD YAML 语法参考](https://docs.gitlab.com/ee/ci/yaml/)

顶层字段

- **stages**(\[]STRING) # 定义要执行哪些阶段，也就是说要执行哪些 Job。`${JOB_NAME}.stage` 字段的值，就是被 stages 识别 JOB 的唯一标识符。
- **${JOB_NAME}**([Jobs](#Jobs)) # 
- **variables**(map\[STRING]STRING) # 定义全局变量，可以被所有 Job 使用。

# Jobs

https://docs.gitlab.com/ee/ci/yaml/#job-keywords

**artifacts**([artifacts](#artifacts)) # 与 [Artifacts(工件)](#Artifacts(工件)) 相关的配置

**image**(STRING) # 运行 Job 要使用的容器镜像。

**release**([release](#release)) # 与 Release 相关的配置。可以创建 Release

**stage**(STRING) # Job 的阶段。顶层字段 stages 所使用的就是该字段定义的值。`默认值: test`

**tags**(\[]STRING) # 运行 Job 要使用的 runner。只有具有相同 Tag 的 runner 才会被调度执行该 Job

## artifacts

**reports**([reports](#reports))

### reports

**dotenv**(STRING) # 符合 `VAR_NAME=VAR_VALUE` 格式的文件。文件中定义的变量可以传递给下游 JOB 作为环境变量使用。

## release


# TODO

内容多了之后，下面的概念相关说明移动到单独的文件中，注意修改 YAML 字段详解的引用。

# Artifacts(工件)

> 参考:
>
> - [官方文档，CI - jobs - job 工件](https://docs.gitlab.com/ee/ci/jobs/job_artifacts.html)

每个 Job 可以输出一些包含文件和目录的 [Archive File(归档文件)](docs/1.操作系统/Filesystem/Archive%20File(归档文件).md)，这些输出称为 **Artifacts(工件)**

> [!Warning]
> Artifacts 通常指我们构建的二进制文件、打包好的归档文件、etc. 。这些文件对项目来说，通常都是提供给使用者的，比如像 [GitHub](docs/2.编程/Programming%20tools/SCM/GitHub/GitHub.md) 的 Release 中的 Assets。
>
> 但是 [<font color="#ff0000">GitLab 不建议将 Artifacts 的作为 Release Assets，因为 Artfacts 通常来说是短暂的，很有可能被轻易删除</font>](https://docs.gitlab.com/ee/user/project/releases/release_fields.html#use-a-generic-package-for-attaching-binaries)。因为对于 GitLab 来说，Release 中的 Assets 并不是被上传的文件，仅仅是一个名称和 URL 连接，指向其他地方。
> 
> 更常见的作法是将 Artifacts 作为 Package 上传到 **[Package Registry](docs/2.编程/Programming%20tools/SCM/GitLab/Packages%20AND%20Registries.md)**，让 Release Assets 设置为指向 Package 的 URL。参考: https://docs.gitlab.com/ee/user/packages/generic_packages/#publish-a-generic-package-by-using-cicd

可以在每个项目的左侧导航栏中，点击 **Build > Arifacts** 查看该项目中所有 Job 产生的 Artifacts。包括如下几类：

- 在 .gitlab-ci.yml 文件中通过 `artifacts` 关键字创建的 Artifacts
- [Report artifacts](https://docs.gitlab.com/ee/ci/yaml/artifacts_reports.html)
- Job 的日志和元数据作为单独的 Artifacts

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