package grpc

import (
	"admin/pkg/discovery"
	"admin/pkg/tools/network"
	"admin/pkg/transport/grpc/middleware/recovery"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerOptions struct {
	Port        int
	EtcdAddr    []string
	ServiceName string
	Timeout     time.Duration
}

type Server struct {
	o      *ServerOptions
	app    string
	host   string
	port   int
	logger *zap.Logger
	server *grpc.Server
}

func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*ServerOptions, error) {
	var (
		err error
		o   = new(ServerOptions)
	)
	if err = v.UnmarshalKey("grpc.server", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal grpc server option error")
	}

	logger.Info("load grpc options success", zap.Any("grpc options", o))

	return o, nil
}

type InitServers func(server *grpc.Server)

func NewServer(o *ServerOptions, logger *zap.Logger, init InitServers, tracer *trace.TracerProvider) (*Server, error) {
	var gs *grpc.Server
	prometheusExporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(prometheusExporter))
	grpc_prometheus.EnableHandlingTimeHistogram()
	rc := recovery.NewInterceptor(recovery.NewZapRecoveryHandler())
	unaryInts := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_recovery.UnaryServerInterceptor(),
		otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithMeterProvider(meterProvider)),
		rc.UnaryServerInterceptor(),
	}
	// if o.Timeout > 0 {
	// 	unaryInts = append(unaryInts, serverinterceptors.UnaryTimeoutInterceptor(o.Timeout))
	// }
	streamInts := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(logger),
		grpc_recovery.StreamServerInterceptor(),
		otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithMeterProvider(meterProvider)),
		rc.StreamServerInterceptor(),
	}
	gs = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInts...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInts...)),
	)
	init(gs)
	grpc_health_v1.RegisterHealthServer(gs, health.NewServer())
	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", o.ServiceName)),
		server: gs,
	}, nil
}

// Application 服务应用
func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = network.GetAvailablePort()
	}

	s.host = network.GetLocalIP4()

	if s.host == "" {
		return errors.New("get local ipv4 error")
	}
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.logger.Info("grpc server starting ...", zap.String("addr", addr))

	//将服务地址注册到etcd中
	etcdRegister := discovery.NewRegister(s.o.EtcdAddr, s.logger)
	//defer etcdRegister.Stop()
	node := discovery.Server{
		Name: s.o.ServiceName,
		Addr: addr,
	}
	go etcdRegister.Register(node, 10)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		_ = etcdRegister.Unregister()
		if i, ok := sig.(syscall.Signal); !ok {
			os.Exit(0)
		} else {
			os.Exit(int(i))
		}

	}()
	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatal("failed to listen: %v", zap.Error(err))
		}
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}

		if err := s.server.Serve(lis); err != nil {
			s.logger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	return nil
}

// Stop  停止GRPC服务
func (s *Server) Stop() error {
	s.logger.Info("grpc server stopping ...")
	s.server.GracefulStop()
	return nil
}
