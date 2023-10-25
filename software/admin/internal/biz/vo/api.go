package vo

import "admin/internal/biz/po"

type ApiList struct {
	Result []*po.Api `json:"result"`
	Count  int64     `json:"count"`
}
