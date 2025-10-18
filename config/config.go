package config

import (
	"PersonalAccountAPI/internal/models"
	"bytes"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//go:embed config.yaml
var config embed.FS

type Config struct {
	App struct {
		Name    string
		Port    string
		Workers struct {
			Number   int
			QueueLen int
		}
		Cache struct {
			Ttl      time.Duration
			Interval time.Duration
		}
		Debug bool
	}

	DB struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}

	Metrics struct {
		Port string
	}

	FileStorage struct {
		Endpoint   string
		User       string
		Password   string
		UseSSL     bool
		BucketName string
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

func (c *Config) GetCacheTTL() time.Duration {
	return c.App.Cache.Ttl
}

func (c *Config) GetCacheInterval() time.Duration {
	return c.App.Cache.Interval
}

func (c *Config) GetWorkersNum() int {
	return c.App.Workers.Number
}

func (c *Config) GetWorkersQueueLen() int {
	return c.App.Workers.QueueLen
}

func (c *Config) GetS3Config() models.S3Config {
	return models.S3Config{
		Endpoint:        c.FileStorage.Endpoint,
		AccessKeyID:     c.FileStorage.User,
		SecretAccessKey: c.FileStorage.Password,
		UseSSL:          c.FileStorage.UseSSL,
		BucketName:      c.FileStorage.BucketName,
	}
}

func (c *Config) GetS3BucketName() string {
	return c.FileStorage.BucketName
}
