package article_default

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerArticle struct {
	model.ArticleArchive
	model.ArticleToday
	ResponseCategoryToday
	ResponseCategoryArchive
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
		category := request.PathValue("category")
		limitStr := request.URL.Query().Get("limit")
		filterArticles := request.URL.Query().Get("onlyArticles")
		withText := request.URL.Query().Get("withText")
		if (filterArticles != "" && filterArticles != "false" && filterArticles != "true") || (withText != "" && withText != "false" && withText != "true") {
			h.ResponseCategoryToday.Error = ErrChoiceArticlesFilter.Error()
			handler_response.HandlerResponse(writer, h.ResponseCategoryToday, http.StatusBadRequest)
			return
		}
		allArticle, errGetAllArticle := h.Dep.ServiceArticle.GetArticlesInCategoryToday(category, limitStr, filterArticles, withText)
		if errGetAllArticle != nil {
			h.ResponseCategoryToday.Error = errGetAllArticle.Error()
			switch errGetAllArticle {
			case ErrLoadArticles:
				handler_response.HandlerResponse(writer, h.ResponseCategoryToday, http.StatusInternalServerError)
			case custom_errors.ErrIncorrectParams, ErrCategory:
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
			case custom_errors.ErrIncorrectId:
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
			case custom_errors.ErrIncorrectParams, custom_errors.ErrIncorrectDate:
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
