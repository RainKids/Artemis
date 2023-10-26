package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type roleRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newRoleRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) RoleRepository {
	return &roleRepository{
		logger: logger.With(zap.String("type", "RoleRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (r *roleRepository) Migrate() error {
	return r.db.AutoMigrate(&po.Role{})
}
