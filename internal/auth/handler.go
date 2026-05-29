package auth

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/handler_request"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type HandlerAuth struct {
	response.Response[any]
	Dep *HandlerAuthDep
	*ServiceAuth
}
type HandlerAuthDep struct {
	*middleware.ManagerMiddleware
	*loggers.Logger
}

func NewHandlerAuth(router *http.ServeMux, service *ServiceAuth, dep *HandlerAuthDep) {
	auth := &HandlerAuth{
		ServiceAuth: service,
		Dep:         dep,
	}
	router.HandleFunc("POST /auth/register", auth.Register())
	router.HandleFunc("POST /auth/login", auth.Login())
	router.Handle("POST /auth/recovery", auth.Recovery())
	router.Handle("POST /auth/confirm", dep.IsTemporaryJWT(auth.Confirm()))
}
func (h *HandlerAuth) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		body, errRequest := handler_request.HandlerRequest[RequestRegister](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "Password" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectPassword.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					} else if errList.Field() == "Name" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectName.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					} else if errList.Field() == "Email" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectEmail.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		respAuth, errAuth := h.ServiceAuth.Register(body)
		if errAuth != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errAuth)
			h.Response.Errors = append(h.Response.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				response.HandlerResponse(writer, h.Response, http.StatusInternalServerError)
			default:
				response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respAuth
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerAuth) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		body, errRequest := handler_request.HandlerRequest[RequestLogin](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "Password" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectPassword.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					} else if errList.Field() == "Email" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectEmail.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		respAuth, errAuth := h.ServiceAuth.Login(body)
		if errAuth != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errAuth)
			h.Response.Errors = append(h.Response.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				response.HandlerResponse(writer, h.Response, http.StatusInternalServerError)
			default:
				response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respAuth
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerAuth) Recovery() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		action := request.URL.Query().Get("action")
		values.DataLog.MapLog["action"] = action
		if action != actionRecoveryPassword && action != actionRecoveryRemove {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectActionRecovery.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestRecovery](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "NewPassword" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectNewPassword.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					} else if errList.Field() == "Email" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectEmail.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		respAuth, errAuth := h.ServiceAuth.Recovery(body.Email, body.NewPassword, action)
		if errAuth != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errAuth)
			h.Response.Errors = append(h.Response.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				response.HandlerResponse(writer, h.Response, http.StatusInternalServerError)
			default:
				response.HandlerResponse(writer, h.Response, errAuth.Status)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respAuth
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerAuth) Confirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.MapLog["session_id"] = values.SessionID
		if len(values.SessionID) != common.LengthSession {
			err := response.Error{
				Message: custom_errors.ErrSession.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		action := request.URL.Query().Get("action")
		values.DataLog.MapLog["action"] = action
		if action != actionRegister && action != actionRecoveryRemove && action != actionRecoveryPassword {
			err := response.Error{
				Message: ErrIncorrectActionConfirm.Error(),
				Status:  http.StatusBadRequest,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		body, errRequest := handler_request.HandlerRequest[common.RequestConfirm](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "Code" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectCode.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
						response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		respConfirm, errConfirm := h.ServiceAuth.Confirm(body.Code, action, values.SessionID)
		if errConfirm != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errConfirm)
			h.Response.Errors = append(h.Response.Errors, *errConfirm)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respConfirm
		response.HandlerResponse(writer, h.Response, http.StatusCreated)
	}
}
