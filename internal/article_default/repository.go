package article_default

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
	fieldHeader    = "header"
	fieldUrl       = "url"
	fieldText      = "text"
	fieldIsArticle = "is_article"

	linkList = "LinkList"
)

func NewRepositoryArticle(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryArticle {
	return &RepositoryArticle{
		PostgresDb: postgres,
		RedisDb:    redis,
	}
}
func (r *RepositoryArticle) GetArticlesInCategoryToday(category string, offset, limit int, filter, withText bool) ([]ResponseCategoryToday, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	keyZ := "Z:" + category + ":" + fmt.Sprint(filter)
	keysArticle, errZRange := r.RedisDb.ZRange(rdbContext, keyZ, int64(offset), int64(offset+limit-1)).Result()
	if errZRange != nil {
		return nil, errZRange
	}
	var sliceArticles []ResponseCategoryToday
	for _, key := range keysArticle {
		id, errParseId := strconv.Atoi(strings.Split(key, ":")[1])
		if errParseId != nil {
			return nil, errParseId
		}
		if !withText {
			dataArticle, errHMGet := r.RedisDb.Client.HMGet(rdbContext, key, fieldHeader, fieldUrl, fieldIsArticle).Result()
			if errHMGet != nil {
				return nil, ErrLoadArticles
			}
			if dataArticle[2] == "1" && filter {
				header, okHeader := dataArticle[0].(string)
				url, okUrl := dataArticle[1].(string)
				if !okHeader || !okUrl {
					return nil, ErrLoadArticles
				}
				isArticle := false
				if dataArticle[2] == "1" {
					isArticle = true
				}
				sliceArticles = append(sliceArticles, ResponseCategoryToday{
					Header:    header,
					URL:       url,
					IDArticle: uint(id),
					IsArticle: isArticle,
				})
			} else if dataArticle[2] == "0" && !filter {
				header, okHeader := dataArticle[0].(string)
				url, okUrl := dataArticle[1].(string)
				if !okHeader || !okUrl {
					return nil, ErrLoadArticles
				}
				isArticle := false
				if dataArticle[2] == "1" {
					isArticle = true
				}
				sliceArticles = append(sliceArticles, ResponseCategoryToday{
					Header:    header,
					URL:       url,
					IDArticle: uint(id),
					IsArticle: isArticle,
				})
			}
		} else {
			mapValueH, errHGetAll := r.RedisDb.HGetAll(rdbContext, key).Result()
			var isArticle bool
			if errHGetAll != nil {
				return nil, errHGetAll
			}
			if mapValueH[fieldIsArticle] == "1" {
				isArticle = true
			}
			sliceArticles = append(sliceArticles, ResponseCategoryToday{
				Header:    mapValueH[fieldHeader],
				URL:       mapValueH[fieldUrl],
				Text:      mapValueH[fieldText],
				IsArticle: isArticle,
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
		return nil, errKey
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
	return nil, ErrLoadArticles
}
func (r *RepositoryArticle) GetArticlesInCategoryArchive(category string, offset, limit int, date time.Time) ([]model.ArticleArchive, error) {
	var archiveArticles []model.ArticleArchive
	res := r.Where("category = ? AND date = ?", category, date).
		Offset(offset).
		Limit(limit).
		Find(&archiveArticles)
	if res.Error != nil || len(archiveArticles) == 0 {
		return nil, ErrNotFoundArticle
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
		keyArticle := fmt.Sprint(art.Category, ":", generate_random.GenerateNumbers(7))
		errHSet := r.RedisDb.HSet(rdbContext, keyArticle, "header", art.Header, "url", art.Url, "text", art.Text, "is_article", art.IsArticle).Err()
		if errHSet != nil {
			return errHSet
		}
		errExp := r.RedisDb.Expire(rdbContext, keyArticle, common.Day).Err()
		if errExp != nil {
			return errExp
		}
		keyZ := "Z:" + art.Category + ":" + fmt.Sprint(art.IsArticle)
		ZCategory := r.RedisDb.ZRevRangeWithScores(rdbContext, keyZ, 0, -1).Val()
		if len(ZCategory) == 0 {
			r.ZAdd(rdbContext, keyZ, redis.Z{
				Member: keyArticle,
				Score:  1,
			})
		} else {
			r.ZAdd(rdbContext, keyZ, redis.Z{
				Member: keyArticle,
				Score:  float64(len(ZCategory) + 1),
			})
		}
		errExpire := r.Expire(rdbContext, keyZ, common.Day).Err()
		if errExpire != nil {
			return errExpire
		}
		return nil
	})
	if errTrans != nil {
		fmt.Println(errTrans)
	}
}
