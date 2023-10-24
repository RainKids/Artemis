package vo

import "blog/global"

type Advert struct {
	ID     int64  `json:"id"`     // 主键ID
	Title  string `json:"title"`  // 显示的标题
	Href   string `json:"href"`   // 跳转链接
	Image  string `json:"image"`  // 图片
	IsShow bool   `json:"isShow"` // 是否展示
	global.OperateBy
	global.ModelTime
}

type AdvertList struct {
	Result []*Advert `json:"result"`
	Count  int64     `json:"count"`
}
