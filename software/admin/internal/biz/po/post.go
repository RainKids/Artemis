package po

import "admin/global"

type Post struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Name   string `gorm:"size:128;" json:"name"`              //岗位名称
	Code   string `gorm:"uniqueIndex;size:128;" json:"code"`  //岗位代码
	Sort   int    `gorm:"size:4;" json:"sort"`                //岗位排序
	Status int    `gorm:"size:4;" json:"status"`              //状态
	Remark string `gorm:"size:255;" json:"remark"`            //描述
	global.OperateBy
	global.ModelTime
}

func (Post) TableName() string {
	return "sys_post"
}
