package article_default

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"log/slog"
	"net/http"
)

type HandlerArticle struct {
	response.Response[any]
	*ServiceArticle
	Dep *HandlerArticleDep
}
type HandlerArticleDep struct {
	Logger *loggers.Logger
}

func NewHandlerArticle(router *http.ServeMux, service *ServiceArticle, dep *HandlerArticleDep) {
	article := &HandlerArticle{
		ServiceArticle: service,
		Dep:            dep,
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
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		filterArticles := request.URL.Query().Get("isArticles")
		withText := request.URL.Query().Get("withText")
		values.DataLog.MapLog["category"] = category
		values.DataLog.MapLog["offset"] = offsetStr
		values.DataLog.MapLog["limit"] = limitStr
		values.DataLog.MapLog["is_articles"] = filterArticles
		values.DataLog.MapLog["with_text"] = withText
		allArticle, errGetAllArticle := h.ServiceArticle.GetArticlesInCategoryToday(category, offsetStr, limitStr, filterArticles, withText)
		if len(errGetAllArticle) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errGetAllArticle...)
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
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		idStr := request.PathValue("id")
		values.DataLog.MapLog["id"] = idStr
		if len(idStr) != lengthIdArticle {
			err := response.Error{
				Message: custom_errors.ErrIncorrectArticleId.Error(),
				Status:  http.StatusBadRequest,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		article, errGetArticle := h.ServiceArticle.GetArticleToday(idStr)
		if errGetArticle != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetArticle)
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
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		category := request.PathValue("category")
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		dateStr := request.URL.Query().Get("date")
		values.DataLog.MapLog["category"] = category
		values.DataLog.MapLog["offset"] = offsetStr
		values.DataLog.MapLog["limit"] = limitStr
		values.DataLog.MapLog["date"] = dateStr
		articlesArchive, errGetArchive := h.ServiceArticle.GetArticlesInCategoryArchive(category, offsetStr, limitStr, dateStr)
		if len(errGetArchive) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, errGetArchive...)
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
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		uuid := request.PathValue("uuid")
		values.DataLog.MapLog["article_uuid"] = uuid
		if uuid == "" {
			err := response.Error{
				Message: ErrIncorrectUUIDArticle.Error(),
				Status:  http.StatusBadRequest,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		archArticle, errGetArchArticle := h.ServiceArticle.GetArchiveArticle(uuid)
		if errGetArchArticle != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetArchArticle)
			h.Response.Errors = append(h.Response.Errors, *errGetArchArticle)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = archArticle
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
