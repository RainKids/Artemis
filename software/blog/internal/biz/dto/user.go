package dto

import (
	validator "blog/pkg/transport/http/middleware/validator"
	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Uuid     string `json:"Uuid" validate:"required"`
}

func (param *UserLoginRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}

type UserRegisterRequest struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"   validate:"required"`
}

func (param *UserRegisterRequest) BindValidParam(c *gin.Context) error {
	return validator.DefaultGetValidParams(c, param)
}
