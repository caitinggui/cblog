package utils

import (
	"time"
)

type ConstantV struct {
	EmptyIntId uint64 // 用来做数据插入的空主键
	EmptyStrId string // 用来做数据插入的空主键

	MaxPageSize     uint64 // 分页的每页最大条数
	DefaultPageSize uint64 // 分页的每页默认条数

	DefaultExpiration time.Duration // cache过期时间
	CleanupInterval   time.Duration // 定期清理cache的间隔时间
}

var V ConstantV = ConstantV{
	EmptyIntId:        0,
	EmptyStrId:        "",
	MaxPageSize:       1000,
	DefaultPageSize:   10,
	DefaultExpiration: 30 * time.Second,
	CleanupInterval:   time.Minute,
}
