package auth

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RepositoryAuth struct {
	*open_Db.RedisDb
}

func NewRepositoryAuth(rdb *open_Db.RedisDb) *RepositoryAuth {
	return &RepositoryAuth{
		RedisDb: rdb,
	}
}

const (
	keySession = "session:"

	fieldName     = "name"
	fieldEmail    = "email"
	fieldPassword = "password"
	fieldCode     = "code"
)

func (r *RepositoryAuth) CreateTemporaryUser(tempUser *model.TemporaryUser) error {
	rdbCtx, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	key := fmt.Sprint(keySession, tempUser.IDSession)
	_, errTrans := r.Client.TxPipelined(rdbCtx, func(pipeliner redis.Pipeliner) error {
		errHSet := r.Client.HSet(rdbCtx, key, fieldName, tempUser.Name, fieldEmail, tempUser.Email, fieldPassword, tempUser.Password, fieldCode, tempUser.TempCode).Err()
		if errHSet != nil {
			return errHSet
		}
		errExpire := r.Client.Expire(rdbCtx, key, time.Minute*5).Err()
		if errExpire != nil {
			return errExpire
		}
		return nil
	})
	if errTrans != nil {
		return errTrans
	}
	return nil
}
func (r *RepositoryAuth) GetTemporaryUser(sessionID string) (*model.TemporaryUser, error) {
	rdbCtx, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	key := fmt.Sprint(keySession, sessionID)
	mapValue, errHGetAll := r.Client.HGetAll(rdbCtx, key).Result()
	if errHGetAll != nil {
		return nil, errHGetAll
	}
	code, errCode := strconv.Atoi(mapValue[fieldCode])
	if errCode != nil {
		return nil, errCode
	}
	return &model.TemporaryUser{
		Name:     mapValue[fieldName],
		Email:    mapValue[fieldEmail],
		Password: mapValue[fieldPassword],
		TempCode: uint(code),
	}, nil
}
