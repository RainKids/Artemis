package api

import (
	"admin/internal/biz/dto"
	"admin/pkg/exception"
	"admin/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthDeptRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthDeptRouter)
}

func registerAuthDeptRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/dept")
	{
		r.GET("", pc.DeptList)
		r.GET("/tree", pc.DeptSetTree)
		r.POST("", pc.DeptCreate)
		r.GET("/:id", pc.DeptRetrieve)
		r.PUT("/:id", pc.DeptUpdate)
		r.DELETE("/:id", pc.ApiDelete)
	}
}

func registerNoAuthDeptRouter(v1 *gin.RouterGroup, pc *Controller) {

}

// SysDept
// @Summary 部门列表数据
// @Description 部门列表
// @Tags SysDept
// @accept application/json
// @Produce application/json
// @Param ddata body dto.DeptSearchParams true "名称, 状态"
// @Success 200 {object} response.Data "{"code": 200, "data": [...]}"
// @Router /api/v1/system/dept [get]
// @Security Bearer
func (c *Controller) DeptList(ctx *gin.Context) {
	deptSearchParams := &dto.DeptSearchParams{}
	err := ctx.ShouldBindQuery(&deptSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Dept().SetDeptPage(ctx, deptSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

// Dept
// @Summary 部门所有数据树
// @Description 部门所有数据树
// @Tags SysDept
// @accept application/json
// @Produce application/json
// @Param ddata body dto.DeptSearchParams true "名称, 状态"
// @Success 200 {object} response.Data "{"code": 200, "data": [...]}"
// @Router /api/v1/system/dept/tree [get]
// @Security Bearer
func (c *Controller) DeptSetTree(ctx *gin.Context) {
	deptSearchParams := &dto.DeptSearchParams{}
	err := ctx.ShouldBindQuery(&deptSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Dept().SetDeptTree(ctx, deptSearchParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

// @Tags SysDept
// @Summary  部门详情
// @Security Bearer
// @accept application/json
// @Produce application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/dept/{id} [get]
func (c *Controller) DeptRetrieve(ctx *gin.Context) {
	var apiUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&apiUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	resp, err := c.service.Dept().Retrieve(ctx, apiUri.ID)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

// SysDept
// @Summary 添加部门
// @Description 获取JSON
// @Tags SysDept
// @Accept  application/json
// @Product application/json
// @Param data body dto.IDUriRequest true "ID"
// @Success 200 {object} response.Data{}"
// @Router /api/v1/system/dept [post]
func (c *Controller) DeptCreate(ctx *gin.Context) {
	deptRequest := new(dto.DeptRequest)
	err := deptRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	resp, err := c.service.Dept().Create(ctx, deptRequest)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

// SysDept
// @Summary 修改部门
// @Description 获取JSON
// @Tags SysDept
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.DeptRequest true "body"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/dept/{id} [put]
// @Security Bearer
func (c *Controller) DeptUpdate(ctx *gin.Context) {
	var deptUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&deptUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	deptRequest := new(dto.DeptRequest)
	err = deptRequest.BindValidParam(ctx)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params error", err))
	}
	err = c.service.Dept().Update(ctx, deptUri.ID, deptRequest)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
	return
}

// SysDept
// @Summary 删除部门
// @Description 删除数据
// @Tags SysDept
// @Param id path int true "id"
// @Success 200 {object} response.Data{}
// @Router /api/v1/system/dept [delete]
// @Security Bearer
func (c *Controller) DeptDelete(ctx *gin.Context) {
	var deptUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&deptUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("api id error", err))
	}
	err = c.service.Dept().Delete(ctx, deptUri.ID)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, nil)
	return
}
