package models

import ()

// 评论
type Comment struct {
	IntIdModelWithoutDeletedAt
	Body      string  `gorm:"type:text" json:"body" binding:"required"` // 内容
	Article   Article `gorm:"ForeignKey:ArticleId;association_autoupdate:false" binding:"-"`
	ArticleId uint64  `json:"article_id"`
	User      User    `gorm:"ForeignKey:UserId;association_autoupdate:false" binding:"-"` // 用户id
	UserId    uint64  `binding:"-"`
	Name      string  `json:"name"`
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

//如果没有id，会删除整个表，所以要检查一下
func (self *Comment) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	if err := DB.First(self, self.ID).Error; err != nil {
		return err
	}
	err := InsertToDeleteDataTable(self)
	return err
}

// 删除
func (self *Comment) Delete() error {
	return DB.Delete(self).Error
}

// 获取最新评论
func GetCommentsByCreatedAt(limit int64) (comment []*Comment, err error) {
	com := Comment{}
	err = DB.Model(&com).Order("created_time desc").Limit(limit).Find(&comment).Error
	return
}

// 根据文章id获取全部评论
func GetCommentByArticleId(articleId uint64) (comment []*Comment, err error) {
	com := Comment{}
	err = DB.Model(&com).Order("created_time desc").Where("article_id = ?", articleId).Find(&comment).Error
	return
}
