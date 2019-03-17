package models

type Role struct {
	IntIdModelWithoutDeletedAt
}

func (self *Role) TableName() string {
	return "role"
}

func (self *Role) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Role) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

// 删除
func (self *Role) Delete() error {
	return DB.Delete(self).Error
}
