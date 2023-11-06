package api

import (
	"admin/internal/service"
	"admin/pkg/database/redis"
	"admin/pkg/transport/http/middleware/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	routerAuth   = make([]func(v1 *gin.RouterGroup, pc *Controller), 0)
	routerNoAuth = make([]func(v1 *gin.RouterGroup, pc *Controller), 0)
)

type Controller struct {
	v        *viper.Viper
	logger   *zap.Logger
	enforcer *casbin.SyncedEnforcer
	rdb      *redis.RedisDB
	service  service.Service
}

func NewController(logger *zap.Logger, v *viper.Viper, rdb *redis.RedisDB, service service.Service) *Controller {
	return &Controller{
		logger:  logger.With(zap.String("type", "Controller")),
		v:       v,
		rdb:     rdb,
		service: service,
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
	v1.Use(auth.AuthenticateToken(), auth.AuthCheckRole(pc.logger, pc.enforcer))
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
