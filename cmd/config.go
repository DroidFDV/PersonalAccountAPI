package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//go:embed config.yaml
var config embed.FS

type Config struct {
	App struct {
		Name  string
		Port  string
		Debug bool
	}

	DB struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}

	Storage struct {
		Path string
	}

	Metrics struct {
		Port string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Читаем файл из embedded FS
	data, err := config.ReadFile("config.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read embedded config file")
	}

	// Загружаем конфигурацию из данных
	if err := viper.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, errors.Wrap(err, "unable to decode embedded config file")
	}

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, errors.Wrap(err, "failed to load .env")
		}
	}

	// Поддержка переменных окружения (префикс, например, APP_PORT)
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "unable to decode into struct")
	}

	return &cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DB.User,
		c.DB.Pass,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}
