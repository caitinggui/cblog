package models

type OtherArticle struct {
	IntIdModelWithoutDeletedAt
	Title     string `gorm:"size:256,index:idx_title" json:"name" form:"title" binding:"lte=128"`                       // 网站名
	ShortUrl  string `gorm:"size:128,unique_index:u_idx_short_url" json:"short_url" form:"short_url" binding:"lte=128"` // 链接地址
	LongUrl   string `gorm:"size:256" json:"long_url" form:"long_url" binding:"lte=128"`                                // 链接地址
	CreatedBy string `gorm:"size:128" json:"created_by" form:"created_by"`
	IfVisited int8   `gorm:"default:-1" json:"if_visited" form:"if_visited" binding:"omitempty,eq=-1|eq=1"` //是否置顶, -1不置顶，1置顶
	Desc      string `gorm:"size:512" json:"desc" form:"desc" binding:"lte=512"`                            // 链接描述
}

type OtherArticleListParam struct {
	// TODO 等gin支持validator.v9时，加上oneof
	Status string `form:"status" binding:"omitempty,eq=d|eq=p"`  //文章状态 d:未发布 p:发布
	Topped int8   `form:"topped" binding:"omitempty,eq=-1|eq=1"` //是否置顶, -1不置顶，1置顶

	Page        uint64 `gorm:"-" form:"page,default=1" binding:"gte=1"`          // 用于分页, start from 1
	PageSize    uint64 `gorm:"-" form:"page_size,default=10" binding:"lte=1000"` // 用于分页
	CategoryId  uint64 `gorm:"-" form:"cate"`
	TagId       uint64 `gorm:"-" form:"tag"` // 用于根据tag查找
	TimeByMonth string `gorm:"-" form:"time_by_month"`
}

func (self *OtherArticle) TableName() string {
	return "other_article"
}

func (self *OtherArticle) Insert() error {
	if self.ID != 0 {
		return ERR_EXIST_ID
	}
	db := DB.Omit("DeletedAt").Create(self)
	return db.Error
}

func (self *OtherArticle) Update() error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(self).Error
}

//如果没有id，会删除整个表，所以要检查一下
func (self *OtherArticle) BeforeDelete() error {
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
func (self *OtherArticle) Delete() error {
	return DB.Delete(self).Error
}

// 更新所有字段时忽略创建时间
func (self *OtherArticle) UpdateAllField() error {
	return DB.Model(self).Omit("CreatedAt", "DeletedAt").Save(self).Error
}

// 更新传进来的字段
// 用struct传进来会忽略掉0值，所以不能用struct
func (self *OtherArticle) UpdateByField(target map[string]interface{}) error {
	return DB.Model(self).Updates(target).Error
}

// 更新时忽略0值
func (self *OtherArticle) UpdateNoneZero(data OtherArticle) error {
	return DB.Model(self).Omit("DeletedAt", "CreatedAt").Updates(data).Error
}

func CountOtherArticleByName(name string) (num int64, err error) {
	err = DB.Model(&OtherArticle{}).Where("name = ?", name).Count(&num).Error
	return
}

func CreateOtherArticle(link *OtherArticle) error {
	return DB.Omit("DeletedAt").Create(link).Error
}

func GetSomeOtherArticles() (links []*OtherArticle, err error) {
	err = DB.Order("if_visited, created_time desc").Limit(2500).Find(&links).Error
	return
}

func GetOtherArticleById(id string) (link OtherArticle, err error) {
	err = DB.First(&link, id).Error
	return
}

func DeleteOtherArticleById(id uint64) error {
	link := OtherArticle{}
	link.ID = id
	return link.Delete()
}
