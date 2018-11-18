package models

import (
	"strings"
	"time"

	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
)

var DB *gorm.DB

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为int
type IntIdModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为string
type StrIdModel struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

type Article struct {
	StrIdModel
	Title         string   `gorm:"size:70" json:"title"`            //文章标题
	Body          string   `gorm:"type:longtext" json:"body"`       //富文本
	Status        uint8    `json:"status" json:"status"`            //文章状态 0:未发布 1:发布
	Abstract      string   `gorm:"size:128" json:"abstract"`        //摘要
	Views         uint     `gorm:"default:0" "json:"views"`         //浏览数
	Likes         uint     `gorm:"default:0" json:"likes"`          //点赞数
	UserLikes     string   `gorm:"type:text" json:"user_likes"`     //点赞的用户
	Weight        uint     `gorm:"default:0" json:"weight"`         //推荐权重
	Topped        bool     `gorm:"default:0" json:"topped"`         //是否置顶
	AttachmentUrl string   `gorm:"type:text" json:"attachment_url"` // 附件地址
	Category      Category `gorm:"ForeignKey:ProfileRefer""`
	CategoryId    uint     `json:"category_id"`
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
func (article *Article) Save() error {
	return DB.Model(article).Save(article).Error
}

// 只更新给定的字段，不用struct是因为它会忽略0,""或者false等
func (article *Article) Update(data map[string]interface{}) error {
	return DB.Model(article).Updates(data).Error
}

// 删除
func (article *Article) Delete() error {
	return DB.Delete(article).Error
}

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
	return DB.Omit("CreatedAt").Save(&category).Error
}

// 更新传进来的字段
// 用struct传进来会忽略掉0值，所以不能用struct
func (category *Category) UpdateByField(target map[string]interface{}) error {
	return DB.Model(&category).Updates(target).Error
}

func FindIfExistCategoryByName(name string) (ifExist bool, err error) {
	var num int
	// 用struct会忽略空字符串，所以少用
	//err = DB.Model(&Category{}).Where(&Category{Name: name}).Count(&num).Error
	err = DB.Model(&Category{}).Where("name = ?", name).Count(&num).Error
	if num > 0 {
		ifExist = true
	} else {
		ifExist = false
	}
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

func GetAllCategories() (cates []Category, err error) {
	err = DB.Find(&cates).Error
	return
}

// 文章标签
type Tag struct {
	IntIdModel
	Name string `gorm:"size:20" json:"name"`
}

// 评论
type Comment struct {
	StrIdModel
	Content string `json:"content"` // 内容
	Article Article
	User    User // 用户id
}

// 用户信息
type User struct {
	StrIdModel
	Email         string    `gorm:"unique_index;default:null"` //邮箱
	Telephone     string    `gorm:"unique_index;default:null"` //手机号码
	Password      string    `gorm:"default:null"`              //密码
	VerifyState   string    `gorm:"default:'0'"`               //邮箱验证状态
	SecretKey     string    `gorm:"default:null"`              //密钥
	OutTime       time.Time //过期时间
	GithubLoginId string    `gorm:"unique_index;default:null"` // github唯一标识
	GithubUrl     string    //github地址
	IsAdmin       bool      //是否是管理员
	AvatarUrl     string    // 头像链接
	NickName      string    // 昵称
	LockState     bool      `gorm:"default:'0'"` //锁定状态
}

// 用uuid代替主键
func (user *User) BeforeCreate(scope *gorm.Scope) error {
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

type VisitorIP struct {
	IntIdModel
	IP      string `gorm:"size:64" json:"ip"`       // 访问者IP
	Country string `gorm:"size:128" json:"country"` // 国家
	City    string `gorm:"size:128" json:"city"`    // 城市
	Referer string `gorm:"size:255" json:"referer"` // 来源地
	Article Article
}

type Link struct {
	IntIdModel
	Name   string `gorm:"size:128" json:"name"` // 网站名
	Url    string `gorm:"size:512" json:"url"`  // 链接地址
	Desc   string `gorm:"size:512" json:"desc"` // 链接描述
	Weight uint   `json:"weight"`               // 排序
	Views  uint   `json:"views"`                // 访问次数
}

type Permission struct {
	IntIdModel
}

type Role struct {
	IntIdModel
}

func InitDB() (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	//db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Asia/Shanghai")
	DB = db
	db.SingularTable(true) //全局设置表名不可以为复数形式。
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour * 6)
	db.LogMode(true)
	db.AutoMigrate(&Article{}, &Category{}, &Comment{}, &Tag{}, &Link{}, &VisitorIP{}, &User{}, &Permission{}, &Role{})
	return
}
