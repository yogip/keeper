package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type IAMConfig struct {
	TokenTTL  time.Duration `env:"TOKEN_TTL" envDefault:"1h"`
	SecretKey string        `env:"SECRET_KEY,unset" envDefault:"local-default-secret"`
}

type S3Config struct {
	BucketName      string `env:"BUCKET_NAME" envDefault:"bucket"`
	Endpoint        string `env:"ENDPOINT" envDefault:"localhost:9000"`
	AccessKeyID     string `env:"AK_ID" envDefault:"admin"`
	SecretAccessKey string `env:"SECRET_AK" envDefault:"password"`
}

type ServerConfig struct {
	Address           string `env:"RUN_ADDRESS" envDefault:"0.0.0.0:8080"`
	DatabaseDSN       string `env:"DATABASE_URI,unset" envDefault:"host=postgres-gophermart port=25432 user=username password=password dbname=gophermart sslmode=disable"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"debug"`
	MasterKey         string `env:"MASTER_KEY,unset" envDefault:"stub-master-secret"`
	EncryptionKeyPath string `env:"ENCRYPTION_KEY_PATH" envDefault:"encryption_keys/"`
}

type Config struct {
	Server ServerConfig
	IAM    IAMConfig `envPrefix:"IAM_"`
	S3     S3Config  `envPrefix:"S3_"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	var serverAddress, databaseDSN, logLevel string
	flag.StringVar(&serverAddress, "a", "", "Address and port to run server")
	flag.StringVar(&databaseDSN, "d", "", "Database URI")
	flag.StringVar(&logLevel, "l", "", "Log levle: debug, info, warn, error, panic, fatal")
	flag.Parse()

	if serverAddress != "" {
		cfg.Server.Address = serverAddress
	}
	if databaseDSN != "" {
		cfg.Server.DatabaseDSN = databaseDSN
	}
	if logLevel != "" {
		cfg.Server.LogLevel = logLevel
	}

	return &cfg, nil
}
