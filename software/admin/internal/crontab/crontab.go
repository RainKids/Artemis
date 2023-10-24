package crontab

import (
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"fmt"
	"github.com/gochore/dcron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DefaultCronJobService CronJob模块service默认对象
type DefaultCronJobService struct {
	logger *zap.Logger
	v      *viper.Viper
	rdb    *redis.RedisDB
	es     *es.Client
	db     *postgres.DB
	mongo  *mongo.MongoDB
}

// NewDefaultCronJobService 初始化
func NewDefaultCronJobService(
	logger *zap.Logger,
	v *viper.Viper,
	rdb *redis.RedisDB,
	es *es.Client,
	db *postgres.DB,
	mongo *mongo.MongoDB,
) *DefaultCronJobService {
	return &DefaultCronJobService{
		logger: logger.With(zap.String("type", "DefaultCronJobService")),
		v:      v,
		rdb:    rdb,
		es:     es,
		db:     db,
		mongo:  mongo,
	}
}

func (s *DefaultCronJobService) CronTest(c context.Context) error {
	//TODO
	if task, ok := dcron.TaskFromContext(c); ok {
		fmt.Print("run test job:", task.Job.Spec(), task.Key)
	}
	fmt.Println("print ---> 123")

	return nil
}
