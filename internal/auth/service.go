package auth

import (
	"app/news-parser/configs"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/di"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/JWT"
	"app/news-parser/pkg/generate_random"
	"app/news-parser/pkg/send_letter"
	"time"

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
func (s *ServiceAuth) Register(body *RequestRegister) (*ResponseAuth, error) {
	if errUserExist := s.IRepoUser.IsUserExist(body.Name, body.Email); errUserExist != nil {
		return nil, errUserExist
	}
	hashPassword, errHashPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errHashPass != nil {
		return nil, ErrFailedSecurity
	}
	respAuth, errAuth := s.authHelper(body.Name, body.Email, string(hashPassword))
	if errAuth != nil {
		return nil, errAuth
	}
	return respAuth, nil
}
func (s *ServiceAuth) Login(body *RequestLogin) (*ResponseAuth, error) {
	user, errGetUser := s.IRepoUser.GetUserByEmail(body.Email)
	if errGetUser != nil {
		return nil, ErrLoginEmailOrPassword
	}
	errComparePassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errComparePassword != nil {
		return nil, ErrLoginEmailOrPassword
	}
	respAuth, errAuth := s.authHelper(user.Name, user.Email, user.Password)
	if errAuth != nil {
		return nil, errAuth
	}
	return respAuth, nil
}

const (
	actionRegister = "register"
	actionLogin    = "login"

	lengthTempCode = 6
	lengthSession  = 9
)

func (s *ServiceAuth) authHelper(name, email, hashPass string) (*ResponseAuth, error) {
	sessionId := generate_random.GenerateString(lengthSession)
	tempCode := generate_random.GenerateNumbers(lengthTempCode)
	if errSendEmail := s.sendEmailLetter(email, uint(tempCode)); errSendEmail != nil {
		return nil, errSendEmail
	}
	if errTempUserCreate := s.Repo.CreateTemporaryUser(&model.TemporaryUser{
		Name:      name,
		Email:     email,
		Password:  hashPass,
		TempCode:  uint(tempCode),
		IDSession: sessionId,
	}); errTempUserCreate != nil {
		return nil, ErrFailedSecurity
	}
	j := JWT.NewJWT(s.Signature)
	token, errToken := j.CreateTemporaryJWT(sessionId)
	if errToken != nil {
		return nil, ErrFailedSecurity
	}
	return &ResponseAuth{
		Message: "we sent a letter to the specified email: " + email,
		JWTTemp: token,
	}, nil
}
func (s *ServiceAuth) sendEmailLetter(userEmail string, tempCode uint) error {
	after := time.After(time.Second * 30)
	letter := send_letter.NewSenderLetter(s.ApiEmail, s.ApiPassword, s.Address, s.AddressHost)
	go letter.SendEmailLetter(userEmail, tempCode)
	select {
	case <-after:
		return ErrSendLetter
	case errSend := <-letter.ChErr:
		return errSend
	}
}
func (s *ServiceAuth) Confirm(tempCode uint, action, sessionId string) (*ResponseConfirm, error) {
	tempUser, errGetTempUser := s.Repo.GetTemporaryUser(sessionId)
	if errGetTempUser != nil {
		return nil, ErrExpiredSession
	}
	if tempUser.TempCode != tempCode {
		return nil, ErrIncorrectCode
	}
	uuId := uuid.New().String()
	j := JWT.NewJWT(s.Signature)
	switch action {
	case actionRegister:
		if errUserExist := s.IRepoUser.IsUserExist(tempUser.Name, tempUser.Email); errUserExist != nil {
			return nil, errUserExist
		}
		user := &model.User{
			Name:     tempUser.Name,
			Email:    tempUser.Email,
			Password: tempUser.Password,
			UUIDUser: uuId,
		}
		if errCreate := s.IRepoUser.CreateUser(user); errCreate != nil {
			return nil, ErrSaveDataUser
		}
		token, errJWTCreate := j.CreateJWT(uuId)
		if errJWTCreate != nil {
			return nil, ErrFailedSecurity
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case actionLogin:
		user, errGetUser := s.IRepoUser.GetUserByEmail(tempUser.Email)
		if errGetUser != nil {
			return nil, custom_errors.ErrRecordNotFound
		}
		token, errJWTCreate := j.ParseTemporaryJWT(user.UUIDUser)
		if errJWTCreate != nil {
			return nil, ErrFailedSecurity
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	default:
		return nil, custom_errors.ErrIncorrectAction
	}
}
