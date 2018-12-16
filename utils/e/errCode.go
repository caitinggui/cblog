package e

import ()

const (
	SUCCESS = 0
	ERROR   = 500

	ERR_INVALID_PARAM = 4001
	ERR_NO_DATA       = 4002

	ERR_CACHE = 5001
	ERR_SQL   = 5002
)

var ErrMsg = map[int]string{
	SUCCESS:           "请求成功",
	ERROR:             "未知错误",
	ERR_INVALID_PARAM: "参数错误",
	ERR_NO_DATA:       "无数据",
	ERR_CACHE:         "缓存异常",
	ERR_SQL:           "数据库异常",
}

func GetMsg(code int) string {
	msg, ok := ErrMsg[code]
	if ok {
		return msg
	}
	return ErrMsg[ERROR]
}
