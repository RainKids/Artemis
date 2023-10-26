package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthDeptRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthDeptRouter)
}

func registerAuthDeptRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/dept")
	{
		r.GET("")
	}
}

func registerNoAuthDeptRouter(v1 *gin.RouterGroup, pc *Controller) {

}
