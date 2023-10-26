package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/common"
	"admin/pkg/database"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type apiRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newApiRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) ApiRepository {
	return &apiRepository{
		logger: logger.With(zap.String("type", "ApiRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (a *apiRepository) List(c context.Context, params *dto.ApiSearchParams,
	p *common.DataPermission) ([]*po.Api, int64, error) {
	var apis []*po.Api
	var count int64
	db := a.db.WithContext(c).Model(&po.Api{})
	if params.Path != "" {
		db = db.Where("path LIKE ?", "%"+params.Path+"%")
	}

	if params.Title != "" {
		db = db.Where("description LIKE ?", "%"+params.Title+"%")
	}

	if params.Method != "" {
		db = db.Where("method = ?", params.Method)
	}

	if params.ApiGroup != "" {
		db = db.Where("api_group = ?", params.ApiGroup)
	}
	orderby := ""
	if params.OrderKey != "" {
		// 设置有效排序key 防止sql注入
		orderMap := make(map[string]bool, 5)
		orderMap["id"] = true
		orderMap["path"] = true
		orderMap["api_group"] = true
		orderMap["description"] = true
		orderMap["method"] = true
		if orderMap[params.OrderKey] {
			if params.Desc {
				orderby = params.OrderKey + " desc"
			} else {
				orderby = params.OrderKey
			}
		}
	} else {
		orderby = "id"
	}
	err := db.Scopes(database.Paginate(params.Page, params.PageSize), common.Permission(po.Api{}.TableName(), p)).Order(orderby).
		Find(&apis).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get api list db error: %s", err)
	}
	return apis, count, nil
}

func (a *apiRepository) Create(c context.Context, api *po.Api) (*po.Api, error) {
	err := a.db.WithContext(c).Where("id != ? and path = ? and method =  ", api.ID, api.Path, api.Method).First(api).Error
	if err != nil {
		return nil, errors.New("存在相同api路径")
	}
	err = a.db.WithContext(c).Model(api.TableName()).Create(api).Error
	if err != nil {
		return nil, errors.Errorf("api create db error: %s", err)
	}

	return api, nil
}

func (a *apiRepository) Retrieve(c context.Context, id int64, p *common.DataPermission) (*po.Api, error) {
	api := &po.Api{}
	err := a.db.WithContext(c).Where("id = ? ", id).Scopes(common.Permission(api.TableName(), p)).First(api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%d]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	return api, nil
}

func (a *apiRepository) Update(c context.Context, api *po.Api, p *common.DataPermission) error {
	err := a.db.WithContext(c).Where("id != ? and path = ? and method =  ", api.ID, api.Path, api.Method).First(api).Error
	if err != nil {
		return errors.New("存在相同api路径")
	}
	DB := a.db.WithContext(c).Model(api.TableName()).Where("id = ?", api.ID).Scopes(common.Permission(api.TableName(), p))
	if err := DB.Error; err != nil {
		return errors.Errorf("Service UpdateSysUser error: %s", err)
	}
	if DB.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	DB = a.db.WithContext(c).Updates(api)
	if err := DB.Error; err != nil {
		return errors.Errorf("db error: %s", err)
	}
	if DB.RowsAffected == 0 {
		return errors.New("db update error")
	}
	return nil
}

func (a *apiRepository) Delete(c context.Context, id interface{}, p *common.DataPermission) error {
	var api po.Api
	DB := a.db.WithContext(c).Model(&api).Scopes(common.Permission(api.TableName(), p)).Delete(&api, id)
	if err := DB.Error; err != nil {
		return errors.Errorf("error found in delete user: %s", err)
	}
	if DB.RowsAffected == 0 {
		return errors.New("no permission to delete api")
	}
	return nil
}

func (a *apiRepository) Migrate() error {
	return a.db.AutoMigrate(&po.Api{})
}
