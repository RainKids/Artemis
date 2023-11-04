package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/pkg/database"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleRepository struct {
	logger   *zap.Logger
	db       *gorm.DB
	rdb      *redis.RedisDB
	enforcer *casbin.SyncedEnforcer
}

func newRoleRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB, enforcer *casbin.SyncedEnforcer) RoleRepository {
	return &roleRepository{
		logger:   logger.With(zap.String("type", "RoleRepository")),
		db:       db.Postgres,
		rdb:      rdb,
		enforcer: enforcer,
	}
}

func (r *roleRepository) Migrate() error {
	return r.db.AutoMigrate(&po.Role{})
}

// GetRoleMenuId 获取角色对应的菜单ids
func (r *roleRepository) GetRoleMenuId(c context.Context, id int64) ([]int64, error) {
	menuIds := make([]int64, 0)
	role := po.Role{}
	role.ID = id
	if err := r.db.WithContext(c).Model(&role).Preload("SysMenu").First(&role).Error; err != nil {
		return nil, err
	}
	l := *role.Menu
	for i := 0; i < len(l); i++ {
		menuIds = append(menuIds, l[i].ID)
	}
	return menuIds, nil
}

func (r *roleRepository) List(c context.Context, params *dto.RoleSearchParams) (list []*po.Role, count int64, err error) {
	db := r.db.WithContext(c).Model(&po.Role{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	db = db.Where("status = ?", params.Status)
	if err = db.Scopes(database.Paginate(params.Page, params.PageSize)).Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error; err != nil {
		return nil, 0, errors.Errorf("get post list db error: %s", err)
	}
	return
}

func (r *roleRepository) Retrieve(c context.Context, id int64) (*po.Role, error) {
	var role po.Role
	err := r.db.WithContext(c).Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%s]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	role.MenuIds, err = r.GetRoleMenuId(c, role.ID)
	if err != nil {
		r.logger.Error("get menuIds error, %s", zap.Error(err))
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(c context.Context, role *po.Role) (*po.Role, error) {
	var err error
	var menu []po.Menu
	err = r.db.WithContext(c).Preload("Api").Where("menu_id in ?", role.MenuIds).Find(&menu).Error
	if err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return nil, err
	}
	role.Menu = &menu
	tx := r.db.WithContext(c)
	var count int64
	err = tx.Model(role).Where("key = ?", role.Key).Count(&count).Error
	if err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return nil, err
	}

	if count > 0 {
		err = errors.New("roleKey已存在，需更换在提交！")
		r.logger.Error("db error:%s", zap.Error(err))
		return nil, err
	}

	err = tx.Create(role).Error
	if err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return nil, err
	}
	return role, nil
}
func (r *roleRepository) Update(c context.Context, id int64, role *po.Role) error {
	var err error
	tx := r.db.WithContext(c)
	var menuList = make([]po.Menu, 0)
	tx.Preload("Menu").First(&role, role.ID)
	tx.Preload("Api").Where("menu_id in ?", role.MenuIds).Find(&menuList)
	err = tx.Model(role).Association("Menu").Delete(role.Menu)
	if err != nil {
		r.logger.Error("delete policy error:%s", zap.Error(err))
		return err
	}
	role.Menu = &menuList
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(role)

	if err = db.Error; err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

func (r *roleRepository) Delete(c context.Context, id int64) error {
	var err error
	tx := r.db.WithContext(c)
	var role = po.Role{}
	tx.Preload("Menu").Preload("Dept").First(&role, id)
	//删除 SysRole 时，同时删除角色所有 关联其它表 记录 (SysMenu 和 SysMenu)
	db := tx.Select(clause.Associations).Delete(&role)

	if err = db.Error; err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

func (r *roleRepository) UpdateDataScope(c context.Context, req *dto.RoleDataScopeRequest) (*po.Role, error) {
	var err error
	tx := r.db.WithContext(c)

	var deptList = make([]po.Dept, 0)
	var role = po.Role{}
	tx.Preload("Dept").First(&role, req.ID)
	tx.Where("id in ?", req.DeptIds).Find(&deptList)
	// 删除SysRole 和 SysDept 的关联关系
	err = tx.Model(&role).Association("Dept").Delete(&po.Dept{})
	if err != nil {
		r.logger.Error("delete SysDept error:%s", zap.Error(err))
		return nil, err
	}
	role.Dept = deptList
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Model(&role).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&role)
	if err = db.Error; err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return nil, err
	}
	if db.RowsAffected == 0 {
		return nil, errors.New("无权更新该数据")
	}
	return &role, nil
}

func (r *roleRepository) UpdateStatus(c context.Context, req *dto.UpdateStatusRequest) error {
	var err error
	tx := r.db.WithContext(c)
	var role = po.Role{}
	tx.First(&role, req.Id)
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&role)
	if err = db.Error; err != nil {
		r.logger.Error("db error:%s", zap.Error(err))
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}
