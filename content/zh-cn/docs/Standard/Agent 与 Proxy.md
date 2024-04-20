---
title: Agent 与 Proxy
---

今天，在阅读 Jolokia 文档的时候，看到其支持两种架构模式： Agent Mode 与 Proxy Mode。从字面上理解，agent 与 proxy 都有代理的意思，那么它们在内涵上到底有什么不同的，值得思考。

Agent Mode

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gu5xfp/1621662188213-22f409de-8912-4cb1-8722-b96d628694b6.png)

Proxy Mode

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gu5xfp/1621662192171-21f74acf-2333-4edd-859d-cc8aba728428.png)
从上面两张图可以看出，Proxy 与 Agent 所处的位置和目的有所不同，Agent 处在被代理对象的内部，而 Proxy 与被代理对象之间处于一种相对独立的状态。

举个例子，Proxy 类似于会计事务所，而 Agent 类似于公司里做会计工作的员工，它们都是处理企业的财务问题，但是角色不同。

> Agent 代理 通常指与用户接口的客户端程序
> 
> Proxy 代理 接收用户请求并将请求发给服务器,然后接收服务器回应并返回给用户 实际上 Proxy 的功能就是代替用户访问服务器,用户被隐藏.

## The Etymology of "Agent" and "Proxy" in Computer Networking Discourse

原文: https://cyber.harvard.edu/archived_content/people/reagle/etymology-agency-proxy-19981217.html

