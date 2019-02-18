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

func SetCache(key string, data interface{}, d time.Duration) {
	Cache.Set(key, data, d)
}

func GetCache(key string) (data interface{}, ok bool) {
	data, ok = Cache.Get(key)
	return
}
