//go:build wireinject
// +build wireinject

package api

import (
	"admin/api/proto"
	"admin/internal/repository"
	"admin/internal/service"
	"admin/pkg/config"
	"admin/pkg/database/redis"
	"admin/pkg/logger"
	"github.com/google/wire"
)

var apiProviderSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	service.ProviderSet,
	//repositories.ProviderSet,
	ProviderSet,
)

func CreateController(cf string, repo repository.Repository, blogRpcSrv proto.BlogServiceClient) (*Controller, error) {
	panic(wire.Build(apiProviderSet))
}
