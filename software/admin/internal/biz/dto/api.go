package dto

import (
	"admin/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type ApiSearchParams struct {
	Title    string `json:"title" `
	Method   string `json:"method"`
	Path     string `json:"path" `
	ApiGroup string `json:"apiGroup"`
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`
	PageInfo
}

type ApiRequest struct {
	Handle      string `json:"handle" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Path        string `json:"path" validate:"required"` // api路径
	Type        string `json:"type" validate:"required"`
	Description string `json:"description" validate:"required"` // api中文描述
	ApiGroup    string `json:"apiGroup" validate:"required"`    // api组
	Method      string `json:"method" validate:"required"`
	OperateBy   int64  `json:"createdBy" validate:"required"`
}

func (param *ApiRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
