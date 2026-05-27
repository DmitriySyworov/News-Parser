package stat

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/internal/response"
	"net/http"
)

type HandlerStat struct {
	response.Response[any]
	*HandlerStatDep
}
type HandlerStatDep struct {
	*ServiceStat
	*middleware.ManagerMiddleware
}

func NewHandlerStat(router *http.ServeMux, dep *HandlerStatDep) {
	stat := &HandlerStat{
		HandlerStatDep: dep,
	}
	router.HandleFunc("GET /stat/article/category", stat.GetStatCategoryByDate())
	router.HandleFunc("GET /stat/article/category/alltime", stat.GetStatCategoryAllTime())
	router.HandleFunc("GET /stat/article", stat.GetStatArticleByDate())
	router.HandleFunc("GET /stat/article/alltime", stat.GetStatArticleAllTime())
	router.Handle("GET /my/stat/article", dep.IsAuthJWT(stat.GetStatUserArticle()))
	router.Handle("GET /my/stat/article/alltime", dep.IsAuthJWT(stat.GetStatUserArticleAllTime()))
}
func (h *HandlerStat) GetStatCategoryByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		respStatDateCategories, errGetStat := h.ServiceStat.GetStatCategoryByDate(dateStr)
		if errGetStat != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetStat)
			switch errGetStat.Message {
			case ErrIncorrectDate.Error():
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
func (h *HandlerStat) GetStatCategoryAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		respStatCategoryAll, errGetAllStat := h.ServiceStat.Repo.GetStatCategoryAllTime()
		if errGetAllStat != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetAllStat)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatCategoryAll
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatArticleByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.Response.Errors = append(h.Response.Errors, response.Error{
				Message: ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			response.HandlerResponse(writer, h.Response, http.StatusBadRequest)
			return
		}
		respStatDateArticles, errGetStat := h.ServiceStat.GetStatArticleByDate(dateStr)
		if errGetStat != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetStat)
			switch errGetStat.Message {
			case ErrIncorrectDate.Error():
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
func (h *HandlerStat) GetStatArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.Response = response.Response[any]{}
		}()
		respStatArticleAll, errGetAllStat := h.ServiceStat.Repo.GetStatArticleAllTime()
		if errGetAllStat != nil {
			h.Response.Errors = append(h.Response.Errors, *errGetAllStat)
			response.HandlerResponse(writer, h.Response, http.StatusNotFound)
			return
		}
		h.Response.Success = true
		h.Response.Data = respStatArticleAll
		response.HandlerResponse(writer, h.Response, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatUserArticle() http.HandlerFunc {
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
		date := request.URL.Query().Get("date")
		respUserStat, sliceErrGetStat := h.ServiceStat.GetUserArticleStat(ctxTokens.UUID, date)
		h.Response.Errors = append(h.Response.Errors, sliceErrGetStat...)
		if len(h.Response.Errors) != 0 {
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
func (h *HandlerStat) GetStatUserArticleAllTime() http.HandlerFunc {
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
		respStat, errGetAllTimeStat := h.ServiceStat.GetUserArticleAllTimeStat(ctxTokens.UUID)
		h.Response.Errors = append(h.Response.Errors, *errGetAllTimeStat)
		if len(h.Response.Errors) != 0 {
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
