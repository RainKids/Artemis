package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type dictRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newDictRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) DictRepository {
	return &dictRepository{
		logger: logger.With(zap.String("type", "DictRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (d *dictRepository) Migrate() error {
	return d.db.AutoMigrate(&po.DictData{}, &po.DictType{})
}
