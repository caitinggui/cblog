package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"cblog/config"
)

var DB *gorm.DB

var EXIST_ID = errors.New("Exist ID")

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为int
// 忽略DeleteAt
type IntIdModel struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `sql:"index" json:"-", form:"-"`
}

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为string
type StrIdModel struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-" form:"-"`
}

// 忽略DeleteAt
type IntIdModelWithoutDeletedAt struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 忽略DeleteAt
type StrIdModelWithoutDeleteAt struct {
	ID        string    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func InitDB() (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	//db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Asia/Shanghai")
	DB = db
	db.SingularTable(true) //全局设置表名不可以为复数形式。
	db.DB().SetMaxIdleConns(config.Config.Mysql.MaxIdle)
	db.DB().SetMaxOpenConns(config.Config.Mysql.MaxOpen)
	db.DB().SetConnMaxLifetime(time.Hour * config.Config.Mysql.MaxLife)
	db.LogMode(config.Config.Mysql.LogMode)
	db.AutoMigrate(&Article{}, &Category{}, &Comment{}, &Tag{}, &Link{}, &Visitor{}, &User{}, &Permission{}, &Role{})
	return
}
