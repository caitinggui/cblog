package models

type Link struct {
	IntIdModel
	Name   string `gorm:"size:128" json:"name" form:"name"` // 网站名
	Url    string `gorm:"size:512" json:"url" form:"url"`   // 链接地址
	Desc   string `gorm:"size:512" json:"desc" form:"desc"` // 链接描述
	Weight uint64 `json:"weight" form:"weight"`             // 排序
}

func (link *Link) UpdateNonzero(data Link) error {
	return DB.Model(link).Updates(data).Error
}

func CountLinkByName(name string) (num int64, err error) {
	err = DB.Model(&Link{}).Where("name = ?", name).Count(&num).Error
	return
}

func CreateLink(link *Link) error {
	return DB.Omit("DeletedAt").Create(link).Error
}

func GetAllLinks() (links []Link, err error) {
	err = DB.Find(&links).Error
	return
}

func GetLinkById(id uint64) (link Link, err error) {
	err = DB.First(&link, id).Error
	return
}

func DeleteLinkById(id uint64) error {
	return DB.Where("id = ?", id).Delete(&Link{}).Error
}
