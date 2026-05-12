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
}
type RequestUpdateBatchArticles struct {
	Domain string `validate:"required"`
}
type ResponseSliceUserArticles struct {
	SliceUserArticles []model.UserArticle `json:"user-articles"`
	Error             string
}
type ResponseUserArticle struct {
	Article model.UserArticle
	Status  int
	Error   string
}
type ResponseRemoveUserArticles struct {
	SliceRemoveUserArticles []RemoveUserArticle `json:"remove-user-articles"`
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
