package po

import "admin/global"

type Menu struct {
	ID          int64        `json:"id" gorm:"primaryKey;autoIncrement"` // 主键ID
	MenuName    string       `json:"name" gorm:"type:varchar(128);"`
	Title       string       `json:"title" gorm:"type:varchar(64);"`
	Icon        string       `json:"icon" gorm:"type:varchar(128);"`
	Path        string       `json:"path" gorm:"type:varchar(128);"`
	Paths       string       `json:"paths" gorm:"size:128;"`
	IsIframe    string       `json:"isIframe" gorm:"type:varchar(1);"`    //是否为内嵌
	IsLink      string       `json:"isLink" gorm:"type:varchar(255);"`    //是否超链接菜单
	IsHide      string       `json:"isHide" gorm:"type:varchar(1);"`      //显示状态（0显示 1隐藏）
	IsKeepAlive string       `json:"isKeepAlive" gorm:"type:varchar(1);"` //是否缓存组件状态（0是 1否）
	IsAffix     string       `json:"isAffix" gorm:"type:varchar(1);"`     //是否固定在 tagsView 栏上（0是 1否）
	Status      string       `json:"status" gorm:"type:varchar(1);"`      // 菜单状态（0正常 1停用）
	Remark      string       `json:"remark"  gorm:"type:varchar(256);"`   // 备注
	Type        string       `json:"type" gorm:"size:1;"`
	Action      string       `json:"action" gorm:"size:16;"`
	Permission  string       `json:"permission" gorm:"size:255;"`
	ParentId    int64        `json:"parentId" gorm:"size:11;"`
	NoCache     bool         `json:"noCache" gorm:"size:8;"`
	Breadcrumb  string       `json:"breadcrumb" gorm:"size:255;"`
	Component   string       `json:"component" gorm:"size:255;"`
	Sort        int          `json:"sort" gorm:"size:4;"`
	Visible     string       `json:"visible" gorm:"size:1;"`
	IsFrame     string       `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	SysApi      []Api        `json:"sysApi" gorm:"many2many:sys_menu_api_rule"`
	Apis        []int64      `json:"apis" gorm:"-"`
	DataScope   string       `json:"dataScope" gorm:"-"`
	Params      string       `json:"params" gorm:"-"`
	RoleId      int          `gorm:"-"`
	Children    []Menu       `json:"children,omitempty" gorm:"-"`
	IsSelect    bool         `json:"is_select" gorm:"-"`
	MenuParams  []MenuParam  `json:"menuParams" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MenuButtons []MenuButton `json:"menuButtons" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	global.OperateBy
	global.ModelTime
}

type SysMenuSlice []Menu

func (x SysMenuSlice) Len() int           { return len(x) }
func (x SysMenuSlice) Less(i, j int) bool { return x[i].Sort < x[j].Sort }
func (x SysMenuSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (*Menu) TableName() string {
	return "sys_menu"
}

type MenuParam struct {
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"` // 主键ID
	MenuID int64  `json:"menuId"`
	Type   string `json:"type" gorm:"comment:地址栏携带参数为params还是query"` // 地址栏携带参数为params还是query
	Key    string `json:"key" gorm:"comment:地址栏携带参数的key"`            // 地址栏携带参数的key
	Value  string `json:"value" gorm:"comment:地址栏携带参数的值"`            // 地址栏携带参数的值
	global.OperateBy
	global.ModelTime
}

func (*MenuParam) TableName() string {
	return "sys_menu_param"
}

type MenuButton struct {
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"` // 主键ID
	Name   string `json:"name" gorm:"comment:按钮关键key"`
	Desc   string `json:"desc" gorm:"按钮备注"`
	MenuID int64  `json:"MenuId" gorm:"comment:菜单ID"`
	global.OperateBy
	global.ModelTime
}

func (*MenuButton) TableName() string {
	return "sys_menu_button"
}
