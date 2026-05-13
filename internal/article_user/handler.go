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
	router.Handle("PATCH /my/article/update/{id}", dep.IsAuthJWT(articleUser.UpdateUserArticle()))
	router.Handle("PATCH /my/article/update/batch", dep.IsAuthJWT(articleUser.UpdateBatchUserArticles()))
	router.Handle("DELETE /my/article/remove/{id}", dep.IsAuthJWT(articleUser.RemoveUserArticle()))
	router.Handle("GET /my/article/remove", dep.IsAuthJWT(articleUser.GetRemoveUserArticle()))
	router.Handle("POST /my/article/recovery/{id}", dep.IsAuthJWT(articleUser.RecoveryUserArticle()))
	router.Handle("GET /my/article/{category}", dep.IsAuthJWT(articleUser.GetAllUserArticles()))
	router.Handle("GET /my/article/text/{id}", dep.IsAuthJWT(articleUser.GetUserArticle()))
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
		sliceUserArticles, errCreateUserArt := h.Dep.CreateUserArticles(body, ctxTokens.UUID, addText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errCreateUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
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
		id := request.URL.Query().Get("id")
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		body, errRequest := handler_request.HandlerRequest[RequestUpdateArticle](request)
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
		updateArticle, errSliceUpdate := h.Dep.UpdateUserArticle(body.Category, ctxTokens.UUID, id, addText, deleteText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errSliceUpdate...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 && h.ResponseError.Errors[0].Message == ErrNotFoundUserArticle.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			} else if len(h.ResponseError.Errors) == 1 && h.ResponseError.Errors[0].Message == ErrFailedUpdateUserArticle.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			} else if len(h.ResponseError.Errors) == 1 && h.ResponseError.Errors[0].Message == ErrFailedParseText.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnprocessableEntity)
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
		updateBatchArticle, errSliceUpdate := h.Dep.UpdateBatchUserArticles(body.Domain, ctxTokens.UUID, addText, deleteText)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errSliceUpdate...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 && h.ResponseError.Errors[0].Message == ErrNotFoundUserArticle.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
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
		idArticleStr := request.PathValue("id")
		userArticle, errGetUserArt := h.Dep.GetUserArticle(ctxTokens.UUID, idArticleStr)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errGetUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
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
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
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
		idArticle := request.PathValue("id")
		allArticleStr := request.URL.Query().Get("allArticle")
		errRemoveUserArt := h.Dep.RemoveUserArticle(idArticle, ctxTokens.UUID, allArticleStr)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errRemoveUserArt...)
		if len(h.ResponseError.Errors) != 0 {
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
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
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
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
		id := request.PathValue("id")
		allArticle := request.URL.Query().Get("allArticle")
		respRecoveryArticle, errRecoveryArticle := h.Dep.RecoveryUserArticle(ctxTokens.UUID, id, allArticle)
		h.ResponseError.Errors = append(h.ResponseError.Errors, errRecoveryArticle...)
		if len(h.ResponseError.Errors) != 0 {
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respRecoveryArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
