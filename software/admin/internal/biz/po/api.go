package po

import "admin/global"

type Api struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Handle      string `json:"handle" gorm:"size:128;comment:handle"`
	Title       string `json:"title" gorm:"size:128;comment:标题"`
	Path        string `json:"path" gorm:"comment:api路径"` // api路径
	Type        string `json:"type" gorm:"size:16;comment:接口类型"`
	Description string `json:"description" gorm:"comment:api中文描述"`   // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`         // api组
	Method      string `json:"method" gorm:"default:GET;comment:方法"` // 方法:创建POST|查看GET(默认)|更新PUT|删除DELETE
	global.OperateBy
	global.ModelTime
}

func (Api) TableName() string {
	return "sys_api"
}
