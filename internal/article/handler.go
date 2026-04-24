package article

import (
	"app/news-parser/internal/model"
	"app/news-parser/pkg/handler_response"
	"errors"
	"net/http"
	"strconv"
)

type HandlerArticle struct {
	model.Article
	ResponseArticle
	Dep *HandlerArticleDep
}
type HandlerArticleDep struct {
	*ServiceArticle
}

func NewHandlerArticle(router *http.ServeMux, dep *HandlerArticleDep) {
	article := &HandlerArticle{
		Dep: dep,
	}
	router.HandleFunc("GET /article/{category}", article.GetAllArticlesInCategory())
	router.HandleFunc("GET /article/{id}", article.GetArticle())
	router.HandleFunc("GET /article/popular", article.PopularCategories())
	router.HandleFunc("GET /article/archive", article.GetArchive())

}
func (h *HandlerArticle) GetAllArticlesInCategory() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		category := request.PathValue("category")
		if category != food && category != politics && category != sport && category != cloth && category != business && category != electronics {
			h.ResponseArticle.Error = ErrCategory.Error()
			handler_response.HandlerResponse(writer, h.ResponseArticle, http.StatusBadRequest)
			return
		}
		allArticle, errGetAllArticle := h.Dep.repo.GetAllArticlesInCategory(category)
		if errGetAllArticle != nil {
			h.ResponseArticle.Error = ErrCategory.Error()
			handler_response.HandlerResponse(writer, h.ResponseArticle, http.StatusInternalServerError)
			return
		}
		handler_response.HandlerResponse(writer, allArticle, http.StatusOK)
	}
}
func (h *HandlerArticle) GetArticle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
		id, errParseId := strconv.Atoi(idStr)
		if errParseId != nil || len(idStr) != lengthIdArticle {
			h.Article.Error = \\
			handler_response.HandlerResponse(writer, h.Article, http.StatusBadRequest)
			return
		}
	article, errGetArticle := h.Dep.repo.GetArticle(id)
	if errGetArticle != nil{
		h.Article.Error = errGetArticle.Error()
		if errors.Is(errGetArticle, ErrLoadArticles){
			handler_response.HandlerResponse(writer, h.Article, http.StatusInternalServerError)
		} else if errors.Is(errGetArticle, ){
			handler_response.HandlerResponse(writer, h.Article, http.StatusNotFound)
		}
		return
	}
		handler_response.HandlerResponse(writer,article, http.StatusOK)
	}
}
func (h *HandlerArticle) PopularCategories() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
func (h *HandlerArticle) GetArchive() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
