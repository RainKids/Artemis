package api

import (
	"admin/internal/biz/dto"
	"admin/internal/common"
	"admin/pkg/exception"
	"admin/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthPostRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthPostRouter)
}

func registerAuthPostRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/post")
	{
		r.GET("", pc.ApiList)
		r.POST("", pc.PostCreate)
		r.GET("/:id", pc.PostRetrieve)
		r.PUT("/:id", pc.PostUpdate)
		r.DELETE("/:id", pc.PostDelete)
	}
}

func registerNoAuthPostRouter(v1 *gin.RouterGroup, pc *Controller) {

}

// @Tags SysPost
// @Summary 岗位列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.PostSearchParams true "岗位名称, 岗位代码, 岗位状态"
// @Success 200 {object} response.Data{data=vo.PostList} "{"code": 200, "data": [...], "message"=""}"
// @Router /api/v1/system/post [get]
func (c *Controller) PostList(ctx *gin.Context) {
	postSearchParams := &dto.PostSearchParams{}
	err := ctx.ShouldBindQuery(&postSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Post().List(ctx, postSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysPost
// @Summary 创建岗位
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.ApiRequest  true "岗位名称, 岗位代码,岗位排序, 岗位状态, 岗位描述"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/post [post]
func (c *Controller) PostCreate(ctx *gin.Context) {
	postRequest := new(dto.PostRequest)
	err := postRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	resp, err := c.service.Post().Create(ctx, postRequest)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysPost
// @Summary 岗位详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/post/{id} [get]
func (c *Controller) PostRetrieve(ctx *gin.Context) {
	var postUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&postUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	p := common.GetPermissionFromContext(ctx)
	resp, err := c.service.Api().Retrieve(ctx, postUri.ID, p)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

// @Tags SysPost
// @Summary 更新岗位
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.PostRequest  true "岗位名称, 岗位代码,岗位排序, 岗位状态, 岗位描述"
// @Success 200 {object} response.Data '{code:200,data={},message=""}'
// @Router /api/v1/system/post/{id} [put]
func (c *Controller) PostUpdate(ctx *gin.Context) {
	var postUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&postUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	postRequest := new(dto.PostRequest)
	err = postRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	err = c.service.Post().Update(ctx, postUri.ID, postRequest)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
}

// @Tags SysPost
// @Summary 删除岗位
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{Code,Message,Data}
// @Router /api/v1/system/post/{id} [delete]
func (c *Controller) PostDelete(ctx *gin.Context) {
	var postUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&postUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	err = c.service.Post().Delete(ctx, postUri.ID)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
}
