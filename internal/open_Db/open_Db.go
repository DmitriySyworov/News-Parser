package open_Db

import (
	"app/news-parser/internal/loggers"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDb struct {
	*gorm.DB
}
type RedisDb struct {
	*redis.Client
}

func OpenPostgres(DSN string, logger *loggers.Logger) *PostgresDb {
	db, errOpen := gorm.Open(postgres.Open(DSN))
	if errOpen != nil {
		logger.SystemLogger(slog.LevelError, "failed to connect PostgreSQL")
		os.Exit(1)
	}
	logger.SystemLogger(slog.LevelInfo, "PostgreSQl connection is successful")
	return &PostgresDb{
		DB: db,
	}
}
func OpenRedis(redisPassword, address string, logger *loggers.Logger) *RedisDb {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: redisPassword,
		DB:       0,
	})
	logger.SystemLogger(slog.LevelInfo, "Redis connection is successful")
	return &RedisDb{
		Client: rdb,
	}
}
