package service

import (
	"admin/global"
	"admin/internal/biz/dto"
	"admin/internal/biz/po"
	"admin/internal/repository"
	"context"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"strconv"
)

type roleService struct {
	logger         *zap.Logger
	roleRepository repository.RoleRepository
	enforcer       *casbin.SyncedEnforcer
}

func newRoleService(logger *zap.Logger, roleRepository repository.RoleRepository) RoleService {
	return &roleService{
		logger:         logger.With(zap.String("type", "RoleService")),
		roleRepository: roleRepository,
	}
}

func (r *roleService) Create(c context.Context, req dto.RoleRequest) (*po.Role, error) {
	role, err := r.roleRepository.Create(c, &po.Role{
		Name:      req.Name,
		Status:    req.Status,
		Key:       req.Key,
		Sort:      req.Sort,
		Flag:      req.Flag,
		Remark:    req.Remark,
		Admin:     req.Admin,
		DataScope: req.DataScope,
		Menu:      &req.Menu,
		Dept:      req.Dept,
		OperateBy: global.OperateBy{
			CreateBy: req.OperateBy,
			UpdateBy: req.OperateBy,
		},
	})
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0)
	roleID := strconv.Itoa(int(role.ID))
	for _, menu := range *role.Menu {
		for _, api := range menu.Api {
			if mp[roleID+"-"+api.Path+"-"+api.Method] != "" {
				mp[roleID+"-"+api.Path+"-"+api.Method] = ""
				polices = append(polices, []string{roleID, api.Path, api.Method})
			}
		}
	}

	if len(polices) <= 0 {
		return role, nil
	}

	// 写入 sys_casbin_rule 权限表里 当前角色数据的记录
	_, err = r.enforcer.AddNamedPolicies("p", polices)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleService) Delete(c context.Context, id int64) error {
	err := r.roleRepository.Delete(c, id)
	if err != nil {
		return err
	}
	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, err = r.enforcer.RemoveFilteredPolicy(0, strconv.Itoa(int(id)))
	return err
}
