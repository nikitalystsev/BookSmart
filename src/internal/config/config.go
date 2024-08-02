package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	Auth     AuthConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type AuthConfig struct {
	JWT          JWTConfig
	PasswordSalt string
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningKey      string
}
type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
}

func Init(configsDir string) (*Config, error) {
	if err := parseConfigFile(configsDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Host = os.Getenv("DB_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_DB_PORT")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB_NAME")
	cfg.Postgres.Username = os.Getenv("POSTGRES_DB_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_DB_PASSWORD")
	cfg.Postgres.SSLMode = "disable"

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Username = os.Getenv("REDIS_USER")
	cfg.Redis.Password = os.Getenv("REDIS_USER_PASSWORD")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

}
