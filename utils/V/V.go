package V

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	EmptyIntId uint64 = 0  // 用来做数据插入的空主键
	EmptyStrId string = "" // 用来做数据插入的空主键

	MaxPageSize     uint64 = 1000 // 分页的每页最大条数
	DefaultPageSize uint   = 10   // 分页的每页默认条数

	DefaultExpiration time.Duration = 30 * time.Second   // cache过期时间
	CleanupInterval   time.Duration = time.Minute        // cache永不过期
	NoExpiration      time.Duration = cache.NoExpiration // 定期清理cache的间隔时间

	CurrentUser string = "CurrentUser" // 保存当前用户的key
)
