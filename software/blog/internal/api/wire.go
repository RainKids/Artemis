//go:build wireinject
// +build wireinject

package api

import (
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/config"
	"blog/pkg/database/redis"
	"blog/pkg/logger"
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

func CreateController(cf string, repo repository.Repository) (*Controller, error) {
	panic(wire.Build(apiProviderSet))
}
