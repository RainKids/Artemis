package grpcserver

import (
	advertPb "blog/api/proto/blog/advertpb"
	bannerPb "blog/api/proto/blog/bannerpb"
	"blog/pkg/transport/grpc"
	"github.com/google/wire"
	stdgrpc "google.golang.org/grpc"
)

func CreateInitGrpcServersFn(
	ps *GrpcServer,
) grpc.InitServers {
	return func(s *stdgrpc.Server) {
		advertPb.RegisterRPCServer(s, ps.advert)
		bannerPb.RegisterRPCServer(s, ps.banner)
	}
}

// ProviderSet 定义grpc service wire
var ProviderSet = wire.NewSet(NewGrpcServer, CreateInitGrpcServersFn)
