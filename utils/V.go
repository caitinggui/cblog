package utils

type ConstantV struct {
	EmptyIntId uint64 // 用来做数据插入的空主键
	EmptyStrId string // 用来做数据插入的空主键
}

var V ConstantV = ConstantV{
	EmptyIntId: 0,
	EmptyStrId: "",
}
