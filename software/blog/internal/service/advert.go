package service

import (
	"blog/internal/repository"
	"go.uber.org/zap"
)

type advertService struct {
	logger           *zap.Logger
	advertRepository repository.AdvertRepository
}

func newAdvertService(logger *zap.Logger, advertRepository repository.AdvertRepository) AdvertService {
	return &advertService{
		logger:           logger.With(zap.String("type", "AdvertService")),
		advertRepository: advertRepository,
	}
}
