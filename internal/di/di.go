package di

import (
	"app/news-parser/internal/model"

	"github.com/redis/go-redis/v9"
)

type IRepoStat interface {
	GetStatCategoryToday() ([]redis.Z, error)
	GetStatArticleToday() ([]redis.Z, error)
	CreateStatCategory(*model.CategoryStat) error
	CreateStatArticle(statArticle *model.ArticleStat) error
}
type IRepoUser interface {
	IsUserExist(name, email string) error
	CreateUser(*model.User) error
	GetUserByEmail(string) (*model.User, error)
	GetUserByUUID(string) (*model.User, error)
}
