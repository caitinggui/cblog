package utils

type ConstantV struct {
	EmptyIntId uint64 // 用来做数据插入的空主键
	EmptyStrId string // 用来做数据插入的空主键

	MaxPageSize     uint64 // 分页的每页最大条数
	DefaultPageSize uint64 // 分页的每页默认条数
}

var V ConstantV = ConstantV{
	EmptyIntId:      0,
	EmptyStrId:      "",
	MaxPageSize:     1000,
	DefaultPageSize: 10,
}
