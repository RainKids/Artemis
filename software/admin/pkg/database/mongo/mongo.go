package mongo

import (
	"context"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/zap"
	"time"
)

type Options struct {
	UserName    string   `toml:"username" json:"username" yaml:"username"  env:"MONGO_USERNAME"`
	Password    string   `toml:"password" json:"password" yaml:"password"  env:"MONGO_PASSWORD"`
	Endpoints   []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	AuthDB      string   `toml:"auth_db" json:"auth_db" yaml:"auth_db"  env:"MONGO_AUTH_DB"`
	EnableTrace bool     `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"MONGO_ENABLE_TRACE"`
	Database    string   `toml:"database" json:"database" yaml:"database"  env:"MONGO_DATABASE"`
	DB          *mongo.Database
	Client      *mongo.Client
}

type MongoDB struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("mongo", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal database mongo option error")
	}

	logger.Info("load database options success", zap.Any("database mongo options", o))
	return o, err
}

func New(o *Options) (*MongoDB, error) {
	if len(o.Endpoints) == 0 {
		return nil, errors.New("缺少mongo配置")
	} else {
		mongodb, err := initDB(o)
		if err != nil {
			return nil, err
		}
		o.DB = mongodb.Database(o.Database)

	}
	return &MongoDB{
		o.DB,
		o.Client,
	}, nil
}

func initDB(m *Options) (*mongo.Client, error) {
	opts := options.Client()
	if m.UserName != "" && m.Password != "" {
		cred := options.Credential{
			AuthSource: m.GetAuthDB(),
		}

		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)
	if m.EnableTrace {
		opts.Monitor = otelmongo.NewMonitor(
			otelmongo.WithCommandAttributeDisabled(true),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "mongo driver open mongodb connection error")
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "mongo ping fail")
	}
	//db := client.Database(m.Database)
	m.Client = client
	return client, nil
}

func (m *Options) GetAuthDB() string {
	if m.AuthDB != "" {
		return m.AuthDB
	}

	return m.Database
}
func (db *Options) Close(ctx context.Context) error {
	if db.Client == nil {
		return nil
	}
	return db.Client.Disconnect(ctx)
}

// ProviderSet dependency injection
var ProviderSet = wire.NewSet(New, NewOptions)
