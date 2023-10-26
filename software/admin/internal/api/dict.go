package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthDictRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthDictRouter)
}

func registerAuthDictRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/dict")
	{
		r.GET("")
	}
}

func registerNoAuthDictRouter(v1 *gin.RouterGroup, pc *Controller) {

}
