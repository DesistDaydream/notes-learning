---
title: "CLI"
linkTitle: "CLI"
weight: 20
---

# 概述

> 参考：

你好，这是必应。我可以用中文回答你的问题。

Python 生态中有很多用于 CLI 的库，例如：

- **pyCLI**：一个用于创建简单、正确的命令行应用的框架，提供了标准的命令行解析、日志、测试和性能分析等功能。
  - 好像不更新了？
- **Python Fire**：一个用于自动从任何 Python 对象生成命令行接口的库，是一种在 Python 中创建 CLI 的简单方法。
- **argparse**：Python 标准库中的一个模块，用于快速地创建 CLI 应用，无需安装第三方库。
- **Click**：一个用于创建美观、可组合的命令行接口的 Python 包，代码量很少，配置灵活，提供了很多默认选项。

你可以根据你的需求和喜好选择合适的库。如果你想要更多的信息，你可以点击下面的链接查看相关文档或源码。

\-
\[1]: [cli — command line tools — pyCLI devel documentation](https://pythonhosted.org/pyCLI/)
\[2]: [GitHub - google/python-fire: Python Fire is a library for automatically ...](https://github.com/google/python-fire)
\[3]: [Build Command-Line Interfaces With Python's argparse](https://realpython.com/command-line-interfaces-python-argparse/)
\[4]: [click · PyPI](https://pypi.org/project/click/)
\[5]: [Welcome to Click — Click Documentation (8.1.x) - Pallets](https://click.palletsprojects.com/)

根据搜索结果，Go 的 cobra 库是一个用于创建强大的现代 CLI 应用的框架，提供了简单的接口和应用脚手架生成工具<sup>1</sup> <sup>2</sup> <sup>3</sup>。

与之类似的 Python 库可能是 **Click**，它也是一个用于创建美观、可组合的命令行接口的框架，提供了简单的接口和应用脚手架生成工具[^4^ [^5^。

你可以比较一下两个库的特性和文档，看看哪个更适合你的需求。

\-
\[1]: [Cobra. Dev](https://cobra.dev/)
\[2]: [GitHub - spf13/cobra: A Commander for modern Go CLI interactions](https://github.com/spf13/cobra)
\[3]: [Cobra documentation - GitHub Pages](https://umarcor.github.io/cobra/)

# Click

> 参考：
> 
> - [GitHub 项目，pallets/click](https://github.com/pallets/click)
> - [官网](https://click.palletsprojects.com/)

 https://juejin.cn/post/7114129474932703268

```python
import click

@click.command()
@click.option('--count', default=1, help='Number of greetings.')
@click.option('--name', prompt='Your name',
              help='The person to greet.')
def hello(count, name):
    """Simple program that greets NAME for a total of COUNT times."""
    for x in range(count):
        click.echo(f"Hello {name}!")

if __name__ == '__main__':
    hello()
```

```bash
~]# python hello.py --count=3
Your name: John
Hello John!
Hello John!
Hello John!
```

```bash
~]# python hello.py --help
Usage: hello.py [OPTIONS]

  Simple program that greets NAME for a total of COUNT times.

Options:
  --count INTEGER  Number of greetings.
  --name TEXT      The person to greet.
  --help           Show this message and exit.
```

