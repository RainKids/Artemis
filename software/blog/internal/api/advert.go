package api

import "github.com/gin-gonic/gin"

func init() {
	routerAuth = append(routerAuth, registerAuthAdvertRouter)
	routerNoAuth = append(routerNoAuth, registerNoAuthAdvertRouter)
}

func registerAuthAdvertRouter(v1 *gin.RouterGroup, pc *Controller) {

}

func registerNoAuthAdvertRouter(v1 *gin.RouterGroup, pc *Controller) {

}
