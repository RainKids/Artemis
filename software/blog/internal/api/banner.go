package api

import (
	"blog/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthBannerRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthBannerRouter)
}

func (c *Controller) BannerList(ctx *gin.Context) {
	resp, err := c.service.Banner().List(ctx)
	if err != nil {
		response.FailedResponse(ctx, err)
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}

func registerAuthBannerRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/banner")
	{
		r.GET("", pc.BannerList)
	}
}

func registerNoAuthBannerRouter(v1 *gin.RouterGroup, pc *Controller) {

}
