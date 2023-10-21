package service

import (
	"blog/internal/repository"
	"go.uber.org/zap"
)

type service struct {
	advert AdvertService
}

func (s *service) Advert() AdvertService {
	return s.Advert()
}

func NewService(log *zap.Logger, repository repository.Repository) Service {
	r := &service{
		advert: newAdvertService(log, repository.Advert()),
	}
	return r
}
