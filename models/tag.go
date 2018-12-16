package models

// 文章标签
// form，json，binding都可用于c.Bind
type Tag struct {
	IntIdModel
	Name string `gorm:"size:20" json:"name" form:"name"`
}

func (tag *Tag) UpdateNonzero(data Tag) error {
	return DB.Model(tag).Updates(data).Error
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
