---
title: Debugger
weight: 8
---

# 概述

> 参考：
>
> - [官方文档，用户指南 - Debugging 任务](https://docs.ansible.com/ansible/latest/user_guide/playbooks_debugger.html)

Ansible 提供了一个任务[调试器](/docs/2.编程/Programming%20tools/Debugger.md)，因此您可以在执行过程中修复错误，而不是编辑 playbook 并再次运行它以查看更改是否有效。您可以在任务上下文中访问调试器的所有功能。您可以检查或设置变量的值，更新模块参数，并使用新的变量和参数重新运行任务。调试器可让您解决故障原因并继续执行 playbook。

### [Enabling the debugger with thedebuggerkeyword](https://docs.ansible.com/ansible/latest/user_guide/playbooks_debugger.html#id8)

You can use the `debugger` keyword to enable (or disable) the debugger for a specific play, role, block, or task. This option is especially useful when developing or extending playbooks, plays, and roles. You can enable the debugger on new or updated tasks. If they fail, you can fix the errors efficiently. The `debugger` keyword accepts five values:

| Value          | Result                                                |
| -------------- | ----------------------------------------------------- |
| always         | Always invoke the debugger, regardless of the outcome |
| never          | Never invoke the debugger, regardless of the outcome  |
| on_failed      | Only invoke the debugger if a task fails              |
| on_unreachable | Only invoke the debugger if a host is unreachable     |
| on_skipped     | Only invoke the debugger if the task is skipped       |
