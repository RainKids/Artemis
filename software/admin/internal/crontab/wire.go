//go:build wireinject
// +build wireinject

package crontab

import (
	"admin/pkg/config"
	"admin/pkg/cron"
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"admin/pkg/logger"
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
