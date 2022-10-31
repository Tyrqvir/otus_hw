package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger    LoggerConf
	DB        DBConf
	HTTP      HTTPConf
	GRPS      GRPSConf
	Publisher PublisherConf
	Consumer  ConsumerConf
	Schedule  ScheduleConf
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
	Port string
}

type PublisherConf struct {
	Dsn          string
	QueueName    string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

type ConsumerConf struct {
	Dsn          string
	QueueName    string
	ExchangeName string
	ExchangeType string
	BindingKey   string
}

type ScheduleConf struct {
	Interval time.Duration
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%v: %w", ErrFailedReadConfigFile, err)
	}

	return &Config{
		LoggerConf{
			Level: viper.GetString("logger.level"),
		},
		DBConf{
			DSN:      viper.GetString("storage.DSN"),
			Provider: viper.GetString("storage.provider"),
		},
		HTTPConf{
			Host:              viper.GetString("http.host"),
			Port:              viper.GetString("http.port"),
			ReadTimeout:       viper.GetDuration("http.read_timeout"),
			WriteTimeout:      viper.GetDuration("http.write_timeout"),
			ReadHeaderTimeout: viper.GetDuration("http.read_header_timeout"),
		},
		GRPSConf{
			Port: viper.GetString("grps.port"),
		},
		PublisherConf{
			Dsn:          viper.GetString("publisher.dsn"),
			QueueName:    viper.GetString("publisher.queue"),
			ExchangeName: viper.GetString("publisher.exchangeName"),
			ExchangeType: viper.GetString("publisher.exchangeType"),
			RoutingKey:   viper.GetString("publisher.routingKey"),
		},
		ConsumerConf{
			Dsn:          viper.GetString("consumer.dsn"),
			QueueName:    viper.GetString("consumer.queue"),
			ExchangeName: viper.GetString("consumer.exchangeName"),
			ExchangeType: viper.GetString("consumer.exchangeType"),
			BindingKey:   viper.GetString("consumer.bindingKey"),
		},
		ScheduleConf{
			Interval: viper.GetDuration("schedule.interval"),
		},
	}, nil
}
