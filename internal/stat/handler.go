package stat

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
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
}

func NewHandlerStat(router *http.ServeMux, dep *HandlerStatDep) {
	stat := &HandlerStat{
		HandlerStatDep: dep,
	}
	router.HandleFunc("GET /stat/article/category", stat.GetStatCategoryByDate())
	router.HandleFunc("GET /stat/article/category/alltime", stat.GetStatCategoryAllTime())
	router.HandleFunc("GET /stat/article", stat.GetStatArticleByDate())
	router.HandleFunc("GET /stat/article/alltime", stat.GetStatArticleAllTime())
	//router.HandleFunc("GET /my/stat/article/create")
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
			case ErrStatLoad.Error():
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
			case ErrStatLoad.Error():
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
