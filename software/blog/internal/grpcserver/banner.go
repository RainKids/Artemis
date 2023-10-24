package grpcserver

import (
	"blog/api/proto"
	"blog/pkg/tools/timeparse"
	"context"
)

func (s *GrpcServer) BannerList(c context.Context, req *proto.BannerListRequest) (*proto.BannerListResponse, error) {
	list, count, err := s.service.Banner().SysList(c, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*proto.Banner, 0, count)
	for _, banner := range list {
		result = append(result, &proto.Banner{
			ID:        banner.ID,
			Path:      banner.Path,
			Hash:      banner.Hash,
			Name:      banner.Name,
			ImageType: int64(banner.ImageType),
			CreatedAt: timeparse.GoTimeToPbTime(banner.CreatedAt),
			UpdatedAt: timeparse.GoTimeToPbTime(banner.UpdatedAt),
		})
	}
	return &proto.BannerListResponse{
		Count:  count,
		Result: result,
	}, nil
}
