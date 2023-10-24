package grpcserver

import (
	"blog/api/proto"
	"blog/internal/service"
	"go.uber.org/zap"
)

type GrpcServer struct {
	logger  *zap.Logger
	service service.Service
	proto.UnimplementedBlogServiceServer
}

func NewGrpcServer(logger *zap.Logger, service service.Service) (*GrpcServer, error) {
	return &GrpcServer{
		logger:  logger,
		service: service,
	}, nil
}
