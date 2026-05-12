package auth

import (
	"app/news-parser/configs"
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
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
func (s *ServiceAuth) Register(body *RequestRegister) (*common.ResponseAuth, *custom_errors.Error) {
	if errUserExist := s.IRepoUser.IsUserExistByNameAndEmail(body.Name, body.Email); errUserExist != nil {
		return nil, &custom_errors.Error{
			Message: errUserExist.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	hashPassword, errHashPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errHashPass != nil {
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrFailedSecurity.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	respAuth, errAuth := s.authHelper(body.Name, body.Email, string(hashPassword))
	if errAuth != nil {
		return nil, &custom_errors.Error{
			Message: errAuth.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	return respAuth, nil
}
func (s *ServiceAuth) Login(body *RequestLogin) (*common.ResponseAuth, *custom_errors.Error) {
	user, errGetUser := s.IRepoUser.GetUserByEmail(body.Email)
	if errGetUser != nil {
		return nil, &custom_errors.Error{
			Message: ErrLoginEmailOrPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	errComparePassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errComparePassword != nil {
		return nil, &custom_errors.Error{
			Message: ErrLoginEmailOrPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	respAuth, errAuth := s.authHelper(user.Name, user.Email, user.Password)
	if errAuth != nil {
		return nil, &custom_errors.Error{
			Message: errAuth.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	return respAuth, nil
}

const (
	actionRegister = "register"
	actionLogin    = "login"
)

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

func (s *ServiceAuth) Confirm(tempCode uint, action, sessionId string) (*ResponseConfirm, *custom_errors.Error) {
	tempUser, errGetTempUser := s.Repo.GetTemporaryUser(sessionId)
	if errGetTempUser != nil {
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrSession.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	if tempUser.TempCode != tempCode {
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrIncorrectCode.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
	uuId := uuid.New().String()
	j := JWT.NewJWT(s.Signature)
	switch action {
	case actionRegister:
		if errUserExist := s.IRepoUser.IsUserExistByNameAndEmail(tempUser.Name, tempUser.Email); errUserExist != nil {
			return nil, &custom_errors.Error{
				Message: errUserExist.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		user := &model.User{
			Name:     tempUser.Name,
			Email:    tempUser.Email,
			Password: tempUser.Password,
			UUIDUser: uuId,
		}
		if errCreate := s.IRepoUser.CreateUser(user); errCreate != nil {
			return nil, &custom_errors.Error{
				Message: ErrSaveDataUser.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		token, errJWTCreate := j.CreateJWT(uuId)
		if errJWTCreate != nil {
			return nil, &custom_errors.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case actionLogin:
		user, errGetUser := s.IRepoUser.GetUserByEmail(tempUser.Email)
		if errGetUser != nil {
			return nil, &custom_errors.Error{
				Message: custom_errors.ErrRecordNotFound.Error(),
				Status:  http.StatusUnauthorized,
			}
		}
		token, errJWTCreate := j.ParseTemporaryJWT(user.UUIDUser)
		if errJWTCreate != nil {
			return nil, &custom_errors.Error{
				Message: custom_errors.ErrFailedSecurity.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	default:
		return nil, &custom_errors.Error{
			Message: custom_errors.ErrIncorrectAction.Error(),
			Status:  http.StatusUnauthorized,
		}
	}
}
