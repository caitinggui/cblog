package models

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"cblog/config"
	"cblog/utils"
)

var DB *gorm.DB

var ERR_EXIST_ID = errors.New("Exist ID")
var ERR_EMPTY_ID = errors.New("Empty ID")
var ERR_INVALID_TIME = errors.New("Invalid Time")

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
	ID        uint64         `gorm:"primary_key" json:"id"`
	CreatedAt utils.JsonTime `json:"created_at"`
	UpdatedAt utils.JsonTime `json:"updated_at"`
	DeletedAt utils.JsonTime `sql:"index" json:"-", form:"-"`
}

// 用来覆盖gorm.Model，主要对json方式做出改变, 主键为string
type StrIdModel struct {
	ID        string          `gorm:"primary_key" json:"id"`
	CreatedAt utils.JsonTime  `json:"created_at"`
	UpdatedAt utils.JsonTime  `json:"updated_at"`
	DeletedAt *utils.JsonTime `sql:"index" json:"-" form:"-"`
}

// 忽略DeleteAt
// TODO 截止到gin 1.3.0,在序列化json时不支持time_format
// TODO 有个bug，没有检查传来的updated_at等参数，等validator升到V9再修复
type IntIdModelWithoutDeletedAt struct {
	ID        uint64    `gorm:"primary_key" json:"id" form:"id"` // 如果用"gorm:bigint"，在sqlite下无法自增
	CreatedAt time.Time `gorm:"column:created_time" json:"created_at" binding:"-" time_format:"2006-01-02T15:04:05"`
	UpdatedAt time.Time `gorm:"column:last_modified_time" json:"updated_at" binding:"-"`
}

// Get struct name by reflect
func (self *IntIdModelWithoutDeletedAt) TableName() string {
	if t := reflect.TypeOf(self); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (self *IntIdModelWithoutDeletedAt) BeforeCreate(scope *gorm.Scope) error {
	//logger.Info("BeforeCreate table: ", self.TableName()) // table还是base
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	if !self.UpdatedAt.IsZero() || !self.CreatedAt.IsZero() {
		return ERR_INVALID_TIME
	}
	err := scope.SetColumn("ID", utils.GenerateId())
	return err
}

func (self *IntIdModelWithoutDeletedAt) BeforeUpdate() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	return nil
}

func (self *IntIdModelWithoutDeletedAt) Insert() error {
	return nil
}

func (self *IntIdModelWithoutDeletedAt) Update() error {
	return nil
}

func (self *IntIdModelWithoutDeletedAt) Delete() error {
	return nil
}

// 忽略DeleteAt
type StrIdModelWithoutDeletedAt struct {
	ID        string         `gorm:"primary_key" json:"id"`
	CreatedAt utils.JsonTime `json:"created_at" binding:"-"`
	UpdatedAt utils.JsonTime `json:"updated_at" binding:"-"`
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

func Ping() error {
	return DB.DB().Ping()
}

func InitDB() (db *gorm.DB) {
	//db, err := gorm.Open("sqlite3", "./foo.db")
	db, err := gorm.Open("mysql", config.Config.Mysql.Server)
	if err != nil {
		panic(err)
	}
	DB = db
	db.SingularTable(true) //全局设置表名不可以为复数形式。
	db.DB().SetMaxIdleConns(config.Config.Mysql.MaxIdle)
	db.DB().SetMaxOpenConns(config.Config.Mysql.MaxOpen)
	db.DB().SetConnMaxLifetime(time.Hour * config.Config.Mysql.MaxLife)
	db.LogMode(config.Config.Mysql.LogMode)
	db.AutoMigrate(&Article{}, &Category{}, &DeletedData{}, &ArticleComment{}, OtherArticle{}, &Tag{}, &Link{}, &Visitor{}, &User{}, &Permission{}, &Role{})
	return
}
