package repository

import (
	"blog/internal/biz/po"
	"blog/internal/biz/vo"
	"blog/pkg/database"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type bannerRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newBannerRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) BannerRepository {
	return &bannerRepository{
		logger: logger.With(zap.String("type", "BannerRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (b *bannerRepository) SysCreate(c context.Context, banner *po.Banner) (*po.Banner, error) {
	if err := b.db.WithContext(c).Create(banner).Error; err != nil {
		return nil, err
	}
	return banner, nil
}

func (b *bannerRepository) List(c context.Context) ([]*vo.Banner, int64, error) {
	var banners []*vo.Banner
	var count int64
	err := b.db.WithContext(c).Model(&po.Banner{}).Find(&banners).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get Banner list db error: %s", err)
	}
	return banners, count, nil
}

func (b *bannerRepository) SysList(c context.Context, page, pageSize int) ([]*po.Banner, int64, error) {
	var banners []*po.Banner
	var count int64
	err := b.db.WithContext(c).Model(&po.Banner{}).Find(&banners).Scopes(database.Paginate(page, pageSize)).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get Banner list db error: %s", err)
	}
	return banners, count, nil
}

func (b *bannerRepository) SysRetrieve(c context.Context, id int64) (*po.Banner, error) {
	banner := new(po.Banner)
	err := b.db.WithContext(c).Where("id = ? ", id).First(banner).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(err, "查看对象不存在在或无权查看[id=%d]", id)
	}
	return banner, nil
}

func (b *bannerRepository) SysUpdate(c context.Context, banner *po.Banner) (*po.Banner, error) {
	tx := b.db.WithContext(c).Model(banner.TableName()).Where("id = ?", banner.ID).Updates(banner)
	if err := tx.Error; err != nil {
		return nil, errors.Errorf("db error: %s", err)
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("无权更新该数据")

	}
	return banner, nil
}

func (b *bannerRepository) SysDelete(c context.Context, id int64) error {
	var Banner po.Banner
	tx := b.db.WithContext(c).Model(&Banner).Delete(&Banner, id)
	if err := tx.Error; err != nil {
		return errors.Errorf("error found in delete user: %s", err)
	}
	if tx.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (b *bannerRepository) Migrate() error {
	return b.db.AutoMigrate(&po.Banner{})
}
