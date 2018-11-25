package models

type VisitorIP struct {
	IntIdModel
	IP      string `gorm:"size:64" json:"ip"`       // 访问者IP
	Country string `gorm:"size:128" json:"country"` // 国家
	City    string `gorm:"size:128" json:"city"`    // 城市
	Referer string `gorm:"size:255" json:"referer"` // 来源地
	Article Article
}
