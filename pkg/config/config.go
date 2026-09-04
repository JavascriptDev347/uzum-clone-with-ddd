package config

import (
	"fmt"
	"net/url"
	"sync"
	"time"

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

type Config struct {
	DB      DBConfig    `envPrefix:"DB_"`
	Redis   RedisConfig `envPrefix:"REDIS_"`
	JWT     JWTConfig   `envPrefix:"JWT_"`
	AppPort int         `env:"APP_PORT" default:"8080"`
	AWS     AWSConfig   `envPrefix:"AWS_"`
}

type DBConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT" default	:"5432"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
	SSLMode  string `env:"SSLMODE" default:"disable"`
}

type RedisConfig struct {
	Host string `env:"HOST" default:"localhost"`
	Port string `env:"PORT" default:"6379"`
}

type JWTConfig struct {
	Secret     string        `env:"SECRET,required"`
	AccessTTL  time.Duration `env:"ACCESS_TTL" default:"15m"`
	RefreshTTL time.Duration `env:"REFRESH_TTL" default:"168h"` // 7 kun
}

type AWSConfig struct {
	AccessKeyID     string `env:"ACCESS_KEY_ID,required"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,required"`
	Region          string `env:"REGION,required"`
	S3BucketName    string `env:"S3_BUCKET_NAME,required"`
	S3Folder        string `env:"S3_FOLDER" default:"uzum-clone"`
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
