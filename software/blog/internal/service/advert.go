package service

import (
	"blog/internal/biz/dto"
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"blog/internal/repository"
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type advertService struct {
	logger           *zap.Logger
	advertRepository repository.AdvertRepository
}

func (a *advertService) Create(ctx context.Context, advert *po.Advert) (*po.Advert, error) {
	return a.advertRepository.Create(ctx, advert)
}

func (a *advertService) Retrieve(ctx context.Context, id int64) (*po.Advert, error) {
	return a.advertRepository.Retrieve(ctx, id)
}

func (a *advertService) List(c context.Context, params dto.AdvertSearchParams, page, pageSize int) (*vo.AdvertList, error) {
	query, err := convertor.StructToMap(params)
	for key, v := range query {
		if v == 0 || v == "" {
			delete(query, key)
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "Get query params error")
	}
	advertList, count, err := a.advertRepository.List(c, query, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &vo.AdvertList{
		advertList,
		count,
	}, nil
}

func (a *advertService) SysList(c context.Context, params dto.AdvertSearchParams, page, pageSize int) ([]*po.Advert, int64, error) {
	query, err := convertor.StructToMap(params)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Get query params error")
	}
	return a.advertRepository.SysList(c, query, page, pageSize)
}

func (a *advertService) Update(ctx context.Context, advert *po.Advert) (*po.Advert, error) {
	return a.advertRepository.Update(ctx, advert)
}

func (a *advertService) Delete(ctx context.Context, id int64) error {
	return a.advertRepository.Delete(ctx, id)
}

func newAdvertService(logger *zap.Logger, advertRepository repository.AdvertRepository) AdvertService {
	return &advertService{
		logger:           logger.With(zap.String("type", "AdvertService")),
		advertRepository: advertRepository,
	}
}
