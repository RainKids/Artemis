package repository

import (
	"blog/internal/biz/dto"
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"context"
)

type Repository interface {
	Advert() AdvertRepository
	Banner() BannerRepository
	Article() ArticleRepository
	Comment() CommentRepository
	Close() error
	Ping(ctx context.Context) error
	Init() error
	Migrant
}

type AdvertRepository interface {
	Create(context.Context, *po.Advert) (*po.Advert, error)
	Retrieve(context.Context, int64) (*po.Advert, error)
	List(c context.Context, query map[string]interface{}, page, pageSize int) ([]*vo.Advert, int64, error)
	SysList(c context.Context, query map[string]interface{}, page, pageSize int) ([]*po.Advert, int64, error)
	Update(context.Context, *po.Advert) (*po.Advert, error)
	Delete(context.Context, int64) error
	Migrate() error
}

type BannerRepository interface {
	SysCreate(context.Context, *po.Banner) (*po.Banner, error)
	SysRetrieve(context.Context, int64) (*po.Banner, error)
	List(c context.Context) ([]*vo.Banner, int64, error)
	SysList(c context.Context, page, pageSize int) ([]*po.Banner, int64, error)
	SysUpdate(context.Context, *po.Banner) (*po.Banner, error)
	SysDelete(context.Context, int64) error
	Migrate() error
}

type ArticleRepository interface {
	Search(search *dto.ArticleSearchParams) (*vo.ArticleSearchList, error)
	Migrate() error
}

type CommentRepository interface {
	Migrate() error
}

type Migrant interface {
	Migrate() error
}
