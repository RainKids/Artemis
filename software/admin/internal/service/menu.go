package service

import (
	"admin/internal/repository"
	"go.uber.org/zap"
)

type menuService struct {
	logger         *zap.Logger
	menuRepository repository.MenuRepository
}

func newMenuService(logger *zap.Logger, menuRepository repository.MenuRepository) MenuService {
	return &menuService{
		logger:         logger.With(zap.String("type", "MenuService")),
		menuRepository: menuRepository,
	}
}
