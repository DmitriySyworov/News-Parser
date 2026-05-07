package article

import "app/news-parser/internal/model"

type ResponseCategoryToday struct {
	Header    string
	URL       string
	IsArticle bool
	IDArticle uint
	Error     string
}
type ResponseCategoryArchive struct {
	Header      string
	URL         string
	UUIDArticle string
	Error       string
}
type RequestCreateArticle struct {
	URL      string `validate:"required,url"`
	Category string `validate:"required"`
}
type ResponseCreateArticle struct {
	UserArticles []model.UserArticle `json:"user-articles"`
	Error        string
}
