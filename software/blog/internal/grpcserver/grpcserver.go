package grpcserver

import (
	"blog/internal/service"
	"go.uber.org/zap"
)

type GrpcServer struct {
	advert *advertGrpcServer
	banner *bannerGrpcServer
}

func NewGrpcServer(logger *zap.Logger, service service.Service) *GrpcServer {
	return &GrpcServer{
		advert: newAdvertGrpcServer(logger, service.Advert()),
		banner: newBannerGrpcServer(logger, service.Banner()),
	}
}
