package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type menuRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newMenuRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) MenuRepository {
	return &menuRepository{
		logger: logger.With(zap.String("type", "MenuRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (m *menuRepository) Migrate() error {
	return m.db.AutoMigrate(&po.MenuParam{}, &po.MenuButton{}, &po.Menu{})
}
