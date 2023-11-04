package dto

import (
	"admin/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type DictDataSearchParams struct {
	Name   string `json:"name" `
	Type   int64  `json:"type"`
	Status int    `json:"status" `
	PageInfo
}

type DictDataRequest struct {
	Name      string `json:"name"  validate:"required"`   //岗位名称
	Type      int64  `json:"sort"  validate:"required"`   //岗位排序
	Status    int    `json:"status"  validate:"required"` //状态
	Remark    string `json:"remark"`                      //描述
	OperateBy int64  `json:"createdBy" validate:"required"`
}

func (param *DictDataRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}

type DictTypeSearchParams struct {
	Name   string `json:"name" `
	Type   string `json:"type"`
	Status int    `json:"status" `
	PageInfo
}

type DictTypeRequest struct {
	Name      string `json:"name"  validate:"required"`   //岗位名称
	Type      string `json:"sort"  validate:"required"`   //岗位排序
	Status    int    `json:"status"  validate:"required"` //状态
	Remark    string `json:"remark"`                      //描述
	OperateBy int64  `json:"createdBy" validate:"required"`
}

func (param *DictTypeRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
