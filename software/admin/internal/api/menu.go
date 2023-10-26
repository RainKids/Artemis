package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthMenuRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthMenuRouter)
}

func registerAuthMenuRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/menu")
	{
		r.GET("")
	}
}

func registerNoAuthMenuRouter(v1 *gin.RouterGroup, pc *Controller) {

}
