package article_default

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerArticle struct {
	custom_errors.ResponseError
	common.ResponseSuccessful
	Dep *HandlerArticleDep
}
type HandlerArticleDep struct {
	*ServiceArticle
}

func NewHandlerArticle(router *http.ServeMux, dep *HandlerArticleDep) {
	article := &HandlerArticle{
		Dep: dep,
	}
	router.HandleFunc("GET /article/today/{category}", article.GetArticlesInCategoryToday())
	router.HandleFunc("GET /article/today/text/{id}", article.GetArticleToday())
	router.HandleFunc("GET /article/archive/{category}", article.GetArticlesInCategoryArchive())
	router.HandleFunc("GET /article/archive/text/{uuid}", article.GetArchiveArticle())
}
func (h *HandlerArticle) GetArticlesInCategoryToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		filterArticles := request.URL.Query().Get("isArticles")
		withText := request.URL.Query().Get("withText")
		allArticle, errGetAllArticle := h.Dep.ServiceArticle.GetArticlesInCategoryToday(category, offsetStr, limitStr, filterArticles, withText)
		if len(errGetAllArticle) != 0 {
			h.ResponseError.Errors = errGetAllArticle
			if len(errGetAllArticle) == 1 && errGetAllArticle[0].Message == ErrLoadArticles.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = allArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticleToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		idStr := request.PathValue("id")
		if len(idStr) != lengthIdArticle {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectArticleId.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		article, errGetArticle := h.Dep.ServiceArticle.GetArticleToday(idStr)
		if errGetArticle != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetArticle)
			if errGetArticle.Message == ErrLoadArticles.Error() {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusInternalServerError)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = article
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticlesInCategoryArchive() http.HandlerFunc {
	defer func() {
		h.ResponseError = custom_errors.ResponseError{}
		h.ResponseSuccessful = common.ResponseSuccessful{}
	}()
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		dateStr := request.URL.Query().Get("date")
		articlesArchive, errGetArchive := h.Dep.ServiceArticle.GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr)
		if len(errGetArchive) != 0 {
			h.ResponseError.Errors = errGetArchive
			if len(errGetArchive) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, errGetArchive[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = articlesArchive
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArchiveArticle() http.HandlerFunc {
	defer func() {
		h.ResponseError = custom_errors.ResponseError{}
		h.ResponseSuccessful = common.ResponseSuccessful{}
	}()
	return func(writer http.ResponseWriter, request *http.Request) {
		uuid := request.PathValue("uuid")
		if uuid == "" {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectUUIDArticle.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		archArticle, errGetArchArticle := h.Dep.ServiceArticle.GetArchiveArticle(uuid)
		if errGetArchArticle != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetArchArticle)
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = archArticle
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
