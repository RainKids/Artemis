//go:build wireinject
// +build wireinject

package server

import (
	"blog/internal/api"
	"blog/internal/app"
	"blog/internal/crontab"
	"blog/internal/grpcserver"
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/application"
	"blog/pkg/config"
	"blog/pkg/cron"
	"blog/pkg/database/es"
	"blog/pkg/database/mongo"
	"blog/pkg/database/postgres"
	"blog/pkg/database/redis"
	"blog/pkg/logger"
	"blog/pkg/trace"
	"blog/pkg/transport/grpc"
	"blog/pkg/transport/http"
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
