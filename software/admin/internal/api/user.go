package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthUserRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthUserRouter)
}

func registerAuthUserRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/User")
	{
		r.GET("")
	}
}

func registerNoAuthUserRouter(v1 *gin.RouterGroup, pc *Controller) {

}
