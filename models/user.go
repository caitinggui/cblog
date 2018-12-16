package models

import (
	"strings"
	"time"

	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

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

func IsAdminExistByUid(uid string) bool {
	return true
}
