package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为int
// 忽略DeleteAt
type IntIdModel struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-", form:"-"`
}

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为string
type StrIdModel struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-" form:"-"`
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
	db.AutoMigrate(&Article{}, &Category{}, &Comment{}, &Tag{}, &Link{}, &Visitor{}, &User{}, &Permission{}, &Role{})
	return
}
