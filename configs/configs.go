package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	*DbConnect
	*EmailConf
	*Token
}
type DbConnect struct {
	DSN           string
	RedisPassword string
	RedisAddress  string
}
type EmailConf struct {
	ApiEmail    string
	ApiPassword string
	AddressHost string
	Address     string
}
type Token struct {
	Signature string
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
			RedisAddress:  os.Getenv("REDIS_ADDRESS"),
		},
		EmailConf: &EmailConf{
			ApiEmail:    os.Getenv("API_EMAIL"),
			ApiPassword: os.Getenv("API_PASSWORD"),
			AddressHost: os.Getenv("ADDRESS_HOST"),
			Address:     os.Getenv("ADDRESS"),
		},
		Token: &Token{
			Signature: os.Getenv("JWT_SIGNATURE"),
		},
	}
}
