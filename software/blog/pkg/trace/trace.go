package trace

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	jaegerp "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type Options struct {
	Endpoint string `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"TRACE_PROVIDER_ENDPOINT"`
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("jaeger", &o); err != nil {
		return nil, errors.Wrap(err, "unmarshal jaeger configuration error")
	}
	logger.Info("load jaeger configuration success")

	return o, nil
}

func New(o *Options) (*sdktrace.TracerProvider, error) {
	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
	if err != nil {
		return nil, errors.Wrap(err, "connection jaeger endpoint error")
	}
	_, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger tracer error")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Default()),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		jaegerp.Jaeger{}),
	)
	return tp, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
