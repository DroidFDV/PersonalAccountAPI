package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppName     string `mapstructure:"app_name"`
	Port        string `mapstructure:"port"`
	Debug       bool   `mapstructure:"debug"`
	DatabaseURL string `mapstructure:"database_url"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("cmd")
	viper.AddConfigPath("./")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Could not read config file, falling back to environment variables: %v\n", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &cfg, nil
}
