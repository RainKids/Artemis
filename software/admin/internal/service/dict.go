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

type dictService struct {
	logger         *zap.Logger
	dictRepository repository.DictRepository
}

func newDictService(logger *zap.Logger, dictRepository repository.DictRepository) DictService {
	return &dictService{
		logger:         logger.With(zap.String("type", "DictService")),
		dictRepository: dictRepository,
	}
}

func (d *dictService) DataList(c context.Context, params *dto.DictDataSearchParams) (*vo.DictDataList, error) {
	list, count, err := d.dictRepository.DataList(c, params)
	if err != nil {
		return nil, err
	}
	return &vo.DictDataList{
		Result: list,
		Count:  count,
	}, nil
}
func (d *dictService) DataRetrieve(c context.Context, id int64) (*po.DictData, error) {
	return d.dictRepository.DataRetrieve(c, id)
}
func (d *dictService) DataCreate(c context.Context, req *dto.DictDataRequest) (*po.DictData, error) {
	return d.dictRepository.DataCreate(c, &po.DictData{
		Name:       req.Name,
		DictTypeID: req.Type,
		Status:     req.Status,
		Remark:     req.Remark,
		OperateBy:  global.OperateBy{UpdateBy: req.OperateBy},
	})
}
func (d *dictService) DataUpdate(c context.Context, id int64, req *dto.DictDataRequest) error {
	return d.dictRepository.DataUpdate(c, &po.DictData{
		ID:         id,
		Name:       req.Name,
		DictTypeID: req.Type,
		Status:     req.Status,
		Remark:     req.Remark,
		OperateBy:  global.OperateBy{UpdateBy: req.OperateBy},
	})

}
func (d *dictService) DataDelete(c context.Context, id int64) error {
	return d.dictRepository.DataDelete(c, id)
}
func (d *dictService) TypeList(c context.Context, params *dto.DictTypeSearchParams) (*vo.DictTypeList, error) {
	list, count, err := d.dictRepository.TypeList(c, params)
	if err != nil {
		return nil, err
	}
	return &vo.DictTypeList{
		Result: list,
		Count:  count,
	}, nil
}
func (d *dictService) TypeRetrieve(c context.Context, id int64) (*po.DictType, error) {
	return d.dictRepository.TypeRetrieve(c, id)
}
func (d *dictService) TypeCreate(c context.Context, req *dto.DictTypeRequest) (*po.DictType, error) {
	return d.dictRepository.TypeCreate(c, &po.DictType{
		Name:      req.Name,
		Type:      req.Type,
		Status:    req.Status,
		Remark:    req.Remark,
		OperateBy: global.OperateBy{UpdateBy: req.OperateBy},
	})
}
func (d *dictService) TypeUpdate(c context.Context, id int64, req *dto.DictTypeRequest) error {
	return d.dictRepository.TypeUpdate(c, &po.DictType{
		ID:        id,
		Name:      req.Name,
		Type:      req.Type,
		Status:    req.Status,
		Remark:    req.Remark,
		OperateBy: global.OperateBy{UpdateBy: req.OperateBy},
	})
}
func (d *dictService) TypeDelete(c context.Context, id int64) error {
	return d.dictRepository.TypeDelete(c, id)
}
