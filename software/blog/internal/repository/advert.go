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

type advertRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newAdvertRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) AdvertRepository {
	return &advertRepository{
		logger: logger.With(zap.String("type", "AdvertRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (a *advertRepository) Create(c context.Context, advert *po.Advert) (*po.Advert, error) {
	if err := a.db.WithContext(c).Create(advert).Error; err != nil {
		return nil, err
	}
	return advert, nil
}

func (a *advertRepository) List(c context.Context, query map[string]interface{}, page, pageSize int) ([]*vo.Advert, int64, error) {
	var adverts []*vo.Advert
	var count int64
	err := a.db.WithContext(c).Model(&po.Advert{}).Where(query).Scopes(database.Paginate(page, pageSize)).Find(&adverts).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get advert list db error: %s", err)
	}
	return adverts, count, nil
}

func (a *advertRepository) SysList(c context.Context, query map[string]interface{}, page, pageSize int) ([]*po.Advert, int64, error) {
	var adverts []*po.Advert
	var count int64
	err := a.db.WithContext(c).Model(&po.Advert{}).Where(query).Scopes(database.Paginate(page, pageSize)).Find(&adverts).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get advert list db error: %s", err)
	}
	return adverts, count, nil
}

func (a *advertRepository) Retrieve(c context.Context, id int64) (*po.Advert, error) {
	advert := new(po.Advert)
	err := a.db.WithContext(c).Where("id = ? ", id).First(advert).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(err, "查看对象不存在在或无权查看[id=%d]", id)
	}
	return advert, nil
}

func (a *advertRepository) Update(c context.Context, advert *po.Advert) (*po.Advert, error) {
	tx := a.db.WithContext(c).Model(advert).Where("id = ?", advert.ID).Updates(advert)
	if err := tx.Error; err != nil {
		return nil, errors.Errorf("db error: %s", err)
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("无权更新该数据")

	}
	return advert, nil
}

func (a *advertRepository) Delete(c context.Context, id int64) error {
	var advert po.Advert
	tx := a.db.WithContext(c).Model(&advert).Delete(&advert, id)
	if err := tx.Error; err != nil {
		return errors.Errorf("error found in delete user: %s", err)
	}
	if tx.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (a *advertRepository) Migrate() error {
	return a.db.AutoMigrate(&po.Advert{})
}
