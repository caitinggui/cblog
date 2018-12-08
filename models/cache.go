package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

func InitCache(fname string) error {
	Cache = cache.New(30*time.Second, time.Minute)
	if fname == "" {
		return nil
	}
	err := Cache.LoadFile(fname)
	return err
}

func DumpCache(fname string) error {
	return Cache.SaveFile(fname)
}

func SetCache(key string, data interface{}, d time.Duration) {
	Cache.Set(key, data, d)
}

func GetCache(key string) (data interface{}, ok bool) {
	data, ok = Cache.Get(key)
	return
}
