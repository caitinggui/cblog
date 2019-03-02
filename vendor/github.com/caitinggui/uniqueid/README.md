# uniqueid

分布式唯一id生成器,参考雪花算法[Twitter's Snowflake](https://blog.twitter.com/2010/announcing-snowflake)

生成的为63位无符号整数(sqlite3只支持int64, 不支持uint64):
```
41位存时间戳(time)，从2019.01.25开始计算，可用69年
10位存机器id(workerId), 每个服务要不同，否则可能会产生相同id
2位预留(ReserveId)，可用于标识业务线
2位存异常标识(abnormalityId)，用于系统时间回退等异常，运行时可自动恢复
8位存同一毫秒内的递增值(sequence number), 意味着系统1s内并发仅有(2^8 - 1)*1000=255000,同1毫秒内此值溢出时，系统会休眠1毫秒，扩大此值能有效提高系统并发
```

# 用法

```
import github.com/caitinggui/uniqueid
// 此处可从MySQL或者zookeeper等获取到当前的worker id
var WorkerId uint16 = 2
var ReserveId uint8 = 0
sf := uniqueid.NewUniqueId(WorkerId, ReserveId)
uid, err := sf.NextId()
if err != nil {
    panic(err)
}
```

# 性能

**仅对比不同位数的sequence对uid生成的影响**

在sequence number为15位时:
```
goos: darwin
goarch: amd64
pkg: uniqueid
BenchmarkNextId-4       20000000            81.3 ns/op
PASS
```

在sequence number为8位时:
```
goos: darwin
goarch: amd64
pkg: uniqueid
BenchmarkNextId-4         300000          4518 ns/op
PASS
```
