---
title: Mock API
---

API 文档虽然满足了阅读者对接接口的需求，但是在 API 起初设计的过程中，通常需要等待几天、甚至数周时间才能实时调用接口，从而真正开始接口调试工作。

而通过创建 Mock API，您即可提前交付模拟真实 API 响应的沙盒环境，以便 API 使用者提前开始调试工作，同时您也可以并行开发接口实现。

此外，您也可以在设计过程中得到 API 使用者对 API 设计的及时反馈，并进行迭代以得到更好的 API。

### 什么是 Mock API&#xA;

Mock API Server 通过提供真实 API 响应的范例数据来模拟真实的 API 服务，它将部署在 CODING 提供的公网服务器上，并且支持路由及参数校验，且在此基础上可限制访问的 IP 或 Token 授权方式。

> 目前 Mock API 支持静态模拟，即基于 API 响应范例的 Mock，暂不支持动态模拟（自定义数据模拟规则）。

### 功能特性

- 基于 API 文件数据生成 Mock API。您仅需完善 API 规则及范例，无需额外设置即可使用 Mock API。

- Mock API 支持参数校验，并可在请求有误时返回响应错误信息。大大降低 API 使用者在对接 API 时的错误率，方便调试及跟踪。

- 自定义 Mock API 响应数据。通过修改 API 范例数据，可让每个 API 的模拟数据趋于完美。

- 每篇文档 Mock API 均有独立域名，并支持 HTTPS / HTTP 双协议。方便统一配置，也方便记忆，地址参考：<http://c3wfvv32.mock.coding.iohttps://c3wfvv32.mock.coding.io>

- 支持 Token / IP 白名单授权方式。由于 Mock API 部署在公网，用户开启上述安全配置可有效防止他人随意调用 Mock API，掌握接口规则。

### 使用场景

1. API 对接/调试通常在公司项目中，API 使用者（如前端、App、自动化 API 测试开发人员）的开发进度会比后端 API 开发人员提前开始开发实现。而使用基于 API 文档的 Mock API 可提供模拟真实 API 响应的沙盒环境，以便 API 使用者提前开始调试工作。另一方面 API 使用者可以及时反馈 API 设计问题，在完成 API 实现之前提早完善 API 设计，使得 API 开发工作更加高效和趋于完美。

2. API 测试当您需要对部分 API 进行统一测试时，您可以替换其他 API 为 Mock API，而无需关心其他依赖 API 是否满足测试条件。

3. 外部 API 服务通常外部 API 服务可能会有不可靠、收费、访问限制等情况，您也可以替换外部 API 服务为 Mock API，通过 Mock 外部 API 服务的真实数据来调试程序逻辑。

### 使用工作流

1. 设计阶段

API 参数/规范设计

- OpenAPI：

  - 在线编辑 Swagger 文件：<http://editor.swagger.io>。

  - Spring 框架编写 API 结构注解。

- Postman：IDE 在线编辑。

- apiDoc：代码注释设计接口规范。

2. 生成范例数据

用于 Mock API 所需的 API 范例数据

- OpenAPI：Swagger Editor 可自动通过参数生成 Example。

- Postman：使用 Postman 内置 Mock Server，针对各 API 编写 Mock 数据规则，访问 Mock API 获取 Mock 数据。

- apiDoc：使用 <https://github.com/cdcabrera/apidoc-mock> 搭建 Mock Server，然后请求各 Mock API 获取 Mock 数据。

3. 生成 Mock API

通过 CODING API 文档导入 API 文件并发布，自动更新 Mock API 规则及数据。
