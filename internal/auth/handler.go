package auth

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/handler_request"
	"net/http"
)

type HandlerAuth struct {
	response.Response[any]
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
			h.Response = response.Response[any]{}
		}()
		body, errRequest := handler_request.HandlerRequest[RequestRegister](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				response.HandlerResponse(writer, h.Response, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Register(body)
		if errAuth != nil {
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
		body, errRequest := handler_request.HandlerRequest[RequestLogin](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				response.HandlerResponse(writer, h.Response, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Login(body)
		if errAuth != nil {
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
		action := request.URL.Query().Get("action")
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
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				response.HandlerResponse(writer, h.Response, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.ServiceAuth.Recovery(body.Email, body.NewPassword, action)
		if errAuth != nil {
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
		action := request.URL.Query().Get("action")
		if action != actionRegister && action != actionRecoveryRemove && action != actionRecoveryPassword {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectActionConfirm.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		valueCtx := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := valueCtx.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[common.RequestConfirm](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				response.HandlerResponse(writer, h.Response, http.StatusUnprocessableEntity)
			}
			return
		}
		respConfirm, errConfirm := h.Dep.Confirm(body.Code, action, ctxTokens.SessionID)
		if errConfirm != nil {
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
