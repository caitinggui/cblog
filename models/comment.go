package models

// 评论
type Comment struct {
	StrIdModel
	Content string `json:"content"` // 内容
	Article Article
	User    User // 用户id
}
