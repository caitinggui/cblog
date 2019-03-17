package models

import ()

type Visitor struct {
	IntIdModelWithoutDeletedAt
	IP      string `gorm:"size:64" json:"ip"`       // 访问者IP
	Country string `gorm:"size:128" json:"country"` // 国家
	City    string `gorm:"size:128" json:"city"`    // 城市
	Referer string `gorm:"size:255" json:"referer"` // 来源地
	Article Article
}

func (self *Visitor) TableName() string {
	return "visitor"
}

func (self *Visitor) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Visitor) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

// 删除
func (self *Visitor) Delete() error {
	return DB.Delete(self).Error
}

func (visitor *Visitor) UpdateNonzero(data Visitor) error {
	return DB.Model(visitor).Updates(data).Error
}

func CountVisitorByIP(ip string) (num int64, err error) {
	err = DB.Model(&Visitor{}).Where("ip = ?", ip).Count(&num).Error
	return
}

func CreateVisitor(visitor *Visitor) error {
	return DB.Omit("DeletedAt").Create(visitor).Error
}

func GetAllVisitors() (visitors []Visitor, err error) {
	err = DB.Find(&visitors).Error
	return
}

// TODO 有问题
func GetVisitors(page, pageSize uint64) (visitors []Visitor, err error) {
	err = DB.Find(&visitors).Order("CreatedAt desc").Offset(pageSize * page).Limit(page).Error
	return
}

func GetVisitorById(id uint64) (visitor Visitor, err error) {
	err = DB.First(&visitor, id).Error
	return
}

func DeleteVisitorById(id uint64) error {
	return DB.Where("id = ?", id).Delete(&Visitor{}).Error
}
