package models

import (
	"errors"
)

// 文章标签
// form，json，binding都可用于c.Bind
type Tag struct {
	IntIdModel
	Name string `gorm:"size:20" json:"name" form:"name"`
}

func (self *Tag) UpdateNoneZero(data Tag) error {
	return DB.Model(self).Updates(data).Error
}

// 更新时忽略0值
// 要检查ID，防止不小心id为空变为批量更新
func (self *Tag) Update() error {
	if self.ID == 0 {
		errors.New("Empty ID")
	}
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

func CountTagByName(name string) (num int64, err error) {
	err = DB.Model(&Tag{}).Where("name = ?", name).Count(&num).Error
	return
}

func CreateTag(tag *Tag) error {
	return DB.Omit("DeletedAt").Create(tag).Error
}

func GetAllTags() (tags []Tag, err error) {
	err = DB.Find(&tags).Error
	return
}

func GetTagById(id uint64) (tag Tag, err error) {
	err = DB.Where("id = ?", id).First(&tag).Error
	return
}

func DeleteTagById(id uint64) error {
	return DB.Where("id = ?", id).Delete(&Tag{}).Error
}
