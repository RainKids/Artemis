package grpcclient

import (
	"admin/api/proto"
	"admin/pkg/transport/grpc"
	"go.uber.org/zap"
	grpcSrv "google.golang.org/grpc"
	"time"
)

type BlogClient struct {
	Client proto.BlogServiceClient
}

func NewBlogClient(client *grpc.Client) (*BlogClient, error) {
	var conn *grpcSrv.ClientConn
	c := new(BlogClient)

	//conn, err := client.DialInsecure("blog-server-grpc")
	//
	//c := proto.NewBlogServiceClient(conn)

	go func(*grpcSrv.ClientConn, *BlogClient) *BlogClient {
		for {
			conn, err := client.DialInsecure("blog-server-grpc")
			if err == nil {
				c.Client = proto.NewBlogServiceClient(conn)
				return c
			} else {
				client.Logger.Error("blog grpc server connect failed", zap.Error(err))
			}
			time.Sleep(time.Second * 5)
		}
	}(conn, c)

	return c, nil
}
