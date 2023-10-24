package crontab

import (
	"admin/pkg/cron"
	"context"
	"github.com/google/wire"
)

func CreateInitServersFn(cronJob *DefaultCronJobService) cron.InitServers {
	return map[string]func(c context.Context) error{
		"sync-crontest-total-per-5s": func(c context.Context) error {
			return cronJob.CronTest(c)
		},
	}
}

// ProviderSet CronJob Service wire 注入
var ProviderSet = wire.NewSet(NewDefaultCronJobService, CreateInitServersFn)
