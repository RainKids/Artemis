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
	dept  DeptService
	dict  DictService
	menu  MenuService
	post  PostService
	role  RoleService
	user  UserService
}

func NewService(log *zap.Logger, repository repository.Repository, blogPpcSvc *grpcclient.BlogClient) Service {
	r := &service{
		blog:  newBlogService(log, blogPpcSvc),
		hello: newHelloService(log, repository.Hello()),
		api:   newApiService(log, repository.Api(), repository.Casbin()),
		dept:  newDeptService(log, repository.Dept()),
		dict:  newDictService(log, repository.Dict()),
		menu:  newMenuService(log, repository.Menu()),
		post:  newPostService(log, repository.Post()),
		role:  newRoleService(log, repository.Role()),
		user:  newUserService(log, repository.User()),
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

func (r *service) Dept() DeptService {
	return r.dept
}

func (r *service) Dict() DictService {
	return r.dict
}

func (r *service) Menu() MenuService {
	return r.menu
}

func (r *service) Post() PostService {
	return r.post
}

func (r *service) Role() RoleService {
	return r.role
}

func (r *service) User() UserService {
	return r.user
}
