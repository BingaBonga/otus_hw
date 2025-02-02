package configs

import (
	"fmt"

	//nolint:depguard
	"github.com/BurntSushi/toml"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Logger   LoggerConf
	DB       DBConfig
	HTTP     HTTPConfig
	Kafka    KafkaConfig
	Schedule ScheduleConfig
}

type LoggerConf struct {
	Level zapcore.Level
	Path  string
}

type DBConfig struct {
	InMemory bool
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

type HTTPConfig struct {
	Host string
	Port int
}

type KafkaConfig struct {
	Url          string
	Group        string
	ConsumeTopic string
	ProduceTopic string
	ServiceName  string
}

type ScheduleConfig struct {
	Cron string
}

func ReadConfig(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	if err != nil {
		return Config{}, fmt.Errorf("error decode configuration file: %w", err)
	}
	return
}
