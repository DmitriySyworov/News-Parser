package user

import "app/news-parser/configs"

type ServiceUser struct {
	Repo *RepositoryUser
	Dep  *HandlerUserDep
}
type ServiceUserDep struct {
	*configs.Configs
}

func NewServiceUser(repo *RepositoryUser, dep *HandlerUserDep) *ServiceUser {
	return &ServiceUser{
		Repo: repo,
		Dep:  dep,
	}
}
func (s *ServiceUser) RemoveMyUser(userUUID string) error {

}
func (s *ServiceUser) helperSendEmail()
