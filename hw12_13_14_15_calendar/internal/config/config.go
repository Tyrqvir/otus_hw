package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	DB     DBConf
	HTTP   HTTPConf
	GRPS   HTTPConf
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	DSN      string
	Provider string
}

type HTTPConf struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
}

type GRPSConf struct {
	Host string
	Port string
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%v: %w", ErrFailedReadConfigFile, err)
	}

	return &Config{
		LoggerConf{
			Level: viper.GetString("logger.level"),
		},
		DBConf{
			DSN:      viper.GetString("db.DSN"),
			Provider: viper.GetString("db.provider"),
		},
		HTTPConf{
			Host:              viper.GetString("http.host"),
			Port:              viper.GetString("http.port"),
			ReadTimeout:       viper.GetDuration("http.read_timeout"),
			WriteTimeout:      viper.GetDuration("http.write_timeout"),
			ReadHeaderTimeout: viper.GetDuration("http.read_header_timeout"),
		},
		HTTPConf{
			Host: viper.GetString("grps.host"),
			Port: viper.GetString("grps.port"),
		},
	}, nil
}