package models

import ()

// 用户信息
type User struct {
	IntIdModelWithoutDeletedAt
	Email       string `gorm:"varchar(128);unique_index;default:null" json:"email" binding:"lte=128,email,required"` //邮箱
	UserName    string `gorm:"varchar(128);unique_index" json:"username" binding:"lte=64,required"`                  // 用户名
	Password    string `gorm:"varchar(128);not null" json:"password" binding:"lte=48,required"`                      //密码
	VerifyState int8   `gorm:"default:-1" json:"verify_state" binding:"oneof=-1 1"`                                  //邮箱验证状态
	AvatarUrl   string `gorm:"varchar(256)" json:"avatar_url" binding:"url"`                                         // 头像链接, 允许和github头像地址重复
	LockState   int8   `gorm:"default:-1" json:"lock_state" binding:"oneof=-1 1"`                                    //锁定状态, 1表示锁定，-1表示未锁定
	IsAdmin     int8   `gorm:"default-1"json:"is_admin" binding:"oneof=-1 1"`                                        //是否是管理员, 1表示管理员，-1表示非管理员

	// github包含5个属性，id, name, email, url, avatar
	GithubLoginId string `gorm:"varchar(128);unique_index;default:null" json:"github_login_id" binding:"required"` // github唯一标识
	GithubUrl     string `gorm:"varchar(256)" json:"github_url" binding:"url"`                                     //github地址
	GithubAvatar  string `gorm:"varchar(256)" json:"github_avatar" binding:"url"`                                  // github头像地址
}

func (self *User) TableName() string {
	return "user"
}

func (self *User) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *User) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *User) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	if err := DB.First(self, self.ID).Error; err != nil {
		return err
	}
	err := InsertToDeleteDataTable(self)
	return err
}

// 删除
func (self *User) Delete() error {
	return DB.Delete(self).Error
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
	if uid == "" {
		return false
	}
	return true
}
