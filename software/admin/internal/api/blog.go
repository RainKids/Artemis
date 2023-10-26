package api

import (
	"admin/internal/biz/dto"
	"admin/pkg/exception"
	"admin/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthBlogRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthBlogRouter)
}

// Blog
// @Summary 广告列表
// @Description   广告列表接口
// @Tags Blog
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} response.Data{}
// @Router /api/v1/blog/advert [get]
func (c *Controller) AdvertList(ctx *gin.Context) {
	var advertParams dto.AdvertParamsRequest
	err := ctx.ShouldBindQuery(advertParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Blog().AdvertList(ctx, &advertParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

// Blog
// @Summary 广告列表
// @Description   广告列表接口
// @Tags Blog
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} response.Data{}
// @Router /api/v1//blog/banner [get]
func (c *Controller) BannerList(ctx *gin.Context) {
	var pageInfo dto.PageInfo
	err := ctx.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Blog().BannerList(ctx, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

func registerAuthBlogRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/blog")
	{
		r.GET("/advert", pc.AdvertList)
		r.GET("/banner", pc.BannerList)
	}
}

func registerNoAuthBlogRouter(v1 *gin.RouterGroup, pc *Controller) {

}
