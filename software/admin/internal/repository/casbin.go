package repository

import (
	"admin/internal/biz/po"
	"admin/pkg/database/postgres"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type casbinRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func newCasbinRepository(logger *zap.Logger, db *postgres.DB) CasbinRepository {
	return &casbinRepository{
		logger: logger.With(zap.String("type", "ApiRepository")),
		db:     db.Postgres,
	}
}

func (c *casbinRepository) UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error {
	err := c.db.WithContext(ctx).Table("casbin_rule").Model(&po.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

func (c *casbinRepository) Migrate() error {
	return c.db.AutoMigrate(&po.CasbinRule{})
}
