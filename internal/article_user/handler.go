package article_user

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_request"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerArticleUser struct {
	custom_errors.ResponseError
	common.ResponseSuccessful
	Dep *HandlerArticleUserDep
}
type HandlerArticleUserDep struct {
	*ServiceArticleUser
	*middleware.ManagerMiddleware
}

func NewHandlerArticleUser(router *http.ServeMux, dep *HandlerArticleUserDep) {
	articleUser := &HandlerArticleUser{
		Dep: dep,
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
		addText := request.URL.Query().Get("addText")
		body, errRequest := handler_request.HandlerRequest[RequestCreateArticle](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat, handler_request.ErrBodyIsEmpty:
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
		sliceUserArticles, errCreateUserArt := h.Dep.CreateUserArticles(body, ctxTokens.UUID, addText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errCreateUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = sliceUserArticles
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusMultiStatus)
	}
}
func (h *HandlerArticleUser) UpdateUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{Message: custom_errors.ErrIncorrectToken.Error(), Status: http.StatusUnauthorized})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		articleUUID := request.PathValue("uuid")
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		body, errRequest := handler_request.HandlerRequest[RequestUpdateArticle](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrBodyIsEmpty:
				body = &RequestUpdateArticle{}
			case handler_request.ErrIncorrectFormat:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
				return
			case handler_request.ErrInvalidData:
				h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
				return
			}
		}
		updateArticle, errSliceUpdate := h.Dep.UpdateUserArticle(body.Category, ctxTokens.UUID, articleUUID, addText, deleteText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errSliceUpdate...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = updateArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}

func (h *HandlerArticleUser) UpdateBatchUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{Message: custom_errors.ErrIncorrectToken.Error(), Status: http.StatusUnauthorized})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		body, errRequest := handler_request.HandlerRequest[RequestUpdateBatchArticles](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat, handler_request.ErrBodyIsEmpty:
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
			return
		}
		updateBatchArticle, errSliceUpdate := h.Dep.UpdateBatchUserArticles(ctxTokens.UUID, body.Domain, body.Category, addText, deleteText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errSliceUpdate...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = updateBatchArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusMultiStatus)
	}
}
func (h *HandlerArticleUser) GetUserArticle() http.HandlerFunc {
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
		idArticleStr := request.PathValue("uuid")
		userArticle, errGetUserArt := h.Dep.GetUserArticle(ctxTokens.UUID, idArticleStr)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errGetUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = userArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticleUser) GetAllUserArticles() http.HandlerFunc {
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
		category := request.PathValue("category")
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		withText := request.URL.Query().Get("withText")
		respUserArticles, errGetArticles := h.Dep.GetAllUserArticles(ctxTokens.UUID, category, offset, limit, withText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errGetArticles...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respUserArticles
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RemoveUserArticle() http.HandlerFunc {
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
		idArticle := request.PathValue("uuid")
		typeRemove := request.URL.Query().Get("type")
		errRemoveUserArt := h.Dep.RemoveUserArticle(idArticle, ctxTokens.UUID, typeRemove)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errRemoveUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (h *HandlerArticleUser) RemoveAllUserArticle() http.HandlerFunc {
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
		typeRemove := request.URL.Query().Get("type")
		errRemoveAll := h.Dep.RemoveAllUserArticle(ctxTokens.UUID, typeRemove)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errRemoveAll...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
func (h *HandlerArticleUser) GetRemoveUserArticle() http.HandlerFunc {
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
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		respRemoveArticles, errGetRemoveArticles := h.Dep.GetRemoveUserArticle(ctxTokens.UUID, offset, limit)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errGetRemoveArticles...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respRemoveArticles
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RecoveryUserArticle() http.HandlerFunc {
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
		articleUUId := request.PathValue("uuid")
		respRecoveryArticle, errRecoveryArticle := h.Dep.RecoveryUserArticle(ctxTokens.UUID, articleUUId)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errRecoveryArticle...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respRecoveryArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RecoveryAllUserArticle() http.HandlerFunc {
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
		respArticles, errAllRecovery := h.Dep.RecoveryAllUserArticle(ctxTokens.UUID)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errAllRecovery...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respArticles
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
