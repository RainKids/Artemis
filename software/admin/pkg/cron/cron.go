package cron

import (
	"admin/pkg/database/redis"
	"context"
	"github.com/gochore/dcron"
	"github.com/google/wire"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type RedisAtomic struct {
	client *goredis.Client
}

func (m *RedisAtomic) SetIfNotExists(ctx context.Context, key, value string) bool {
	ret := m.client.SetNX(ctx, key, value, time.Hour)
	return ret.Err() == nil && ret.Val()
}

type Options struct {
	ServerName string
	Projects   map[string]string
}

type ServerOptional struct {
	spec string
	f    func(c context.Context) error
}

type Server struct {
	app    string
	o      *Options
	logger *zap.Logger
	cron   *dcron.Cron
	jobs   map[string]ServerOptional
}

type InitServers map[string]func(c context.Context) error

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("cron", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal cron option error")
	}
	logger.Info("load cron options success", zap.Any("cron options", o))
	return o, err
}

func New(o *Options, logger *zap.Logger, rdb *redis.RedisDB, init InitServers) (*Server, error) {
	optionals := make(map[string]ServerOptional)
	for name, spec := range o.Projects {
		if jobFunc, ok := init[name]; ok {
			optionals[name] = ServerOptional{
				spec: spec,
				f:    jobFunc,
			}
		} else {
			logger.Error("定时任务不存在", zap.String("name", name))
			return nil, errors.New("定时任务不存在")
		}
	}
	atomic := &RedisAtomic{
		client: rdb.Client,
	}
	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", "cronServer")),
		cron:   dcron.NewCron(dcron.WithKey("Cron"), dcron.WithAtomic(atomic)),
		jobs:   optionals,
	}, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	go func() {
		if err := s.register(); err != nil {
			s.logger.Fatal("failed to register cron: %v", zap.Error(err))
		}
		s.cron.Start()
	}()
	return nil
}

func (s *Server) register() error {
	for name, obj := range s.jobs {
		job := dcron.NewJob(name, obj.spec, obj.f)
		err := s.cron.AddJobs(job)
		if err != nil {
			s.logger.Error("注册job失败", zap.Error(err))
			return err
		}
		s.logger.Info("注册cron任务成功", zap.String("name", name))
	}
	return nil
}

func (s *Server) deRegister() error {
	<-s.cron.Stop().Done()
	s.logger.Info("deregister cron services success", zap.String("name", s.o.ServerName))
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("cron server stopping ...")
	if err := s.deRegister(); err != nil {
		return errors.Wrap(err, "deregister cron server error")
	}
	return nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
