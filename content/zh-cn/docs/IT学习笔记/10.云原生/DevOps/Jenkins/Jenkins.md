---
title: Jenkins
---

# 概述

> 参考：
> - 官方文档：<https://jenkins.io/zh/>

Jenkins 是一款开源 CI\&CD 软件，用于自动化各种任务，包括构建、测试和部署软件。Jenkins 支持各种运行方式，可通过系统包、Docker 或者通过一个独立的 Java 程序。

Jenkins 自己本身并没有什么功能，仅提供一个基本的 web 页面和监听端口(用来接收 webhook 等)，CI/CD 功能的实现依赖于 Jenkins 的各种插件来实现。

Jenkins 说白了就是一个运行在系统上的程序，通过编写各种 Jenkins Pipelines(流水线)，来在系统中进行代码构建、拷贝文件、执行 shell 命令等等。使用的是系统上环境，比如系统有 go 编译器，那么 Jenkins 就可以执行 go 的相关命令；如果系统上没有安装 docker ，那么 Jenkins 也就没法执行 docker 相关的命令。可以把 Jenkins 当作一个复杂的、可以执行更多功能的 shell 脚本。

当然，既然系统上无法满足各类要求的话，那么 Pipeline 还可以通过使用容器的方式来运行，以便规避环境影响。

## Pipe-line(流水线) 介绍

官方文档：<https://www.jenkins.io/doc/book/pipeline/>

Pipe-line 是 Jenkins 实现 CI/CD 的必备程序。在部署完 Jenkins 之后，一般情况都会首先安装 Pipe-line 插件，而日常使用中，最常用的也是 Pipeline。

顾名思义，Pipeline 流水线，就好像工厂中的流水线一样，第一步应该做什么，第二步应该做什么......都是有明确规定的，这是一个自动化，明确的任务线。

### Jenkinsfile

Pipeline 通过一个文本文件(称为 Jenkinsfile)来决定其每一步应该执行什么样的操作(比如第一步下载代码、第二步构建、第三步部署...等等)。这个 Jenkinsfile 文件可以提交到项目的代码存储库中，这是 Pipeline as code(流水线即代码) 的基础。

Jenskinsfile 的内容可以直接通过 Jenkins 的 web 页面编写，也可以放在项目的根目录，来让 Jenkins 自动读取该文件。将 Jenkinsfile 提交到代码仓库有以下好处

1. 为所有分支和拉取请求自动创建管道构建过程。
2. 管道上的代码审查/迭代（以及其余的源代码）。
3. 管道的审计跟踪。
4. 管道的单一事实来源 \[ 3 ]，可以由项目的多个成员查看和编辑。

## Jenkinsfile 基本示例

Jenkinsfile 有自己的一套语法，一共有两种分类，声明式语法和脚本式语法。一般常用声明式语法。

    pipeline {
        agent any
        options {
            skipStagesAfterUnstable()
        }
        stages {
            stage('Build') {
                steps {
                    sh 'make'
                }
            }
            stage('Test'){
                steps {
                    sh 'make check'
                    junit 'reports/**/*.xml'
                }
            }
            stage('Deploy') {
                steps {
                    sh 'make publish'
                }
            }
        }
    }

1. pipeline 是特定于声明式管道的语法，该语法定义了一个“块”，其中包含用于执行整个管道的所有内容和指令
2. agent 是声明式管道特定的语法，指示 Jenkins 为整个管道分配执行器（在节点上）和工作空间。可以指定代理为 docker，那么会创建一个容器，并在该容器中执行后续所有步骤。
3. stage 是描述 此 Pipeline 阶段的 语法块 。了解更多关于 stage 在在声明管道语法块管道语法页。如所提到的上述，stage 块在脚本管道语法可选的。
4. steps 是声明性管道特定的语法，描述了在 this 中要运行的步骤 stage。
5. sh 是执行给定 shell 命令的 Pipeline 步骤（由 Pipeline：Nodes and Processes 插件提供）。
6. junit 是另一个管道步骤（由 JUnit 插件提供），用于汇总测试报告

