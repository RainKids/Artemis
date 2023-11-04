package vo

import "admin/internal/biz/po"

type DeptLabel struct {
	ID       int64       `json:"id"`
	Label    string      `json:"label"`
	Children []DeptLabel `json:"children"`
}

type DeptId struct {
	ID int `json:"id"`
}

type DeptList struct {
	Result []po.Dept `json:"result"`
}
type DeptLabelList struct {
	Result []DeptLabel `json:"result"`
}
