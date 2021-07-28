package config

import (
	"fmt"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppAddr string `env:"APP_ADDR" env-default:"0.0.0.0"`
	AppPort string `env:"APP_PORT" env-default:"8080"`

	PostgresHost string `env:"POSTGRES_HOST" env-default:"127.0.0.1"`
	PostgresPort string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPass string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	PostgresDB   string `env:"POSTGRES_DB" env-default:"todo"`

	Mode string `env:"APP_MODE" env-default:"debug"`

	Email         string `env:"EMAIL"`
	EmailPassword string `env:"EMAIL_PASSWORD"`

	Domain string `env:"DOMAIN"`

	RedisHost      string `env:"REDIS_HOST" env-default:"127.0.0.1"`
	RedisPort      string `env:"REDIS_PORT" env-default:"6379"`
	MaxLoggedInStr string `env:"MAX_LOGGED_IN" env-default:"6"`
	MaxLoggedIn    int

	AccessKey  string `env:"JWT_ACCESS_KEY" env-default:"access_key"`
	RefreshKey string `env:"JWT_REFRESH_KEY" env-default:"refresh_key"`
}

func New(configName string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(configName, cfg)
	if err != nil {
		return nil, err
	}

	cfg.MaxLoggedIn, err = strconv.Atoi(cfg.MaxLoggedInStr)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

// Return addr:port for server
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.AppAddr, c.AppPort)
}
