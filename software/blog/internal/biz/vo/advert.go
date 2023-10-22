package vo

type Advert struct {
	ID    int64  `json:"id"`
	Title string `json:"title"` // 显示的标题
	Href  string `json:"href"`  // 跳转链接
	Image string `json:"image"` // 图片
}

type AdvertList struct {
	Result []*Advert `json:"result"`
	Count  int64
}
