package utils

import (
	"strconv"
)

func StrToUnit64(s string) (n uint64) {
	// 10进制，64位
	n, _ = strconv.ParseUint(s, 10, 64)
	return n
}

func StrToInt64(s string) (n int64) {
	n, _ = strconv.ParseInt(s, 10, 64)
	return
}
