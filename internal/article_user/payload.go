package article_user

import (
	"app/news-parser/internal/model"
	"time"
)

type RequestCreateArticle struct {
	URL      string `validate:"required,url"`
	Category string `validate:"required,min=2,max=20"`
}
type RequestUpdateArticle struct {
	Category string `validate:"omitempty,min=2,max=20"`
}
type RequestUpdateBatchArticles struct {
	Category string `validate:"omitempty,min=2,max=20"`
	Domain   string `validate:"required"`
}
type ResponseSliceUserArticles struct {
	SliceUserArticles []model.UserArticle `json:"user-articles"`
}
type ResponseUserArticle struct {
	Article          model.UserArticle
	SuccessOperation SuccessOfTheOperation `json:"success-of-the-operation"`
}
type ResponseCreateUserArticle struct {
	Header           string
	URL              string
	Text             string
	Category         string
	ArticleUUID      string
	UserUUID         string
	SuccessOperation SuccessOfTheOperation `json:"success-of-the-operation"`
}
type SuccessOfTheOperation struct {
	Success bool
	Message string
	Status  int
}
type ResponseRemoveUserArticles struct {
	SliceRemoveUserArticles []RemoveUserArticle `json:"remove-user-articles"`
}
type RemoveUserArticle struct {
	*model.BaseModel
	ExpiredAt   time.Time
	Header      string
	URL         string
	Text        string
	Category    string
	ArticleUUID string
	UserUUID    string
}
