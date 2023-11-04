package vo

import "admin/internal/biz/po"

type MenuLabel struct {
	ID       int64       `json:"id,omitempty"`
	Label    string      `json:"label,omitempty"`
	Children []MenuLabel `json:"children,omitempty"`
}

type MenuRole struct {
	po.Menu
	IsSelect bool `json:"is_select" gorm:"-"`
}
