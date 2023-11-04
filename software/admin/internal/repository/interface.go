package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
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
	SetDeptPage(context.Context, *dto.DeptSearchParams) ([]po.Dept, error)
	SetDeptTree(context.Context, *dto.DeptSearchParams) ([]vo.DeptLabel, error)
	Retrieve(context.Context, int64) (*po.Dept, error)
	Create(context.Context, *po.Dept) (*po.Dept, error)
	Update(context.Context, *po.Dept) error
	Delete(context.Context, int64) error
	Migrate() error
}

type DictRepository interface {
	DataList(context.Context, *dto.DictDataSearchParams) ([]*po.DictData, int64, error)
	DataRetrieve(context.Context, int64) (*po.DictData, error)
	DataCreate(context.Context, *po.DictData) (*po.DictData, error)
	DataUpdate(context.Context, *po.DictData) error
	DataDelete(context.Context, int64) error
	TypeList(context.Context, *dto.DictTypeSearchParams) ([]*po.DictType, int64, error)
	TypeRetrieve(context.Context, int64) (*po.DictType, error)
	TypeCreate(context.Context, *po.DictType) (*po.DictType, error)
	TypeUpdate(context.Context, *po.DictType) error
	TypeDelete(context.Context, int64) error
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
	Create(context.Context, *po.Role) (*po.Role, error)
	Delete(context.Context, int64) error
	Migrate() error
}

type UserRepository interface {
	Migrate() error
}
