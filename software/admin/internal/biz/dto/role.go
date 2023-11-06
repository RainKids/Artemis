package dto

import (
	"admin/internal/biz/po"
	"admin/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type RoleSearchParams struct {
	Name   string `json:"name" `
	Status int    `json:"status" `
	PageInfo
}

type RoleRequest struct {
	Name      string    `form:"name" comment:"角色名称" validate:"required"` // 角色名称
	Status    string    `form:"status" comment:"状态" validate:"required"` // 状态
	Key       string    `form:"key" comment:"角色代码" validate:"required"`  // 角色代码
	Sort      int       `form:"sort" comment:"角色排序"`                     // 角色排序
	Flag      string    `form:"flag" comment:"标记"`                       // 标记
	Remark    string    `form:"remark" comment:"备注"`                     // 备注
	Admin     bool      `form:"admin" comment:"是否管理员"`
	DataScope string    `form:"dataScope"`
	Menu      []po.Menu `form:"menu"`
	MenuIds   []int     `form:"menuIds"`
	Dept      []po.Dept `form:"dept"`
	DeptIds   []int     `form:"deptIds"`
	OperateBy int64     `json:"createdBy" validate:"required"`
}

func (param *RoleRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}

type RoleDataScopeRequest struct {
	ID        int    `json:"id" validate:"required"`
	DataScope string `json:"dataScope" validate:"required"`
	DeptIds   []int  `json:"deptIds"`
}

func (param *RoleDataScopeRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}

type UpdateStatusRequest struct {
	Id     int64 `form:"id" comment:"角色编码"`   // 角色编码
	Status int   `form:"status" comment:"状态"` // 状态

}

func (param *UpdateStatusRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
