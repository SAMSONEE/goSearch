# SearchEngine
`SearchEngine`一个基于[悟空数据集](https://wukong-dataset.github.io/wukong-dataset/benchmark.html)的图文搜索引擎。

 - [基于中文词典树的高效索引和搜索](https://github.com/MingweiGuo/Goland/tree/main/SearchEngine/trie)
 - 支持词典树的[持久化](https://github.com/MingweiGuo/Goland/blob/main/SearchEngine/trie/serialize.go)
 - 支持中文分词（使用[sego分词包](https://github.com/huichen/sego)进行并发分词）
 - 支持[BM25算法](https://github.com/MingweiGuo/Goland/tree/main/SearchEngine/rank)
 - 支持搜索数据持久化（使用[leveldb](https://github.com/google/leveldb)进行索引存储）
 - 支持用户登录注册搜索数据（使用[Gin框架](https://github.com/gin-gonic/gin)实现）
 - 基于快排和二分法实现搜索结果[排序](https://github.com/MingweiGuo/Goland/blob/main/SearchEngine/core/sorts.go)
# 运行
## Mysql 配置
需要一张表来存储用户信息
```sql
mysql> CREATE TABLE IF NOT EXISTS `user`( 
      `user_id` INT UNSIGNED AUTO_INCREMENT,    
      `user_name` VARCHAR(100) NOT NULL,   
      `user_passowrd` VARCHAR(40) NOT NULL,     
      PRIMARY KEY ( `user_id` ) );
```
具体细节在[mysql](https://github.com/MingweiGuo/Goland/tree/main/SearchEngine/router/mysql)

## main.go 
运行需要配置[PictureEngine](https://github.com/MingweiGuo/Goland/tree/main/SearchEngine/core)和[router/api](https://github.com/MingweiGuo/Goland/tree/main/SearchEngine/router/api)
已经持久化了部分数据，直接运行main.go即可（部分路径需要自行修改）


