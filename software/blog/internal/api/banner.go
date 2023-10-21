package api

import "blog/global"

type Banner struct {
	ID        int64            `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Path      string           `json:"path"`                               // 图片路径
	Hash      string           `json:"hash"`                               // 图片的hash值，用于判断重复图片
	Name      string           `gorm:"size:38" json:"name"`                // 图片名称
	ImageType global.ImageType `gorm:"default:1" json:"image_type"`        // 图片的类型，本地还是七牛
	global.OperateBy
	global.ModelTime
}

func (Banner) TableName() string {
	return "blog_banner"
}