## 名词介绍

# Jenkins 部署

## 通过 WAR 文件安装 Jenkins

通过该链接下载最新的 war 文件

运行命令 java -jar jenkins.war 即可

jenkins 启动后，会自动创建 /root/.jenkins/ 目录作为自己的存储数据的路径($JENKINS_HOME)，并将 $JENKINS_HOME 添加到环境变量中。默认监听在 8080 端口，通过浏览器访问即可

## 通过 docker 运行 jenkins

    docker volume create jenkins-data
    docker run \
      -u root \
      -d \
      --name jenkins \
      -p 8080:8080 \
      -p 50000:50000 \
      -v jenkins-data:/var/jenkins_home \
      -v /var/run/docker.sock:/var/run/docker.sock \
      jenkinsci/blueocean

Note：

1. Jenkins 运行在 docker 环境中时，是无法获取宿主机的相关信息的，比如宿主机上 $PATH 路径下的命令文件、宿主机的环境变量等等信息，都是无法获取的，在执行 shell 流程时，需要注意该情况。
2. -v /var/run/docker.sock:/var/run/docker.sock 参数就是为了解决上述问题所使用的其中一个解决办法，将宿主机 docker 的 socker 挂载进容器中，那么通过容器启动的 Jenkins 在对 docker 操作时，就可以使用宿主机的 docker 命令，获取宿主机 docker 的镜像了。

# Jenkins 初始化

Jenkins 启动后，在访问 Jenkins 的时候，会触发初始化配置内容，页面如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cfhrw6/1616077873079-89b54590-c498-4fd0-89a6-c1d10d87373d.jpeg)

并在后台日志中输入密码

                ......*************************************************************Jenkins initial setup is required. An admin user has been created and a password generated.Please use the following password to proceed to installation:3c80b98b9aef457a87c7d202f6bd4dd7# 这个密码也会存在 Jenkins 数据路径的 /secrets/initialAdminPassword 文件中。This may also be found at: /var/jenkins_home/secrets/initialAdminPassword*************************************************************......

输入密码后，可以选择是根据向导继续初始化，还是点击下图右上角红框中的 X 来跳过向导。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cfhrw6/1616077873043-73556142-9cba-4495-a70c-bd637acea581.jpeg)

Note:

1. 如果跳过了该向导中创建用户的步骤，则日志中输出的密码可以作为 admin 用户的密码使用。如果创建了用户，则 initialAdminPassword 文件会被自动删除。因为新创建的用户就是具有最高权限的 admin 用户。
2. 如果输入密码后没有开始初始化向导页面。有可能是因为 Jenkins 的升级中心网站被墙的原因导致的，可以修改 $JENKINS_HOME/hudson.model.UpdateCenter.xml 文件中的 字段，将 Jenkins 的升级中心网址改为其他即可


    <?xml version='1.0' encoding='UTF-8'?>
     <sites>
       <site>
        <id>default</id>
        <url>http://updates.jenkins-ci.org/update-center.json</url>
       </site>
    </sites>

- 可用的 URL 有如下这些
- <http://mirror.xmission.com/jenkins/updates/update-center.json>
- <https://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/update-center.json>
- <http://updates.jenkins-ci.org/download/plugins/>
- <https://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/current/update-center.json>
- <http://mirror.esuni.jp/jenkins/updates/update-center.json>
- Localization: Chinese (Simplified) # Jenkins 中文插件

# Jenkins 配置

根据不同的安装情况，Jenkins 的配置文件存储路径都不同，不过 Jenkins 启动后，都会生成一个路径来保存自己的配置信心、插件、数据文件等等内容。可以将该目录称为 Jenkins 的数据目录。后文用 ${JENKINS_HOME} 来表示 Jenkins 的数据存储目录。

${JENKINS_HOME} # Jenkins 数据存储目录

1. ./jobs/\* #
2. ./workspace/\* #
