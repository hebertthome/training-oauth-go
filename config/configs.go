package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"bitbucket.org/hebertthome/traning-oauth-go/logger"

	path "path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configuration *Configuration
)

type Configuration struct {
	Bind   string               `json:"bind" mapstructure:"bind"`
	Logger logger.Configuration `json:"logger" mapstructure:"logger"`
	Redis  RedisConfig          `json:"redis" mapstructure:"redis"`
	Auth   AuthConfig           `json:"auth" mapstructure:"auth"`
}

type RedisConfig struct {
	GoRedis GoRedisConfig `mapstructure:"go-redis"`
	Redigo  RedigoConfig  `mapstructure:"redigo"`
}

type GoRedisConfig struct {
	Addresses    map[string]string `mapstructure:"addresses"`
	TtlInSeconds int               `mapstructure:"ttlInSeconds"`
	Context      context.Context
}

func (rc GoRedisConfig) IsStructureEmpty() bool {
	return reflect.DeepEqual(rc, GoRedisConfig{})
}

type RedigoConfig struct {
	Address      string `mapstructure:"address"`
	TtlInSeconds int    `mapstructure:"ttlInSeconds"`
	ExpireCount  int    `mapstructure:"expireCount"`
}

func (rc RedigoConfig) IsStructureEmpty() bool {
	return reflect.DeepEqual(rc, RedigoConfig{})
}

type AuthConfig struct {
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
}

func Setup(configurationPath string) error {
	// Configure Viper
	ext := path.Ext(configurationPath)
	name := strings.TrimSuffix(path.Base(configurationPath), ext)
	path := path.Dir(configurationPath)
	var v = viper.New()
	v.SetConfigName(name)
	v.SetConfigType(ext[1:])
	v.AddConfigPath(path)
	v.WatchConfig()
	// Viper event to reload configuration if file has change
	v.OnConfigChange(func(e fsnotify.Event) {
		mountConfiguration(configurationPath, configuration, v)
	})
	// Mount configuration
	config, err := mountConfiguration(configurationPath, configuration, v)
	if err != nil {
		return err
	}
	configuration = config
	if !config.Redis.GoRedis.IsStructureEmpty() {
		config.Redis.GoRedis.Context = context.Background()
	}
	return nil
}

func Get() Configuration {
	return *configuration
}

func mountConfiguration(configurationPath string, configuration *Configuration, v *viper.Viper) (*Configuration, error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.SetupErr[Path=%s Err=%s]", configurationPath, err.Error())
	}
	if err := v.Unmarshal(&configuration); err != nil {
		return nil, fmt.Errorf("config.DecodeErr[Path=%s Err=%s]", configurationPath, err.Error())
	}
	if err := logger.Setup(configuration.Logger); err != nil {
		return nil, err
	}
	fmt.Printf("File loaded at time: %v \n", time.Now())
	fmt.Printf("%v \n", configurationPath)
	return configuration, nil
}
