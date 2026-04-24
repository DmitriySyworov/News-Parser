package open_Db

import (
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

func OpenPostgres(DSN string) *PostgresDb {
	db, errOpen := gorm.Open(postgres.Open(DSN))
	if errOpen != nil {
		panic(errOpen)
	}
	return &PostgresDb{
		DB: db,
	}
}
func OpenRedis(redisPassword, address string)*RedisDb{
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
		Password: redisPassword,
		DB: 0,
	})
	return &RedisDb{
		Client : rdb,
	}
}
