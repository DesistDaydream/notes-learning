---
title: helm template 模板相关命令
---

# template # 在本地渲染 chart 模板，并展示输出

**helm template \[NAME] \[CHART] \[FLAGS]**

## FLAGS

1. -a, --api-versions stringArray     Kubernetes api versions used for Capabilities.APIVersions
2. --atomic                       if set, the installation process deletes the installation on failure. The --wait flag will be set automatically if --atomic is used
3. --ca-file string               verify certificates of HTTPS-enabled servers using this CA bundle
4. --cert-file string             identify HTTPS client using this SSL certificate file
5. --create-namespace             create the release namespace if not present
6. --dependency-update            run helm dependency update before installing the chart
7. --description string           add a custom description
8. --devel                        use development versions, too. Equivalent to version '>0.0.0-0'. If --version is set, this is ignored
9. --disable-openapi-validation   if set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema
10. --dry-run                      simulate an install
11. -g, --generate-name                generate the name (and omit the NAME parameter)
12. -h, --help                         help for template
13. --include-crds                 include CRDs in the templated output
14. --is-upgrade                   set .Release.IsUpgrade instead of .Release.IsInstall
15. --key-file string              identify HTTPS client using this SSL key file
16. --keyring string               location of public keys used for verification (default "/root/.gnupg/pubring.gpg")
17. --name-template string         specify template used to name the release
18. --no-hooks                     prevent hooks from running during install
19. **--output-dir**# 将输出的内容写入指定的文件中，默认是输出到标准输出
20. --password string              chart repository password where to locate the requested chart
21. --post-renderer postrenderer   the path to an executable to be used for post rendering. If it exists in $PATH, the binary will be used, otherwise it will try to look for the executable at the given path (default exec)
22. --release-name                 use release name in the output-dir path.
23. --render-subchart-notes        if set, render subchart notes along with the parent
24. --replace                      re-use the given name, only if that name is a deleted release which remains in the history. This is unsafe in production
25. --repo string                  chart repository url where to locate the requested chart
26. --set stringArray              set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
27. --set-file stringArray         set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
28. --set-string stringArray       set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
29. **-s, --show-only PATH/TO/TEMP**# 只显示指定的模板进行渲染出来的 manifest 文件
30. --skip-crds                    if set, no CRDs will be installed. By default, CRDs are installed if not already present
31. --timeout duration             time to wait for any individual Kubernetes operation (like Jobs for hooks) (default 5m0s)
32. --username string              chart repository username where to locate the requested chart
33. --validate                     validate your manifests against the Kubernetes cluster you are currently pointing at. This is the same validation performed on an install
34. -f, --values strings               specify values in a YAML file or a URL (can specify multiple)
35. --verify                       verify the package before installing it
36. --version string               specify the exact chart version to install. If this is not specified, the latest version is installed
37. --wait                         if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment, StatefulSet, or ReplicaSet are in a ready state before marking the release as successful. It will wait for as long as --timeout

## EXAMPLE

1. helm template rancher ./rancher-2.4.5.tgz --output-dir . #
2. helm template test ./platform/ --output-dir ./test #
