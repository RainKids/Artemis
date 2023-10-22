package api

import (
	"blog/pkg/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func CreateInitControllersFn(
	pc *Controller,
) http.InitControllers {
	return func(r *gin.Engine) {
		InitRouter(r, pc)
	}
}

// ProviderSet controllers wire
var ProviderSet = wire.NewSet(NewController, CreateInitControllersFn)
