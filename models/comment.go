package models

import (
	"github.com/jinzhu/gorm"

	"cblog/utils"
)

// 评论
type Comment struct {
	IntIdModelWithoutDeletedAt
	Content string `json:"content"` // 内容
	Article Article
	User    User // 用户id
}

func (self *Comment) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("ID", utils.GenerateId())
	return err
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

func (self *Comment) BeforeUpdate() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	return nil
}

func (self *Comment) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

func (self *Comment) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	err := InsertToDeleteDataTable(self)
	return err
}

// 删除
func (self *Comment) Delete() error {
	return DB.Delete(self).Error
}
