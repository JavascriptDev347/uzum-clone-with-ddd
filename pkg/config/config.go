package config

import (
	"fmt"
	"net/url"
)

type Config struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT" default:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSLMODE" default:"disable"`

	RedisHost string `env:"REDIS_HOST" default:"localhost"`
	RedisPort string `env:"REDIS_PORT" default:"6379"`

	AppPort int `env:"APP_PORT" default:"8080"`
}

func (c Config) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   fmt.Sprintf("%s:%s", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
