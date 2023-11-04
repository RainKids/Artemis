package casbin

import (
	"admin/pkg/database/postgres"
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

type Options struct {
	URL         string // host:port
	MaxIdle     int    // 最大空闲连接数
	MaxActive   int    // 最大连接数
	IdleTimeout int    // 空闲连接超时时间
	Timeout     int    // 操作超时时间
	Network     string // tcp or udp
	Password    string
}

var (
	log      *zap.Logger
	enforcer *casbin.SyncedEnforcer
	once     sync.Once
)

// NewOptions for redis
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("redis", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}

	logger.Info("load redis options success", zap.Any("redis options", o))
	log = logger
	return o, err
}

func New(o *Options, logger *zap.Logger, postgresDB *postgres.DB) (*casbin.SyncedEnforcer, error) {
	once.Do(func() {
		db := postgresDB.Postgres
		Apter, err := gormadapter.NewAdapterByDBUseTableName(db.WithContext(context.Background()), "sys", "casbin_rule")
		if err != nil {
			logger.Error("read db casbin rule failed!", zap.Error(err))
			panic(err)
		}
		m, err := model.NewModelFromString(text)
		if err != nil {
			logger.Error("字符串加载模型失败!", zap.Error(err))
			panic(err)
		}
		enforcer, err = casbin.NewSyncedEnforcer(m, Apter)
		if err != nil {
			logger.Error("创建数据库同步执行器失败!", zap.Error(err))
			panic(err)
		}
		if err = enforcer.LoadPolicy(); err != nil {
			logger.Error("casbin rbac_model or policy init error, message: %v \r\n", zap.Error(err))
			panic(err)
		}
		w, err := rediswatcher.NewWatcher(o.URL, rediswatcher.WatcherOptions{Options: redis.Options{
			Addr:         o.URL,
			Password:     o.Password, // no password set
			DB:           0,          // use default DB
			MaxIdleConns: o.MaxIdle,
			PoolSize:     o.MaxActive,
			DialTimeout:  time.Duration(o.IdleTimeout) * time.Second,
			Network:      "tcp",
		},
			Channel:    "/casbin",
			IgnoreSelf: false,
		})
		if err != nil {
			logger.Error("use the Redis host as parameter error: %v \r\n", zap.Error(err))
			panic(err)
		}
		err = w.SetUpdateCallback(updateCallback)
		if err != nil {
			logger.Error("set callback error: %v \r\n", zap.Error(err))
			panic(err)
		}

		err = enforcer.SetWatcher(w)
		if err != nil {
			logger.Error("set the watcher for the enforcer error: %v \r\n", zap.Error(err))
			panic(err)
		}
		enforcer.EnableLog(true)
	})

	return enforcer, nil
}

func updateCallback(msg string) {
	log.Sugar().Infof("casbin updateCallback msg: %v", msg)
	err := enforcer.LoadPolicy()
	if err != nil {
		log.Sugar().Errorf("casbin LoadPolicy err: %v", err)
	}
}

var ProviderSet = wire.NewSet(New, NewOptions)
