package service

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/vo"
	"context"
)

type Service interface {
	Blog() BlogService
	Hello() HelloService
}

type BlogService interface {
	AdvertList(c context.Context, req *dto.AdvertParamsRequest) (*vo.AdvertList, error)
	BannerList(c context.Context, page, pageSize int) (*vo.BannerList, error)
}

type HelloService interface {
	Hello(context.Context, string) (string, error)
}
