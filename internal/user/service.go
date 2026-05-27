package user

import (
	"app/news-parser/configs"
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/JWT"
	"app/news-parser/pkg/generate_random"
	"app/news-parser/pkg/handler_request"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser struct {
	Repo *RepositoryUser
	Dep  *ServiceUserDep
}
type ServiceUserDep struct {
	*configs.Configs
}

func NewServiceUser(repo *RepositoryUser, dep *ServiceUserDep) *ServiceUser {
	return &ServiceUser{
		Repo: repo,
		Dep:  dep,
	}
}

const (
	actionRemove = "soft-delete"
	actionDelete = "hard-delete"
	actionUpdate = "update"
)

func (s *ServiceUser) RemoveMyUser(userUUID, password, action string) (*common.ResponseAuth, *response.Error) {
	user, errGetUser := s.Repo.GetUserByUUID(userUUID)
	if errGetUser != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		}
	}
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errPassword != nil {
		return nil, &response.Error{
			Message: ErrIncorrectPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	token, errSecurity := s.helperSecurity(user.Email, action, nil)
	if errSecurity != nil {
		return nil, &response.Error{
			Message: errSecurity.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &common.ResponseAuth{
		Message: common.MessageEmail + user.Email,
		JWTTemp: token,
	}, nil
}
func (s *ServiceUser) UpdateMyUser(body *RequestUpdateUser, userUUID string) (*ResponseUser, *common.ResponseAuth, *response.Error) {
	user, errGetUser := s.Repo.GetUserByUUID(userUUID)
	if errGetUser != nil {
		return nil, nil, &response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		}
	}
	if body.NewEmail != "" || body.NewPassword != "" {
		errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if errPassword != nil {
			return nil, nil, &response.Error{
				Message: ErrIncorrectPassword.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
	}
	if body.NewPassword != "" {
		dataPassword, errHashed := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if errHashed != nil {
			return nil, nil, &response.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		body.NewPassword = string(dataPassword)
	}
	if body.Name != "" && body.NewEmail == "" && body.NewPassword == "" {
		updateUser, errUpdateUser := s.Repo.UpdateMyUserOneColumn(userUUID, "name", body.Name)
		if errUpdateUser != nil {
			return nil, nil, &response.Error{
				Message: ErrUpdateUser.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return &ResponseUser{
			CreatedAt: updateUser.CreatedAt,
			Name:      updateUser.Name,
			Email:     updateUser.Email,
			UserUUID:  updateUser.UserUUID,
		}, nil, nil
	} else if body.NewEmail != "" {
		token, errSecurity := s.helperSecurity(body.NewEmail, actionUpdate, body)
		if errSecurity != nil {
			return nil, nil, &response.Error{
				Message: errSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return nil, &common.ResponseAuth{
			Message: common.MessageEmail + body.NewEmail,
			JWTTemp: token,
		}, nil
	} else {
		token, errSecurity := s.helperSecurity(user.Email, actionUpdate, body)
		if errSecurity != nil {
			return nil, nil, &response.Error{
				Message: errSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return nil, &common.ResponseAuth{
			Message: common.MessageEmail + user.Email,
			JWTTemp: token,
		}, nil
	}
}
func (s *ServiceUser) helperSecurity(email, action string, bodyUpdate *RequestUpdateUser) (string, error) {
	code := generate_random.GenerateNumbers(common.LengthTempCode)
	sessionId := generate_random.GenerateString(common.LengthSession)
	errSend := common.SendEmailLetter(email, uint(code), s.Dep.Configs)
	if errSend != nil {
		return "", errSend
	}
	j := JWT.NewJWT(s.Dep.Signature)
	token, errCreateJWT := j.CreateTemporaryJWT(sessionId)
	if errCreateJWT != nil {
		return "", custom_errors.ErrFailedSecurity
	}
	if action == actionDelete || action == actionRemove {
		errCreateSession := s.Repo.CreateSessionDeleteOrRemove(uint(code), sessionId, action)
		if errCreateSession != nil {
			return "", custom_errors.ErrFailedSecurity
		}
		return token, nil
	} else if action == actionUpdate {
		errCreateSession := s.Repo.CreateSessionUpdate(&model.TemporaryData{
			TempCode:  uint(code),
			IDSession: sessionId,
			Name:      bodyUpdate.Name,
			Email:     bodyUpdate.NewEmail,
			Password:  bodyUpdate.NewPassword,
		})
		if errCreateSession != nil {
			return "", custom_errors.ErrFailedSecurity
		}
		return token, nil
	}
	return "", ErrIncorrectAction
}
func (s *ServiceUser) ConfirmMyUser(userUUID, sessionID, action string, code uint) (*model.User, *response.Error) {
	user, errGetUser := s.Repo.GetUserByUUID(userUUID)
	if errGetUser != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrUserNotExist.Error(),
			Status:  http.StatusNotFound,
		}
	}
	tempData, errGetSession := s.Repo.GetSession(sessionID, action)
	if errGetSession != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrSession.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	if code != tempData.TempCode {
		return nil, &response.Error{
			Message: custom_errors.ErrIncorrectCode.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	switch action {
	case actionDelete:
		errDeleteUser := s.Repo.DeleteMyUser(userUUID)
		if errDeleteUser != nil {
			return nil, &response.Error{
				Message: ErrFailedDeleteUser.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return nil, nil
	case actionRemove:
		errRemoveUser := s.Repo.RemoveMyUser(userUUID)
		if errRemoveUser != nil {
			return nil, &response.Error{
				Message: ErrFailedRemoveUser.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return nil, nil
	case actionUpdate:
		if tempData.Name != "" && tempData.Email != "" && tempData.Password != "" {
			resUser := model.User{
				Name:     tempData.Name,
				Email:    tempData.Email,
				Password: tempData.Password,
				UserUUID: userUUID,
			}
			errUpdate := s.Repo.UpdateMyUserFull(&resUser)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return &resUser, nil
		} else if tempData.Name != "" && tempData.Email == "" && tempData.Password != "" {
			resUser := model.User{
				Name:     tempData.Name,
				Email:    user.Email,
				Password: tempData.Password,
				UserUUID: userUUID,
			}
			errUpdate := s.Repo.UpdateMyUserFull(&resUser)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return &resUser, nil
		} else if tempData.Name != "" && tempData.Email != "" && tempData.Password == "" {
			resUser := model.User{
				Name:     tempData.Name,
				Email:    tempData.Email,
				Password: user.Password,
				UserUUID: userUUID,
			}
			errUpdate := s.Repo.UpdateMyUserFull(&resUser)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return &resUser, nil
		} else if tempData.Name == "" && tempData.Email != "" && tempData.Password != "" {
			resUser := model.User{
				Name:     user.Name,
				Email:    tempData.Email,
				Password: tempData.Password,
				UserUUID: userUUID,
			}
			errUpdate := s.Repo.UpdateMyUserFull(&resUser)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return &resUser, nil
		} else if tempData.Name == "" && tempData.Email != "" && tempData.Password == "" {
			resUser, errUpdate := s.Repo.UpdateMyUserOneColumn(userUUID, "email", tempData.Email)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return resUser, nil
		} else if tempData.Name == "" && tempData.Email == "" && tempData.Password != "" {
			resUser, errUpdate := s.Repo.UpdateMyUserOneColumn(userUUID, "password", tempData.Password)
			if errUpdate != nil {
				return nil, &response.Error{
					Message: ErrUpdateUser.Error(),
					Status:  http.StatusInternalServerError,
				}
			}
			return resUser, nil
		}
	default:
		return nil, &response.Error{
			Message: ErrIncorrectAction.Error(),
			Status:  http.StatusBadRequest,
		}
	}
	return nil, &response.Error{
		Message: handler_request.ErrInvalidData.Error(),
		Status:  http.StatusBadRequest,
	}
}
func (s *ServiceUser) DeletingRemoveUser() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	select {
	case <-ticker.C:
		s.Repo.deleteExpiredUser()
	}
}
