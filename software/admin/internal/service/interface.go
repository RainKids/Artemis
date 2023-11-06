package service

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/internal/common"
	"context"
)

type Service interface {
	Blog() BlogService
	Hello() HelloService
	Api() ApiService
	Dept() DeptService
	Dict() DictService
	Menu() MenuService
	Post() PostService
	Role() RoleService
	User() UserService
}

type BlogService interface {
	AdvertList(c context.Context, req *dto.AdvertParamsRequest) (*vo.AdvertList, error)
	BannerList(c context.Context, page, pageSize int) (*vo.BannerList, error)
}

type HelloService interface {
	Hello(context.Context, string) (string, error)
}

type ApiService interface {
	List(context.Context, *dto.ApiSearchParams, *common.DataPermission) (*vo.ApiList, error)
	Create(context.Context, *dto.ApiRequest) (*po.Api, error)
	Retrieve(context.Context, int64, *common.DataPermission) (*po.Api, error)
	Update(context.Context, int64, *dto.ApiRequest, *common.DataPermission) error
	Delete(context.Context, int64, *common.DataPermission) error
}

type DeptService interface {
	SetDeptPage(context.Context, *dto.DeptSearchParams) (*vo.DeptList, error)
	SetDeptTree(context.Context, *dto.DeptSearchParams) (*vo.DeptLabelList, error)
	Retrieve(context.Context, int64) (*po.Dept, error)
	Create(context.Context, *dto.DeptRequest) (*po.Dept, error)
	Update(context.Context, int64, *dto.DeptRequest) error
	Delete(context.Context, int64) error
}

type DictService interface {
	DataList(context.Context, *dto.DictDataSearchParams) (*vo.DictDataList, error)
	DataRetrieve(context.Context, int64) (*po.DictData, error)
	DataCreate(context.Context, *dto.DictDataRequest) (*po.DictData, error)
	DataUpdate(context.Context, int64, *dto.DictDataRequest) error
	DataDelete(context.Context, int64) error
	TypeList(context.Context, *dto.DictTypeSearchParams) (*vo.DictTypeList, error)
	TypeRetrieve(context.Context, int64) (*po.DictType, error)
	TypeCreate(context.Context, *dto.DictTypeRequest) (*po.DictType, error)
	TypeUpdate(context.Context, int64, *dto.DictTypeRequest) error
	TypeDelete(context.Context, int64) error
}

type MenuService interface {
}

type PostService interface {
	List(c context.Context, params *dto.PostSearchParams) (*vo.PostList, error)
	Create(context.Context, *dto.PostRequest) (*po.Post, error)
	Retrieve(context.Context, int64) (*po.Post, error)
	Update(context.Context, int64, *dto.PostRequest) error
	Delete(context.Context, int64) error
}

type RoleService interface {
	Create(c context.Context, req dto.RoleRequest) (*po.Role, error)
	Delete(context.Context, int64) error
}

type UserService interface {
}
