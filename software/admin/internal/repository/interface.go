package repository

import (
	"context"
)

type Repository interface {
	Hello() HelloRepository
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
