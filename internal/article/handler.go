package article

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/handler_request"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerArticle struct {
	model.ArticleArchive
	model.ArticleToday
	model.UserArticle
	ResponseUserArticles
	ResponseCategoryToday
	ResponseCategoryArchive
	ResponseUserDelete
	Dep *HandlerArticleDep
}
type HandlerArticleDep struct {
	*ServiceArticle
	*middleware.ManagerMiddleware
}

func NewHandlerArticle(router *http.ServeMux, dep *HandlerArticleDep) {
	article := &HandlerArticle{
		Dep: dep,
	}
	router.HandleFunc("GET /article/today/{category}", article.GetArticlesInCategoryToday())
	router.HandleFunc("GET /article/today/text/{id}", article.GetArticleToday())
	router.HandleFunc("GET /article/archive/{category}", article.GetArticlesInCategoryArchive())
	router.HandleFunc("GET /article/archive/text/{uuid}", article.GetArchiveArticle())
	router.Handle("POST /my/add/article", dep.IsAuthJWT(article.CreateUserArticles()))
	//router.HandleFunc("PATCH /my/update/article/{category}", article.UpdateUserArticle())
	router.Handle("DELETE /my/delete/article/{id}", dep.IsAuthJWT(article.DeleteUserArticle()))
	router.Handle("GET /my/article/{category}", dep.IsAuthJWT(article.GetAllUserArticles()))
	router.Handle("GET /my/article/text/{id}", dep.IsAuthJWT(article.GetUserArticle()))

}
func (h *HandlerArticle) GetArticlesInCategoryToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		limitStr := request.URL.Query().Get("limit")
		filterArticles := request.URL.Query().Get("onlyArticles")
		if filterArticles != "false" && filterArticles != "true" {
			h.ResponseCategoryToday.Error = ErrChoiceArticlesFilter.Error()
			handler_response.HandlerResponse(writer, h.ResponseCategoryToday, http.StatusBadRequest)
			return
		}
		allArticle, errGetAllArticle := h.Dep.ServiceArticle.GetArticlesInCategoryToday(category, limitStr, filterArticles)
		if errGetAllArticle != nil {
			h.ResponseCategoryToday.Error = errGetAllArticle.Error()
			switch errGetAllArticle {
			case ErrLoadArticles:
				handler_response.HandlerResponse(writer, h.ResponseCategoryToday, http.StatusInternalServerError)
			case ErrIncorrectParams, ErrCategory:
				handler_response.HandlerResponse(writer, h.ResponseCategoryToday, http.StatusBadRequest)
			}
			return
		}
		handler_response.HandlerResponse(writer, allArticle, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticleToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idStr := request.PathValue("id")
		if len(idStr) != lengthIdArticle {
			handler_response.HandlerResponse(writer, h.ArticleToday, http.StatusBadRequest)
			return
		}
		article, errGetArticle := h.Dep.ServiceArticle.GetArticleToday(idStr)
		if errGetArticle != nil {
			h.ArticleToday.Error = errGetArticle.Error()
			switch errGetArticle {
			case ErrLoadArticles:
				handler_response.HandlerResponse(writer, h.ArticleToday, http.StatusInternalServerError)
			case custom_errors.ErrRecordNotFound:
				handler_response.HandlerResponse(writer, h.ArticleToday, http.StatusNotFound)
			case ErrIncorrectId:
				handler_response.HandlerResponse(writer, h.ArticleToday, http.StatusBadRequest)
			}
			return
		}
		handler_response.HandlerResponse(writer, article, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticlesInCategoryArchive() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		dateStr := request.URL.Query().Get("date")
		articlesArchive, errGetArchive := h.Dep.ServiceArticle.GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr)
		if errGetArchive != nil {
			switch errGetArchive {
			case ErrIncorrectParams, custom_errors.ErrIncorrectDate:
				handler_response.HandlerResponse(writer, h.ResponseCategoryArchive, http.StatusBadRequest)
			case ErrLoadArticles:
				handler_response.HandlerResponse(writer, h.ResponseCategoryArchive, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, articlesArchive, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArchiveArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		uuid := request.PathValue("uuid")
		if uuid == "" {
			handler_response.HandlerResponse(writer, h.ArticleArchive, http.StatusBadRequest)
			return
		}
		archArticle, errGetArchArticle := h.Dep.ServiceArticle.GetArchiveArticle(uuid)
		if errGetArchArticle != nil {
			h.ArticleArchive.Error = archArticle.Error
			handler_response.HandlerResponse(writer, h.ArticleArchive, http.StatusNotFound)
			return
		}
		handler_response.HandlerResponse(writer, archArticle, http.StatusOK)
	}
}
func (h *HandlerArticle) CreateUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		valueCtx := request.Context().Value(middleware.KeyAuthToken)
		userUUID, ok := valueCtx.(string)
		if !ok {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectToken.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			return
		}
		addText := request.URL.Query().Get("addText")
		var isAddText bool
		if addText == "false" {
			isAddText = false
		} else if addText == "true" {
			isAddText = true
		} else {
			h.ResponseUserArticles.Error = custom_errors.ErrIncorrectAction.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			return
		}
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
		sliceUserArticles, errCreateUserArt := h.Dep.CreateUserArticles(body, userUUID, isAddText)
		if errCreateUserArt != nil {
			h.ResponseUserArticles.Error = errCreateUserArt.Error()
			switch errCreateUserArt {
			case custom_errors.ErrUserNotFound:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnauthorized)
			case ErrFailedToParse:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusUnprocessableEntity)
			}
			return
		}
		handler_response.HandlerResponse(writer, sliceUserArticles, http.StatusMultiStatus)
	}
}
func (h *HandlerArticle) DeleteUserArticle() http.HandlerFunc {
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
		if (idArticle == "" || len(idArticle) != 11) && (allArticleStr != "true" && allArticleStr != "false") {
			h.ResponseUserDelete.Error = ErrIncorrectParams.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusBadRequest)
			return
		}
		errDeleteUserArt := h.Dep.DeleteUserArticle(idArticle, userUUID, idArticle)
		if errDeleteUserArt != nil {
			h.ResponseUserDelete.Error = errDeleteUserArt.Error()
			switch errDeleteUserArt {
			case custom_errors.ErrUserNotFound:
				handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusNotFound)
			case ErrIncorrectId:
				handler_response.HandlerResponse(writer, h.ResponseUserDelete, http.StatusBadRequest)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
func (h *HandlerArticle) GetUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idArticleStr := request.PathValue("id")
		if idArticleStr == "" {
			h.UserArticle.Error = ErrIncorrectParams.Error()
			handler_response.HandlerResponse(writer, h.UserArticle, http.StatusBadRequest)
			return
		}
		userArticle, errGetUserArt := h.Dep.GetUserArticle(idArticleStr)
		if errGetUserArt != nil {
			h.UserArticle.Error = errGetUserArt.Error()
			switch errGetUserArt {
			case ErrIncorrectId:
				handler_response.HandlerResponse(writer, h.UserArticle, http.StatusBadRequest)
			case custom_errors.ErrRecordNotFound:
				handler_response.HandlerResponse(writer, h.UserArticle, http.StatusNotFound)
			}
			return
		}
		handler_response.HandlerResponse(writer, userArticle, http.StatusOK)
	}
}
func (h *HandlerArticle) GetAllUserArticles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		offset := request.URL.Query().Get("offset")
		limit := request.URL.Query().Get("limit")
		withText := request.URL.Query().Get("withText")
		if category == "" && offset == "" && limit == "" && withText == "" {
			h.ResponseUserArticles.Error = ErrIncorrectParams.Error()
			handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			return
		}
		respUserArticles, errGetArticles := h.Dep.GetAllUserArticles(category, offset, limit, withText)
		if errGetArticles != nil {
			h.ResponseUserArticles.Error = errGetArticles.Error()
			switch errGetArticles {
			case ErrIncorrectParams:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusBadRequest)
			case custom_errors.ErrRecordNotFound:
				handler_response.HandlerResponse(writer, h.ResponseUserArticles, http.StatusNotFound)
			}
			return
		}
		handler_response.HandlerResponse(writer, respUserArticles, http.StatusOK)
	}
}
