package article_default

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	_ "app/news-parser/internal/model"
	"app/news-parser/internal/response"
	"log/slog"
	"net/http"
)

type HandlerArticle struct {
	response.Response
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
	router.HandleFunc("GET /api/v1/article/today/{category}", article.GetArticlesInCategoryToday())
	router.HandleFunc("GET /api/v1/article/today/text/{id}", article.GetArticleToday())
	router.HandleFunc("GET /api/v1/article/archive/{category}", article.GetArticlesInCategoryArchive())
	router.HandleFunc("GET /api/v1/article/archive/text/{uuid}", article.GetArchiveArticle())
}

// GetArticlesInCategoryToday godoc
// @Summary      get today articles
// @Description	 get today's articles by category
// @Tags         article_default
// @Produce      json
// @Param        category  path  string true "food, politics, sport, cloth, electronics, the following categories are valid, as well as all to get all entries"
// @Param 		 offset query  int  false  "offset default 0"
// @Param 		 limit  query  int false "limit default 50"
// @Param 		 is_articles query  bool  false "is_articles default true"
// @Param 		 with_text query bool false  "with_text default false"
// @Success      200     {object} response.Response{data=[]ResponseCategoryToday}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `is_articles` - the isArticles must be a boolean is true or false | `with_text` - the withText must be a boolean is true or false | `offset` - the offset must be a positive integer | `limit` - the limit must be a positive integer | `category` - this category of articles does not exist
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_articles` - article not found
// @Router       /api/v1/article/today/{category} [get]
func (h *HandlerArticle) GetArticlesInCategoryToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
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

// GetArticleToday godoc
// @Summary      get article with text
// @Description	  get one today article with text
// @Tags         article_default
// @Produce      json
// @Param 		 id    path   int  true "id of today's article must be 7 digits"
// @Success      200     {object} response.Response{data=model.ArticleToday}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `id` - id must be 7 digits long
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_articles` - article not found
// @Router      /api/v1/article/today/text/{id} [get]
func (h *HandlerArticle) GetArticleToday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
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
				Message: ErrIncorrectIDArticleToday.Error(),
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
			if errGetArticle.Message == ErrNotFoundArticle.Error() {
				response.HandlerResponse(writer, h.Response, http.StatusNotFound)
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

// GetArticlesInCategoryArchive godoc
// @Summary      get archive articles
// @Description	 get archive articles by category
// @Tags         article_default
// @Produce      json
// @Param        category  path  string true  "food, politics, sport, cloth, electronics, the following categories are valid, as well as all to get all entries"
// @Param 		 offset query  int false   "offset default 0"
// @Param 		 limit  query  int  false "limit default 50"
// @Param 		 date   query  string  true "date must be YYYY-MM-DD"
// @Success      200     {object} response.Response{data=[]ResponseCategoryArchive}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `offset` - the offset must be a positive integer | `limit` - the limit must be a positive integer | `category` - this category of articles does not exist | `date` - the date format must be YYYY-MM-DD
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_articles` - article not found
// @Router       /api/v1/article/archive/{category} [get]
func (h *HandlerArticle) GetArticlesInCategoryArchive() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
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

// GetArchiveArticle godoc
// @Summary      get archive article
// @Description	  get one archive article
// @Tags         article_default
// @Produce      json
// @Param 		 uuid    path    int true  "uuid  archive article must be 36 characters"
// @Success      200     {object} response.Response{data=model.ArticleArchive}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `uuid` - article_uuid must be 36 characters long
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_articles` - article not found
// @Router      /api/v1/article/archive/text/{uuid} [get]
func (h *HandlerArticle) GetArchiveArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		uuid := request.PathValue("uuid")
		values.DataLog.MapLog["article_uuid"] = uuid
		if len(uuid) != 36 {
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
