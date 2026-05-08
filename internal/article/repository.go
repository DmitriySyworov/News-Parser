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
func (r *RepositoryArticle) GetArticlesInCategoryToday(category, filter string, limit int) ([]ResponseCategoryToday, error) {
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
			dataArticle, errHMGet := r.RedisDb.Client.HMGet(rdbContext, key, fieldHeader, fieldUrl, fieldIsArticle).Result()
			if errHMGet != nil {
				return nil, ErrLoadArticles
			}
			if dataArticle[2] == "1" && filter == "true" {
				header, okHeader := dataArticle[0].(string)
				url, okUrl := dataArticle[1].(string)
				if !okHeader || !okUrl {
					return nil, ErrLoadArticles
				}
				id, errParseId := strconv.Atoi(strings.Split(key, ":")[1])
				if errParseId != nil {
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
			} else if dataArticle[2] == "0" && filter == "false" {
				header, okHeader := dataArticle[0].(string)
				url, okUrl := dataArticle[1].(string)
				if !okHeader || !okUrl {
					return nil, ErrLoadArticles
				}
				id, errParseId := strconv.Atoi(strings.Split(key, ":")[1])
				if errParseId != nil {
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
func (r *RepositoryArticle) GetArticlesInCategoryArchive(category string, offset, limit int, date time.Time) ([]model.ArticleArchive, error) {
	var archiveArticles []model.ArticleArchive
	res := r.Where("category = ? AND date = ?", category, date).
		Offset(offset).
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
func (r *RepositoryArticle) GetUserArticle(idArticle uint) (*model.UserArticle, error) {
	var userArticle model.UserArticle
	res := r.PostgresDb.Where("id_article = ?", idArticle).First(&userArticle)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userArticle, nil
}
func (r *RepositoryArticle) DeleteAllUserArticle(uuidUser string) error {
	res := r.PostgresDb.Where("uuid_user = ?", uuidUser).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryArticle) GetAllUserArticlesWithoutText(category string, offset, limit int) (*ResponseUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, category, id_article, uuid_user FROM user_articles
					OFFSET ?
					LIMIT ?
`, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, category, id_article, uuid_user FROM user_articles
                    WHERE  category = ?
					OFFSET ?
					LIMIT ?
`, category, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	if len(sliceUserArticle) == 0 {
		return nil, custom_errors.ErrUserNotFound
	}
	return &ResponseUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticle) GetAllUserArticlesWithText(category string, offset, limit int) (*ResponseUserArticles, error) {
	var sliceUserArticle []model.UserArticle
	if category == "" {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, text, category, id_article, uuid_user FROM user_articles
					OFFSET ?
					LIMIT ?
`, category, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := r.PostgresDb.
			Raw(`SELECT id, created_at, updated_at, deleted_at, header, url, text, category, id_article, uuid_user FROM user_articles
                    WHERE  category = ?
					OFFSET ?
					LIMIT ?
`, category, offset, limit).
			Find(&sliceUserArticle)
		if res.Error != nil {
			return nil, res.Error
		}
	}
	if len(sliceUserArticle) == 0 {
		return nil, custom_errors.ErrUserNotFound
	}
	return &ResponseUserArticles{
		SliceUserArticles: sliceUserArticle,
	}, nil
}
func (r *RepositoryArticle) DeleteUserArticleByID(idArticle uint) error {
	res := r.PostgresDb.Where("id_article = ?", idArticle).Delete(&model.UserArticle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
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
func (r *RepositoryArticle) createUserNewArticle(art *model.UserArticle) error {
	res := r.PostgresDb.Create(&art)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
type idArticlesSlice struct{
	IdArticle uint
	DeletedAt time.Time
}
func (r *RepositoryArticle) deleteUserArticles() []model.UserArticle {
	var sliceIdArticles []model.UserArticle
	res := r.PostgresDb.Raw(`SELECT id_article, deleted_at FROM user_articles
							  WHERE deleted_at IS NOT NULL`).
		Scan(&sliceIdArticles)
	if res.Error != nil {
		log.Println(res.Error)
	}
	for _, articles := range sliceIdArticles{
		if time.Now().Compare() < articles.DeletedAt
		r.PostgresDb.
	}
}
