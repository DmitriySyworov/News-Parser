package stat

import (
	"app/news-parser/internal/common"
	"app/news-parser/internal/model"
	"app/news-parser/internal/open_Db"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RepositoryStat struct {
	*open_Db.RedisDb
	*open_Db.PostgresDb
}

func NewRepositoryStat(postgres *open_Db.PostgresDb, redis *open_Db.RedisDb) *RepositoryStat {
	return &RepositoryStat{
		RedisDb:    redis,
		PostgresDb: postgres,
	}
}
func (r *RepositoryStat) GetStatCategoryAllTime() (*ResponseStatCategoryAll, error) {
	var dbStat []CategoryDbAll
	res := r.DB.Model(&model.CategoryStat{}).Raw(`
SELECT category, sum(click) as sum_click, row_number() over (Order by category) AS place
FROM category_stats
GROUP BY category
ORDER BY sum_click
`).Scan(&dbStat)
	if res.Error != nil || len(dbStat) == 0 {
		return nil, ErrStatLoad
	}
	return &ResponseStatCategoryAll{
		Categories: dbStat,
	}, nil
}

func (r *RepositoryStat) GetStatCategoryByDate(date time.Time) (*ResponseStatCategoryDate, error) {
	var dbStat []CategoryDbDate
	res := r.DB.Raw(`
	SELECT ROW_NUMBER() OVER (ORDER BY category) as place,
	*FROM category_stats
	WHERE date = ?
	ORDER BY click
`, date).
		Scan(&dbStat)
	if res.Error != nil || len(dbStat) == 0 {
		return nil, ErrStatLoad
	}
	return &ResponseStatCategoryDate{
		Categories: dbStat,
	}, nil
}
func (r *RepositoryStat) GetStatArticleByDate(date time.Time) (*ResponseStatArticleDate, error) {
	var dbStat []ArticleDbDate
	res := r.DB.Raw(`
	SELECT ROW_NUMBER() OVER (ORDER BY url) as place,
	*FROM article_stats
	WHERE date = ?
	ORDER BY click
`, date).
		Scan(&dbStat)
	if res.Error != nil || len(dbStat) == 0 {
		return nil, ErrStatLoad
	}
	return &ResponseStatArticleDate{
		Articles: dbStat,
	}, nil
}
func (r *RepositoryStat) GetStatArticleAllTime() (*ResponseStatArticleAll, error) {
	var dbStat []ArticleDbAll
	res := r.DB.Model(&model.CategoryStat{}).Raw(`
SELECT url, sum(click) as sum_click, row_number() over (Order by url) AS place
FROM article_stats
GROUP BY url
ORDER BY sum_click
`).Scan(&dbStat)
	if res.Error != nil || len(dbStat) == 0 {
		return nil, ErrStatLoad
	}
	return &ResponseStatArticleAll{
		Articles: dbStat,
	}, nil
}
func (r *RepositoryStat) CreateStatCategory(statCategory *model.CategoryStat) error {
	res := r.PostgresDb.Create(statCategory)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryStat) CreateStatArticle(statArticle *model.ArticleStat) error {
	res := r.PostgresDb.Create(statArticle)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

const (
	RdbKeyStatCategory = "stat:categories"
	RdbKeyStatArticle  = "stat:articles"
	TimeDeleted        = 24*time.Hour + 10*time.Minute
)

func (r *RepositoryStat) GetStatCategoryToday() ([]redis.Z, error) {
	rdbCtx, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	statCategories, errZRev := r.RedisDb.ZRevRangeWithScores(rdbCtx, RdbKeyStatCategory, 0, -1).Result()
	if errZRev != nil {
		return nil, errZRev
	}
	return statCategories, nil
}
func (r *RepositoryStat) GetStatArticleToday() ([]redis.Z, error) {
	rdbCtx, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	statCategories, errZRev := r.RedisDb.ZRevRangeWithScores(rdbCtx, RdbKeyStatArticle, 0, -1).Result()
	if errZRev != nil {
		return nil, errZRev
	}
	return statCategories, nil
}

func (r *RepositoryStat) addClickCategory(category string) {
	rdbCTX, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	errZIncr := r.RedisDb.ZIncrBy(rdbCTX, RdbKeyStatCategory, 1, category).Err()
	if errZIncr != nil {
		log.Println(errZIncr)
	}
	errExpire := r.RedisDb.Expire(rdbCTX, RdbKeyStatCategory, TimeDeleted).Err()
	if errExpire != nil {
		log.Println(errExpire)
	}
}
func (r *RepositoryStat) addClickArticle(url string) {
	rdbCTX, cancel := context.WithTimeout(context.Background(), common.RdbTimeout)
	defer cancel()
	_, errTrans := r.Client.TxPipelined(rdbCTX, func(pipeliner redis.Pipeliner) error {
		errZIncr := r.RedisDb.ZIncrBy(rdbCTX, RdbKeyStatArticle, 1, url).Err()
		if errZIncr != nil {
			return errZIncr
		}
		errExpire := r.RedisDb.Expire(rdbCTX, RdbKeyStatArticle, TimeDeleted).Err()
		if errExpire != nil {
			return errExpire
		}
		return nil
	})
	if errTrans != nil {
		log.Println(errTrans)
		return
	}
}
