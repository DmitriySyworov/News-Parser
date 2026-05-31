package stat

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/loggers"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"log/slog"
	"net/http"
)

type HandlerStat struct {
	response.Response
	*ServiceStat
	Dep *HandlerStatDep
}
type HandlerStatDep struct {
	*loggers.Logger
	*middleware.ManagerMiddleware
}

func NewHandlerStat(router *http.ServeMux, service *ServiceStat, dep *HandlerStatDep) {
	stat := &HandlerStat{
		ServiceStat: service,
		Dep:         dep,
	}
	router.HandleFunc("GET /api/v1/stat/article/category", stat.GetStatCategoryByDate())
	router.HandleFunc("GET /api/v1/stat/article/category/alltime", stat.GetStatCategoryAllTime())
	router.HandleFunc("GET /api/v1/stat/article", stat.GetStatArticleByDate())
	router.HandleFunc("GET /api/v1/stat/article/alltime", stat.GetStatArticleAllTime())
	router.Handle("GET /api/v1/my/stat/article", dep.IsAuthJWT(stat.GetStatUserArticleByDate()))
	router.Handle("GET /api/v1/my/stat/article/alltime", dep.IsAuthJWT(stat.GetStatUserArticleAllTime()))
}

// GetStatCategoryByDate godoc
// @Summary      get stat categories  by date
// @Description	 get stat popular categories  by date
// @Tags         stat
// @Produce      json
// @Param 		 date   query  string true  "date must be YYYY-MM-D'"
// @Success      200     {object} response.Response{data=ResponseStatCategoryDate}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `date` - the date format must be YYYY-MM-DD
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty
// @Router       /api/v1/stat/article/category [get]
func (h *HandlerStat) GetStatCategoryByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		dateStr := request.URL.Query().Get("date")
		values.DataLog.MapLog["date"] = dateStr
		if dateStr == "" {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		respStatDateCategories, errGetStat := h.ServiceStat.GetStatCategoryByDate(dateStr)
		if errGetStat != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetStat)
			h.Response.Errors = append(h.Response.Errors, *errGetStat)
			switch errGetStat.Message {
			case custom_errors.ErrIncorrectDate.Error():
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case ErrStatNotFound.Error():
				response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatDateCategories
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

// GetStatCategoryAllTime godoc
// @Summary      get stat categories all time
// @Description	 get stat popular categories all-time
// @Tags         stat
// @Produce      json
// @Success      200     {object} response.Response{data=ResponseStatCategoryAll}
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty
// @Router      /api/v1/stat/article/category/alltime [get]
func (h *HandlerStat) GetStatCategoryAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		respStatCategoryAll, errGetAllStat := h.ServiceStat.Repo.GetStatCategoryAllTime()
		if errGetAllStat != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetAllStat)
			h.Response.Errors = append(h.Response.Errors, *errGetAllStat)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatCategoryAll
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

// GetStatArticleByDate godoc
// @Summary      get stat articles  by date
// @Description	 get stat popular articles  by date
// @Tags         stat
// @Produce      json
// @Param 		 date   query  string  true "date must be YYYY-MM-DD"
// @Success      200     {object} response.Response{data=ResponseStatArticleDate}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `date` - the date format must be YYYY-MM-DD
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty
// @Router       /api/v1/stat/article [get]
func (h *HandlerStat) GetStatArticleByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		dateStr := request.URL.Query().Get("date")
		values.DataLog.MapLog["date"] = dateStr
		if dateStr == "" {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: custom_errors.ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		respStatDateArticles, errGetStat := h.ServiceStat.GetStatArticleByDate(dateStr)
		if errGetStat != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetStat)
			h.Response.Errors = append(h.Response.Errors, *errGetStat)
			switch errGetStat.Message {
			case custom_errors.ErrIncorrectDate.Error():
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			case ErrStatNotFound.Error():
				response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatDateArticles
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

// GetStatArticleAllTime godoc
// @Summary      get stat articles all time
// @Description	 get stat popular articles  all-time
// @Tags         stat
// @Produce      json
// @Success      200     {object} response.Response{data=ResponseStatArticleAll}
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty
// @Router     /api/v1/stat/article/alltime [get]
func (h *HandlerStat) GetStatArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		respStatArticleAll, errGetAllStat := h.ServiceStat.Repo.GetStatArticleAllTime()
		if errGetAllStat != nil {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetAllStat)
			h.Response.Errors = append(h.Response.Errors, *errGetAllStat)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatArticleAll
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

// GetStatUserArticleByDate godoc
// @Summary      get stat user by date
// @Description	 get stat actions user by date
// @Tags         stat
// @Produce      json
// @Param 		 date   query string true "date must be YYYY-MM-DD"
// @Param 		 Authorization  header string true "JWT storing user_uuid"
// @Success      200     {object} response.Response{data=ResponseUserArticleStat}
// @Failure      400      {object}  response.Response[string] "validation errors:" | `date` - the date format must be YYYY-MM-DD
// @Failure		 401       {object}  response.Response[string] "unauthorized error:" | `token` - incorrect token
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty | `not_found_user` - such user does not exist
// @Router     /api/v1/my/stat/article [get]
func (h *HandlerStat) GetStatUserArticleByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		date := request.URL.Query().Get("date")
		values.DataLog.MapLog["date"] = date
		respUserStat, sliceErrGetStat := h.ServiceStat.GetUserArticleStatByDate(values.UserUUID, date)
		h.Response.Errors = append(h.Response.Errors, sliceErrGetStat...)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, sliceErrGetStat...)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respUserStat
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}

// GetStatUserArticleAllTime godoc
// @Summary      get stat user all time
// @Description	 get stat actions user all-time
// @Tags         stat
// @Produce      json
// @Param 		 Authorization  header  string true  "JWT storing user_uuid"
// @Success      200     {object} response.Response{data=ResponseUserArticleAllTimeStat}
// @Failure		 401       {object}  response.Response[string] "unauthorized error:" | `token` - incorrect token
// @Failure      404      {object}  response.Response[string] "not found:" | `not_found_statistic` - statistic not found or empty | `not_found_user` - such user does not exist
// @Router     /api/v1/my/stat/article/alltime [get]
func (h *HandlerStat) GetStatUserArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response{}
		}()
		ctxValues := request.Context().Value(middleware.KeyContextValues)
		values, ok := ctxValues.(*middleware.ContextValues)
		if !ok {
			h.Dep.Logger.SystemLogger(slog.LevelError, custom_errors.ErrFailedTypeContextValues.Error()+request.Pattern)
		}
		values.DataLog.UserUUID = values.UserUUID
		if len(values.UserUUID) != 36 {
			err := response.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			}
			values.DataLog.Errors = append(values.DataLog.Errors, err)
			h.Response.Errors = append(h.Response.Errors, err)
			response.HandlerResponse(writer, h.Response, http.StatusUnauthorized)
			return
		}
		respStat, errGetAllTimeStat := h.ServiceStat.GetUserArticleAllTimeStat(values.UserUUID)
		h.Response.Errors = append(h.Response.Errors, *errGetAllTimeStat)
		if len(h.Response.Errors) != 0 {
			values.DataLog.Errors = append(values.DataLog.Errors, *errGetAllTimeStat)
			if len(h.Response.Errors) == 1 {
				response.HandlerResponse(writer, h.Response, h.Response.Errors[0].Status)
			} else {
				response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			}
			return
		}
		h.Response.Success = true
		h.Response.Data = respStat
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
