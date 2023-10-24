package grpcserver

import (
	"admin/api/proto"
	"admin/internal/service"
	"go.uber.org/zap"
)

type GrpcServer struct {
	logger  *zap.Logger
	service service.Service
	proto.UnimplementedAdminServiceServer
}

func NewGrpcServer(logger *zap.Logger, service service.Service) (*GrpcServer, error) {
	return &GrpcServer{
		logger:  logger,
		service: service,
	}, nil
}
