package article_user

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/handler_request"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type HandlerArticleUser struct {
	response.Response
	Dep     *HandlerArticleUserDep
	Service *ServiceArticleUser
}
type HandlerArticleUserDep struct {
	*middleware.ManagerMiddleware
	Logger *loggers.Logger
}

func NewHandlerArticleUser(router *http.ServeMux, service *ServiceArticleUser, dep *HandlerArticleUserDep) {
	articleUser := &HandlerArticleUser{
		Service: service,
		Dep:     dep,
	}
	router.Handle("POST /my/article/add", dep.IsAuthJWT(articleUser.CreateUserArticles()))
	router.Handle("PATCH /my/article/update/{uuid}", dep.IsAuthJWT(articleUser.UpdateUserArticle()))
	router.Handle("PATCH /my/article/update/batch", dep.IsAuthJWT(articleUser.UpdateBatchUserArticles()))
	router.Handle("DELETE /my/article/remove/{uuid}", dep.IsAuthJWT(articleUser.RemoveUserArticle()))
	router.Handle("DELETE /my/article/remove/all", dep.IsAuthJWT(articleUser.RemoveAllUserArticle()))
	router.Handle("GET /my/article/remove", dep.IsAuthJWT(articleUser.GetRemoveUserArticle()))
	router.Handle("POST /my/article/recovery/{uuid}", dep.IsAuthJWT(articleUser.RecoveryUserArticle()))
	router.Handle("POST /my/article/recovery/all", dep.IsAuthJWT(articleUser.RecoveryAllUserArticle()))
	router.Handle("GET /my/article/{category}", dep.IsAuthJWT(articleUser.GetAllUserArticles()))
	router.Handle("GET /my/article/text/{uuid}", dep.IsAuthJWT(articleUser.GetUserArticle()))
}
func (h *HandlerArticleUser) CreateUserArticles() http.HandlerFunc {
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
		addText := request.URL.Query().Get("addText")
		values.DataLog.MapLog["add_text"] = addText
		body, errRequest := handler_request.HandlerRequest[RequestCreateArticle](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "URL" {
						err := response.Error{
							Message: ErrIncorrectURL.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["url"] = body.URL
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					} else if errList.Field() == "Category" {
						err := response.Error{
							Message: ErrIncorrectCategory.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["category"] = body.Category
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
		sliceUserArticles, errCreateUserArt := h.Service.CreateUserArticles(body, values.UserUUID, addText)
		h.Response.Errors = append(h.Response.Errors, errCreateUserArt...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errCreateUserArt...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = sliceUserArticles
		response.HandlerResponse(writer, h.Response, http.StatusMultiStatus)
	}
}
func (h *HandlerArticleUser) UpdateUserArticle() http.HandlerFunc {
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
		articleUUID := request.PathValue("uuid")
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		values.DataLog.MapLog["article_uuid"] = articleUUID
		values.DataLog.MapLog["add_text"] = addText
		values.DataLog.MapLog["delete_text"] = deleteText
		body, errRequest := handler_request.HandlerRequest[RequestUpdateArticle](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "Category" {
						err := response.Error{
							Message: ErrIncorrectCategory.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["category"] = body.Category
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
		updateArticle, errSliceUpdate := h.Service.UpdateUserArticle(body.Category, values.UserUUID, articleUUID, addText, deleteText)
		h.Response.Errors = append(h.Response.Errors, errSliceUpdate...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errSliceUpdate...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = updateArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

func (h *HandlerArticleUser) UpdateBatchUserArticles() http.HandlerFunc {
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
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		values.DataLog.MapLog["add_text"] = addText
		values.DataLog.MapLog["delete_text"] = deleteText
		body, errRequest := handler_request.HandlerRequest[RequestUpdateBatchArticles](request)
		if errRequest != nil {
			if errValid, isValidErr := errRequest.(validator.ValidationErrors); isValidErr {
				for _, errList := range errValid {
					if errList.Field() == "Category" {
						err := response.Error{
							Message: ErrIncorrectCategory.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["category"] = body.Category
						values.DataLog.Errors = append(values.DataLog.Errors, err)
						h.Response.Errors = append(h.Response.Errors, err)
					} else if errList.Field() == "Domain" {
						err := response.Error{
							Message: ErrIncorrectDomain.Error(),
							Status:  http.StatusBadRequest,
						}
						values.DataLog.MapLog["domain"] = body.Domain
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
		updateBatchArticle, errSliceUpdate := h.Service.UpdateBatchUserArticles(values.UserUUID, body.Domain, body.Category, addText, deleteText)
		values.DataLog.Errors = append(values.DataLog.Errors, errSliceUpdate...)
		h.Response.Errors = append(h.Response.Errors, errSliceUpdate...)
		if len(h.Response.Errors) != 0 {
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = updateBatchArticle
		response.HandlerResponse(writer, h.Response, http.StatusMultiStatus)
	}
}
func (h *HandlerArticleUser) GetUserArticle() http.HandlerFunc {
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
		idArticleStr := request.PathValue("uuid")
		values.DataLog.MapLog["article_uuid"] = idArticleStr
		userArticle, errGetUserArt := h.Service.GetUserArticle(values.UserUUID, idArticleStr)
		h.Response.Errors = append(h.Response.Errors, errGetUserArt...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errGetUserArt...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = userArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticleUser) GetAllUserArticles() http.HandlerFunc {
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
		category := request.PathValue("category")
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		withText := request.URL.Query().Get("withText")
		values.DataLog.MapLog["category"] = category
		values.DataLog.MapLog["offset"] = offset
		values.DataLog.MapLog["limit"] = limit
		values.DataLog.MapLog["with_text"] = withText
		respUserArticles, errGetArticles := h.Service.GetAllUserArticles(values.UserUUID, category, offset, limit, withText)
		h.Response.Errors = append(h.Response.Errors, errGetArticles...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errGetArticles...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respUserArticles
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RemoveUserArticle() http.HandlerFunc {
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
		idArticle := request.PathValue("uuid")
		typeRemove := request.URL.Query().Get("type")
		values.DataLog.MapLog["article_uuid"] = idArticle
		values.DataLog.MapLog["type_remove"] = typeRemove
		errRemoveUserArt := h.Service.RemoveUserArticle(idArticle, values.UserUUID, typeRemove)
		h.Response.Errors = append(h.Response.Errors, errRemoveUserArt...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errRemoveUserArt...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (h *HandlerArticleUser) RemoveAllUserArticle() http.HandlerFunc {
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
		typeRemove := request.URL.Query().Get("type")
		values.DataLog.MapLog["type_remove"] = typeRemove
		errRemoveAll := h.Service.RemoveAllUserArticle(values.UserUUID, typeRemove)
		h.Response.Errors = append(h.Response.Errors, errRemoveAll...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errRemoveAll...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
func (h *HandlerArticleUser) GetRemoveUserArticle() http.HandlerFunc {
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
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		values.DataLog.MapLog["offset"] = offset
		values.DataLog.MapLog["limit"] = limit
		respRemoveArticles, errGetRemoveArticles := h.Service.GetRemoveUserArticle(values.UserUUID, offset, limit)
		h.Response.Errors = append(h.Response.Errors, errGetRemoveArticles...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errGetRemoveArticles...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respRemoveArticles
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RecoveryUserArticle() http.HandlerFunc {
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
		articleUUId := request.PathValue("uuid")
		values.DataLog.MapLog["article_uuid"] = articleUUId
		respRecoveryArticle, errRecoveryArticle := h.Service.RecoveryUserArticle(values.UserUUID, articleUUId)
		h.Response.Errors = append(h.Response.Errors, errRecoveryArticle...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errRecoveryArticle...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respRecoveryArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RecoveryAllUserArticle() http.HandlerFunc {
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
		respArticles, errAllRecovery := h.Service.RecoveryAllUserArticle(values.UserUUID)
		h.Response.Errors = append(h.Response.Errors, errAllRecovery...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errAllRecovery...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respArticles
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
