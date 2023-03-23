---
title: 关于 Template 的其他说明
---

# 如何 Debug Templates

官方文档：[**https://helm.sh/docs/chart_template_guide/debugging/**](https://helm.sh/docs/chart_template_guide/debugging/)

调试模板可能会很棘手，因为渲染的模板已发送到 Kubernetes API 服务器，该服务器可能会出于格式化以外的其他原因而拒绝 YAML 文件。

有一些命令可以帮助您调试。

- `helm lint` 是验证图表是否遵循最佳做法的首选工具
- `helm install --dry-run --debug`或`helm template --debug`：我们已经看到了这个技巧。这是让服务器呈现模板，然后返回生成的清单文件的好方法。
- `helm get manifest`：这是查看服务器上安装了哪些模板的好方法。
- `helm template` ：用于调试模板渲染结果

当您的 YAML 无法解析，但您想查看生成的内容时，检索 YAML 的一种简单方法是在模板中注释掉问题部分，然后重新运行`helm install --dry-run --debug`：

```yaml
apiVersion: v2
# some: problem section
# {{ .Values.foo | quote }}
```

上面的内容将呈现并返回完整的注释：

```yaml
apiVersion: v2
# some: problem section
#  "bar"
```

这提供了一种查看生成的内容的快速方法，而不会阻止 YAML 分析错误。

# .helmignore 文件

[**https://helm.sh/docs/chart_template_guide/helm_ignore_file/**](https://helm.sh/docs/chart_template_guide/helm_ignore_file/)

# NOTES.txt 文件

[**https://helm.sh/docs/chart_template_guide/notes_files/**](https://helm.sh/docs/chart_template_guide/notes_files/)

# 其他

[**https://helm.sh/docs/chart_template_guide/wrapping_up/**](https://helm.sh/docs/chart_template_guide/wrapping_up/)

# 关于 YAML 与 Go 数据类型 和 Go 模板的说明

[**https://helm.sh/docs/chart_template_guide/yaml_techniques/**](https://helm.sh/docs/chart_template_guide/yaml_techniques/)

[**https://helm.sh/docs/chart_template_guide/data_types/**](https://helm.sh/docs/chart_template_guide/data_types/)
