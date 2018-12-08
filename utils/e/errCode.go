package e

import ()

const (
	SUCCESS = 0
	ERROR   = 500

	ERR_INVALID_PARAM = 40001
)

var ErrMsg = map[int]string{
	SUCCESS:           "请求成功",
	ERROR:             "未知错误",
	ERR_INVALID_PARAM: "参数错误",
}

func GetMsg(code int) string {
	msg, ok := ErrMsg[code]
	if ok {
		return msg
	}
	return ErrMsg[ERROR]
}
