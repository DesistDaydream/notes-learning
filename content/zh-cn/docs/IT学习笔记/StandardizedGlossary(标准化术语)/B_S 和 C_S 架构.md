---
title: B/S 和 C/S 架构
---

# Service Side 与 Client Side

> 参考：[Wiki，Server Side](https://en.wikipedia.org/wiki/Server-side)、[Wiki，Client Side](https://en.wikipedia.org/wiki/Client-side)

Client/Server 结构(C/S 结构) 是大家熟知的 **Client Side(客户端) 与 Server Side(服务端)** 结构。它是软件系统体系结构，通过它可以充分利用两端硬件环境的优势，将任务合理分配到 Client Side 和 Server Side 来实现，降低了系统的通讯开销。目前大多数应用软件系统都是 Client/Server 形式的两层结构，由于现在的软件应用系统正在向分布式的 Web 应用发展，Web 和 Client/Server 应用都可以进行同样的业务处理，应用不同的模块共享逻辑组件；因此，内部的和外部的用户都可以访问新的和现有的应用系统，通过现有应用系统中的逻辑可以扩展出新的应用系统。这也就是目前应用系统的发展方向。

B/S 结构（Browser/Server，浏览器/服务器模式），是 WEB 兴起后的一种网络结构模式，WEB 浏览器是客户端最主要的应用软件。这种模式统一了客户端，将系统功能实现的核心部分集中到服务器上，简化了系统的开发、维护和使用。客户机上只要安装一个浏览器（Browser 英 \['braʊzə]美 \['braʊzɚ]），如 Netscape Navigator 或 Internet Explorer，服务器安装 SQL Server、Oracle、MYSQL 等数据库。浏览器通过 Web Server 同数据库进行数据交互。

Apache 是普通服务器，本身只支持 html 即普通网页。不过可以通过插件支持 php,还可以与 Tomcat 连通(单向 Apache 连接 Tomcat,就是说通过 Apache 可以访问 Tomcat 资源。反之不然)。Apache 只支持静态网页，但像 php,cgi,jsp 等动态网页就需要 Tomcat 来处理。 Tomcat 是由 Apache 软件基金会下属的 Jakarta 项目开发的一个 Servlet 容器，按照 Sun Microsystems 提供的技术规范，实现了对 Servlet 和 JavaServer Page（JSP）的支持，并提供了作为 Web 服务器的一些特有功能，如 Tomcat 管理和控制平台、安全域管理和 Tomcat 阀等。由于 Tomcat 本身也内含了一个 HTTP 服务器，它也可以被视作一个单独的 Web 服务器。但是，不能将 Tomcat 和 Apache Web 服务器混淆，Apache Web Server 是一个用 C 语言实现的 HTTP web server；这两个 HTTP web server 不是捆绑在一起的。Apache Tomcat 包含了一个配置管理工具，也可以通过编辑 XML 格式的配置文件来进行配置。Apache，nginx，tomcat 并称为网页服务三剑客，可见其应用度之广泛。（说白了，tomcat 就是个底层设施软件服务，网页上所有的东西就要放在 tomcat 上，别人才能通过 tomcat 访问，tomcat 占用 80 端口）

例子：当通过浏览器，访问一个网站的时候，这时候就是一个 B/S 的架构，因为网站肯定是运行在服务器上的，这个服务器的系统上又装了 tomcat 这个服务，并占用 80 端口，因此，人们通过浏览器访问网站，由于自动使用 80 端口，那么就直接访问到 comcat 服务，然后由 comcat 来调取网站的页面资源给客户展示出来。

至于 C/S 就相当于通过 PC 端的一个软件，通过 URL 网址访问到服务器上的应用程序。
