package dto

import (
	"admin/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type PostSearchParams struct {
	Name   string `json:"name" `
	Code   string `json:"code"`
	Status int    `json:"status" `
	PageInfo
}

type PostRequest struct {
	Name      string `json:"name"  validate:"required"`   //岗位名称
	Code      string `json:"code"  validate:"required"`   //岗位代码
	Sort      int    `json:"sort"  validate:"required"`   //岗位排序
	Status    int    `json:"status"  validate:"required"` //状态
	Remark    string `json:"remark"`                      //描述
	OperateBy int64  `json:"createdBy" validate:"required"`
}

func (param *PostRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
