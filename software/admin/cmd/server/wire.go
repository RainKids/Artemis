//go:build wireinject
// +build wireinject

package server

import (
	"admin/internal/api"
	"admin/internal/app"
	"admin/internal/crontab"
	"admin/internal/grpcclient"
	"admin/internal/grpcserver"
	"admin/internal/repository"
	"admin/internal/service"
	"admin/pkg/application"
	"admin/pkg/config"
	"admin/pkg/cron"
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"admin/pkg/logger"
	"admin/pkg/trace"
	"admin/pkg/transport/grpc"
	"admin/pkg/transport/http"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	logger.ProviderSet,
	config.ProviderSet,
	mongo.ProviderSet,
	postgres.ProviderSet,
	es.ProviderSet,
	redis.ProviderSet,
	trace.ProviderSet,
	repository.ProviderSet,
	service.ProviderSet,
	grpcclient.ProviderSet,
	grpcserver.ProviderSet,
	http.ProviderSet,
	grpc.ProviderSet,
	cron.ProviderSet,
	crontab.ProviderSet,
	api.ProviderSet,
	app.ProviderSet)

func CreateApp(cf string) (*application.Application, error) {
	panic(wire.Build(providerSet))
}
