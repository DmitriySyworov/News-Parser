package user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/handler_request"
	"net/http"
)

type HandlerUser struct {
	response.Response[any]
	*HandlerUserDep
}
type HandlerUserDep struct {
	*ServiceUser
	*middleware.ManagerMiddleware
}

func NewHandlerUser(router *http.ServeMux, dep *HandlerUserDep) {
	user := &HandlerUser{
		HandlerUserDep: dep,
	}
	router.Handle("GET /my/user/get", dep.IsAuthJWT(user.GetMyUser()))
	router.Handle("PATCH /my/user/update", dep.IsAuthJWT(user.UpdateMyUser()))
	router.Handle("DELETE /my/user/remove", dep.IsAuthJWT(user.RemoveMyUser()))
	router.Handle("POST /my/user/confirm", dep.IsAuthJWT(dep.IsTemporaryJWT(user.ConfirmMyUser())))
}
func (h *HandlerUser) GetMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		myUser, errGetMyUser := h.ServiceUser.Repo.GetMyUser(ctxTokens.UUID)
		if errGetMyUser != nil {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			})
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = myUser
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerUser) RemoveMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestRemoveOrDelete](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		typeRemove := request.URL.Query().Get("type")
		if typeRemove != actionRemove && typeRemove != actionDelete {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectType.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respAuth, errRemove := h.ServiceUser.RemoveMyUser(ctxTokens.UUID, body.Password, typeRemove)
		if errRemove != nil {
			h.Response.Errors = append(h.Response.Errors, *errRemove)
			if len(h.Response.Errors) != 0 {
				if len(h.Response.Errors) == 1 {
					response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
				} else {
					response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
				}
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respAuth
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerUser) UpdateMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestUpdateUser](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		updateUser, respAuth, errUpdate := h.ServiceUser.UpdateMyUser(body, ctxTokens.UUID)
		if errUpdate != nil {
			h.Response.Errors = append(h.Response.Errors, *errUpdate)
			if len(h.Response.Errors) != 0 {
				if len(h.Response.Errors) == 1 {
					response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
				} else {
					response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
				}
				return
			}
		}
		h.Response.Success = true
		if respAuth != nil {
			h.Response.Data = respAuth
		} else {
			h.Response.Data = updateUser
		}
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerUser) ConfirmMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
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
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		action := request.URL.Query().Get("action")
		if action != actionRemove && action != actionUpdate && action != actionDelete {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectAction.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respConfirm, errConfirm := h.ServiceUser.ConfirmMyUser(ctxTokens.UUID, ctxTokens.SessionID, action, body.Code)
		if errConfirm != nil {
			h.Response.Errors = append(h.Response.Errors, *errConfirm)
			if len(h.Response.Errors) != 0 {
				if len(h.Response.Errors) == 1 {
					response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
				} else {
					response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
				}
				return
			}
		}
		if action == actionUpdate {
			h.Response.Success = true
			h.Response.Data = respConfirm
			response.HandlerResponse(writer, h.Response, http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	}
}
