package config

import (
	"blog/pkg/tools/uuid"
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func New(path string) (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)
	v.AddConfigPath(".")
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		return nil, errors.Wrap(err, "read config file error")
	}
	if v.GetString("server.uuid") == "" {
		uuidstr, err := uuid.GetHostUuid()
		if err != nil || uuidstr == "" {
			fmt.Println("new uuid")
			uuidstr = uuid.NewUUID()
		}
		v.Set("server.uuid", uuidstr)
	}
	_ = v.WriteConfig()
	return v, err
}

var ProviderSet = wire.NewSet(New)
