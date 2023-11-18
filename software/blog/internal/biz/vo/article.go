package vo

type ArticleSearch struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"` // 显示的标题
	Content string `json:"content"`
}

type ArticleSearchList struct {
	Result []*ArticleSearch `json:"result"`
	Count  int64            `json:"count"`
}
