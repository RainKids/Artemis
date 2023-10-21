package repository

import (
	"blog/pkg/database/es"
	"blog/pkg/database/mongo"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db       *gorm.DB
	rdb      *redis.RedisDB
	es       *es.Client
	mongo    *mongo.MongoDB
	advert   AdvertRepository
	migrants []Migrant
}

func (r *repository) Init() error {
	//TODO implement me
	return nil
}

func NewRepository(log *zap.Logger, db *postgres.DB, rdb *redis.RedisDB, es *es.Client, mongo *mongo.MongoDB) Repository {
	r := &repository{
		db:     db.Postgres,
		rdb:    rdb,
		es:     es,
		mongo:  mongo,
		advert: newAdvertRepository(log, db, rdb),
	}
	r.migrants = getMigrants(
		r.advert,
	)

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

func (r *repository) Advert() AdvertRepository {
	return r.advert
}