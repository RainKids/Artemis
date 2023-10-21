package api

import (
	"blog/internal/service"
	"blog/pkg/database/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	routerAuth   = make([]func(v1 *gin.RouterGroup, pc *Controller), 0)
	routerNoAuth = make([]func(v1 *gin.RouterGroup, pc *Controller), 0)
)

type Controller struct {
	logger        *zap.Logger
	rdb           *redis.RedisDB
	advertService service.AdvertService
}

func NewController(logger *zap.Logger, rdb *redis.RedisDB, advertService service.AdvertService) *Controller {
	return &Controller{
		logger:        logger.With(zap.String("type", "Controller")),
		rdb:           rdb,
		advertService: advertService,
	}
}

// 路由示例
func InitRouter(r *gin.Engine, pc *Controller) *gin.Engine {

	// 无需认证的路由
	NoAuthRouter(r, pc)
	// 需要认证的路由
	AuthRouter(r, pc)

	return r
}

// 无需认证的路由示例
func NoAuthRouter(r *gin.Engine, pc *Controller) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1")
	for _, f := range routerNoAuth {
		f(v1, pc)
	}
}

// 需要认证的路由示例
func AuthRouter(r *gin.Engine, pc *Controller) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1")
	for _, f := range routerAuth {
		f(v1, pc)
	}
}
