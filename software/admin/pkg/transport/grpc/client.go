package grpc

import (
	"admin/pkg/discovery"
	"admin/pkg/transport/grpc/middleware/exception"
	"context"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	grpcInsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	golog "log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ClientOptions struct {
	Timeout         time.Duration
	Tag             string
	GrpcDialOptions []grpc.DialOption
	logger          *zap.Logger
	EtcdAddr        []string
	ServerName      string
}
type Client struct {
	o      *ClientOptions
	Logger *zap.Logger
}

type ClientOptional func(o *ClientOptions)

func NewClientOptions(v *viper.Viper, logger *zap.Logger, tracer *trace.TracerProvider) (*ClientOptions, error) {
	var (
		err error
		o   = new(ClientOptions)
	)
	if err = v.UnmarshalKey("grpc.client", o); err != nil {
		return nil, err
	}

	logger.Info("load grpc.client options success", zap.Any("grpc.client options", o))
	prometheusExporter, err := prometheus.New()
	if err != nil {
		golog.Fatal(err)
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(prometheusExporter))
	grpc_prometheus.EnableClientHandlingTimeHistogram()
	streamInts := []grpc.StreamClientInterceptor{
		grpc_prometheus.StreamClientInterceptor,
		grpc_zap.StreamClientInterceptor(logger),
		otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithMeterProvider(meterProvider)),
	}
	unaryInts := []grpc.UnaryClientInterceptor{
		grpc_prometheus.UnaryClientInterceptor,
		grpc_zap.UnaryClientInterceptor(logger),
		otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithMeterProvider(meterProvider)),
		//clientinterceptors.TimeoutInterceptor(o.Timeout),
	}
	o.GrpcDialOptions = append(o.GrpcDialOptions,
		//grpc.WithTransportCredentials(grpcInsecure.NewCredentials()),
		// 将异常转化为 API Exception
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(exception.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(streamInts...)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInts...)),
	)

	o.logger = logger.With(zap.String("type", o.ServerName))
	return o, nil
}

func NewClient(o *ClientOptions) (*Client, error) {
	return &Client{
		o:      o,
		Logger: o.logger,
	}, nil
}

// WithTimeout grpc client time out
func WithTimeout(d time.Duration) ClientOptional {
	return func(o *ClientOptions) {
		o.Timeout = d
	}
}

// WithTag grpc client tag
func WithTag(tag string) ClientOptional {
	return func(o *ClientOptions) {
		o.Tag = tag
	}
}

func WithEndpoint(endpoints []string) ClientOptional {
	return func(o *ClientOptions) {
		o.EtcdAddr = endpoints
	}
}

func WithLogger(logger *zap.Logger) ClientOptional {
	return func(o *ClientOptions) {
		o.logger = logger
	}
}

func WithGrpcDialOptions(options ...grpc.DialOption) ClientOptional {
	return func(o *ClientOptions) {
		o.GrpcDialOptions = append(o.GrpcDialOptions, options...)
	}
}

func (c *Client) DialInsecure(service string, options ...ClientOptional) (*grpc.ClientConn, error) {
	return c.dial(service, true, options...)
}

// Dial grpc client dail
func (c *Client) Dial(service string, options ...ClientOptional) (*grpc.ClientConn, error) {
	return c.dial(service, false, options...)
}

func (c *Client) dial(service string, insecure bool, options ...ClientOptional) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	o := &ClientOptions{
		Timeout:         c.o.Timeout,
		Tag:             c.o.Tag,
		GrpcDialOptions: c.o.GrpcDialOptions,
		EtcdAddr:        c.o.EtcdAddr,
		ServerName:      c.o.ServerName,
		logger:          c.o.logger,
	}
	//options = append(options, WithGrpcDialOptions(grpc.WithInsecure()))
	options = append(options, WithGrpcDialOptions(grpc.WithTransportCredentials(grpcInsecure.NewCredentials())))

	options = append(options, WithGrpcDialOptions(grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))))
	for _, option := range options {
		option(o)
	}
	go func() {

	}()
	etcdRegister := discovery.NewResolver(o.EtcdAddr, o.logger)
	resolver.Register(etcdRegister)
	addr := fmt.Sprintf("%s:///%s", etcdRegister.Scheme(), service)
	conn, err := grpc.DialContext(ctx, addr, o.GrpcDialOptions...)
	// conn, err := grpc.DialContext(ctx, service, o.GrpcDialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial error")
	}
	return conn, nil
}
