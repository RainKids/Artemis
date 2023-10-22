//go:build wireinject
// +build wireinject

package service

import (
	"blog/internal/repository"
	"blog/pkg/config"
	"blog/pkg/logger"
	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	ProviderSet)

func CreateService(cf string,
	rpo repository.Repository,
) (Service, error) {
	panic(wire.Build(ServiceProviderSet))
}
