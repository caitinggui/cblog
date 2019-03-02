package models

import (
	"errors"
)

// 文章标签
// form，json，binding都可用于c.Bind
type Tag struct {
	IntIdModelWithoutDeletedAt
	Name     string    `gorm:"size:20;unique_index:uk_name" json:"name"`
	Articles []Article `gorm:"many2many:article_tag;association_autoupdate:false" json:"tags"`
}

func (self *Tag) Insert() error {
	return DB.Omit("DeletedAt").Create(self).Error
}

// 删除
func (self *Tag) Delete() error {
	return DB.Delete(self).Error
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
