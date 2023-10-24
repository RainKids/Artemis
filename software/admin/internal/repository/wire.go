//go:build wireinject
// +build wireinject

package repository

import (
	"admin/pkg/config"
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"admin/pkg/logger"
	"admin/pkg/mq/kafka"
	"admin/pkg/trace"
	"github.com/google/wire"
)

var repositoryProviderSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	mongo.ProviderSet,
	postgres.ProviderSet,
	es.ProviderSet,
	redis.ProviderSet,
	trace.ProviderSet,
	kafka.ProviderSet,
	ProviderSet)

func CreateRepository(f string) (Repository, error) {
	panic(wire.Build(repositoryProviderSet))
}
