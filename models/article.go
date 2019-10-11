package models

import (
	"unicode/utf8"

	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
)

type Article struct {
	IntIdModelWithoutDeletedAt
	Title         string   `gorm:"size:70" json:"title" form:"title" binding:"lte=70,required"`               //文章标题
	Body          string   `gorm:"type:longtext" json:"body" form:"editormd-markdown-doc" binding:"required"` //富文本
	Status        int8     `gorm:"default:-1" json:"status" form:"status" binding:"omitempty,eq=-1|eq=1"`     //文章状态 -1:未发布 1:发布
	Abstract      string   `gorm:"size:128" json:"abstract" form:"abstract" binding:"lte=128"`                //摘要
	Views         uint64   `gorm:"default:0" json:"views" form:"views" binding:"-"`                           //浏览数
	Likes         uint64   `gorm:"default:0" json:"likes" form:"likes"`                                       //点赞数
	UserLikes     string   `gorm:"type:text" json:"user_likes"`                                               //点赞的用户
	Weight        uint64   `gorm:"default:0" json:"weight" form:"weight"`                                     //推荐权重
	Topped        int8     `gorm:"default:-1" json:"topped" form:"topped" binding:"omitempty,eq=-1|eq=1"`     //是否置顶, -1不置顶，1置顶
	AttachmentUrl string   `gorm:"type:text" json:"attachment_url"`                                           //附件地址
	Category      Category `gorm:"ForeignKey:CategoryId;save_associations:false" binding:"-"`
	CategoryId    uint64   `json:"category_id" form:"category_id"`
	Tags          []Tag    `gorm:"many2many:article_tag;save_associations:false" binding:"-" json:"tags"`

	TagsId []uint64 `gorm:"-" json:"tags_id" binding:"dive,omitempty" form:"tags_id"`
}

type ArticleListParam struct {
	// TODO 等gin支持validator.v9时，加上oneof
	Status int8 `form:"status" binding:"omitempty,eq=-1|eq=1"` //文章状态 -1:未发布 1:发布
	Topped int8 `form:"topped" binding:"omitempty,eq=-1|eq=1"` //是否置顶, -1不置顶，1置顶

	Page       uint64 `gorm:"-" form:"page,default=1" binding:"gte=1"`          // 用于分页, start from 0
	PageSize   uint64 `gorm:"-" form:"page_size,default=10" binding:"lte=1000"` // 用于分页
	CategoryId uint64 `gorm:"-" form:"category_id"`
	TagId      uint64 `gorm:"-" form:"tag_id"` // 用于根据tag查找
}

func (self *Article) TableName() string {
	return "article"
}

// 增删改查在业务端记录log
func (self *Article) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	if self.Abstract == "" {
		bodyLen := utf8.RuneCountInString(self.Body)
		if bodyLen > 128 {
			bodyLen = 128
		}
		self.Abstract = self.Body[:bodyLen]
	}
	if self.CategoryId != 0 {
		self.Category = Category{}
		self.Category.ID = self.CategoryId
	}
	for _, tag_id := range self.TagsId {
		tag := Tag{}
		tag.ID = tag_id
		self.Tags = append(self.Tags, tag)
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Article) Update() error {
	if self.CategoryId != 0 {
		self.Category.ID = self.CategoryId
	}
	// TODO 要清除已有的tag_id
	for _, tag_id := range self.TagsId {
		tag := Tag{}
		tag.ID = tag_id
		self.Tags = append(self.Tags, tag)
	}
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *Article) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	if err := DB.First(self, self.ID).Error; err != nil {
		return err
	}
	err := InsertToDeleteDataTable(self)
	return err
}

//删除前要完整的查一遍数据
func (self *Article) Delete() error {
	// TODO 可以记录一下tag有哪些
	return DB.Delete(self).Error
}

// 替换tags
func (self *Article) ReplaceTags(tags []Tag) error {
	return DB.Model(&self).Association("Tags").Replace(tags).Error
}

func (self *Article) GetInfoColumn() string {
	return "id, title, abstract, likes, status, topped, views, weight, created_at, updated_at"
}

func (self *Article) GetDefaultOrder() string {
	return "topped desc, id desc"
}

// 所有字段都更新
func (article *Article) UpdateAllField() error {
	return DB.Model(article).Omit("CreatedAt", "DeletedAt").Save(article).Error
}

// 只更新给定的字段，不用struct是因为它会忽略0,""或者false等
func (article *Article) UpdateByField(data map[string]interface{}) error {
	return DB.Model(article).Updates(data).Error
}

// 获取所有文章名
func GetAllArticleNames() (articles []*Article, err error) {
	err = DB.Select("id, title").Find(&articles).Error
	return
}

// 分页获取文章简单信息
// Omit不在查询中生效，仅在Update中生效
// form不用引用是为了规范一下，毕竟form不能再修改
func GetArticleInfos(form ArticleListParam) (articles []*Article, err error) {
	arti := Article{}
	err = DB.Where(&form).Select(arti.GetInfoColumn()).Limit(form.PageSize).Offset((form.Page - 1) * form.PageSize).Order(arti.GetDefaultOrder()).Find(&articles).Error
	return
}

// 获取所有文章简单信息
func GetAllArticleInfos() (articles []*Article, err error) {
	arti := Article{}
	err = DB.Model(&arti).Select(arti.GetInfoColumn()).Order(arti.GetDefaultOrder()).Find(&articles).Error
	return
}

func GetArticleById(id string) (article Article, err error) {
	err = DB.Where("id = ?", id).First(&article).Error
	return
}

func GetFullArticleById(id string) (article Article, err error) {

	//cate := Category{}
	//tag := Tag{}
	// 这里不能用TableName，否则字段名不对
	err = DB.Where("id = ?", id).Preload("Category").Preload("Tags").First(&article).Error
	return
}

// get all article detail for search index
func GetFullArticle() (articles []Article, err error) {
	err = DB.Preload("Category").Preload("Tags").Find(&articles).Error
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

// TODO tag_id要清除
func DeleteArticleById(id uint64) error {
	arti := Article{}
	arti.ID = id
	return arti.Delete()
}

func CountArticle() (n uint64, err error) {
	err = DB.Model(&Article{}).Count(&n).Error
	return
}
