## 简介

该工具用于将oec leveldb data转换成 rocksdb data

### 使用方式
```
go build --tags rocksdb main.go

nohup ./main {oec}/data {oec}/data_rocksdb/ &
```
以下是转换效果及耗时：

**转换前**
```
du -sh data/*
11G	data/application.db
1.2M	data/blockstore.db
4.0K	data/priv_validator_state.json
4.6M	data/state.db
```

**转换后**
```
du -sh data_rocksdb/*
14G	data_rocksdb/application.db
4.2M	data_rocksdb/blockstore.db
12M	data_rocksdb/state.db
```

**输出及耗时**
```
[../data ../data_rocksdb/]
2021/09/30 05:32:53 convert state start...
2021/09/30 05:32:53 convert application start...
2021/09/30 05:32:53 convert blockstore start...
2021/09/30 05:32:53 convert blockstore end.
2021/09/30 05:32:53 compact blockstore start...
2021/09/30 05:32:53 compact blockstore end.
2021/09/30 05:32:53 convert state end.
2021/09/30 05:32:53 compact state start...
2021/09/30 05:32:53 compact state end.
2021/09/30 05:38:36 convert application end.
2021/09/30 05:38:36 compact application start...
2021/09/30 05:40:24 compact application end.
```
