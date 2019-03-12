package models

import (
	"encoding/json"
	"errors"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"cblog/config"
)

var DB *gorm.DB

var ERR_EXIST_ID = errors.New("Exist ID")
var ERR_EMPTY_ID = errors.New("Empty ID")

// 定义数据表的接口
type DataTable interface {
	TableName() string
	Insert() error
	Update() error
	Delete() error
}

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
	ID        uint64    `gorm:"primary_key" json:"id"` // 如果用"gorm:bigint"，在sqlite下无法自增
	CreatedAt time.Time `json:"created_at" binding:"-"`
	UpdatedAt time.Time `json:"updated_at" binding:"-"`
}

// 忽略DeleteAt
type StrIdModelWithoutDeletedAt struct {
	ID        string    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at" binding:"-"`
	UpdatedAt time.Time `json:"updated_at" binding:"-"`
}

type DeletedData struct {
	IntIdModelWithoutDeletedAt
	TableName string `gorm:"size:255"`
	Data      string `gorm:"type:longtext"`
}

func (self *DeletedData) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func InsertToDeleteDataTable(data DataTable) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("序列化要删除的数据失败: ", data, err)
	}
	deletedData := DeletedData{
		TableName: data.TableName(),
		Data:      string(jsonData),
	}
	err = deletedData.Insert()
	return err
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
	db.AutoMigrate(&Article{}, &Category{}, &DeletedData{}, &Comment{}, &Tag{}, &Link{}, &Visitor{}, &User{}, &Permission{}, &Role{})
	return
}
