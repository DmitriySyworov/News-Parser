package article_user

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/handler_request"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerArticleUser struct {
	ResponseRemoveUserArticles
	ResponseUserArticles
	ResponseUserDelete
	model.UserArticle
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
	router.Handle("DELETE /my/article/remove/{id}", dep.IsAuthJWT(articleUser.RemoveUserArticle()))
	router.Handle("GET /my/article/remove", dep.IsAuthJWT(articleUser.GetRemoveUserArticle()))
	router.Handle("POST /my/article/recovery/{id}", dep.IsAuthJWT(articleUser.RecoveryUserArticle()))
	router.Handle("GET /my/article/{category}", dep.IsAuthJWT(articleUser.GetAllUserArticles()))
	router.Handle("GET /my/article/text/{id}", dep.IsAuthJWT(articleUser.GetUserArticle()))
}
func (h *HandlerArticleUser) CreateUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		addText := request.URL.Query().Get("addText")
		body, errRequest := handler_request.HandlerRequest[RequestCreateArticle](request)
		if errRequest != nil {
			h.ResponseUserArticles.Error = errRequest.Error()
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnprocessableEntity)
			}
			return
		}
		sliceUserArticles, errCreateUserArt := h.Dep.CreateUserArticles(body, userUUID, addText)
		if errCreateUserArt != nil {
			h.ResponseUserArticles.Error = errCreateUserArt.Error()
			switch errCreateUserArt {
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			case ErrIncorrectAddText:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			case ErrFailedToParse:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnprocessableEntity)
			}
			return
		}
		handler_response.HandlerResponse(writer, sliceUserArticles, http.StatusMultiStatus)
	}
}
func (h *HandlerArticleUser) UpdateUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		id := request.URL.Query().Get("id")
		addText := request.URL.Query().Get("addText")
		deleteText := request.URL.Query().Get("deleteText")
		allArticle := request.URL.Query().Get("allArticle")
		body, errRequest := handler_request.HandlerRequest[RequestUpdateArticle](request)
		if errRequest != nil {
			h.ResponseUserArticles.Error = errRequest.Error()
			switch errRequest {
			case handler_request.ErrIncorrectFormat:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			case handler_request.ErrInvalidData:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnprocessableEntity)
			}
			return
		}
		h.Dep.UpdateUserArticle(body, userUUID, id, addText, deleteText, allArticle)
	}
}
func (h *HandlerArticleUser) GetUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		idArticleStr := request.PathValue("id")
		if idArticleStr == "" {
			h.UserArticle.Error = ErrIncorrectArticleId.Error()
			handler_response.HandlerResponse(writer, h.UserArticle, http.StatusBadRequest)
			return
		}
		userArticle, errGetUserArt := h.Dep.GetUserArticle(userUUID, idArticleStr)
		if errGetUserArt != nil {
			h.UserArticle.Error = errGetUserArt.Error()
			switch errGetUserArt {
			case ErrIncorrectArticleId:
				handler_response.HandlerResponse(writer, h.UserArticle, http.StatusBadRequest)
			case ErrNotFoundUserArticle:
				handler_response.HandlerResponse(writer, h.UserArticle, http.StatusNotFound)
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.UserArticle, http.StatusUnauthorized)
			}
			return
		}
		handler_response.HandlerResponse(writer, userArticle, http.StatusOK)
	}
}
func (h *HandlerArticleUser) GetAllUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		category := request.PathValue("category")
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		withText := request.URL.Query().Get("withText")
		respUserArticles, errGetArticles := h.Dep.GetAllUserArticles(userUUID, category, offset, limit, withText)
		if errGetArticles != nil {
			h.ResponseUserArticles.Error = errGetArticles.Error()
			switch errGetArticles {
			case custom_errors.ErrIncorrectOffsetAndLimit, custom_errors.ErrIncorrectLimit, custom_errors.ErrIncorrectOffset, ErrIncorrectWithText:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			case ErrNotFoundUserArticle:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusNotFound)
			}
			return
		}
		handler_response.HandlerResponse(writer, respUserArticles, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RemoveUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		idArticle := request.PathValue("id")
		allArticleStr := request.URL.Query().Get("allArticle")
		errRemoveUserArt := h.Dep.RemoveUserArticle(idArticle, userUUID, allArticleStr)
		if errRemoveUserArt != nil {
			h.ResponseUserDelete.Error = errRemoveUserArt.Error()
			switch errRemoveUserArt {
			case ErrFailedRemoveArticle:
				handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusNotFound)
			case ErrIncorrectArticleId, ErrIncorrectAllArticle:
				handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusBadRequest)
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusUnauthorized)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
func (h *HandlerArticleUser) GetRemoveUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseRemoveUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusUnauthorized)
			return
		}
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		respRemoveArticles, errGetRemoveArticles := h.Dep.GetRemoveUserArticle(userUUID, offset, limit)
		if errGetRemoveArticles != nil {
			h.ResponseRemoveUserArticles.Error = errGetRemoveArticles.Error()
			switch errGetRemoveArticles {
			case custom_errors.ErrIncorrectOffsetAndLimit, custom_errors.ErrIncorrectLimit, custom_errors.ErrIncorrectOffset:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusBadRequest)
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusUnauthorized)
			case ErrNotFoundRemoveArticles:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusNotFound)
			}
			return
		}
		handler_response.HandlerResponse(writer, respRemoveArticles, http.StatusOK)
	}
}
func (h *HandlerArticleUser) RecoveryUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseRemoveUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusUnauthorized)
			return
		}
		id := request.PathValue("id")
		allArticle := request.URL.Query().Get("allArticle")
		respRecoveryArticle, errRecoveryArticle := h.Dep.RecoveryUserArticle(userUUID, id, allArticle)
		if errRecoveryArticle != nil {
			h.ResponseUserArticles.Error = errRecoveryArticle.Error()
			switch errRecoveryArticle {
			case ErrIncorrectArticleId, ErrIncorrectAllArticle, ErrIdAndAllArticleParams:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusBadRequest)
			case ErrFailedRecoveryArticle:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusNotFound)
			case custom_errors.ErrUserNotExist:
				handler_response.HandlerResponse(writer, h.ResponseRemoveUserArticles, http.StatusUnauthorized)
			}
			return
		}
		handler_response.HandlerResponse(writer, respRecoveryArticle, http.StatusOK)
	}
}
