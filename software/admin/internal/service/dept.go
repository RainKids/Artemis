package service

import (
	"admin/internal/repository"
	"go.uber.org/zap"
)

type deptService struct {
	logger         *zap.Logger
	deptRepository repository.DeptRepository
}

func newDeptService(logger *zap.Logger, deptRepository repository.DeptRepository) DeptService {
	return &deptService{
		logger:         logger.With(zap.String("type", "DeptService")),
		deptRepository: deptRepository,
	}
}
