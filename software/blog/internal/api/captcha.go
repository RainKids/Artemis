package api

import (
	"blog/internal/biz/po"
	"blog/pkg/captcha"
	"blog/pkg/database/redis"
	"blog/pkg/transport/http/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CaptchaController struct {
	logger *zap.Logger
	rdb    *redis.RedisDB
}

func NewCaptchaController(logger *zap.Logger, rdb *redis.RedisDB) *CaptchaController {
	return &CaptchaController{
		logger: logger.With(zap.String("type", "AdvertController")),
		rdb:    rdb,
	}
}

func (c *CaptchaController) Captcha(ctx *gin.Context) {
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
