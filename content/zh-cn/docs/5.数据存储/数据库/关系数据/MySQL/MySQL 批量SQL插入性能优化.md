---
title: MySQL 批量SQL插入性能优化
---

对于一些数据量较大的系统，数据库面临的问题除了查询效率低下，还有就是数据入库时间长。特别像报表系统，每天花费在数据导入上的时间可能会长达几个小时或十几个小时之久。因此，优化数据库插入性能是很有意义的。

经过对 MySQL InnoDB 的一些性能测试，发现一些可以提高 insert 效率的方法，供大家参考参考。

## 1、一条 SQL 语句插入多条数据

常用的插入语句如：

```sql
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('0', 'userid_0', 'content_0', 0);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('1', 'userid_1', 'content_1', 1);
```

修改成：

```sql
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('0', 'userid_0', 'content_0', 0), ('1', 'userid_1', 'content_1', 1);
```

修改后的插入操作能够提高程序的插入效率。这里第二种 SQL 执行效率高的主要原因是合并后日志量（MySQL 的 binlog 和 innodb 的事务让日志）减少了，**降低日志刷盘的数据量和频率，从而提高效率。通过合并 SQL 语句，同时也能减少 SQL 语句解析的次数，减少网络传输的 IO**。

这里提供一些测试对比数据，分别是进行单条数据的导入与转化成一条 SQL 语句进行导入，分别测试 1 百、1 千、1 万条数据记录。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/osqpuu/1616133530263-c38130ab-121b-4b17-a60a-cfab979b341d.png)

## 2、在事务中进行插入处理。

把插入修改成：

```sql
START TRANSACTION;
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('0', 'userid_0', 'content_0', 0);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('1', 'userid_1', 'content_1', 1);
...
COMMIT;
```

使用事务可以提高数据的插入效率，这是因为进行一个 INSERT 操作时，MySQL 内部会建立一个事务，在事务内才进行真正插入处理操作。通过使用事务可以减少创建事务的消耗，`所有插入都在执行后才进行提交操作`。

这里也提供了测试对比，分别是不使用事务与使用事务在记录数为 1 百、1 千、1 万的情况。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/osqpuu/1616133530313-c27041f7-5ee1-40ad-b700-39240dcf2a81.png)

## 3、数据有序插入。

数据有序的插入是指插入记录在主键上是有序排列，例如 datetime 是记录的主键：

```sql
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('1', 'userid_1', 'content_1', 1);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('0', 'userid_0', 'content_0', 0);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('2', 'userid_2', 'content_2',2);
```

修改成：

```sql
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('0', 'userid_0', 'content_0', 0);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('1', 'userid_1', 'content_1', 1);
INSERT INTO `insert_table` (`datetime`, `uid`, `content`, `type`)
    VALUES ('2', 'userid_2', 'content_2',2);
```

由于数据库插入时，需要维护索引数据，`无序的记录会增大维护索引的成本`。我们可以参照 InnoDB 使用的 B+tree 索引，如果每次插入记录都在索引的最后面，索引的定位效率很高，并且对索引调整较小；如果插入的记录在索引中间，需要 B+tree 进行分裂合并等处理，会消耗比较多计算资源，并且插入记录的索引定位效率会下降，数据量较大时会有频繁的磁盘操作。

下面提供随机数据与顺序数据的性能对比，分别是记录为 1 百、1 千、1 万、10 万、100 万。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/osqpuu/1616133530291-b48ef56b-2dff-401b-8ede-bf3dc7e64b61.png)

从测试结果来看，该优化方法的性能有所提高，但是提高并不是很明显。

## 4、性能综合测试

这里提供了同时使用上面三种方法进行 INSERT 效率优化的测试。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/osqpuu/1616133530367-110c58c4-622e-41ff-8096-d1c640627053.png)

从测试结果可以看到，合并数据+事务的方法在较小数据量时，性能提高是很明显的，数据量较大时（1 千万以上），性能会急剧下降，这是由于此时数据量超过了 innodb_buffer 的容量，每次定位索引涉及较多的磁盘读写操作，性能下降较快。而使用合并数据+事务+有序数据的方式在数据量达到千万级以上表现依旧是良好，在数据量较大时，有序数据索引定位较为方便，不需要频繁对磁盘进行读写操作，所以可以维持较高的性能。

**注意事项：**

1. `SQL语句是有长度限制`，在进行数据合并在同一 SQL 中务必不能超过 SQL 长度限制，通过 max_allowed_packet 配置可以修改，默认是 1M，测试时修改为 8M。
2. `事务需要控制大小`，事务太大可能会影响执行的效率。MySQL 有 innodb_log_buffer_size 配置项，超过这个值会把 innodb 的数据刷到磁盘中，这时，效率会有所下降。所以比较好的做法是，在数据达到这个这个值前进行事务提交。

参考文档：MySQL 批量 SQL 插入性能优化
