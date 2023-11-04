package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/pkg/database"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
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

func (d *dictRepository) DataList(c context.Context, params *dto.DictDataSearchParams) ([]*po.DictData, int64, error) {
	var dictDatas []*po.DictData
	var count int64
	db := d.db.WithContext(c).Model(&po.Post{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Type != 0 {
		db = db.Where("dict_type_id =  ?", params.Type)
	}
	db = db.Where("status = ?", params.Status)
	err := db.Scopes(database.Paginate(params.Page, params.PageSize)).Preload("DictType").Find(&dictDatas).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get dictData list db error: %s", err)
	}
	return dictDatas, count, nil
}
func (d *dictRepository) DataRetrieve(c context.Context, id int64) (*po.DictData, error) {
	dictData := &po.DictData{}
	err := d.db.WithContext(c).Model(dictData).Where("id = ? ", id).Preload("DictType").First(dictData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%d]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	return dictData, nil
}
func (d *dictRepository) DataCreate(c context.Context, dictData *po.DictData) (*po.DictData, error) {
	if err := d.db.WithContext(c).Create(dictData).Error; err != nil {
		return nil, errors.Errorf("create dictData db error: %s", err)
	}
	return dictData, nil
}
func (d *dictRepository) DataUpdate(c context.Context, dictData *po.DictData) error {
	if err := d.db.WithContext(c).Updates(dictData).Error; err != nil {
		return errors.Errorf("update dictData db error: %s", err)
	}
	return nil
}
func (d *dictRepository) DataDelete(c context.Context, id int64) error {
	var dictData po.DictData
	db := d.db.WithContext(c).Model(&dictData).Delete(&dictData, id)
	err := db.Error
	if err != nil {
		return errors.Errorf("delete dictData db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete dictData")
	}
	return nil
}

func (d *dictRepository) TypeList(c context.Context, params *dto.DictTypeSearchParams) ([]*po.DictType, int64, error) {
	var dicTypes []*po.DictType
	var count int64
	db := d.db.WithContext(c).Model(&po.Post{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Type != "" {
		db = db.Where("type LIKE ?", "%"+params.Type+"%")
	}
	db = db.Where("status = ?", params.Status)
	err := db.Scopes(database.Paginate(params.Page, params.PageSize)).Preload("DictType").Find(&dicTypes).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get dictData list db error: %s", err)
	}
	return dicTypes, count, nil
}

func (d *dictRepository) TypeAll(c context.Context) ([]*po.DictType, int64, error) {
	var dicTypes []*po.DictType
	var count int64
	err := d.db.WithContext(c).Model(&po.DictType{}).Find(&dicTypes).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get dictData list db error: %s", err)
	}
	return dicTypes, count, nil

}
func (d *dictRepository) TypeRetrieve(c context.Context, id int64) (*po.DictType, error) {
	dictType := &po.DictType{}
	err := d.db.WithContext(c).Model(dictType).Where("id = ? ", id).Preload("DictType").First(dictType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%d]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	return dictType, nil

}
func (d *dictRepository) TypeCreate(c context.Context, dictType *po.DictType) (*po.DictType, error) {
	if err := d.db.WithContext(c).Create(dictType).Error; err != nil {
		return nil, errors.Errorf("create dictData db error: %s", err)
	}
	return dictType, nil
}

func (d *dictRepository) TypeUpdate(c context.Context, dictType *po.DictType) error {
	if err := d.db.WithContext(c).Updates(dictType).Error; err != nil {
		return errors.Errorf("update dictData db error: %s", err)
	}
	return nil
}
func (d *dictRepository) TypeDelete(c context.Context, id int64) error {
	var dictType po.DictType
	db := d.db.WithContext(c).Model(&dictType).Delete(&dictType, id)
	err := db.Error
	if err != nil {
		return errors.Errorf("delete dictData db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete dictData")
	}
	return nil
}
