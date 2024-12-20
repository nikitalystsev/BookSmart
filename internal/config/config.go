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
	Port     string
	Mongo    MongoConfig
	DBType   string
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

type MongoConfig struct {
	URI      string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func Init(configsDir string) (*Config, error) {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("db.type", &cfg.DBType); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB_NAME")
	cfg.Postgres.Username = os.Getenv("POSTGRES_DB_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_DB_PASSWORD")
	cfg.Postgres.SSLMode = os.Getenv("POSTGRES_SSL_MODE")

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Host = os.Getenv("MONGO_DB_HOST")
	cfg.Mongo.Port = os.Getenv("MONGO_DB_PORT")
	cfg.Mongo.DBName = os.Getenv("MONGO_DB_NAME")
	cfg.Mongo.Username = os.Getenv("MONGO_DB_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_DB_PASSWORD")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Username = os.Getenv("REDIS_USER")
	cfg.Redis.Password = os.Getenv("REDIS_USER_PASSWORD")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	cfg.Port = os.Getenv("APP_PORT")
}
