//go:build wireinject
// +build wireinject

package crontab

import (
	"blog/pkg/config"
	"blog/pkg/cron"
	"blog/pkg/database/es"
	"blog/pkg/database/mongo"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"blog/pkg/logger"
	"github.com/google/wire"
)

var testProviderSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	redis.ProviderSet,
	es.ProviderSet,
	postgres.ProviderSet,
	mongo.ProviderSet,
	cron.ProviderSet,
	ProviderSet)

func CreateDefaultCronJobService(cf string,
) (*DefaultCronJobService, error) {
	panic(wire.Build(testProviderSet))
}
