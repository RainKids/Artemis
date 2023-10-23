package service

import (
	"blog/internal/biz/dto"
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"context"
)

type Service interface {
	Advert() AdvertService
	Banner() BannerService
}

type AdvertService interface {
	Create(context.Context, *po.Advert) (*po.Advert, error)
	Retrieve(context.Context, int64) (*po.Advert, error)
	List(c context.Context, params dto.AdvertSearchParams, page, pageSize int) (*vo.AdvertList, error)
	SysList(c context.Context, params dto.AdvertSearchParams, page, pageSize int) ([]*po.Advert, int64, error)
	Update(context.Context, *po.Advert) (*po.Advert, error)
	Delete(context.Context, int64) error
}

type BannerService interface {
	Create(context.Context, *po.Banner) (*po.Banner, error)
	Retrieve(context.Context, int64) (*po.Banner, error)
	List(c context.Context) (*vo.BannerList, error)
	SysList(c context.Context, page, pageSize int) ([]*po.Banner, int64, error)
	Update(context.Context, *po.Banner) (*po.Banner, error)
	Delete(context.Context, int64) error
}
