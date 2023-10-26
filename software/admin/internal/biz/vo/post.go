package vo

import "admin/internal/biz/po"

type PostList struct {
	Result []*po.Post `json:"result"`
	Count  int64      `json:"count"`
}
