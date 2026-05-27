package auth

import (
	"app/news-parser/configs"
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/JWT"
	"app/news-parser/pkg/generate_random"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth struct {
	Repo *RepositoryAuth
	*ServiceAuthDep
}

type ServiceAuthDep struct {
	di.IRepoUser
	*configs.Configs
}

func NewServiceAuth(repo *RepositoryAuth, dep *ServiceAuthDep) *ServiceAuth {
	return &ServiceAuth{
		Repo:           repo,
		ServiceAuthDep: dep,
	}
}
func (s *ServiceAuth) Register(body *RequestRegister) (*common.ResponseAuth, *response.Error) {
	if errUserExist := s.IRepoUser.IsUserExistByNameAndEmail(body.Name, body.Email); errUserExist != nil {
		return nil, &response.Error{
			Message: errUserExist.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	hashPassword, errHashPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errHashPass != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrFailedSecurity.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	respAuth, errAuth := s.authHelper(body.Name, body.Email, string(hashPassword))
	if errAuth != nil {
		return nil, &response.Error{
			Message: errAuth.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	return respAuth, nil
}
func (s *ServiceAuth) Login(body *RequestLogin) (*ResponseConfirm, *response.Error) {
	user, errGetUser := s.IRepoUser.GetUserByEmail(body.Email)
	if errGetUser != nil {
		return nil, &response.Error{
			Message: ErrLoginEmailOrPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	errComparePassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errComparePassword != nil {
		return nil, &response.Error{
			Message: ErrLoginEmailOrPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	j := JWT.NewJWT(s.Signature)
	token, errJWT := j.CreateJWT(user.UserUUID)
	if errJWT != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrFailedSecurity.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	return &ResponseConfirm{
		JWT: token,
	}, nil
}

const (
	actionRegister         = "register"
	actionRecoveryPassword = "recovery-password"
	actionRecoveryRemove   = "recovery-remove"
)

func (s *ServiceAuth) Recovery(email, newPassword, action string) (*common.ResponseAuth, *response.Error) {
	switch action {
	case actionRecoveryRemove:
		user, errGetUserByEmail := s.IRepoUser.GetRemoveUserByEmail(email)
		if errGetUserByEmail != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			}
		}
		respAuth, errAuth := s.authHelper(user.Name, user.Email, user.Password)
		if errAuth != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return respAuth, nil
	case actionRecoveryPassword:
		if newPassword == "" {
			return nil, &response.Error{
				Message: ErrNotNewPassword.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		hashedPass, errHashes := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if errHashes != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		user, errGetByEmail := s.IRepoUser.GetUserByEmail(email)
		if errGetByEmail != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			}
		}
		respAuth, errAuth := s.authHelper(user.Name, user.Email, string(hashedPass))
		if errAuth != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return respAuth, nil
	default:
		return nil, &response.Error{
			Message: ErrIncorrectActionRecovery.Error(),
			Status:  http.StatusBadRequest,
		}
	}
}
func (s *ServiceAuth) authHelper(name, email, hashPass string) (*common.ResponseAuth, error) {
	sessionId := generate_random.GenerateString(common.LengthSession)
	tempCode := generate_random.GenerateNumbers(common.LengthTempCode)
	if errSendEmail := common.SendEmailLetter(email, uint(tempCode), s.Configs); errSendEmail != nil {
		return nil, errSendEmail
	}
	if errTempUserCreate := s.Repo.CreateTemporaryUser(&model.TemporaryData{
		Name:      name,
		Email:     email,
		Password:  hashPass,
		TempCode:  uint(tempCode),
		IDSession: sessionId,
	}); errTempUserCreate != nil {
		return nil, custom_errors.ErrFailedSecurity
	}
	j := JWT.NewJWT(s.Signature)
	token, errToken := j.CreateTemporaryJWT(sessionId)
	if errToken != nil {
		return nil, custom_errors.ErrFailedSecurity
	}
	return &common.ResponseAuth{
		Message: common.MessageEmail + email,
		JWTTemp: token,
	}, nil
}

func (s *ServiceAuth) Confirm(code uint, action, sessionId string) (*ResponseConfirm, *response.Error) {
	tempUser, errGetTempUser := s.Repo.GetTemporaryUser(sessionId)
	if errGetTempUser != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrSession.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	if tempUser.TempCode != code {
		return nil, &response.Error{
			Message: custom_errors.ErrIncorrectCode.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	j := JWT.NewJWT(s.Signature)
	var UUID string
	switch action {
	case actionRegister:
		if errUserExist := s.IRepoUser.IsUserExistByNameAndEmail(tempUser.Name, tempUser.Email); errUserExist != nil {
			return nil, &response.Error{
				Message: errUserExist.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		uuId := uuid.New().String()
		user := &model.User{
			Name:     tempUser.Name,
			Email:    tempUser.Email,
			Password: tempUser.Password,
			UserUUID: uuId,
		}
		if errCreate := s.IRepoUser.CreateUser(user); errCreate != nil {
			return nil, &response.Error{
				Message: ErrSaveDataUser.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		UUID = uuId
	case actionRecoveryRemove:
		user, errGetUserByEmail := s.IRepoUser.GetRemoveUserByEmail(tempUser.Email)
		if errGetUserByEmail != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			}
		}
		errRecovery := s.IRepoUser.RecoveryUser(user.UserUUID)
		if errRecovery != nil {
			return nil, &response.Error{
				Message: ErrFailedRecovery.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		UUID = user.UserUUID
	case actionRecoveryPassword:
		user, errGetByEmail := s.IRepoUser.GetUserByEmail(tempUser.Email)
		if errGetByEmail != nil {
			return nil, &response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			}
		}
		_, errUpdate := s.IRepoUser.UpdateMyUserOneColumn(user.UserUUID, "password", tempUser.Password)
		if errUpdate != nil {
			return nil, &response.Error{
				Message: ErrFailedChangePassword.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		UUID = user.UserUUID
	default:
		return nil, &response.Error{
			Message: ErrIncorrectActionConfirm.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	token, errJWTCreate := j.CreateJWT(UUID)
	if errJWTCreate != nil {
		return nil, &response.Error{
			Message: custom_errors.ErrFailedSecurity.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &ResponseConfirm{
		JWT: token,
	}, nil
}
