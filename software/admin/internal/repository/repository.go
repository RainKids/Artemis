package repository

import (
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db       *gorm.DB
	rdb      *redis.RedisDB
	es       *es.Client
	mongo    *mongo.MongoDB
	hello    HelloRepository
	migrants []Migrant
}

func (r *repository) Init() error {
	return r.Migrate()
}

func NewRepository(log *zap.Logger, db *postgres.DB, rdb *redis.RedisDB, es *es.Client, mongo *mongo.MongoDB) Repository {
	r := &repository{
		db:    db.Postgres,
		rdb:   rdb,
		es:    es,
		mongo: mongo,
		hello: newHelloRepository(log),
	}
	r.migrants = getMigrants()
	err := r.Init()
	if err != nil {
		log.Error("DB Table Migrate Error", zap.Error(err))
	}
	return r
}

func getMigrants(objs ...interface{}) []Migrant {
	var migrants []Migrant
	for _, obj := range objs {
		if m, ok := obj.(Migrant); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}

	if r.rdb != nil {
		if err := r.rdb.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	if err = db.PingContext(ctx); err != nil {
		return err
	}

	if r.rdb == nil {
		return nil
	}
	if _, err := r.rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func (r *repository) Migrate() error {
	for _, m := range r.migrants {
		if err := m.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Hello() HelloRepository {
	return r.hello
}
