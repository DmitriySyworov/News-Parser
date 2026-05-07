package user

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
)

type RepositoryUser struct {
	*open_Db.PostgresDb
}

func NewRepositoryUser(postgres *open_Db.PostgresDb) *RepositoryUser {
	return &RepositoryUser{
		PostgresDb: postgres,
	}
}
func (r *RepositoryUser) IsUserExist(name, email string) error {
	resName := r.PostgresDb.Where("name = ?", name).First(&model.User{})
	resEmail := r.PostgresDb.Where("email = ?", email).First(&model.User{})
	if resName.Error == nil || resEmail.Error == nil {
		return custom_errors.ErrUserExist
	}
	return nil
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
