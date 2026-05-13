package user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_request"
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
	router.Handle("POST /my/user/confirm", dep.IsAuthJWT(dep.IsTemporaryJWT(user.ConfirmMyUser())))
}
func (h *HandlerUser) GetMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		myUser, errGetMyUser := h.ServiceUser.Repo.GetMyUser(ctxTokens.UUID)
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
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestRemoveOrDelete](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		typeRemove := request.URL.Query().Get("type")
		if typeRemove != actionRemove && typeRemove != actionDelete {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectType.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respAuth, errRemove := h.ServiceUser.RemoveMyUser(ctxTokens.UUID, body.Password, typeRemove)
		if errRemove != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errRemove)
			if len(h.ResponseError.Errors) != 0 {
				if len(h.ResponseError.Errors) == 1 {
					handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
				} else {
					handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
				}
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respAuth
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerUser) UpdateMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestUpdateUser](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		updateUser, respAuth, errUpdate := h.ServiceUser.UpdateMyUser(body, ctxTokens.UUID)
		if errUpdate != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errUpdate)
			if len(h.ResponseError.Errors) != 0 {
				if len(h.ResponseError.Errors) == 1 {
					handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
				} else {
					handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
				}
				return
			}
		}
		h.ResponseSuccessful.Success = true
		if respAuth != nil {
			h.ResponseSuccessful.Data = respAuth
		} else {
			h.ResponseSuccessful.Data = updateUser
		}
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerUser) ConfirmMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
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
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
			}
		}
		action := request.URL.Query().Get("action")
		if action != actionRemove && action != actionUpdate && action != actionDelete {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectAction.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respConfirm, errConfirm := h.ServiceUser.ConfirmMyUser(ctxTokens.UUID, ctxTokens.SessionID, action, body.Code)
		if errConfirm != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errConfirm)
			if len(h.ResponseError.Errors) != 0 {
				if len(h.ResponseError.Errors) == 1 {
					handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
				} else {
					handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
				}
				return
			}
		}
		if action == actionUpdate {
			h.ResponseSuccessful.Success = true
			h.ResponseSuccessful.Data = respConfirm
			handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	}
}
