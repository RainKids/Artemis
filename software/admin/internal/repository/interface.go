package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/common"
	"context"
)

type Repository interface {
	Hello() HelloRepository
	Api() ApiRepository
	Casbin() CasbinRepository
	Dept() DeptRepository
	Dict() DictRepository
	Menu() MenuRepository
	Post() PostRepository
	Role() RoleRepository
	User() UserRepository
	Close() error
	Ping(ctx context.Context) error
	Init() error
	Migrant
}

type Migrant interface {
	Migrate() error
}

type HelloRepository interface {
	Hello(context.Context, string) (string, error)
}

type ApiRepository interface {
	Create(c context.Context, api *po.Api) (*po.Api, error)
	List(c context.Context, params *dto.ApiSearchParams, p *common.DataPermission) ([]*po.Api, int64, error)
	Retrieve(c context.Context, id int64, p *common.DataPermission) (*po.Api, error)
	Update(c context.Context, api *po.Api, p *common.DataPermission) error
	Delete(c context.Context, id interface{}, p *common.DataPermission) error
	Migrate() error
}

type CasbinRepository interface {
	UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error
	Migrate() error
}

type DeptRepository interface {
	Migrate() error
}

type DictRepository interface {
	Migrate() error
}

type MenuRepository interface {
	Migrate() error
}

type PostRepository interface {
	List(context.Context, *dto.PostSearchParams) ([]*po.Post, int64, error)
	Retrieve(context.Context, int64) (*po.Post, error)
	Create(context.Context, *po.Post) (*po.Post, error)
	Update(context.Context, *po.Post) error
	Delete(context.Context, int64) error
	Migrate() error
}

type RoleRepository interface {
	Migrate() error
}

type UserRepository interface {
	Migrate() error
}
