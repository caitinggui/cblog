package models

import ()

// 评论
type Comment struct {
	IntIdModelWithoutDeletedAt
	Content   string  `gorm:"type:text" json:"content" binding:"required"` // 内容
	Article   Article `gorm:"ForeignKey:ArticleId;association_autoupdate:false"`
	ArticleId uint64  `json:"article_id"`
	User      User    `gorm:"ForeignKey:UserId;association_autoupdate:false"` // 用户id
	UserId    uint64
}

func (self *Comment) TableName() string {
	return "comment"
}

func (self *Comment) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Comment) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

// 删除
func (self *Comment) Delete() error {
	return DB.Delete(self).Error
}
