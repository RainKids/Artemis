package api

import (
	"blog/internal/biz/dto"
	"blog/pkg/exception"
	"blog/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
)

func init() {
	routerAuth = append(routerAuth, registerAuthAdvertRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthAdvertRouter)
}

func (c *Controller) AdvertList(ctx *gin.Context) {
	var advertParams dto.AdvertParamsRequest
	err := ctx.ShouldBindQuery(&advertParams)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Advert().List(ctx, dto.AdvertSearchParams{advertParams.Title}, advertParams.Page, advertParams.PageSize)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
	return
}
func (c *Controller) AdvertRetrieve(ctx *gin.Context) {
	var advertUri dto.IDUriRequest
	err := ctx.ShouldBindUri(&advertUri)
	if err != nil {
		response.FailedResponse(ctx, exception.NewInternalServerError("params err", err))
	}
	resp, err := c.service.Advert().Retrieve(ctx, advertUri.ID)
	if err != nil {
		response.FailedResponse(ctx, exception.NewAPIExceptionFromError(err))
		return
	}
	response.SuccessResponse(ctx, resp)
}

func registerAuthAdvertRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/advert")
	{
		r.GET("", pc.AdvertList)
		r.GET("/:id", pc.AdvertRetrieve)
	}
}

func registerNoAuthAdvertRouter(v1 *gin.RouterGroup, pc *Controller) {

}
