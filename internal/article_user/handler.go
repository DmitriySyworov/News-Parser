package article_user

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"app/news-parser/pkg/handler_request"
	"net/http"
)

type HandlerArticleUser struct {
	response.Response[any]
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
		addText := request.URL.Query().Get("addText")
		body, errRequest := handler_request.HandlerRequest[RequestCreateArticle](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat, handler_request.ErrBodyIsEmpty:
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
		sliceUserArticles, errCreateUserArt := h.Dep.CreateUserArticles(body, ctxTokens.UUID, addText)
		h.Response.Errors = append(h.Response.Errors, errCreateUserArt...)
		if len(h.Response.Errors) != 0 {
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
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{Message: custom_errors.ErrIncorrectToken.Error(), Status: http.StatusUnauthorized})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
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
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusBadRequest,
				})
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
				return
			case handler_request.ErrInvalidData:
				h.Response.Errors = append(h.Response.Errors, response.Error{
					Message: errRequest.Error(),
					Status:  http.StatusUnprocessableEntity,
				})
				response.HandlerResponse(writer, h.Response, http.StatusUnprocessableEntity)
				return
			}
		}
		updateArticle, errSliceUpdate := h.Dep.UpdateUserArticle(body.Category, ctxTokens.UUID, articleUUID, addText, deleteText)
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
		h.Response.Data = updateArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

func (h *HandlerArticleUser) UpdateBatchUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.Response.Errors = append(h.Response.Errors, response.Error{Message: custom_errors.ErrIncorrectToken.Error(), Status: http.StatusUnauthorized})
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		body, errRequest := handler_request.HandlerRequest[RequestUpdateBatchArticles](request)
		if errRequest != nil {
			switch errRequest {
			case handler_request.ErrIncorrectFormat, handler_request.ErrBodyIsEmpty:
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
		updateBatchArticle, errSliceUpdate := h.Dep.UpdateBatchUserArticles(ctxTokens.UUID, body.Domain, body.Category, addText, deleteText)
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
		idArticleStr := request.PathValue("uuid")
		userArticle, errGetUserArt := h.Dep.GetUserArticle(ctxTokens.UUID, idArticleStr)
		h.Response.Errors = append(h.Response.Errors, errGetUserArt...)
		if len(h.Response.Errors) != 0 {
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
		category := request.PathValue("category")
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		withText := request.URL.Query().Get("withText")
		respUserArticles, errGetArticles := h.Dep.GetAllUserArticles(ctxTokens.UUID, category, offset, limit, withText)
		h.Response.Errors = append(h.Response.Errors, errGetArticles...)
		if len(h.Response.Errors) != 0 {
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
		idArticle := request.PathValue("uuid")
		typeRemove := request.URL.Query().Get("type")
		errRemoveUserArt := h.Dep.RemoveUserArticle(idArticle, ctxTokens.UUID, typeRemove)
		h.Response.Errors = append(h.Response.Errors, errRemoveUserArt...)
		if len(h.Response.Errors) != 0 {
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
		typeRemove := request.URL.Query().Get("type")
		errRemoveAll := h.Dep.RemoveAllUserArticle(ctxTokens.UUID, typeRemove)
		h.Response.Errors = append(h.Response.Errors, errRemoveAll...)
		if len(h.Response.Errors) != 0 {
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
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		respRemoveArticles, errGetRemoveArticles := h.Dep.GetRemoveUserArticle(ctxTokens.UUID, offset, limit)
		h.Response.Errors = append(h.Response.Errors, errGetRemoveArticles...)
		if len(h.Response.Errors) != 0 {
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
		articleUUId := request.PathValue("uuid")
		respRecoveryArticle, errRecoveryArticle := h.Dep.RecoveryUserArticle(ctxTokens.UUID, articleUUId)
		h.Response.Errors = append(h.Response.Errors, errRecoveryArticle...)
		if len(h.Response.Errors) != 0 {
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
		respArticles, errAllRecovery := h.Dep.RecoveryAllUserArticle(ctxTokens.UUID)
		h.Response.Errors = append(h.Response.Errors, errAllRecovery...)
		if len(h.Response.Errors) != 0 {
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
