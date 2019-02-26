package models

import (
	"errors"
	"strings"

	logger "github.com/caitinggui/seelog"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// 用户信息
type User struct {
	StrIdModelWithoutDeleteAt
	Email       string `gorm:"varchar(128);unique_index;default:null"`   //邮箱
	UserName    string `gorm:"varchar(64);unique_index" json:"userName"` // 用户名
	Password    string `gorm:"varchar(48);not null"`                     //密码
	VerifyState int8   `gorm:"default:-1"`                               //邮箱验证状态
	AvatarUrl   string `gorm:"varchar(256)"`                             // 头像链接, 允许和github头像地址重复
	LockState   int8   `gorm:"default:-1"`                               //锁定状态, 1表示锁定，-1表示未锁定
	IsAdmin     int8   `gorm:"default-1"json:"isAdmin"`                  //是否是管理员, 1表示管理员，-1表示非管理员

	// github包含5个属性，id, name, email, url, avatar
	GithubLoginId string `gorm:"varchar(128);unique_index;default:null"` // github唯一标识
	GithubUrl     string `gorm:"varchar(256)"`                           //github地址
	GithubAvatar  string `gorm:"varchar(256)"`                           // github头像地址
}

// 用uuid代替主键
func (self *User) BeforeCreate(scope *gorm.Scope) error {
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

func (self *User) Update() error {
	if self.ID == "" {
		errors.New("Empty ID")
	}
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

func CreateUser(user *User) error {
	return DB.Omit("DeletedAt").Create(user).Error
}

func GetUserById(id string) (user User, err error) {
	err = DB.Where("id = ?", id).First(&user).Error
	return
}

func DeleteUserById(id string) error {
	return DB.Where("id = ?", id).Delete(&User{}).Error
}

func GetAllUsers() (users []User, err error) {
	err = DB.Order("UpdatedAt desc").Find(&users).Error
	return
}

func IsAdminByUid(uid string) bool {
	return true
}
