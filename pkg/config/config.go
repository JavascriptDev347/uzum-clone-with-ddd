package config

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/bakhod1r/oneenv"
)

var baseOptions = []oneenv.Option{
	oneenv.WithRequired(),
	oneenv.WithEnvFiles(),
	oneenv.WithSecretReveal(4),
	oneenv.WithOnce(),
	oneenv.WithTable(),
	oneenv.WithWriteExample(),
	oneenv.WithMutex(&sync.RWMutex{}),
	oneenv.WithWatch(),
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE" default:"disable"`
}

type RedisConfig struct {
	Host string `env:"REDIS_HOST" default:"localhost"`
	Port string `env:"REDIS_PORT" default:"6379"`
}

type Config struct {
	DB      DBConfig
	Redis   RedisConfig
	AppPort int `env:"APP_PORT" default:"8080"`
}

func (c Config) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DB.User, c.DB.Password),
		Host:   fmt.Sprintf("%s:%s", c.DB.Host, c.DB.Port),
		Path:   c.DB.Name,
	}
	q := u.Query()
	q.Set("sslmode", c.DB.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func Load() (*Config, error) {
	cfg, err := oneenv.Parse[Config](baseOptions...)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
