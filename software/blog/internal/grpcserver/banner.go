package grpcserver

import (
	bannerPb "blog/api/proto/blog/bannerpb"
	"blog/internal/service"
	"blog/pkg/tools/timeparse"
	"context"
	"go.uber.org/zap"
)

type bannerGrpcServer struct {
	logger  *zap.Logger
	service service.BannerService
	bannerPb.UnimplementedRPCServer
}

func newBannerGrpcServer(logger *zap.Logger, service service.BannerService) *bannerGrpcServer {
	return &bannerGrpcServer{
		logger:  logger.With(zap.String("type", "AdvertService")),
		service: service,
	}
}

func (s *bannerGrpcServer) BannerList(c context.Context, req *bannerPb.BannerListRequest) (*bannerPb.BannerListResponse, error) {
	list, count, err := s.service.SysList(c, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*bannerPb.Banner, 0, count)
	for _, banner := range list {
		result = append(result, &bannerPb.Banner{
			ID:        banner.ID,
			Path:      banner.Path,
			Hash:      banner.Hash,
			Name:      banner.Name,
			ImageType: int64(banner.ImageType),
			CreatedAt: timeparse.GoTimeToPbTime(banner.CreatedAt),
			UpdatedAt: timeparse.GoTimeToPbTime(banner.UpdatedAt),
		})
	}
	return &bannerPb.BannerListResponse{
		Count:  count,
		Result: result,
	}, nil
}
