package api

import (
	"blog/internal/service"
	"go.uber.org/zap"
)

type AdvertController struct {
	logger        *zap.Logger
	advertService service.AdvertService
}

func NewAdvertController(logger *zap.Logger, advertService service.AdvertService) *AdvertController {
	return &AdvertController{
		logger:        logger.With(zap.String("type", "AdvertController")),
		advertService: advertService,
	}
}
