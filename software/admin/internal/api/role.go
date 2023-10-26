package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthRoleRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthRoleRouter)
}

func registerAuthRoleRouter(v1 *gin.RouterGroup, pc *Controller) {
	r := v1.Group("/system/role")
	{
		r.GET("")
	}
}

func registerNoAuthRoleRouter(v1 *gin.RouterGroup, pc *Controller) {

}
