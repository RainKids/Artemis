package service

import (
	"admin/api/proto"
	"admin/internal/repository"
	"go.uber.org/zap"
)

type service struct {
	blog  BlogService
	hello HelloService
}

func NewService(log *zap.Logger, repository repository.Repository, blogPpcSvc proto.BlogServiceClient) Service {
	r := &service{
		blog:  newBlogService(log, blogPpcSvc),
		hello: newHelloService(log, repository.Hello()),
	}
	return r
}

func (s *service) Blog() BlogService {
	return s.blog
}

func (s *service) Hello() HelloService {
	return s.hello
}
