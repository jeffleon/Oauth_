package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

var Config config

func init() {
	if err := env.Parse(&Config); err != nil {
		logrus.Fatalf("Error initializing: %s", err.Error())
	}
}

const (
	EmailWelcomeSubject = "Welcome Email"
	EmailWelcomeMsg     = "Welcome to the jungle"
)

type config struct {
	Environment       string        `env:"APP_ENV"`
	Port              string        `env:"PORT" envDefault:"8080"`
	DbUser            string        `env:"DB_USER" envDefault:""`
	DbPassword        string        `env:"DB_PASSWORD" envDefault:""`
	DbHost            string        `env:"DB_HOST" envDefault:""`
	DbName            string        `env:"DB_NAME" envDefault:""`
	DbOptions         string        `env:"DB_OPTIONS" envDefault:""`
	DbTimeZone        string        `env:"DB_TIME_ZONE" envDefault:"America/Bogota"`
	Timeout           time.Duration `env:"DB_TIMEOUT" envDefault:"10s"`
	TokenExp          time.Duration `env:"TOKEN_EXP" envDefault:"2h"`
	EmailFrom         string        `env:"EMAIL_FROM" envDefault:"example@gmail.com"`
	EmailFromName     string        `env:"EMAIL_FROM_NAME" envDefault:"example"`
	TokenSecret       string        `env:"TOKEN_SECRET"`
	TokenExpFP        time.Duration `env:"TOKEN_EXP_FP" envDefault:"48h"`
	TokenRefreshExp   time.Duration `env:"REFRESH_TOKEN_EXP" envDefault:"120h"`
	TokenSecretFP     string        `env:"TOKEN_SECRET_FP"`
	TokenBlackListSet string        `env:"REDIS_TOKEN_BLACK_LIST_SET"`
	RPCHost           string        `env:"RPC_HOST" envDefault:"localhost"`
	RPCPort           string        `env:"RPC_PORT" envDefault:"5001"`
}

func (b *config) DBConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=%s",
		b.DbHost,
		b.DbUser,
		b.DbPassword,
		b.DbName,
		b.DbTimeZone,
	)
}
