package article

import (
	"app/news-parser/internal/model"
	"app/news-parser/pkg/custom_errors"
	"app/news-parser/pkg/open_Db"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RepositoryArticle struct {
	*open_Db.PostgresDb
	*open_Db.RedisDb
}

const (
	fieldHeader = "header"
	fieldUrl    = "url"
	fieldText   = "text"
)

func NewRepositoryArticle(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryArticle {
	return &RepositoryArticle{
		PostgresDb: postgres,
		RedisDb:    redis,
	}
}
func (r *RepositoryArticle) GetArticlesInCategoryToday(category string, limit int) ([]ResponseCategoryToday, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keys, errKeys := r.RedisDb.Client.Keys(rdbContext, "*").Result()
	if errKeys != nil {
		return nil, ErrLoadArticles
	}
	var sliceArticles []ResponseCategoryToday
	for _, key := range keys {
		if len(sliceArticles) >= limit {
			return sliceArticles, nil
		}
		if strings.Contains(key, category) {
			dataArticle, errHMGet := r.RedisDb.Client.HMGet(rdbContext, key, fieldHeader, fieldUrl).Result()
			if errHMGet != nil {
				return nil, ErrLoadArticles
			}
			header, okHeader := dataArticle[0].(string)
			url, okUrl := dataArticle[1].(string)
			if !okHeader || !okUrl {
				return nil, ErrLoadArticles
			}
			id, errParseId := strconv.Atoi(strings.Split(key, ":")[1])
			if errParseId != nil {
				return nil, ErrLoadArticles
			}
			sliceArticles = append(sliceArticles, ResponseCategoryToday{
				Header:    header,
				URL:       url,
				IDArticle: uint(id),
			})
		}
	}
	if len(sliceArticles) == 0 {
		return nil, ErrLoadArticles
	}
	return sliceArticles, nil
}
func (r *RepositoryArticle) GetArticleToday(id int) (*model.ArticleToday, error) {
	idStr := fmt.Sprint(id)
	rdbContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keys, errKey := r.RedisDb.Client.Keys(rdbContext, "*").Result()
	if errKey != nil {
		return nil, ErrLoadArticles
	}
	for _, key := range keys {
		if strings.Contains(key, idStr) {
			mapValue, errHGetAll := r.RedisDb.Client.HGetAll(rdbContext, key).Result()
			if errHGetAll != nil {
				return nil, custom_errors.ErrRecordNotFound
			}
			return &model.ArticleToday{
				Header:    mapValue[fieldHeader],
				URL:       mapValue[fieldUrl],
				Text:      mapValue[fieldText],
				IDArticle: uint(id),
			}, nil
		}
	}
	return nil, custom_errors.ErrRecordNotFound
}
func (r *RepositoryArticle) GetArticlesInCategoryArchive(category string, limit int, date time.Time) ([]model.ArticleArchive, error) {
	var archiveArticles []model.ArticleArchive
	res := r.Where("category = ? AND date = ?", category, date).
		Limit(limit).
		Find(&archiveArticles)
	if res.Error != nil || len(archiveArticles) == 0 {
		return nil, ErrLoadArticles
	}
	return archiveArticles, nil
}
func (r *RepositoryArticle) GetArchiveArticle(uuid string) (*model.ArticleArchive, error) {
	var articleArch model.ArticleArchive
	res := r.Where("uuid_article = ?", uuid).First(&articleArch)
	if res.Error != nil {
		return nil, res.Error
	}
	return &articleArch, nil
}
func (r *RepositoryArticle)PopularCategories(){

}