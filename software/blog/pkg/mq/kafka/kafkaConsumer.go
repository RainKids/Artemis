package kafka

import (
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type ConsumerOptions struct {
	Brokers        []string       `toml:"brokers" json:"brokers" yaml:"brokers"  env:"KAFKA_BROKERS"`
	ScramAlgorithm ScramAlgorithm `toml:"scram_algorithm" json:"scram_algorithm" yaml:"scram_algorithm"  env:"KAFKA_SCRAM_ALGORITHM"`
	Topics         []string       `toml:"topics" json:"topics" yaml:"topics"  env:"KAFKA_TOPICS"`
	GroupID        string         `toml:"groupId" json:"groupId" yaml:"groupId"  env:"KAFKA_GROUPID"`
	UserName       string         `toml:"username" json:"username" yaml:"username"  env:"KAFKA_USERNAME"`
	Password       string         `toml:"password" json:"password" yaml:"password"  env:"KAFKA_PASSWORD"`
	logger         *zap.Logger
}

func NewConsumerOptions(v *viper.Viper, logger *zap.Logger) (*ConsumerOptions, error) {
	var (
		err error
		o   = new(ConsumerOptions)
	)
	if err = v.UnmarshalKey("kafka.consumerOptions", &o); err != nil {
		return nil, errors.Wrap(err, "unmarshal kafka consumerOptions option error")
	}
	logger.Info("load kafka consumerOptions options success", zap.Any("kafka consumerOptions options", o))
	return o, err
}

func NewConsumer(o *ConsumerOptions) (*kafka.Reader, error) {
	mechanism, err := scram.Mechanism(scramAlgorithm(o.ScramAlgorithm), o.UserName, o.Password)
	if err != nil {
		return nil, err
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     o.Brokers,
		Dialer:      dialer,
		GroupID:     o.GroupID,
		GroupTopics: o.Topics,
		Logger:      kafka.LoggerFunc(o.logger.Sugar().Debugf),
		ErrorLogger: kafka.LoggerFunc(o.logger.Sugar().Errorf),
	}), nil
}
