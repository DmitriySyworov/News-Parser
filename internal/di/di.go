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
