package repository

import (
	"admin/global"
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
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
	return m.db.AutoMigrate(&po.Menu{}, &po.MenuParam{}, &po.MenuButton{})
}

func (m *menuRepository) Create(c context.Context, menu *po.Menu) (*po.Menu, error) {
	var err error
	tx := m.db.WithContext(c).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Where("id in ?", menu.Apis).Find(&menu.Apis).Error
	if err != nil {
		tx.Rollback()
		m.logger.Error("create dept db error: %s", zap.Error(err))
	}
	err = tx.Create(&menu).Error
	if err != nil {
		tx.Rollback()
		m.logger.Error("create dept db error: %s", zap.Error(err))
	}
	err = m.initPaths(tx, menu)
	if err != nil {
		tx.Rollback()
		m.logger.Error("create dept db error: %s", zap.Error(err))
	}
	tx.Commit()
	return menu, err
}

func (m *menuRepository) initPaths(tx *gorm.DB, menu *po.Menu) error {
	var err error
	var data po.Menu
	parentMenu := new(po.Menu)
	if menu.ParentId != 0 {
		err = tx.Model(&data).First(parentMenu, menu.ParentId).Error
		if err != nil {
			return err
		}
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return err
		}
		menu.Paths = parentMenu.Paths + "/" + strconv.FormatInt(menu.ID, 10)
	} else {
		menu.Paths = "/0/" + strconv.FormatInt(menu.ID, 10)
	}
	err = tx.Model(&data).Where("id = ?", menu.ID).Update("paths", menu.Paths).Error
	return err
}

func menuConstructor(menuList *[]po.Menu, menu po.Menu) po.Menu {
	list := *menuList

	m := make([]po.Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.ID != list[j].ParentId {
			continue
		}
		mi := po.Menu{}
		mi.ID = list[j].ID
		mi.Name = list[j].Name
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.Type = list[j].Type
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentId = list[j].ParentId
		mi.NoCache = list[j].NoCache
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.Visible = list[j].Visible
		mi.CreatedAt = list[j].CreatedAt
		mi.Api = list[j].Api
		mi.IsIframe = list[j].IsIframe
		mi.IsLink = list[j].IsLink
		mi.IsHide = list[j].IsHide
		mi.IsKeepAlive = list[j].IsKeepAlive
		mi.IsAffix = list[j].IsAffix
		mi.Status = list[j].Status
		mi.Children = []po.Menu{}

		if mi.Type != global.Button {
			ms := menuConstructor(menuList, mi)
			m = append(m, ms)
		} else {
			m = append(m, mi)
		}
	}
	menu.Children = m
	return menu
}

func (m *menuRepository) List(c context.Context, params *dto.MenuSearchParams) (menus *[]po.Menu, err error) {
	var menu = make([]po.Menu, 0)
	err = m.list(c, params, &menu)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(menu); i++ {
		if menu[i].ParentId != 0 {
			continue
		}
		menusInfo := menuConstructor(&menu, menu[i])
		*menus = append(*menus, menusInfo)
	}
	return
}

