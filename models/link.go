package models

type Link struct {
	IntIdModel
	Name   string `gorm:"size:128" json:"name" form:"name"` // 网站名
	Url    string `gorm:"size:512" json:"url" form:"url"`   // 链接地址
	Desc   string `gorm:"size:512" json:"desc" form:"desc"` // 链接描述
	Weight uint64 `json:"weight" form:"weight"`             // 排序
}

func (self *Link) Insert() error {
	return DB.Omit("DeletedAt").Create(self).Error
}

// 更新所有字段时忽略创建时间
func (self *Link) UpdateAllField() error {
	return DB.Model(self).Omit("CreatedAt", "DeletedAt").Save(self).Error
}

// 更新传进来的字段
// 用struct传进来会忽略掉0值，所以不能用struct
func (self *Link) UpdateByField(target map[string]interface{}) error {
	return DB.Model(self).Updates(target).Error
}

// 更新时忽略0值
func (self *Link) UpdateNoneZero(data Link) error {
	return DB.Model(self).Updates(data).Error
}

func CountLinkByName(name string) (num int64, err error) {
	err = DB.Model(&Link{}).Where("name = ?", name).Count(&num).Error
	return
}

func CreateLink(link *Link) error {
	return DB.Omit("DeletedAt").Create(link).Error
}

func GetAllLinks() (links []Link, err error) {
	err = DB.Order("weight desc").Find(&links).Error
	return
}

func GetLinkById(id uint64) (link Link, err error) {
	err = DB.First(&link, id).Error
	return
}

func DeleteLinkById(id uint64) error {
	return DB.Where("id = ?", id).Delete(&Link{}).Error
}
