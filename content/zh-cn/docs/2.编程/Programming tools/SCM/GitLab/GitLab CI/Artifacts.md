---
title: Artifacts
linkTitle: Artifacts
date: 2024-07-18T10:24
weight: 20
---

# 概述

> 参考:
>
> - [官方文档，CI - jobs - job 工件](https://docs.gitlab.com/ee/ci/jobs/job_artifacts.html)

每个 Job 可以输出一些 包含文件和目录的 [Archive File(归档文件)](/docs/1.操作系统/Filesystem/Archive%20File(归档文件).md)、元数据、etc. ，这些输出称为 **Artifacts(工件)**。

> [!Warning]
> Artifacts 通常指我们构建的二进制文件、打包好的归档文件、etc. 。这些文件对项目来说，通常都是提供给使用者的，比如像 [GitHub](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub.md) 的 Release 中的 Assets。
>
> 但是 [<font color="#ff0000">GitLab 不建议将 Artifacts 的作为 Release Assets，因为 Artfacts 通常来说是短暂的，很有可能被轻易删除</font>](https://docs.gitlab.com/ee/user/project/releases/release_fields.html#use-a-generic-package-for-attaching-binaries)。因为对于 GitLab 来说，Release 中的 Assets 并不是被上传的文件，仅仅是一个名称和 URL 连接，指向其他地方。
>
> 更常见的作法是将 Artifacts 作为 Package 上传到 **[Package Registry](/docs/2.编程/Programming%20tools/SCM/GitLab/Packages%20AND%20Registries.md)**，让 Release Assets 设置为指向 Package 的 URL。参考: https://docs.gitlab.com/ee/user/packages/generic_packages/#publish-a-generic-package-by-using-cicd

创建 Artifact 的最简单示例:

```yaml
pdf:
  script: xelatex mycv.tex
  artifacts:
    paths:
      - mycv.pdf
```

在这个示例中，一个名为 pdf 的 JOB 调用 xelatex 命令从 LaTeX 源文件 mycv.tex 构建 PDF 文件，生成了名为 mycv.pdf 的文件，通过 paths 关键字指定了 mycv.pdf 文件作为 Artifact。

可以通过 WebUI 或者 [API](https://docs.gitlab.com/ee/api/job_artifacts.html#get-job-artifacts) 管理（查看、下载、删除、etch.）这些 Artifacts。

在 WebUI 中，点击项目左侧导航栏中 **Build > Arifacts**，查看该项目中所有 Job 产生的 Artifacts。包括如下几类：

- 在 [gitlab-ci.yml](/docs/2.编程/Programming%20tools/SCM/GitLab/GitLab%20CI/gitlab-ci.yml.md) 文件中通过 `artifacts` 关键字创建的 Artifacts
- [Report artifacts](#Report%20artifacts)
- Job 的日志和元数据作为单独的 Artifacts

# Report Artifacts

> 参考:
>
> - [官方文档，artifacts reports](https://docs.gitlab.com/ee/ci/yaml/artifacts_reports.html)

reports 功能可以收集一些复合 GitLab 内置标准模板的文件作为 Artifaces，比如 测试报告、代码质量报告、环境变量、etc.