func (m *menuRepository) list(c context.Context, params *dto.MenuSearchParams, menus *[]po.Menu) error {
	db := m.db.WithContext(c).Model(&po.Menu{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	err := db.Preload("Api").Find(menus).Error
	if err != nil {
		return errors.Errorf("get menu list db error: %s", err)
	}
	return nil
}

func (m *menuRepository) Retrieve(c context.Context, id int64) (*po.Menu, error) {
	menu := &po.Menu{}
	err := m.db.WithContext(c).Preload("Api").Where("id = ?", id).First(menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%d]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	apis := make([]int64, 0)
	for _, v := range menu.Api {
		apis = append(apis, v.ID)
	}
	menu.Apis = apis
	return menu, nil
}

func (m *menuRepository) Update(c context.Context, menu *po.Menu) error {
	var err error
	tx := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var apiList = make([]po.Api, 0)
	tx.Preload("Api").First(menu, menu.ID)
	oldPath := menu.Paths
	tx.Where("id in ?", menu.Apis).Find(&apiList)
	err = tx.Model(menu).Association("Api").Delete(menu.Api)
	if err != nil {
		m.logger.Error("delete policy error:%s", zap.Error(err))
		return errors.Errorf("delete policy error:%s", err)
	}
	menu.Api = apiList
	db := tx.Model(&menu).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&menu)
	if err = db.Error; err != nil {
		m.logger.Error("db error:%s", zap.Error(err))
		return errors.Errorf("db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	if err = tx.Save(menu).Error; err != nil {
		m.logger.Error("update dept db error: %s", zap.Error(err))
		return err
	}
	var menuList []po.Menu
	tx.Where("paths like ?", oldPath+"%").Find(&menuList)
	for _, v := range menuList {
		v.Paths = strings.Replace(v.Paths, oldPath, menu.Paths, 1)
		tx.Model(&v).Update("paths", v.Paths)
	}
	return err
}

func (m *menuRepository) Delete(c context.Context, id int64) error {
	var menu po.Menu
	db := m.db.WithContext(c).Model(&menu).Delete(&menu, id)
	err := db.Error
	if err != nil {
		return errors.Errorf("delete menu db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete menu")
	}
	return nil
}

func (m *menuRepository) SetLabel(c context.Context) (menuLabels []vo.MenuLabel, err error) {
	list := make([]po.Menu, 0)
	err = m.db.WithContext(c).Find(&list).Error
	if err != nil {
		m.logger.Error("find menu list error, %s", zap.Error(err))
		return
	}
	menuLabels = make([]vo.MenuLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		label := vo.MenuLabel{}
		label.ID = list[i].ID
		label.Label = list[i].Title
		deptsInfo := menuLabelConstructor(&list, label)

		menuLabels = append(menuLabels, deptsInfo)
	}
	return
}

// menuLabelConstructor 递归构造组织数据
func menuLabelConstructor(menuList *[]po.Menu, dept vo.MenuLabel) vo.MenuLabel {
	list := *menuList

	menuLabels := make([]vo.MenuLabel, 0)
	for j := 0; j < len(list); j++ {

		if dept.ID != list[j].ParentId {
			continue
		}
		mi := vo.MenuLabel{}
		mi.ID = list[j].ID
		mi.Label = list[j].Title
		mi.Children = []vo.MenuLabel{}
		if list[j].Type != "F" {
			ms := menuLabelConstructor(menuList, mi)
			menuLabels = append(menuLabels, ms)
		} else {
			menuLabels = append(menuLabels, mi)
		}
	}
	if len(menuLabels) > 0 {
		dept.Children = menuLabels
	} else {
		dept.Children = nil
	}
	return dept
}

func (m *menuRepository) GetMenuByRoleName(c context.Context, roleName ...string) (list []po.Menu, err error) {

	var role po.Role
	admin := false
	for _, s := range roleName {
		if s == "admin" {
			admin = true
		}
	}
	if len(roleName) > 0 && admin {
		var data []po.Menu
		err = m.db.WithContext(c).Where("type in ('M','C')").Order("sort").Find(&data).Error
		list = data
	} else {
		err = m.db.WithContext(c).Preload("Menu", func(db *gorm.DB) *gorm.DB {
			return db.Where(" type in ('M','C')").Order("sort")
		}).Where("name in ?", roleName).Find(&role).Error
		list = *role.Menu
	}
	return list, err
}

func menuDistinct(menuList []po.Menu) (result []po.Menu) {
	distinctMap := make(map[int64]struct{}, len(menuList))
	for _, menu := range menuList {
		if _, ok := distinctMap[menu.ID]; !ok {
			distinctMap[menu.ID] = struct{}{}
			result = append(result, menu)
		}
	}
	return result
}

func recursiveSetMenu(orm *gorm.DB, menuIds []int64, menus *[]po.Menu) error {
	if len(menuIds) == 0 || menus == nil {
		return nil
	}
	var subMenus []po.Menu
	err := orm.Where(fmt.Sprintf(" type in ('%s', '%s', '%s') and id in ?",
		global.Directory, global.Menu, global.Button), menuIds).Order("sort").Find(&subMenus).Error
	if err != nil {
		return err
	}

	subIds := make([]int64, 0)
	for _, menu := range subMenus {
		if menu.ParentId != 0 {
			subIds = append(subIds, menu.ParentId)
		}
		if menu.Type != global.Button {
			*menus = append(*menus, menu)
		}
	}
	return recursiveSetMenu(orm, subIds, menus)
}

func (m *menuRepository) getByRoleName(c context.Context, roleName string) ([]po.Menu, error) {
	var role po.Role
	var err error
	data := make([]po.Menu, 0)

	if roleName == "admin" {
		err = m.db.WithContext(c).Where(" type in ('M','C') and deleted_at is null").
			Order("sort").
			Find(&data).
			Error
		err = errors.WithStack(err)
	} else {
		role.Key = roleName
		err = m.db.WithContext(c).Model(&role).Where("key = ? ", roleName).Preload("Menu").First(&role).Error

		if role.Menu != nil {
			mIds := make([]int64, 0)
			for _, menu := range *role.Menu {
				mIds = append(mIds, menu.ID)
			}
			if err := recursiveSetMenu(m.db, mIds, &data); err != nil {
				return nil, err
			}

			data = menuDistinct(data)
		}
	}

	sort.Sort(po.MenuSlice(data))
	return data, err
}

func (m *menuRepository) SetMenuRole(c context.Context, roleName string) (menuList []po.Menu, err error) {
	menus, err := m.getByRoleName(c, roleName)
	menuList = make([]po.Menu, 0)
	for i := 0; i < len(menus); i++ {
		if menus[i].ParentId != 0 {
			continue
		}
		menusInfo := menuConstructor(&menus, menus[i])
		menuList = append(menuList, menusInfo)
	}
	return
}
