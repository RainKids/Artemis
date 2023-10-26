package service

import (
	"admin/internal/repository"
	"go.uber.org/zap"
)

type roleService struct {
	logger         *zap.Logger
	roleRepository repository.RoleRepository
}

func newRoleService(logger *zap.Logger, roleRepository repository.RoleRepository) RoleService {
	return &roleService{
		logger:         logger.With(zap.String("type", "RoleService")),
		roleRepository: roleRepository,
	}
}
