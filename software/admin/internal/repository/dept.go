package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type deptRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newDeptRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) DeptRepository {
	return &deptRepository{
		logger: logger.With(zap.String("type", "DeptRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (d *deptRepository) Migrate() error {
	return d.db.AutoMigrate(&po.Dept{})
}
