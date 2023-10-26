package po

import "admin/global"

type Dept struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement;"` //部门编码
	ParentId int64  `json:"parentId" gorm:""`                    //上级部门
	Path     string `json:"deptPath" gorm:"size:255;"`           //
	Name     string `json:"name"  gorm:"size:128;"`              //部门名称
	Sort     int    `json:"sort" gorm:"size:4;"`                 //排序
	Leader   int64  `json:"leader" gorm:"size:128;"`             //负责人
	Phone    string `json:"phone" gorm:"size:11;"`               //手机
	Email    string `json:"email" gorm:"size:64;"`               //邮箱
	Status   int    `json:"status" gorm:"size:4;"`               //状态
	global.OperateBy
	global.ModelTime
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
	Children  []Dept `json:"children" gorm:"-"`
}

func (*Dept) TableName() string {
	return "sys_dept"
}
