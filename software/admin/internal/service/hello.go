package service

import (
	"admin/internal/repository"
	"context"
	"go.uber.org/zap"
)

type helloService struct {
	logger          *zap.Logger
	helloRepository repository.HelloRepository
}

func newHelloService(logger *zap.Logger, helloRepository repository.HelloRepository) HelloService {
	return &helloService{
		logger:          logger.With(zap.String("type", "BlogService")),
		helloRepository: helloRepository,
	}
}

func (h *helloService) Hello(c context.Context, message string) (string, error) {
	return h.helloRepository.Hello(c, message)
}
