---
title: Tree Data Structure(树型数据结构)
---

# 概述

> 参考：[Wiki-TreeDataStructure](<https://en.wikipedia.org/wiki/Tree_(data_structure)>)

在[计算机科学中](https://en.wikipedia.org/wiki/Computer_science)，**Tree** 是一种广泛使用的[抽象数据类型](https://en.wikipedia.org/wiki/Abstract_data_type)，它模拟分层[树结构](https://en.wikipedia.org/wiki/Tree_structure)，其根值和具有[父节点的](<https://en.wikipedia.org/wiki/Tree_(data_structure)#Terminology>)子级子树表示为一组链接[节点](<https://en.wikipedia.org/wiki/Node_(computer_science)>)。

可以将树数据结构[递归](https://en.wikipedia.org/wiki/Recursion)定义为节点的集合（从根节点开始），其中每个节点都是由值组成的数据结构，以及对节点（“子级”）的引用列表，其中约束，即没有重复的引用，也没有指向根的约束。或者，可以将树抽象为一个整体（全局地）定义为[有序树](https://en.wikipedia.org/wiki/Ordered_tree)，并为每个节点分配一个值。这两种观点都很有用：虽然一棵树可以作为一个整体进行数学分析，但是当实际上表示为数据结构时，它通常由节点表示和使用（而不是作为一组节点和节点之间的[邻接表）](https://en.wikipedia.org/wiki/Adjacency_list)，例如表示一个[有向图](<https://en.wikipedia.org/wiki/Tree_(data_structure)#Digraphs>)）。例如，从整体上看一棵树，可以谈论给定节点的“父节点”，但是通常，给定节点作为数据结构仅包含其子节点列表，但不包含引用。给它的父母（如果有的话）。

# 术语

1、结点(Node)：表示树中的数据元素，由数据项和数据元素之间的关系组成。在图 1 中，共有 10 个结点。
2、结点的度(Degree of Node)：结点所拥有的子树的个数，在图 1 中，结点 A 的度为 3。
3、树的度(Degree of Tree)：树中各结点度的最大值。在图 1 中，树的度为 3。
4、叶子结点(Leaf Node)：度为 0 的结点，也叫终端结点。在图 1 中，结点 E、F、G、H、I、J 都是叶子结点。
5、分支结点(Branch Node)：度不为 0 的结点，也叫非终端结点或内部结点。在图 1 中，结点 A、B、C、D 是分支结点。
6、孩子(Child)：结点子树的根。在图 1 中，结点 B、C、D 是结点 A 的孩子。
7、双亲(Parent)：结点的上层结点叫该结点的双亲。在图 1 中，结点 B、C、D 的双亲是结点 A。
8、祖先(Ancestor)：从根到该结点所经分支上的所有结点。在图 1 中，结点 E 的祖先是 A 和 B。
9、子孙(Descendant)：以某结点为根的子树中的任一结点。在图 1 中，除 A 之外的所有结点都是 A 的子孙。
10、兄弟(Brother)：同一双亲的孩子。在图 1 中，结点 B、C、D 互为兄弟。
11、结点的层次(Level of Node)：从根结点到树中某结点所经路径上的分支数称为该结点的层次。根结点的层次规定为 1，其余结点的层次等于其双亲结点的层次加 1。
12、堂兄弟(Sibling)：同一层的双亲不同的结点。在图 1 中，G 和 H 互为堂兄弟。
13、树的深度(Depth of Tree)：树中结点的最大层次数。在图 1 中，树的深度为 3。
14、无序树(Unordered Tree)：树中任意一个结点的各孩子结点之间的次序构成无关紧要的树。通常树指无序树。
15、有序树(Ordered Tree)：树中任意一个结点的各孩子结点有严格排列次序的树。二叉树是有序树，因为二叉树中每个孩子结点都确切定义为是该结点的左孩子结点还是右孩子结点。
16、森林(Forest)：m(m≥0)棵树的集合。自然界中的树和森林的概念差别很大，但在数据结构中树和森林的概念差别很小。从定义可知，一棵树有根结点和 m 个子树构成，若把树的根结点删除，则树变成了包含 m 棵树的森林。当然，根据定义，一棵树也可以称为森林。
