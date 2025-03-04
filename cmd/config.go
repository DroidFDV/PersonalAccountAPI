package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// `mapstructure: "filed_name"` нужен для того, чтобы сопоставить поля структуры и конфига
type Config struct {
	AppName         string `mapstructure:"app_name"`
	Port            string `mapstructure:"port"`
	Debug           bool   `mapstructure:"debug"`
	DatabaseURL     string `mapstructure:"database_url"`
	FileStoragePatg string `mapstructure:"file_storage_path"`
	MetricsPort     string `mapstructure:"metrics_port"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd")
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Warning: Could not read config file, falling back to environment variables")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "Unable to decode into struct")
	}

	return &cfg, nil
}
