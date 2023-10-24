package service

import (
	"admin/api/proto"
	"admin/internal/biz/dto"
	"admin/internal/biz/vo"
	"admin/pkg/tools/timeparse"
	"blog/global"
	"context"
	"go.uber.org/zap"
)

type blogService struct {
	logger     *zap.Logger
	blogPpcSvc proto.BlogServiceClient
}

func newBlogService(logger *zap.Logger, blogPpcSvc proto.BlogServiceClient) BlogService {
	return &blogService{
		logger:     logger.With(zap.String("type", "BlogService")),
		blogPpcSvc: blogPpcSvc,
	}
}

func (b *blogService) AdvertList(c context.Context, req *dto.AdvertParamsRequest) (*vo.AdvertList, error) {
	resp, err := b.blogPpcSvc.AdvertList(c, &proto.AdvertListRequest{
		Title:    req.Title,
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	advertArray := make([]*vo.Advert, 0, len(resp.Result))
	for _, v := range resp.Result {
		advertArray = append(advertArray, &vo.Advert{
			ID:     v.ID,
			Title:  v.Title,
			Href:   v.Href,
			Image:  v.Image,
			IsShow: v.IsShow,
			ModelTime: global.ModelTime{CreatedAt: timeparse.PbTimeToGoTime(v.CreatedAt),
				UpdatedAt: timeparse.PbTimeToGoTime(v.UpdatedAt)},
		})
	}
	return &vo.AdvertList{
		Result: advertArray,
		Count:  resp.Count,
	}, nil
}

func (b *blogService) BannerList(c context.Context, page, pageSize int) (*vo.BannerList, error) {

	resp, err := b.blogPpcSvc.BannerList(c, &proto.BannerListRequest{Page: int64(page), PageSize: int64(pageSize)})
	if err != nil {
		return nil, err
	}
	bannerArray := make([]*vo.Banner, 0, len(resp.Result))
	for _, v := range resp.Result {
		bannerArray = append(bannerArray, &vo.Banner{
			ID:        v.ID,
			Path:      v.Path,
			Hash:      v.Hash,
			Name:      v.Name,
			ImageType: global.ImageType(v.ImageType),
			ModelTime: global.ModelTime{CreatedAt: timeparse.PbTimeToGoTime(v.CreatedAt),
				UpdatedAt: timeparse.PbTimeToGoTime(v.UpdatedAt)},
		})
	}
	return &vo.BannerList{
		Result: bannerArray,
		Count:  resp.Count,
	}, nil
}
