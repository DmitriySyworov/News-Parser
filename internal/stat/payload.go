package stat

import (
	"app/news-parser/internal/model"
	"time"
)

type ResponseStatCategoryDate struct {
	Categories []CategoryDbDate
}
type ResponseStatCategoryAll struct {
	Categories []CategoryDbAll
}
type CategoryDbDate struct {
	Category string
	Click    uint
	Date     time.Time
	Place    uint
}
type CategoryDbAll struct {
	Category string
	SumClick uint
	Place    uint
}
type ResponseStatArticleDate struct {
	Articles []ArticleDbDate
}
type ResponseStatArticleAll struct {
	Articles []ArticleDbAll
}
type ArticleDbDate struct {
	URL   string
	Click uint
	Date  time.Time
	Place uint
}
type ResponseUserArticleStat struct {
	NowExist               int
	*model.UserArticleStat `json:"user-article-statistic"`
}
type ResponseUserArticleAllTimeStat struct {
	NowExist               int
	*UserArticleAllTimeStat `json:"user-article-statistic-all-time"`
}
type UserArticleAllTimeStat struct {
	AllTimeCreated int
	AllTimeUpdated int
	AllTimeSoftDeleted int
	AllTimeHardDeleted int
	AllTimeRecovered int
}
type ArticleDbAll struct {
	URL      string
	SumClick uint
	Place    uint
}
