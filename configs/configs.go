package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	*DbConnect
}
type DbConnect struct {
	DSN           string
	RedisPassword string
	RedisAddress string
}

func NewConfigs() *Configs {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		panic(errEnv)
	}
	return &Configs{
		DbConnect: &DbConnect{
			DSN:           os.Getenv("DSN"),
			RedisPassword: os.Getenv("REDIS"),
			RedisAddress: os.Getenv("REDIS_ADDRESS"),
		},
	}
}
