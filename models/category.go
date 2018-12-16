package models

import (
	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
)

// 文章类型
type Category struct {
	IntIdModel
	Name string `gorm:"size:20" json:"name"`
}

func (category *Category) Insert() error {
	logger.Info("insert category")
	db := DB.Create(category)
	logger.Info("category after insert:", category)
	if db.Error != nil {
		logger.Error("insert category error:", db.Error)
	}
	return db.Error
}

// 更新所有字段时忽略创建时间
func (category *Category) UpdateAllField() error {
	return DB.Omit("CreatedAt", "DeletedAt").Save(&category).Error
}

// 更新传进来的字段
// 用struct传进来会忽略掉0值，所以不能用struct
func (category *Category) UpdateByField(target map[string]interface{}) error {
	return DB.Model(&category).Updates(target).Error
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

func GetCategoryById(id string) (cate Category, err error) {
	err = DB.First(&cate, id).Error
	return
}

// 找不到会返回record not find的错误
func GetCategoryByName(name string) (cate Category, err error) {
	err = DB.Where("name = ?", name).First(&cate).Error
	return
}

func DeleteCategoryById(id string) error {
	return DB.Where("id = ?", id).Delete(&Category{}).Error
}

func GetAllCategories() (cates []Category, err error) {
	err = DB.Find(&cates).Error
	return
}
