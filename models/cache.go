package models

import (
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/patrickmn/go-cache"

	"cblog/utils/V"
)

var Cache *cache.Cache

func InitCache(fname string) error {
	Cache = cache.New(V.DefaultExpiration, V.CleanupInterval)
	if fname == "" {
		return nil
	}
	err := Cache.LoadFile(fname)
	return err
}

func DumpCache(fname string) {
	err := Cache.SaveFile(fname)
	logger.Info("dump cache to file ", fname, " result: ", err)
}

// 设置缓存，如果时间为0表示使用默认过期时间
func SetCache(key string, data interface{}, d time.Duration) {
	Cache.Set(key, data, d)
}

// 获取缓存，如果ok为true表示有数据
func GetCache(key string) (data interface{}, ok bool) {
	data, ok = Cache.Get(key)
	return
}

// add one a time
func IncrUint(key string) (n uint, err error) {
	n, err = Cache.IncrementUint(key, 1)
	return
}
