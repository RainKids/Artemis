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
