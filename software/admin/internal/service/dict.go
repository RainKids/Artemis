package service

import (
	"admin/internal/repository"
	"go.uber.org/zap"
)

type dictService struct {
	logger         *zap.Logger
	dictRepository repository.DictRepository
}

func newDictService(logger *zap.Logger, dictRepository repository.DictRepository) DictService {
	return &dictService{
		logger:         logger.With(zap.String("type", "DictService")),
		dictRepository: dictRepository,
	}
}
