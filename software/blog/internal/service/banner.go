package service

import (
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"blog/internal/repository"
	"context"
	"go.uber.org/zap"
)

type bannerService struct {
	logger           *zap.Logger
	bannerRepository repository.BannerRepository
}

func newBannerService(logger *zap.Logger, bannerRepository repository.BannerRepository) BannerService {
	return &bannerService{
		logger:           logger.With(zap.String("type", "AdvertService")),
		bannerRepository: bannerRepository,
	}
}

func (b *bannerService) Create(ctx context.Context, banner *po.Banner) (*po.Banner, error) {
	return b.bannerRepository.SysCreate(ctx, banner)
}

func (b *bannerService) Retrieve(ctx context.Context, id int64) (*po.Banner, error) {
	return b.bannerRepository.SysRetrieve(ctx, id)
}

func (b *bannerService) List(ctx context.Context) (*vo.BannerList, error) {
	bannerList, count, err := b.bannerRepository.List(ctx)
	if err != nil {
		return nil, err
	}
	return &vo.BannerList{
		bannerList,
		count,
	}, nil
}

func (b *bannerService) SysList(ctx context.Context, page, pageSize int) ([]*po.Banner, int64, error) {
	bannerList, count, err := b.bannerRepository.SysList(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return bannerList, count, nil
}

func (b *bannerService) Update(ctx context.Context, banner *po.Banner) (*po.Banner, error) {
	return b.bannerRepository.SysUpdate(ctx, banner)
}

func (b *bannerService) Delete(ctx context.Context, id int64) error {
	return b.bannerRepository.SysDelete(ctx, id)
}
