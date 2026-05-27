package article_default

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/response"
	"net/http"
)

type HandlerArticle struct {
	response.Response[any]
	*ServiceArticle
	//	Dep *HandlerArticleDep
}

//type HandlerArticleDep struct {
//	*ServiceArticle
//}

func NewHandlerArticle(router *http.ServeMux, service *ServiceArticle) { //dep *HandlerArticleDep) {
	article := &HandlerArticle{
		ServiceArticle: service,
	}
	router.HandleFunc("GET /article/today/{category}", article.GetArticlesInCategoryToday())
	router.HandleFunc("GET /article/today/text/{id}", article.GetArticleToday())
	router.HandleFunc("GET /article/archive/{category}", article.GetArticlesInCategoryArchive())
	router.HandleFunc("GET /article/archive/text/{uuid}", article.GetArchiveArticle())
}
func (h *HandlerArticle) GetArticlesInCategoryToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		filterArticles := request.URL.Query().Get("isArticles")
		withText := request.URL.Query().Get("withText")
		allArticle, errGetAllArticle := h.ServiceArticle.GetArticlesInCategoryToday(category, offsetStr, limitStr, filterArticles, withText)
		if len(errGetAllArticle) != 0 {
			h.Response.Errors = errGetAllArticle
			if len(errGetAllArticle) == 1 && errGetAllArticle[0].Message == ErrNotFoundArticle.Error() {
				response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = allArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticleToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		idStr := request.PathValue("id")
		if len(idStr) != lengthIdArticle {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectArticleId.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		article, errGetArticle := h.ServiceArticle.GetArticleToday(idStr)
		if errGetArticle != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetArticle)
			if errGetArticle.Message == ErrLoadArticles.Error() {
				response.HandlerResponse(writer, h.Response, http.StatusInternalServerError)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = article
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticlesInCategoryArchive() http.HandlerFunc {
	defer func() {
		h.Response = response.Response[any]{}
	}()
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		dateStr := request.URL.Query().Get("date")
		articlesArchive, errGetArchive := h.ServiceArticle.GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr)
		if len(errGetArchive) != 0 {
			h.Response.Errors = errGetArchive
			if len(errGetArchive) == 1 {
				response.HandlerResponse(writer, h.Response, errGetArchive[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = articlesArchive
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArchiveArticle() http.HandlerFunc {
	defer func() {
		h.Response = response.Response[any]{}
	}()
	return func(writer http.ResponseWriter, request *http.Request) {
		uuid := request.PathValue("uuid")
		if uuid == "" {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectUUIDArticle.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		archArticle, errGetArchArticle := h.ServiceArticle.GetArchiveArticle(uuid)
		if errGetArchArticle != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetArchArticle)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = archArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
