package auth

import (
	"admin/pkg/exception"
	"admin/pkg/jwt"
	"admin/pkg/transport/http/response"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UrlInfo struct {
	Url    string
	Method string
}

// CasbinExclude casbin 排除的路由列表
var CasbinExclude = []UrlInfo{}

// AuthCheckRole 权限检查中间件
func AuthCheckRole(log *zap.Logger, enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var errMsg string
		data := c.GetString("JWT_PAYLOAD")
		v, err := jwt.ParseToken(data)
		if err != nil {
			response.FailedResponse(c, exception.NewUnauthorized("Parse Authorization Token error, %s", err))
		}
		var res, casbinExclude bool

		//检查权限
		if v.SuperAdmin {
			res = true
			c.Next()
			return
		}
		for _, i := range CasbinExclude {
			if util.KeyMatch2(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				casbinExclude = true
				break
			}
		}
		if casbinExclude {
			errMsg = fmt.Sprintf("Casbin exclusion, no validation method:%s path:%s", c.Request.Method, c.Request.URL.Path)
			log.Info(errMsg)
			c.Next()
			return
		}
		res, err = enforcer.Enforce(v.ID, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			errMsg = fmt.Sprintf("AuthCheckRole method:%s path:%s", c.Request.Method, c.Request.URL.Path)
			log.Error(errMsg, zap.Error(err))
			response.FailedResponse(c, exception.NewPermissionDeny(errMsg))
			return
		}

		if res {
			errMsg = fmt.Sprintf("isTrue: %v role: %s method: %s path: %s", res, v.ID, c.Request.Method, c.Request.URL.Path)
			log.Info(errMsg)
			c.Next()
		} else {
			errMsg = fmt.Sprintf("isTrue: %v role: %s method: %s path: %s message: %s", res, v.ID, c.Request.Method, c.Request.URL.Path, "当前request无权限，请管理员确认！")
			log.Warn(errMsg)
			response.FailedResponse(c, exception.NewPermissionDeny(errMsg))
			c.Abort()
			return
		}

	}
}
