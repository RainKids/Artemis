package app

import (
	"admin/pkg/application"
	"admin/pkg/cron"
	"admin/pkg/transport/grpc"
	"admin/pkg/transport/http"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options 定义User类配置选项
type Options struct {
	Name string
}

// NewOptions 初始化Options
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("app", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal app option error")
	}
	logger.Info("load application options success")
	return o, err
}

// NewApp 初始化app
func NewApp(o *Options, logger *zap.Logger, hs *http.Server, gs *grpc.Server, con *cron.Server) (*application.Application, error) {
	a, err := application.New(o.Name, logger, application.GrpcServerOptions(gs), application.HTTPServerOption(hs),
		application.CronServerOptions(con))
	if err != nil {
		return nil, errors.Wrap(err, "new app error")
	}
	return a, nil
}

// ProviderSet user模块wire NewSet
var ProviderSet = wire.NewSet(NewApp, NewOptions)
