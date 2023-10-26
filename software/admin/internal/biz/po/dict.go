package po

import "admin/global"

type DictData struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Sort      int    `json:"sort" gorm:"size:20;comment:DictSort"`
	Label     string `json:"label" gorm:"size:128;comment:DictLabel"`
	Value     string `json:"value" gorm:"size:255;comment:DictValue"`
	Type      string `json:"type" gorm:"size:64;comment:DictType"`
	CssClass  string `json:"cssClass" gorm:"size:128;comment:CssClass"`
	ListClass string `json:"listClass" gorm:"size:128;comment:ListClass"`
	IsDefault string `json:"isDefault" gorm:"size:8;comment:IsDefault"`
	Status    int    `json:"status" gorm:"size:4;comment:Status"`
	Default   string `json:"default" gorm:"size:8;comment:Default"`
	Remark    string `json:"remark" gorm:"size:255;comment:Remark"`
	global.OperateBy
	global.ModelTime
}

func (DictData) TableName() string {
	return "sys_dict_data"
}

type DictType struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Name   string `json:"name" gorm:"size:128;comment:DictName"`
	Type   string `json:"type" gorm:"size:128;comment:DictType"`
	Status int    `json:"status" gorm:"size:4;comment:Status"`
	Remark string `json:"remark" gorm:"size:255;comment:Remark"`
	global.OperateBy
	global.ModelTime
}

func (DictType) TableName() string {
	return "sys_dict_type"
}
