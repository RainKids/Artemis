package service

import (
	"blog/internal/biz/dto"
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"context"
)

type Service interface {
	Advert() AdvertService
}

type AdvertService interface {
	Create(context.Context, *po.Advert) (*po.Advert, error)
	Retrieve(context.Context, int64) (*po.Advert, error)
	List(c context.Context, params dto.AdvertSearchParams, page, pageSize int) (*vo.AdvertList, error)
	SysList(c context.Context, params dto.AdvertSearchParams, page, pageSize int) ([]*po.Advert, int64, error)
	Update(context.Context, *po.Advert) (*po.Advert, error)
	Delete(context.Context, int64) error
}
