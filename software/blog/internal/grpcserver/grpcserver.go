package grpcserver

import (
	"blog/api/proto/advertpb"
	"blog/internal/service"
	"go.uber.org/zap"
)

type GrpcServer struct {
	logger  *zap.Logger
	service service.Service
	advertPb.UnimplementedRPCServer
}

func NewGrpcServer(logger *zap.Logger, ps service.Service) *GrpcServer {
	return &GrpcServer{
		logger:  logger,
		service: ps,
	}
}
