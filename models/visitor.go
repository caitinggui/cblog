package models

import (
	"cblog/utils"
	"cblog/utils/V"
	"errors"
	logger "github.com/caitinggui/seelog"
)

var ERR_EMPTY_IP = errors.New("Empty IP")

type Visitor struct {
	IntIdModelWithoutDeletedAt
	IP        string   `gorm:"size:64" json:"ip"`        // 访问者IP
	Country   string   `gorm:"size:128" json:"country"`  // 国家
	Province  string   `gorm:"size:128" json:"province"` // 省份
	City      string   `gorm:"size:128" json:"city"`     // 城市
	Isp       string   `gorm:"size:128" json:"isp"`
	Referer   string   `gorm:"size:255" json:"referer"` // 来源地
	Article   *Article `gorm:"ForeignKey:ArticleId;association_autoupdate:false" binding:"-"`
	ArticleId string   `json:"article_id"`
}

func (self *Visitor) TableName() string {
	return "visitor"
}

func (self *Visitor) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	go func() {
		// 新增失败就重新统计
		if _, err := IncrCacheUint(V.VisitorSum); err != nil {
			logger.Warnf("there doesn't exist %s in cache", V.VisitorSum)
			visitorSum, _ := CountVisitor()
			SetCache(V.VisitorSum, visitorSum, 0)
		}
	}()
	return db.Error
}

func (self *Visitor) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *Visitor) BeforeDelete() error {
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
func (self *Visitor) Delete() error {
	return DB.Delete(self).Error
}

func (self *Visitor) PraseIp() error {
	if self.IP == "" {
		return ERR_EMPTY_IP
	}
	ip2Region := utils.Ip2Region{IP: self.IP}
	err := ip2Region.PraseIp()
	if err != nil {
		return err
	}
	self.Country = ip2Region.Country
	self.Province = ip2Region.Province
	self.City = ip2Region.City
	self.Isp = ip2Region.Isp
	return nil
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
	err = DB.Order("CreatedAt desc").Find(&visitors).Error
	return
}

// page start from 0
func GetVisitors(page, pageSize uint) (visitors []Visitor, err error) {
	// Find要放order之类的后面
	err = DB.Order("ID desc").Offset(pageSize * page).Limit(pageSize).Find(&visitors).Error
	return
}

func GetVisitorById(id uint64) (visitor Visitor, err error) {
	err = DB.First(&visitor, id).Error
	return
}

func DeleteVisitorById(id uint64) error {
	return DB.Where("id = ?", id).Delete(&Visitor{}).Error
}

func GetVisitorsByArticle(articleId string) (visitors []Visitor, err error) {
	err = DB.Order("ID desc").Limit(5).Where("article_id=?", articleId).Find(&visitors).Error
	return
}

func CountVisitor() (n uint, err error) {
	err = DB.Model(&Visitor{}).Count(&n).Error
	//visitorSum, _ := GetCache(V.VisitorSum)
	return
}
