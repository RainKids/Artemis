package grpcserver

import (
	"blog/api/proto"
	"blog/pkg/transport/grpc"
	"github.com/google/wire"
	stdgrpc "google.golang.org/grpc"
)

func CreateInitGrpcServersFn(
	ps *GrpcServer,
) grpc.InitServers {
	return func(s *stdgrpc.Server) {
		proto.RegisterBlogServiceServer(s, ps)
	}
}

// ProviderSet 定义grpc service wire
var ProviderSet = wire.NewSet(NewGrpcServer, CreateInitGrpcServersFn)
