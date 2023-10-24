package application

import (
	"admin/pkg/cron"
	"admin/pkg/transport/grpc"
	"admin/pkg/transport/http"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// Application app server
type Application struct {
	name       string
	logger     *zap.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
	cronServer *cron.Server
}

// Option app option
type Option func(*Application) error

// HTTPServerOption app http server option
func HTTPServerOption(svr *http.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.httpServer = svr
		return nil
	}
}

// GrpcServerOptions app grpc server option
func GrpcServerOptions(svr *grpc.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.grpcServer = svr
		return nil
	}
}

// CronServerOptions
func CronServerOptions(svr *cron.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.cronServer = svr
		return nil
	}
}

// new app
func New(name string, logger *zap.Logger, options ...Option) (*Application, error) {
	app := &Application{
		name:   name,
		logger: logger.With(zap.String("type", "Application")),
	}

	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

// Start start app server
func (a *Application) Start() error {
	if a.grpcServer != nil {
		if err := a.grpcServer.Start(); err != nil {
			return errors.Wrap(err, "grpc server start error")
		}
	}

	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return errors.Wrap(err, "http server start error")
		}
	}

	if a.cronServer != nil {
		if err := a.cronServer.Start(); err != nil {
			return errors.Wrap(err, "cron server start error")
		}
	}

	return nil
}

// AwaitSignal await signal for exit app server
func (a *Application) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	s := <-c
	a.logger.Info("receive a signal", zap.String("signal", s.String()))
	if a.httpServer != nil {
		if err := a.httpServer.Stop(); err != nil {
			a.logger.Error("stop http server error", zap.Error(err))
		}
	}
	if a.grpcServer != nil {
		if err := a.grpcServer.Stop(); err != nil {
			a.logger.Error("stop grpc server error", zap.Error(err))
		}
	}
	if a.cronServer != nil {
		if err := a.cronServer.Stop(); err != nil {
			a.logger.Error("stop cron server error", zap.Error(err))
		}
	}

	os.Exit(0)
}

// ProviderSet wire 注入
var ProviderSet = wire.NewSet(New)
