package api

import (
	"admin/internal/biz/dto"
	"admin/internal/common"
	"admin/pkg/exception"
	"admin/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthApiRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthApiRouter)
}

// @Tags SysApi
// @Summary api列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.ApiSearchParams true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Data{data=vo.ApiList} "{"code": 200, "data": [...], "message"=""}"
// @Router /api/v1/system/api [get]
func (c *Controller) ApiList(ctx *gin.Context) {
	apiSearchParams := &dto.ApiSearchParams{}
	err := ctx.ShouldBindQuery(&apiSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	p := common.GetPermissionFromContext(ctx)
	resp, err := c.service.Api().List(ctx, apiSearchParams, p)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysApi
// @Summary 创建基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.ApiRequest true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/api [post]
func (c *Controller) ApiCreate(ctx *gin.Context) {
	apiRequest := new(dto.ApiRequest)
	err := apiRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	resp, err := c.service.Api().Create(ctx, apiRequest)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysApi
// @Summary api	详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/api/{id} [get]
func (c *Controller) ApiRetrieve(ctx *gin.Context) {
	var apiUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&apiUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	p := common.GetPermissionFromContext(ctx)
	resp, err := c.service.Api().Retrieve(ctx, apiUri.ID, p)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysApi
// @Summary 更新api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.ApiRequest true "api路径, api中文描述, api组, 方法"
// @Success 200 {object} response.Data '{code:200,data={},message=""}'
// @Router /api/v1/system/api/{id} [put]
func (c *Controller) ApiUpdate(ctx *gin.Context) {
	var apiUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&apiUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	apiRequest := new(dto.ApiRequest)
	err = apiRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	p := common.GetPermissionFromContext(ctx)
	err = c.service.Api().Update(ctx, apiUri.ID, apiRequest, p)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
}

// @Tags SysApi
// @Summary 删除api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{Code,Message,Data}
// @Router /api/v1/system/api/{id} [delete]
func (c *Controller) ApiDelete(ctx *gin.Context) {
	var apiUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&apiUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	p := common.GetPermissionFromContext(ctx)
	err = c.service.Api().Delete(ctx, apiUri.ID, p)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
}

func registerAuthApiRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/api")
	{
		r.GET("", pc.ApiList)
		r.POST("", pc.ApiCreate)
		r.GET("/:id", pc.ApiRetrieve)
		r.PUT("/:id", pc.ApiUpdate)
		r.DELETE("/:id", pc.ApiDelete)
	}
}

func registerNoAuthApiRouter(v1 *gin.RouterGroup, pc *Controller) {

}
