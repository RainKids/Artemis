package grpcserver

import (
	"blog/api/proto"
	"blog/internal/biz/dto"
	"blog/pkg/tools/timeparse"
	"context"
)

func (s *GrpcServer) AdvertList(c context.Context, req *proto.AdvertListRequest) (*proto.AdvertListResponse, error) {
	list, count, err := s.service.Advert().SysList(c, dto.AdvertSearchParams{req.Title}, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*proto.Advert, 0, count)
	for _, advert := range list {
		result = append(result, &proto.Advert{ID: advert.ID,
			Title:     advert.Title,
			Href:      advert.Href,
			Image:     advert.Image,
			IsShow:    advert.IsShow,
			UpdatedAt: timeparse.GoTimeToPbTime(advert.UpdatedAt),
			CreatedAt: timeparse.GoTimeToPbTime(advert.CreatedAt),
		})
	}
	return &proto.AdvertListResponse{
		Count:  count,
		Result: result,
	}, nil
}
