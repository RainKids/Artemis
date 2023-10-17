package postgres

import (
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"time"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Options database option
type Options struct {
	Dsn         string `toml:"dsn" json:"dsn" yaml:"dsn" env:"POSTGRES_DSN"`
	Debug       bool   `toml:"debug" json:"debug" yaml:"debug" env:"POSTGRES_DEBUG"`
	EnableTrace bool   `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"POSTGRES_ENABLE_TRACE"`
}

// NewOptions new database option
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("postgres", &o); err != nil {
		return nil, errors.Wrap(err, "unmarshal database option error")
	}
	logger.Info("load database options success", zap.Any("database options", o))
	return o, err
}

// Database 定义数据库struct
type DB struct {
	Postgres *gorm.DB
}

// DBClient
var DBClient DB

// New database
func New(o *Options) (*DB, error) {
	var d = new(DB)
	if o.Dsn == "" {
		return nil, errors.New("缺少postgresql配置")
	} else {
		Postgres, err := postgresql(o)
		if err != nil {
			return nil, err
		}
		d.Postgres = Postgres
	}
	DBClient.Postgres = d.Postgres
	return d, nil
}

// postgresql 定义timescaledb连接信息
func postgresql(o *Options) (*gorm.DB, error) {
	Psql, err := gorm.Open(postgres.Open(o.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return Psql, errors.Wrap(err, "gorm open postgresql connection error")
	}
	db, _ := Psql.DB()
	err = db.Ping()
	if err != nil {
		return Psql, errors.Wrap(err, "postgresql ping fail")
	}
	if o.Debug {
		Psql = Psql.Debug()
	}
	if o.EnableTrace {
		if err := Psql.Use(otelgorm.NewPlugin()); err != nil {
			return Psql, errors.Wrap(err, "gorm opentelemetry error")
		}
	}

	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	// 自动迁移模式
	// Psql.AutoMigrate()

	return Psql, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
