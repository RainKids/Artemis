package vo

import "admin/internal/biz/po"

type DictDataList struct {
	Result []*po.DictData `json:"result"`
	Count  int64          `json:"count"`
}

type DictTypeList struct {
	Result []*po.DictType `json:"result"`
	Count  int64          `json:"count"`
}
