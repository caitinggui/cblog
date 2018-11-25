package models

import (
	"strings"

	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Article struct {
	StrIdModel
	Title         string   `gorm:"size:70" json:"title"`            //文章标题
	Body          string   `gorm:"type:longtext" json:"body"`       //富文本
	Status        uint8    `json:"status" json:"status"`            //文章状态 0:未发布 1:发布
	Abstract      string   `gorm:"size:128" json:"abstract"`        //摘要
	Views         uint64   `gorm:"default:0" "json:"views"`         //浏览数
	Likes         uint64   `gorm:"default:0" json:"likes"`          //点赞数
	UserLikes     string   `gorm:"type:text" json:"user_likes"`     //点赞的用户
	Weight        uint64   `gorm:"default:0" json:"weight"`         //推荐权重
	Topped        bool     `gorm:"default:0" json:"topped"`         //是否置顶
	AttachmentUrl string   `gorm:"type:text" json:"attachment_url"` // 附件地址
	Category      Category `gorm:"ForeignKey:ProfileRefer""`
	CategoryId    uint64   `json:"category_id"`
	Tags          []Tag    `gorm:"many2many:article_tags" json:"tags"`
}

// 用uuid代替主键
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	logger.Info("set uuid to id")
	uuid_s := uuid.NewV1().String()
	logger.Debug("uuid.NewV1: ", uuid_s)
	uuid_s = strings.Replace(uuid_s, "-", "", -1)
	err := scope.SetColumn("ID", uuid_s)
	if err != nil {
		logger.Info("set uuid to id failed: ", err)
	}
	return err
}

// 增删改查在业务端记录log
func (article *Article) Insert() error {
	logger.Info("insert article")
	db := DB.Create(article)
	if db.Error != nil {
		logger.Error("insert article error: ", db.Error)
	}
	return db.Error
}

// 所有字段都更新
func (article *Article) UpdateAllField() error {
	return DB.Omit("CreatedAt", "DeletedAt").Save(article).Error
}

// 只更新给定的字段，不用struct是因为它会忽略0,""或者false等
func (article *Article) UpdateByField(data map[string]interface{}) error {
	return DB.Model(article).Updates(data).Error
}

// 删除
func (article *Article) Delete() error {
	return DB.Delete(article).Error
}
