package service

import (
	"admin/global"
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/internal/repository"
	"context"
	"go.uber.org/zap"
)

type deptService struct {
	logger         *zap.Logger
	deptRepository repository.DeptRepository
}

func newDeptService(logger *zap.Logger, deptRepository repository.DeptRepository) DeptService {
	return &deptService{
		logger:         logger.With(zap.String("type", "DeptService")),
		deptRepository: deptRepository,
	}
}

func (d *deptService) SetDeptPage(c context.Context, params *dto.DeptSearchParams) (*vo.DeptList, error) {
	list, err := d.deptRepository.SetDeptPage(c, params)
	if err != nil {
		return nil, err
	}
	return &vo.DeptList{
		Result: list,
	}, nil
}

func (d *deptService) SetDeptTree(c context.Context, params *dto.DeptSearchParams) (*vo.DeptLabelList, error) {
	list, err := d.deptRepository.SetDeptTree(c, params)
	if err != nil {
		return nil, err
	}
	return &vo.DeptLabelList{
		Result: list,
	}, nil
}

func (d *deptService) Retrieve(c context.Context, id int64) (*po.Dept, error) {
	return d.deptRepository.Retrieve(c, id)
}

func (d *deptService) Create(c context.Context, req *dto.DeptRequest) (*po.Dept, error) {
	return d.deptRepository.Create(c, &po.Dept{
		ParentId:  req.ParentId,
		Name:      req.Name,
		Sort:      req.Sort,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		OperateBy: global.OperateBy{UpdateBy: req.OperateBy},
	})
}

func (d *deptService) Update(c context.Context, id int64, req *dto.DeptRequest) error {
	return d.deptRepository.Update(c, &po.Dept{
		ID:        id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Sort:      req.Sort,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		OperateBy: global.OperateBy{UpdateBy: req.OperateBy},
	})
}

func (d *deptService) Delete(c context.Context, id int64) error {
	return d.deptRepository.Delete(c, id)
}
