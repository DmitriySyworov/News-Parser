package configs

import (
	"app/news-parser/internal/loggers"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	*DbConnect
	*APIConf
	*Token
}
type DbConnect struct {
	DSN           string
	RedisPassword string
	RedisAddress  string
}
type APIConf struct {
	ApiPort     string
	ApiEmail    string
	ApiPassword string
	AddressHost string
	Address     string
	RodBin      string
}
type Token struct {
	Signature string
}

func NewConfigs(logger *loggers.Logger) *Configs {
	godotenv.Load(".env")
	dsn := os.Getenv("DSN")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddress := os.Getenv("REDIS_ADDRESS")
	apiPort := os.Getenv("API_EXTERNAL_PORT")
	apiEmail := os.Getenv("API_EMAIL")
	apiPassword := os.Getenv("API_PASSWORD")
	addressHost := os.Getenv("SMTP_ADDRESS_HOST")
	address := os.Getenv("SMTP_ADDRESS")
	signature := os.Getenv("JWT_SIGNATURE")
	rodBin := os.Getenv("ROD_BIN")
	counterEmpty := 0
	if dsn == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable DSN is empty")
	}
	if redisPassword == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable REDIS_PASSWORD is empty")
	}
	if redisAddress == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable REDIS_ADDRESS is empty")
	}
	if apiPort == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable API_EXTERNAL_PORT is empty")
	}
	if apiEmail == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable API_EMAIL is empty")
	}
	if apiPassword == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable API_PASSWORD is empty")
	}
	if addressHost == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable SMTP_ADDRESS_HOST is empty")
	}
	if address == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable SMTP_ADDRESS is empty")
	}
	if signature == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable JWT_SIGNATURE is empty")
	}
	if rodBin == "" {
		counterEmpty++
		logger.SystemLogger(slog.LevelError, "environment variable ROD_BIN is empty")
	}
	if counterEmpty != 0 {
		os.Exit(1)
	}
	logger.SystemLogger(slog.LevelInfo, "environment variables loaded successfully")
	return &Configs{
		DbConnect: &DbConnect{
			DSN:           dsn,
			RedisPassword: redisPassword,
			RedisAddress:  redisAddress,
		},
		APIConf: &APIConf{
			ApiPort:     apiPort,
			ApiEmail:    apiEmail,
			ApiPassword: apiPassword,
			AddressHost: addressHost,
			Address:     address,
			RodBin:      rodBin,
		},
		Token: &Token{
			Signature: signature,
		},
	}
}
