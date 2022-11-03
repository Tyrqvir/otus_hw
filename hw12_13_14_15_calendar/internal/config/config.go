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
	GRPC      GRPCConf
	Publisher PublisherConf
	Consumer  ConsumerConf
	Schedule  ScheduleConf
	RMQ       RMQConf
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

type GRPCConf struct {
	Port string
}

type PublisherConf struct {
	Tag            string
	ConnectionName string
}

type RMQConf struct {
	Dsn          string
	QueueName    string
	ExchangeName string
	ExchangeType string
	BindingKey   string
}

type ConsumerConf struct {
	Tag            string
	ConnectionName string
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
		GRPCConf{
			Port: viper.GetString("grpc.port"),
		},
		PublisherConf{
			ConnectionName: viper.GetString("publisher.connectionName"),
			Tag:            viper.GetString("publisher.tag"),
		},
		ConsumerConf{
			ConnectionName: viper.GetString("consumer.connectionName"),
			Tag:            viper.GetString("consumer.tag"),
		},
		ScheduleConf{
			Interval: viper.GetDuration("schedule.interval"),
		},
		RMQConf{
			Dsn:          viper.GetString("rmq.dsn"),
			QueueName:    viper.GetString("rmq.queue"),
			ExchangeName: viper.GetString("rmq.exchangeName"),
			ExchangeType: viper.GetString("rmq.exchangeType"),
			BindingKey:   viper.GetString("rmq.bindingKey"),
		},
	}, nil
}
