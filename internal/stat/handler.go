package stat

import (
	"app/news-parser/internal/custom_errors"
	"app/news-parser/pkg/handler_response"
	"net/http"
)

type HandlerStat struct {
	ResponseStatCategoryDate
	ResponseStatCategoryAll
	ResponseStatArticleDate
	ResponseStatArticleAll
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
	//router.HandleFunc("GET /stat/user") //auth!!!
}
func (h *HandlerStat) GetStatCategoryByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.ResponseStatCategoryDate.Error = custom_errors.ErrIncorrectDate.Error()
			handler_response.HandlerResponse(writer, h.ResponseStatCategoryDate, http.StatusBadRequest)
			return
		}
		respStatDateCategories, errGetStat := h.ServiceStat.GetStatCategoryByDate(dateStr)
		if errGetStat != nil {
			h.ResponseStatCategoryDate.Error = errGetStat.Error()
			switch errGetStat {
			case custom_errors.ErrIncorrectDate:
				handler_response.HandlerResponse(writer, h.ResponseStatCategoryDate, http.StatusBadRequest)
			case ErrStatLoad:
				handler_response.HandlerResponse(writer, h.ResponseStatCategoryDate, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, respStatDateCategories, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatCategoryAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		respStatCategoryAll, errGetAllStat := h.ServiceStat.Repo.GetStatCategoryAllTime()
		if errGetAllStat != nil {
			h.ResponseStatCategoryAll.Error = errGetAllStat.Error()
			handler_response.HandlerResponse(writer, h.ResponseStatCategoryAll, http.StatusInternalServerError)
			return
		}
		handler_response.HandlerResponse(writer, respStatCategoryAll, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatArticleByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		dateStr := request.URL.Query().Get("date")
		if dateStr == "" {
			h.ResponseStatArticleDate.Error = custom_errors.ErrIncorrectDate.Error()
			handler_response.HandlerResponse(writer, h.ResponseStatArticleDate, http.StatusBadRequest)
			return
		}
		respStatDateArticles, errGetStat := h.ServiceStat.GetStatArticleByDate(dateStr)
		if errGetStat != nil {
			h.ResponseStatArticleDate.Error = errGetStat.Error()
			switch errGetStat {
			case custom_errors.ErrIncorrectDate:
				handler_response.HandlerResponse(writer, h.ResponseStatArticleDate, http.StatusBadRequest)
			case ErrStatLoad:
				handler_response.HandlerResponse(writer, h.ResponseStatArticleDate, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, respStatDateArticles, http.StatusOK)
	}
}
func (h *HandlerStat) GetStatArticleAllTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		respStatArticleAll, errGetAllStat := h.ServiceStat.Repo.GetStatArticleAllTime()
		if errGetAllStat != nil {
			h.ResponseStatArticleAll.Error = errGetAllStat.Error()
			handler_response.HandlerResponse(writer, h.ResponseStatArticleAll, http.StatusInternalServerError)
			return
		}
		handler_response.HandlerResponse(writer, respStatArticleAll, http.StatusOK)
	}
}
