//go:build wireinject
// +build wireinject

package service

import (
	"admin/internal/grpcclient"
	"admin/internal/repository"
	"admin/pkg/config"
	"admin/pkg/logger"
	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	ProviderSet)

func CreateService(cf string,
	rpo repository.Repository,
	blogRpcSvc *grpcclient.BlogClient,
) (Service, error) {
	panic(wire.Build(ServiceProviderSet))
}
