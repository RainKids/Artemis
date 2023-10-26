package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newUserRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) UserRepository {
	return &userRepository{
		logger: logger.With(zap.String("type", "UserRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (p *userRepository) Migrate() error {
	return p.db.AutoMigrate(&po.User{}, &po.UserToken{})
}
