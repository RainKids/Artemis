package gocache

import (
	"github.com/bluele/gcache"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct {
	Size int `json:"size" yaml:"size" toml:"size" env:"GO_CACHE_SIZE"`
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("gocache", &o); err != nil {
		return nil, errors.Wrap(err, "unmarshal gocache option error")
	}
	logger.Info("load gocache options success", zap.Any("gocache options", o))
	return o, err
}

func New(o *Options) (gcache.Cache, error) {
	gc := gcache.New(20).LRU().Build()
	return gc, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
