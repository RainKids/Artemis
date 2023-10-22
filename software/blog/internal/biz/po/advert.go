package po

import "blog/global"

type Advert struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Title  string `gorm:"size:32" json:"title"`               // 显示的标题
	Href   string `json:"href"`                               // 跳转链接
	Image  string `json:"image"`                              // 图片
	IsShow bool   `json:"isShow"`                             // 是否展示
	global.OperateBy
	global.ModelTime
}

func (Advert) TableName() string {
	return "blog_advert"
}
