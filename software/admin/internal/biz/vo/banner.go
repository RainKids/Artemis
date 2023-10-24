package vo

import "blog/global"

type Banner struct {
	ID        int64            `json:"id"`        // 主键ID
	Path      string           `json:"path"`      // 图片路径
	Hash      string           `json:"hash"`      // 图片的hash值，用于判断重复图片
	Name      string           `json:"name"`      // 图片名称
	ImageType global.ImageType `json:"imageType"` // 图片的类型，本地还是七牛
	global.OperateBy
	global.ModelTime
}

type BannerList struct {
	Result []*Banner `json:"result"`
	Count  int64     `json:"count"`
}
