package user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerUser struct {
	custom_errors.ResponseError
	common.ResponseSuccessful
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
	router.Handle("DELETE /my/user/delete", dep.IsAuthJWT(user.DeleteMyUser()))
	router.Handle("POST /my/user/confirm", dep.IsAuthJWT(user.ConfirmMyUser()))
}
func (h *HandlerUser) GetMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctxValue := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := ctxValue.(string)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		myUser, errGetMyUser := h.ServiceUser.Repo.GetUserByUUID(userUUID)
		if errGetMyUser != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = myUser
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerUser) RemoveMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctxValue := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := ctxValue.(string)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
	}
}
func (h *HandlerUser) DeleteMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctxValue := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := ctxValue.(string)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
	}
}
func (h *HandlerUser) UpdateMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctxValue := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := ctxValue.(string)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
	}
}
