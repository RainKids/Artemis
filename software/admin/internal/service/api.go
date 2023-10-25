package service

import (
	"admin/global"
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/biz/vo"
	"admin/internal/common"
	"admin/internal/repository"
	"context"
	"go.uber.org/zap"
)

type apiService struct {
	logger           *zap.Logger
	apiRepository    repository.ApiRepository
	casbinRepository repository.CasbinRepository
}

func newApiService(logger *zap.Logger, apiRepository repository.ApiRepository, casbinRepository repository.CasbinRepository) ApiService {
	return &apiService{
		logger:           logger.With(zap.String("type", "AdvertService")),
		apiRepository:    apiRepository,
		casbinRepository: casbinRepository,
	}
}

func (a *apiService) List(c context.Context, params *dto.ApiSearchParams, p *common.DataPermission) (*vo.ApiList, error) {
	list, count, err := a.apiRepository.List(c, params, p)
	if err != nil {
		return nil, err
	}
	return &vo.ApiList{
		Result: list,
		Count:  count,
	}, nil
}
func (a *apiService) Create(c context.Context, req *dto.ApiRequest) (*po.Api, error) {
	return a.apiRepository.Create(c, &po.Api{
		Title:       req.Title,
		Handle:      req.Handle,
		Path:        req.Path,
		Method:      req.Method,
		Type:        req.Type,
		Description: req.Description,
		ApiGroup:    req.ApiGroup,
		OperateBy:   global.OperateBy{CreateBy: req.OperateBy},
	})
}
func (a *apiService) Retrieve(c context.Context, id int64, p *common.DataPermission) (*po.Api, error) {
	return a.apiRepository.Retrieve(c, id, p)
}
func (a *apiService) Update(c context.Context, id int64, req *dto.ApiRequest, p *common.DataPermission) error {
	return a.apiRepository.Update(c, &po.Api{
		ID:          id,
		Title:       req.Title,
		Handle:      req.Handle,
		Path:        req.Path,
		Method:      req.Method,
		Type:        req.Type,
		Description: req.Description,
		ApiGroup:    req.ApiGroup,
		OperateBy:   global.OperateBy{UpdateBy: req.OperateBy},
	}, p)
}
func (a *apiService) Delete(c context.Context, id int64, p *common.DataPermission) error {
	return a.apiRepository.Delete(c, id, p)
}
