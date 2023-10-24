package grpcclient

import (
	"admin/api/proto"
	"admin/pkg/transport/grpc"
	"github.com/pkg/errors"
)

func NewBlogClient(client *grpc.Client) (proto.BlogServiceClient, error) {
	conn, err := client.DialInsecure("blog-server-grpc")
	if err != nil {
		return nil, errors.Wrap(err, "blog grpc client dial error")
	}
	c := proto.NewBlogServiceClient(conn)
	return c, nil
}
