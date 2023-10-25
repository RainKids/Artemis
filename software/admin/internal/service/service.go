package service

import (
	"admin/internal/grpcclient"
	"admin/internal/repository"
	"go.uber.org/zap"
)

type service struct {
	blog  BlogService
	hello HelloService
	api   ApiService
}

func NewService(log *zap.Logger, repository repository.Repository, blogPpcSvc *grpcclient.BlogClient) Service {
	r := &service{
		blog:  newBlogService(log, blogPpcSvc),
		hello: newHelloService(log, repository.Hello()),
		api:   newApiService(log, repository.Api(), repository.Casbin()),
	}
	return r
}

func (s *service) Blog() BlogService {
	return s.blog
}

func (s *service) Hello() HelloService {
	return s.hello
}

func (s *service) Api() ApiService {
	return s.api
}
