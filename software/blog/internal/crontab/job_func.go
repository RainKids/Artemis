package crontab

import (
	"blog/pkg/cron"
	"context"
	"github.com/google/wire"
)

func CreateInitServersFn(cronJob *DefaultCronJobService) cron.InitServers {
	return map[string]func(c context.Context) error{
		"sync.es.total.per.5s": func(c context.Context) error {
			return cronJob.RedisToES(c)
		},
		"sync.mongo.total.per.5s": func(c context.Context) error {
			return cronJob.RedisToMongo(c)
		},
		"sync.postgres.total.per.5s": func(c context.Context) error {
			return cronJob.RedisToPostgres(c)
		},
		"sync.DBToes.total.per.5s": func(c context.Context) error {
			return cronJob.PostgresToEs(c)
		},
	}
}

// ProviderSet CronJob Service wire 注入
var ProviderSet = wire.NewSet(NewDefaultCronJobService, CreateInitServersFn)
