package models

import (
	"fmt"
	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
)

// 文章类型
type Category struct {
	IntIdModelWithoutDeletedAt
	Name string `gorm:"size:20;unique_index:uk_name" binding:"lte=20,required" json:"name" form:"name"`
}

func (self *Category) TableName() string {
	return "category"
}

func (self *Category) Insert() error {
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Category) AfterUpdate() error {
	logger.Infof("index articles after category(%v) change", self.ID)
	go func() {
		artis, _ := GetArticleIdsByCategory(self.Name)
		ids := make([]string, len(artis))
		for _, arti := range artis {
			ids = append(ids, fmt.Sprint(arti.ID))
		}
		IndexArticleByIds(ids)
	}()
	return nil
}

func (self *Category) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *Category) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	if err := DB.First(self, self.ID).Error; err != nil {
		return err
	}
	err := InsertToDeleteDataTable(self)
	go func() {
		logger.Infof("Update index after Category %v deleted", self.ID)
		artis, _ := GetArticleIdsByCategory(self.Name)
		ids := make([]string, len(artis))
		for _, arti := range artis {
			ids = append(ids, fmt.Sprint(arti.ID))
		}
		IndexArticleByIds(ids)
	}()
	return err
}

// 删除
func (self *Category) Delete() error {
	return DB.Delete(self).Error
}

// 更新所有字段时忽略创建时间
func (self *Category) UpdateAllField() error {
	return DB.Model(self).Omit("CreatedAt", "DeletedAt").Save(self).Error
}

// 更新传进来的字段
// 用struct传进来会忽略掉0值，所以不能用struct
func (self *Category) UpdateByField(target map[string]interface{}) error {
	return DB.Model(self).Updates(target).Error
}

// 更新时忽略0值
func (self *Category) UpdateNoneZero(data Category) error {
	return DB.Model(self).Updates(data).Error
}

// 找不到就返回false，包括数据库异常也是false
func CheckIsExistCategoryByName(name string) bool {
	cate := Category{}
	err := DB.Select("name").Where("name = ?", name).First(&cate).Error
	if err == nil {
		return true
	}
	// 如果不是没找到，说明是操作数据库失败，要记录，但也是表示找不到
	if !gorm.IsRecordNotFoundError(err) {
		logger.Error("check category error: ", err)
	}
	return false
}

func CountCategoryByName(name string) (num int64, err error) {
	// 用struct会忽略空字符串，所以少用
	//err = DB.Model(&Category{}).Where(&Category{Name: name}).Count(&num).Error
	err = DB.Model(&Category{}).Where("name = ?", name).Count(&num).Error
	return
}

func GetCategoryById(id uint64) (cate Category, err error) {
	err = DB.First(&cate, id).Error
	return
}

// 找不到会返回record not find的错误
func GetCategoryByName(name string) (cate Category, err error) {
	err = DB.Where("name = ?", name).First(&cate).Error
	return
}

func DeleteCategoryById(id uint64) error {
	data := Category{}
	data.ID = id
	return data.Delete()
}

func GetAllCategories() (cates []Category, err error) {
	err = DB.Find(&cates).Error
	return
}
