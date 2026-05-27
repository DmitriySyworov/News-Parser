package stat

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/middleware"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerStat struct {
	custom_errors.ResponseError
	common.ResponseSuccessful
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
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		respStatDateCategories, errGetStat := h.ServiceStat.GetStatCategoryByDate(dateStr)
		if errGetStat != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetStat)
			switch errGetStat.Message {
			case ErrIncorrectDate.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case ErrStatNotFound.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respStatDateCategories
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatCategoryAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		respStatCategoryAll, errGetAllStat := h.ServiceStat.Repo.GetStatCategoryAllTime()
		if errGetAllStat != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetAllStat)
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respStatCategoryAll
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatArticleByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: ErrIncorrectDate.Error(),
				Status:  http.StatusBadRequest,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			return
		}
		respStatDateArticles, errGetStat := h.ServiceStat.GetStatArticleByDate(dateStr)
		if errGetStat != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetStat)
			switch errGetStat.Message {
			case ErrIncorrectDate.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			case ErrStatNotFound.Error():
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respStatDateArticles
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		respStatArticleAll, errGetAllStat := h.ServiceStat.Repo.GetStatArticleAllTime()
		if errGetAllStat != nil {
			h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetAllStat)
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusNotFound)
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respStatArticleAll
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatUserArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		date := request.URL.Query().Get("date")
		respUserStat, sliceErrGetStat := h.ServiceStat.GetUserArticleStat(ctxTokens.UUID, date)
		h.ResponseError.Errors = append(h.ResponseError.Errors, sliceErrGetStat...)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respUserStat
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatUserArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			h.ResponseError = custom_errors.ResponseError{}
			h.ResponseSuccessful = common.ResponseSuccessful{}
		}()
		ctxValue := request.Context().Value(middleware.KeyContext)
		ctxTokens, ok := ctxValue.(middleware.ContextToken)
		if !ok {
			h.ResponseError.Errors = append(h.ResponseError.Errors, custom_errors.Error{
				Message: custom_errors.ErrIncorrectToken.Error(),
				Status:  http.StatusUnauthorized,
			})
			handler_response.HandlerResponse(writer, h.ResponseError, http.StatusUnauthorized)
			return
		}
		respStat, errGetAllTimeStat := h.ServiceStat.GetUserArticleAllTimeStat(ctxTokens.UUID)
		h.ResponseError.Errors = append(h.ResponseError.Errors, *errGetAllTimeStat)
		if len(h.ResponseError.Errors) != 0 {
			if len(h.ResponseError.Errors) == 1 {
				handler_response.HandlerResponse(writer, h.ResponseError, h.ResponseError.Errors[0].Status)
			} else {
				handler_response.HandlerResponse(writer, h.ResponseError, http.StatusBadRequest)
			}
			return
		}
		h.ResponseSuccessful.Success = true
		h.ResponseSuccessful.Data = respStat
		handler_response.HandlerResponse(writer, h.ResponseSuccessful, http.StatusOK)
	}
}
