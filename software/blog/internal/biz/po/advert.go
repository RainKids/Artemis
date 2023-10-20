package po

import "blog/global"

type Advert struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Title  string `gorm:"size:32" json:"title"`               // 显示的标题
	Href   string `json:"href"`                               // 跳转链接
	Images string `json:"images"`                             // 图片
	IsShow bool   `json:"is_show"`                            // 是否展示
	global.OperateBy
	global.ModelTime
}

func (Advert) TableName() string {
	return "blog_advert"
}
