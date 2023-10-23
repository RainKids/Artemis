// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func CreateApp(cf string) (*application.Application, error) {
	viper, err2 := config.New(cf)
	if err2 != nil {
		return nil, err
	}
	options, err2 := logger.NewOptions(viper)
	if err2 != nil {
		return nil, err
	}
	zapLogger, err2 := logger.New(options)
	if err2 != nil {
		return nil, err
	}
	appOptions, err2 := app.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	httpOptions, err2 := http.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	redisOptions, err2 := redis.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	redisDB, err2 := redis.New(redisOptions)
	if err2 != nil {
		return nil, err
	}
	postgresOptions, err2 := postgres.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	db, err2 := postgres.New(postgresOptions)
	if err2 != nil {
		return nil, err
	}
	esOptions, err2 := es.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	client, err2 := es.New(esOptions)
	if err2 != nil {
		return nil, err
	}
	mongoOptions, err2 := mongo.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	mongoDB, err2 := mongo.New(mongoOptions)
	if err2 != nil {
		return nil, err
	}
	repositoryRepository := repository.NewRepository(zapLogger, db, redisDB, client, mongoDB)
	serviceService := service.NewService(zapLogger, repositoryRepository)
	controller := api.NewController(zapLogger, redisDB, serviceService)
	initControllers := api.CreateInitControllersFn(controller)
	traceOptions, err2 := trace.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	tracerProvider, err2 := trace.New(traceOptions)
	if err2 != nil {
		return nil, err
	}
	engine := http.NewRouter(httpOptions, zapLogger, initControllers, tracerProvider)
	httpServer, err2 := http.New(httpOptions, zapLogger, engine)
	if err2 != nil {
		return nil, err
	}
	serverOptions, err2 := grpc.NewServerOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	grpcServer := grpcserver.NewGrpcServer(zapLogger, serviceService)
	initServers := grpcserver.CreateInitGrpcServersFn(grpcServer)
	server2, err2 := grpc.NewServer(serverOptions, zapLogger, initServers, tracerProvider)
	if err2 != nil {
		return nil, err
	}
	cronOptions, err2 := cron.NewOptions(viper, zapLogger)
	if err2 != nil {
		return nil, err
	}
	defaultCronJobService := crontab.NewDefaultCronJobService(zapLogger, viper, redisDB, client, db, mongoDB)
	cronInitServers := crontab.CreateInitServersFn(defaultCronJobService)
	cronServer, err2 := cron.New(cronOptions, zapLogger, redisDB, cronInitServers)
	if err2 != nil {
		return nil, err
	}
	applicationApplication, err2 := app.NewApp(appOptions, zapLogger, httpServer, server2, cronServer)
	if err2 != nil {
		return nil, err
	}
	return applicationApplication, nil
}

// wire.go:

var providerSet = wire.NewSet(logger.ProviderSet, config.ProviderSet, mongo.ProviderSet, postgres.ProviderSet, es.ProviderSet, redis.ProviderSet, trace.ProviderSet, repository.ProviderSet, service.ProviderSet, grpcserver.ProviderSet, http.ProviderSet, grpc.ProviderSet, cron.ProviderSet, crontab.ProviderSet, api.ProviderSet, app.ProviderSet)