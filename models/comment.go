package models

import ()

// 评论
type ArticleComment struct {
	IntIdModelWithoutDeletedAt
	Body      string `gorm:"type:text" json:"body" binding:"required"` // 内容
	ArticleId uint64 `json:"article_id"`
	UserId    uint64 `binding:"-"`
	UserName  string `json:"username"`
}

type ArticleCommentWithTitle struct {
	ArticleComment
	Title string `json:"title"`
}

func (self *ArticleComment) TableName() string {
	return "article_comment"
}

func (self *ArticleComment) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *ArticleComment) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *ArticleComment) BeforeDelete() error {
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
func (self *ArticleComment) Delete() error {
	return DB.Delete(self).Error
}

// 获取最新评论
func GetCommentsByCreatedAt(limit int64) (comment []*ArticleComment, err error) {
	com := ArticleComment{}
	err = DB.Model(&com).Order("created_time desc").Limit(limit).Find(&comment).Error
	return
}

// 获取最新评论,带文章标题
func GetCommentsWithTitleByCreatedAt(limit int64) (comment []*ArticleCommentWithTitle, err error) {
	err = DB.Table("article_comment ac").Select("ac.*, article.title").Joins("left join article on ac.article_id=article.id").Order("ac.created_time desc").Limit(limit).Find(&comment).Error
	return
}

// 根据文章id获取全部评论
func GetCommentByArticleId(articleId uint64) (comment []*ArticleComment, err error) {
	com := ArticleComment{}
	err = DB.Model(&com).Order("created_time desc").Where("article_id = ?", articleId).Find(&comment).Error
	return
}

func DeleteCommentById(id uint64) error {
	comment := ArticleComment{}
	comment.ID = id
	return comment.Delete()
}
