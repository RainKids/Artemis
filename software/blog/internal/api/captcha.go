package api

import (
	"blog/internal/biz/po"
	"blog/pkg/captcha"
	"blog/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	routerAuth = append(routerAuth, registerAuthCaptchaRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthCaptchaRouter)
}

func (c *Controller) Captcha(ctx *gin.Context) {
	// dirver := base64Captcha.NewDriverDigit(global.CaptchaImgHeight, global.CaptchaImgWidth, global.CaptchaKeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(dirver, store.UseWithCtx(c))
	if id, b64s, err := captcha.DriverDigitFunc(); err != nil {
		c.logger.Error("验证码获取失败!", zap.Error(err))
		response.FailedResponse(ctx, err)
		return
	} else {
		response.SuccessResponse(ctx, po.SysCaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: captcha.CaptchaKeyLong,
		})
		return
	}
}

func registerAuthCaptchaRouter(v1 *gin.RouterGroup, pc *Controller) {

}

func registerNoAuthCaptchaRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/captcha")
	{
		r.GET("", pc.Captcha)
	}

}
