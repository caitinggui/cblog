package models

import (
	"strings"

	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Article struct {
	StrIdModelWithoutDeleteAt
	Title         string   `gorm:"size:70" json:"title"`            //文章标题
	Body          string   `gorm:"type:longtext" json:"body"`       //富文本
	Status        int8     `json:"status" json:"status"`            //文章状态 -1:未发布 1:发布
	Abstract      string   `gorm:"size:128" json:"abstract"`        //摘要
	Views         uint64   `gorm:"default:0" "json:"views"`         //浏览数
	Likes         uint64   `gorm:"default:0" json:"likes"`          //点赞数
	UserLikes     string   `gorm:"type:text" json:"user_likes"`     //点赞的用户
	Weight        uint64   `gorm:"default:0" json:"weight"`         //推荐权重
	Topped        int8     `gorm:"default:-1" json:"topped"`        //是否置顶, -1不置顶，1置顶
	AttachmentUrl string   `gorm:"type:text" json:"attachment_url"` // 附件地址
	Category      Category `gorm:"ForeignKey:CategoryId;association_autoupdate:false"`
	CategoryId    uint64   `json:"category_id"`
	Tags          []Tag    `gorm:"many2many:article_tag;association_autoupdate:false" json:"tags"`
}

// 用uuid代替主键
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	uuid_s := uuid.NewV1().String()
	uuid_s = strings.Replace(uuid_s, "-", "", -1)
	err := scope.SetColumn("ID", uuid_s)
	return err
}

func (self *Article) TableName() string {
	return "article"
}

// 增删改查在业务端记录log
func (self *Article) Insert() error {
	if self.ID != "" {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Article) BeforeUpdate() error {
	if self.ID == "" {
		return ERR_EMPTY_ID
	}
	return nil
}

func (self *Article) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

// 如果没有id，会删除整个表，所以要检查一下
func (self *Article) BeforeDelete() error {
	if self.ID == "" {
		return ERR_EMPTY_ID
	}
	err := InsertToDeleteDataTable(self)
	return err
}

// 删除
func (self *Article) Delete() error {
	return DB.Delete(self).Error
}

// 所有字段都更新
func (article *Article) UpdateAllField() error {
	return DB.Model(article).Omit("CreatedAt", "DeletedAt").Save(article).Error
}

// 只更新给定的字段，不用struct是因为它会忽略0,""或者false等
func (article *Article) UpdateByField(data map[string]interface{}) error {
	return DB.Model(article).Updates(data).Error
}

func GetArticleById(id string) (article Article, err error) {
	err = DB.Where("id = ?", id).First(&article).Error
	return
}

func GetArticlesByCategory(category string) (articles []Article, err error) {
	err = DB.Table("article ").Select("article.*").Where("cg.name = ?", category).Joins("join category cg on article.category_id=cg.id").Find(&articles).Error
	return
}

func GetArticleByTag(tagName string) (articles []*Article, err error) {
	tag := Tag{}
	//articles = make([]Article, 3)
	logger.Info(articles, &articles, tag, &tag)
	err = DB.First(&tag, "name = ?", tagName).Error
	if gorm.IsRecordNotFoundError(err) {
		logger.Warn("找不到tag: ", tagName, err)
		return nil, err
	}
	logger.Info(articles, &articles, tag, &tag)
	//err = DB.Model(&tag).Association("article_tag").Find(&articles).Error
	//err = DB.Model(&tag).Related(&articles, "article_tag").Error
	err = DB.Table("article").Select("article.*").Where("ag.tag_id = ?", tag.ID).Joins("join article_tag ag on article.id=ag.article_id").Find(&articles).Error
	return
}
