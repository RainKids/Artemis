package repository

import (
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/pkg/database"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type deptRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.RedisDB
}

func newDeptRepository(logger *zap.Logger, db *postgres.DB, rdb *redis.RedisDB) DeptRepository {
	return &deptRepository{
		logger: logger.With(zap.String("type", "DeptRepository")),
		db:     db.Postgres,
		rdb:    rdb,
	}
}

func (d *deptRepository) Migrate() error {
	return d.db.AutoMigrate(&po.Dept{})
}

func (d *dictRepository) List(c context.Context, params *dto.DeptSearchParams) ([]*po.Dept, int64, error) {
	var depts []*po.Dept
	var count int64
	db := d.db.WithContext(c).Model(&po.Dept{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	db = db.Where("status = ?", params.Status)
	err := db.Scopes(database.Paginate(params.Page, params.PageSize)).Find(&depts).Count(&count).Error
	if err != nil {
		return nil, 0, errors.Errorf("get dictData list db error: %s", err)
	}
	return depts, count, nil
}

func (d *deptRepository) list(c context.Context, params *dto.DeptSearchParams, deptList *[]po.Dept) error {
	db := d.db.WithContext(c).Model(&po.Dept{})
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	db = db.Where("status = ?", params.Status)
	err := db.Find(deptList).Error
	if err != nil {
		return errors.Errorf("get dictData list db error: %s", err)
	}
	return nil
}

func (d *deptRepository) Create(c context.Context, dept *po.Dept) (*po.Dept, error) {
	var err error
	tx := d.db.WithContext(c).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err = tx.Create(dept).Error; err != nil {
		return nil, errors.Errorf("create dept db error: %s", err)
	}
	path := strconv.FormatInt(dept.ID, 10) + "/"
	if dept.ParentId != 0 {
		var deptP po.Dept
		tx.First(&deptP, dept.ParentId)
		path = deptP.Path + path
	} else {
		path = "/0/" + path
	}
	var mp = map[string]string{}
	mp["dept_path"] = path
	if err = tx.Model(&dept).Update("path", path).Error; err != nil {
		d.logger.Error("create dept db to update dept_path error: %s", zap.Error(err))
		return nil, err
	}
	return dept, nil
}

func (d *deptRepository) Retrieve(c context.Context, id int64) (*po.Dept, error) {
	dept := &po.Dept{}
	err := d.db.WithContext(c).Model(dept).Where("id = ? ", id).First(dept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "查看对象不存在[id=%d]", id)
		}
		return nil, errors.Errorf("db error: %s", err)
	}
	return dept, nil
}

func (d *deptRepository) Update(c context.Context, dept *po.Dept) error {
	var err error
	tx := d.db.WithContext(c).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	path := strconv.FormatInt(dept.ID, 10) + "/"
	if dept.ParentId != 0 {
		var deptP po.Dept
		tx.First(&deptP, dept.ParentId)
		path = deptP.Path + path
	} else {
		path = "/0/" + path
	}
	dept.Path = path
	if err = tx.Save(dept).Error; err != nil {
		d.logger.Error("update dept db error: %s", zap.Error(err))
		return err
	}
	return nil
}

func (d *deptRepository) Delete(c context.Context, id int64) error {
	var dept po.Dept
	db := d.db.WithContext(c).Model(&dept).Delete(&dept, id)
	err := db.Error
	if err != nil {
		return errors.Errorf("delete dept db error: %s", err)
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete dept")
	}
	return nil
}

func (d *deptRepository) SetDeptTree(c context.Context, params *dto.DeptSearchParams) (m []vo.DeptLabel, err error) {
	var list []po.Dept
	err = d.list(c, params, &list)
	if err != nil {
		return
	}
	m = make([]vo.DeptLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := vo.DeptLabel{}
		e.ID = list[i].ID
		e.Label = list[i].Name
		deptsInfo := deptTreeConstructor(&list, e)
		m = append(m, deptsInfo)
	}
	return
}

func deptTreeConstructor(deptList *[]po.Dept, dept vo.DeptLabel) vo.DeptLabel {
	list := *deptList
	deptLabelArray := make([]vo.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.ID != list[j].ParentId {
			continue
		}
		mi := vo.DeptLabel{ID: list[j].ID, Label: list[j].Name, Children: []vo.DeptLabel{}}
		ms := deptTreeConstructor(deptList, mi)
		deptLabelArray = append(deptLabelArray, ms)
	}
	dept.Children = deptLabelArray
	return dept
}

func (d *deptRepository) SetDeptPage(c context.Context, params *dto.DeptSearchParams) (m []po.Dept, err error) {
	var list []po.Dept
	err = d.list(c, params, &list)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := deptPageConstructor(&list, list[i])
		m = append(m, info)
	}
	return
}

func deptPageConstructor(deptList *[]po.Dept, menu po.Dept) po.Dept {
	list := *deptList
	deptArray := make([]po.Dept, 0)
	for j := 0; j < len(list); j++ {
		if menu.ID != list[j].ParentId {
			continue
		}
		mi := po.Dept{}
		mi.ID = list[j].ID
		mi.ParentId = list[j].ParentId
		mi.Path = list[j].Path
		mi.Name = list[j].Name
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []po.Dept{}
		ms := deptPageConstructor(deptList, mi)
		deptArray = append(deptArray, ms)
	}
	menu.Children = deptArray
	return menu
}

func (d *deptRepository) GetWithRoleId(c context.Context, roleId int) ([]int, error) {
	deptIds := make([]int, 0)
	deptList := make([]vo.DeptId, 0)
	if err := d.db.WithContext(c).Table("sys_role_dept").
		Select("sys_role_dept.dept_id").
		Joins("LEFT JOIN sys_dept on sys_dept.id=sys_role_dept.dept_id").
		Where("role_id = ? ", roleId).
		Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.id=sys_role_dept.dept_id where role_id =? )", roleId).
		Find(&deptList).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].ID)
	}
	return deptIds, nil
}

func (d *deptRepository) SetDeptLabel(c context.Context) (m []vo.DeptLabel, err error) {
	list := make([]po.Dept, 0)
	err = d.db.WithContext(c).Find(&list).Error
	if err != nil {
		d.logger.Error("find dept list error, %s", zap.Error(err))
		return
	}
	m = make([]vo.DeptLabel, 0)
	var item vo.DeptLabel
	for i := range list {
		if list[i].ParentId != 0 {
			continue
		}
		item = vo.DeptLabel{}
		item.ID = list[i].ID
		item.Label = list[i].Name
		deptInfo := deptLabelConstructor(&list, item)
		m = append(m, deptInfo)
	}
	return
}

func deptLabelConstructor(deptList *[]po.Dept, dept vo.DeptLabel) vo.DeptLabel {
	list := *deptList
	var mi vo.DeptLabel
	deptLabelArray := make([]vo.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.ID != list[j].ParentId {
			continue
		}
		mi = vo.DeptLabel{ID: list[j].ID, Label: list[j].Name, Children: []vo.DeptLabel{}}
		ms := deptLabelConstructor(deptList, mi)
		deptLabelArray = append(deptLabelArray, ms)
	}
	dept.Children = deptLabelArray
	return dept
}
