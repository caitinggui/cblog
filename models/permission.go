package models

type Permission struct {
	IntIdModelWithoutDeletedAt
}

func (self *Permission) TableName() string {
	return "permission"
}

func (self *Permission) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *Permission) BeforeUpdate() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	return nil
}

func (self *Permission) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

func (self *Permission) BeforeDelete() error {
	if self.ID == 0 {
		return ERR_EMPTY_ID
	}
	err := InsertToDeleteDataTable(self)
	return err
}

// 删除
func (self *Permission) Delete() error {
	return DB.Delete(self).Error
}
