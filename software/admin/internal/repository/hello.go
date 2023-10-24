package repository

import (
	"context"
	"go.uber.org/zap"
)

type helloRepository struct {
	logger *zap.Logger
}

func newHelloRepository(logger *zap.Logger) HelloRepository {
	return &helloRepository{
		logger: logger.With(zap.String("type", "AdvertRepository")),
	}
}

func (h *helloRepository) Hello(c context.Context, message string) (string, error) {
	return "reply -->" + message, nil
}
