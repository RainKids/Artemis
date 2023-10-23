package grpcserver

import (
	"blog/api/proto/blog/advertpb"
	"blog/internal/biz/dto"
	"blog/internal/service"
	"blog/pkg/tools/timeparse"
	"context"
	"go.uber.org/zap"
)

type advertGrpcServer struct {
	logger  *zap.Logger
	service service.AdvertService
	advertPb.UnimplementedRPCServer
}

func newAdvertGrpcServer(logger *zap.Logger, service service.AdvertService) *advertGrpcServer {
	return &advertGrpcServer{
		logger:  logger.With(zap.String("type", "AdvertService")),
		service: service,
	}
}

func (s *advertGrpcServer) AdvertList(c context.Context, req *advertPb.AdvertListRequest) (*advertPb.AdvertListResponse, error) {
	list, count, err := s.service.SysList(c, dto.AdvertSearchParams{req.Title}, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*advertPb.Advert, 0, count)
	for _, advert := range list {
		result = append(result, &advertPb.Advert{ID: advert.ID,
			Title:     advert.Title,
			Href:      advert.Href,
			Image:     advert.Image,
			IsShow:    advert.IsShow,
			UpdatedAt: timeparse.GoTimeToPbTime(advert.UpdatedAt),
			CreatedAt: timeparse.GoTimeToPbTime(advert.CreatedAt),
		})
	}
	return &advertPb.AdvertListResponse{
		Count:  count,
		Result: result,
	}, nil
}
