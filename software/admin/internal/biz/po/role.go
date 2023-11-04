package po

import "admin/global"

type Role struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"roleId"` // 主键ID
	Name      string  `json:"name" gorm:"size:128;"`                  // 角色名称
	Status    string  `json:"status" gorm:"size:4;"`                  //
	Key       string  `json:"key" gorm:"size:128;"`                   //角色代码
	Sort      int     `json:"sort" gorm:""`                           //角色排序
	Flag      string  `json:"flag" gorm:"size:128;"`                  //
	Remark    string  `json:"remark" gorm:"size:255;"`                //备注
	Admin     bool    `json:"admin" gorm:"size:4;"`
	DataScope string  `json:"dataScope" gorm:"size:128;"`
	MenuIds   []int64 `json:"menuIds" gorm:"-"`
	DeptIds   []int   `json:"deptIds" gorm:"-"`
	Dept      []Dept  `json:"Dept" gorm:"many2many:sys_role_dept;"`
	Menu      *[]Menu `json:"Menu" gorm:"many2many:sys_role_menu;"`
	global.OperateBy
	global.ModelTime
}

func (Role) TableName() string {
	return "sys_role"
}
