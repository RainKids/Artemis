package service

import (
	"admin/internal/repository"
	"go.uber.org/zap"
)

type userService struct {
	logger         *zap.Logger
	userRepository repository.UserRepository
}

func newUserService(logger *zap.Logger, userRepository repository.UserRepository) UserService {
	return &userService{
		logger:         logger.With(zap.String("type", "UserService")),
		userRepository: userRepository,
	}
}
