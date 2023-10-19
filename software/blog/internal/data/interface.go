package data

import (
	"blog/internal/domain/po"
	"context"
)

type Repository interface {
	Advert() AdvertRepository
	Close() error
	Ping(ctx context.Context) error
	Init() error
	Migrant
}

type AdvertRepository interface {
	Create(context.Context, *po.Advert) (*po.Advert, error)
	Retrieve(context.Context, int64) (*po.Advert, error)
	List(c context.Context, title string, page, pageSize int) ([]*po.Advert, int64, error)
	Update(context.Context, *po.Advert) (*po.Advert, error)
	Delete(context.Context, int64) error
	Migrate() error
}

type Migrant interface {
	Migrate() error
}
