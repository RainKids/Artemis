package logger

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// ClientType 定义日志 client 结构体
type Client struct {
	*zap.Logger
}

// Client  logger连接类型
var Logger Client

// Options log option
type Options struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      string
	Stdout     bool
}

// NewOptions new log option
func NewOptions(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("log", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal log option error")
	}

	return o, err
}

// New log
func New(o *Options) (*zap.Logger, error) {
	var (
		err   error
		level = zap.NewAtomicLevel()
		l     *zap.Logger
	)
	if err = level.UnmarshalText([]byte(o.Level)); err != nil {
		return nil, err
	}

	fw := zapcore.AddSync(&lumberjack.Logger{
		Filename:   o.Filename,
		MaxSize:    o.MaxSize,
		MaxBackups: o.MaxBackups,
		MaxAge:     o.MaxAge,
	})
	cw := zapcore.Lock(os.Stdout)

	// file core 采用jsonEncoder
	cores := make([]zapcore.Core, 0, 2)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	je := zapcore.NewJSONEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(je, fw, level))

	// stdout core 采用consoleEncoder
	if o.Stdout {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		ce := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(ce, cw, level))
	}

	core := zapcore.NewTee(cores...)
	l = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(l)
	logger := otelzap.New(l)
	Logger.Logger = logger.Logger
	return logger.Logger, err
}

var ProviderSet = wire.NewSet(New, NewOptions)
