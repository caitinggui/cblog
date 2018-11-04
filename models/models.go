package models

import "time"

type Article struct {
	Id         uint32     `gorm:"primary_key;auto_increment"  json:"id"`
	UniqueId   string     `json:"unique_id"`                               //唯一id 防止遍历
	CategoryId uint32     `json:"category_id"`                             //所属目录
	Title      string     `gorm:"size:50" json:"title"`                    //文章标题
	Label      string     `gorm:"size:100" json:"label"`                   //文章标签
	Status     uint8      `json:"status"`                                  //文章状态 1:完结，2:更新
	Body       string     `gorm:"type:longtext;not null" json:"body"`      //富文本
	PureText   string     `gorm:"type:longtext;not null" json:"pure_text"` //纯粹的文章文本
	ImageURL   string     `gorm:"size:100" json:"image_url"`               //图片地址
	FileURL    string     `gorm:":size:100" json:"file_url"`               //附件文件地址
	Original   uint8      `gorm:"default:0" json:"original"`               //是否原创 1原创，2转载
	Source     string     `gorm:"size:100" json:"source"`                  //文章资源地址，如果原创为空，非原创必须带有地址
	Operator   string     `gorm:"size:50" json:"operator"`                 //创建者
	Disabled   bool       `gorm:"default:true" json:"disabled"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (Article) TableName() string {
	return "article"
}

type Category struct {
	Id   uint32 `gorm:"primary_key;auto_increment"  json:"id"`
	Name string `gorm:"size:20" json:"name"`
}

func MyName() (Name string) {
	Name = "models"
	return Name
}
