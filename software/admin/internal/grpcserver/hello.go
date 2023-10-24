package grpcserver

import (
	"admin/api/proto"
	"context"
)

func (s *GrpcServer) Hello(c context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	msg, err := s.service.Hello().Hello(c, req.Greeting)
	return &proto.HelloResponse{
		Reply: msg,
	}, err
}
