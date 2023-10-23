package service

import (
	"blog/internal/repository"
	"go.uber.org/zap"
)

type service struct {
	advert AdvertService
	banner BannerService
}

func NewService(log *zap.Logger, repository repository.Repository) Service {
	r := &service{
		advert: newAdvertService(log, repository.Advert()),
		banner: newBannerService(log, repository.Banner()),
	}
	return r
}

func (s *service) Advert() AdvertService {
	return s.advert
}

func (s *service) Banner() BannerService {
	return s.banner
}
