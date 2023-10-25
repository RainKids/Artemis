package service

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/internal/common"
	"context"
)

type Service interface {
	Blog() BlogService
	Hello() HelloService
	Api() ApiService
}

type BlogService interface {
	AdvertList(c context.Context, req *dto.AdvertParamsRequest) (*vo.AdvertList, error)
	BannerList(c context.Context, page, pageSize int) (*vo.BannerList, error)
}

type HelloService interface {
	Hello(context.Context, string) (string, error)
}

type ApiService interface {
	List(context.Context, *dto.ApiSearchParams, *common.DataPermission) (*vo.ApiList, error)
	Create(context.Context, *dto.ApiRequest) (*po.Api, error)
	Retrieve(context.Context, int64, *common.DataPermission) (*po.Api, error)
	Update(context.Context, int64, *dto.ApiRequest, *common.DataPermission) error
	Delete(context.Context, int64, *common.DataPermission) error
}
