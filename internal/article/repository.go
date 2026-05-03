package article

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/custom_errors"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"app/news-parser/pkg/generate_random"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RepositoryArticle struct {
	*open_Db.PostgresDb
	*open_Db.RedisDb
}

const (
	fieldHeader = "header"
	fieldUrl    = "url"
	fieldText   = "text"

	linkList = "LinkList"
)

func NewRepositoryArticle(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryArticle {
	return &RepositoryArticle{
		PostgresDb: postgres,
		RedisDb:    redis,
	}
}
func (r *RepositoryArticle) GetArticlesInCategoryToday(category string, limit int) ([]ResponseCategoryToday, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
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
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
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
func (r *RepositoryArticle) CreateArchiveArticle(article *model.ArticleArchive) error {
	res := r.PostgresDb.Create(&article)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticle) isRedisArticleExist() bool {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	keys, errKey := r.RedisDb.Keys(rdbContext, "*").Result()
	if errKey != nil {
		return false
	}
	for _, key := range keys {
		for _, category := range StorageCategories {
			if strings.Contains(key, category) {
				return true
			}
		}
	}
	return false
}
func (r *RepositoryArticle) allArticlesRedis() ([]model.ArticleArchive, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	keys, errKey := r.RedisDb.Keys(rdbContext, "*").Result()
	if errKey != nil {
		return nil, errKey
	}
	var sliceArticles []model.ArticleArchive
	for _, key := range keys {
		for _, category := range StorageCategories {
			if strings.Contains(key, category) {
				mapValue, errHGetAll := r.RedisDb.HGetAll(rdbContext, key).Result()
				if errHGetAll != nil {
					log.Println(errHGetAll)
				}
				sliceArticles = append(sliceArticles, model.ArticleArchive{
					Header:      mapValue[fieldHeader],
					URL:         mapValue[fieldUrl],
					Text:        mapValue[fieldText],
					Category:    category,
					Date:        common.DateNow(),
					UUIDArticle: uuid.New().String(),
				})
			}
		}
	}
	if len(sliceArticles) == 0 {
		return nil, errors.New("failed to get all the articles")
	}
	return sliceArticles, nil
}
func (r *RepositoryArticle) loadLinkList() ([]string, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	list, errGetList := r.RedisDb.LRange(rdbContext, linkList, 0, -1).Result()
	if errGetList != nil || len(list) == 0 {
		return nil, errors.New("failed to load LinkList")
	}
	return list, nil
}
func (r *RepositoryArticle) createNewArticle(art *ArticlesGoroutines) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	_, errTrans := r.Client.TxPipelined(rdbContext, func(pipeliner redis.Pipeliner) error {
		key := fmt.Sprint(art.Category, ":", generate_random.GenerateNumbers(7))
		errHSet := r.RedisDb.HSet(rdbContext, key, "header", art.Header, "url", art.Url, "text", art.Text, "is_article", art.IsArticle).Err()
		if errHSet != nil {
			return errHSet
		}
		errExp := r.RedisDb.Expire(rdbContext, key, 24*time.Hour).Err()
		if errExp != nil {
			return errExp
		}
		return nil
	})
	if errTrans != nil {
		fmt.Println(errTrans)
	}
}
