//go:build wireinject
// +build wireinject

package repository

import (
	"blog/pkg/config"
	"blog/pkg/database/es"
	"blog/pkg/database/mongo"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"blog/pkg/logger"
	"blog/pkg/mq/kafka"
	"blog/pkg/trace"
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