_September 18, 1998._ _[Joseph Reagle](https://link.zhihu.com/?target=mailto%3A//reagle%40mit.edu)
_
Revised: January 15, 1999 .

Given that the topic of this paper addresses both computer and legal agency, we believe an examination of the usage of the terms "agent" and "proxy" within each field is instructive. The technical use of the term \_[agent](https://link.zhihu.com/?target=http%3A//ai.miningco.com/)\_can be characterized into three overlapping periods. During the [earliest period](https://link.zhihu.com/?target=http%3A//ai.miningco.com/library/weekly/aa080397.htm) (1960 onward) the term was used by the artificial intelligence community. During the same period, the community of Internet network designers used the term in a way similar to our own: a service that acted on the behalf of another. Finally, the 1990's has been the decade of the agents: where the concept of potentially intelligent, autonomous computer programs that interact with each other within a networked community is hyped. It was also during this latter period that the term _proxy_ became widely used.

The technical concept of an agent was apparently first used in the artificial intelligence community in 1959 to describe the constituents (agents or demons) of a larger self-organizing behavior. \[RR, Self] However, within AI its usage did not always relate to computer programs, but to linguistic analysis, particularly in Fillmore's case grammar \[Fill]:

An agent is one who performs a particular action and this can be reflected in > _syntax_. In English, in an active sentence, the agent is usually referred to by the subject.... A form used to indicate the agent is called 'agentive.' In case grammar, agentive or agent is one of the specified sets of cases. \[Bear]

Agents were the subject of papers on intelligent information sharing systems and artificial intelligence throughout the 70 and 80s'. Occasionally, IT related papers used the term as any other discipline would, such as "Computers as an Agent of Change." \[Dieb]. Consequently, the use of the term to denote an initiator of action both within the AI/linguistic community and its typical English usage led the term to be used as a substitute of the term "program" when the connotation of autonomy and network interaction were required.

<> [networking](https://link.zhihu.com/?target=http%3A//wombat.doc.ic.ac.uk/foldoc/contents/networking.html)> In the > [client-server](https://link.zhihu.com/?target=http%3A//wombat.doc.ic.ac.uk/foldoc/foldoc.cgi%3Fclient-server) model, the part of the system that performs information preparation and exchange on behalf of a > [client](https://link.zhihu.com/?target=http%3A//wombat.doc.ic.ac.uk/foldoc/foldoc.cgi%3Fclient) or > [server](https://link.zhihu.com/?target=http%3A//wombat.doc.ic.ac.uk/foldoc/foldoc.cgi%3Fserver). Especially in the phrase "intelligent agent" it implies some kind of automatic process which can communicate with other agents to perform some collective task on behalf of one or more humans. \[> [Free On-line Dictionary of Computing](https://link.zhihu.com/?target=http%3A//wombat.doc.ic.ac.uk/foldoc/index.html)]

While computer agents were the subject of papers for nearly two decades, it was in the 1990's that the term became worthy of appearing in the title. Previously, the key words of related publications were information sharing, computer mediation, collaboration, and cooperation. \[Wood] Starting in 1994, specific venues for agent research and discussion, \[EWMA, CIKM] as well as a special issue of the Communications of the ACM placed the term prominently in popular discourse. \[Comm]

One of the best sources of information on how the terms _agents_ and _proxies_ came to be used in the context of the Internet and Web is the IETF's Request for Comments (RFC) archive. While the first occurrence of \_agent*in an RFC was administrative, \[RFC95] its technical debut was instructive in its use with respect to qualifying an agent as a program that resides on and communicates with other programs on a network.

The above premise is that the program (or agent) is doing the communicating with an NCP and thus needs to be identified for message traffic routing from an NCP. \[RFC129]

During the 70's at the IETF the only other time \_agent* was used in a technical sense was to refer to refer to the sender of an email. \[RFC733] However, in 1980, with the arrival of the Transmission Control Protocol (TCP) \[RFC761, RFC793] its usage was firmly established as a keyword in network protocol vernacular.

By 1987, the Internet TCP/IP protocol had matured sufficiently that the larger issue of network architecture and management became a focus of work at the IETF. Three critical RFCs related to this work include:

- [RFC1067](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1067.txt) -- Simple Network Management Protocol. \[RFC1067]
- [RFC1027](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1027.txt) -- Using ARP to implement transparent subnet gateways. \[RFC1027]
- [RFC1009](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1009.txt) -- Requirements for Internet gateways. \[RFC1009]

A concept integral to _proxy_ was that of a _gateway_. **Gateways** often provide an interface between two networks -- each of which may use a different networking protocol. A _proxy_ is often run on the gateway and acts as a go between, representing requests or services on behalf of one party in terms the second party can understand. The term proxy was seemingly first used in the network context by Shapiro \[Shap] in 1986 to designate one object as a local representative of a remote object. However, both of our terms, _agent_ and _proxy_, received their fullest treatment in 1989 in RFC 1095. _Proxy_ is explicitly placed in a legal context.

3.1. Architectural OverviewThe basic concepts underlying OSI network management are quite simple \[16]. There reside application processes called "managers" on managing systems (or management stations). There reside application processes called "agents" on managed systems (or network elements being managed). Network management occurs when managers and agents conspire (via protocols and a shared conceptual schema) to exchange monitoring and control information useful to the management of a network and its components. The terms "manager" and "agent" are also used in a loose and popular sense to refer to the managing and managed system, respectively.... The terms "manager" and "agent" are used to denote the asymmetric relationship between management application processes in which the manager plays the superior role and the agent plays the subordinate.4.3. Proxy ManagementProxy is a term that originated in the legal community to indicate an entity empowered to perform actions on behalf of another. In our context, a proxy is a manager empowered to perform actions on behalf of another manager. This may be necessary because the manager cannot communicate directly with the managed devices either for security or other administrative reasons or because of incompatible communication mechanisms or protocols. In either case, the proxy assumes the agent role with respect to the requesting manager and the manager role with respect to the managed device. \[RFC1090]

Subsequently, the terms agents and proxy have become integral to the nomenclature of networks and their applications. This includes email \[ref], ftp \[ref], the Web \[HTTP], and firewalls \[ref]. Professors \[Maes, Mins, Negr] and students associated with the [Agent Group](https://link.zhihu.com/?target=http%3A//agents.www.media.mit.edu/groups/agents/) at the [MIT Media Lab](https://link.zhihu.com/?target=http%3A//www.media.mit.edu/) have significantly contributed to the advancement of software agents in the popular scientific press. Others -- beyond the authors of this paper -- have begin to think about the social implications of autonomous network agents \[Fried1-3, Niss]

### References

\[Bear] Beardon C. (Editor) _[Artificial Intelligence Terminology : A Reference Guide (Ellis Horwood Series in Artificial Intelligence Foundations and Concepts)](https://link.zhihu.com/?target=http%3A//www.amazon.com/exec/obidos/ASIN/0130482994/qid%3D906405323/sr%3D1-8/002-1402323-9200248)_ (1989).
\[CIKM] Third International Conference on Information and Knowledge Management (CIKM'94)
\[Comm] Communications of the ACM. v.37 n.7, July1994.
\[Dieb] Diebold, J. _Man and the Computer; Technology as an Agent of Social Change._ F. A. Praeger, New York. (1969).
\[EWMA] European Workshop on Modeling Autonomous Agents in a Multi-Agent World (7th : 1996 : Eindhoven, Netherlands)
\[Fill] Fillmore, "_The Case for Case_." Universals in Linguistic Theory. (Bach, E. and Harms, R., eds.) Holt, Rinehart, and Winston, New York. (1968) pp. 1-90.
\[Fried1] Friedman B., and Millett L. _"It's the computer's fault" -- Reasoning about computers as moral agents_. Conference companion of the conference on Human Factors in Computing Systems, CHI '95. New York: Association for Computing Machinery. (May 1995) pp. 226- 227.
\[Fried2] Friedman B, and Nissenbaum H. [Bias in computer systems ](https://link.zhihu.com/?target=http%3A//www.acm.org/pubs/citations/journals/tois/1996-14-3/p330-friedman/). ACM Trans. Inf. Syst. 14, 3 (Jul. 1996) pp. 330 - 347.
\[Fried3] Friedman B, and Nissenbaum H [Software agents and user autonomy](https://link.zhihu.com/?target=http%3A//www.acm.org/pubs/citations/proceedings/ai/267658/p466-friedman/). Autonomous agents. (1997) pp. 466 - 469.
\[Maes] Maes, P. [Agents that Reduce Work and Information Overload.](https://link.zhihu.com/?target=http%3A//pattie.www.media.mit.edu/people/pattie/CACM-94/CACM-94.p1.html) Communications of the ACM. Vol. 37 No.7 (July 1994), pp. 31-40.
\[Mins] Minsky, M, and D. Riecken. [A Conversation with Marvin Minsky about Agents](https://link.zhihu.com/?target=http%3A//www.acm.org/pubs/citations/journals/cacm/1994-37-7/p22-minsky/). _[Communications of the ACM](https://link.zhihu.com/?target=http%3A//www.acm.org/cacm/)_ Vol. 37, No. 7 (July 1994) pp. 22-29.
\[Negr] Nicholas Negroponte. _Agents: From Direct Manipulation to Delegation_. Software Agents. (Jeffrey M. Bradshaw ed.), MIT Press 1997.
\[Niss] Nissenbaum, H. [Computing and accountability](https://link.zhihu.com/?target=http%3A//www.acm.org/pubs/citations/journals/cacm/1994-37-1/p72-nissenbaum/). Communications of the ACM. Vol. 37, No. 1 (Jan. 1994) pp. 72-80.
\[RR] Rome, B. and Rome, S. _Leviathan: A Simulation of Behavioral Systems, to Operate Dynamically on a Digital Computer_, System Development Corporation report no. SP-50, 6 (Nov. 1959), pp 15.
\[RFC1095] _[RFC1095](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1095.txt)_ _--_ _Common Management Information Services and Protocol over TCP/IP (CMOT)._ U.S. Warrier, L. Besaw. Apr-01-1989. (Obsoleted by RFC1189)
\[RFC1067] _[RFC1067](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1067.txt)_ _--_ _Simple Network Management Protocol_. J.D. Case, M. Fedor, M.L. Schoffstall, J. Davin. Aug-01-1988. (Obsoleted by RFC1098)
\[RFC1027] _[RFC1027](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1027.txt)_ _--_ Using ARP to implement transparent subnet gateways. S. Carl-Mitchell, J.S. Quarterman. Oct-01-1987.(Status: UNKNOWN)
\[RFC1009] _[RFC1009](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc1009.txt)_ _-- Requirements for Internet gateways_. R.T. Braden, J. Postel. Jun-01-1987. (Obsoletes RFC0985) (Obsoleted by RFC1812)
\[RFC733] _[RFC733](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc733.txt)_ _-- Standard for the format of ARPA network text messages_. D. Crocker, J. Vittal, K.T. Pogran, D.A. Henderson. Nov-21-1977. (Obsoletes RFC0724) (Obsoleted by RFC0822)
\[RFC761] _[RFC761](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc761.txt)_ _-- DoD standard Transmission Control Protocol_. J. Postel. Jan-01-1980.
\[RFC95] [RFC95](https://link.zhihu.com/?target=http%3A//info.internet.isi.edu/in-notes/rfc/files/rfc95.txt) -- _Distribution of NWG/RFC's through the NIC_. S.D. Crocker. Feb-04-1971. (Obsoleted by RFC0155)
\[Self] Selfridge, O. "Pandemonium: A Paradigm for Learning." _Mechanisation of Thought Processes_. London: H. M. Stationery Off., 1959, pp. 511-527.
\[Shap] Shapiro, M. _Structure and encapsulation in distributed systems: The proxy principle_. 6th International Conference on Distributed Computing Systems Proceedings (Cat. No. 86CH2293-9). (6th International Conference on Distributed Computing Systems Proceedings (Cat. No. 86CH2293-9), Cambridge, MA, USA, 19-23 May 1986). Washington, DC, USA: IEEE Comput. Soc. Press, 1986. pp. 198-204.
\[Wood] [Andy Wood](https://link.zhihu.com/?target=https%3A//cyber.harvard.edu/archived_content/people/reagle/amw%40cs.bham.ac.uk). [Agent Information and References](https://link.zhihu.com/?target=http%3A//www.cs.bham.ac.uk/~amw/agents/index.html). Available at [http://www.cs.bham.ac.uk/~amw/agen](https://link.zhihu.com/?target=http%3A//www.cs.bham.ac.uk/~amw/agents/index.html)
