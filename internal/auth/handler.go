package auth

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_request"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerAuth struct {
	custom_errors.ResponseError
	common.ResponseSuccessful
	Dep *HandlerAuthDep
}
type HandlerAuthDep struct {
	*ServiceAuth
	*middleware.ManagerMiddleware
}

func NewHandlerAuth(router *http.ServeMux, dep *HandlerAuthDep) {
	auth := &HandlerAuth{
		Dep: dep,
	}
	router.HandleFunc("POST /auth/register", auth.Register())
	router.HandleFunc("POST /auth/login", auth.Login())
	router.Handle("POST /auth/recovery", auth.Recovery())
	router.Handle("POST /auth/confirm", dep.IsTemporaryJWT(auth.Confirm()))
}
func (h *HandlerAuth) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		body, errRequest := handler_request.HandlerRequest[RequestRegister](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Register(body)
		if errAuth != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respAuth
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerAuth) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		body, errRequest := handler_request.HandlerRequest[RequestLogin](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Login(body)
		if errAuth != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respAuth
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerAuth) Recovery() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		action := request.URL.Query().Get("action")
		if action != actionRecoveryPassword && action != actionRecoveryRemove {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectActionRecovery.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestRecovery](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.ServiceAuth.Recovery(body.Email, body.NewPassword, action)
		if errAuth != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errAuth)
			switch errAuth.Message {
			case custom_errors.ErrFailedSecurity.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseError, errAuth.Status)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respAuth
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerAuth) Confirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		action := request.URL.Query().Get("action")
		if action != actionRegister && action != actionRecoveryRemove && action != actionRecoveryPassword {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectActionConfirm.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		valueCtx := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := valueCtx.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[common.RequestConfirm](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
			}
			return
		}
		respConfirm, errConfirm := h.Dep.Confirm(body.Code, action, ctxTokens.SessionID)
		if errConfirm != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errConfirm)
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respConfirm
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusCreated)
	}
}
