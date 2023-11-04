package dto

import (
	"admin/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type DeptSearchParams struct {
	Name   string `json:"name" `
	Status int    `json:"status" `
	PageInfo
}

type DeptRequest struct {
	ParentId  int64  `json:"parentId"`                   //上级部门
	Leader    string `json:"leader" validate:"required"` //负责人
	Phone     string `json:"phone" `                     //手机
	Email     string `json:"email" `                     //邮箱
	Name      string `json:"name"  validate:"required"`  //部门名称
	Sort      int64  `json:"sort"  validate:"required"`  //部门排序
	Status    int    `json:"status" validate:"required"` //状态
	OperateBy int64  `json:"createdBy" validate:"required"`
}

func (param *DeptRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
