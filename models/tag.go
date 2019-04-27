package models

import ()

// 文章标签
// form，json，binding都可用于c.Bind
type Tag struct {
	IntIdModelWithoutDeletedAt
	Name     string    `gorm:"size:20;unique_index" json:"name" binding:"lte=20,required"`
	Articles []Article `gorm:"many2many:article_tag;association_autoupdate:false" json:"tags"`
}

func (self *Tag) TableName() string {
	return "tag"
}

func (self *Tag) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Tag) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *Tag) BeforeDelete() error {
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
func (self *Tag) Delete() error {
	return DB.Delete(self).Error
}

func (self *Tag) UpdateNoneZero(data Tag) error {
	return DB.Model(self).Updates(data).Error
}

func CountTagByName(name string) (num int64, err error) {
	err = DB.Model(&Tag{}).Where("name = ?", name).Count(&num).Error
	return
}

func GetTagByName(name string) (tag Tag, err error) {
	err = DB.Where("name = ?", name).First(&tag).Error
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
	tag := Tag{}
	tag.ID = id
	return tag.Delete()
}
