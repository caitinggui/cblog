package models

type Link struct {
	IntIdModel
	Name   string `gorm:"size:128" json:"name"` // 网站名
	Url    string `gorm:"size:512" json:"url"`  // 链接地址
	Desc   string `gorm:"size:512" json:"desc"` // 链接描述
	Weight uint64 `json:"weight"`               // 排序
	Views  uint64 `json:"views"`                // 访问次数
}
