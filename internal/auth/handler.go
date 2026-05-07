package auth

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_request"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerAuth struct {
	ResponseAuth
	ResponseConfirm
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
	router.Handle("POST /auth/confirm", dep.IsTemporaryJWT(auth.Confirm()))
}
func (h *HandlerAuth) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errRequest := handler_request.HandlerRequest[RequestRegister](request)
		if errRequest != nil {
			h.ResponseAuth.Error = errRequest.Error()
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Register(body)
		if errAuth != nil {
			h.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case ErrFailedSecurity:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusUnauthorized)
			}
			return
		}
		handler_response.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (h *HandlerAuth) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errRequest := handler_request.HandlerRequest[RequestLogin](request)
		if errRequest != nil {
			h.ResponseAuth.Error = errRequest.Error()
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusUnprocessableEntity)
			}
			return
		}
		respAuth, errAuth := h.Dep.Login(body)
		if errAuth != nil {
			h.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case ErrFailedSecurity:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseAuth, http.StatusUnauthorized)
			}
			return
		}
		handler_response.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (h *HandlerAuth) Confirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		action := request.URL.Query().Get("action")
		if action != actionLogin && action != actionRegister {
			h.ResponseConfirm.Error = custom_errors.ErrIncorrectAction.Error()
			handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusBadRequest)
			return
		}
		valueCtx := request.Context().Value(middleware.KeySessionJWT)
		sessionId, ok := valueCtx.(string)
		if !ok {
			h.ResponseConfirm.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestConfirm](request)
		if errRequest != nil {
			h.ResponseConfirm.Error = errRequest.Error()
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusUnprocessableEntity)
			}
			return
		}
		respConfirm, errConfirm := h.Dep.Confirm(body.Code, action, sessionId)
		if errConfirm != nil {
			h.ResponseConfirm.Error = errConfirm.Error()
			switch errConfirm {
			case ErrSaveDataUser, ErrFailedSecurity:
				handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusInternalServerError)
			default:
				handler_response.HandlerResponse(writer, h.ResponseConfirm, http.StatusUnauthorized)
			}
			return
		}
		handler_response.HandlerResponse(writer, respConfirm, http.StatusCreated)
	}
}
