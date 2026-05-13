package user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RepositoryUser struct {
	*open_Db.PostgresDb
	*open_Db.RedisDb
}

func NewRepositoryUser(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryUser {
	return &RepositoryUser{
		PostgresDb: postgres,
		RedisDb:    redis,
	}
}
func (r *RepositoryUser) IsUserExistByNameAndEmail(name, email string) error {
	resName := r.PostgresDb.Where("name = ?", name).First(&model.User{})
	if resName == nil {
		return custom_errors.ErrUserExist
	}
	resEmail := r.PostgresDb.Where("email = ?", email).First(&model.User{})
	if resEmail.Error == nil {
		return custom_errors.ErrUserExist
	}
	return nil
}
func (r *RepositoryUser) IsUserExistByUUID(uuid string) bool {
	res := r.PostgresDb.Where("uuid_user = ?", uuid).First(&model.User{})
	if res.Error != nil {
		return false
	}
	return true
}
func (r *RepositoryUser) CreateUser(user *model.User) error {
	if res := r.PostgresDb.Create(&user); res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUser) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := r.PostgresDb.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUser) GetUserByUUID(uuid string) (*model.User, error) {
	var user model.User
	res := r.PostgresDb.Where("uuid_user = ?", uuid).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUser) GetMyUser(uuid string) (*ResponseUser, error) {

	var getUser ResponseUser
	res := r.PostgresDb.
		Model(&model.User{}).
		Raw(`SELECT created_at, name, email, uuid_user FROM users
				  WHERE uuid_user = ?`, uuid).
		First(&getUser)
	if res.Error != nil {
		return nil, res.Error
	}
	return &getUser, nil
}

func (r *RepositoryUser) UpdateMyUserOneColumn(userUUID, columnName, value string) (*model.User, error) {
	var user model.User
	res := r.PostgresDb.Model(&model.User{}).
		Where("uuid_user = ?", userUUID).
		Update(columnName, value).
		First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUser) UpdateMyUserFull(user *model.User) error {
	res := r.PostgresDb.Where("uuid_user = ?", user.UUIDUser).Updates(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

const (
	fieldCode     = "code"
	fieldEmail    = "email"
	fieldName     = "name"
	fieldPassword = "password"
)

func (r *RepositoryUser) CreateSessionUpdate(temp *model.TemporaryData) error { //!!!
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	key := common.KeySession + actionUpdate + ":" + temp.IDSession
	_, errTrans := r.RedisDb.TxPipelined(rdbContext, func(pipeliner redis.Pipeliner) error {
		errHSet := r.RedisDb.HSet(rdbContext, key, fieldCode, temp.TempCode, fieldName, temp.Name, fieldEmail, temp.Email, fieldPassword, temp.Password).Err()
		if errHSet != nil {
			return errHSet
		}
		errExpired := r.RedisDb.Expire(rdbContext, key, time.Minute*5).Err()
		if errExpired != nil {
			return errExpired
		}
		return nil
	})
	if errTrans != nil {
		return errTrans
	}
	return nil
}
func (r *RepositoryUser) CreateSessionDeleteOrRemove(code uint, sessionId, action string) error {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	key := common.KeySession + action + ":" + sessionId
	errSet := r.RedisDb.Set(rdbContext, key, code, time.Minute*5).Err()
	if errSet != nil {
		return errSet
	}
	return nil
}
func (r *RepositoryUser) GetSession(sessionId, action string) (*model.TemporaryData, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	key := common.KeySession + action + ":" + sessionId
	if action == actionDelete || action == actionRemove {
		codeStr, errGet := r.RedisDb.Get(rdbContext, key).Result()
		if errGet != nil {
			return nil, errGet
		}
		code, _ := strconv.Atoi(codeStr)
		return &model.TemporaryData{
			TempCode: uint(code),
		}, nil
	}
	mapValue, errHGetAll := r.RedisDb.HGetAll(rdbContext, key).Result()
	if errHGetAll != nil {
		return nil, errHGetAll
	}
	code, _ := strconv.Atoi(mapValue[fieldCode])
	return &model.TemporaryData{
		Name:     mapValue[fieldName],
		Email:    mapValue[fieldEmail],
		Password: mapValue[fieldPassword],
		TempCode: uint(code),
	}, nil

}
func (r *RepositoryUser) RemoveMyUser(userUUID string) error {
	res := r.PostgresDb.
		Where("uuid_user = ?", userUUID).
		Delete(&model.User{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUser) DeleteMyUser(userUUID string) error {
	res := r.PostgresDb.
		Unscoped().
		Where("uuid_user = ?", userUUID).
		Delete(&model.User{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
