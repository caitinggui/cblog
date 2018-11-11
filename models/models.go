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
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为string
type StrIdModel struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Article struct {
	IntIdModel
	UniqueId   string `json:"unique_id"`                               //唯一id 防止遍历
	CategoryId uint32 `json:"category_id"`                             //所属目录
	Title      string `gorm:"size:50" json:"title"`                    //文章标题
	Label      string `gorm:"size:100" json:"label"`                   //文章标签
	Status     uint8  `json:"status"`                                  //文章状态 1:完结，2:更新
	Body       string `gorm:"type:longtext;not null" json:"body"`      //富文本
	PureText   string `gorm:"type:longtext;not null" json:"pure_text"` //纯粹的文章文本
	ImageURL   string `gorm:"size:100" json:"image_url"`               //图片地址
	FileURL    string `gorm:":size:100" json:"file_url"`               //附件文件地址
	Original   uint8  `gorm:"default:0" json:"original"`               //是否原创 1原创，2转载
	Source     string `gorm:"size:100" json:"source"`                  //文章资源地址，如果原创为空，非原创必须带有地址
	Operator   string `gorm:"size:50" json:"operator"`                 //创建者
	Disabled   bool   `gorm:"default:true" json:"disabled"`
}

// 用uuid代替主键
// TODO 这里要参考category中的
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV1())
	return nil
}

func (article *Article) Insert() error {
	logger.Info("insert article")
	db := DB.Create(article)
	if db.Error != nil {
		logger.Error("insert article error: ", db.Error)
	}
	return db.Error

}

// 文章类型
type Category struct {
	StrIdModel
	Name string `gorm:"size:20" json:"name"`
}

// 用uuid代替主键
func (category *Category) BeforeCreate(scope *gorm.Scope) error {
	logger.Info("set id to uuid")
	uuid_s := uuid.NewV1().String()
	logger.Debug("uuid.NewV1: ", uuid_s)
	uuid_s = strings.Replace(uuid_s, "-", "", -1)
	err := scope.SetColumn("ID", uuid_s)
	if err != nil {
		logger.Info("set id to uuid failed: ", err)
	}
	return err
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
	//db.LogMode(true)
	db.AutoMigrate(&Article{}, Category{})
	//db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
	return
}
