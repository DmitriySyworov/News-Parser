package article

import (
	"app/news-parser/internal/model"
	"app/news-parser/pkg/open_Db"
	"context"
	"fmt"
	"net/http"
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
	fieldText = "text"


)

func NewRepositoryArticle(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryArticle {
	return &RepositoryArticle{
		PostgresDb: postgres,
		RedisDb:    redis,
	}
}
func (r *RepositoryArticle) GetAllArticlesInCategory(category string) ([]ResponseArticle, error) {
	rdbContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keys, errKeys := r.RedisDb.Client.Keys(rdbContext, "*").Result()
	if errKeys != nil {
		return nil, ErrLoadArticles
	}
	var sliceArticles []ResponseArticle
	for _, key := range keys {
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
			sliceArticles = append(sliceArticles, ResponseArticle{
				Header:    header,
				URL:       url,
				IdArticle: uint(id),
			})
		}
	}
	if len(sliceArticles) == 0 {
		return nil, ErrLoadArticles
	}
	return sliceArticles, nil
}
func (r *RepositoryArticle) GetArticle(id int) (*model.Article, error) {
	idStr := fmt.Sprint(id)
	rdbContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	keys, errKey :=r.RedisDb.Client.Keys(rdbContext, "*").Result()
	if errKey != nil {
		return nil, ErrLoadArticles
	}
	for _, key := range keys{
		if strings.Contains(key, idStr){
			mapValue, errHGetAll := r.RedisDb.Client.HGetAll(rdbContext, key).Result()
			if errHGetAll != nil {
				return nil,
			}
			return &model.Article{
			Header: mapValue[fieldHeader],
			URL: mapValue[fieldUrl],
			Text: mapValue[fieldText],
			IdArticle: uint(id),
			}, nil
		}
	}
}