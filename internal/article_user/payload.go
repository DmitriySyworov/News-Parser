package article_user

import (
	"app/news-parser/internal/model"
	"time"

	"gorm.io/gorm"
)

type RequestCreateArticle struct {
	URL      string `validate:"required,url"`
	Category string `validate:"required"`
}
type RequestUpdateArticle struct {
	Category string
	Domain   string
}
type ResponseUserArticles struct {
	SliceUserArticles []model.UserArticle `json:"user-articles"`
	Error             string
}
type ResponseUserDelete struct {
	Error string
}
type ResponseRemoveUserArticles struct {
	SliceRemoveUserArticles []RemoveUserArticle `json:"remove-user-articles"`
	Error                   string
}
type RemoveUserArticle struct {
	*gorm.Model
	ExpiredAt time.Time
	Header    string
	URL       string
	Text      string
	Category  string
	IDArticle uint
	UUIDUser  string
}
