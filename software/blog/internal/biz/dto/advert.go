package dto

type AdvertParamsRequest struct {
	AdvertSearchParams
	PageInfo
}
type AdvertSearchParams struct {
	Title string `json:"title"`
}
