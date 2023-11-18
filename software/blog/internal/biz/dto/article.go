package dto

type ArticleSearchParams struct {
	Keyword string `json:"keyword" from:"Keyword"`
	PageInfo
}
