package user

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

type HandlerUser struct {
	response.Response
	Dep *HandlerUserDep
	*ServiceUser
}
type HandlerUserDep struct {
	*loggers.Logger
	*middleware.ManagerMiddleware
}

func NewHandlerUser(router *http.ServeMux, service *ServiceUser, dep *HandlerUserDep) {
	user := &HandlerUser{
		ServiceUser: service,
		Dep:         dep,
	}
	router.Handle("GET /my/user/get", dep.IsAuthJWT(user.GetMyUser()))
	router.Handle("PATCH /my/user/update", dep.IsAuthJWT(user.UpdateMyUser()))
	router.Handle("DELETE /my/user/remove", dep.IsAuthJWT(user.RemoveMyUser()))
	router.Handle("POST /my/user/confirm", dep.IsAuthJWT(dep.IsTemporaryJWT(user.ConfirmMyUser())))
}
func (h *HandlerUser) GetMyUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		myUser, errGetMyUser := h.ServiceUser.Repo.GetMyUser(values.UserUUID)
		if errGetMyUser != nil {
			err := response.Error{
				Message: custom_errors.ErrUserNotExist.Error(),
				Status:  http.StatusNotFound,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
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
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestRemoveOrDelete](request)
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
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
			}
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		typeRemove := request.URL.Query().Get("type")
		values.DataLog.MapLog["type_remove"] = typeRemove
		if typeRemove != actionRemove && typeRemove != actionDelete {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectType.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respAuth, errRemove := h.ServiceUser.RemoveMyUser(values.UserUUID, body.Password, typeRemove)
		if errRemove != nil {
			h.Response.Errors = append(h.Response.Errors, *errRemove)
			if len(h.Response.Errors) != 0 {
				values.DataLog.Errors = append(values.DataLog.Errors, *errRemove)
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
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		body, errRequest := handler_request.HandlerRequest[RequestUpdateUser](request)
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
					} else if errList.Field() == "NewPassword" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectNewPassword.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					} else if errList.Field() == "Name" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectName.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["name"] = body.Name
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					} else if errList.Field() == "NewEmail" {
						err := response.Error{
							Message: custom_errors.ErrIncorrectEmail.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["new_email"] = body.NewEmail
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
			}
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		updateUser, respAuth, errUpdate := h.ServiceUser.UpdateMyUser(body, values.UserUUID)
		if errUpdate != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errUpdate)
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
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
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
						values.DataLog.MapLog["code"] = body.Code
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					}
				}
			} else {
				err := response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				}
				values.DataLog.Errors = append(values.DataLog.Errors, err)
				h.Response.Errors = append(h.Response.Errors, err)
			}
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		action := request.URL.Query().Get("action")
		values.DataLog.MapLog["action"] = action
		if action != actionRemove && action != actionUpdate && action != actionDelete {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectAction.Error(),
				Status:  http.StatusBadRequest,
			})
		}
		respConfirm, errConfirm := h.ServiceUser.ConfirmMyUser(values.UserUUID, values.SessionID, action, body.Code)
		if errConfirm != nil {
			h.Response.Errors = append(h.Response.Errors, *errConfirm)
			if len(h.Response.Errors) != 0 {
				values.DataLog.Errors = append(values.DataLog.Errors, *errConfirm)
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
