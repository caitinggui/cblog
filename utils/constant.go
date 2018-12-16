package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type ConstantV struct {
	EmptyIntId uint64 // 用来做数据插入的空主键
	EmptyStrId string // 用来做数据插入的空主键

	MaxPageSize     uint64 // 分页的每页最大条数
	DefaultPageSize uint64 // 分页的每页默认条数

	DefaultExpiration time.Duration // cache过期时间
	NoExpiration      time.Duration // cache永不过期
	CleanupInterval   time.Duration // 定期清理cache的间隔时间

	CurrentUser string // 保存当前用户的key
}

var V ConstantV = ConstantV{
	EmptyIntId: 0,
	EmptyStrId: "",

	MaxPageSize:     1000,
	DefaultPageSize: 10,

	DefaultExpiration: 30 * time.Second,
	CleanupInterval:   time.Minute,
	NoExpiration:      cache.NoExpiration,

	CurrentUser: "CurrentUser",
}
